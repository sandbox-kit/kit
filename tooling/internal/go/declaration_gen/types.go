package declarationgen

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/go/schema"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func (e *Emitter) Type(t spec.TypeRef) (string, error) {
	if bound, ok := e.TypeBindings[t.Ref]; ok {
		local := *e
		local.TypeBindings = map[string]spec.TypeRef{}
		for id, value := range e.TypeBindings {
			if id != t.Ref {
				local.TypeBindings[id] = value
			}
		}
		bound.Optional = bound.Optional || t.Optional
		bound.Nullable = bound.Nullable || t.Nullable
		bound.Reference = bound.Reference || t.Reference
		if t.Ownership != "" {
			bound.Ownership = t.Ownership
		}
		return local.Type(bound)
	}
	if t.Tuple != nil {
		return "", fmt.Errorf("Go has no first-class tuple type; use named object fields or callable results")
	}
	if len(t.Union) > 0 {
		return "", fmt.Errorf("Go profile cannot represent union types; declare a tagged schema alternative")
	}
	if t.Optional && t.Nullable {
		return "", fmt.Errorf("Go profile cannot distinguish optional from explicit null in the same slot")
	}
	var value string
	nilable := false
	switch {
	case t.Builtin != "":
		value = e.Profile.Builtins[t.Builtin]
		switch value {
		case "string", "bool", "any", "int", "int32", "uint32", "float32", "int64", "uint64", "float64", "[]byte":
		default:
			return "", fmt.Errorf("unsupported Go builtin %s", t.Builtin)
		}
		nilable = value == "any" || value == "[]byte"
	case t.Sequence != nil:
		element, err := e.Type(*t.Sequence)
		if err != nil {
			return "", err
		}
		value = "[]" + element
		nilable = true
	case t.Map != nil:
		if t.Map.Key.Optional || t.Map.Key.Nullable || t.Map.Key.Reference || (t.Map.Key.Builtin != "string" && t.Map.Key.Builtin != "integer" && t.Map.Key.Builtin != "unsigned" && t.Map.Key.Builtin != "boolean") {
			return "", fmt.Errorf("Go map key must be a supported comparable scalar")
		}
		key, err := e.Type(t.Map.Key)
		if err != nil {
			return "", err
		}
		element, err := e.Type(t.Map.Value)
		if err != nil {
			return "", err
		}
		value = "map[" + key + "]" + element
		nilable = true
	case t.Function != nil:
		c := spec.Callable{Name: "callback", Params: t.Function.Params, Results: t.Function.Results, Fallible: t.Function.Fallible, Cancellable: t.Function.Cancellable}
		signature, err := e.Signature("", c, true)
		if err != nil {
			return "", err
		}
		value = "func" + strings.TrimPrefix(signature, "callback")
		nilable = true
	default:
		if variable, ok := e.TypeVariables[t.Ref]; ok {
			value = variable
		} else if typ, ok := e.Module.TypesByID[t.Ref]; ok {
			value = typ.Name
			nilable = typ.Kind == "interface" || typ.Kind == "callback"
			if typ.Kind == "alias" {
				var err error
				nilable, err = e.aliasNilable(typ, t.Arguments, map[string]bool{})
				if err != nil {
					return "", err
				}
			}
			if len(t.Arguments) != len(typ.Generics) {
				return "", fmt.Errorf("type %s requires %d type arguments", typ.Name, len(typ.Generics))
			}
		} else if strings.HasPrefix(t.Ref, "schema_enum.") {
			full := strings.TrimPrefix(t.Ref, "schema_enum.")
			var found *protogen.Enum
			var visit func([]*protogen.Message)
			visit = func(messages []*protogen.Message) {
				for _, m := range messages {
					for _, enum := range m.Enums {
						if string(enum.Desc.FullName()) == full {
							found = enum
						}
					}
					visit(m.Messages)
				}
			}
			for _, file := range e.Plugin.Files {
				for _, enum := range file.Enums {
					if string(enum.Desc.FullName()) == full {
						found = enum
					}
				}
				visit(file.Messages)
			}
			if found == nil {
				return "", fmt.Errorf("unknown schema enum %s", full)
			}
			value = e.G.QualifiedGoIdent(found.GoIdent)
		} else if strings.HasPrefix(t.Ref, "schema.") {
			typ, err := schema.Message(e.Plugin, strings.TrimPrefix(t.Ref, "schema."))
			if err != nil {
				return "", err
			}
			value = e.G.QualifiedGoIdent(typ.GoIdent)
		} else if native, ok := e.Native[t.Ref]; ok {
			value = e.G.QualifiedGoIdent(native)
		} else if external, ok := e.Profile.Externals[t.Ref]; ok {
			if err := goIdentifier(external.Name); err != nil {
				return "", err
			}
			value = external.Name
			if external.Import != "" {
				value = e.G.QualifiedGoIdent(protogen.GoIdent{GoName: external.Name, GoImportPath: protogen.GoImportPath(external.Import)})
			}
			nilable = external.Kind == "interface"
		} else {
			return "", fmt.Errorf("unknown API type reference %s", t.Ref)
		}
	}
	if len(t.Arguments) > 0 {
		arguments := []string{}
		for _, argument := range t.Arguments {
			rendered, err := e.Type(argument)
			if err != nil {
				return "", err
			}
			arguments = append(arguments, rendered)
		}
		value += "[" + strings.Join(arguments, ",") + "]"
	}
	if t.Reference && !nilable || ((t.Optional || t.Nullable) && !nilable) {
		value = "*" + value
	}
	return value, nil
}

func (e *Emitter) aliasNilable(typ spec.APIType, arguments []spec.TypeRef, seen map[string]bool) (bool, error) {
	if seen[typ.ID] {
		return false, fmt.Errorf("cyclic type alias %s", typ.Name)
	}
	seen[typ.ID] = true
	defer delete(seen, typ.ID)
	if typ.Underlying == nil {
		return false, fmt.Errorf("alias has no underlying type")
	}
	underlying := *typ.Underlying
	for i, parameter := range typ.Generics {
		if underlying.Ref == parameter.ID && i < len(arguments) {
			underlying = arguments[i]
		}
	}
	if underlying.Optional || underlying.Nullable || underlying.Reference || underlying.Sequence != nil || underlying.Map != nil || underlying.Function != nil || underlying.Builtin == "any" || underlying.Builtin == "bytes" {
		return true, nil
	}
	if target, ok := e.Module.TypesByID[underlying.Ref]; ok {
		if target.Kind == "alias" {
			return e.aliasNilable(target, underlying.Arguments, seen)
		}
		return target.Kind == "interface" || target.Kind == "callback", nil
	}
	if external, ok := e.Profile.Externals[underlying.Ref]; ok {
		return external.Kind == "interface", nil
	}
	return false, nil
}

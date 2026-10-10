package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

// API contains validated declarations. Language profiles are selected separately.
type API struct{ Modules map[string]APIModule }
type APIModule struct {
	Attachments   []spec.APIAttachment
	Types         []spec.APIType
	Functions     []spec.Callable
	TypesByID     map[string]spec.APIType
	FunctionsByID map[string]spec.Callable
}

var declarationIdentifier = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func CompileAPI(input spec.API) (API, error) {
	result := API{Modules: map[string]APIModule{}}
	for module, raw := range input.Modules {
		if !declarationIdentifier.MatchString(module) {
			return result, fmt.Errorf("invalid API module %q", module)
		}
		compiled := APIModule{Attachments: raw.Attachments, Types: raw.Types, Functions: raw.Functions, TypesByID: map[string]spec.APIType{}, FunctionsByID: map[string]spec.Callable{}}
		for _, attachment := range raw.Attachments {
			if !strings.HasPrefix(attachment.Target, "schema.") {
				return result, fmt.Errorf("attachments require a schema target")
			}
			if err := validateSlots(attachment.Fields, true); err != nil {
				return result, err
			}
			for _, field := range attachment.Fields {
				if !field.RuntimeOnly {
					return result, fmt.Errorf("shared serialized fields belong in protobuf")
				}
			}
		}
		names := map[string]bool{}
		for _, typ := range raw.Types {
			if !declarationIdentifier.MatchString(typ.ID) || !declarationIdentifier.MatchString(typ.Name) || names[typ.Name] {
				return result, fmt.Errorf("invalid or duplicate type %q", typ.Name)
			}
			if _, ok := compiled.TypesByID[typ.ID]; ok {
				return result, fmt.Errorf("duplicate type ID %s", typ.ID)
			}
			if typ.Kind != "object" && typ.Kind != "interface" && typ.Kind != "alias" && typ.Kind != "enum" && typ.Kind != "callback" {
				return result, fmt.Errorf("unsupported API type kind %s", typ.Kind)
			}
			if typ.Kind != "object" && (len(typ.Fields) > 0 || typ.FieldsFrom != "") {
				return result, fmt.Errorf("interfaces cannot declare storage")
			}
			if typ.FieldsFrom != "" && typ.FieldsFrom != "provider_state" {
				return result, fmt.Errorf("unsupported field source %s", typ.FieldsFrom)
			}
			switch typ.Visibility {
			case "", "public", "private":
			default:
				return result, fmt.Errorf("unknown type visibility %s", typ.Visibility)
			}
			if err := validateSlots(typ.Generics, true); err != nil {
				return result, err
			}
			if typ.Kind == "alias" || typ.Kind == "enum" {
				if typ.Underlying == nil {
					return result, fmt.Errorf("%s requires underlying type", typ.Kind)
				}
				if err := validateTypeRef(*typ.Underlying); err != nil {
					return result, err
				}
			} else if typ.Underlying != nil {
				return result, fmt.Errorf("unexpected underlying type")
			}
			if typ.Kind == "enum" {
				values := map[string]bool{}
				for _, value := range typ.Values {
					if !declarationIdentifier.MatchString(value.Name) || values[value.Name] {
						return result, fmt.Errorf("invalid or duplicate enum member")
					}
					values[value.Name] = true
				}
			} else if len(typ.Values) > 0 {
				return result, fmt.Errorf("enum values require enum kind")
			}
			if typ.Kind == "callback" {
				if typ.Signature == nil {
					return result, fmt.Errorf("callback needs a signature")
				}
				if err := validateFunctionType(*typ.Signature); err != nil {
					return result, err
				}
			} else if typ.Signature != nil {
				return result, fmt.Errorf("unexpected callback signature")
			}
			names[typ.Name] = true
			compiled.TypesByID[typ.ID] = typ
		}
		for _, typ := range raw.Types {
			if err := validateSlots(typ.Fields, true); err != nil {
				return result, err
			}
			methods := map[string]bool{}
			for _, field := range typ.Fields {
				if field.Variadic {
					return result, fmt.Errorf("object fields cannot be variadic")
				}
				if field.Promote && field.Type.Ref == "" {
					return result, fmt.Errorf("promoted fields require a named type")
				}
				methods[field.Name] = true
			}
			ids := map[string]bool{}
			for _, method := range typ.Methods {
				if methods[method.Name] || ids[method.ID] {
					return result, fmt.Errorf("duplicate method %s.%s", typ.Name, method.Name)
				}
				if err := validateCallable(method, typ.Kind == "interface"); err != nil {
					return result, err
				}
				methods[method.Name] = true
				ids[method.ID] = true
			}
		}
		for _, function := range raw.Functions {
			if function.Static {
				return result, fmt.Errorf("static callable must belong to a type")
			}
			if function.Kind == "constructor" {
				typ, ok := compiled.TypesByID[function.Constructs]
				if !ok || typ.Kind != "object" {
					return result, fmt.Errorf("constructor needs a declared object target")
				}
				if len(function.Results) != 1 || function.Results[0].Type.Ref != function.Constructs {
					return result, fmt.Errorf("constructor must return its declared object")
				}
			}
			if names[function.Name] {
				return result, fmt.Errorf("duplicate declaration %s", function.Name)
			}
			if _, ok := compiled.FunctionsByID[function.ID]; ok {
				return result, fmt.Errorf("duplicate function ID %s", function.ID)
			}
			if function.Receiver != "" {
				return result, fmt.Errorf("free function cannot declare a receiver")
			}
			if err := validateCallable(function, false); err != nil {
				return result, err
			}
			names[function.Name] = true
			compiled.FunctionsByID[function.ID] = function
		}
		result.Modules[module] = compiled
	}
	return result, nil
}
func validateCallable(c spec.Callable, contract bool) error {
	switch c.Kind {
	case "", "function", "constructor":
	default:
		return fmt.Errorf("unknown callable kind %s", c.Kind)
	}
	if c.Kind != "constructor" && c.Constructs != "" {
		return fmt.Errorf("constructs requires constructor kind")
	}
	switch c.Effect {
	case "", "pure", "io":
	default:
		return fmt.Errorf("unknown callable effect %s", c.Effect)
	}
	if contract && (c.Static || c.Kind == "constructor") {
		return fmt.Errorf("interface methods cannot be static or constructors")
	}
	if !declarationIdentifier.MatchString(c.ID) || !declarationIdentifier.MatchString(c.Name) {
		return fmt.Errorf("invalid callable %q", c.Name)
	}
	if c.Receiver != "" && !declarationIdentifier.MatchString(c.Receiver) {
		return fmt.Errorf("invalid receiver %s", c.Receiver)
	}
	if contract && c.Behavior != "" {
		return fmt.Errorf("interface methods cannot define behavior")
	}
	if !contract && c.Behavior == "" {
		return fmt.Errorf("missing behavior for %s", c.Name)
	}
	if err := validateSlots(c.Params, true); err != nil {
		return err
	}
	if err := validateSlots(c.Results, c.NamedResults); err != nil {
		return err
	}
	for _, slot := range c.Params {
		if slot.Promote {
			return fmt.Errorf("parameters cannot promote object fields")
		}
	}
	for _, slot := range c.Results {
		if slot.Promote || slot.Variadic {
			return fmt.Errorf("results cannot be promoted or variadic")
		}
	}
	return validateSlots(c.Generics, true)
}
func validateSlots(slots []spec.APISlot, named bool) error {
	names := map[string]bool{}
	ids := map[string]bool{}
	for i, slot := range slots {
		if named && (!declarationIdentifier.MatchString(slot.Name) || !declarationIdentifier.MatchString(slot.ID)) {
			return fmt.Errorf("slot needs a valid name and ID")
		}
		if slot.Name != "" {
			if names[slot.Name] {
				return fmt.Errorf("duplicate slot %s", slot.Name)
			}
			names[slot.Name] = true
		}
		if slot.ID != "" {
			if ids[slot.ID] {
				return fmt.Errorf("duplicate slot ID %s", slot.ID)
			}
			ids[slot.ID] = true
		}
		if slot.Variadic && (i != len(slots)-1 || slot.Type.Sequence == nil) {
			return fmt.Errorf("variadic parameter must be a final sequence")
		}
		if err := validateTypeRef(slot.Type); err != nil {
			return err
		}
	}
	return nil
}
func validateTypeRef(t spec.TypeRef) error {
	shapes := 0
	if t.Function != nil {
		shapes++
		if err := validateFunctionType(*t.Function); err != nil {
			return err
		}
	}
	if t.Tuple != nil {
		shapes++
		if len(t.Tuple) == 0 {
			return fmt.Errorf("empty tuple")
		}
		for _, item := range t.Tuple {
			if err := validateTypeRef(item); err != nil {
				return err
			}
		}
	}
	if t.Ref != "" {
		shapes++
		for _, part := range strings.Split(t.Ref, ".") {
			if !declarationIdentifier.MatchString(part) {
				return fmt.Errorf("invalid type reference %s", t.Ref)
			}
		}
	}
	if t.Builtin != "" {
		shapes++
	}
	if t.Sequence != nil {
		shapes++
		if err := validateTypeRef(*t.Sequence); err != nil {
			return err
		}
	}
	if t.Map != nil {
		shapes++
		if err := validateTypeRef(t.Map.Key); err != nil {
			return err
		}
		if err := validateTypeRef(t.Map.Value); err != nil {
			return err
		}
	}
	if t.Union != nil {
		shapes++
		if len(t.Union) < 2 {
			return fmt.Errorf("union needs at least two alternatives")
		}
		for _, member := range t.Union {
			if err := validateTypeRef(member); err != nil {
				return err
			}
		}
	}
	if shapes != 1 {
		return fmt.Errorf("type reference needs exactly one shape")
	}
	if len(t.Arguments) > 0 && t.Ref == "" {
		return fmt.Errorf("type arguments require a named reference")
	}
	for _, argument := range t.Arguments {
		if err := validateTypeRef(argument); err != nil {
			return err
		}
	}
	switch t.Ownership {
	case "", "owned", "borrowed", "shared":
	default:
		return fmt.Errorf("unknown ownership %s", t.Ownership)
	}
	return nil
}

func validateFunctionType(signature spec.FunctionType) error {
	return validateCallable(spec.Callable{ID: "callback", Name: "callback", Behavior: "callback", Params: signature.Params, Results: signature.Results}, false)
}

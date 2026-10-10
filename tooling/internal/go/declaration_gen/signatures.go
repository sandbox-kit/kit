package declarationgen

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func (e *Emitter) Signature(owner string, c spec.Callable, contract bool) (string, error) {
	if !contract && (c.Behavior == "initialize" || c.Behavior == "create" || c.Behavior == "close") && !c.NamedResults {
		return "", fmt.Errorf("Go behavior %s requires named results for deferred error context", c.Behavior)
	}
	if c.Static {
		if c.Receiver != "" {
			return "", fmt.Errorf("static methods cannot have a receiver variable")
		}
		if typ, ok := e.Module.TypesByID[owner]; ok {
			c.Generics = append(append([]spec.APISlot{}, typ.Generics...), c.Generics...)
		}
	}
	local := *e
	e = &local
	inherited := e.TypeVariables
	e.TypeVariables = map[string]string{}
	for id, name := range inherited {
		e.TypeVariables[id] = name
	}
	if typ, ok := e.Module.TypesByID[owner]; ok {
		for _, parameter := range typ.Generics {
			e.TypeVariables[parameter.ID] = parameter.Name
		}
	}
	if err := goIdentifier(c.Name); err != nil {
		return "", err
	}
	var generic string
	if len(c.Generics) > 0 {
		if owner != "" && !c.Static {
			return "", fmt.Errorf("Go methods cannot declare independent type parameters")
		}
		parts := []string{}
		for _, parameter := range c.Generics {
			if err := goIdentifier(parameter.Name); err != nil {
				return "", err
			}
			e.TypeVariables[parameter.ID] = parameter.Name
			constraint, err := e.Type(parameter.Type)
			if err != nil {
				return "", err
			}
			parts = append(parts, parameter.Name+" "+constraint)
		}
		generic = "[" + strings.Join(parts, ",") + "]"
	}
	params := []string{}
	seen := map[string]bool{}
	if c.Cancellable {
		name := e.Profile.Cancellation.Name
		if err := goIdentifier(name); err != nil {
			return "", err
		}
		typ := e.Profile.Cancellation.Type
		value := e.G.QualifiedGoIdent(protogen.GoIdent{GoName: typ.Name, GoImportPath: protogen.GoImportPath(typ.Import)})
		if !contract {
			value = name + " " + value
		}
		params = append(params, value)
		seen[name] = true
	}
	for _, parameter := range c.Params {
		if err := goIdentifier(parameter.Name); err != nil {
			return "", err
		}
		if seen[parameter.Name] {
			return "", fmt.Errorf("duplicate Go parameter %s", parameter.Name)
		}
		seen[parameter.Name] = true
		value, err := e.Type(parameter.Type)
		if err != nil {
			return "", err
		}
		if parameter.Variadic {
			if parameter.Type.Sequence == nil {
				return "", fmt.Errorf("variadic parameter must be a sequence")
			}
			value, err = e.Type(*parameter.Type.Sequence)
			if err != nil {
				return "", err
			}
			value = "..." + value
		}
		if !contract {
			value = parameter.Name + " " + value
		}
		params = append(params, value)
	}
	results := []string{}
	for _, result := range c.Results {
		value, err := e.Type(result.Type)
		if err != nil {
			return "", err
		}
		if c.NamedResults && !contract {
			if err := goIdentifier(result.Name); err != nil {
				return "", err
			}
			if seen[result.Name] {
				return "", fmt.Errorf("result name collides with parameter %s", result.Name)
			}
			seen[result.Name] = true
			value = result.Name + " " + value
		}
		results = append(results, value)
	}
	if c.Fallible {
		value := "error"
		if c.NamedResults && !contract {
			if seen[e.Profile.ErrorResultName] {
				return "", fmt.Errorf("error result name collision")
			}
			value = e.Profile.ErrorResultName + " error"
		}
		results = append(results, value)
	}
	returned := ""
	if len(results) > 0 {
		returned = " " + strings.Join(results, ",")
		if len(results) > 1 || c.NamedResults && !contract {
			returned = " (" + strings.Join(results, ",") + ")"
		}
	}
	prefix := ""
	if !contract {
		prefix = "func "
		if owner != "" && !c.Static {
			typ, ok := e.Module.TypesByID[owner]
			if !ok || (typ.Kind != "object" && typ.Kind != "enum") {
				return "", fmt.Errorf("method receiver must be a declared object")
			}
			receiver := typ.Name
			if typ.Kind == "object" {
				receiver = "*" + receiver
			}
			if len(typ.Generics) > 0 {
				arguments := []string{}
				for _, parameter := range typ.Generics {
					arguments = append(arguments, parameter.Name)
				}
				receiver += "[" + strings.Join(arguments, ",") + "]"
			}
			if c.Receiver != "" {
				if err := goIdentifier(c.Receiver); err != nil {
					return "", err
				}
				receiver = c.Receiver + " " + receiver
			}
			prefix += "(" + receiver + ")"
		}
	}
	name := c.Name
	if c.Static {
		typ, ok := e.Module.TypesByID[owner]
		if !ok {
			return "", fmt.Errorf("static method requires an owner")
		}
		name = typ.Name + c.Name
	}
	return prefix + name + generic + "(" + strings.Join(params, ",") + ")" + returned, nil
}

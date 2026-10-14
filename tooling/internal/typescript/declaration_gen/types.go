package declarationgen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
)

type Renderer struct {
	Module  model.APIModule
	Profile spec.LanguageProfile
	vars    map[string]string
}

func NewRenderer(module model.APIModule, profile spec.LanguageProfile) (*Renderer, error) {
	if err := ValidateProfile(profile); err != nil {
		return nil, err
	}
	return &Renderer{Module: module, Profile: profile, vars: map[string]string{}}, nil
}

// ValidateModule proves that every declaration in the compiled API module has
// a lossless TypeScript representation before runtime source is emitted.
func ValidateModule(module model.APIModule, profile spec.LanguageProfile) error {
	renderer, err := NewRenderer(module, profile)
	if err != nil {
		return err
	}
	for _, declared := range module.Types {
		for _, generic := range declared.Generics {
			renderer.vars[generic.ID] = generic.Name
		}
		if declared.Underlying != nil {
			if _, err := renderer.Type(*declared.Underlying); err != nil {
				return fmt.Errorf("type %s: %w", declared.Name, err)
			}
		}
		for _, field := range declared.Fields {
			if _, err := renderer.Type(field.Type); err != nil {
				return fmt.Errorf("field %s.%s: %w", declared.Name, field.Name, err)
			}
		}
		if declared.Signature != nil {
			if _, err := renderer.function(*declared.Signature); err != nil {
				return fmt.Errorf("callback %s: %w", declared.Name, err)
			}
		}
		for _, method := range declared.Methods {
			if err := renderer.validateCallable(method); err != nil {
				return fmt.Errorf("method %s.%s: %w", declared.Name, method.Name, err)
			}
		}
		for _, generic := range declared.Generics {
			delete(renderer.vars, generic.ID)
		}
	}
	for _, function := range module.Functions {
		if err := renderer.validateCallable(function); err != nil {
			return fmt.Errorf("function %s: %w", function.Name, err)
		}
	}
	return nil
}

func (r *Renderer) validateCallable(callable spec.Callable) error {
	for _, generic := range callable.Generics {
		r.vars[generic.ID] = generic.Name
	}
	defer func() {
		for _, generic := range callable.Generics {
			delete(r.vars, generic.ID)
		}
	}()
	for _, parameter := range callable.Params {
		if _, err := r.Type(parameter.Type); err != nil {
			return err
		}
	}
	for _, result := range callable.Results {
		if _, err := r.Type(result.Type); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) Type(value spec.TypeRef) (string, error) {
	var rendered string
	switch {
	case value.Builtin != "":
		rendered = r.Profile.Builtins[value.Builtin]
		if rendered == "" {
			return "", fmt.Errorf("unsupported TypeScript builtin %s", value.Builtin)
		}
	case value.Sequence != nil:
		element, err := r.Type(*value.Sequence)
		if err != nil {
			return "", err
		}
		rendered = "readonly " + element + "[]"
	case value.Map != nil:
		key, err := r.Type(value.Map.Key)
		if err != nil {
			return "", err
		}
		if key != "string" && key != "number" && key != "symbol" {
			return "", fmt.Errorf("TypeScript map key must be string, number, or symbol")
		}
		element, err := r.Type(value.Map.Value)
		if err != nil {
			return "", err
		}
		rendered = "Readonly<Record<" + key + ", " + element + ">>"
	case value.Tuple != nil:
		items := make([]string, 0, len(value.Tuple))
		for _, item := range value.Tuple {
			part, err := r.Type(item)
			if err != nil {
				return "", err
			}
			items = append(items, part)
		}
		rendered = "readonly [" + strings.Join(items, ", ") + "]"
	case value.Union != nil:
		items := make([]string, 0, len(value.Union))
		for _, item := range value.Union {
			part, err := r.Type(item)
			if err != nil {
				return "", err
			}
			items = append(items, part)
		}
		rendered = strings.Join(items, " | ")
	case value.Function != nil:
		signature, err := r.function(*value.Function)
		if err != nil {
			return "", err
		}
		rendered = signature
	case value.Ref != "":
		if variable := r.vars[value.Ref]; variable != "" {
			rendered = variable
		} else if declared, ok := r.Module.TypesByID[value.Ref]; ok {
			rendered = declared.Name
		} else if strings.HasPrefix(value.Ref, "schema.") || strings.HasPrefix(value.Ref, "schema_enum.") {
			parts := strings.Split(value.Ref, ".")
			rendered = parts[len(parts)-1]
		} else if external, ok := r.Profile.Externals[value.Ref]; ok {
			rendered = external.Name
		} else {
			return "", fmt.Errorf("unknown TypeScript API type reference %s", value.Ref)
		}
	default:
		return "", fmt.Errorf("TypeScript type reference has no shape")
	}
	if len(value.Arguments) > 0 {
		arguments := make([]string, 0, len(value.Arguments))
		for _, argument := range value.Arguments {
			part, err := r.Type(argument)
			if err != nil {
				return "", err
			}
			arguments = append(arguments, part)
		}
		rendered += "<" + strings.Join(arguments, ", ") + ">"
	}
	if value.Nullable {
		rendered += " | null"
	}
	if value.Optional {
		rendered += " | undefined"
	}
	return rendered, nil
}

func (r *Renderer) function(value spec.FunctionType) (string, error) {
	params := make([]string, 0, len(value.Params)+1)
	if value.Cancellable {
		params = append(params, r.Profile.Cancellation.Name+"?: "+r.Profile.Cancellation.Type.Name)
	}
	for _, parameter := range value.Params {
		typeName, err := r.Type(parameter.Type)
		if err != nil {
			return "", err
		}
		marker := ""
		if parameter.Type.Optional {
			marker = "?"
			typeName = strings.TrimSuffix(typeName, " | undefined")
		}
		params = append(params, naming.Camel(parameter.Name)+marker+": "+typeName)
	}
	result := "void"
	if len(value.Results) == 1 {
		var err error
		result, err = r.Type(value.Results[0].Type)
		if err != nil {
			return "", err
		}
	} else if len(value.Results) > 1 {
		fields := make([]string, 0, len(value.Results))
		for _, output := range value.Results {
			typeName, err := r.Type(output.Type)
			if err != nil {
				return "", err
			}
			fields = append(fields, "readonly "+naming.Camel(output.Name)+": "+typeName)
		}
		result = "{ " + strings.Join(fields, "; ") + " }"
	}
	if value.Fallible {
		result = "Promise<" + result + ">"
	}
	return "(" + strings.Join(params, ", ") + ") => " + result, nil
}

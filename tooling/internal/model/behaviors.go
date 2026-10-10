package model

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"reflect"
)

type Behaviors struct {
	Signatures map[string]spec.BehaviorSignature
}

func CompileBehaviors(raw spec.Behaviors) (Behaviors, error) {
	result := Behaviors{Signatures: raw.Signatures}
	if raw.Version != 1 || len(raw.Signatures) == 0 {
		return result, fmt.Errorf("unsupported or empty behavior contracts")
	}
	for id, signature := range raw.Signatures {
		if !declarationIdentifier.MatchString(id) {
			return result, fmt.Errorf("invalid behavior ID %s", id)
		}
		if signature.Owner != "method" && signature.Owner != "function" && signature.Owner != "either" {
			return result, fmt.Errorf("invalid behavior owner %s", id)
		}
		for _, slots := range [][]spec.BehaviorSlot{signature.Params, signature.Results, signature.Generics} {
			for _, slot := range slots {
				if err := validateTypeRef(slot.Type); err != nil {
					return result, fmt.Errorf("behavior %s: %w", id, err)
				}
			}
		}
	}
	return result, nil
}
func (b Behaviors) Validate(c spec.Callable, method bool) error {
	if c.Behavior == "" {
		return nil
	}
	expected, ok := b.Signatures[c.Behavior]
	if !ok {
		return fmt.Errorf("unknown behavior %s", c.Behavior)
	}
	if expected.Owner == "method" && !method || expected.Owner == "function" && method {
		return fmt.Errorf("behavior %s has incompatible declaration owner", c.Behavior)
	}
	if c.Fallible != expected.Fallible || c.Cancellable != expected.Cancellable {
		return fmt.Errorf("behavior %s requires fallible=%t cancellable=%t", c.Behavior, expected.Fallible, expected.Cancellable)
	}
	for _, slots := range []struct {
		name     string
		actual   []spec.APISlot
		required []spec.BehaviorSlot
	}{{"params", c.Params, expected.Params}, {"results", c.Results, expected.Results}, {"generics", c.Generics, expected.Generics}} {
		if len(slots.actual) != len(slots.required) {
			return fmt.Errorf("behavior %s requires %d %s", c.Behavior, len(slots.required), slots.name)
		}
		for i, actual := range slots.actual {
			required := slots.required[i]
			// IDs connect behavior references; names are intentionally free to change.
			if required.ID != "" && actual.ID != required.ID || actual.Variadic != required.Variadic || !reflect.DeepEqual(actual.Type, required.Type) {
				return fmt.Errorf("behavior %s: incompatible %s[%d] (%s)", c.Behavior, slots.name, i, required.ID)
			}
		}
	}
	return nil
}
func (b Behaviors) ValidateAPI(api API) error {
	for module, surface := range api.Modules {
		for _, typ := range surface.Types {
			for _, method := range typ.Methods {
				if err := b.Validate(method, true); err != nil {
					return fmt.Errorf("api.%s.%s.%s: %w", module, typ.ID, method.ID, err)
				}
			}
		}
		for _, function := range surface.Functions {
			if err := b.Validate(function, false); err != nil {
				return fmt.Errorf("api.%s.%s: %w", module, function.ID, err)
			}
		}
	}
	return nil
}
func (b Behaviors) ValidateTemplates(templates Templates) error {
	for id, template := range templates.Callables {
		if err := b.Validate(template.Declaration, template.Owner != ""); err != nil {
			return fmt.Errorf("template.%s: %w", id, err)
		}
	}
	return nil
}

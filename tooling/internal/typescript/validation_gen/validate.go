// Package validationgen validates compiled rules before TypeScript emission.
package validationgen

import (
	"fmt"
	"go/token"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func ValidateProvider(provider model.Provider) error {
	seen := map[string]bool{}
	for _, group := range provider.Groups {
		name := group.Rule.Name
		if !token.IsIdentifier(name) || seen[name] {
			return fmt.Errorf("invalid/duplicate TypeScript mapping group %q", name)
		}
		seen[name] = true
		if group.Rule.Direction != "request" && group.Rule.Direction != "response" {
			return fmt.Errorf("%s: invalid TypeScript mapping direction %q", name, group.Rule.Direction)
		}
		if err := validateRepresentations(name, group.Rule.Target.Bindings["typescript"]); err != nil {
			return err
		}
		if err := validateRepresentations(name, group.Rule.Source.Bindings["typescript"]); err != nil {
			return err
		}
	}
	return nil
}

func validateRepresentations(group string, binding spec.Binding) error {
	for field, representation := range binding.Representations {
		if _, ok := binding.Fields[field]; !ok {
			return fmt.Errorf("%s: representation references unknown native field %q", group, field)
		}
		switch representation.Transform {
		case "", "whole", "scaled_integer", "duration_seconds", "divide_exactly", "multiply":
		default:
			return fmt.Errorf("%s.%s: unknown TypeScript representation transform %q", group, field, representation.Transform)
		}
		if (representation.Transform == "scaled_integer" || representation.Transform == "divide_exactly" || representation.Transform == "multiply") && representation.Factor == 0 {
			return fmt.Errorf("%s.%s: TypeScript representation transform requires a positive factor", group, field)
		}
		switch representation.Cast {
		case "", "integer", "number", "unsigned", "string":
		default:
			return fmt.Errorf("%s.%s: unknown TypeScript representation cast %q", group, field, representation.Cast)
		}
	}
	return nil
}

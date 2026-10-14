package mappinggen

import (
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func TestValidateBindingsUsesTargetForRequestsAndSourceForResponses(t *testing.T) {
	provider := model.Provider{
		Rules: spec.Provider{Provider: "example"},
		Groups: []model.Group{
			{Rule: spec.Group{Name: "request", Target: boundType("native", map[string]string{"value": "requestValue"})}},
			{Rule: spec.Group{Name: "response", Direction: "response", Source: boundType("native", map[string]string{"value": "responseValue"})}},
		},
	}
	required := map[string]Binding{
		"request":  {Import: "native", Fields: map[string]string{"value": "requestValue"}},
		"response": {Import: "native", Fields: map[string]string{"value": "responseValue"}},
	}

	if err := ValidateBindings(provider, required); err != nil {
		t.Fatalf("ValidateBindings() error = %v", err)
	}
	if len(required) != 2 {
		t.Fatalf("ValidateBindings() mutated required bindings: got %d entries", len(required))
	}
}

func TestValidateBindingsReportsSortedMissingGroups(t *testing.T) {
	provider := model.Provider{Rules: spec.Provider{Provider: "example"}}
	err := ValidateBindings(provider, map[string]Binding{
		"zeta":  {Import: "native"},
		"alpha": {Import: "native"},
	})
	if err == nil || !strings.Contains(err.Error(), "missing mapping groups alpha, zeta") {
		t.Fatalf("ValidateBindings() error = %v", err)
	}
}

func TestValidateBindingsRejectsWrongNativeMember(t *testing.T) {
	provider := model.Provider{
		Rules: spec.Provider{Provider: "example"},
		Groups: []model.Group{{Rule: spec.Group{
			Name:   "request",
			Target: boundType("native", map[string]string{"value": "wrong"}),
		}}},
	}
	err := ValidateBindings(provider, map[string]Binding{
		"request": {Import: "native", Fields: map[string]string{"value": "expected"}},
	})
	if err == nil || !strings.Contains(err.Error(), "request.value must bind to expected") {
		t.Fatalf("ValidateBindings() error = %v", err)
	}
}

func boundType(importPath string, fields map[string]string) spec.Type {
	return spec.Type{Bindings: map[string]spec.Binding{
		"typescript": {Import: importPath, Fields: fields},
	}}
}

package validationgen

import (
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func TestValidateProviderRejectsUnknownRepresentationField(t *testing.T) {
	provider := model.Provider{
		Rules: spec.Provider{Provider: "example"},
		Groups: []model.Group{{Rule: spec.Group{
			Name:      "mapValue",
			Direction: "request",
			Target: spec.Type{Bindings: map[string]spec.Binding{"typescript": {
				Fields:          map[string]string{"value": "nativeValue"},
				Representations: map[string]spec.Representation{"missing": {}},
			}}},
		}}},
	}

	err := ValidateProvider(provider)
	if err == nil || !strings.Contains(err.Error(), "unknown native field") {
		t.Fatalf("ValidateProvider() error = %v", err)
	}
}

func TestValidateProviderAcceptsLanguageRepresentationOverride(t *testing.T) {
	maximum := uint64(86_400_000)
	provider := model.Provider{
		Rules: spec.Provider{Provider: "example"},
		Groups: []model.Group{{Rule: spec.Group{
			Name:      "mapDuration",
			Direction: "request",
			Target: spec.Type{Bindings: map[string]spec.Binding{"typescript": {
				Fields:          map[string]string{"timeout": "timeoutMs"},
				Representations: map[string]spec.Representation{"timeout": {Maximum: &maximum}},
			}}},
		}}},
	}

	if err := ValidateProvider(provider); err != nil {
		t.Fatalf("ValidateProvider() error = %v", err)
	}
}

package declarationgen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func ValidateProfile(profile spec.LanguageProfile) error {
	if profile.Language != "typescript" || profile.ObjectKind != "class" || profile.InterfaceKind != "interface" || profile.Optional != "property" || profile.Fallible != "rejected_promise" || profile.MultipleResults != "object" || profile.Unions != "discriminated" {
		return fmt.Errorf("unsupported TypeScript language profile")
	}
	for _, role := range []string{"shared", "client", "provider"} {
		if !strings.HasSuffix(profile.Outputs[role], ".ts") {
			return fmt.Errorf("TypeScript output %s must use a .ts suffix", role)
		}
	}
	return nil
}

package declarationgen

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"strings"
)

// ValidateProfile owns Go capabilities and output conventions. Shared loaders
// accept other languages without enforcing Go file extensions or syntax.
func ValidateProfile(profile spec.LanguageProfile) error {
	if profile.Language != "go" || profile.ObjectKind != "struct" || profile.InterfaceKind != "interface" || profile.Reference != "pointer" || profile.Optional != "pointer" || profile.Nullable != "nil" || profile.Fallible != "error_result" || profile.MultipleResults != "tuple" || profile.Unions != "reject" {
		return fmt.Errorf("unsupported Go declaration profile")
	}
	for _, role := range []string{"client", "provider", "errors", "copies"} {
		if profile.Outputs[role] == "" {
			return fmt.Errorf("missing Go output role %s", role)
		}
	}
	for role, suffix := range profile.Outputs {
		if !strings.HasSuffix(suffix, ".gen.go") {
			return fmt.Errorf("Go output %s must end in .gen.go", role)
		}
	}
	return nil
}

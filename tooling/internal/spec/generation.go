package spec

import (
	"fmt"
	"path/filepath"
)

// Generation selects semantic output phases, independent of language syntax.
type Generation struct {
	Version    int                 `yaml:"version"`
	Shared     []string            `yaml:"shared"`
	Client     []string            `yaml:"client"`
	Provider   []string            `yaml:"provider"`
	Operations map[string][]string `yaml:"operations"`
}

func LoadGeneration(root string) (Generation, error) {
	var result Generation
	if err := decode(filepath.Join(root, "generation.yaml"), &result); err != nil {
		return result, err
	}
	if result.Version != 1 {
		return result, fmt.Errorf("unsupported generation version %d", result.Version)
	}
	roles := []struct {
		name    string
		phases  []string
		allowed []string
	}{
		{"shared", result.Shared, []string{"types", "copies", "diagnostic_paths", "field_paths", "validation"}},
		{"client", result.Client, []string{"types", "diagnostic_paths", "field_paths", "validation"}},
		{"provider", result.Provider, []string{"checks", "mappings"}},
	}
	for _, role := range roles {
		if len(role.phases) != len(role.allowed) {
			return result, fmt.Errorf("%s: missing generation phase", role.name)
		}
		seen := map[string]bool{}
		for _, phase := range role.phases {
			valid := false
			for _, allowed := range role.allowed {
				if phase == allowed {
					valid = true
				}
			}
			if !valid || seen[phase] {
				return result, fmt.Errorf("%s: unknown or duplicate phase %q", role.name, phase)
			}
			seen[phase] = true
		}
	}
	if err := ValidateOperations(result.Operations); err != nil {
		return result, err
	}
	return result, nil
}

// ValidateOperations checks the version 1 instruction vocabulary and required
// dependencies. Instructions express behavior, never target-language snippets.
func ValidateOperations(operations map[string][]string) error {
	required := map[string][]string{
		"initialize": {"require_provider", "identify_provider", "validate_configuration", "require_provider_name", "capture_timeout", "initialize_provider", "require_backend", "match_backend_identity", "return_client"},
		"create":     {"require_context", "require_client", "prepare_request", "resolve_deadline", "invoke_creation", "return_handle"},
		"close":      {"require_context", "close_backend"},
	}
	for name, expected := range required {
		steps, ok := operations[name]
		if !ok {
			return fmt.Errorf("missing client operation %s", name)
		}
		if len(steps) != len(expected) {
			return fmt.Errorf("%s: missing or duplicate instruction", name)
		}
		for i, step := range steps {
			if step != expected[i] {
				return fmt.Errorf("%s: instruction %q violates vocabulary or dependencies; expected %q", name, step, expected[i])
			}
		}
	}
	clientCapability := []string{"require_context", "require_client", "require_identity", "resolve_deadline", "require_capability", "invoke_capability", "return_handle"}
	handleUpdate := []string{"require_context", "require_handle", "resolve_deadline", "require_capability", "invoke_capability", "update_handle"}
	handleResult := []string{"require_context", "require_handle", "resolve_deadline", "require_capability", "invoke_capability", "return_result"}
	for name, steps := range operations {
		if _, base := required[name]; base {
			continue
		}
		valid := sameInstructions(steps, clientCapability) || sameInstructions(steps, handleUpdate) || sameInstructions(steps, handleResult)
		if !valid {
			return fmt.Errorf("%s: instruction sequence violates lifecycle vocabulary", name)
		}
	}
	return nil
}

func sameInstructions(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if actual[i] != expected[i] {
			return false
		}
	}
	return true
}

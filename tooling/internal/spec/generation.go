package spec

import (
	"fmt"
	"path/filepath"
)

// Generation selects semantic output phases, independent of language syntax.
type Generation struct {
	Version  int      `yaml:"version"`
	Shared   []string `yaml:"shared"`
	Client   []string `yaml:"client"`
	Provider []string `yaml:"provider"`
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
		name     string
		phases   []string
		allowed  []string
		required bool
	}{
		{"shared", result.Shared, nil, false},
		{"client", result.Client, nil, false},
		{"provider", result.Provider, []string{"mappings"}, true},
	}
	for _, role := range roles {
		if role.required && len(role.phases) != len(role.allowed) {
			return result, fmt.Errorf("%s: missing generation phase", role.name)
		}
		if !role.required && len(role.phases) != 0 {
			return result, fmt.Errorf("%s: no emission phases remain; shared create validation, provider checks, and client helpers are handwritten", role.name)
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
	return result, nil
}

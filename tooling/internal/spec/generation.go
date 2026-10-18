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
		name    string
		phases  []string
		allowed []string
	}{
		{"shared", result.Shared, []string{"diagnostic_paths", "field_paths", "validation"}},
		{"client", result.Client, []string{"diagnostic_paths", "field_paths"}},
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
	return result, nil
}

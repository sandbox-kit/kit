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
		name   string
		phases []string
	}{
		{"shared", result.Shared},
		{"client", result.Client},
		{"provider", result.Provider},
	}
	for _, role := range roles {
		if len(role.phases) != 0 {
			return result, fmt.Errorf("%s: no emission phases remain; validation, checks, mappings, and helpers are handwritten", role.name)
		}
	}
	return result, nil
}

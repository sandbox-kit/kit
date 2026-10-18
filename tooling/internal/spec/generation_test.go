package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerationRejectsUnsupportedOrUnsafePlans(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "..", "specs", "generation.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, old, new string }{
		{"unknown phase", "shared: [diagnostic_paths, field_paths, validation]", "shared: [diagnostic_paths, field_paths, arbitrary_code]"},
		{"missing phase", "shared: [diagnostic_paths, field_paths, validation]", "shared: [diagnostic_paths, validation]"},
		{"duplicate phase", "shared: [diagnostic_paths, field_paths, validation]", "shared: [diagnostic_paths, field_paths, field_paths]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "generation.yaml"), []byte(strings.Replace(string(body), tc.old, tc.new, 1)), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := LoadGeneration(root); err == nil {
				t.Fatal("invalid generation plan accepted")
			}
		})
	}
}

func TestGenerationLoadsPhasePlan(t *testing.T) {
	plan, err := LoadGeneration(filepath.Join("..", "..", "..", "specs"))
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Shared) == 0 || len(plan.Client) == 0 || len(plan.Provider) == 0 {
		t.Fatal("expected shared, client, and provider phases")
	}
}

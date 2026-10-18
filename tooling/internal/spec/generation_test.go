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
		{"unknown phase", "shared: [validation]", "shared: [arbitrary_code]"},
		{"missing phase", "shared: [validation]", "shared: []"},
		{"duplicate phase", "provider: [checks, mappings]", "provider: [checks, mappings, checks]"},
		{"client phases removed", "provider: [checks, mappings]", "client: [validation]\nprovider: [checks, mappings]"},
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
	if len(plan.Shared) != 1 || plan.Shared[0] != "validation" {
		t.Fatalf("shared phases = %#v", plan.Shared)
	}
	if len(plan.Client) != 0 {
		t.Fatalf("client phases should be empty, got %#v", plan.Client)
	}
	if len(plan.Provider) != 2 {
		t.Fatalf("provider phases = %#v", plan.Provider)
	}
}

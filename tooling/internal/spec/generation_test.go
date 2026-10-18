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
		{"unknown shared phase", "shared: []", "shared: [validation]"},
		{"unknown provider phase", "provider: [mappings]", "provider: [mappings, checks]"},
		{"missing provider phase", "provider: [mappings]", "provider: []"},
		{"duplicate phase", "provider: [mappings]", "provider: [mappings, mappings]"},
		{"client phases removed", "provider: [mappings]", "client: [validation]\nprovider: [mappings]"},
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
	if len(plan.Shared) != 0 {
		t.Fatalf("shared phases should be empty, got %#v", plan.Shared)
	}
	if len(plan.Client) != 0 {
		t.Fatalf("client phases should be empty, got %#v", plan.Client)
	}
	if len(plan.Provider) != 1 || plan.Provider[0] != "mappings" {
		t.Fatalf("provider phases = %#v", plan.Provider)
	}
}

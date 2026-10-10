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
		{"unknown phase", "error_runtime", "arbitrary_code"},
		{"missing phase", "copies, ", ""},
		{"duplicate phase", "copies, diagnostic_paths", "copies, copies"},
		{"unknown instruction", "require_provider_name", "execute_go_expression"},
		{"missing instruction", "    - require_backend\n", ""},
		{"dependency ordering", "    - prepare_request\n    - resolve_deadline", "    - resolve_deadline\n    - prepare_request"},
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

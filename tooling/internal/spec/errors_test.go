package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestErrorPolicy(t *testing.T) {
	root := filepath.Join("..", "..", "..", "specs")
	p, err := LoadErrorPolicy(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.CauseRules) != 2 || !p.ExistingError.ReuseMatchingContext {
		t.Fatal("missing shared error policy")
	}
	body, err := os.ReadFile(filepath.Join(root, "errors.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range []struct{ old, new string }{
		{"match: canceled", "match: arbitrary"},
		{"kind: timeout", "kind: arbitrary"},
		{"preserve_details: true", "preserve_details: false"},
		{"native_type, http_status", "http_status, native_type"},
	} {
		t.Run(change.new, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "errors.yaml"), []byte(strings.Replace(string(body), change.old, change.new, 1)), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := LoadErrorPolicy(dir); err == nil {
				t.Fatal("unsupported policy accepted")
			}
		})
	}
}

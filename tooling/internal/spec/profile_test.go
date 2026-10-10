package spec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSharedProfileLoaderAcceptsNonGoConventions(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "languages"), 0700); err != nil {
		t.Fatal(err)
	}
	body := "version: 1\nlanguage: typescript\nobject_kind: class\nfallible: rejected_promise\noutputs: {client: .gen.ts}\n"
	file := filepath.Join(root, "languages", "typescript.yaml")
	if err := os.WriteFile(file, []byte(body), 0600); err != nil {
		t.Fatal(err)
	}
	profile, err := LoadLanguageProfile(root, "typescript")
	if err != nil || profile.Outputs["client"] != ".gen.ts" {
		t.Fatalf("non-Go profile rejected: %v", err)
	}
	if err := os.WriteFile(file, []byte("version: 1\nlanguage: typescript\noutputs: {client: ../outside.ts}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadLanguageProfile(root, "typescript"); err == nil {
		t.Fatal("unsafe output path accepted")
	}
}

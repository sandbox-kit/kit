package spec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositorySpecs(t *testing.T) {
	root := filepath.Join("..", "..", "..", "specs")
	validation, err := LoadValidation(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(validation.Messages) == 0 {
		t.Fatal("missing validation rules")
	}
	client, err := LoadClient(root)
	if err != nil {
		t.Fatal(err)
	}
	if !client.Validation.Messages["APIKeyCredentials"].Fields["key"].Sensitive || client.ProviderContract != "Provider" {
		t.Fatal("missing initialization or credential metadata")
	}

	for _, name := range []string{"modal", "daytona"} {
		p, err := LoadProvider(root, name)
		if err != nil {
			t.Fatal(err)
		}
		if len(p.Groups) < 2 {
			t.Fatal("missing provider mappings")
		}
		for _, g := range p.Groups {
			for _, typ := range []Type{g.Source, g.Target} {
				if len(typ.Bindings) > 0 {
					if typ.Bindings["go"].Import == "" {
						t.Fatal("missing Go binding")
					}
				}
			}
		}
	}
}
func TestRejectsInvalidDocuments(t *testing.T) {
	for _, body := range []string{"version: 9\nmessages: {}\n", "version: 1\nunknown: true\n", "version: 1\nmessages: {}\n---\nversion: 1\n"} {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "validation.yaml"), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := LoadValidation(root); err == nil {
			t.Fatalf("accepted invalid document %q", body)
		}
	}
}

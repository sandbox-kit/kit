package model

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func TestTemplateBindingAndValidation(t *testing.T) {
	raw, err := spec.LoadTemplates(filepath.Join("..", "..", "..", "specs"))
	if err != nil {
		t.Fatal(err)
	}
	templates, err := CompileTemplates(raw)
	if err != nil {
		t.Fatal(err)
	}
	mapping, err := templates.Instantiate("request_mapping", map[string]string{"mapping": "convertResources"})
	if err != nil {
		t.Fatal(err)
	}
	if mapping.Declaration.Name != "convertResources" || !strings.Contains(mapping.Declaration.Doc, "convertResources") {
		t.Fatal("template did not preserve mapping name")
	}
	if _, err := templates.Instantiate("request_mapping", nil); err == nil {
		t.Fatal("missing binding accepted")
	}
	if _, err := templates.Instantiate("request_mapping", map[string]string{"mapping": "Value;panic()"}); err == nil {
		t.Fatal("code accepted in a name binding")
	}
	pattern := raw.Callables["request_mapping"]
	pattern.Declaration.Name = "func arbitrary()"
	raw.Callables["request_mapping"] = pattern
	if _, err := CompileTemplates(raw); err == nil {
		t.Fatal("code accepted as a declaration")
	}
}

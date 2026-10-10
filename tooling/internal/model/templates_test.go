package model

import (
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"path/filepath"
	"testing"
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
	getter, err := templates.Instantiate("getter", map[string]string{"field": "CPUCores"})
	if err != nil {
		t.Fatal(err)
	}
	if getter.Declaration.Name != "GetCPUCores" || getter.Declaration.Results[0].Type.Ref != "field_value" {
		t.Fatal("template did not preserve projection")
	}
	if _, err := templates.Instantiate("getter", nil); err == nil {
		t.Fatal("missing binding accepted")
	}
	if _, err := templates.Instantiate("getter", map[string]string{"field": "Value;panic()"}); err == nil {
		t.Fatal("code accepted in a name binding")
	}
	pattern := raw.Callables["getter"]
	pattern.Declaration.Name = "func arbitrary()"
	raw.Callables["getter"] = pattern
	if _, err := CompileTemplates(raw); err == nil {
		t.Fatal("code accepted as a declaration")
	}
}

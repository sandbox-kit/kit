package model

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func TestTemplatesLoadEmptyCallablesAndFieldPaths(t *testing.T) {
	raw, err := spec.LoadTemplates(filepath.Join("..", "..", "..", "specs"))
	if err != nil {
		t.Fatal(err)
	}
	templates, err := CompileTemplates(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(templates.Callables) != 0 {
		t.Fatalf("expected empty callables, got %#v", templates.Callables)
	}
	if templates.FieldPaths.Name != "{prefix}{fields}" || templates.FieldPaths.Roots["CreateOptions"] != "CreateField" {
		t.Fatalf("field-path inventory lost: %#v", templates.FieldPaths)
	}
}

func TestTemplateBindingAndValidation(t *testing.T) {
	raw := spec.Templates{
		Version: 1,
		Callables: map[string]spec.DeclarationTemplate{
			"request_mapping": {
				Declaration: spec.Callable{
					ID: "request_mapping", Name: "{mapping}", Behavior: "request_mapping", Fallible: true,
					Doc: "maps {mapping}",
					Params: []spec.APISlot{
						{ID: "source", Name: "source", Type: spec.TypeRef{Ref: "source", Reference: true, Nullable: true}},
					},
					Results: []spec.APISlot{
						{Type: spec.TypeRef{Ref: "target"}},
						{Type: spec.TypeRef{Ref: "source"}},
					},
				},
			},
		},
		FieldPaths: spec.FieldPathTemplates{
			Name:  "{prefix}{fields}",
			Roots: map[string]string{"CreateOptions": "CreateField"},
		},
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

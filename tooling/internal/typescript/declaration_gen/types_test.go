package declarationgen

import (
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func TestRendererLowersTypeScriptOnlyShapes(t *testing.T) {
	profile := spec.LanguageProfile{
		Language: "typescript", ObjectKind: "class", InterfaceKind: "interface",
		Optional: "property", Fallible: "rejected_promise", MultipleResults: "object", Unions: "discriminated",
		Builtins: map[string]string{"string": "string", "integer": "bigint"},
		Outputs:  map[string]string{"shared": ".gen.ts", "client": ".gen.ts", "provider": ".gen.ts"},
	}
	renderer, err := NewRenderer(model.APIModule{TypesByID: map[string]spec.APIType{}}, profile)
	if err != nil {
		t.Fatal(err)
	}
	value, err := renderer.Type(spec.TypeRef{Union: []spec.TypeRef{{Builtin: "string"}, {Tuple: []spec.TypeRef{{Builtin: "integer"}, {Builtin: "string"}}}}})
	if err != nil {
		t.Fatal(err)
	}
	if want := "string | readonly [bigint, string]"; value != want {
		t.Fatalf("Type() = %q, want %q", value, want)
	}
}

func TestSignatureUsesSemanticIDAndPromiseFailure(t *testing.T) {
	profile := spec.LanguageProfile{
		Language: "typescript", ObjectKind: "class", InterfaceKind: "interface",
		Optional: "property", Fallible: "rejected_promise", MultipleResults: "object", Unions: "discriminated",
		Builtins: map[string]string{"string": "string"}, Outputs: map[string]string{"shared": ".gen.ts", "client": ".gen.ts", "provider": ".gen.ts"},
	}
	renderer, err := NewRenderer(model.APIModule{TypesByID: map[string]spec.APIType{}}, profile)
	if err != nil {
		t.Fatal(err)
	}
	got, err := renderer.Signature(spec.Callable{ID: "attach_sandbox", Cancellable: true, Fallible: true, Params: []spec.APISlot{{ID: "id", Name: "id", Type: spec.TypeRef{Builtin: "string"}}}, Results: []spec.APISlot{{Type: spec.TypeRef{Builtin: "string"}}}})
	if err != nil {
		t.Fatal(err)
	}
	if want := "attach(id: string, options?: OperationOptions): Promise<string>"; got != want {
		t.Fatalf("Signature() = %q, want %q", got, want)
	}
}

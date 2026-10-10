package model

import (
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"testing"
)

func TestAPIDeclarationsArePortableAndValidated(t *testing.T) {
	module := spec.APIModule{Functions: []spec.Callable{{ID: "create", Name: "create", Behavior: "create", Params: []spec.APISlot{{ID: "options", Name: "options", Type: spec.TypeRef{Union: []spec.TypeRef{{Builtin: "string"}, {Builtin: "integer"}}}}}, Fallible: true}}}
	// Unions are valid semantic types even when a particular emitter rejects them.
	if _, err := CompileAPI(spec.API{Version: 1, Modules: map[string]spec.APIModule{"test": module}}); err != nil {
		t.Fatal(err)
	}
	module.Functions[0].Params[0].Type = spec.TypeRef{Ref: "Config", Builtin: "string"}
	if _, err := CompileAPI(spec.API{Version: 1, Modules: map[string]spec.APIModule{"test": module}}); err == nil {
		t.Fatal("ambiguous shape accepted")
	}
	module.Functions[0].Params[0].Type = spec.TypeRef{Builtin: "string"}
	module.Functions[0].Name = "create; execute()"
	if _, err := CompileAPI(spec.API{Version: 1, Modules: map[string]spec.APIModule{"test": module}}); err == nil {
		t.Fatal("code accepted as a declaration name")
	}
}

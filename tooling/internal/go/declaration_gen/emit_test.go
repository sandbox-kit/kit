package declarationgen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func fixture(t *testing.T, module spec.APIModule) (*Emitter, *protogen.Plugin) {
	t.Helper()
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"test.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("test.proto"), Package: proto.String("test"), Syntax: proto.String("proto3"), Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/test;test")}}}})
	if err != nil {
		t.Fatal(err)
	}
	api, err := model.CompileAPI(spec.API{Version: 1, Modules: map[string]spec.APIModule{"test": module}})
	if err != nil {
		t.Fatal(err)
	}
	profile, err := spec.LoadLanguageProfile("../../../../specs", "go")
	if err != nil {
		t.Fatal(err)
	}
	g := p.NewGeneratedFile("test.gen.go", p.Files[0].GoImportPath)
	g.P("package test")
	e, err := New(g, p, api.Modules["test"], profile, nil)
	if err != nil {
		t.Fatal(err)
	}
	return e, p
}
func TestGoPresenceAndContainers(t *testing.T) {
	e, _ := fixture(t, spec.APIModule{Types: []spec.APIType{{ID: "contract", Name: "Contract", Kind: "interface"}}})
	for _, tc := range []struct {
		ref  spec.TypeRef
		want string
	}{
		{spec.TypeRef{Builtin: "integer", Optional: true}, "*int64"},
		{spec.TypeRef{Builtin: "boolean", Nullable: true}, "*bool"},
		{spec.TypeRef{Ref: "contract", Reference: true}, "Contract"},
		{spec.TypeRef{Sequence: &spec.TypeRef{Builtin: "string"}, Optional: true}, "[]string"},
		{spec.TypeRef{Map: &spec.MapType{Key: spec.TypeRef{Builtin: "string"}, Value: spec.TypeRef{Builtin: "integer"}}}, "map[string]int64"},
	} {
		got, err := e.Type(tc.ref)
		if err != nil || got != tc.want {
			t.Fatalf("wanted %s got %s (%v)", tc.want, got, err)
		}
	}
	for _, ref := range []spec.TypeRef{
		{Union: []spec.TypeRef{{Builtin: "string"}, {Builtin: "integer"}}},
		{Builtin: "integer", Optional: true, Nullable: true},
		{Ref: "unknown"},
	} {
		if _, err := e.Type(ref); err == nil {
			t.Fatalf("unsupported representation accepted: %+v", ref)
		}
	}
}
func TestGenericAndMultipleResultDeclarationsTypeCheck(t *testing.T) {
	functions := []spec.Callable{
		{ID: "identity", Name: "Identity", Behavior: "test", Generics: []spec.APISlot{{ID: "T", Name: "T", Type: spec.TypeRef{Builtin: "any"}}}, Params: []spec.APISlot{{ID: "value", Name: "value", Type: spec.TypeRef{Ref: "T"}}}, Results: []spec.APISlot{{Type: spec.TypeRef{Ref: "T"}}}},
		{ID: "pair", Name: "Pair", Behavior: "test", Results: []spec.APISlot{{Type: spec.TypeRef{Builtin: "string"}}, {Type: spec.TypeRef{Builtin: "integer"}}}, Fallible: true},
	}
	e, p := fixture(t, spec.APIModule{Functions: functions})
	c, err := e.BeginBehavior("", "identity", "test")
	if err != nil {
		t.Fatal(err)
	}
	e.Body(c, "return ", e.Param(c, "value"))
	e.G.P("}")
	c, err = e.BeginBehavior("", "pair", "test")
	if err != nil {
		t.Fatal(err)
	}
	e.Body(c, "return \"ok\",1,nil")
	e.G.P("}")
	response := p.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}
	source := response.File[0].GetContent()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.gen.go", source, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := (&types.Config{}).Check("example.com/test", fset, []*ast.File{file}, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(source, "func Pair() (string, int64, error)") {
		t.Fatal(source)
	}
}
func TestBodyNamesFollowDeclarationsWithoutChangingStrings(t *testing.T) {
	typ := spec.APIType{ID: "client", Name: "Client", Kind: "object", Fields: []spec.APISlot{{ID: "backend", Name: "integration", Type: spec.TypeRef{Builtin: "any"}}}}
	method := spec.Callable{ID: "create", Name: "Create", Behavior: "test", Receiver: "self", Cancellable: true, Params: []spec.APISlot{{ID: "request", Name: "options", Type: spec.TypeRef{Builtin: "any"}}}}
	typ.Methods = []spec.Callable{method}
	e, p := fixture(t, spec.APIModule{Types: []spec.APIType{typ}})
	e.Profile.Cancellation.Name = "callContext"
	if err := e.TypeDeclaration("client"); err != nil {
		t.Fatal(err)
	}
	c, err := e.Begin("client", "create")
	if err != nil {
		t.Fatal(err)
	}
	e.Body(c, "_ = ", e.Receiver(c), ".", e.Field("client", "backend"), ";_ = ", e.Param(c, "request"), ";_ = ", e.Context(c), ";_ = \"request c.backend ctx\"")
	e.G.P("}")
	source := p.Response().File[0].GetContent()
	for _, want := range []string{"self.integration", "_ = options", "_ = callContext", `"request c.backend ctx"`} {
		if !strings.Contains(source, want) {
			t.Fatalf("missing %s: %s", want, source)
		}
	}
}

func TestAliasesEnumsCallbacksAndGenericObjects(t *testing.T) {
	module := spec.APIModule{Types: []spec.APIType{
		{ID: "alias", Name: "Labels", Kind: "alias", Underlying: &spec.TypeRef{Map: &spec.MapType{Key: spec.TypeRef{Builtin: "string"}, Value: spec.TypeRef{Builtin: "string"}}}},
		{ID: "status", Name: "Status", Kind: "enum", Underlying: &spec.TypeRef{Builtin: "integer"}, Values: []spec.EnumValue{{Name: "StatusUnknown", Number: 0}, {Name: "StatusReady", Number: 7}}},
		{ID: "callback", Name: "Callback", Kind: "callback", Signature: &spec.FunctionType{Params: []spec.APISlot{{ID: "value", Name: "value", Type: spec.TypeRef{Builtin: "integer"}}}, Results: []spec.APISlot{{Type: spec.TypeRef{Builtin: "boolean"}}}}},
		{ID: "box", Name: "Box", Kind: "object", Generics: []spec.APISlot{{ID: "T", Name: "T", Type: spec.TypeRef{Builtin: "any"}}}, Fields: []spec.APISlot{{ID: "value", Name: "value", Type: spec.TypeRef{Ref: "T"}}}, Methods: []spec.Callable{{ID: "get", Name: "Get", Behavior: "get", Receiver: "b", Results: []spec.APISlot{{Type: spec.TypeRef{Ref: "T"}}}}}},
	}}
	e, p := fixture(t, module)
	for _, typ := range module.Types {
		if err := e.TypeDeclaration(typ.ID); err != nil {
			t.Fatal(err)
		}
	}
	method, err := e.BeginBehavior("box", "get", "get")
	if err != nil {
		t.Fatal(err)
	}
	e.Body(method, "return b.value")
	e.G.P("}")
	if got, err := e.Type(spec.TypeRef{Ref: "box", Arguments: []spec.TypeRef{{Builtin: "string"}}, Reference: true}); err != nil || got != "*Box[string]" {
		t.Fatalf("generic instantiation: %s, %v", got, err)
	}
	response := p.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}
	source := response.File[0].GetContent()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "types.gen.go", source, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := (&types.Config{}).Check("example.com/test", fset, []*ast.File{file}, nil); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"type Labels = map[string]string", "StatusReady   Status = 7", "type Callback func(int64) bool", "type Box[T any] struct", "func (b *Box[T]) Get() T"} {
		if !strings.Contains(source, want) {
			t.Fatalf("missing %s: %s", want, source)
		}
	}
}

func TestVariadicAndStaticDeclarationsTypeCheck(t *testing.T) {
	count := spec.Callable{ID: "count", Name: "Count", Behavior: "count", Params: []spec.APISlot{{ID: "values", Name: "values", Variadic: true, Type: spec.TypeRef{Sequence: &spec.TypeRef{Builtin: "string"}}}}, Results: []spec.APISlot{{Type: spec.TypeRef{Builtin: "integer"}}}}
	static := spec.Callable{ID: "empty", Name: "Empty", Static: true, Behavior: "empty", Results: []spec.APISlot{{Type: spec.TypeRef{Builtin: "boolean"}}}}
	e, p := fixture(t, spec.APIModule{Types: []spec.APIType{{ID: "bag", Name: "Bag", Kind: "object", Methods: []spec.Callable{static}}}, Functions: []spec.Callable{count}})
	if err := e.TypeDeclaration("bag"); err != nil {
		t.Fatal(err)
	}
	c, err := e.BeginBehavior("", "count", "count")
	if err != nil {
		t.Fatal(err)
	}
	e.Body(c, "return int64(len(values))")
	e.G.P("}")
	c, err = e.BeginBehavior("bag", "empty", "empty")
	if err != nil {
		t.Fatal(err)
	}
	e.Body(c, "return true")
	e.G.P("}")
	response := p.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}
	source := response.File[0].GetContent()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", source, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := (&types.Config{}).Check("example.com/test", fset, []*ast.File{file}, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(source, "func BagEmpty() bool") || !strings.Contains(source, "values ...string") {
		t.Fatal(source)
	}
}

func TestExplicitReferencesPreserveShadowedLocals(t *testing.T) {
	function := spec.Callable{ID: "calculate", Name: "Calculate", Behavior: "calculate", Params: []spec.APISlot{{ID: "request", Name: "options", Type: spec.TypeRef{Builtin: "integer"}}}, Results: []spec.APISlot{{Type: spec.TypeRef{Builtin: "integer"}}}}
	e, p := fixture(t, spec.APIModule{Functions: []spec.Callable{function}})
	c, err := e.BeginBehavior("", "calculate", "calculate")
	if err != nil {
		t.Fatal(err)
	}
	e.Body(c, "outer:=", e.Param(c, "request"), ";{", e.Local("request"), ":=int64(7);outer+=", e.Local("request"), "};return outer")
	e.G.P("}")
	response := p.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}
	source := response.File[0].GetContent()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", source, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := (&types.Config{}).Check("example.com/test", fset, []*ast.File{file}, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(source, "outer := options") || !strings.Contains(source, "request := int64(7)") {
		t.Fatal(source)
	}
}

func TestMissingExplicitSymbolFailsGeneration(t *testing.T) {
	e, p := fixture(t, spec.APIModule{})
	e.Param(spec.Callable{ID: "test"}, "missing")
	if p.Response().GetError() == "" {
		t.Fatal("missing symbol accepted")
	}
}

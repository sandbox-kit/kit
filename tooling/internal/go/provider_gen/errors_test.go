package providergen

import (
	"github.com/sandbox-kit/kit/tooling/internal/go/schema"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"go/parser"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"strings"
	"testing"
)

func errorFixture(t *testing.T) *protogen.Plugin {
	t.Helper()
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"errors.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("errors.proto"), Syntax: proto.String("proto3"), Package: proto.String("kit.sandbox.v1"), Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/provider;provider")}, EnumType: []*descriptorpb.EnumDescriptorProto{{Name: proto.String("ErrorKind"), Value: []*descriptorpb.EnumValueDescriptorProto{{Name: proto.String("ERROR_KIND_UNKNOWN"), Number: proto.Int32(0)}, {Name: proto.String("ERROR_KIND_AUTHENTICATION"), Number: proto.Int32(1)}}}}, MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("ErrorInfo"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("kind"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(), TypeName: proto.String(".kit.sandbox.v1.ErrorKind"), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}}}}}}})
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestErrorSpecsControlNativeClassification(t *testing.T) {
	p := errorFixture(t)
	s := spec.Provider{Provider: "test", ErrorStatuses: map[int]string{401: "authentication"}, Errors: map[string]spec.ErrorBinding{"go": {Target: spec.Binding{Import: "example.com/native", Name: "NativeError", Fields: map[string]string{"status": "Status", "code": "Code", "source": "Source"}}}}}
	if err := emitErrorFixture(p, p.Files[0], s); err != nil {
		t.Fatal(err)
	}
	response := p.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}
	src := response.File[0].GetContent()
	for _, want := range []string{"case 401:", "sandbox.ErrorKindAuthentication", "errors.As(err, &native)", "Message: err.Error()", "ProviderCode: providerCode", "sandbox.NewError(sandbox.ErrorInfo{"} {
		if !strings.Contains(src, want) {
			t.Fatalf("missing %s", want)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "errors.gen.go", src, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	s.ErrorStatuses = map[int]string{401: "missing_kind"}
	p = errorFixture(t)
	if err := emitErrorFixture(p, p.Files[0], s); err == nil {
		t.Fatal("unknown contract kind accepted")
	}
	s.ErrorStatuses = map[int]string{999: "authentication"}
	p = errorFixture(t)
	if err := emitErrorFixture(p, p.Files[0], s); err == nil {
		t.Fatal("invalid HTTP status accepted")
	}
}

func TestNativeErrorBindingUsesPortableRule(t *testing.T) {
	s := spec.Provider{Provider: "test", ErrorRules: map[string]string{"credential_rejected": "authentication"}, Errors: map[string]spec.ErrorBinding{"go": {Types: []spec.ErrorType{{Import: "example.com/native", Name: "AuthError", Rule: "credential_rejected", Pointer: true}}}}}
	p := errorFixture(t)
	if err := emitErrorFixture(p, p.Files[0], s); err != nil {
		t.Fatal(err)
	}
	src := p.Response().File[0].GetContent()
	if !strings.Contains(src, "native.AuthError") || !strings.Contains(src, "sandbox.ErrorKindAuthentication") {
		t.Fatal("portable rule not bound to native type")
	}
	delete(s.ErrorRules, "credential_rejected")
	p = errorFixture(t)
	if err := emitErrorFixture(p, p.Files[0], s); err == nil {
		t.Fatal("missing portable rule accepted")
	}
}

func emitErrorFixture(p *protogen.Plugin, f *protogen.File, s spec.Provider) error {
	compiled, err := model.CompileProvider(schema.Index(p), s)
	if err != nil {
		return err
	}
	return GenerateErrors(p, f, compiled)
}

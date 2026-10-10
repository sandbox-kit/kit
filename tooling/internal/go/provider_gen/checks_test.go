package providergen

import (
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"github.com/sandbox-kit/kit/tooling/internal/testutil"
	"go/parser"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"strings"
	"testing"
)

func checksFixture(t *testing.T) *protogen.Plugin {
	t.Helper()
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"checks.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("checks.proto"), Syntax: proto.String("proto2"), Package: proto.String("kit.sandbox.v1"), Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/provider;provider")}, MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("CreateOptions"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("runtime"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".kit.sandbox.v1.RuntimeConfig"), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}}}, {Name: proto.String("RuntimeConfig"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("language"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}}}}}}})
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestProviderChecksPreservePresenceAndAllowedValues(t *testing.T) {
	templates, profile := testutil.Generation(t, "../../../../specs")
	p := checksFixture(t)
	if err := GenerateChecks(p, p.Files[0], spec.Provider{Provider: "test", Checks: []spec.Check{{Path: "runtime.language", Allowed: []string{"python", "javascript"}}}}, templates, profile); err != nil {
		t.Fatal(err)
	}
	r := p.Response()
	if r.GetError() != "" {
		t.Fatal(r.GetError())
	}
	source := r.File[0].GetContent()
	for _, want := range []string{"request.Runtime != nil", "request.Runtime.Language != nil", `*request.Runtime.Language != "python" && *request.Runtime.Language != "javascript"`} {
		if !strings.Contains(source, want) {
			t.Fatalf("missing %s", want)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "checks.gen.go", source, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	for _, check := range []spec.Check{{Path: "runtime.missing", Allowed: []string{"x"}}, {Path: "runtime.language", Format: "unknown"}, {Path: "runtime.language", Minimum: proto.Float64(1)}, {Path: "runtime.language"}} {
		p := checksFixture(t)
		if err := GenerateChecks(p, p.Files[0], spec.Provider{Checks: []spec.Check{check}}, templates, profile); err == nil {
			t.Fatal("invalid check accepted")
		}
	}
}

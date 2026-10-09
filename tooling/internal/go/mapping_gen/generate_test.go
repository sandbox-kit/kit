package mappinggen

import (
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func fixture(t *testing.T) *protogen.Plugin {
	t.Helper()
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"sandbox.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("sandbox.proto"), Syntax: proto.String("proto2"), Package: proto.String("kit.sandbox.v1"), Options: &descriptorpb.FileOptions{GoPackage: proto.String(sandbox + ";sandbox")}, MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("Resources"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("memory_mib"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_UINT64.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}}}}}}})
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func mapping() spec.Provider {
	return spec.Provider{Version: 1, Provider: "test", Groups: []spec.Group{{Name: "mapResources", Direction: "request", Source: spec.Type{Name: "Resources"}, Target: spec.Type{Name: "NativeResources", Bindings: map[string]spec.Binding{"go": {Import: "example.com/sdk", Name: "Resources", Fields: map[string]string{"memory": "Memory"}}}}, Fields: []spec.Mapping{{From: "memory_mib", To: "memory", Transform: "divide_exactly", Factor: 1024, Cast: "integer"}}}}}
}
func TestGeneratesPresenceAndPrecisionChecks(t *testing.T) {
	p := fixture(t)
	if err := Generate(p, p.Files[0], mapping()); err != nil {
		t.Fatal(err)
	}
	r := p.Response()
	if r.GetError() != "" {
		t.Fatal(r.GetError())
	}
	s := r.File[0].GetContent()
	for _, want := range []string{"source.MemoryMiB != nil", "value%1024 != 0", "int(value / 1024)", "remaining.MemoryMiB = nil"} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %q", want)
		}
	}
}
func TestRejectsInvalidMappings(t *testing.T) {
	for _, edit := range []func(*spec.Provider){
		func(s *spec.Provider) { s.Groups[0].Fields[0].From = "misspelled" },
		func(s *spec.Provider) { s.Groups[0].Fields[0].To = "unbound" },
		func(s *spec.Provider) { s.Groups[0].Fields[0].Factor = 0 },
		func(s *spec.Provider) { s.Groups[0].Fields[0].Transform = "go_expression" },
		func(s *spec.Provider) {
			s.Groups[0].Fields[0].Transform = "policy_minutes"
			s.Groups[0].Fields[0].Policy = &spec.PolicyMapping{Disabled: "zero"}
		},
		func(s *spec.Provider) { s.Groups[0].Fields[0].Policy = &spec.PolicyMapping{Disabled: "zero"} },
		func(s *spec.Provider) { s.Groups[0].Fields = append(s.Groups[0].Fields, s.Groups[0].Fields[0]) },
	} {
		p := fixture(t)
		s := mapping()
		edit(&s)
		if err := Generate(p, p.Files[0], s); err == nil {
			t.Fatal("accepted invalid mapping")
		}
	}
}

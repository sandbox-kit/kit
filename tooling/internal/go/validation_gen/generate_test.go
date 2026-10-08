package validationgen

import (
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func testPlugin(t *testing.T) *protogen.Plugin {
	t.Helper()
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"sandbox.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("sandbox.proto"), Syntax: proto.String("proto2"), Package: proto.String("kit.sandbox.v1"), Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/sandbox;sandbox")}, MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("Resources"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("cpu_cores"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}, {Name: proto.String("cpu_limit_cores"), Number: proto.Int32(2), Type: descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}}}}}}})
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestGeneratesLimitRule(t *testing.T) {
	p := testPlugin(t)
	rules := spec.Validation{Version: 1, Messages: map[string]spec.Message{"Resources": {Constraints: []spec.Constraint{{Op: "limit", Field: "cpu_limit_cores", Other: "cpu_cores"}}}}}
	if err := Generate(p, p.Files[0], rules, "CreateOptions"); err != nil {
		t.Fatal(err)
	}
	response := p.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}
	s := response.File[0].GetContent()
	if !strings.Contains(s, "x.GetCPULimitCores() < x.GetCPUCores()") {
		t.Fatal("limit rule was not generated")
	}
}
func TestRejectsUnknownRulesAndFields(t *testing.T) {
	for _, rule := range []spec.Message{
		{Fields: map[string]spec.Field{"unknown": {Finite: true}}},
		{Fields: map[string]spec.Field{"cpu_cores": {Each: true}}},
		{Constraints: []spec.Constraint{{Op: "at_most_one", Fields: []string{"unknown"}}}},
		{Constraints: []spec.Constraint{{Op: "raw_go_code"}}},
	} {
		p := testPlugin(t)
		if err := Generate(p, p.Files[0], spec.Validation{Messages: map[string]spec.Message{"Resources": rule}}, "CreateOptions"); err == nil {
			t.Fatal("invalid validation rule accepted")
		}
	}
}

func TestOptionalNumericTags(t *testing.T) {
	zero := 0.0
	if tags := Tags(spec.Field{Finite: true, ExclusiveMinimum: &zero}, true, false); tags != "omitnil,finite,gt=0" {
		t.Fatalf("unexpected tags %s", tags)
	}
}

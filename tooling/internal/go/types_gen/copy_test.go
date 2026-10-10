package typesgen

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"strings"
	"testing"
)

func TestCopyEmitterUsesSchemaPresenceAndDepth(t *testing.T) {
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"copy.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("copy.proto"), Syntax: proto.String("proto2"), Package: proto.String("kit.sandbox.v1"), Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/sandbox;sandbox")}, MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("Resources"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("memory_mib"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_UINT64.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}}}}}}})
	if err != nil {
		t.Fatal(err)
	}
	if err := GenerateCopies(p, p.Files[0], 7); err != nil {
		t.Fatal(err)
	}
	r := p.Response()
	if r.GetError() != "" {
		t.Fatal(r.GetError())
	}
	src := r.File[0].GetContent()
	for _, want := range []string{"depth > 7", "x.MemoryMiB != nil", "value := *x.MemoryMiB", "out.MemoryMiB = &value"} {
		if !strings.Contains(src, want) {
			t.Fatalf("missing %s", want)
		}
	}
	if strings.Contains(src, "encoding/json") {
		t.Fatal("copy still serializes")
	}
}

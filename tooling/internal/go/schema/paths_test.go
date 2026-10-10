package schema

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"testing"
)

func TestCanonicalPathResolvesAgainstSchemaIdentity(t *testing.T) {
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"path.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("path.proto"), Syntax: proto.String("proto2"), Package: proto.String("kit.sandbox.v1"), Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/sandbox;sandbox")}, MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("CreateOptions"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("cpu_cores"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_DOUBLE.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}}}}}}})
	if err != nil {
		t.Fatal(err)
	}
	fields, err := Resolve(p, "CreateOptions", "cpu_cores")
	if err != nil || Constant("CreateOptions", fields) != "CreateFieldCPUCores" {
		t.Fatal("canonical path resolution failed", err)
	}
	for _, path := range []string{"cpu_core", "cpu_cores.child"} {
		if _, err := Resolve(p, "CreateOptions", path); err == nil {
			t.Fatalf("invalid path %s accepted", path)
		}
	}
}

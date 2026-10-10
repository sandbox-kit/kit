package errorgen

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"strings"
	"testing"
)

func TestSharedExpressionPreservesCauseAndDetails(t *testing.T) {
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"example.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("example.proto"), Syntax: proto.String("proto3"), Package: proto.String("example"), Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/provider;provider")}}}})
	if err != nil {
		t.Fatal(err)
	}
	g := p.NewGeneratedFile("example.go", "example.com/provider")
	expression := Expression(g, Details{Kind: Kind(g, "ErrorKindAuthentication"), Provider: `"modal"`, Message: "native.Error()", Cause: "native", StatusCode: "code"})
	for _, want := range []string{"sandbox.NewError(sandbox.ErrorInfo{", "Kind:sandbox.ErrorKindAuthentication", "Provider:\"modal\"", "StatusCode:code", "},native)"} {
		if !strings.Contains(expression, want) {
			t.Fatalf("missing %s in %s", want, expression)
		}
	}
}

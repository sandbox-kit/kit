package clientgen

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	codegenv1 "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func generateForTest(t *testing.T, declaration *codegenv1.ClientDeclaration) (string, error) {
	t.Helper()
	options := &descriptorpb.FileOptions{
		GoPackage: proto.String("github.com/sandbox-kit/kit/core;core"),
	}
	if declaration != nil {
		proto.SetExtension(options, codegenv1.E_Client, declaration)
	}
	plugin, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"client.proto"},
		ProtoFile: []*descriptorpb.FileDescriptorProto{{
			Name: proto.String("client.proto"), Syntax: proto.String("proto3"),
			Package: proto.String("kit.core.v1"), Options: options,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := Generate(plugin, plugin.Files[0]); err != nil {
		return "", err
	}
	response := plugin.Response()
	if response.GetError() != "" {
		t.Fatalf("generator response: %s", response.GetError())
	}
	if len(response.File) == 0 {
		return "", nil
	}
	return response.File[0].GetContent(), nil
}

func TestGeneratesConfiguredNamesAndValidGo(t *testing.T) {
	source, err := generateForTest(t, &codegenv1.ClientDeclaration{
		GoName: "Session", GoProviderInterface: "Backend",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "client.kit.go", source, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"type Session struct", "func NewSession(provider Backend)", "ProviderName() string"} {
		if !strings.Contains(source, expected) {
			t.Fatalf("missing configured declaration: %s", expected)
		}
	}
	if strings.Contains(source, "func (c *Session) Provider()") || strings.Contains(source, "[P ") {
		t.Fatal("generated client still exposes provider-specific typing")
	}
}

func TestRejectsInvalidOrCollidingNames(t *testing.T) {
	for _, name := range []string{"", "client", "type", "Client;panic()", "Provider"} {
		t.Run(name, func(t *testing.T) {
			_, err := generateForTest(t, &codegenv1.ClientDeclaration{
				GoName: name, GoProviderInterface: "Provider",
			})
			if err == nil {
				t.Fatal("invalid declaration was accepted")
			}
		})
	}
}

func TestSkipsFilesWithoutClientAnnotation(t *testing.T) {
	source, err := generateForTest(t, nil)
	if err != nil || source != "" {
		t.Fatalf("unexpected generation: %q, %v", source, err)
	}
}

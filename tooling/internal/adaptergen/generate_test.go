package adaptergen

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

func validSpec() *codegenv1.AdapterDeclaration {
	return &codegenv1.AdapterDeclaration{
		ProviderName: "custom", GoName: "Bridge", GoConstructor: "Attach",
	}
}

func generateForTest(t *testing.T, spec *codegenv1.AdapterDeclaration, withClient bool) (string, error) {
	t.Helper()
	options := &descriptorpb.FileOptions{GoPackage: proto.String("example.com/kit/custom;custom")}
	if spec != nil {
		proto.SetExtension(options, codegenv1.E_Adapter, spec)
	}
	if withClient {
		proto.SetExtension(options, codegenv1.E_Client, &codegenv1.ClientDeclaration{})
	}
	plugin, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"adapter.proto"},
		ProtoFile: []*descriptorpb.FileDescriptorProto{{
			Name: proto.String("adapter.proto"), Syntax: proto.String("proto3"),
			Package: proto.String("kit.custom"), Options: options,
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
		t.Fatal(response.GetError())
	}
	if len(response.File) == 0 {
		return "", nil
	}
	return response.File[0].GetContent(), nil
}

func TestGeneratesSDKIndependentAdapterAndConfiguredNames(t *testing.T) {
	source, err := generateForTest(t, validSpec(), false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "adapter.kit.go", source, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"type Bridge[C any] struct", "func Attach[C any](client *C)",
		"client: client", `return "custom"`,
		"core.Provider",
	} {
		if !strings.Contains(source, expected) {
			t.Fatalf("missing generated binding: %s", expected)
		}
	}
	if strings.Contains(source, "SDK()") || strings.Contains(source, "Native()") {
		t.Fatal("adapter exposes its internal SDK client")
	}
	if strings.Contains(source, "example.com/native/sdk") || strings.Contains(source, "modal-client") {
		t.Fatal("adapter imports an official SDK")
	}
}

func TestRejectsInvalidDeclarations(t *testing.T) {
	for _, change := range []struct {
		name string
		edit func(*codegenv1.AdapterDeclaration)
	}{
		{"unexported name", func(s *codegenv1.AdapterDeclaration) { s.GoName = "adapter" }},
		{"invalid constructor", func(s *codegenv1.AdapterDeclaration) { s.GoConstructor = "New()" }},
		{"duplicate declaration", func(s *codegenv1.AdapterDeclaration) { s.GoConstructor = s.GoName }},
		{"missing provider", func(s *codegenv1.AdapterDeclaration) { s.ProviderName = "" }},
	} {
		t.Run(change.name, func(t *testing.T) {
			spec := validSpec()
			change.edit(spec)
			if _, err := generateForTest(t, spec, false); err == nil {
				t.Fatal("invalid specification accepted")
			}
		})
	}
	if _, err := generateForTest(t, validSpec(), true); err == nil {
		t.Fatal("client and adapter declared in the same file")
	}
}

func TestSkipsUnannotatedFile(t *testing.T) {
	if source, err := generateForTest(t, nil, false); err != nil || source != "" {
		t.Fatalf("unexpected output: %q, %v", source, err)
	}
}

package providergen

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	codegenv1 "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func validSpec() *codegenv1.ProviderDeclaration {
	return &codegenv1.ProviderDeclaration{
		ProviderName: "custom", Name: "Bridge",
		Sdks: map[string]*codegenv1.TargetSDK{"go": {ImportPath: "example.com/native/sdk", ClientType: "Client"}},
	}
}

func generateForTest(t *testing.T, declaration *codegenv1.ProviderDeclaration, withClient bool) (string, error) {
	return generateRuntimeForTest(t, declaration, withClient, spec.RuntimeBinding{Cleanup: spec.Cleanup{Method: "Close", ReturnsError: true}})
}

func generateRuntimeForTest(t *testing.T, declaration *codegenv1.ProviderDeclaration, withClient bool, runtime spec.RuntimeBinding) (string, error) {
	t.Helper()
	options := &descriptorpb.FileOptions{GoPackage: proto.String("example.com/kit/custom;custom")}
	if declaration != nil {
		proto.SetExtension(options, codegenv1.E_Provider, declaration)
	}
	if withClient {
		proto.SetExtension(options, codegenv1.E_Client, &codegenv1.ClientDeclaration{})
	}
	plugin, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"provider.proto"},
		ProtoFile: []*descriptorpb.FileDescriptorProto{{
			Name: proto.String("config.proto"), Syntax: proto.String("proto2"), Package: proto.String("kit.sandbox.v1"),
			Options:     &descriptorpb.FileOptions{GoPackage: proto.String("github.com/sandbox-kit/kit/sdks/go/sandbox;sandbox")},
			MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("Config"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("region"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}}}},
		}, {
			Name: proto.String("provider.proto"), Syntax: proto.String("proto3"),
			Dependency: []string{"config.proto"},
			Package:    proto.String("kit.custom"), Options: options,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := Generate(plugin, plugin.Files[1], runtime); err != nil {
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

func TestConstructorBindingsAndStateCapture(t *testing.T) {
	runtime := spec.RuntimeBinding{
		Constructor: spec.Constructor{Function: "NewConfiguredClient"},
		Cleanup:     spec.Cleanup{Method: "Close", ReturnsError: true},
		State:       []spec.StateField{{Name: "region", Type: "string", Optional: true, Capture: "region"}},
	}
	source, err := generateRuntimeForTest(t, validSpec(), false, runtime)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"sdk.NewConfiguredClient(&params)", "result.region = &value", "config.Region != nil", "config = &sandbox.Config{}", "return nil, err"} {
		if !strings.Contains(source, want) {
			t.Fatalf("missing %q", want)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "provider.gen.go", source, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	runtime.State[0].Capture = "unknown"
	if _, err := generateRuntimeForTest(t, validSpec(), false, runtime); err == nil {
		t.Fatal("unknown capture accepted")
	}
	runtime.State = nil
	runtime.Constructor.Function = "NewClient()"
	if _, err := generateRuntimeForTest(t, validSpec(), false, runtime); err == nil {
		t.Fatal("constructor expression accepted")
	}
}

func TestGeneratesProviderBackendAndConfiguredNames(t *testing.T) {
	source, err := generateForTest(t, validSpec(), false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "provider.gen.go", source, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"type bridge struct", "func New() *Provider",
		"NewClient(config *sandbox.Config)", `return "custom"`, "a.create(ctx, request)",
		"sandbox.Provider",
	} {
		if !strings.Contains(source, expected) {
			t.Fatalf("missing generated binding: %s", expected)
		}
	}
	if strings.Contains(source, "SDK()") || strings.Contains(source, "Native()") {
		t.Fatal("backend exposes its internal SDK client")
	}
	if !strings.Contains(source, "example.com/native/sdk") {
		t.Fatal("provider did not import its typed SDK binding")
	}
}

func TestRejectsInvalidDeclarations(t *testing.T) {
	for _, change := range []struct {
		name string
		edit func(*codegenv1.ProviderDeclaration)
	}{
		{"unexported name", func(s *codegenv1.ProviderDeclaration) { s.Name = "adapter" }},
		{"invalid name", func(s *codegenv1.ProviderDeclaration) { s.Name = "backend()" }},
		{"duplicate declaration", func(s *codegenv1.ProviderDeclaration) { s.Name = "New" }},
		{"missing provider", func(s *codegenv1.ProviderDeclaration) { s.ProviderName = "" }},
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
		t.Fatal("client and provider declared in the same file")
	}
}

func TestSkipsUnannotatedFile(t *testing.T) {
	if source, err := generateForTest(t, nil, false); err != nil || source != "" {
		t.Fatalf("unexpected output: %q, %v", source, err)
	}
}

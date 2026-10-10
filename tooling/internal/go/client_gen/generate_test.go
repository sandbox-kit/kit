package clientgen

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	codegenv1 "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func generateForTest(t *testing.T, declaration *codegenv1.ClientDeclaration, edit func(*spec.API)) (string, error) {
	t.Helper()
	options := &descriptorpb.FileOptions{
		GoPackage: proto.String("github.com/sandbox-kit/kit/sdks/go/sandbox;sandbox"),
	}
	if declaration != nil {
		declaration.CreationService = "SandboxCreation"
		declaration.ConfigMessage = "Config"
		proto.SetExtension(options, codegenv1.E_Client, declaration)
	}
	plugin, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"client.proto"},
		ProtoFile: []*descriptorpb.FileDescriptorProto{{
			Name: proto.String("client.proto"), Syntax: proto.String("proto3"),
			Package: proto.String("kit.sandbox.v1"), Options: options,
			MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("Config")}, {Name: proto.String("CreateOptions")}, {Name: proto.String("CreateResult")}},
			Service:     []*descriptorpb.ServiceDescriptorProto{{Name: proto.String("SandboxCreation"), Method: []*descriptorpb.MethodDescriptorProto{{Name: proto.String("Create"), InputType: proto.String(".kit.sandbox.v1.CreateOptions"), OutputType: proto.String(".kit.sandbox.v1.CreateResult")}}}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	compiled, err := model.CompileClient(plugin.Files[0].Desc)
	if err != nil {
		return "", err
	}
	generation, err := spec.LoadGeneration("../../../../specs")
	if err != nil {
		t.Fatal(err)
	}
	compiled.Operations = generation.Operations
	raw, err := spec.LoadAPI("../../../../specs")
	if err != nil {
		t.Fatal(err)
	}
	if edit != nil {
		edit(&raw)
	}
	api, err := model.CompileAPI(raw)
	if err != nil {
		return "", err
	}
	compiled.Surface = api.Modules["client"]
	profile, err := spec.LoadLanguageProfile("../../../../specs", "go")
	if err != nil {
		t.Fatal(err)
	}
	if err := Generate(plugin, plugin.Files[0], compiled, profile); err != nil {
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
	source, err := generateForTest(t, &codegenv1.ClientDeclaration{}, func(raw *spec.API) {
		module := raw.Modules["client"]
		for i, typ := range module.Types {
			if typ.ID == "client" {
				typ.Name = "Session"
			}
			module.Types[i] = typ
		}
		for i, function := range module.Functions {
			if function.ID == "constructor" {
				function.Name = "NewSession"
			}
			module.Functions[i] = function
		}
		raw.Modules["client"] = module
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "client.gen.go", source, parser.AllErrors); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"type Session struct", "func NewSession(config Config)", "ProviderName() string"} {
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
			_, err := generateForTest(t, &codegenv1.ClientDeclaration{}, func(raw *spec.API) {
				module := raw.Modules["client"]
				for i, typ := range module.Types {
					if typ.ID == "client" {
						typ.Name = name
					}
					module.Types[i] = typ
				}
				raw.Modules["client"] = module
			})
			if err == nil {
				t.Fatal("invalid declaration was accepted")
			}
		})
	}
}

func TestSkipsFilesWithoutClientAnnotation(t *testing.T) {
	source, err := generateForTest(t, nil, nil)
	if err != nil || source != "" {
		t.Fatalf("unexpected generation: %q, %v", source, err)
	}
}

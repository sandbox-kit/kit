package typesgen

import (
	"github.com/sandbox-kit/kit/tooling/internal/testutil"
	"strings"
	"testing"

	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestPreservesEmptyCollectionsAndLimitsCreationHelpers(t *testing.T) {
	options := &descriptorpb.FileOptions{GoPackage: proto.String("example.com/sandbox;sandbox")}
	proto.SetExtension(options, codegen.E_NativeTypes, true)
	plugin, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{FileToGenerate: []string{"runtime.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{Name: proto.String("runtime.proto"), Syntax: proto.String("proto3"), Package: proto.String("kit.sandbox.v1"), Options: options, MessageType: []*descriptorpb.DescriptorProto{{Name: proto.String("RuntimeConfig"), Field: []*descriptorpb.FieldDescriptorProto{{Name: proto.String("entrypoint"), Number: proto.Int32(1), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum()}}}}}}})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := spec.LoadAPI("../../../../specs")
	if err != nil {
		t.Fatal(err)
	}
	api, err := model.CompileAPI(raw)
	if err != nil {
		t.Fatal(err)
	}
	profile, err := spec.LoadLanguageProfile("../../../../specs", "go")
	if err != nil {
		t.Fatal(err)
	}
	templates, _ := testutil.Generation(t, "../../../../specs")
	if err := Generate(plugin, plugin.Files[0], spec.Validation{}, model.APIModule{}, api.Modules["helpers"], profile, templates); err != nil {
		t.Fatal(err)
	}
	response := plugin.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}
	source := response.File[0].GetContent()
	if !strings.Contains(source, `json:"entrypoint"`) {
		t.Fatal("empty list presence would be lost during copying")
	}
	if strings.Contains(source, "func Value[") || strings.Contains(source, "type CreateOptions") {
		t.Fatal("non-creation files duplicate creation helper declarations")
	}
}

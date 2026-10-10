package model

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// No go_package or native SDK descriptors: the model uses schema identities only.
func fixture(t *testing.T) protoreflect.FileDescriptor {
	t.Helper()
	field := func(name string, number int32, kind descriptorpb.FieldDescriptorProto_Type, typeName string) *descriptorpb.FieldDescriptorProto {
		f := &descriptorpb.FieldDescriptorProto{Name: proto.String(name), Number: proto.Int32(number), Type: kind.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}
		if typeName != "" {
			f.TypeName = proto.String(typeName)
		}
		return f
	}
	file, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{Name: proto.String("shared.proto"), Package: proto.String("kit.sandbox.v1"), Syntax: proto.String("proto2"),
		EnumType: []*descriptorpb.EnumDescriptorProto{{Name: proto.String("ErrorKind"), Value: []*descriptorpb.EnumValueDescriptorProto{{Name: proto.String("ERROR_KIND_UNKNOWN"), Number: proto.Int32(0)}, {Name: proto.String("ERROR_KIND_AUTHENTICATION"), Number: proto.Int32(3)}}}},
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: proto.String("ErrorInfo"), Field: []*descriptorpb.FieldDescriptorProto{field("kind", 1, descriptorpb.FieldDescriptorProto_TYPE_ENUM, ".kit.sandbox.v1.ErrorKind")}},
			{Name: proto.String("CreateOptions"), Field: []*descriptorpb.FieldDescriptorProto{field("resources", 1, descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, ".kit.sandbox.v1.Resources")}},
			{Name: proto.String("Resources"), Field: []*descriptorpb.FieldDescriptorProto{field("cpu_cores", 1, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, "")}},
			{Name: proto.String("Config"), Field: []*descriptorpb.FieldDescriptorProto{field("endpoint", 1, descriptorpb.FieldDescriptorProto_TYPE_STRING, "")}},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func TestPortablePathIdentityAndPresence(t *testing.T) {
	s := NewSchema([]protoreflect.FileDescriptor{fixture(t)})
	path, err := s.Resolve("CreateOptions", "resources.cpu_cores")
	if err != nil {
		t.Fatal(err)
	}
	if path.Root != "kit.sandbox.v1.CreateOptions" || len(path.Fields) != 2 || !path.Fields[1].HasPresence() {
		t.Fatalf("identity/presence lost: %+v", path)
	}
	for _, path := range []string{"resources.cpu", "resources.cpu_cores.child"} {
		if _, err := s.Resolve("CreateOptions", path); err == nil {
			t.Fatalf("invalid path %s accepted", path)
		}
	}
	if _, err := s.Message("other.Resources"); err == nil {
		t.Fatal("short-name collision accepted")
	}
}

func TestPortableProviderCompilation(t *testing.T) {
	s := NewSchema([]protoreflect.FileDescriptor{fixture(t)})
	rules := spec.Provider{Provider: "test", ErrorRules: map[string]string{"credentials": "authentication"}, ErrorStatuses: map[int]string{403: "authentication", 401: "authentication"},
		Errors: map[string]spec.ErrorBinding{"go": {Types: []spec.ErrorType{{Name: "GoCredentialError", Rule: "credentials"}}}, "typescript": {Types: []spec.ErrorType{{Name: "JavaScriptCredentialError", Rule: "credentials"}}}},
		Groups: []spec.Group{{Name: "resources", Direction: "request", Source: spec.Type{Name: "Resources"}, Target: spec.Type{Name: "NativeResources", Bindings: map[string]spec.Binding{"go": {Name: "Resources"}}}, Fields: []spec.Mapping{{From: "cpu_cores", To: "cpu", Transform: "whole"}}}},
	}
	compiled, err := CompileProvider(s, rules)
	if err != nil {
		t.Fatal(err)
	}
	if compiled.Errors.HTTP[0].Code != 401 || compiled.Errors.Native[0].Kind.Number != 3 {
		t.Fatal("classification identity or ordering lost")
	}
	mapping := compiled.Groups[0].Fields[0]
	if mapping.ErrorPath != "resources.cpu_cores" || mapping.Shared.Fields[0].Kind() != protoreflect.DoubleKind {
		t.Fatal("mapping was not resolved")
	}
	// Changing a language binding cannot alter the semantic plan.
	rules.Errors["typescript"] = spec.ErrorBinding{Types: []spec.ErrorType{{Name: "DifferentNativeName", Rule: "credentials"}}}
	other, err := CompileProvider(s, rules)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(compiled.Errors, other.Errors) {
		t.Fatal("language binding changed semantic classification")
	}
	rules.Groups[0].Fields[0].Transform = "arbitrary_code"
	if _, err := CompileProvider(s, rules); err == nil {
		t.Fatal("unknown operation accepted")
	}
}

func TestCompileResolvesSpecsBeforeEmission(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"generation.yaml": "version: 1\nshared: [types, error_runtime, copies, diagnostic_paths, field_paths, validation]\nclient: [types, diagnostic_paths, field_paths, validation, client]\nprovider: [provider, configuration, checks, errors, mappings]\n",
		"validation.yaml": "version: 1\nmessages:\n  Resources:\n    fields:\n      cpu_cores: {finite: true}\n",
		"client.yaml":     "version: 1\nconfig_message: Config\nvalidation:\n  version: 1\n  messages:\n    Config:\n      fields:\n        endpoint: {nonblank: true}\n",
		"contracts.yaml":  "version: 1\ncopy_max_depth: 64\ncopy_max_nodes: 16384\ncopy_max_bytes: 1048576\nmetadata_message: MetadataObject\nerror_message: ErrorInfo\norigin_enum: ValueOrigin\n",
		"errors.yaml":     "version: 1\nfallback_kind: unknown\ncause_rules: []\nexisting_error: {reuse_matching_context: true, preserve_details: true}\nclassification_precedence: [existing_error, native_type, http_status, grpc_code, fallback]\n",
	}
	generation, readErr := os.ReadFile(filepath.Join("..", "..", "..", "specs", "generation.yaml"))
	if readErr != nil {
		t.Fatal(readErr)
	}
	files["generation.yaml"] = string(generation)
	api, err := os.ReadFile(filepath.Join("..", "..", "..", "specs", "api.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files["api.yaml"] = string(api)
	templates, err := os.ReadFile(filepath.Join("..", "..", "..", "specs", "templates.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files["templates.yaml"] = string(templates)
	behaviors, err := os.ReadFile(filepath.Join("..", "..", "..", "specs", "behaviors.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	files["behaviors.yaml"] = string(behaviors)
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := Compile(root, []protoreflect.FileDescriptor{fixture(t)}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "validation.yaml"), []byte("version: 1\nmessages:\n  Resources:\n    fields:\n      misspelled_cpu: {finite: true}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Compile(root, []protoreflect.FileDescriptor{fixture(t)}); err == nil {
		t.Fatal("bad spec deferred to emitter")
	}
}

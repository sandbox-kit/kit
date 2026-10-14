package providergen

import (
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestValidateConfigPlanRequiresRetainedStateCapture(t *testing.T) {
	t.Parallel()
	schema := configTestSchema(t)
	provider := spec.Provider{
		Provider: "example",
		Client:   &spec.ClientMapping{Retained: []string{"region"}},
		Runtime:  map[string]spec.RuntimeBinding{"typescript": {}},
	}

	err := validateConfigPlan(provider, schema)
	if err == nil || !strings.Contains(err.Error(), "region has no TypeScript state capture") {
		t.Fatalf("validateConfigPlan() error = %v", err)
	}
}

func TestValidateConfigPlanAcceptsTypedStateAndNestedRetention(t *testing.T) {
	t.Parallel()
	schema := configTestSchema(t)
	provider := spec.Provider{
		Provider: "example",
		Client: &spec.ClientMapping{
			Retained: []string{"region"},
			Scope:    &spec.ClientComponent{Field: "scope", Retained: []string{"app_name"}},
		},
		Runtime: map[string]spec.RuntimeBinding{"typescript": {State: []spec.StateField{
			{Name: "scope", Type: "Scope", Optional: true, Capture: "scope"},
			{Name: "region", Type: "string", Optional: true, Capture: "region"},
		}}},
	}

	if err := validateConfigPlan(provider, schema); err != nil {
		t.Fatalf("validateConfigPlan() error = %v", err)
	}
}

func TestValidateConfigPlanRejectsStateTypeMismatch(t *testing.T) {
	t.Parallel()
	provider := spec.Provider{
		Provider: "example",
		Client:   &spec.ClientMapping{},
		Runtime: map[string]spec.RuntimeBinding{"typescript": {State: []spec.StateField{
			{Name: "region", Type: "Scope", Optional: true, Capture: "region"},
		}}},
	}

	err := validateConfigPlan(provider, configTestSchema(t))
	if err == nil || !strings.Contains(err.Error(), "disagrees with Config.region") {
		t.Fatalf("validateConfigPlan() error = %v", err)
	}
}

func configTestSchema(t *testing.T) *model.Schema {
	t.Helper()
	optional := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	stringType := descriptorpb.FieldDescriptorProto_TYPE_STRING
	messageType := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
	file, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:    proto.String("kit/sandbox/v1/config_test.proto"),
		Package: proto.String("kit.sandbox.v1"),
		Syntax:  proto.String("proto2"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: proto.String("Scope"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: proto.String("app_name"), Number: proto.Int32(1), Label: &optional, Type: &stringType},
			}},
			{Name: proto.String("Config"), Field: []*descriptorpb.FieldDescriptorProto{
				{Name: proto.String("scope"), Number: proto.Int32(1), Label: &optional, Type: &messageType, TypeName: proto.String(".kit.sandbox.v1.Scope")},
				{Name: proto.String("region"), Number: proto.Int32(2), Label: &optional, Type: &stringType},
			}},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return model.NewSchema([]protoreflect.FileDescriptor{file})
}

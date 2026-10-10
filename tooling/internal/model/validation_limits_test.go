package model

import (
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/reflect/protoreflect"
	"testing"
)

func TestPortableSecurityFormats(t *testing.T) {
	schema := NewSchema([]protoreflect.FileDescriptor{fixture(t)})
	for _, format := range []string{"secure_endpoint", "image_reference"} {
		rules := spec.Validation{Messages: map[string]spec.Message{"Config": {Fields: map[string]spec.Field{"endpoint": {Format: format}}}}}
		if err := schema.ValidateRules(rules); err != nil {
			t.Fatal(err)
		}
		rules.Messages["Config"].Fields["endpoint"] = spec.Field{Format: "unrecognized"}
		if err := schema.ValidateRules(rules); err == nil {
			t.Fatal("unknown format accepted")
		}
		rules.Messages = map[string]spec.Message{"Resources": {Fields: map[string]spec.Field{"cpu_cores": {Format: format}}}}
		if err := schema.ValidateRules(rules); err == nil {
			t.Fatal("format on numeric field accepted")
		}
	}
}

func TestPortableClientEnvironmentRules(t *testing.T) {
	schema := NewSchema([]protoreflect.FileDescriptor{fixture(t)})
	for _, rule := range []spec.ClientEnvironment{
		{Field: "unknown", Variables: []string{"ENDPOINT"}, Default: "https://example.test"},
		{Field: "endpoint", Variables: []string{"BAD=NAME"}, Default: "https://example.test"},
		{Field: "endpoint", Variables: []string{"ENDPOINT"}},
	} {
		_, err := CompileProvider(schema, spec.Provider{Client: &spec.ClientMapping{Environment: []spec.ClientEnvironment{rule}}})
		if err == nil {
			t.Fatal("invalid resolution accepted")
		}
	}
}

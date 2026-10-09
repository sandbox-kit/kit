package providergen

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func configFixture(t *testing.T) *protogen.Plugin {
	t.Helper()
	p, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"client.proto"}, ProtoFile: []*descriptorpb.FileDescriptorProto{{
			Name: proto.String("client.proto"), Syntax: proto.String("proto2"), Package: proto.String("kit.sandbox.v1"),
			Options:     &descriptorpb.FileOptions{GoPackage: proto.String("example.com/provider;provider")},
			MessageType: configMessages(),
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func configMessages() []*descriptorpb.DescriptorProto {
	stringField := func(name string, n int32) *descriptorpb.FieldDescriptorProto {
		return &descriptorpb.FieldDescriptorProto{Name: proto.String(name), Number: proto.Int32(n), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}
	}
	messageField := func(name, typ string, n int32) *descriptorpb.FieldDescriptorProto {
		return &descriptorpb.FieldDescriptorProto{Name: proto.String(name), Number: proto.Int32(n), Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".kit.sandbox.v1." + typ), Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum()}
	}
	return []*descriptorpb.DescriptorProto{
		{Name: proto.String("Config"), Field: []*descriptorpb.FieldDescriptorProto{stringField("endpoint", 1), messageField("scope", "Scope", 2), messageField("auth", "AuthConfig", 3)}},
		{Name: proto.String("Scope"), Field: []*descriptorpb.FieldDescriptorProto{stringField("environment", 1)}},
		{Name: proto.String("AuthConfig"), Field: []*descriptorpb.FieldDescriptorProto{messageField("api_key", "APIKeyCredentials", 1), messageField("bearer_token", "BearerTokenCredentials", 2)}},
		{Name: proto.String("APIKeyCredentials"), Field: []*descriptorpb.FieldDescriptorProto{stringField("key", 1)}},
		{Name: proto.String("BearerTokenCredentials"), Field: []*descriptorpb.FieldDescriptorProto{stringField("token", 1)}},
	}
}

func TestCompositionRejectsConflictsAndNativeTypeMismatch(t *testing.T) {
	fixture := func() spec.Provider {
		native := spec.Type{Name: "Params", Bindings: map[string]spec.Binding{"go": {Import: "example.com/sdk", Name: "Params", Fields: map[string]string{"endpoint": "Endpoint", "environment": "Environment", "key": "Key"}}}}
		return spec.Provider{Provider: "custom", Client: &spec.ClientMapping{Target: native, Settings: "mapSettings", Scope: &spec.ClientComponent{Field: "scope", Group: "mapScope"}}, Groups: []spec.Group{
			{Name: "mapSettings", Direction: "request", Source: spec.Type{Name: "Config"}, Target: native, Fields: []spec.Mapping{{From: "endpoint", To: "endpoint"}}},
			{Name: "mapScope", Direction: "request", Source: spec.Type{Name: "Scope"}, Target: native, Fields: []spec.Mapping{{From: "environment", To: "environment"}}},
		}}
	}
	for _, edit := range []func(*spec.Provider){
		func(s *spec.Provider) { s.Groups[1].Fields[0].To = "endpoint" },
		func(s *spec.Provider) {
			s.Groups[1].Target = spec.Type{Name: "Params", Bindings: map[string]spec.Binding{"go": {Import: "other/sdk", Name: "Params"}}}
		},
		func(s *spec.Provider) { s.Client.Managed = []string{"endpoint"} },
		func(s *spec.Provider) { s.Client.Retained = []string{"endpoint"} },
		func(s *spec.Provider) { s.Client.Rejected = map[string]string{"endpoint": "no"} },
		func(s *spec.Provider) { s.Client.Target.Bindings["go"] = spec.Binding{Import: "sdk", Name: "Params()"} },
		func(s *spec.Provider) {
			s.Groups[0].Source.Bindings = map[string]spec.Binding{"go": {Import: "other/sdk", Name: "Config"}}
		},
		func(s *spec.Provider) { s.Groups[0].Fields[0].From = "missing" },
		func(s *spec.Provider) { s.Client.Scope.Retained = []string{"environment"} },
	} {
		s := fixture()
		edit(&s)
		p := configFixture(t)
		if err := validateComposition(p, s); err == nil {
			t.Fatal("invalid composition accepted")
		}
	}
	s := fixture()
	p := configFixture(t)
	if err := validateComposition(p, s); err != nil {
		t.Fatal(err)
	}
	// Exclusive authentication branches may use the same native destination.
	s.Client.Auth = []spec.ClientComponent{{Field: "api_key", Group: "mapKey"}, {Field: "bearer_token", Group: "mapToken"}}
	s.Groups = append(s.Groups, spec.Group{Name: "mapKey", Direction: "request", Source: spec.Type{Name: "APIKeyCredentials"}, Target: s.Client.Target, Fields: []spec.Mapping{{From: "key", To: "key"}}}, spec.Group{Name: "mapToken", Direction: "request", Source: spec.Type{Name: "BearerTokenCredentials"}, Target: s.Client.Target, Fields: []spec.Mapping{{From: "token", To: "key"}}})
	if err := validateComposition(p, s); err != nil {
		t.Fatal(err)
	}
	s.Groups[2].Fields[0].To = "endpoint"
	if err := validateComposition(p, s); err == nil {
		t.Fatal("auth overwrote settings")
	}
	s.Groups[2].Fields[0].To = "key"
	s.Client.Auth[0].Destination = "credentials"
	b := s.Client.Target.Bindings["go"]
	b.Fields["credentials"] = "Credentials"
	b.Objects = map[string]spec.Binding{"credentials": {Import: "other/sdk", Name: "Params"}}
	s.Client.Target.Bindings["go"] = b
	if err := validateComposition(p, s); err == nil {
		t.Fatal("incompatible attached object accepted")
	}
}

func TestClientSpecControlsRejectedFields(t *testing.T) {
	for _, message := range []string{"endpoint unavailable", "use default endpoint"} {
		p := configFixture(t)
		s := spec.Provider{Provider: "custom", Client: &spec.ClientMapping{Target: spec.Type{Name: "Params", Bindings: map[string]spec.Binding{"go": {Import: "example.com/sdk", Name: "Params"}}}, Rejected: map[string]string{"endpoint": message}}}
		if err := GenerateConfig(p, p.Files[0], s); err != nil {
			t.Fatal(err)
		}
		r := p.Response()
		if r.GetError() != "" {
			t.Fatal(r.GetError())
		}
		source := r.File[0].GetContent()
		if !strings.Contains(source, message) || !strings.Contains(source, "config.Endpoint != nil") || !strings.Contains(source, "remaining.Endpoint = nil") {
			t.Fatal("spec did not control generated rejection")
		}
		if _, err := parser.ParseFile(token.NewFileSet(), "client.gen.go", source, parser.AllErrors); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRejectsUnknownClientAssemblyBindings(t *testing.T) {
	for _, edit := range []func(*spec.ClientMapping){
		func(c *spec.ClientMapping) { c.Managed = []string{"misspelled"} },
		func(c *spec.ClientMapping) { c.Settings = "missingMapping" },
		func(c *spec.ClientMapping) { c.Rejected = map[string]string{"unknown": "unsupported"} },
		func(c *spec.ClientMapping) {
			c.Auth = []spec.ClientComponent{{Field: "missing", Group: "missingMapping"}}
		},
	} {
		p := configFixture(t)
		c := &spec.ClientMapping{Target: spec.Type{Name: "Params", Bindings: map[string]spec.Binding{"go": {Import: "example.com/sdk", Name: "Params"}}}}
		edit(c)
		if err := GenerateConfig(p, p.Files[0], spec.Provider{Provider: "custom", Client: c}); err == nil {
			t.Fatal("invalid binding accepted")
		}
	}
}

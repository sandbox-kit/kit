package declarationgen

import (
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"github.com/sandbox-kit/kit/tooling/internal/testutil"
	"strings"
	"testing"
)

func TestDescriptorTemplateNamesParamsAndTypes(t *testing.T) {
	e, p := fixture(t, spec.APIModule{})
	templates, profile := testutil.Generation(t, "../../../../specs")
	pattern := templates.Callables["request_mapping"]
	pattern.Declaration.Params[0].Name = "input"
	templates.Callables["request_mapping"] = pattern
	emitter := TemplateEmitter{G: e.G, Plugin: p, Templates: templates, Profile: profile}
	d, c, err := emitter.Begin("request_mapping", TemplateBindings{Names: map[string]string{"mapping": "convertResources"}, Types: map[string]spec.TypeRef{"source": {Ref: "native_source"}, "target": {Ref: "native_target"}}, Externals: map[string]spec.ExternalType{"native_source": {Name: "Source"}, "native_target": {Name: "Target"}}})
	if err != nil {
		t.Fatal(err)
	}
	d.Body(c, "remaining:=*", d.Param(c, "source"), ";return Target{},remaining,nil")
	e.G.P("}")
	source := p.Response().File[0].GetContent()
	for _, want := range []string{"func convertResources(input *Source) (Target, Source, error)", "remaining := *input"} {
		if !strings.Contains(source, want) {
			t.Fatalf("missing %s: %s", want, source)
		}
	}
	if _, _, err := emitter.Begin("getter", TemplateBindings{Names: map[string]string{"subject": "Resources", "field": "CPU"}}); err == nil {
		t.Fatal("missing projected return type accepted")
	}
}

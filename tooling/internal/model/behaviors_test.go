package model

import (
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func behaviorFixture(t *testing.T) (Behaviors, spec.Callable) {
	t.Helper()
	raw, err := spec.LoadBehaviors("../../../specs")
	if err != nil {
		t.Fatal(err)
	}
	contracts, err := CompileBehaviors(raw)
	if err != nil {
		t.Fatal(err)
	}
	return contracts, spec.Callable{
		ID: "request_mapping", Name: "mapCreate", Behavior: "request_mapping", Fallible: true,
		Params:  []spec.APISlot{{ID: "source", Name: "source", Type: spec.TypeRef{Ref: "source", Reference: true, Nullable: true}}},
		Results: []spec.APISlot{{Type: spec.TypeRef{Ref: "target"}}, {Type: spec.TypeRef{Ref: "source"}}},
	}
}

func TestBehaviorContractsPermitNativeNames(t *testing.T) {
	b, c := behaviorFixture(t)
	c.Name = "MapRequest"
	c.Params = append([]spec.APISlot{}, c.Params...)
	c.Params[0].Name = "options"
	if err := b.Validate(c, false); err != nil {
		t.Fatal(err)
	}
}

func TestBehaviorContractsRejectIncompatibleDeclarations(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*spec.Callable)
	}{
		{"missing request", func(c *spec.Callable) { c.Params = nil }},
		{"wrong type", func(c *spec.Callable) { c.Params[0].Type = spec.TypeRef{Builtin: "string"} }},
		{"wrong ID", func(c *spec.Callable) { c.Params[0].ID = "unrelated" }},
		{"missing failure", func(c *spec.Callable) { c.Fallible = false }},
		{"extra parameter", func(c *spec.Callable) {
			c.Params = append(c.Params, spec.APISlot{ID: "extra", Name: "extra", Type: spec.TypeRef{Builtin: "string"}})
		}},
		{"unknown behavior", func(c *spec.Callable) { c.Behavior = "arbitrary_execution" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, c := behaviorFixture(t)
			tc.edit(&c)
			if err := b.Validate(c, false); err == nil {
				t.Fatal("incompatible behavior accepted")
			}
		})
	}
}

package model

import (
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"path/filepath"
	"testing"
)

func behaviorFixture(t *testing.T) (Behaviors, spec.Callable) {
	t.Helper()
	root := filepath.Join("..", "..", "..", "specs")
	raw, err := spec.LoadBehaviors(root)
	if err != nil {
		t.Fatal(err)
	}
	contracts, err := CompileBehaviors(raw)
	if err != nil {
		t.Fatal(err)
	}
	api, err := spec.LoadAPI(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, typ := range api.Modules["client"].Types {
		if typ.ID == "client" {
			for _, method := range typ.Methods {
				if method.ID == "create" {
					return contracts, method
				}
			}
		}
	}
	t.Fatal("missing creation declaration")
	return contracts, spec.Callable{}
}
func TestBehaviorContractsPermitNativeNames(t *testing.T) {
	b, c := behaviorFixture(t)
	c.Name = "Launch"
	c.Receiver = "self"
	c.Params = append([]spec.APISlot{}, c.Params...)
	c.Params[0].Name = "options"
	if err := b.Validate(c, true); err != nil {
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
		{"missing cancellation", func(c *spec.Callable) { c.Cancellable = false }},
		{"wrong output", func(c *spec.Callable) { c.Results[0].Type = spec.TypeRef{Builtin: "boolean"} }},
		{"extra parameter", func(c *spec.Callable) {
			c.Params = append(c.Params, spec.APISlot{ID: "extra", Name: "extra", Type: spec.TypeRef{Builtin: "string"}})
		}},
		{"unknown behavior", func(c *spec.Callable) { c.Behavior = "arbitrary_execution" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, c := behaviorFixture(t)
			tc.edit(&c)
			if err := b.Validate(c, true); err == nil {
				t.Fatal("incompatible behavior accepted")
			}
		})
	}
}

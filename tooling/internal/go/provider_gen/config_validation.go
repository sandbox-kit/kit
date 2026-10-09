package providergen

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
)

func validBinding(b spec.Binding) bool {
	return b.Import != "" && token.IsIdentifier(b.Name) && !token.Lookup(b.Name).IsKeyword()
}
func sameNative(a, b spec.Type) bool {
	x, y := a.Bindings["go"], b.Bindings["go"]
	return validBinding(x) && validBinding(y) && x.Import == y.Import && x.Name == y.Name
}
func validateComposition(p *protogen.Plugin, provider spec.Provider) error {
	c := provider.Client
	if c == nil {
		return nil
	}
	if !validBinding(c.Target.Bindings["go"]) {
		return fmt.Errorf("invalid client native target binding")
	}
	groups := map[string]spec.Group{}
	for _, g := range provider.Groups {
		if !token.IsIdentifier(g.Name) || token.Lookup(g.Name).IsKeyword() {
			return fmt.Errorf("invalid group identifier %q", g.Name)
		}
		if _, ok := groups[g.Name]; ok {
			return fmt.Errorf("duplicate group %q", g.Name)
		}
		groups[g.Name] = g
	}
	group := func(name string) (spec.Group, error) {
		g, ok := groups[name]
		if !ok || g.Direction != "request" {
			return g, fmt.Errorf("unknown request group %q", name)
		}
		if len(g.When) > 0 {
			return g, fmt.Errorf("conditional groups cannot assemble client config")
		}
		if _, native := g.Source.Bindings["go"]; native {
			return g, fmt.Errorf("client mapping source must be a shared schema type")
		}
		if !validBinding(g.Target.Bindings["go"]) {
			return g, fmt.Errorf("invalid native target for %s", name)
		}
		return g, nil
	}
	owners := map[string]string{}
	own := func(path, owner string) error {
		if prev, ok := owners[path]; ok {
			return fmt.Errorf("field %s owned by both %s and %s", path, prev, owner)
		}
		owners[path] = owner
		return nil
	}
	retained := func(path string) error {
		for _, state := range provider.Runtime["go"].State {
			if state.Capture == path || strings.HasPrefix(path, state.Capture+".") && state.Capture != "" {
				return nil
			}
		}
		return fmt.Errorf("retained field %s has no backend capture", path)
	}
	for _, f := range c.Managed {
		if f != "provider" && f != "timeout" {
			return fmt.Errorf("field %s is not Kit-managed", f)
		}
		if err := own(f, "Kit"); err != nil {
			return err
		}
	}
	for _, f := range c.Retained {
		if _, err := sharedField(p, "Config", f); err != nil {
			return err
		}
		if err := retained(f); err != nil {
			return err
		}
		if err := own(f, "retained"); err != nil {
			return err
		}
	}
	for f := range c.Rejected {
		if _, err := sharedField(p, "Config", f); err != nil {
			return err
		}
		if err := own(f, "rejected"); err != nil {
			return err
		}
	}
	base := map[string]bool{}
	add := func(g spec.Group, dest map[string]bool) error {
		for _, m := range g.Fields {
			member := g.Target.Bindings["go"].Fields[m.To]
			if !token.IsIdentifier(member) || token.Lookup(member).IsKeyword() {
				return fmt.Errorf("invalid destination member %q", member)
			}
			if dest[member] {
				return fmt.Errorf("conflicting native destination %s", member)
			}
			dest[member] = true
		}
		return nil
	}
	if c.Settings != "" {
		g, err := group(c.Settings)
		if err != nil {
			return err
		}
		m, err := sharedMessage(p, g.Source.Name)
		if err != nil {
			return err
		}
		if string(m.Desc.FullName()) != "kit.sandbox.v1.Config" || !sameNative(g.Target, c.Target) {
			return fmt.Errorf("incompatible settings native/schema identity")
		}
		for _, f := range g.Fields {
			if _, err := sharedField(p, "Config", f.From); err != nil {
				return err
			}
			if err := own(f.From, "settings"); err != nil {
				return err
			}
		}
		if err := add(g, base); err != nil {
			return err
		}
	}
	component := func(comp spec.ClientComponent, parent string, dest map[string]bool) error {
		f, err := sharedField(p, parent, comp.Field)
		if err != nil {
			return err
		}
		g, err := group(comp.Group)
		if err != nil {
			return err
		}
		m, err := sharedMessage(p, g.Source.Name)
		if err != nil {
			return err
		}
		if f.Message == nil || f.Message.Desc.FullName() != m.Desc.FullName() {
			return fmt.Errorf("incompatible component schema identity")
		}
		used := map[string]bool{}
		for _, mapping := range g.Fields {
			if _, err := sharedField(p, g.Source.Name, mapping.From); err != nil {
				return err
			}
			used[mapping.From] = true
		}
		for _, field := range comp.Retained {
			if used[field] {
				return fmt.Errorf("field %s is both mapped and retained", field)
			}
			if _, err := sharedField(p, g.Source.Name, field); err != nil {
				return err
			}
			if err := retained(comp.Field + "." + field); err != nil {
				return err
			}
			used[field] = true
		}
		if comp.Destination != "" {
			expected := c.Target.Bindings["go"].Objects[comp.Destination]
			actual := g.Target.Bindings["go"]
			if !validBinding(expected) || expected.Import != actual.Import || expected.Name != actual.Name {
				return fmt.Errorf("incompatible object destination type %s", comp.Destination)
			}
			member := c.Target.Bindings["go"].Fields[comp.Destination]
			if !token.IsIdentifier(member) || token.Lookup(member).IsKeyword() {
				return fmt.Errorf("invalid destination binding")
			}
			if dest[member] {
				return fmt.Errorf("conflicting native destination %s", member)
			}
			dest[member] = true
		} else {
			if !sameNative(g.Target, c.Target) {
				return fmt.Errorf("incompatible component native identity")
			}
			if err := add(g, dest); err != nil {
				return err
			}
		}
		return nil
	}
	if c.Scope != nil {
		if c.Scope.Field != "scope" {
			return fmt.Errorf("scope component must select scope")
		}
		if err := own("scope", "scope component"); err != nil {
			return err
		}
		if err := component(*c.Scope, "Config", base); err != nil {
			return err
		}
	}
	// Authentication branches are exclusive; each must be compatible with the
	// common settings/scope destinations, but may share destinations with each other.
	if err := own("auth", "auth component"); err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, a := range c.Auth {
		if seen[a.Field] {
			return fmt.Errorf("duplicate auth variant")
		}
		seen[a.Field] = true
		branch := map[string]bool{}
		for k, v := range base {
			branch[k] = v
		}
		if err := component(a, "AuthConfig", branch); err != nil {
			return err
		}
	}
	return nil
}

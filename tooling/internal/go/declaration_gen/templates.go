package declarationgen

import (
	"fmt"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
)

// TemplateBindings retains type shapes and native symbol identities separately.
// Types contain no emitted Go syntax; externals are selected language bindings.
type TemplateBindings struct {
	Names     map[string]string
	Types     map[string]spec.TypeRef
	Externals map[string]spec.ExternalType
}
type TemplateEmitter struct {
	G         *protogen.GeneratedFile
	Plugin    *protogen.Plugin
	Templates model.Templates
	Profile   spec.LanguageProfile
}

func (t TemplateEmitter) Bind(id string, b TemplateBindings) (*Emitter, spec.DeclarationTemplate, error) {
	pattern, err := t.Templates.Instantiate(id, b.Names)
	if err != nil {
		return nil, pattern, err
	}
	profile := t.Profile
	profile.Externals = map[string]spec.ExternalType{}
	for id, value := range t.Profile.Externals {
		profile.Externals[id] = value
	}
	for id, value := range b.Externals {
		if _, exists := profile.Externals[id]; exists {
			return nil, pattern, fmt.Errorf("external binding collision %s", id)
		}
		profile.Externals[id] = value
	}
	module := model.APIModule{TypesByID: map[string]spec.APIType{}, FunctionsByID: map[string]spec.Callable{}}
	if pattern.Owner != "" {
		name := b.Names[pattern.Owner]
		if name == "" {
			return nil, pattern, fmt.Errorf("missing receiver binding %s", pattern.Owner)
		}
		module.TypesByID[pattern.Owner] = spec.APIType{ID: pattern.Owner, Name: name, Kind: "object", Methods: []spec.Callable{pattern.Declaration}}
	} else {
		module.FunctionsByID[pattern.Declaration.ID] = pattern.Declaration
	}
	d, err := New(t.G, t.Plugin, module, profile, nil)
	if err != nil {
		return nil, pattern, err
	}
	d.TypeBindings = b.Types
	if _, err := d.Signature(pattern.Owner, pattern.Declaration, false); err != nil {
		return nil, pattern, err
	}
	return d, pattern, nil
}
func (t TemplateEmitter) Begin(id string, b TemplateBindings) (*Emitter, spec.Callable, error) {
	d, pattern, err := t.Bind(id, b)
	if err != nil {
		return nil, pattern.Declaration, err
	}
	c, err := d.BeginBehavior(pattern.Owner, pattern.Declaration.ID, id)
	return d, c, err
}

// NamedBinding describes an SDK symbol without formatting or parsing its type.
func NamedBinding(id protogen.GoIdent) spec.ExternalType {
	return spec.ExternalType{Import: string(id.GoImportPath), Name: id.GoName}
}

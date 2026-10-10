// Package declarationgen lowers portable declarations into idiomatic Go source.
package declarationgen

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"unicode"
	"unicode/utf8"
)

type Emitter struct {
	G             *protogen.GeneratedFile
	Plugin        *protogen.Plugin
	Module        model.APIModule
	Profile       spec.LanguageProfile
	Native        map[string]protogen.GoIdent
	TypeBindings  map[string]spec.TypeRef
	TypeVariables map[string]string
}

func New(g *protogen.GeneratedFile, p *protogen.Plugin, module model.APIModule, profile spec.LanguageProfile, native map[string]protogen.GoIdent) (*Emitter, error) {
	if err := ValidateProfile(profile); err != nil {
		return nil, err
	}
	e := &Emitter{G: g, Plugin: p, Module: module, Profile: profile, Native: native}
	for _, typ := range module.Types {
		local, err := e.genericScope(typ.Generics)
		if err != nil {
			return nil, err
		}
		first, _ := utf8.DecodeRuneInString(typ.Name)
		if typ.Visibility == "public" && !unicode.IsUpper(first) {
			return nil, fmt.Errorf("public Go type must be exported: %s", typ.Name)
		}
		if err := goIdentifier(typ.Name); err != nil {
			return nil, err
		}
		for _, field := range typ.Fields {
			if err := goIdentifier(field.Name); err != nil {
				return nil, err
			}
			if _, err := local.Type(field.Type); err != nil {
				return nil, err
			}
		}
		for _, method := range typ.Methods {
			if _, err := e.Signature(typ.ID, method, typ.Kind == "interface"); err != nil {
				return nil, err
			}
		}
	}
	for _, function := range module.Functions {
		if _, err := e.Signature("", function, false); err != nil {
			return nil, err
		}
	}
	return e, nil
}
func goIdentifier(name string) error {
	if !token.IsIdentifier(name) {
		return fmt.Errorf("invalid Go identifier %q", name)
	}
	return nil
}

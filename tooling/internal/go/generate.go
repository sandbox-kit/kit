// Package goemit dispatches compiled generation phases to native Go emitters.
package goemit

import (
	"fmt"

	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	mappinggen "github.com/sandbox-kit/kit/tooling/internal/go/mapping_gen"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// Generate consumes a language-neutral program. Only this layer chooses Go
// bindings, native naming, imports, and output file conventions.
func Generate(plugin *protogen.Plugin, file *protogen.File, program *model.Program) error {
	profile, err := spec.LoadLanguageProfile("specs", "go")
	if err != nil {
		return err
	}
	if err := declarationgen.ValidateProfile(profile); err != nil {
		return err
	}
	var phases []string
	var provider model.Provider
	switch {
	case file.Desc.Path() == "kit/sandbox/v1/sandbox.proto":
		// Shared create validation is handwritten; no shared emission phases remain.
		return nil
	case proto.HasExtension(file.Desc.Options(), codegen.E_Client):
		// Client path/error helpers are handwritten; no emission phases remain.
		return nil
	case proto.HasExtension(file.Desc.Options(), codegen.E_Provider):
		declaration := proto.GetExtension(file.Desc.Options(), codegen.E_Provider).(*codegen.ProviderDeclaration)
		provider = program.Providers[declaration.ProviderName]
		phases = program.Generation.Provider
	default:
		return fmt.Errorf("%s: unsupported Go generation target", file.Desc.Path())
	}
	for _, phase := range phases {
		var err error
		switch phase {
		case "mappings":
			err = mappinggen.Generate(plugin, file, provider, program.Templates, profile)
		default:
			return fmt.Errorf("unsupported Go emission phase %s", phase)
		}
		if err != nil {
			return fmt.Errorf("%s: %s: %w", file.Desc.Path(), phase, err)
		}
	}
	return nil
}

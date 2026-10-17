// Package goemit dispatches compiled generation phases to native Go emitters.
package goemit

import (
	"fmt"
	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	mappinggen "github.com/sandbox-kit/kit/tooling/internal/go/mapping_gen"
	providergen "github.com/sandbox-kit/kit/tooling/internal/go/provider_gen"
	typesgen "github.com/sandbox-kit/kit/tooling/internal/go/types_gen"
	validationgen "github.com/sandbox-kit/kit/tooling/internal/go/validation_gen"
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
	rules := spec.Validation{}
	representation := spec.RepresentationPolicy{}
	root := "CreateOptions"
	var provider model.Provider
	switch {
	case file.Desc.Path() == "kit/sandbox/v1/sandbox.proto":
		phases = program.Generation.Shared
		rules = program.Validation
	case proto.HasExtension(file.Desc.Options(), codegen.E_Client):
		phases = program.Generation.Client
		rules = program.Client.Validation
		representation = program.Client.Representation
		root = program.Client.ConfigMessage
	case proto.HasExtension(file.Desc.Options(), codegen.E_Provider):
		declaration := proto.GetExtension(file.Desc.Options(), codegen.E_Provider).(*codegen.ProviderDeclaration)
		provider = program.Providers[declaration.ProviderName]
		phases = program.Generation.Provider
	default:
		return typesgen.Generate(plugin, file, rules, representation, profile, program.Templates)
	}
	for _, phase := range phases {
		// Client data is emitted by protoc into the contracts module. The runtime
		// package consumes those messages directly; it must not create a second
		// native representation of Config, credentials, scope, or request envelopes.
		if proto.HasExtension(file.Desc.Options(), codegen.E_Client) && (phase == "types" || phase == "copies" || phase == "validation") {
			continue
		}
		var err error
		switch phase {
		case "types":
			err = typesgen.Generate(plugin, file, rules, representation, profile, program.Templates)
		case "copies":
			err = typesgen.GenerateCopies(plugin, file, program.Contracts.CopyMaxDepth, program.Templates, profile)
		case "diagnostic_paths":
			err = typesgen.GenerateErrorPaths(plugin, file)
		case "field_paths":
			err = typesgen.GeneratePaths(plugin, file, program.Templates, profile)
		case "validation":
			err = validationgen.Generate(plugin, file, rules, root, program.Templates, profile)
		case "checks":
			err = providergen.GenerateChecks(plugin, file, provider.Rules, program.Templates, profile)
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

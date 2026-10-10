// Package goemit dispatches compiled generation phases to native Go emitters.
package goemit

import (
	"fmt"
	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	clientgen "github.com/sandbox-kit/kit/tooling/internal/go/client_gen"
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
	surface := model.APIModule{}
	root := "CreateOptions"
	var provider model.Provider
	switch {
	case file.Desc.Path() == "kit/sandbox/v1/sandbox.proto":
		phases = program.Generation.Shared
		rules = program.Validation
	case proto.HasExtension(file.Desc.Options(), codegen.E_Client):
		phases = program.Generation.Client
		rules = program.Client.Validation
		surface = program.API.Modules["client"]
		root = program.Client.ConfigMessage
	case proto.HasExtension(file.Desc.Options(), codegen.E_Provider):
		declaration := proto.GetExtension(file.Desc.Options(), codegen.E_Provider).(*codegen.ProviderDeclaration)
		provider = program.Providers[declaration.ProviderName]
		phases = program.Generation.Provider
	default:
		return typesgen.Generate(plugin, file, rules, surface, program.API.Modules["helpers"], profile, program.Templates)
	}
	for _, phase := range phases {
		var err error
		switch phase {
		case "types":
			err = typesgen.Generate(plugin, file, rules, surface, program.API.Modules["helpers"], profile, program.Templates)
		case "error_runtime":
			err = typesgen.GenerateRuntime(plugin, file, program.Contracts, program.API.Modules["errors"], profile)
		case "copies":
			err = typesgen.GenerateCopies(plugin, file, program.Contracts.CopyMaxDepth, program.Templates, profile)
		case "diagnostic_paths":
			err = typesgen.GenerateErrorPaths(plugin, file)
		case "field_paths":
			err = typesgen.GeneratePaths(plugin, file, program.Templates, profile)
		case "validation":
			err = validationgen.Generate(plugin, file, rules, root, program.Templates, profile)
		case "client":
			err = clientgen.Generate(plugin, file, program.Clients[file.Desc.Path()], profile)
		case "provider":
			err = providergen.Generate(plugin, file, provider.Rules.Runtime["go"], program.API.Modules["provider"], profile)
		case "configuration":
			err = providergen.GenerateConfig(plugin, file, provider.Rules, program.Templates, profile)
		case "checks":
			err = providergen.GenerateChecks(plugin, file, provider.Rules, program.Templates, profile)
		case "errors":
			err = providergen.GenerateErrors(plugin, file, provider)
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

// Package typescriptemit dispatches the language-neutral program to TypeScript emitters.
package typescriptemit

import (
	"fmt"
	"strings"

	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	clientgen "github.com/sandbox-kit/kit/tooling/internal/typescript/client_gen"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	providergen "github.com/sandbox-kit/kit/tooling/internal/typescript/provider_gen"
	typesgen "github.com/sandbox-kit/kit/tooling/internal/typescript/types_gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// Generate emits one TypeScript source file for a declared generation target.
func Generate(plugin *protogen.Plugin, file *protogen.File, program *model.Program) error {
	profile, err := spec.LoadLanguageProfile("specs", "typescript")
	if err != nil {
		return err
	}
	if err := declarationgen.ValidateProfile(profile); err != nil {
		return err
	}
	var phases []string
	providerName := ""
	switch {
	case file.Desc.Path() == "kit/sandbox/v1/sandbox.proto":
		phases = program.Generation.Shared
	case proto.HasExtension(file.Desc.Options(), codegen.E_Client):
		phases = program.Generation.Client
	case proto.HasExtension(file.Desc.Options(), codegen.E_Provider):
		declaration := proto.GetExtension(file.Desc.Options(), codegen.E_Provider).(*codegen.ProviderDeclaration)
		if strings.TrimSpace(declaration.ProviderName) == "" {
			return fmt.Errorf("%s: provider name is required", file.Desc.Path())
		}
		providerName = declaration.ProviderName
		phases = program.Generation.Provider
	default:
		return fmt.Errorf("%s: unsupported TypeScript generation target", file.Desc.Path())
	}
	for _, phase := range phases {
		var phaseErr error
		switch phase {
		case "types":
			// Client declarations are emitted with the client runtime because both
			// share client.proto; shared protobuf declarations are emitted here.
			if file.Desc.Path() == "kit/sandbox/v1/sandbox.proto" {
				phaseErr = typesgen.Generate(plugin, file, program, profile)
			}
		case "client":
			phaseErr = clientgen.Generate(plugin, file, program, profile)
		case "provider":
			phaseErr = providergen.Generate(plugin, file, program.Providers[providerName], profile)
		case "mappings":
			switch providerName {
			case "daytona":
				phaseErr = generateDaytonaMappings(plugin, file, program.Providers[providerName])
			case "modal":
				phaseErr = generateModalMappings(plugin, file, program.Providers[providerName])
			default:
				phaseErr = fmt.Errorf("unsupported TypeScript provider %s", providerName)
			}
		case "errors":
			if providerName != "" {
				phaseErr = providergen.GenerateErrors(plugin, file, program.Providers[providerName])
			}
		case "checks":
			if providerName != "" {
				phaseErr = providergen.GenerateChecks(plugin, file, program.Providers[providerName].Rules, program.Schema)
			}
		case "configuration":
			if providerName != "" {
				phaseErr = providergen.GenerateConfig(plugin, file, program.Providers[providerName].Rules, program.Schema)
			}
		case "error_runtime", "copies", "diagnostic_paths", "field_paths", "validation", "lifecycle":
			// These concerns currently share an output unit with types, client, or
			// mappings. Keeping them in phase dispatch makes omissions explicit and
			// allows each unit to be split without changing the specification.
		default:
			return fmt.Errorf("%s: unsupported TypeScript emission phase %s", file.Desc.Path(), phase)
		}
		if phaseErr != nil {
			return fmt.Errorf("%s: %s: %w", file.Desc.Path(), phase, phaseErr)
		}
	}
	return nil
}

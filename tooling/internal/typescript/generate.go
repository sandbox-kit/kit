// Package typescriptemit dispatches the language-neutral program to TypeScript emitters.
package typescriptemit

import (
	"fmt"
	"strings"

	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
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
		// Shared Go owns validation helpers; TypeScript has no shared emission yet.
		return nil
	case proto.HasExtension(file.Desc.Options(), codegen.E_Client):
		// Client helpers are handwritten; no TypeScript emission phases remain.
		return nil
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
		case "mappings":
			switch providerName {
			case "daytona":
				phaseErr = generateDaytonaMappings(plugin, file, program.Providers[providerName])
			case "modal":
				phaseErr = generateModalMappings(plugin, file, program.Providers[providerName])
			default:
				phaseErr = fmt.Errorf("unsupported TypeScript provider %s", providerName)
			}
		default:
			return fmt.Errorf("%s: unsupported TypeScript emission phase %s", file.Desc.Path(), phase)
		}
		if phaseErr != nil {
			return fmt.Errorf("%s: %s: %w", file.Desc.Path(), phase, phaseErr)
		}
	}
	return nil
}

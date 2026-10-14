// Package providergen emits TypeScript provider declarations from the compiled
// provider model. Native SDK orchestration remains in the provider module.
package providergen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	"google.golang.org/protobuf/compiler/protogen"
)

// Generate emits the stable provider identity after validating the compiled
// provider and the TypeScript SDK binding selected by the protobuf declaration.
func Generate(plugin *protogen.Plugin, file *protogen.File, provider model.Provider, profile spec.LanguageProfile) error {
	name := strings.TrimSpace(provider.Rules.Provider)
	if name == "" {
		return fmt.Errorf("provider name is required")
	}
	if _, ok := provider.Rules.Runtime[profile.Language]; !ok {
		return fmt.Errorf("provider %s has no %s runtime binding", name, profile.Language)
	}
	g := declarationgen.GeneratedFile(plugin, "provider"+profile.Outputs["provider"], string(file.Desc.Path())+" + specs/providers/"+name+".yaml")
	g.P("export const providerName = ", fmt.Sprintf("%q", name), " as const;")
	return nil
}

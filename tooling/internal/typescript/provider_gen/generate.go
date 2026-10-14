// Package providergen emits TypeScript provider declarations from the compiled
// provider model. Native SDK orchestration remains in the provider module.
package providergen

import (
	"fmt"
	"go/token"
	"strings"

	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
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
	declaration := proto.GetExtension(file.Desc.Options(), codegen.E_Provider).(*codegen.ProviderDeclaration)
	sdk := declaration.Sdks[profile.Language]
	if sdk == nil || strings.TrimSpace(sdk.ImportPath) == "" || !token.IsIdentifier(sdk.ClientType) {
		return fmt.Errorf("provider %s has no valid %s SDK declaration", name, profile.Language)
	}
	g := declarationgen.GeneratedFile(plugin, "provider"+profile.Outputs["provider"], string(file.Desc.Path())+" + specs/providers/"+name+".yaml")
	g.P(`import type { Backend, Config, CreateOptions, CreateResult, OperationOptions, Provider } from "@sandbox-kit/sandbox";`)
	g.P()
	g.P("export const providerName = ", fmt.Sprintf("%q", name), " as const;")
	g.P()
	g.P("/** @internal Shared generated provider identity wiring. */")
	g.P("export abstract class ProviderBase implements Provider {")
	g.P("  readonly name = providerName;")
	g.P("  abstract initialize(config: Config): Backend;")
	g.P("}")
	g.P()
	g.P("/** @internal Shared generated backend identity and required operation surface. */")
	g.P("export abstract class BackendBase implements Backend {")
	g.P("  readonly name = providerName;")
	g.P("  abstract create(request: CreateOptions, options: OperationOptions): Promise<CreateResult>;")
	g.P("  abstract close(options: OperationOptions): Promise<void>;")
	g.P("}")
	return nil
}

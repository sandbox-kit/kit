// Command protoc-gen-kit-go generates local Go SDK clients from annotations.
package main

import (
	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	clientgen "github.com/sandbox-kit/kit/tooling/internal/go/client_gen"
	mappinggen "github.com/sandbox-kit/kit/tooling/internal/go/mapping_gen"
	providergen "github.com/sandbox-kit/kit/tooling/internal/go/provider_gen"
	typesgen "github.com/sandbox-kit/kit/tooling/internal/go/types_gen"
	validationgen "github.com/sandbox-kit/kit/tooling/internal/go/validation_gen"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, file := range plugin.Files {
			if file.Generate {
				if file.Desc.Path() == "kit/sandbox/v1/sandbox.proto" {
					rules, err := spec.LoadValidation("specs")
					if err != nil {
						return err
					}
					if err = typesgen.Generate(plugin, file, rules, spec.Client{}); err != nil {
						return err
					}
					contracts, err := spec.LoadContracts("specs")
					if err != nil {
						return err
					}
					if err = typesgen.GenerateRuntime(plugin, file, contracts); err != nil {
						return err
					}
					if err = typesgen.GenerateCopies(plugin, file, contracts.CopyMaxDepth); err != nil {
						return err
					}
					if err = typesgen.GenerateErrorPaths(plugin, file); err != nil {
						return err
					}
					if err = typesgen.GeneratePaths(plugin, file); err != nil {
						return err
					}
					if err = validationgen.Generate(plugin, file, rules, "CreateOptions"); err != nil {
						return err
					}
				} else if proto.HasExtension(file.Desc.Options(), codegen.E_Client) {
					client, err := spec.LoadClient("specs")
					if err != nil {
						return err
					}
					if err = typesgen.Generate(plugin, file, client.Validation, client); err != nil {
						return err
					}
					if err = typesgen.GenerateErrorPaths(plugin, file); err != nil {
						return err
					}
					if err = typesgen.GeneratePaths(plugin, file); err != nil {
						return err
					}
					if err = validationgen.Generate(plugin, file, client.Validation, client.ConfigMessage); err != nil {
						return err
					}
				} else if err := typesgen.Generate(plugin, file, spec.Validation{}, spec.Client{}); err != nil {
					return err
				}
				if err := clientgen.Generate(plugin, file); err != nil {
					return err
				}
				if proto.HasExtension(file.Desc.Options(), codegen.E_Provider) {
					declaration := proto.GetExtension(file.Desc.Options(), codegen.E_Provider).(*codegen.ProviderDeclaration)
					mappings, err := spec.LoadProvider("specs", declaration.ProviderName)
					if err != nil {
						return err
					}
					if err = providergen.Generate(plugin, file, mappings.Runtime["go"]); err != nil {
						return err
					}
					if err = providergen.GenerateConfig(plugin, file, mappings); err != nil {
						return err
					}
					if err = providergen.GenerateChecks(plugin, file, mappings); err != nil {
						return err
					}
					if _, err = spec.LoadErrorPolicy("specs"); err != nil {
						return err
					}
					if err = providergen.GenerateErrors(plugin, file, mappings); err != nil {
						return err
					}
					if err = mappinggen.Generate(plugin, file, mappings); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

package typescriptemit

import (
	"fmt"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	mappinggen "github.com/sandbox-kit/kit/tooling/internal/typescript/mapping_gen"
	"google.golang.org/protobuf/compiler/protogen"
)

func generateModalMappings(plugin *protogen.Plugin, file *protogen.File, provider model.Provider) error {
	if provider.Rules.Client == nil || provider.Rules.Client.Target.Bindings["typescript"].Import != "modal" {
		return fmt.Errorf("modal TypeScript: ModalClient binding is required")
	}
	if err := validateModalTypeScriptBindings(provider); err != nil {
		return err
	}
	return mappinggen.Generate(plugin, file, provider)
}

func validateModalTypeScriptBindings(provider model.Provider) error {
	required := map[string]mappinggen.Binding{
		"mapSandboxIdentity":      {Import: "modal", Fields: map[string]string{"id": "sandboxId"}},
		"mapClientScope":          {Import: "modal", Fields: map[string]string{"environment": "environment"}},
		"mapCommonCreateFields":   {Import: "modal", Fields: map[string]string{"name": "name", "environment": "env", "labels": "tags"}},
		"mapResources":            {Import: "modal", Fields: map[string]string{"cpu": "cpu", "cpu_limit": "cpuLimit", "memory": "memoryMiB", "memory_limit": "memoryLimitMiB"}},
		"mapLifetimeDuration":     {Import: "modal", Fields: map[string]string{"timeout": "timeoutMs"}},
		"mapIdleDuration":         {Import: "modal", Fields: map[string]string{"timeout": "idleTimeoutMs"}},
		"mapTokenPairCredentials": {Import: "modal", Fields: map[string]string{"id": "tokenId", "secret": "tokenSecret"}},
		"mapOAuthAuth":            {Import: "modal", Fields: map[string]string{"refresh_token": "oauthRefreshToken", "client_id": "oauthClientId", "client_secret": "oauthClientSecret", "jwt_key": "oauthJwtKey"}},
	}
	return mappinggen.ValidateBindings(provider, required)
}

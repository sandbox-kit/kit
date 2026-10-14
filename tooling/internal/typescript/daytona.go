package typescriptemit

import (
	"fmt"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	mappinggen "github.com/sandbox-kit/kit/tooling/internal/typescript/mapping_gen"
	"google.golang.org/protobuf/compiler/protogen"
)

func generateDaytonaMappings(plugin *protogen.Plugin, file *protogen.File, provider model.Provider) error {
	if err := validateDaytonaTypeScriptBindings(provider); err != nil {
		return err
	}
	return mappinggen.Generate(plugin, file, provider)
}

func validateDaytonaTypeScriptBindings(provider model.Provider) error {
	if provider.Rules.Client == nil {
		return fmt.Errorf("daytona TypeScript: client mapping is required")
	}
	required := map[string]mappinggen.Binding{
		"mapLifetimePolicies":   {Import: "@daytona/sdk", Fields: map[string]string{"stop": "autoStopInterval", "pause": "autoPauseInterval", "archive": "autoArchiveInterval", "delete": "autoDeleteInterval"}},
		"mapCommonCreateFields": {Import: "@daytona/sdk", Fields: map[string]string{"name": "name", "environment": "envVars", "labels": "labels"}},
		"mapResources":          {Import: "@daytona/sdk", Fields: map[string]string{"cpu": "cpu", "memory": "memory", "disk": "disk"}},
		"mapSandboxInfo":        {Import: "@daytona/sdk", Fields: map[string]string{"id": "id", "name": "name", "state": "state", "region": "target", "labels": "labels"}},
		"mapAllocatedResources": {Import: "@daytona/sdk", Fields: map[string]string{"cpu": "cpu", "memory": "memory", "disk": "disk"}},
		"mapClientSettings":     {Import: "@daytona/sdk", Fields: map[string]string{"endpoint": "apiUrl", "region": "target"}},
		"mapAPIKeyCredentials":  {Import: "@daytona/sdk", Fields: map[string]string{"key": "apiKey"}},
		"mapBearerAuth":         {Import: "@daytona/sdk", Fields: map[string]string{"token": "jwtToken"}},
		"mapScope":              {Import: "@daytona/sdk", Fields: map[string]string{"organization_id": "organizationId"}},
	}
	return mappinggen.ValidateBindings(provider, required)
}

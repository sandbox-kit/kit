// Native Daytona client setup is handwritten.
// mapClientConfig turns shared client settings into the native SDK constructor
// options.
//
// Example:
//
//	native, err := mapClientConfig(config)
package daytona

import (
	types "github.com/daytona/clients/sdk-go/pkg/types"
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
	os "os"
)

// mapClientConfig maps shared client configuration into the native SDK
// options. Unsupported supplied settings fail before the SDK is constructed.
func mapClientConfig(config *sandboxv1.Config) (types.DaytonaConfig, error) {
	var params types.DaytonaConfig
	if config == nil {
		config = &sandboxv1.Config{}
	}
	resolved := *config
	if resolved.Endpoint == nil {
		value := "https://app.daytona.io/api"
		if candidate := os.Getenv("DAYTONA_SERVER_URL"); candidate != "" {
			value = candidate
		}
		if candidate := os.Getenv("DAYTONA_API_URL"); candidate != "" {
			value = candidate
		}
		resolved.Endpoint = &value
	}
	if err := sandbox.ValidateConfig(&resolved); err != nil {
		return params, err
	}
	params.APIUrl = *resolved.Endpoint
	if resolved.Region != nil {
		params.Target = *resolved.Region
	}
	if resolved.Scope != nil {
		if resolved.Scope.AppName != nil || resolved.Scope.Environment != nil || resolved.Scope.ProjectId != nil {
			return params, unsupportedClientField("scope", "Daytona only supports organization_id scope")
		}
		if resolved.Scope.OrganizationId != nil {
			params.OrganizationID = *resolved.Scope.OrganizationId
		}
	}
	if auth := resolved.GetAuth(); auth != nil {
		switch credentials := auth.Credentials.(type) {
		case *sandbox.APIKey:
			if credentials.ApiKey == nil || credentials.ApiKey.Key == "" {
				return params, invalidClientField("auth.api_key", "api_key requires key")
			}
			params.APIKey = credentials.ApiKey.Key
		case *sandbox.BearerToken:
			if credentials.BearerToken == nil || credentials.BearerToken.Token == "" {
				return params, invalidClientField("auth.bearer_token", "bearer_token requires token")
			}
			params.JWTToken = credentials.BearerToken.Token
		default:
			return params, sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "initialize", Field: sandbox.ConfigFieldAuth, Message: "sandbox-kit daytona: unsupported authentication mode"}, nil)
		}
	}
	return params, nil
}

func invalidClientField(field, message string) error {
	return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "initialize", Field: field, Message: "sandbox-kit daytona: " + message}, nil)
}

func unsupportedClientField(field, message string) error {
	return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "initialize", Field: field, Message: "sandbox-kit daytona: " + message}, nil)
}

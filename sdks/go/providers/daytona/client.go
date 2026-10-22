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
)

const defaultEndpoint = "https://app.daytona.io/api"

// mapClientConfig maps shared client configuration into the native SDK
// options. Unsupported supplied settings fail before the SDK is constructed.
// Shared Config validation runs in sandbox.NewClient; this maps only.
// An omitted endpoint uses defaultEndpoint.
func mapClientConfig(config *sandboxv1.Config) (types.DaytonaConfig, error) {
	if config == nil {
		config = &sandboxv1.Config{}
	}
	params, _, err := mapClientSettings(config)
	if err != nil {
		return params, err
	}
	if config.Endpoint == nil {
		params.APIUrl = defaultEndpoint
	}
	if config.Scope != nil {
		if config.Scope.AppName != nil || config.Scope.Environment != nil || config.Scope.ProjectId != nil {
			return params, unsupportedClientField("scope", "Daytona only supports organization_id scope")
		}
		mapped, _, err := mapScope(config.Scope)
		if err != nil {
			return params, err
		}
		params.OrganizationID = mapped.OrganizationID
	}
	if auth := config.GetAuth(); auth != nil {
		switch credentials := auth.Credentials.(type) {
		case *sandbox.APIKey:
			if credentials.ApiKey == nil || credentials.ApiKey.Key == "" {
				return params, invalidClientField("auth.api_key", "api_key requires key")
			}
			mapped, _, err := mapAPIKeyCredentials(credentials.ApiKey)
			if err != nil {
				return params, err
			}
			params.APIKey = mapped.APIKey
		case *sandbox.BearerToken:
			if credentials.BearerToken == nil || credentials.BearerToken.Token == "" {
				return params, invalidClientField("auth.bearer_token", "bearer_token requires token")
			}
			mapped, _, err := mapBearerAuth(credentials.BearerToken)
			if err != nil {
				return params, err
			}
			params.JWTToken = mapped.JWTToken
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

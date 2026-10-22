// Native Modal client setup is handwritten.
// mapClientConfig turns shared client settings into the native SDK constructor
// options.
//
// Example:
//
//	native, err := mapClientConfig(config)
package modal

import (
	_go "github.com/modal-labs/modal-client/go"
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// mapClientConfig maps shared client configuration into the native SDK
// options. Unsupported supplied settings fail before the SDK is constructed.
func mapClientConfig(config *sandboxv1.Config) (_go.ClientParams, error) {
	var params _go.ClientParams
	if config == nil {
		config = &sandboxv1.Config{}
	}
	if config.Endpoint != nil {
		return params, sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "modal", Operation: "initialize", Field: sandbox.ConfigFieldEndpoint, Message: "sandbox-kit modal: Endpoint cannot be configured through this SDK's public client constructor"}, nil)
	}
	if config.Scope != nil {
		if config.Scope.OrganizationId != nil || config.Scope.ProjectId != nil {
			return params, unsupportedClientField("scope", "Modal only supports app_name and environment scope")
		}
		mapped, _, err := mapClientScope(config.Scope)
		if err != nil {
			return params, err
		}
		params.Environment = mapped.Environment
	}
	if auth := config.GetAuth(); auth != nil {
		switch credentials := auth.Credentials.(type) {
		case *sandbox.TokenPair:
			if credentials.TokenPair == nil || credentials.TokenPair.Id == "" || credentials.TokenPair.Secret == "" {
				return params, invalidClientField("auth.token_pair", "token_pair requires id and secret")
			}
			mapped, _, err := mapTokenPairCredentials(credentials.TokenPair)
			if err != nil {
				return params, err
			}
			params.TokenID = mapped.TokenID
			params.TokenSecret = mapped.TokenSecret
		case *sandbox.OAuthRefresh:
			if credentials.OauthRefresh == nil || credentials.OauthRefresh.RefreshToken == "" || credentials.OauthRefresh.ClientId == "" {
				return params, invalidClientField("auth.oauth_refresh", "oauth_refresh requires refresh_token and client_id")
			}
			if credentials.OauthRefresh.ClientSecret != nil && credentials.OauthRefresh.JwtKey != nil {
				return params, invalidClientField("auth.oauth_refresh", "oauth_refresh accepts only one of client_secret or jwt_key")
			}
			mapped, _, err := mapOAuthAuth(credentials.OauthRefresh)
			if err != nil {
				return params, err
			}
			params.OAuthCredentials = &mapped
		default:
			return params, sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "modal", Operation: "initialize", Field: sandbox.ConfigFieldAuth, Message: "sandbox-kit modal: unsupported authentication mode"}, nil)
		}
	}
	return params, nil
}

func invalidClientField(field, message string) error {
	return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "initialize", Field: field, Message: "sandbox-kit modal: " + message}, nil)
}

func unsupportedClientField(field, message string) error {
	return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "modal", Operation: "initialize", Field: field, Message: "sandbox-kit modal: " + message}, nil)
}

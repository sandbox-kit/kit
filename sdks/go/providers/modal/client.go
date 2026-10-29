package modal

import (
	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// mapClientConfig maps shared client configuration into the native SDK
// options. Unsupported supplied settings fail before the SDK is constructed.
// clientSettings translates shared Config into native Modal client options.
type clientSettings struct {
	source *sandbox.Config
	native sdk.ClientParams
}

func mapClientConfig(config *sandbox.Config) (sdk.ClientParams, error) {
	if config == nil {
		config = &sandbox.Config{}
	}
	settings := &clientSettings{source: config}
	if err := settings.applyEndpoint(); err != nil {
		return settings.native, err
	}
	if err := settings.applyScope(); err != nil {
		return settings.native, err
	}
	if err := settings.applyAuth(); err != nil {
		return settings.native, err
	}
	return settings.native, nil
}

func (s *clientSettings) applyEndpoint() error {
	if s.source.Endpoint == nil {
		return nil
	}
	return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "modal", Operation: "initialize", Field: sandbox.ConfigFieldEndpoint, Message: "sandbox-kit modal: Endpoint cannot be configured through this SDK's public client constructor"}, nil)
}

func (s *clientSettings) applyScope() error {
	if s.source.Scope == nil {
		return nil
	}
	if s.source.Scope.OrganizationId != nil {
		return unsupportedClientField(sandbox.ConfigFieldScopeOrganizationId, "unsupported client configuration field")
	}
	if s.source.Scope.ProjectId != nil {
		return unsupportedClientField(sandbox.ConfigFieldScopeProjectId, "unsupported client configuration field")
	}
	mapped, _, err := mapScope(s.source.Scope)
	if err != nil {
		return err
	}
	s.native.Environment = mapped.Environment
	return nil
}

func (s *clientSettings) applyAuth() error {
	auth := s.source.GetAuth()
	if auth == nil {
		return nil
	}
	switch credentials := auth.Credentials.(type) {
	case *sandbox.TokenPair:
		if credentials.TokenPair == nil {
			return invalidClientField(sandbox.ConfigFieldAuthTokenPair, "token_pair requires id and secret")
		}
		if credentials.TokenPair.Id == "" {
			return invalidClientField(sandbox.ConfigFieldAuthTokenPairId, "token_pair requires id and secret")
		}
		if credentials.TokenPair.Secret == "" {
			return invalidClientField(sandbox.ConfigFieldAuthTokenPairSecret, "token_pair requires id and secret")
		}
		mapped, _, err := mapTokenPairCredentials(credentials.TokenPair)
		if err != nil {
			return err
		}
		s.native.TokenID = mapped.TokenID
		s.native.TokenSecret = mapped.TokenSecret
	case *sandbox.OAuthRefresh:
		if credentials.OauthRefresh == nil {
			return invalidClientField(sandbox.ConfigFieldAuthOauthRefresh, "oauth_refresh requires refresh_token and client_id")
		}
		if credentials.OauthRefresh.RefreshToken == "" {
			return invalidClientField(sandbox.ConfigFieldAuthOauthRefreshRefreshToken, "oauth_refresh requires refresh_token and client_id")
		}
		if credentials.OauthRefresh.ClientId == "" {
			return invalidClientField(sandbox.ConfigFieldAuthOauthRefreshClientId, "oauth_refresh requires refresh_token and client_id")
		}
		if credentials.OauthRefresh.ClientSecret != nil && credentials.OauthRefresh.JwtKey != nil {
			return invalidClientField(sandbox.ConfigFieldAuthOauthRefresh, "oauth_refresh accepts only one of client_secret or jwt_key")
		}
		mapped, _, err := mapOAuthCredentials(credentials.OauthRefresh)
		if err != nil {
			return err
		}
		s.native.OAuthCredentials = &mapped
	default:
		return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "modal", Operation: "initialize", Field: sandbox.ConfigFieldAuth, Message: "sandbox-kit modal: unsupported authentication mode"}, nil)
	}
	return nil
}

func invalidClientField(field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "initialize", Field: field, Message: "sandbox-kit modal: " + message}, nil)
}

func unsupportedClientField(field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "modal", Operation: "initialize", Field: field, Message: "sandbox-kit modal: " + message}, nil)
}

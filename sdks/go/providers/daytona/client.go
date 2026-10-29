package daytona

import (
	types "github.com/daytona/clients/sdk-go/pkg/types"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

const defaultEndpoint = "https://app.daytona.io/api"

// mapClientConfig maps shared client configuration into the native SDK
// options. Unsupported supplied settings fail before the SDK is constructed.
// Shared Config validation runs in sandbox.NewClient; this maps only.
// An omitted endpoint uses defaultEndpoint.
// clientSettings translates shared Config into native Daytona client options.
type clientSettings struct {
	source *sandbox.Config
	native types.DaytonaConfig
}

func mapClientConfig(config *sandbox.Config) (types.DaytonaConfig, error) {
	if config == nil {
		config = &sandbox.Config{}
	}
	native, _, err := mapClientSettings(config)
	if err != nil {
		return native, err
	}
	settings := &clientSettings{source: config, native: native}
	if config.Endpoint == nil {
		settings.native.APIUrl = defaultEndpoint
	}
	if err := settings.applyScope(); err != nil {
		return settings.native, err
	}
	if err := settings.applyAuth(); err != nil {
		return settings.native, err
	}
	return settings.native, nil
}

func (s *clientSettings) applyScope() error {
	if s.source.Scope == nil {
		return nil
	}
	if s.source.Scope.AppName != nil {
		return unsupportedClientField(sandbox.ConfigFieldScopeAppName, "unsupported client configuration field")
	}
	if s.source.Scope.Environment != nil {
		return unsupportedClientField(sandbox.ConfigFieldScopeEnvironment, "unsupported client configuration field")
	}
	if s.source.Scope.ProjectId != nil {
		return unsupportedClientField(sandbox.ConfigFieldScopeProjectId, "unsupported client configuration field")
	}
	mapped, _, err := mapScope(s.source.Scope)
	if err != nil {
		return err
	}
	s.native.OrganizationID = mapped.OrganizationID
	return nil
}

func (s *clientSettings) applyAuth() error {
	auth := s.source.GetAuth()
	if auth == nil {
		return nil
	}
	switch credentials := auth.Credentials.(type) {
	case *sandbox.APIKey:
		if credentials.ApiKey == nil || credentials.ApiKey.Key == "" {
			return invalidClientField("auth.api_key", "api_key requires key")
		}
		mapped, _, err := mapAPIKeyCredentials(credentials.ApiKey)
		if err != nil {
			return err
		}
		s.native.APIKey = mapped.APIKey
	case *sandbox.BearerToken:
		if credentials.BearerToken == nil || credentials.BearerToken.Token == "" {
			return invalidClientField("auth.bearer_token", "bearer_token requires token")
		}
		mapped, _, err := mapBearerTokenCredentials(credentials.BearerToken)
		if err != nil {
			return err
		}
		s.native.JWTToken = mapped.JWTToken
	default:
		return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "initialize", Field: sandbox.ConfigFieldAuth, Message: "sandbox-kit daytona: unsupported authentication mode"}, nil)
	}
	return nil
}

func invalidClientField(field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "initialize", Field: field, Message: "sandbox-kit daytona: " + message}, nil)
}

func unsupportedClientField(field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "initialize", Field: field, Message: "sandbox-kit daytona: " + message}, nil)
}

package blaxel

import (
	"github.com/blaxel-ai/sdk-go/option"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// mapClientConfig maps shared client configuration into native SDK options.
// Endpoint is unsupported. Scope.Workspace maps to the native workspace
// option when set. Unsupported supplied settings fail before the SDK is
// constructed. Shared Config validation runs in sandbox.NewClient; this
// maps only.
type clientSettings struct {
	source *sandbox.Config
	opts   []option.RequestOption
}

func mapClientConfig(config *sandbox.Config) ([]option.RequestOption, error) {
	if config == nil {
		config = &sandbox.Config{}
	}
	settings := &clientSettings{source: config}
	if err := settings.applyEndpoint(); err != nil {
		return nil, err
	}
	if err := settings.applyScope(); err != nil {
		return nil, err
	}
	if err := settings.applyAuth(); err != nil {
		return nil, err
	}
	return settings.opts, nil
}

func (s *clientSettings) applyEndpoint() error {
	if s.source.Endpoint == nil {
		return nil
	}
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      sandbox.ErrorKindUnsupported,
		Provider:  "blaxel",
		Operation: "initialize",
		Field:     sandbox.ConfigFieldEndpoint,
		Message:   "sandbox-kit blaxel: Endpoint cannot be configured through this SDK's public client constructor",
	}, nil)
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
	if s.source.Scope.OrganizationId != nil {
		return unsupportedClientField(sandbox.ConfigFieldScopeOrganizationId, "unsupported client configuration field")
	}
	if s.source.Scope.ProjectId != nil {
		return unsupportedClientField(sandbox.ConfigFieldScopeProjectId, "unsupported client configuration field")
	}
	if s.source.Scope.Workspace == nil {
		return nil
	}
	workspace := s.source.Scope.GetWorkspace()
	if workspace == "" {
		return invalidClientField(sandbox.ConfigFieldScopeWorkspace, "workspace requires a non-blank value")
	}
	s.opts = append(s.opts, option.WithWorkspace(workspace))
	return nil
}

func (s *clientSettings) applyAuth() error {
	auth := s.source.GetAuth()
	if auth == nil {
		return nil
	}
	switch credentials := auth.Credentials.(type) {
	case *sandbox.APIKey:
		if credentials.ApiKey == nil {
			return invalidClientField(sandbox.ConfigFieldAuthApiKey, "api_key requires key")
		}
		if credentials.ApiKey.Key == "" {
			return invalidClientField(sandbox.ConfigFieldAuthApiKeyKey, "api_key requires key")
		}
		s.opts = append(s.opts, option.WithAPIKey(credentials.ApiKey.Key))
	default:
		return sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindUnsupported,
			Provider:  "blaxel",
			Operation: "initialize",
			Field:     sandbox.ConfigFieldAuth,
			Message:   "sandbox-kit blaxel: unsupported authentication mode",
		}, nil)
	}
	return nil
}

func invalidClientField(field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      sandbox.ErrorKindInvalidArgument,
		Provider:  "blaxel",
		Operation: "initialize",
		Field:     field,
		Message:   "sandbox-kit blaxel: " + message,
	}, nil)
}

func unsupportedClientField(field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      sandbox.ErrorKindUnsupported,
		Provider:  "blaxel",
		Operation: "initialize",
		Field:     field,
		Message:   "sandbox-kit blaxel: " + message,
	}, nil)
}

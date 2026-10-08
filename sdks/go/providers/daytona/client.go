package daytona

import (
	"fmt"

	sdk "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func newBackend(config *sandbox.Config) (sandbox.Backend, error) {
	params, err := mapClientConfig(config)
	if err != nil {
		return nil, err
	}
	native, err := sdk.NewClientWithConfig(&params)
	if err != nil {
		return nil, err
	}
	return &backend{client: native}, nil
}
func mapClientConfig(config *sandbox.Config) (types.DaytonaConfig, error) {
	if config == nil {
		config = &sandbox.Config{}
	}
	if err := sandbox.ValidateConfig(config); err != nil {
		return types.DaytonaConfig{}, err
	}
	params, remaining, err := mapClientSettings(config)
	if err != nil {
		return params, err
	}
	remaining.Provider = nil // Native factory attachment, not SDK settings.
	remaining.Timeout = nil  // Kit operation deadlines, not HTTP transport timeouts.
	if config.Scope != nil {
		mapped, rest, err := mapScope(config.Scope)
		if err != nil {
			return params, err
		}
		if err := sandbox.RejectUnmapped("daytona client context", &rest); err != nil {
			return params, err
		}
		params.OrganizationID = mapped.OrganizationID
	}
	remaining.Scope = nil
	auth := config.GetAuth()
	if auth != nil {
		switch {
		case auth.APIKey != nil:
			mapped, rest, err := mapAPIKeyCredentials(auth.APIKey)
			if err != nil {
				return params, err
			}
			if err := sandbox.RejectUnmapped("daytona API key", &rest); err != nil {
				return params, err
			}
			params.APIKey = mapped.APIKey
		case auth.BearerToken != nil:
			mapped, rest, err := mapBearerAuth(auth.BearerToken)
			if err != nil {
				return params, err
			}
			if err := sandbox.RejectUnmapped("daytona bearer token", &rest); err != nil {
				return params, err
			}
			params.JWTToken = mapped.JWTToken
		default:
			return params, fmt.Errorf("sandbox-kit daytona: use API-key or bearer-token authentication")
		}
	}
	remaining.Auth = nil
	if err := sandbox.RejectUnmapped("daytona client", &remaining); err != nil {
		return params, err
	}
	return params, nil
}

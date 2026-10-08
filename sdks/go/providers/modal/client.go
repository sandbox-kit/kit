package modal

import (
	"fmt"

	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func newBackend(config *sandbox.Config) (sandbox.Backend, error) {
	params, err := mapClientConfig(config)
	if err != nil {
		return nil, err
	}
	native, err := sdk.NewClientWithOptions(&params)
	if err != nil {
		return nil, err
	}
	return &backend{client: native, scope: &sandbox.Scope{AppName: sandbox.Value(config.GetScope().GetAppName()), Environment: sandbox.Value(config.GetScope().GetEnvironment())}, region: config.GetRegion()}, nil
}
func mapClientConfig(config *sandbox.Config) (sdk.ClientParams, error) {
	var params sdk.ClientParams
	if config == nil {
		config = &sandbox.Config{}
	}
	if err := sandbox.ValidateConfig(config); err != nil {
		return params, err
	}
	if config.Endpoint != nil {
		return params, fmt.Errorf("sandbox-kit modal: Endpoint cannot be configured through this SDK's public client constructor")
	}
	remaining := *config
	remaining.Provider, remaining.Auth, remaining.Endpoint, remaining.Region, remaining.Scope, remaining.Timeout = nil, nil, nil, nil, nil, nil
	if err := sandbox.RejectUnmapped("modal client", &remaining); err != nil {
		return params, err
	}
	scope := config.GetScope()
	if scope != nil {
		remaining := *scope
		remaining.AppName, remaining.Environment = nil, nil
		if err := sandbox.RejectUnmapped("modal client context", &remaining); err != nil {
			return params, err
		}
		params.Environment = scope.GetEnvironment()
	}
	auth := config.GetAuth()
	if auth == nil {
		return params, nil
	}
	switch {
	case auth.TokenPair != nil:
		mapped, remaining, err := mapTokenPairCredentials(auth.TokenPair)
		if err != nil {
			return params, err
		}
		if err := sandbox.RejectUnmapped("modal token pair", &remaining); err != nil {
			return params, err
		}
		params.TokenID, params.TokenSecret = mapped.TokenID, mapped.TokenSecret
	case auth.OAuthRefresh != nil:
		mapped, remaining, err := mapOAuthAuth(auth.OAuthRefresh)
		if err != nil {
			return params, err
		}
		if err := sandbox.RejectUnmapped("modal OAuth", &remaining); err != nil {
			return params, err
		}
		params.OAuthCredentials = &mapped
	default:
		return params, fmt.Errorf("sandbox-kit modal: use token-pair or OAuth-refresh authentication")
	}
	return params, nil
}

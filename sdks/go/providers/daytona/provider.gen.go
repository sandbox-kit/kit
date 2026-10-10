// Code generated from API declarations and provider bindings. DO NOT EDIT.
// source: kit/providers/daytona/v1/provider.proto
package daytona

import (
	context "context"
	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// Provider initializes the provider SDK from shared configuration.
type Provider struct {
}

// backend owns the initialized native SDK and captured configuration.
type backend struct {
	client *daytona.Client
}

func New() *Provider {
	return &Provider{}
}
func (*Provider) Name() string {
	return "daytona"
}
func (*backend) Name() string {
	return "daytona"
}
func (*Provider) NewClient(config *sandbox.Config) (sandbox.Backend, error) {
	backend, err := newBackend(config)
	return backend, mapProviderError(err, "initialize")
}
func newBackend(config *sandbox.Config) (sandbox.Backend, error) {
	if config == nil {
		config = &sandbox.Config{}
	}
	params, err := mapClientConfig(config)
	if err != nil {
		return nil, err
	}
	native, err := daytona.NewClientWithConfig(&params)
	if err != nil {
		return nil, err
	}
	result := &backend{client: native}
	captureBackendState(result, config)
	return result, nil
}
func captureBackendState(result *backend, config *sandbox.Config) {
	if config == nil {
		return
	}
}
func (a *backend) Create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	response, err := a.create(ctx, request)
	return response, mapProviderError(err, "create")
}
func (a *backend) Close(ctx context.Context) error {
	return mapProviderError(a.client.Close(ctx), "close")
}

var _ sandbox.Provider = (*Provider)(nil)
var _ sandbox.Backend = (*backend)(nil)

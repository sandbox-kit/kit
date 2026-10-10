// Code generated from API declarations and provider bindings. DO NOT EDIT.
// source: kit/providers/modal/v1/provider.proto
package modal

import (
	context "context"
	_go "github.com/modal-labs/modal-client/go"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// Provider initializes the provider SDK from shared configuration.
type Provider struct {
}

// backend owns the initialized native SDK and captured configuration.
type backend struct {
	client *_go.Client
	scope  *sandbox.Scope
	region *string
}

func New() *Provider {
	return &Provider{}
}
func (*Provider) Name() string {
	return "modal"
}
func (*backend) Name() string {
	return "modal"
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
	native, err := _go.NewClientWithOptions(&params)
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
	if config.Scope != nil {
		result.scope = &sandbox.Scope{}
		if config.Scope.AppName != nil {
			value := *config.Scope.AppName
			result.scope.AppName = &value
		}
		if config.Scope.Environment != nil {
			value := *config.Scope.Environment
			result.scope.Environment = &value
		}
		if config.Scope.OrganizationID != nil {
			value := *config.Scope.OrganizationID
			result.scope.OrganizationID = &value
		}
		if config.Scope.ProjectID != nil {
			value := *config.Scope.ProjectID
			result.scope.ProjectID = &value
		}
	}
	if config.Region != nil {
		value := *config.Region
		result.region = &value
	}
}
func (a *backend) Create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	response, err := a.create(ctx, request)
	return response, mapProviderError(err, "create")
}
func (a *backend) Close(ctx context.Context) error {
	a.client.Close()
	return nil
}

var _ sandbox.Provider = (*Provider)(nil)
var _ sandbox.Backend = (*backend)(nil)

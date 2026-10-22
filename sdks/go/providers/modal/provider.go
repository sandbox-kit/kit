// Provider construction and native SDK orchestration are handwritten.
// New returns a provider factory. Pass it to sandbox.NewClient. The public
// client owns SDK cleanup.
//
// Example:
//
//	client, err := sandbox.NewClient(New(), &sandbox.Config{})
package modal

import (
	context "context"

	_go "github.com/modal-labs/modal-client/go"
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// Provider initializes the provider SDK from shared configuration.
type Provider struct {
}

// backend owns the initialized native SDK and captured configuration.
type backend struct {
	// client is the initialized native SDK client. It is not part of the public
	// sandbox API.
	client *_go.Client
	// scope is the captured app and environment. Later changes to Config do not
	// affect it.
	scope *sandboxv1.Scope
	// region is the captured default region preference. Nil means the caller
	// omitted it.
	region *string
}

// New returns a provider factory. Pass it to sandbox.NewClient.
//
// Example:
//
//	client, err := sandbox.NewClient(New(), &sandbox.Config{})
func New() *Provider {
	return &Provider{}
}

// Name returns this integration's provider identifier.
func (*Provider) Name() string {
	return "modal"
}

// Name returns this integration's provider identifier.
func (*backend) Name() string {
	return "modal"
}

// NewClient maps config into the official SDK and returns a private backend.
// Pass the factory from New to sandbox.NewClient instead of calling this
// directly.
//
// Example:
//
//	backend, err := New().NewClient(config)
func (*Provider) NewClient(config *sandboxv1.Config) (sandbox.Backend, error) {
	backend, err := newBackend(config)
	return backend, mapProviderError(err, "initialize")
}

// newBackend maps config, constructs the native SDK, and captures provider
// state before returning the backend.
func newBackend(config *sandboxv1.Config) (sandbox.Backend, error) {
	if config == nil {
		config = &sandboxv1.Config{}
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

// captureBackendState copies retained configuration onto the backend. Later
// changes to the caller config do not affect the copy.
func captureBackendState(result *backend, config *sandboxv1.Config) {
	if config == nil {
		return
	}
	if config.Scope != nil {
		result.scope = &sandboxv1.Scope{}
		if config.Scope.AppName != nil {
			value := *config.Scope.AppName
			result.scope.AppName = &value
		}
		if config.Scope.Environment != nil {
			value := *config.Scope.Environment
			result.scope.Environment = &value
		}
		if config.Scope.OrganizationId != nil {
			value := *config.Scope.OrganizationId
			result.scope.OrganizationId = &value
		}
		if config.Scope.ProjectId != nil {
			value := *config.Scope.ProjectId
			result.scope.ProjectId = &value
		}
	}
	if config.Region != nil {
		value := *config.Region
		result.region = &value
	}
}

// Create creates one sandbox through the native SDK and returns shared
// metadata.
func (a *backend) Create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	response, err := a.create(ctx, request)
	return response, mapProviderError(err, "create")
}

// Close releases the native SDK client. It does not stop or delete the cloud
// sandbox. The Modal SDK Close call does not return an error; context
// cancellation is honored before release.
func (a *backend) Close(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return mapProviderError(err, "close")
	}
	a.client.Close()
	return nil
}

var _ sandbox.Provider = (*Provider)(nil)
var _ sandbox.Backend = (*backend)(nil)
var _ sandbox.SandboxCreator = (*backend)(nil)

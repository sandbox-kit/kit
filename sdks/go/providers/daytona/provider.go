// Provider construction and native SDK orchestration are handwritten.
// New returns a provider factory. Pass it to sandbox.NewClient. The public
// client owns SDK cleanup.
//
// Example:
//
//	client, err := sandbox.NewClient(New(), &sandbox.Config{})
package daytona

import (
	"context"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// Provider initializes the provider SDK from shared configuration.
type Provider struct {
}

// backend owns the initialized native SDK and captured configuration.
type backend struct {
	// client is the initialized native SDK client. It is not part of the public
	// sandbox API.
	client *daytona.Client
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
	return "daytona"
}

// Name returns this integration's provider identifier.
func (*backend) Name() string {
	return "daytona"
}

// NewClient maps config into the official SDK and returns a private backend.
// Pass the factory from New to sandbox.NewClient instead of calling this
// directly.
//
// Example:
//
//	backend, err := New().NewClient(config)
func (*Provider) NewClient(config *sandbox.Config) (sandbox.Backend, error) {
	backend, err := newBackend(config)
	return backend, mapProviderError(err, "initialize")
}

// newBackend maps config and constructs the native SDK.
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
	return &backend{client: native}, nil
}

// Create creates one sandbox through the native SDK and returns shared
// metadata.
func (a *backend) Create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	response, err := a.create(ctx, request)
	return response, mapProviderError(err, "create")
}

// Close releases the native SDK client. It does not stop or delete the cloud
// sandbox.
func (a *backend) Close(ctx context.Context) error {
	return mapProviderError(a.client.Close(ctx), "close")
}

var _ sandbox.Provider = (*Provider)(nil)
var _ sandbox.Backend = (*backend)(nil)
var _ sandbox.SandboxCreator = (*backend)(nil)

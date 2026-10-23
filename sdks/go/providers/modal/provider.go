// Provider construction and native SDK orchestration are handwritten.
// New returns a provider factory. Pass it to sandbox.NewClient. The public
// client owns SDK cleanup.
//
// Example:
//
//	client, err := sandbox.NewClient(New(), &sandbox.Config{})
package modal

import (
	"context"

	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/proto"
)

// Provider initializes the provider SDK from shared configuration.
type Provider struct {
}

// backend owns the initialized native SDK and captured configuration.
type backend struct {
	// client is the initialized native SDK client. It is not part of the public
	// sandbox API.
	client *sdk.Client
	// scope is the captured app and environment. Later changes to Config do not
	// affect it.
	scope *sandbox.Scope
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
func (*Provider) NewClient(config *sandbox.Config) (sandbox.Backend, error) {
	backend, err := newBackend(config)
	return backend, mapProviderError(err, "initialize")
}

// newBackend maps config, constructs the native SDK, and captures provider
// state before returning the backend.
func newBackend(config *sandbox.Config) (sandbox.Backend, error) {
	if config == nil {
		config = &sandbox.Config{}
	}
	params, err := mapClientConfig(config)
	if err != nil {
		return nil, err
	}
	native, err := sdk.NewClientWithOptions(&params)
	if err != nil {
		return nil, err
	}
	result := &backend{client: native}
	captureBackendState(result, config)
	return result, nil
}

// captureBackendState copies retained configuration onto the backend. Later
// changes to the caller config do not affect the copy.
func captureBackendState(result *backend, config *sandbox.Config) {
	if config == nil {
		return
	}
	if config.Scope != nil {
		result.scope = proto.Clone(config.Scope).(*sandbox.Scope)
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

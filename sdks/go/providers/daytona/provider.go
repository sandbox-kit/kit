package daytona

import (
	"context"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// Provider is the Daytona provider. Pass New to sandbox.NewClient.
type Provider struct {
}

// backend owns the initialized native SDK and captured configuration.
type backend struct {
	// client is the initialized native SDK client. It is not part of the public
	// sandbox API.
	client *daytona.Client
}

// New returns the Daytona provider. Pass it to sandbox.NewClient.
//
// Example:
//
//	client, err := sandbox.NewClient(New(), &sandbox.Config{})
func New() *Provider {
	return &Provider{}
}

// Name returns "daytona".
func (*Provider) Name() string {
	return "daytona"
}

// Name returns this integration's provider identifier.
func (*backend) Name() string {
	return "daytona"
}

// NewClient opens the Daytona SDK from config.
// Pass New to sandbox.NewClient. Calling NewClient directly skips the shared client.
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
func (b *backend) Create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	response, err := b.create(ctx, request)
	return response, mapProviderError(err, "create")
}

// Close releases the native SDK client. It does not stop or delete the cloud
// sandbox.
func (b *backend) Close(ctx context.Context) error {
	return mapProviderError(b.client.Close(ctx), "close")
}

// CloseSandbox is the handle close. Daytona has no sandbox close, so this
// returns without a provider request.
func (*backend) CloseSandbox(context.Context, string) error {
	return nil
}

var _ sandbox.Provider = (*Provider)(nil)
var _ sandbox.Backend = (*backend)(nil)
var _ sandbox.SandboxCreator = (*backend)(nil)

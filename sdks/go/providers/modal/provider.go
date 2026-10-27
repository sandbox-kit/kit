package modal

import (
	"context"

	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/proto"
)

// Provider is the Modal provider. Pass New to sandbox.NewClient.
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

// New returns the Modal provider. Pass it to sandbox.NewClient.
//
// Example:
//
//	client, err := sandbox.NewClient(New(), &sandbox.Config{})
func New() *Provider {
	return &Provider{}
}

// Name returns "modal".
func (*Provider) Name() string {
	return "modal"
}

// Name returns this integration's provider identifier.
func (*backend) Name() string {
	return "modal"
}

// NewClient opens the Modal SDK from config.
// Pass New to sandbox.NewClient. Calling NewClient directly skips the shared client.
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
	result.capture(config)
	return result, nil
}

// capture copies retained configuration onto the backend. Later changes to the
// caller config do not affect the copy.
func (b *backend) capture(config *sandbox.Config) {
	if config == nil {
		return
	}
	if config.Scope != nil {
		b.scope = proto.Clone(config.Scope).(*sandbox.Scope)
	}
	if config.Region != nil {
		value := *config.Region
		b.region = &value
	}
}

// Create creates one sandbox through the native SDK and returns shared
// metadata.
func (b *backend) Create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	response, err := b.create(ctx, request)
	return response, mapProviderError(err, "create")
}

// Close releases the native SDK client. It does not stop or delete the cloud
// sandbox. The Modal SDK Close call does not return an error; context
// cancellation is honored before release.
func (b *backend) Close(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return mapProviderError(err, "close")
	}
	b.client.Close()
	return nil
}

// CloseSandbox disconnects the local Modal sandbox object for this id.
// It does not terminate the cloud sandbox.
func (b *backend) CloseSandbox(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return mapProviderError(err, "close")
	}
	instance, err := b.client.Sandboxes.FromID(ctx, id, nil)
	if err != nil {
		return mapProviderError(err, "close")
	}
	return mapProviderError(instance.Detach(), "close")
}

var _ sandbox.Provider = (*Provider)(nil)
var _ sandbox.Backend = (*backend)(nil)
var _ sandbox.SandboxCreator = (*backend)(nil)

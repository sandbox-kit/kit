package blaxel

import (
	"context"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// Provider is the Blaxel provider. Pass New to sandbox.NewClient.
type Provider struct {
}

// backend owns the initialized native SDK and captured configuration.
type backend struct {
	// client is the initialized native SDK client. It is not part of the public
	// sandbox API.
	client *blaxel.Client
	// region is the captured default region. Nil means the caller omitted it;
	// create requires a region.
	region *string
}

// New returns the Blaxel provider. Pass it to sandbox.NewClient.
//
// Example:
//
//	client, err := sandbox.NewClient(New(), &sandbox.Config{
//		Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
//			ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("BL_API_KEY")},
//		}},
//		Scope: &sandbox.Scope{
//			Workspace: sandbox.Value(os.Getenv("BL_WORKSPACE")),
//		},
//		Region: sandbox.Value("us-pdx-1"),
//	})
func New() *Provider {
	return &Provider{}
}

// Name returns "blaxel".
func (*Provider) Name() string {
	return "blaxel"
}

// Name returns this integration's provider identifier.
func (*backend) Name() string {
	return "blaxel"
}

// NewClient opens the Blaxel SDK from config.
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
	opts, err := mapClientConfig(config)
	if err != nil {
		return nil, err
	}
	native := blaxel.NewClient(opts...)
	result := &backend{client: &native}
	result.capture(config)
	return result, nil
}

// capture copies retained configuration onto the backend. Later changes to the
// caller config do not affect the copy.
func (b *backend) capture(config *sandbox.Config) {
	if config == nil {
		return
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

// Close releases local provider state. It does not stop or delete the cloud
// sandbox. The Blaxel SDK client has no Close method; context cancellation is
// honored before return.
func (b *backend) Close(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return mapProviderError(err, "close")
	}
	return nil
}

// CloseSandbox does not contact Blaxel and does not change the sandbox.
func (*backend) CloseSandbox(context.Context, string) error {
	return nil
}

var _ sandbox.Provider = (*Provider)(nil)
var _ sandbox.Backend = (*backend)(nil)
var _ sandbox.SandboxCreator = (*backend)(nil)

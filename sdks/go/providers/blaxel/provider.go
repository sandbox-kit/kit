package blaxel

import (
	"context"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// Provider is the Blaxel provider. Pass New to sandbox.NewClient.
type Provider struct{}

// backend owns the initialized native SDK and captured configuration.
type backend struct {
	client *blaxel.Client
	// region is the captured default region; create requires it when set on Config.
	region *string
}

// New returns the Blaxel provider for sandbox.NewClient.
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
func New() *Provider { return &Provider{} }

func (*Provider) Name() string { return "blaxel" }
func (*backend) Name() string  { return "blaxel" }

// NewClient opens the Blaxel SDK from config.
// Prefer sandbox.NewClient(New(), config); calling this directly skips shared client wiring.
func (*Provider) NewClient(config *sandbox.Config) (sandbox.Backend, error) {
	backend, err := newBackend(config)
	return backend, mapProviderError(err, "initialize")
}

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

func (b *backend) capture(config *sandbox.Config) {
	if config == nil || config.Region == nil {
		return
	}
	value := *config.Region
	b.region = &value
}

func (b *backend) Create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	response, err := b.create(ctx, request)
	return response, mapProviderError(err, "create")
}

// Close releases local provider state only. It does not stop or delete cloud sandboxes.
func (b *backend) Close(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return mapProviderError(err, "close")
	}
	return nil
}

func (*backend) CloseSandbox(context.Context, string) error { return nil }

var _ sandbox.Provider = (*Provider)(nil)
var _ sandbox.Backend = (*backend)(nil)
var _ sandbox.SandboxCreator = (*backend)(nil)

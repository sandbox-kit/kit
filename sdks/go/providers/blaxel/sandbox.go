package blaxel

import (
	"context"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// create calls the Blaxel SDK and returns shared metadata.
// Create requires a name and an image source. Native create is used; create-if-
// not-exists is not.
func (b *backend) create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	if request == nil {
		request = &sandbox.CreateOptions{}
	}
	region := ""
	if b.region != nil {
		region = *b.region
	}
	plan, err := planCreate(request, region)
	if err != nil {
		return nil, err
	}
	instance, err := b.client.Sandboxes.NewInstance(ctx, plan.params)
	if err != nil {
		return nil, err
	}
	native, err := requireSandbox(instance, "create")
	if err != nil {
		return nil, err
	}
	info, err := sandboxInfo(native)
	if err != nil {
		return nil, err
	}
	return &sandbox.CreateResult{Sandbox: info}, nil
}

// RefreshSandbox retrieves current metadata for an existing Blaxel sandbox by name.
func (b *backend) RefreshSandbox(ctx context.Context, id string) (*sandbox.SandboxInfo, error) {
	instance, err := b.getInstance(ctx, id, "refresh")
	if err != nil {
		return nil, err
	}
	native, err := requireSandbox(instance, "refresh")
	if err != nil {
		return nil, err
	}
	return sandboxInfo(native)
}

var _ sandbox.SandboxRefresher = (*backend)(nil)

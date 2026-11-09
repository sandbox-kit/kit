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
	if instance == nil || instance.Sandbox == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: "create",
			Field:     sandbox.InfoFieldId,
			Message:   "sandbox-kit blaxel: SDK returned no sandbox",
		}, nil)
	}
	info, err := sandboxInfo(instance.Sandbox)
	if err != nil {
		return nil, err
	}
	return &sandbox.CreateResult{Sandbox: info}, nil
}

// RefreshSandbox retrieves current metadata for an existing Blaxel sandbox by name.
func (b *backend) RefreshSandbox(ctx context.Context, id string) (*sandbox.SandboxInfo, error) {
	instance, err := b.client.Sandboxes.GetInstance(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, "refresh")
	}
	if instance == nil || instance.Sandbox == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: "refresh",
			Field:     sandbox.InfoFieldId,
			Message:   "sandbox-kit blaxel: SDK returned no sandbox",
		}, nil)
	}
	info, err := sandboxInfo(instance.Sandbox)
	if err != nil {
		return nil, err
	}
	return info, nil
}

var _ sandbox.SandboxRefresher = (*backend)(nil)

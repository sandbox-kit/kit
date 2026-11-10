package blaxel

import (
	"context"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

var _ sandbox.SandboxDeleter = (*backend)(nil)

// DeleteSandbox permanently deletes a Blaxel sandbox by name.
// Start, stop, pause, archive, and terminate are not implemented.
func (b *backend) DeleteSandbox(ctx context.Context, id string, _ *sandbox.DeleteOptions) (*sandbox.SandboxInfo, error) {
	instance, err := b.getInstance(ctx, id, "delete")
	if err != nil {
		return nil, err
	}
	native, err := requireSandbox(instance, "delete")
	if err != nil {
		return nil, err
	}
	info, err := sandboxInfo(native)
	if err != nil {
		return nil, err
	}
	if err := instance.Delete(ctx); err != nil {
		return nil, mapProviderError(err, "delete")
	}
	info.State = sandbox.SandboxStateDeleted
	info.Terminal = sandbox.Value(true)
	info.ProviderState = sandbox.Value(string(blaxel.StatusTerminated))
	return info, nil
}

package blaxel

import (
	"context"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

var _ sandbox.SandboxDeleter = (*backend)(nil)

// DeleteSandbox permanently deletes a Blaxel sandbox by name.
// Start, stop, pause, archive, and terminate are not implemented.
// Blaxel's DELETE may still finish asynchronously; a successful call is
// reported as Deleted so portable callers get the same terminal outcome.
func (b *backend) DeleteSandbox(ctx context.Context, id string, _ *sandbox.DeleteOptions) (*sandbox.SandboxInfo, error) {
	instance, err := b.client.Sandboxes.DeleteInstance(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, "delete")
	}
	native, err := requireSandbox(instance, "delete")
	if err != nil {
		return nil, err
	}
	info, err := sandboxInfo(native)
	if err != nil {
		return nil, err
	}
	info.State = sandbox.SandboxStateDeleted
	info.Terminal = sandbox.Value(true)
	info.ProviderState = sandbox.Value(string(blaxel.StatusTerminated))
	return info, nil
}

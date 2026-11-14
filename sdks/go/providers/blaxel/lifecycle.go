package blaxel

import (
	"context"
	"errors"
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

var _ sandbox.SandboxDeleter = (*backend)(nil)

const deletePollInterval = time.Second

// DeleteSandbox permanently deletes a Blaxel sandbox by name.
// Start, stop, pause, archive, and terminate are not implemented.
// When Wait is omitted or false, a successful DELETE returns Deleted with
// Terminal true (Blaxel may still finish asynchronously). When Wait is true,
// Kit polls Get until the name is gone or Deleted.
func (b *backend) DeleteSandbox(ctx context.Context, id string, options *sandbox.DeleteOptions) (*sandbox.SandboxInfo, error) {
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
	if options != nil && options.GetWait() {
		// Blaxel has no native wait-for-delete; Kit polls until gone.
		return b.waitUntilDeleted(ctx, id)
	}
	info.State = sandbox.SandboxStateDeleted
	info.Terminal = sandbox.Value(true)
	info.ProviderState = sandbox.Value(string(blaxel.StatusTerminated))
	return info, nil
}

func (b *backend) waitUntilDeleted(ctx context.Context, id string) (*sandbox.SandboxInfo, error) {
	ticker := time.NewTicker(deletePollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, mapProviderError(ctx.Err(), "delete")
		case <-ticker.C:
			info, err := b.RefreshSandbox(ctx, id)
			if err != nil {
				var detail *sandbox.Error
				if errors.As(err, &detail) && detail.Kind == sandbox.ErrorKindNotFound {
					return &sandbox.SandboxInfo{
						Id:            id,
						Provider:      "blaxel",
						State:         sandbox.SandboxStateDeleted,
						Terminal:      sandbox.Value(true),
						ProviderState: sandbox.Value(string(blaxel.StatusTerminated)),
					}, nil
				}
				return nil, err
			}
			if info.State == sandbox.SandboxStateDeleted {
				info.Terminal = sandbox.Value(true)
				info.ProviderState = sandbox.Value(string(blaxel.StatusTerminated))
				return info, nil
			}
		}
	}
}

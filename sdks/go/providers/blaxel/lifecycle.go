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
	instance, err := b.client.Sandboxes.GetInstance(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, "delete")
	}
	if instance == nil || instance.Sandbox == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: "delete",
			Field:     sandbox.InfoFieldId,
			Message:   "sandbox-kit blaxel: SDK returned no sandbox",
		}, nil)
	}
	info, err := sandboxInfo(instance.Sandbox)
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

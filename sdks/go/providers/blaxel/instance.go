package blaxel

import (
	"context"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// getInstance loads a sandbox handle. Callers check FS / Process / Sandbox.
func (b *backend) getInstance(ctx context.Context, id, operation string) (*blaxel.SandboxInstance, error) {
	instance, err := b.client.Sandboxes.GetInstance(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	return instance, nil
}

func requireSandbox(instance *blaxel.SandboxInstance, operation string) (*blaxel.Sandbox, error) {
	if instance == nil || instance.Sandbox == nil {
		return nil, invalidResponse(operation, sandbox.InfoFieldId, "SDK returned no sandbox")
	}
	return instance.Sandbox, nil
}

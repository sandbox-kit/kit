package blaxel

import (
	"context"
	"strconv"
	"time"

	"github.com/blaxel-ai/sdk-go/option"
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
	instance, err := b.client.Sandboxes.NewInstance(ctx, plan.params, blaxelCreateOptions(request)...)
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

// blaxelCreateOptions forwards Blaxel's create timeout when Kit's create
// deadline is a whole number of seconds in 1..50. Outside that range Kit keeps
// the context deadline only (same bounds as @blaxel/core).
func blaxelCreateOptions(request *sandbox.CreateOptions) []option.RequestOption {
	seconds := blaxelCreateTimeoutSeconds(request)
	if seconds == 0 {
		return nil
	}
	return []option.RequestOption{
		option.WithHeader("X-Blaxel-Creation-Timeout", strconv.Itoa(seconds)),
	}
}

func blaxelCreateTimeoutSeconds(request *sandbox.CreateOptions) int {
	if request == nil || request.Provisioning == nil || request.Provisioning.Timeout == nil {
		return 0
	}
	timeout := request.Provisioning.Timeout.AsDuration()
	if timeout <= 0 || timeout%time.Second != 0 {
		return 0
	}
	seconds := int(timeout / time.Second)
	if seconds < 1 || seconds > 50 {
		return 0
	}
	return seconds
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

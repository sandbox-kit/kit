package modal

import (
	"context"
	"math"

	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/proto"
)

var _ sandbox.SandboxAttacher = (*backend)(nil)
var _ sandbox.SandboxTerminator = (*backend)(nil)

// create looks up the app, secrets, and image, then creates the sandbox.
// A nil request becomes an empty request. Modal still requires an image.
// The returned metadata is shared, not a Modal object.
func (b *backend) create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	if request == nil {
		request = &sandbox.CreateOptions{}
	}
	region := ""
	if b.region != nil {
		region = *b.region
	}
	plan, err := planCreate(request, b.scope, region)
	if err != nil {
		return nil, err
	}
	app, err := b.client.Apps.FromName(ctx, plan.app, &sdk.AppFromNameParams{Environment: plan.environment})
	if err != nil {
		return nil, err
	}
	for _, ref := range plan.secrets {
		secret, err := b.client.Secrets.FromName(ctx, ref.Reference, &sdk.SecretFromNameParams{Environment: plan.environment})
		if err != nil {
			return nil, err
		}
		plan.params.Secrets = append(plan.params.Secrets, secret)
	}
	if plan.proxy != "" {
		proxy, err := b.client.Proxies.FromName(ctx, plan.proxy, &sdk.ProxyFromNameParams{Environment: plan.environment})
		if err != nil {
			return nil, err
		}
		plan.params.Proxy = proxy
	}
	image := b.client.Images.FromRegistry(plan.image, nil)
	instance, err := b.client.Sandboxes.Create(ctx, app, image, &plan.params)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "modal",
			Operation: "create",
			Field:     sandbox.InfoFieldId,
			Message:   "sandbox-kit modal: SDK returned no sandbox",
		}, nil)
	}
	if plan.ready {
		if err := instance.WaitUntilReady(ctx, plan.readyTimeout, nil); err != nil {
			return nil, err
		}
	}
	info, err := mapSandboxIdentity(instance)
	if err != nil {
		return nil, err
	}
	info.Name = request.Name
	info.State = sandbox.SandboxStateRunning
	info.Labels = plan.params.Tags
	if request.Isolation != nil {
		info.Isolation = proto.Clone(request.Isolation).(*sandbox.IsolationConfig)
		info.Origins[sandbox.InfoFieldIsolation] = sandbox.ValueOriginRequest
	}
	if info.Name != nil {
		info.Origins["name"] = sandbox.ValueOriginRequest
	}
	if info.Labels != nil {
		info.Origins["labels"] = sandbox.ValueOriginRequest
	}
	return &sandbox.CreateResult{Sandbox: info}, nil
}

// AttachSandbox builds a handle from an ID. Modal does not look the sandbox up,
// so this does not refresh metadata.
func (b *backend) AttachSandbox(ctx context.Context, id string) (*sandbox.SandboxInfo, error) {
	instance, err := b.client.Sandboxes.FromID(ctx, id, nil)
	if err != nil {
		return nil, mapProviderError(err, "attach")
	}
	info, err := mapSandboxIdentity(instance)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// TerminateSandbox permanently finishes a Modal sandbox execution.
func (b *backend) TerminateSandbox(ctx context.Context, id string, options *sandbox.TerminateOptions) (*sandbox.ExitResult, error) {
	instance, err := b.client.Sandboxes.FromID(ctx, id, nil)
	if err != nil {
		return nil, mapProviderError(err, "terminate")
	}
	wait := options != nil && options.GetWait()
	code, err := instance.Terminate(ctx, &sdk.SandboxTerminateParams{Wait: wait})
	if err != nil {
		return nil, mapProviderError(err, "terminate")
	}
	result := &sandbox.ExitResult{}
	if wait {
		if code < math.MinInt32 || code > math.MaxInt32 {
			return nil, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidResponse, Provider: "modal", Operation: "terminate", Field: "exit_code", Message: "sandbox-kit modal: exit code exceeds the shared int32 range"}, nil)
		}
		result.ExitCode = sandbox.Value(int32(code))
	}
	return result, nil
}

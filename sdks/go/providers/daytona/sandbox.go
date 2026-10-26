// This file calls the Daytona SDK. creation.go builds the arguments, and
// metadata.go turns the SDK sandbox into shared metadata.
//
// Example:
//
//	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
//	instance, err := client.Create(ctx, nil)

package daytona

import (
	"context"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// create calls the Daytona SDK and returns shared metadata.
// A nil request uses the default snapshot and provider defaults.
func (b *backend) create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	if request == nil {
		request = &sandbox.CreateOptions{}
	}
	plan, err := planCreate(request)
	if err != nil {
		return nil, err
	}
	instance, err := b.client.Create(ctx, plan.params, plan.options...)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "daytona",
			Operation: "create",
			Field:     sandbox.InfoFieldId,
			Message:   "sandbox-kit daytona: SDK returned no sandbox",
		}, nil)
	}
	// Use the same complete resolver as Get and lifecycle responses so a newly
	// created handle immediately has the class and metadata needed for accurate
	// capability decisions.
	info, err := sandboxInfo(instance)
	if err != nil {
		return nil, err
	}
	return &sandbox.CreateResult{Sandbox: info}, nil
}

// RefreshSandbox retrieves current metadata for an existing Daytona sandbox.
func (b *backend) RefreshSandbox(ctx context.Context, id string) (*sandbox.SandboxInfo, error) {
	instance, err := b.client.Get(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, "refresh")
	}
	info, err := sandboxInfo(instance)
	if err != nil {
		return nil, err
	}
	return info, nil
}

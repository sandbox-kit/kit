package sandbox

import (
	"context"
	"strings"

	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
)

func (c *Client) resolveSandbox(ctx context.Context, id, operation string, invoke func(context.Context) (*SandboxInfo, error)) (*Sandbox, error) {
	provider := ""
	if c != nil && c.backend != nil {
		provider = c.backend.Name()
	}
	if ctx == nil {
		return nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: operation, Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c == nil || c.backend == nil {
		return nil, NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: provider, Operation: operation, Field: "client", Message: "sandbox-kit: client is not initialized"}, nil)
	}
	if strings.TrimSpace(id) == "" {
		return nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: operation, Field: InfoFieldId, Message: "sandbox-kit: sandbox ID is required"}, nil)
	}
	operationContext, cancel := ctx, func() {}
	if c.timeout != nil && *c.timeout > 0 {
		operationContext, cancel = context.WithTimeout(ctx, *c.timeout)
	}
	defer cancel()
	info, err := invoke(operationContext)
	if err != nil {
		return nil, sdk.WithErrorContext(err, provider, operation)
	}
	instance, err := sandboxFromResponse(c.backend, &CreateResult{Sandbox: info}, operation)
	if err != nil {
		return nil, err
	}
	instance.timeout = c.timeout
	return instance, nil
}

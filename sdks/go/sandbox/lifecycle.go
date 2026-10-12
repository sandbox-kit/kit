package sandbox

import "context"

// Get returns a handle for an existing sandbox and retrieves its current metadata.
func (c *Client) Get(ctx context.Context, id string) (*Sandbox, error) {
	provider := ""
	if c != nil && c.backend != nil {
		provider = c.backend.Name()
	}
	if ctx == nil {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "get", Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c == nil || c.backend == nil {
		return nil, NewError(ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: provider, Operation: "get", Field: "client", Message: "sandbox-kit: client is not initialized"}, nil)
	}
	if id == "" {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "get", Field: InfoFieldID, Message: "sandbox-kit: sandbox ID is required"}, nil)
	}
	refresher, ok := c.backend.(SandboxRefresher)
	if !ok {
		return nil, NewError(ErrorInfo{Kind: ErrorKindUnsupported, Provider: provider, Operation: "get", Message: "sandbox-kit: provider does not support sandbox retrieval"}, nil)
	}
	operationCtx := ctx
	cancel := func() {}
	if c.timeout != nil && *c.timeout > 0 {
		operationCtx, cancel = context.WithTimeout(ctx, *c.timeout)
	}
	defer cancel()
	info, err := refresher.RefreshSandbox(operationCtx, id)
	if err != nil {
		return nil, WithErrorContext(err, provider, "get")
	}
	instance, err := sandboxFromResponse(c.backend, &CreateResult{Sandbox: info}, "get")
	if err != nil {
		return nil, err
	}
	instance.timeout = c.timeout
	return instance, nil
}

// Attach returns a handle for an existing sandbox without claiming that the
// provider retrieved or validated its current metadata.
func (c *Client) Attach(ctx context.Context, id string) (*Sandbox, error) {
	provider := ""
	if c != nil && c.backend != nil {
		provider = c.backend.Name()
	}
	if ctx == nil {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "attach", Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c == nil || c.backend == nil {
		return nil, NewError(ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: provider, Operation: "attach", Field: "client", Message: "sandbox-kit: client is not initialized"}, nil)
	}
	if id == "" {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "attach", Field: InfoFieldID, Message: "sandbox-kit: sandbox ID is required"}, nil)
	}
	attacher, ok := c.backend.(SandboxAttacher)
	if !ok {
		return nil, NewError(ErrorInfo{Kind: ErrorKindUnsupported, Provider: provider, Operation: "attach", Message: "sandbox-kit: provider does not support sandbox attachment"}, nil)
	}
	operationCtx, cancel := ctx, func() {}
	if c.timeout != nil && *c.timeout > 0 {
		operationCtx, cancel = context.WithTimeout(ctx, *c.timeout)
	}
	defer cancel()
	info, err := attacher.AttachSandbox(operationCtx, id)
	if err != nil {
		return nil, WithErrorContext(err, provider, "attach")
	}
	instance, err := sandboxFromResponse(c.backend, &CreateResult{Sandbox: info}, "attach")
	if err != nil {
		return nil, err
	}
	instance.timeout = c.timeout
	return instance, nil
}

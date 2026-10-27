package sandbox

import (
	"context"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
)

func (s *Sandbox) operationContext(ctx context.Context, operation string, timeout *time.Duration) (context.Context, context.CancelFunc, string, error) {
	if ctx == nil {
		return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: operation, Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, "", err
	}
	if s == nil || s.operations == nil {
		return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Operation: operation, Message: "sandbox-kit: sandbox handle is not initialized"}, nil)
	}
	if s.currentID() == "" {
		return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Operation: operation, Message: "sandbox-kit: sandbox handle is not initialized"}, nil)
	}
	if timeout == nil {
		timeout = s.timeout
	}
	if timeout != nil && *timeout < 0 {
		return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: s.ProviderName(), Operation: operation, Field: "timeout", Message: "sandbox-kit: timeout cannot be negative"}, nil)
	}
	operationCtx := ctx
	baseCancel := func() {}
	if timeout != nil && *timeout > 0 {
		child, cancel := context.WithTimeout(ctx, *timeout)
		operationCtx, baseCancel = child, cancel
	}
	select {
	case s.operationGate <- struct{}{}:
		// Read the id after the gate is held. replaceInfo cannot run until this
		// operation releases the gate, so the provider call uses this snapshot.
		return operationCtx, func() { <-s.operationGate; baseCancel() }, s.currentID(), nil
	case <-operationCtx.Done():
		baseCancel()
		return nil, nil, "", operationCtx.Err()
	}
}

// Close does not stop, terminate, or delete the cloud sandbox. You can keep
// using the handle. A nil handle returns nil. A nil context returns
// ErrorKindInvalidArgument.
func (s *Sandbox) Close(ctx context.Context) error {
	if s == nil {
		if ctx == nil {
			return NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: "close_handle", Field: "context", Message: "sandbox-kit: context is required"}, nil)
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		return nil
	}
	operationCtx, release, id, err := s.operationContext(ctx, "close_handle", nil)
	if err != nil {
		return err
	}
	defer release()
	return sdk.WithErrorContext(s.operations.CloseSandbox(operationCtx, id), s.ProviderName(), "close_handle")
}

func (s *Sandbox) unsupported(operation string) error {
	return NewError(&ErrorInfo{Kind: ErrorKindUnsupported, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: provider does not support this sandbox operation"}, nil)
}

func (s *Sandbox) replaceInfo(operation string, info *SandboxInfo) error {
	if info == nil || info.Id == "" {
		return NewError(&ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: s.ProviderName(), Operation: operation, Field: InfoFieldId, Message: "sandbox-kit: provider returned no sandbox identity"}, nil)
	}
	provider := s.ProviderName()
	if info.Provider != "" && info.Provider != provider {
		return NewError(&ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: provider, Operation: operation, Field: InfoFieldProvider, Message: "sandbox-kit: provider returned inconsistent sandbox metadata"}, nil)
	}
	info.Provider = provider
	s.mu.Lock()
	s.info = cloneSandboxInfo(info)
	s.mu.Unlock()
	return nil
}

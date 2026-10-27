package sandbox

import (
	"context"
	"time"
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
	view := s.operationView()
	if view.id == "" {
		return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Operation: operation, Message: "sandbox-kit: sandbox handle is not initialized"}, nil)
	}
	if view.closed {
		return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox handle is closed"}, nil)
	}
	clientRelease := func() {}
	if s.owner != nil {
		var err error
		clientRelease, err = s.owner.beginOperation(operation)
		if err != nil {
			return nil, nil, "", err
		}
	}
	view = s.operationView()
	if view.terminal {
		clientRelease()
		return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox lifecycle is terminal"}, nil)
	}
	if timeout == nil {
		timeout = s.timeout
	}
	if timeout != nil && *timeout < 0 {
		clientRelease()
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
		// Close acquires this gate before marking the handle closed. Recheck
		// after the gate is held so a call that passed the earlier check cannot
		// start once Close has finished.
		view = s.operationView()
		if view.closed {
			<-s.operationGate
			baseCancel()
			clientRelease()
			return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox handle is closed"}, nil)
		}
		if view.terminal {
			<-s.operationGate
			baseCancel()
			clientRelease()
			return nil, nil, "", NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox lifecycle is terminal"}, nil)
		}
		// Read the id after the gate is held. replaceInfo cannot run until this
		// operation releases the gate, so the provider call uses this snapshot.
		return operationCtx, func() { <-s.operationGate; baseCancel(); clientRelease() }, view.id, nil
	case <-operationCtx.Done():
		baseCancel()
		clientRelease()
		return nil, nil, "", operationCtx.Err()
	}
}

// Close releases this local handle. It does not stop, terminate, or delete the cloud sandbox.
// A second call is safe. A nil handle returns nil. A nil context returns ErrorKindInvalidArgument.
func (s *Sandbox) Close(ctx context.Context) error {
	if ctx == nil {
		return NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: "close_handle", Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if s == nil {
		return nil
	}
	if s.operationGate == nil {
		return nil
	}
	select {
	case s.operationGate <- struct{}{}:
		defer func() { <-s.operationGate }()
	case <-ctx.Done():
		return ctx.Err()
	}
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
	return nil
}
func (s *Sandbox) markTerminal(state SandboxState, providerState string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.info == nil {
		return
	}
	s.info.State = state
	s.info.Terminal = Value(true)
	s.info.ProviderState = Value(providerState)
	s.capabilityReady = false
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
	s.capabilityReady = false
	s.mu.Unlock()
	return nil
}

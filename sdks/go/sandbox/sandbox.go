package sandbox

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Sandbox is the handle returned by Client.Create.
// It exposes identity and a copy of metadata. It does not expose the native SDK object.
// Lifecycle methods dispatch through private provider capabilities without exposing native SDK objects.
//
// Example:
//
//	instance, err := client.Create(ctx, nil)
//	fmt.Println(instance.ID(), instance.ProviderName())
//	info := instance.Info()
type Sandbox struct {
	mu            sync.RWMutex
	operations    Backend
	timeout       *time.Duration
	operationGate chan struct{}
	owner         *Client
	closed        bool
	// info is an owned metadata copy. Info returns another copy so callers cannot mutate this one.
	info *SandboxInfo
}

// sandboxFromResponse checks that creation returned an identifier
// and consistent provider metadata.
func sandboxFromResponse(backend Backend, response *CreateResult, operation string) (*Sandbox, error) {
	if response == nil || response.Sandbox == nil || response.Sandbox.Id == "" {
		return nil, NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  backend.Name(),
			Operation: operation,
			Field:     InfoFieldId,
			Message:   "sandbox-kit: provider returned no sandbox identity",
		}, nil)
	}
	info := cloneSandboxInfo(response.Sandbox)
	name := backend.Name()
	if info.Provider != "" && info.Provider != name {
		return nil, NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  name,
			Operation: operation,
			Field:     InfoFieldProvider,
			Message:   fmt.Sprintf("sandbox-kit: response provider %q differs from selected provider %q", info.Provider, name),
		}, nil)
	}
	info.Provider = name
	if err := validateResponseOrigins(info, operation); err != nil {
		return nil, err
	}
	return &Sandbox{operations: backend, info: info, operationGate: make(chan struct{}, 1)}, nil
}

// ID returns the provider sandbox identifier.
func (s *Sandbox) ID() string {
	if s == nil {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.info == nil {
		return ""
	}
	return s.info.Id
}

// ProviderName returns the selected provider, such as "daytona" or "modal".
func (s *Sandbox) ProviderName() string {
	if s == nil {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.info == nil {
		return ""
	}
	return s.info.Provider
}

// Info returns an owned metadata copy. Changing the result does not change the handle.
//
// Example:
//
//	info := instance.Info()
//	fmt.Println(info.GetID(), info.Origins[InfoFieldId])
func (s *Sandbox) Info() *SandboxInfo {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.info == nil {
		return nil
	}
	info := cloneSandboxInfo(s.info)
	return info
}

type sandboxCapabilityResolver interface {
	SandboxCapabilities(*SandboxInfo) []LifecycleCapability
}

// Capabilities returns the operations supported by this provider and resolved
// sandbox class. A terminal handle has no further lifecycle capabilities.
func (s *Sandbox) Capabilities() []LifecycleCapability {
	if s == nil || s.operations == nil {
		return nil
	}
	s.mu.RLock()
	closed := s.closed
	s.mu.RUnlock()
	if closed {
		return nil
	}
	if s.owner != nil && s.owner.isClosed() {
		return nil
	}
	info := s.Info()
	if info == nil || info.GetTerminal() {
		return nil
	}
	if resolver, ok := s.operations.(sandboxCapabilityResolver); ok {
		values := resolver.SandboxCapabilities(info)
		return append([]LifecycleCapability(nil), values...)
	}
	return staticLifecycleCapabilities(s.operations)
}

// Supports reports whether capability is available for this resolved handle.
func (s *Sandbox) Supports(capability LifecycleCapability) bool {
	for _, candidate := range s.Capabilities() {
		if candidate == capability {
			return true
		}
	}
	return false
}

func (s *Sandbox) operationContext(ctx context.Context, operation string, timeout *time.Duration) (context.Context, context.CancelFunc, error) {
	if ctx == nil {
		return nil, nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: operation, Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if s == nil || s.operations == nil || s.ID() == "" {
		return nil, nil, NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Operation: operation, Message: "sandbox-kit: sandbox handle is not initialized"}, nil)
	}
	s.mu.RLock()
	closed := s.closed
	s.mu.RUnlock()
	if closed {
		return nil, nil, NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox handle is closed"}, nil)
	}
	clientRelease := func() {}
	if s.owner != nil {
		var err error
		clientRelease, err = s.owner.beginOperation(operation)
		if err != nil {
			return nil, nil, err
		}
	}
	if info := s.Info(); info != nil && info.GetTerminal() {
		clientRelease()
		return nil, nil, NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox lifecycle is terminal"}, nil)
	}
	if timeout == nil {
		timeout = s.timeout
	}
	if timeout != nil && *timeout < 0 {
		clientRelease()
		return nil, nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: s.ProviderName(), Operation: operation, Field: "timeout", Message: "sandbox-kit: timeout cannot be negative"}, nil)
	}
	operationCtx := ctx
	baseCancel := func() {}
	if timeout != nil && *timeout > 0 {
		child, cancel := context.WithTimeout(ctx, *timeout)
		operationCtx, baseCancel = child, cancel
	}
	select {
	case s.operationGate <- struct{}{}:
		if info := s.Info(); info != nil && info.GetTerminal() {
			<-s.operationGate
			baseCancel()
			clientRelease()
			return nil, nil, NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox lifecycle is terminal"}, nil)
		}
		return operationCtx, func() { <-s.operationGate; baseCancel(); clientRelease() }, nil
	case <-operationCtx.Done():
		baseCancel()
		clientRelease()
		return nil, nil, operationCtx.Err()
	}
}

// Close releases this local handle. It never changes the remote sandbox. Close
// is idempotent.
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

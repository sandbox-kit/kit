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
	operations    any
	timeout       *time.Duration
	operationGate chan struct{}
	// info is an owned metadata copy. Info returns another copy so callers cannot mutate this one.
	info *SandboxInfo
}

// sandboxFromResponse checks that creation returned an identifier
// and consistent provider metadata.
func sandboxFromResponse(backend Backend, response *CreateResult, operation string) (*Sandbox, error) {
	if response == nil || response.Sandbox == nil || response.Sandbox.ID == "" {
		return nil, NewError(ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  backend.Name(),
			Operation: operation,
			Field:     InfoFieldID,
			Message:   "sandbox-kit: provider returned no sandbox identity",
		}, nil)
	}
	info, err := response.Sandbox.Clone()
	if err != nil {
		return nil, NewError(ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  backend.Name(),
			Operation: operation,
			Field:     InfoFieldProviderMetadata,
			Message:   err.Error(),
		}, err)
	}
	name := backend.Name()
	if info.Provider != "" && info.Provider != name {
		return nil, NewError(ErrorInfo{
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
	return s.info.ID
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
//	fmt.Println(info.GetID(), info.Origins[InfoFieldID])
func (s *Sandbox) Info() *SandboxInfo {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.info == nil {
		return nil
	}
	info, _ := s.info.Clone()
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
	info := s.Info()
	if info == nil || info.GetTerminal() {
		return nil
	}
	if resolver, ok := s.operations.(sandboxCapabilityResolver); ok {
		values := resolver.SandboxCapabilities(info)
		return append([]LifecycleCapability(nil), values...)
	}
	values := make([]LifecycleCapability, 0, 7)
	if _, ok := s.operations.(SandboxRefresher); ok {
		values = append(values, LifecycleCapabilityRefresh)
	}
	if _, ok := s.operations.(SandboxStarter); ok {
		values = append(values, LifecycleCapabilityStart)
	}
	if _, ok := s.operations.(SandboxStopper); ok {
		values = append(values, LifecycleCapabilityStop)
	}
	if _, ok := s.operations.(SandboxPauser); ok {
		values = append(values, LifecycleCapabilityPause)
	}
	if _, ok := s.operations.(SandboxTerminator); ok {
		values = append(values, LifecycleCapabilityTerminate)
	}
	if _, ok := s.operations.(SandboxArchiver); ok {
		values = append(values, LifecycleCapabilityArchive)
	}
	if _, ok := s.operations.(SandboxDeleter); ok {
		values = append(values, LifecycleCapabilityDelete)
	}
	return values
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
		return nil, nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: operation, Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if s == nil || s.operations == nil || s.ID() == "" {
		return nil, nil, NewError(ErrorInfo{Kind: ErrorKindFailedPrecondition, Operation: operation, Message: "sandbox-kit: sandbox handle is not initialized"}, nil)
	}
	if info := s.Info(); info != nil && info.GetTerminal() {
		return nil, nil, NewError(ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox lifecycle is terminal"}, nil)
	}
	if timeout == nil {
		timeout = s.timeout
	}
	if timeout != nil && *timeout < 0 {
		return nil, nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: s.ProviderName(), Operation: operation, Field: "timeout", Message: "sandbox-kit: timeout cannot be negative"}, nil)
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
			return nil, nil, NewError(ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: sandbox lifecycle is terminal"}, nil)
		}
		return operationCtx, func() { <-s.operationGate; baseCancel() }, nil
	case <-operationCtx.Done():
		baseCancel()
		return nil, nil, operationCtx.Err()
	}
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
	return NewError(ErrorInfo{Kind: ErrorKindUnsupported, Provider: s.ProviderName(), Operation: operation, Message: "sandbox-kit: provider does not support this sandbox operation"}, nil)
}

func (s *Sandbox) replaceInfo(operation string, info *SandboxInfo) error {
	if info == nil || info.ID == "" {
		return NewError(ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: s.ProviderName(), Operation: operation, Field: InfoFieldID, Message: "sandbox-kit: provider returned no sandbox identity"}, nil)
	}
	provider := s.ProviderName()
	if info.Provider != "" && info.Provider != provider {
		return NewError(ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: provider, Operation: operation, Field: InfoFieldProvider, Message: "sandbox-kit: provider returned inconsistent sandbox metadata"}, nil)
	}
	info.Provider = provider
	copy, err := info.Clone()
	if err != nil {
		return NewError(ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: provider, Operation: operation, Field: InfoFieldProviderMetadata, Message: err.Error()}, err)
	}
	s.mu.Lock()
	s.info = copy
	s.mu.Unlock()
	return nil
}

// Refresh replaces cached metadata with the provider's current metadata.
func (s *Sandbox) Refresh(ctx context.Context) error {
	ctx, cancel, err := s.operationContext(ctx, "refresh", nil)
	if err != nil {
		return err
	}
	defer cancel()
	op, ok := s.operations.(SandboxRefresher)
	if !ok {
		return s.unsupported("refresh")
	}
	info, err := op.RefreshSandbox(ctx, s.ID())
	if err != nil {
		return WithErrorContext(err, s.ProviderName(), "refresh")
	}
	if info == nil || info.ID == "" {
		return NewError(ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: s.ProviderName(), Operation: "refresh", Field: InfoFieldID, Message: "sandbox-kit: refresh returned no sandbox identity"}, nil)
	}
	return s.replaceInfo("refresh", info)
}

func (s *Sandbox) Start(ctx context.Context, o *StartOptions) error {
	var t *time.Duration
	if o != nil {
		t = o.Timeout
	}
	ctx, c, e := s.operationContext(ctx, "start", t)
	if e != nil {
		return e
	}
	defer c()
	op, ok := s.operations.(SandboxStarter)
	if !ok {
		return s.unsupported("start")
	}
	info, err := op.StartSandbox(ctx, s.ID(), o)
	if err != nil {
		return WithErrorContext(err, s.ProviderName(), "start")
	}
	return s.replaceInfo("start", info)
}
func (s *Sandbox) Stop(ctx context.Context, o *StopOptions) error {
	var t *time.Duration
	if o != nil {
		t = o.Timeout
	}
	ctx, c, e := s.operationContext(ctx, "stop", t)
	if e != nil {
		return e
	}
	defer c()
	op, ok := s.operations.(SandboxStopper)
	if !ok {
		return s.unsupported("stop")
	}
	info, err := op.StopSandbox(ctx, s.ID(), o)
	if err != nil {
		return WithErrorContext(err, s.ProviderName(), "stop")
	}
	return s.replaceInfo("stop", info)
}
func (s *Sandbox) Pause(ctx context.Context, o *PauseOptions) error {
	var t *time.Duration
	if o != nil {
		t = o.Timeout
	}
	ctx, c, e := s.operationContext(ctx, "pause", t)
	if e != nil {
		return e
	}
	defer c()
	op, ok := s.operations.(SandboxPauser)
	if !ok {
		return s.unsupported("pause")
	}
	info, err := op.PauseSandbox(ctx, s.ID(), o)
	if err != nil {
		return WithErrorContext(err, s.ProviderName(), "pause")
	}
	return s.replaceInfo("pause", info)
}
func (s *Sandbox) Terminate(ctx context.Context, o *TerminateOptions) (*ExitResult, error) {
	var t *time.Duration
	if o != nil {
		t = o.Timeout
	}
	ctx, c, e := s.operationContext(ctx, "terminate", t)
	if e != nil {
		return nil, e
	}
	defer c()
	op, ok := s.operations.(SandboxTerminator)
	if !ok {
		return nil, s.unsupported("terminate")
	}
	r, e := op.TerminateSandbox(ctx, s.ID(), o)
	if e != nil {
		return r, WithErrorContext(e, s.ProviderName(), "terminate")
	}
	s.markTerminal(SandboxStateTerminated, "terminated")
	return r, nil
}
func (s *Sandbox) Archive(ctx context.Context, o *ArchiveOptions) error {
	var t *time.Duration
	if o != nil {
		t = o.Timeout
	}
	ctx, c, e := s.operationContext(ctx, "archive", t)
	if e != nil {
		return e
	}
	defer c()
	op, ok := s.operations.(SandboxArchiver)
	if !ok {
		return s.unsupported("archive")
	}
	info, err := op.ArchiveSandbox(ctx, s.ID(), o)
	if err != nil {
		return WithErrorContext(err, s.ProviderName(), "archive")
	}
	return s.replaceInfo("archive", info)
}
func (s *Sandbox) Delete(ctx context.Context, o *DeleteOptions) error {
	var t *time.Duration
	if o != nil {
		t = o.Timeout
	}
	ctx, c, e := s.operationContext(ctx, "delete", t)
	if e != nil {
		return e
	}
	defer c()
	op, ok := s.operations.(SandboxDeleter)
	if !ok {
		return s.unsupported("delete")
	}
	info, err := op.DeleteSandbox(ctx, s.ID(), o)
	if err != nil {
		return WithErrorContext(err, s.ProviderName(), "delete")
	}
	if err := s.replaceInfo("delete", info); err != nil {
		return err
	}
	s.markTerminal(SandboxStateDeleted, "deleted")
	return nil
}

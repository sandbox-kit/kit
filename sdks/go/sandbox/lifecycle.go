// Sandbox lifecycle behavior is handwritten. Portable option/result types come
// from protobuf; Backend capability interfaces are optional provider adapters.
package sandbox

import (
	context "context"
	reflect "reflect"
	time "time"

	"google.golang.org/protobuf/proto"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

func staticLifecycleCapabilities(operations Backend) []LifecycleCapability {
	values := make([]LifecycleCapability, 0, 7)
	if _, ok := operations.(SandboxRefresher); ok {
		values = append(values, LifecycleCapabilityRefresh)
	}
	if _, ok := operations.(SandboxStarter); ok {
		values = append(values, LifecycleCapabilityStart)
	}
	if _, ok := operations.(SandboxStopper); ok {
		values = append(values, LifecycleCapabilityStop)
	}
	if _, ok := operations.(SandboxPauser); ok {
		values = append(values, LifecycleCapabilityPause)
	}
	if _, ok := operations.(SandboxTerminator); ok {
		values = append(values, LifecycleCapabilityTerminate)
	}
	if _, ok := operations.(SandboxArchiver); ok {
		values = append(values, LifecycleCapabilityArchive)
	}
	if _, ok := operations.(SandboxDeleter); ok {
		values = append(values, LifecycleCapabilityDelete)
	}
	return values
}

func (s *Sandbox) runInfoOp(ctx context.Context, operation string, timeout *time.Duration, invoke func(context.Context) (*SandboxInfo, error)) error {
	operationContext, cancel, err := s.operationContext(ctx, operation, timeout)
	if err != nil {
		return err
	}
	defer cancel()
	info, err := invoke(operationContext)
	if err != nil {
		return WithErrorContext(err, s.ProviderName(), operation)
	}
	return s.replaceInfo(operation, info)
}

func (s *Sandbox) Refresh(ctx context.Context) error {
	return s.runInfoOp(ctx, "refresh", nil, func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxRefresher)
		if !ok {
			return nil, s.unsupported("refresh")
		}
		return capability.RefreshSandbox(operationContext, s.ID())
	})
}

func (s *Sandbox) Start(ctx context.Context, options *StartOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	return s.runInfoOp(ctx, "start", timeout, func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxStarter)
		if !ok {
			return nil, s.unsupported("start")
		}
		return capability.StartSandbox(operationContext, s.ID(), providerOptions(options, timeout, func() *StartOptions { return &StartOptions{} }, func(cloned *StartOptions) {
			cloned.Timeout = Duration(*timeout)
		}))
	})
}

func (s *Sandbox) Stop(ctx context.Context, options *StopOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	return s.runInfoOp(ctx, "stop", timeout, func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxStopper)
		if !ok {
			return nil, s.unsupported("stop")
		}
		return capability.StopSandbox(operationContext, s.ID(), providerOptions(options, timeout, func() *StopOptions { return &StopOptions{} }, func(cloned *StopOptions) {
			cloned.Timeout = Duration(*timeout)
		}))
	})
}

func (s *Sandbox) Pause(ctx context.Context, options *PauseOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	return s.runInfoOp(ctx, "pause", timeout, func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxPauser)
		if !ok {
			return nil, s.unsupported("pause")
		}
		return capability.PauseSandbox(operationContext, s.ID(), providerOptions(options, timeout, func() *PauseOptions { return &PauseOptions{} }, func(cloned *PauseOptions) {
			cloned.Timeout = Duration(*timeout)
		}))
	})
}

func (s *Sandbox) Terminate(ctx context.Context, options *TerminateOptions) (*ExitResult, error) {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	operationContext, cancel, err := s.operationContext(ctx, "terminate", timeout)
	if err != nil {
		return nil, err
	}
	defer cancel()
	capability, ok := s.operations.(SandboxTerminator)
	if !ok {
		return nil, s.unsupported("terminate")
	}
	outcome, err := capability.TerminateSandbox(operationContext, s.ID(), providerOptions(options, timeout, func() *TerminateOptions { return &TerminateOptions{} }, func(cloned *TerminateOptions) {
		cloned.Timeout = Duration(*timeout)
	}))
	if err != nil {
		return outcome, WithErrorContext(err, s.ProviderName(), "terminate")
	}
	s.markTerminal(SandboxStateTerminated, "terminated")
	return outcome, nil
}

func (s *Sandbox) Archive(ctx context.Context, options *ArchiveOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	return s.runInfoOp(ctx, "archive", timeout, func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxArchiver)
		if !ok {
			return nil, s.unsupported("archive")
		}
		return capability.ArchiveSandbox(operationContext, s.ID(), providerOptions(options, timeout, func() *ArchiveOptions { return &ArchiveOptions{} }, func(cloned *ArchiveOptions) {
			cloned.Timeout = Duration(*timeout)
		}))
	})
}

func (s *Sandbox) Delete(ctx context.Context, options *DeleteOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	operationContext, cancel, err := s.operationContext(ctx, "delete", timeout)
	if err != nil {
		return err
	}
	defer cancel()
	capability, ok := s.operations.(SandboxDeleter)
	if !ok {
		return s.unsupported("delete")
	}
	info, err := capability.DeleteSandbox(operationContext, s.ID(), providerOptions(options, timeout, func() *DeleteOptions { return &DeleteOptions{} }, func(cloned *DeleteOptions) {
		cloned.Timeout = Duration(*timeout)
	}))
	if err != nil {
		return WithErrorContext(err, s.ProviderName(), "delete")
	}
	if err := s.replaceInfo("delete", info); err != nil {
		return err
	}
	// Mark the handle terminal before cancel releases the operation gate.
	s.markTerminal(SandboxStateDeleted, "deleted")
	return nil
}

func (s *Sandbox) inheritedTimeout() *time.Duration {
	if s == nil {
		return nil
	}
	return s.timeout
}

// effectiveTimeout prefers an explicit operation timeout, including explicit
// zero. A nil explicit value inherits the client timeout.
func effectiveTimeout(explicit *durationpb.Duration, inherited *time.Duration) *time.Duration {
	if timeout := OptionalDuration(explicit); timeout != nil {
		return timeout
	}
	return inherited
}

// providerOptions copies options and records timeout for the provider. A nil
// timeout leaves the caller value unchanged, including a nil options pointer.
func providerOptions[T proto.Message](options T, timeout *time.Duration, blank func() T, apply func(T)) T {
	if timeout == nil {
		return options
	}
	cloned := blank()
	value := reflect.ValueOf(any(options))
	if value.Kind() == reflect.Pointer && !value.IsNil() {
		cloned = proto.Clone(options).(T)
	}
	apply(cloned)
	return cloned
}

func (c *Client) Get(ctx context.Context, id string) (*Sandbox, error) {
	return c.resolveSandbox(ctx, id, "get", func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := c.backend.(SandboxRefresher)
		if !ok {
			return nil, NewError(&ErrorInfo{Kind: ErrorKindUnsupported, Provider: c.backend.Name(), Operation: "get", Message: "sandbox-kit: provider does not support sandbox get"}, nil)
		}
		return capability.RefreshSandbox(operationContext, id)
	})
}

func (c *Client) Attach(ctx context.Context, id string) (*Sandbox, error) {
	return c.resolveSandbox(ctx, id, "attach", func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := c.backend.(SandboxAttacher)
		if !ok {
			return nil, NewError(&ErrorInfo{Kind: ErrorKindUnsupported, Provider: c.backend.Name(), Operation: "attach", Message: "sandbox-kit: provider does not support sandbox attach"}, nil)
		}
		return capability.AttachSandbox(operationContext, id)
	})
}

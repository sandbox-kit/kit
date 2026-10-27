package sandbox

import (
	context "context"
	time "time"

	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
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

func (s *Sandbox) runInfoOp(ctx context.Context, operation string, timeout *time.Duration, invoke func(context.Context, string) (*SandboxInfo, error)) error {
	operationContext, cancel, id, err := s.operationContext(ctx, operation, timeout)
	if err != nil {
		return err
	}
	defer cancel()
	info, err := invoke(operationContext, id)
	if err != nil {
		return sdk.WithErrorContext(err, s.ProviderName(), operation)
	}
	return s.replaceInfo(operation, info)
}

// Refresh loads the sandbox's current metadata from the provider and stores it on this handle.
// Daytona supports Refresh. Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Refresh(ctx context.Context) error {
	return s.runInfoOp(ctx, "refresh", nil, func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxRefresher)
		if !ok {
			return nil, s.unsupported("refresh")
		}
		return capability.RefreshSandbox(operationContext, id)
	})
}

// Start starts a stopped sandbox and stores the updated metadata on this handle.
// A nil options value uses the client timeout. Set StartOptions.Timeout to override it.
// Daytona supports Start. Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Start(ctx context.Context, options *StartOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	return s.runInfoOp(ctx, "start", timeout, func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxStarter)
		if !ok {
			return nil, s.unsupported("start")
		}
		return capability.StartSandbox(operationContext, id, startOptionsForCall(options, timeout))
	})
}

// Stop stops a sandbox and stores the updated metadata on this handle.
// Set StopOptions.Force with Value(true) to request an immediate stop. A nil options value uses the client timeout.
// Daytona supports Stop. Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Stop(ctx context.Context, options *StopOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	return s.runInfoOp(ctx, "stop", timeout, func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxStopper)
		if !ok {
			return nil, s.unsupported("stop")
		}
		return capability.StopSandbox(operationContext, id, stopOptionsForCall(options, timeout))
	})
}

// Pause pauses a sandbox and stores the updated metadata on this handle.
// Call Supports with LifecycleCapabilityPause first. Some Daytona sandbox classes cannot pause.
// Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Pause(ctx context.Context, options *PauseOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	return s.runInfoOp(ctx, "pause", timeout, func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxPauser)
		if !ok {
			return nil, s.unsupported("pause")
		}
		return capability.PauseSandbox(operationContext, id, pauseOptionsForCall(options, timeout))
	})
}

// Terminate ends a sandbox. The returned ExitResult includes the exit code when Wait is set.
// Set TerminateOptions.Wait with Value(true) to wait. A nil options value uses the client timeout.
// Modal supports Terminate. Daytona returns ErrorKindUnsupported. Use Delete or Stop for Daytona.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Terminate(ctx context.Context, options *TerminateOptions) (*ExitResult, error) {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	operationContext, cancel, id, err := s.operationContext(ctx, "terminate", timeout)
	if err != nil {
		return nil, err
	}
	defer cancel()
	capability, ok := s.operations.(SandboxTerminator)
	if !ok {
		return nil, s.unsupported("terminate")
	}
	outcome, err := capability.TerminateSandbox(operationContext, id, terminateOptionsForCall(options, timeout))
	if err != nil {
		return outcome, sdk.WithErrorContext(err, s.ProviderName(), "terminate")
	}
	s.markTerminal(SandboxStateTerminated, "terminated")
	return outcome, nil
}

// Archive archives a sandbox and stores the updated metadata on this handle.
// Call Supports with LifecycleCapabilityArchive first. Some Daytona sandbox classes cannot archive.
// Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Archive(ctx context.Context, options *ArchiveOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	return s.runInfoOp(ctx, "archive", timeout, func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxArchiver)
		if !ok {
			return nil, s.unsupported("archive")
		}
		return capability.ArchiveSandbox(operationContext, id, archiveOptionsForCall(options, timeout))
	})
}

// Delete deletes the cloud sandbox and marks this handle deleted.
// Close only releases the local handle.
// Daytona supports Delete. Modal returns ErrorKindUnsupported. Use Terminate for Modal.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Delete(ctx context.Context, options *DeleteOptions) error {
	timeout := effectiveTimeout(options.GetTimeout(), s.inheritedTimeout())
	operationContext, cancel, id, err := s.operationContext(ctx, "delete", timeout)
	if err != nil {
		return err
	}
	defer cancel()
	capability, ok := s.operations.(SandboxDeleter)
	if !ok {
		return s.unsupported("delete")
	}
	info, err := capability.DeleteSandbox(operationContext, id, deleteOptionsForCall(options, timeout))
	if err != nil {
		return sdk.WithErrorContext(err, s.ProviderName(), "delete")
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
	if timeout := sdk.OptionalDuration(explicit); timeout != nil {
		return timeout
	}
	return inherited
}

// optionsWithTimeout passes options through when no timeout is recorded, including a nil pointer.
// Otherwise build returns a new message so the provider cannot mutate the caller.
func optionsWithTimeout[T any](options *T, timeout *time.Duration, build func(*durationpb.Duration, *T) *T) *T {
	if timeout == nil {
		return options
	}
	return build(Duration(*timeout), options)
}

func startOptionsForCall(options *StartOptions, timeout *time.Duration) *StartOptions {
	return optionsWithTimeout(options, timeout, func(recorded *durationpb.Duration, _ *StartOptions) *StartOptions {
		return &StartOptions{Timeout: recorded}
	})
}

func stopOptionsForCall(options *StopOptions, timeout *time.Duration) *StopOptions {
	return optionsWithTimeout(options, timeout, func(recorded *durationpb.Duration, options *StopOptions) *StopOptions {
		cloned := &StopOptions{Timeout: recorded}
		if options != nil {
			cloned.Force = copyOptionalBool(options.Force)
		}
		return cloned
	})
}

func pauseOptionsForCall(options *PauseOptions, timeout *time.Duration) *PauseOptions {
	return optionsWithTimeout(options, timeout, func(recorded *durationpb.Duration, _ *PauseOptions) *PauseOptions {
		return &PauseOptions{Timeout: recorded}
	})
}

func terminateOptionsForCall(options *TerminateOptions, timeout *time.Duration) *TerminateOptions {
	return optionsWithTimeout(options, timeout, func(recorded *durationpb.Duration, options *TerminateOptions) *TerminateOptions {
		cloned := &TerminateOptions{Timeout: recorded}
		if options != nil {
			cloned.Wait = copyOptionalBool(options.Wait)
		}
		return cloned
	})
}

func archiveOptionsForCall(options *ArchiveOptions, timeout *time.Duration) *ArchiveOptions {
	return optionsWithTimeout(options, timeout, func(recorded *durationpb.Duration, _ *ArchiveOptions) *ArchiveOptions {
		return &ArchiveOptions{Timeout: recorded}
	})
}

func deleteOptionsForCall(options *DeleteOptions, timeout *time.Duration) *DeleteOptions {
	return optionsWithTimeout(options, timeout, func(recorded *durationpb.Duration, options *DeleteOptions) *DeleteOptions {
		cloned := &DeleteOptions{Timeout: recorded}
		if options != nil {
			cloned.Wait = copyOptionalBool(options.Wait)
		}
		return cloned
	})
}

func copyOptionalBool(value *bool) *bool {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}

// Get looks up a sandbox by ID and returns a handle the provider has verified.
// An empty ID returns ErrorKindInvalidArgument with Field InfoFieldId.
// Daytona supports Get. Modal returns ErrorKindUnsupported. Use Attach for Modal.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (c *Client) Get(ctx context.Context, id string) (*Sandbox, error) {
	return c.resolveSandbox(ctx, id, "get", func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := c.backend.(SandboxRefresher)
		if !ok {
			return nil, NewError(&ErrorInfo{Kind: ErrorKindUnsupported, Provider: c.backend.Name(), Operation: "get", Message: "sandbox-kit: provider does not support sandbox get"}, nil)
		}
		return capability.RefreshSandbox(operationContext, id)
	})
}

// Attach returns a handle for an existing sandbox ID without asking the provider to look it up.
// An empty ID returns ErrorKindInvalidArgument with Field InfoFieldId.
// Modal supports Attach. Daytona returns ErrorKindUnsupported. Use Get for Daytona.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (c *Client) Attach(ctx context.Context, id string) (*Sandbox, error) {
	return c.resolveSandbox(ctx, id, "attach", func(operationContext context.Context) (*SandboxInfo, error) {
		capability, ok := c.backend.(SandboxAttacher)
		if !ok {
			return nil, NewError(&ErrorInfo{Kind: ErrorKindUnsupported, Provider: c.backend.Name(), Operation: "attach", Message: "sandbox-kit: provider does not support sandbox attach"}, nil)
		}
		return capability.AttachSandbox(operationContext, id)
	})
}

package sandbox

import (
	context "context"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"google.golang.org/protobuf/proto"
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
// A nil timeout uses the client timeout as this call's deadline. The provider receives a timeout only when one is set.
// Daytona supports Start. Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Start(ctx context.Context, options *StartOptions) error {
	return s.runInfoOp(ctx, "start", sdk.OptionalDuration(options.GetTimeout()), func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxStarter)
		if !ok {
			return nil, s.unsupported("start")
		}
		return capability.StartSandbox(operationContext, id, callOptions(options))
	})
}

// Stop stops a sandbox and stores the updated metadata on this handle.
// Set StopOptions.Force with Value(true) to request an immediate stop.
// A nil timeout uses the client timeout as this call's deadline. The provider receives a timeout only when one is set.
// Daytona supports Stop. Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Stop(ctx context.Context, options *StopOptions) error {
	return s.runInfoOp(ctx, "stop", sdk.OptionalDuration(options.GetTimeout()), func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxStopper)
		if !ok {
			return nil, s.unsupported("stop")
		}
		return capability.StopSandbox(operationContext, id, callOptions(options))
	})
}

// Pause pauses a sandbox and stores the updated metadata on this handle.
// A nil timeout uses the client timeout as this call's deadline. The provider receives a timeout only when one is set.
// Daytona exposes Pause and accepts or rejects the sandbox. Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Pause(ctx context.Context, options *PauseOptions) error {
	return s.runInfoOp(ctx, "pause", sdk.OptionalDuration(options.GetTimeout()), func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxPauser)
		if !ok {
			return nil, s.unsupported("pause")
		}
		return capability.PauseSandbox(operationContext, id, callOptions(options))
	})
}

// Terminate ends a sandbox. The returned ExitResult includes the exit code when Wait is set.
// Set TerminateOptions.Wait with Value(true) to wait until the sandbox exits.
// A nil timeout uses the client timeout as this call's deadline. Modal waits on that deadline and does not take a separate terminate duration.
// Modal supports Terminate. Daytona returns ErrorKindUnsupported. Use Delete or Stop for Daytona.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Terminate(ctx context.Context, options *TerminateOptions) (*ExitResult, error) {
	operationContext, cancel, id, err := s.operationContext(ctx, "terminate", sdk.OptionalDuration(options.GetTimeout()))
	if err != nil {
		return nil, err
	}
	defer cancel()
	capability, ok := s.operations.(SandboxTerminator)
	if !ok {
		return nil, s.unsupported("terminate")
	}
	outcome, err := capability.TerminateSandbox(operationContext, id, callOptions(options))
	if err != nil {
		return outcome, sdk.WithErrorContext(err, s.ProviderName(), "terminate")
	}
	return outcome, nil
}

// Archive archives a sandbox and stores the updated metadata on this handle.
// A nil timeout uses the client timeout as this call's deadline. Daytona archive does not take a separate duration.
// Daytona exposes Archive and accepts or rejects the sandbox. Modal returns ErrorKindUnsupported.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Archive(ctx context.Context, options *ArchiveOptions) error {
	return s.runInfoOp(ctx, "archive", sdk.OptionalDuration(options.GetTimeout()), func(operationContext context.Context, id string) (*SandboxInfo, error) {
		capability, ok := s.operations.(SandboxArchiver)
		if !ok {
			return nil, s.unsupported("archive")
		}
		return capability.ArchiveSandbox(operationContext, id, callOptions(options))
	})
}

// Delete deletes the cloud sandbox and stores the metadata the provider returns.
// Set DeleteOptions.Wait with Value(true) to wait until deletion finishes. Value(false) returns after the provider accepts the request.
// A nil timeout uses the client timeout as this call's deadline. The provider receives a timeout only when one is set.
// Close leaves the cloud sandbox in place. Daytona supports Delete. Modal returns ErrorKindUnsupported. Use Terminate for Modal.
// A nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
func (s *Sandbox) Delete(ctx context.Context, options *DeleteOptions) error {
	operationContext, cancel, id, err := s.operationContext(ctx, "delete", sdk.OptionalDuration(options.GetTimeout()))
	if err != nil {
		return err
	}
	defer cancel()
	capability, ok := s.operations.(SandboxDeleter)
	if !ok {
		return s.unsupported("delete")
	}
	info, err := capability.DeleteSandbox(operationContext, id, callOptions(options))
	if err != nil {
		return sdk.WithErrorContext(err, s.ProviderName(), "delete")
	}
	return s.replaceInfo("delete", info)
}

// callOptions returns a copy so the provider cannot change the caller's message.
// A nil value stays nil. Fields the caller left unset stay unset.
func callOptions[T interface {
	proto.Message
	*StartOptions | *StopOptions | *PauseOptions | *TerminateOptions | *ArchiveOptions | *DeleteOptions
}](options T) T {
	if options == nil {
		return options
	}
	return proto.CloneOf(options)
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

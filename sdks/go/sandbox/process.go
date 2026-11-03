package sandbox

import (
	"context"
	"strings"

	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Process runs commands inside one sandbox and waits for each to finish.
// Obtain it from Sandbox.Process. Set exactly one of Argv or Shell.
// A working directory, when set, must be absolute. A nil timeout inherits
// Config.Timeout. Explicit zero adds no Kit deadline.
type Process struct {
	owner *Sandbox
}

// Process returns the command runner for this sandbox.
// Methods on a nil handle return ErrorKindFailedPrecondition.
func (s *Sandbox) Process() *Process {
	if s == nil {
		return &Process{}
	}
	return &Process{owner: s}
}

// Run runs one command and returns its exit code and collected output.
// A non-zero exit code is not an error. Set exactly one of options.Argv or
// options.Shell. A nil context returns ErrorKindInvalidArgument.
func (p *Process) Run(ctx context.Context, options *ExecOptions) (*ExecResult, error) {
	var result *ExecResult
	err := p.call(ctx, "run", options, func(operationContext context.Context, id string) error {
		runner, err := p.runner("run")
		if err != nil {
			return err
		}
		result, err = runner.RunCommand(operationContext, id, cloneMessage(options))
		return err
	})
	if err != nil {
		return nil, err
	}
	return cloneMessage(result), nil
}

func (p *Process) call(ctx context.Context, operation string, options *ExecOptions, invoke func(context.Context, string) error) error {
	owner := p.owner
	if owner == nil {
		if ctx == nil {
			return NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: operation, Field: "context", Message: "sandbox-kit: context is required"}, nil)
		}
		return NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Operation: operation, Message: "sandbox-kit: sandbox handle is not initialized"}, nil)
	}
	if err := validateExecOptions(owner, operation, options); err != nil {
		return err
	}
	var timeout *durationpb.Duration
	if options != nil {
		timeout = options.Timeout
	}
	operationContext, cancel, id, err := owner.operationContext(ctx, operation, sdk.OptionalDuration(timeout))
	if err != nil {
		return err
	}
	defer cancel()
	if err := invoke(operationContext, id); err != nil {
		return sdk.WithErrorContext(err, owner.ProviderName(), operation)
	}
	return nil
}

func (p *Process) runner(operation string) (SandboxProcess, error) {
	runner, ok := p.owner.operations.(SandboxProcess)
	if !ok {
		return nil, p.owner.unsupported(operation)
	}
	return runner, nil
}

func validateExecOptions(owner *Sandbox, operation string, options *ExecOptions) error {
	if options == nil {
		return NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidArgument,
			Provider:  owner.ProviderName(),
			Operation: operation,
			Field:     ExecFieldArgv,
			Message:   "sandbox-kit: set argv or shell",
		}, nil)
	}
	hasArgv := len(options.Argv) > 0
	hasShell := options.Shell != nil
	switch {
	case hasArgv && hasShell:
		return NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidArgument,
			Provider:  owner.ProviderName(),
			Operation: operation,
			Field:     ExecFieldArgv,
			Message:   "sandbox-kit: set exactly one of argv or shell",
		}, nil)
	case !hasArgv && !hasShell:
		return NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidArgument,
			Provider:  owner.ProviderName(),
			Operation: operation,
			Field:     ExecFieldArgv,
			Message:   "sandbox-kit: set argv or shell",
		}, nil)
	case hasShell && strings.TrimSpace(*options.Shell) == "":
		return NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidArgument,
			Provider:  owner.ProviderName(),
			Operation: operation,
			Field:     ExecFieldShell,
			Message:   "sandbox-kit: shell must not be blank",
		}, nil)
	}
	if options.WorkingDirectory != nil {
		path := *options.WorkingDirectory
		if path == "" || !strings.HasPrefix(path, "/") || strings.ContainsRune(path, 0) {
			return NewError(&ErrorInfo{
				Kind:      ErrorKindInvalidArgument,
				Provider:  owner.ProviderName(),
				Operation: operation,
				Field:     ExecFieldWorkingDirectory,
				Message:   "sandbox-kit: working directory must be absolute",
			}, nil)
		}
	}
	return nil
}

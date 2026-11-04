package modal

import (
	"context"
	"io"
	"sync"
	"time"

	modalsdk "github.com/modal-labs/modal-client/go"
	internalsdk "github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/types/known/durationpb"
)

var _ sandbox.SandboxProcess = (*backend)(nil)

func (b *backend) RunCommand(ctx context.Context, id string, request *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
	instance, err := b.processSandbox(ctx, id, "run")
	if err != nil {
		return nil, err
	}
	command, err := modalCommand(request)
	if err != nil {
		return nil, err
	}
	timeout, err := commandTimeout("run", request.GetTimeout())
	if err != nil {
		return nil, err
	}
	params := &modalsdk.SandboxExecParams{
		Stdout:  modalsdk.Pipe,
		Stderr:  modalsdk.Pipe,
		Timeout: timeout,
	}
	if request.WorkingDirectory != nil {
		params.Workdir = *request.WorkingDirectory
	}
	if request.Environment != nil {
		params.Env = request.Environment
	}
	process, err := instance.Exec(ctx, command, params)
	if err != nil {
		return nil, mapProviderError(err, "run")
	}
	if process == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "modal",
			Operation: "run",
			Message:   "sandbox-kit modal: SDK returned no process",
		}, nil)
	}
	if process.Stdin != nil {
		_ = process.Stdin.Close()
	}
	stdout, stderr, exitCode, err := collectProcess(ctx, process)
	if err != nil {
		return nil, mapProviderError(err, "run")
	}
	return &sandbox.ExecResult{
		ExitCode: int32(exitCode),
		Stdout:   stdout,
		Stderr:   stderr,
	}, nil
}

func (b *backend) processSandbox(ctx context.Context, id, operation string) (*modalsdk.Sandbox, error) {
	instance, err := b.client.Sandboxes.FromID(ctx, id, nil)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	if instance == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "modal",
			Operation: operation,
			Message:   "sandbox-kit modal: SDK returned no sandbox",
		}, nil)
	}
	return instance, nil
}

func modalCommand(request *sandbox.ExecOptions) ([]string, error) {
	if request.Shell != nil {
		return []string{"bash", "-c", *request.Shell}, nil
	}
	if len(request.Argv) == 0 {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "modal",
			Operation: "run",
			Field:     sandbox.ExecFieldArgv,
			Message:   "sandbox-kit: set argv or shell",
		}, nil)
	}
	return append([]string(nil), request.Argv...), nil
}

// commandTimeout maps Kit ExecOptions.timeout to Modal's command timeout.
// Nil and explicit zero become 0 (no command timeout). Non-zero values must
// be whole seconds; Modal rejects fractional durations.
func commandTimeout(operation string, value *durationpb.Duration) (time.Duration, error) {
	timeout, ok := internalsdk.WholeSecondDuration(value)
	if !ok {
		return 0, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "modal",
			Operation: operation,
			Field:     sandbox.ExecFieldTimeout,
			Message:   "sandbox-kit modal: timeout must be a whole number of seconds",
		}, nil)
	}
	return timeout, nil
}

func collectProcess(ctx context.Context, process *modalsdk.ContainerProcess) (stdout, stderr []byte, exitCode int, err error) {
	var wg sync.WaitGroup
	var stdoutErr, stderrErr error
	wg.Add(2)
	go func() {
		defer wg.Done()
		if process.Stdout == nil {
			return
		}
		stdout, stdoutErr = io.ReadAll(process.Stdout)
	}()
	go func() {
		defer wg.Done()
		if process.Stderr == nil {
			return
		}
		stderr, stderrErr = io.ReadAll(process.Stderr)
	}()
	exitCode, err = process.Wait(ctx, nil)
	wg.Wait()
	if err != nil {
		return nil, nil, 0, err
	}
	if stdoutErr != nil {
		return nil, nil, 0, stdoutErr
	}
	if stderrErr != nil {
		return nil, nil, 0, stderrErr
	}
	if stdout == nil {
		stdout = []byte{}
	}
	if stderr == nil {
		stderr = []byte{}
	}
	return stdout, stderr, exitCode, nil
}

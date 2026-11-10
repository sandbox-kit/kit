package blaxel

import (
	"context"
	"strings"
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/types/known/durationpb"
)

// blaxelCompletionHold is the longest positive Kit timeout that stays on
// Blaxel's waitForCompletion request. Blaxel's public docs limit that hold
// to 60 seconds; the edge in front of the sandbox can close a longer one.
const blaxelCompletionHold = 60 * time.Second

// blaxelWaitUnlimited is Process.Wait's maxWait sentinel for no added
// deadline. The SDK treats 0 as one minute and only -1 as unlimited.
const blaxelWaitUnlimited = time.Duration(-1)

// processPollInterval is the gap between process status reads. It matches
// the Blaxel SDK default.
const processPollInterval = time.Second

var _ sandbox.SandboxProcess = (*backend)(nil)

func (b *backend) RunCommand(ctx context.Context, id string, request *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
	instance, err := b.processSandbox(ctx, id, "run")
	if err != nil {
		return nil, err
	}
	command, err := blaxelCommand(request)
	if err != nil {
		return nil, err
	}
	timeout, err := commandTimeout("run", request.GetTimeout())
	if err != nil {
		return nil, err
	}
	body := blaxel.ProcessRequestParam{
		Command: command,
		// Explicit zero disables Blaxel's default process kill window.
		Timeout:           blaxel.Int(int64(timeout / time.Second)),
		WaitForCompletion: blaxel.Bool(completionHold(timeout)),
	}
	if request.WorkingDirectory != nil {
		body.WorkingDir = blaxel.String(*request.WorkingDirectory)
	}
	if request.Environment != nil {
		body.Env = request.Environment
	}
	started, err := instance.Process.New(ctx, body)
	if err != nil {
		return nil, mapProviderError(err, "run")
	}
	if started == nil || completionHold(timeout) || processFinished(started) {
		return execResult(started)
	}
	identifier := processIdentifier(started)
	if identifier == "" {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: "run",
			Message:   "sandbox-kit blaxel: SDK returned no process identifier",
		}, nil)
	}
	// Wait polls until a terminal status or until ctx ends. A gateway 504 on
	// one read is retried inside Wait; the deadline surfaces as Timeout.
	finished, err := instance.Process.Wait(ctx, identifier, waitBudget(timeout), processPollInterval)
	if err != nil {
		return nil, mapProviderError(err, "run")
	}
	return execResult(finished)
}

// completionHold reports whether this timeout can stay on Blaxel's
// waitForCompletion request. Omit, zero, and anything above 60 seconds
// start the process and poll instead.
func completionHold(timeout time.Duration) bool {
	return timeout > 0 && timeout <= blaxelCompletionHold
}

// waitBudget is the Process.Wait maxWait for a Kit timeout. Omit and zero
// add no deadline beyond ctx.
func waitBudget(timeout time.Duration) time.Duration {
	if timeout <= 0 {
		return blaxelWaitUnlimited
	}
	return timeout
}

func processFinished(response *blaxel.ProcessResponse) bool {
	if response == nil {
		return false
	}
	switch response.Status {
	case blaxel.ProcessResponseStatusCompleted,
		blaxel.ProcessResponseStatusFailed,
		blaxel.ProcessResponseStatusKilled,
		blaxel.ProcessResponseStatusStopped:
		return true
	default:
		return false
	}
}

func processIdentifier(response *blaxel.ProcessResponse) string {
	if response == nil {
		return ""
	}
	if response.Pid != "" {
		return response.Pid
	}
	return response.Name
}

func execResult(response *blaxel.ProcessResponse) (*sandbox.ExecResult, error) {
	if response == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: "run",
			Message:   "sandbox-kit blaxel: SDK returned no exec result",
		}, nil)
	}
	return &sandbox.ExecResult{
		ExitCode: int32(response.ExitCode),
		Stdout:   []byte(response.Stdout),
		Stderr:   []byte(response.Stderr),
	}, nil
}

func (b *backend) processSandbox(ctx context.Context, id, operation string) (*blaxel.SandboxInstance, error) {
	instance, err := b.client.Sandboxes.GetInstance(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	if instance == nil || instance.Process == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: operation,
			Message:   "sandbox-kit blaxel: SDK returned no process API",
		}, nil)
	}
	return instance, nil
}

func blaxelCommand(request *sandbox.ExecOptions) (string, error) {
	if request.Shell != nil {
		return *request.Shell, nil
	}
	if len(request.Argv) == 0 {
		return "", sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "blaxel",
			Operation: "run",
			Field:     sandbox.ExecFieldArgv,
			Message:   "sandbox-kit: set argv or shell",
		}, nil)
	}
	return posixShellJoin(request.Argv), nil
}

// commandTimeout maps Kit ExecOptions.timeout to Blaxel's whole-second timeout.
// Nil and explicit zero become 0 so Blaxel adds no process auto-kill.
func commandTimeout(operation string, value *durationpb.Duration) (time.Duration, error) {
	timeout, ok := sdk.WholeSecondDuration(value)
	if !ok {
		return 0, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "blaxel",
			Operation: operation,
			Field:     sandbox.ExecFieldTimeout,
			Message:   "sandbox-kit blaxel: timeout must be a whole number of seconds",
		}, nil)
	}
	return timeout, nil
}

func posixShellJoin(argv []string) string {
	parts := make([]string, len(argv))
	for i, arg := range argv {
		parts[i] = posixQuote(arg)
	}
	return strings.Join(parts, " ")
}

func posixQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}

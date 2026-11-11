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

const (
	// Process.Wait treats 0 as 1m; only -1 is unlimited.
	blaxelWaitUnlimited = time.Duration(-1)
)

var _ sandbox.SandboxProcess = (*backend)(nil)

func (b *backend) RunCommand(ctx context.Context, id string, request *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
	instance, err := b.processSandbox(ctx, id, "run")
	if err != nil {
		return nil, err
	}
	timeout, err := commandTimeout("run", request.GetTimeout())
	if err != nil {
		return nil, err
	}
	body, err := processRequest(request, timeout)
	if err != nil {
		return nil, err
	}
	started, err := instance.Process.New(ctx, body)
	if err != nil {
		return nil, mapProviderError(err, "run")
	}
	return finishProcess(ctx, instance, started, timeout)
}

// processRequest starts without waitForCompletion; Kit waits via Process.Wait.
// Timeout 0 disables Blaxel process auto-kill.
func processRequest(request *sandbox.ExecOptions, timeout time.Duration) (blaxel.ProcessRequestParam, error) {
	command, err := blaxelCommand(request)
	if err != nil {
		return blaxel.ProcessRequestParam{}, err
	}
	body := blaxel.ProcessRequestParam{
		Command:           command,
		Timeout:           blaxel.Int(int64(timeout / time.Second)),
		WaitForCompletion: blaxel.Bool(false),
	}
	if request.WorkingDirectory != nil {
		body.WorkingDir = blaxel.String(*request.WorkingDirectory)
	}
	if request.Environment != nil {
		body.Env = request.Environment
	}
	return body, nil
}

// finishProcess returns a terminal start as-is; otherwise one native Wait.
func finishProcess(ctx context.Context, instance *blaxel.SandboxInstance, started *blaxel.ProcessResponse, timeout time.Duration) (*sandbox.ExecResult, error) {
	if started == nil || processFinished(started) {
		return execResult(started)
	}
	identifier := processIdentifier(started)
	if identifier == "" {
		return nil, invalidResponse("run", "", "SDK returned no process identifier")
	}
	budget := blaxelWaitUnlimited
	if timeout > 0 {
		budget = timeout
	}
	finished, err := instance.Process.Wait(ctx, identifier, budget, 0)
	if err != nil {
		return nil, mapProviderError(err, "run")
	}
	return execResult(finished)
}

var terminalProcessStatus = map[blaxel.ProcessResponseStatus]struct{}{
	blaxel.ProcessResponseStatusCompleted: {},
	blaxel.ProcessResponseStatusFailed:    {},
	blaxel.ProcessResponseStatusKilled:    {},
	blaxel.ProcessResponseStatusStopped:   {},
}

func processFinished(response *blaxel.ProcessResponse) bool {
	if response == nil {
		return false
	}
	_, ok := terminalProcessStatus[response.Status]
	return ok
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
		return nil, invalidResponse("run", "", "SDK returned no exec result")
	}
	return &sandbox.ExecResult{
		ExitCode: int32(response.ExitCode),
		Stdout:   []byte(response.Stdout),
		Stderr:   []byte(response.Stderr),
	}, nil
}

func (b *backend) processSandbox(ctx context.Context, id, operation string) (*blaxel.SandboxInstance, error) {
	instance, err := b.getInstance(ctx, id, operation)
	if err != nil {
		return nil, err
	}
	if instance == nil || instance.Process == nil {
		return nil, invalidResponse(operation, "", "SDK returned no process API")
	}
	return instance, nil
}

func blaxelCommand(request *sandbox.ExecOptions) (string, error) {
	if request.Shell != nil {
		return *request.Shell, nil
	}
	if len(request.Argv) == 0 {
		// Shared kit wording (no provider infix) — keep identical across providers.
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

func commandTimeout(operation string, value *durationpb.Duration) (time.Duration, error) {
	timeout, ok := sdk.WholeSecondDuration(value)
	if !ok {
		return 0, kitError(kitFault{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Operation: operation,
			Field:     sandbox.ExecFieldTimeout,
			Message:   "timeout must be a whole number of seconds",
		})
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

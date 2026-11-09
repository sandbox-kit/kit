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
		Command:           command,
		WaitForCompletion: blaxel.Bool(true),
	}
	if request.WorkingDirectory != nil {
		body.WorkingDir = blaxel.String(*request.WorkingDirectory)
	}
	if request.Environment != nil {
		body.Env = request.Environment
	}
	if timeout > 0 {
		body.Timeout = blaxel.Int(int64(timeout / time.Second))
	} else {
		// Explicit zero disables Blaxel's keep-alive default kill window.
		body.Timeout = blaxel.Int(0)
	}
	response, err := instance.Process.New(ctx, body)
	if err != nil {
		return nil, mapProviderError(err, "run")
	}
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
// Nil and explicit zero become 0 so Blaxel adds no command limit.
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

package daytona

import (
	"context"
	"strings"
	"time"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	daytonaopts "github.com/daytona/clients/sdk-go/pkg/options"
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
	command, err := daytonaCommand(request)
	if err != nil {
		return nil, err
	}
	timeout, err := commandTimeout("run", request.GetTimeout())
	if err != nil {
		return nil, err
	}
	opts := []func(*daytonaopts.ExecuteCommand){
		daytonaopts.WithExecuteTimeout(timeout),
	}
	if request.WorkingDirectory != nil {
		opts = append(opts, daytonaopts.WithCwd(*request.WorkingDirectory))
	}
	if request.Environment != nil {
		opts = append(opts, daytonaopts.WithCommandEnv(request.Environment))
	}
	response, err := instance.Process.ExecuteCommand(ctx, command, opts...)
	if err != nil {
		return nil, mapProviderError(err, "run")
	}
	if response == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "daytona",
			Operation: "run",
			Message:   "sandbox-kit daytona: SDK returned no exec result",
		}, nil)
	}
	return &sandbox.ExecResult{
		ExitCode: int32(response.ExitCode),
		Stdout:   []byte(response.Result),
		Stderr:   []byte{},
	}, nil
}

func (b *backend) processSandbox(ctx context.Context, id, operation string) (*daytona.Sandbox, error) {
	instance, err := b.nativeSandbox(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	if instance == nil || instance.Process == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "daytona",
			Operation: operation,
			Message:   "sandbox-kit daytona: SDK returned no process API",
		}, nil)
	}
	return instance, nil
}

func daytonaCommand(request *sandbox.ExecOptions) (string, error) {
	if request.Shell != nil {
		return *request.Shell, nil
	}
	if len(request.Argv) == 0 {
		return "", sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "daytona",
			Operation: "run",
			Field:     sandbox.ExecFieldArgv,
			Message:   "sandbox-kit: set argv or shell",
		}, nil)
	}
	return posixShellJoin(request.Argv), nil
}

// commandTimeout maps Kit ExecOptions.timeout to Daytona's command timeout.
// Nil and explicit zero become 0 so Daytona adds no command limit (its omitted
// default is 10 seconds). Non-zero values must be whole seconds.
func commandTimeout(operation string, value *durationpb.Duration) (time.Duration, error) {
	timeout, ok := sdk.WholeSecondDuration(value)
	if !ok {
		return 0, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "daytona",
			Operation: operation,
			Field:     sandbox.ExecFieldTimeout,
			Message:   "sandbox-kit daytona: timeout must be a whole number of seconds",
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

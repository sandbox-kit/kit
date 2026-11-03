package sandbox_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type processBackend struct {
	run func(context.Context, string, *sandbox.ExecOptions) (*sandbox.ExecResult, error)
}

func (*processBackend) Name() string                { return "process" }
func (*processBackend) Close(context.Context) error { return nil }
func (*processBackend) CloseSandbox(context.Context, string) error {
	return nil
}
func (*processBackend) Create(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{Id: "created"}}, nil
}
func (b *processBackend) RunCommand(ctx context.Context, id string, options *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
	return b.run(ctx, id, options)
}

type processProvider struct{ backend sandbox.Backend }

func (*processProvider) Name() string { return "process" }
func (p *processProvider) NewClient(*sandbox.Config) (sandbox.Backend, error) {
	return p.backend, nil
}

func TestProcessNilHandle(t *testing.T) {
	_, err := (*sandbox.Sandbox)(nil).Process().Run(context.Background(), &sandbox.ExecOptions{
		Argv: []string{"echo", "hi"},
	})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindFailedPrecondition {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProcessRejectsMissingCommand(t *testing.T) {
	calls := 0
	backend := &processBackend{run: func(context.Context, string, *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
		calls++
		return &sandbox.ExecResult{}, nil
	}}
	client, err := sandbox.NewClient(&processProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instance.Process().Run(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.ExecFieldArgv || detail.Operation != "run" {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 0 {
		t.Fatal("provider was called without a command")
	}
}

func TestProcessRejectsBothArgvAndShell(t *testing.T) {
	client, err := sandbox.NewClient(&processProvider{backend: &processBackend{
		run: func(context.Context, string, *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
			t.Fatal("provider was called")
			return nil, nil
		},
	}}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Argv:  []string{"echo"},
		Shell: sandbox.Value("echo hi"),
	})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.ExecFieldArgv {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProcessRejectsBlankShell(t *testing.T) {
	client, err := sandbox.NewClient(&processProvider{backend: &processBackend{
		run: func(context.Context, string, *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
			t.Fatal("provider was called")
			return nil, nil
		},
	}}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, shell := range []string{"", "   "} {
		_, err = instance.Process().Run(context.Background(), &sandbox.ExecOptions{
			Shell: sandbox.Value(shell),
		})
		var detail *sandbox.Error
		if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.ExecFieldShell || detail.Operation != "run" {
			t.Fatalf("shell %q: unexpected error: %v", shell, err)
		}
	}
}

func TestProcessRejectsRelativeWorkingDirectory(t *testing.T) {
	client, err := sandbox.NewClient(&processProvider{backend: &processBackend{
		run: func(context.Context, string, *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
			t.Fatal("provider was called")
			return nil, nil
		},
	}}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Argv:             []string{"pwd"},
		WorkingDirectory: sandbox.Value("relative"),
	})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.ExecFieldWorkingDirectory {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProcessRunsArgv(t *testing.T) {
	backend := &processBackend{run: func(_ context.Context, id string, options *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
		if id != "created" || len(options.Argv) != 2 || options.Argv[0] != "echo" || options.Argv[1] != "hello" {
			t.Fatalf("unexpected run %s %#v", id, options)
		}
		return &sandbox.ExecResult{ExitCode: 0, Stdout: []byte("hello\n")}, nil
	}}
	client, err := sandbox.NewClient(&processProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Argv: []string{"echo", "hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.GetExitCode() != 0 || string(result.GetStdout()) != "hello\n" {
		t.Fatalf("unexpected result %#v", result)
	}
}

func TestProcessNonZeroExitIsSuccess(t *testing.T) {
	backend := &processBackend{run: func(context.Context, string, *sandbox.ExecOptions) (*sandbox.ExecResult, error) {
		return &sandbox.ExecResult{ExitCode: 7, Stderr: []byte("boom\n")}, nil
	}}
	client, err := sandbox.NewClient(&processProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Shell: sandbox.Value("exit 7"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.GetExitCode() != 7 {
		t.Fatalf("unexpected exit %d", result.GetExitCode())
	}
}

func TestProcessUnsupportedWithoutCapability(t *testing.T) {
	client, err := sandbox.NewClient(&lifecycleProvider{backend: &lifecycleBackend{}}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instance.Process().Run(context.Background(), &sandbox.ExecOptions{Argv: []string{"true"}})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Operation != "run" {
		t.Fatalf("unexpected error: %v", err)
	}
}

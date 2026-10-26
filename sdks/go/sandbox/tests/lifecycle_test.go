package sandbox_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type lifecycleBackend struct {
	start     func(context.Context, string, *sandbox.StartOptions) (*sandbox.SandboxInfo, error)
	get       func(context.Context, string) (*sandbox.SandboxInfo, error)
	terminate func(context.Context, string, *sandbox.TerminateOptions) (*sandbox.ExitResult, error)
}

func (b *lifecycleBackend) TerminateSandbox(ctx context.Context, id string, options *sandbox.TerminateOptions) (*sandbox.ExitResult, error) {
	if b.terminate == nil {
		return &sandbox.ExitResult{}, nil
	}
	return b.terminate(ctx, id, options)
}

func (*lifecycleBackend) Name() string                { return "lifecycle" }
func (*lifecycleBackend) Close(context.Context) error { return nil }
func (*lifecycleBackend) Create(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{Id: "created"}}, nil
}
func (b *lifecycleBackend) StartSandbox(ctx context.Context, id string, options *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
	return b.start(ctx, id, options)
}
func (b *lifecycleBackend) RefreshSandbox(ctx context.Context, id string) (*sandbox.SandboxInfo, error) {
	return b.get(ctx, id)
}

type lifecycleProvider struct{ backend sandbox.Backend }

func (*lifecycleProvider) Name() string { return "lifecycle" }
func (p *lifecycleProvider) NewClient(*sandbox.Config) (sandbox.Backend, error) {
	return p.backend, nil
}

func TestLifecycleDispatchAndUnsupportedOperation(t *testing.T) {
	called := false
	backend := &lifecycleBackend{start: func(ctx context.Context, id string, options *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		called = true
		if id != "created" || options.GetTimeout() == nil || options.GetTimeout().AsDuration() != time.Second {
			t.Fatal("operation input changed")
		}
		return &sandbox.SandboxInfo{Id: id, ProviderState: sandbox.Value("started")}, nil
	}}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := instance.Start(context.Background(), &sandbox.StartOptions{Timeout: sandbox.Duration(time.Second)}); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("start was not dispatched")
	}
	if !instance.Supports(sandbox.LifecycleCapabilityStart) || !instance.Supports(sandbox.LifecycleCapabilityTerminate) || instance.Supports(sandbox.LifecycleCapabilityPause) {
		t.Fatal("capability discovery is incorrect")
	}
	if instance.Info().GetProviderState() != "started" {
		t.Fatal("lifecycle result did not update handle metadata")
	}
	err = instance.Pause(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Operation != "pause" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTerminateMakesHandleTerminal(t *testing.T) {
	backend := &lifecycleBackend{
		start: func(_ context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
			return &sandbox.SandboxInfo{Id: id}, nil
		},
		terminate: func(context.Context, string, *sandbox.TerminateOptions) (*sandbox.ExitResult, error) {
			return &sandbox.ExitResult{ExitCode: sandbox.Value(int32(0))}, nil
		},
	}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := instance.Terminate(context.Background(), &sandbox.TerminateOptions{Wait: sandbox.Value(true)}); err != nil {
		t.Fatal(err)
	}
	if !instance.Info().GetTerminal() || instance.Info().GetState() != sandbox.SandboxStateTerminated || len(instance.Capabilities()) != 0 {
		t.Fatal("terminal state was not retained")
	}
	err = instance.Start(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindFailedPrecondition {
		t.Fatalf("terminal operation was dispatched: %v", err)
	}
}

func TestLifecycleOperationsSerializeAndRespectWaitingContext(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	backend := &lifecycleBackend{start: func(_ context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		select {
		case entered <- struct{}{}:
			<-release
		default:
			t.Fatal("concurrent lifecycle call reached provider")
		}
		return &sandbox.SandboxInfo{Id: id}, nil
	}}
	client, _ := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	instance, _ := client.Create(context.Background(), nil)
	done := make(chan error, 1)
	go func() { done <- instance.Start(context.Background(), nil) }()
	<-entered
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if err := instance.Start(ctx, nil); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("waiting operation ignored context: %v", err)
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestGetAndRefreshOwnMetadata(t *testing.T) {
	name := "first"
	backend := &lifecycleBackend{
		start: func(_ context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
			return &sandbox.SandboxInfo{Id: id}, nil
		},
		get: func(_ context.Context, id string) (*sandbox.SandboxInfo, error) {
			return &sandbox.SandboxInfo{Id: id, Name: sandbox.Value(name)}, nil
		},
	}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Get(context.Background(), "existing")
	if err != nil {
		t.Fatal(err)
	}
	name = "second"
	if instance.Info().GetName() != "first" {
		t.Fatal("handle did not own retrieved metadata")
	}
	if err := instance.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	if instance.Info().GetName() != "second" {
		t.Fatal("refresh did not replace metadata")
	}
}

func TestLifecycleInheritsClientTimeout(t *testing.T) {
	backend := &lifecycleBackend{start: func(ctx context.Context, _ string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{Timeout: sandbox.Duration(time.Millisecond)})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Start(context.Background(), nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout not inherited: %v", err)
	}
}

func TestDeleteMarksHandleTerminalBeforeAnotherOperation(t *testing.T) {
	backend := &deleteBackend{lifecycleBackend: lifecycleBackend{start: func(context.Context, string, *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		t.Fatal("start reached the provider after delete")
		return nil, nil
	}}}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := instance.Delete(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if !instance.Info().GetTerminal() || instance.Info().GetState() != sandbox.SandboxStateDeleted {
		t.Fatal("delete did not retain a terminal handle")
	}
	err = instance.Start(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindFailedPrecondition {
		t.Fatalf("operation ran after delete: %v", err)
	}
}

func TestCloseWaitsForTheCurrentOperation(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	var calls int
	backend := &lifecycleBackend{start: func(_ context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		calls++
		if calls == 1 {
			close(entered)
			<-release
		}
		return &sandbox.SandboxInfo{Id: id}, nil
	}}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- instance.Start(context.Background(), nil) }()
	<-entered
	closed := make(chan error, 1)
	go func() { closed <- instance.Close(context.Background()) }()
	select {
	case err := <-closed:
		t.Fatalf("close returned while start still held the handle: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if err := <-closed; err != nil {
		t.Fatal(err)
	}
	err = instance.Start(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindFailedPrecondition || calls != 1 {
		t.Fatalf("operation ran after close: %v calls=%d", err, calls)
	}
}

func TestClosedHandleAndClientRejectLifecycleOperations(t *testing.T) {
	backend := &lifecycleBackend{start: func(_ context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		return &sandbox.SandboxInfo{Id: id}, nil
	}}
	client, _ := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	instance, _ := client.Create(context.Background(), nil)
	if err := instance.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := instance.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(instance.Capabilities()) != 0 {
		t.Fatal("closed handle still advertises capabilities")
	}
	err := instance.Start(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindFailedPrecondition {
		t.Fatalf("closed handle accepted operation: %v", err)
	}

	client2, _ := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	instance2, _ := client2.Create(context.Background(), nil)
	if err := client2.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(instance2.Capabilities()) != 0 {
		t.Fatal("handle whose client is closed still advertises capabilities")
	}
	err = instance2.Start(context.Background(), nil)
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindFailedPrecondition {
		t.Fatalf("handle outlived its client: %v", err)
	}
}

func TestGetRejectsBlankID(t *testing.T) {
	backend := &lifecycleBackend{get: func(context.Context, string) (*sandbox.SandboxInfo, error) {
		t.Fatal("blank id reached the provider")
		return nil, nil
	}}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		id   string
	}{
		{"empty", ""},
		{"spaces", "   "},
		{"tab and newline", "\t\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.Get(context.Background(), tc.id)
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.InfoFieldId {
				t.Fatalf("blank id accepted: %v", err)
			}
		})
	}
}

func TestLifecycleStampsInheritedTimeoutWithoutMutatingCaller(t *testing.T) {
	var seen time.Duration
	backend := &lifecycleBackend{start: func(_ context.Context, id string, options *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		if options.GetTimeout() == nil {
			t.Fatal("inherited timeout was not forwarded")
		}
		seen = options.GetTimeout().AsDuration()
		return &sandbox.SandboxInfo{Id: id}, nil
	}}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{Timeout: sandbox.Duration(2 * time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := instance.Start(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if seen != 2*time.Second {
		t.Fatalf("provider timeout = %s", seen)
	}
}

type deleteBackend struct {
	lifecycleBackend
	del func(context.Context, string, *sandbox.DeleteOptions) (*sandbox.SandboxInfo, error)
}

func (b *deleteBackend) DeleteSandbox(ctx context.Context, id string, options *sandbox.DeleteOptions) (*sandbox.SandboxInfo, error) {
	if b.del == nil {
		return &sandbox.SandboxInfo{Id: id}, nil
	}
	return b.del(ctx, id, options)
}

type stopBackend struct {
	lifecycleBackend
	stop func(context.Context, string, *sandbox.StopOptions) (*sandbox.SandboxInfo, error)
}

func (b *stopBackend) StopSandbox(ctx context.Context, id string, options *sandbox.StopOptions) (*sandbox.SandboxInfo, error) {
	return b.stop(ctx, id, options)
}

func TestExplicitZeroTimeoutOverridesClientTimeout(t *testing.T) {
	backend := &stopBackend{stop: func(ctx context.Context, id string, options *sandbox.StopOptions) (*sandbox.SandboxInfo, error) {
		if _, hasDeadline := ctx.Deadline(); hasDeadline {
			t.Fatal("explicit zero added a deadline")
		}
		if options.GetTimeout() == nil || options.GetTimeout().AsDuration() != 0 || !options.GetForce() {
			t.Fatalf("provider options = timeout %v force %v", options.GetTimeout(), options.GetForce())
		}
		return &sandbox.SandboxInfo{Id: id}, nil
	}}
	client, err := sandbox.NewClient(&lifecycleProvider{backend: backend}, &sandbox.Config{Timeout: sandbox.Duration(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	caller := &sandbox.StopOptions{Timeout: sandbox.Duration(0), Force: sandbox.Value(true)}
	timeout := caller.Timeout
	if err := instance.Stop(context.Background(), caller); err != nil {
		t.Fatal(err)
	}
	if caller.Timeout != timeout || !caller.GetForce() {
		t.Fatal("caller options were mutated")
	}
}

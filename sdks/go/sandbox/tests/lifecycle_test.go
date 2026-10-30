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
	closes    int
}

func (b *lifecycleBackend) TerminateSandbox(ctx context.Context, id string, options *sandbox.TerminateOptions) (*sandbox.ExitResult, error) {
	if b.terminate == nil {
		return &sandbox.ExitResult{}, nil
	}
	return b.terminate(ctx, id, options)
}

func (*lifecycleBackend) Name() string                { return "lifecycle" }
func (*lifecycleBackend) Close(context.Context) error { return nil }
func (b *lifecycleBackend) CloseSandbox(context.Context, string) error {
	b.closes++
	return nil
}
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

func TestTerminateLeavesTheNextCallToTheProvider(t *testing.T) {
	var starts int
	backend := &lifecycleBackend{
		start: func(_ context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
			starts++
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
	if instance.Info().GetTerminal() || !instance.Supports(sandbox.LifecycleCapabilityStart) {
		t.Fatal("terminate changed the handle into a local terminal state")
	}
	if err := instance.Start(context.Background(), nil); err != nil || starts != 1 {
		t.Fatalf("start did not reach the provider: %v calls=%d", err, starts)
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

func TestDeleteKeepsProviderMetadataAndCallsAgain(t *testing.T) {
	var starts int
	backend := &deleteBackend{
		del: func(_ context.Context, id string, _ *sandbox.DeleteOptions) (*sandbox.SandboxInfo, error) {
			return &sandbox.SandboxInfo{Id: id, State: sandbox.SandboxStateDeleted, Terminal: sandbox.Value(true)}, nil
		},
		lifecycleBackend: lifecycleBackend{start: func(_ context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
			starts++
			return &sandbox.SandboxInfo{Id: id}, nil
		}},
	}
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
	if !instance.Info().GetTerminal() || instance.Info().GetState() != sandbox.SandboxStateDeleted || !instance.Supports(sandbox.LifecycleCapabilityStart) {
		t.Fatal("delete did not keep the provider metadata")
	}
	if err := instance.Start(context.Background(), nil); err != nil || starts != 1 {
		t.Fatalf("start did not reach the provider: %v calls=%d", err, starts)
	}
}

func TestCloseLeavesTheNextCallToTheProvider(t *testing.T) {
	var starts int
	backend := &lifecycleBackend{start: func(_ context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		starts++
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
	if err := instance.Close(context.Background()); err != nil || backend.closes != 1 {
		t.Fatalf("handle close did not reach the provider: %v calls=%d", err, backend.closes)
	}
	if !instance.Supports(sandbox.LifecycleCapabilityStart) {
		t.Fatal("handle close cleared provider capabilities")
	}
	if err := instance.Start(context.Background(), nil); err != nil || starts != 1 {
		t.Fatalf("start did not reach the provider after handle close: %v calls=%d", err, starts)
	}
	if err := client.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !instance.Supports(sandbox.LifecycleCapabilityStart) {
		t.Fatal("client close cleared provider capabilities")
	}
	if err := instance.Start(context.Background(), nil); err != nil || starts != 2 {
		t.Fatalf("start did not reach the provider after client close: %v calls=%d", err, starts)
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

func TestOmittedLifecycleTimeoutStaysOmitted(t *testing.T) {
	backend := &lifecycleBackend{start: func(ctx context.Context, id string, options *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("client deadline was not applied")
		}
		if options.GetTimeout() != nil {
			t.Fatalf("provider timeout = %s", options.GetTimeout().AsDuration())
		}
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
	caller := &sandbox.StartOptions{}
	if err := instance.Start(context.Background(), caller); err != nil {
		t.Fatal(err)
	}
	if caller.Timeout != nil {
		t.Fatal("caller options were mutated")
	}

	stopper := &stopBackend{stop: func(ctx context.Context, id string, options *sandbox.StopOptions) (*sandbox.SandboxInfo, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("client deadline was not applied")
		}
		if options.GetTimeout() != nil || !options.GetForce() {
			t.Fatalf("provider options = timeout %v force %v", options.GetTimeout(), options.GetForce())
		}
		return &sandbox.SandboxInfo{Id: id}, nil
	}}
	stopClient, err := sandbox.NewClient(&lifecycleProvider{backend: stopper}, &sandbox.Config{Timeout: sandbox.Duration(2 * time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	stopped, err := stopClient.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	stopCaller := &sandbox.StopOptions{Force: sandbox.Value(true)}
	if err := stopped.Stop(context.Background(), stopCaller); err != nil {
		t.Fatal(err)
	}
	if stopCaller.Timeout != nil || !stopCaller.GetForce() {
		t.Fatal("caller options were mutated")
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

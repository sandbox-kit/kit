package daytona

import (
	"context"
	"time"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

var _ sandbox.SandboxRefresher = (*backend)(nil)
var _ sandbox.SandboxStarter = (*backend)(nil)
var _ sandbox.SandboxStopper = (*backend)(nil)
var _ sandbox.SandboxPauser = (*backend)(nil)
var _ sandbox.SandboxArchiver = (*backend)(nil)
var _ sandbox.SandboxDeleter = (*backend)(nil)

func (*backend) SandboxCapabilities(info *sandbox.SandboxInfo) []sandbox.LifecycleCapability {
	capabilities := []sandbox.LifecycleCapability{sandbox.LifecycleCapabilityRefresh, sandbox.LifecycleCapabilityStart, sandbox.LifecycleCapabilityStop, sandbox.LifecycleCapabilityDelete}
	if info == nil || info.GetTerminal() {
		return nil
	}
	gpu := false
	if metadata, err := sandbox.MetadataObjectAsMap(info.GetProviderMetadata()); err == nil {
		if count, ok := metadata["gpu_count"].(int64); ok {
			gpu = count > 0
		}
	}
	isolation := info.GetIsolation()
	if !gpu && isolation != nil && (isolation.Kind == sandbox.IsolationKindLinuxVM || isolation.GetProviderKind() == "windows") {
		capabilities = append(capabilities, sandbox.LifecycleCapabilityPause)
	}
	if !gpu && isolation != nil && isolation.Kind == sandbox.IsolationKindContainer {
		capabilities = append(capabilities, sandbox.LifecycleCapabilityArchive)
	}
	return capabilities
}

func (b *backend) nativeSandbox(ctx context.Context, id string) (*daytona.Sandbox, error) {
	return b.client.Get(ctx, id)
}

func (b *backend) withSandbox(ctx context.Context, id, operation string, mutate func(context.Context, *daytona.Sandbox) error) (*sandbox.SandboxInfo, error) {
	instance, err := b.nativeSandbox(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	if err := mutate(ctx, instance); err != nil {
		return nil, err
	}
	info, err := sandboxInfo(instance)
	return info, err
}

func (b *backend) StartSandbox(ctx context.Context, id string, options *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
	timeout := nativeTimeout(options.GetTimeout())
	return b.withSandbox(ctx, id, "start", func(ctx context.Context, instance *daytona.Sandbox) error {
		if err := instance.StartWithTimeout(ctx, timeout); err != nil {
			return mapProviderError(err, "start")
		}
		return nil
	})
}

func (b *backend) StopSandbox(ctx context.Context, id string, options *sandbox.StopOptions) (*sandbox.SandboxInfo, error) {
	timeout := nativeTimeout(options.GetTimeout())
	return b.withSandbox(ctx, id, "stop", func(ctx context.Context, instance *daytona.Sandbox) error {
		force := options != nil && options.GetForce()
		if err := instance.StopWithTimeout(ctx, timeout, force); err != nil {
			return mapProviderError(err, "stop")
		}
		return nil
	})
}

func (b *backend) PauseSandbox(ctx context.Context, id string, options *sandbox.PauseOptions) (*sandbox.SandboxInfo, error) {
	timeout := nativeTimeout(options.GetTimeout())
	return b.withSandbox(ctx, id, "pause", func(ctx context.Context, instance *daytona.Sandbox) error {
		if instance.SandboxClass == nil || (*instance.SandboxClass != types.SandboxClassLinuxVM && string(*instance.SandboxClass) != "windows") {
			return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "pause", Field: sandbox.InfoFieldIsolation, Message: "sandbox-kit daytona: pause requires a Linux VM or Windows sandbox"}, nil)
		}
		if err := instance.PauseWithTimeout(ctx, timeout); err != nil {
			return mapProviderError(err, "pause")
		}
		return nil
	})
}

func (b *backend) ArchiveSandbox(ctx context.Context, id string, _ *sandbox.ArchiveOptions) (*sandbox.SandboxInfo, error) {
	return b.withSandbox(ctx, id, "archive", func(ctx context.Context, instance *daytona.Sandbox) error {
		if instance.SandboxClass == nil || *instance.SandboxClass != types.SandboxClassContainer {
			return sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "archive", Field: sandbox.InfoFieldIsolation, Message: "sandbox-kit daytona: archive requires a container sandbox"}, nil)
		}
		if err := instance.Archive(ctx); err != nil {
			return mapProviderError(err, "archive")
		}
		return nil
	})
}

func (b *backend) DeleteSandbox(ctx context.Context, id string, options *sandbox.DeleteOptions) (*sandbox.SandboxInfo, error) {
	timeout := nativeTimeout(options.GetTimeout())
	return b.withSandbox(ctx, id, "delete", func(ctx context.Context, instance *daytona.Sandbox) error {
		if options != nil && options.GetWait() {
			if err := instance.DeleteAndWait(ctx, timeout); err != nil {
				return mapProviderError(err, "delete")
			}
			return nil
		}
		if err := instance.DeleteWithTimeout(ctx, timeout); err != nil {
			return mapProviderError(err, "delete")
		}
		return nil
	})
}

func nativeTimeout(value *durationpb.Duration) time.Duration {
	timeout := sdk.OptionalDuration(value)
	if timeout == nil {
		return 0
	}
	return *timeout
}

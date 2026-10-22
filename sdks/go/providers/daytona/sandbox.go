// This file turns a shared creation request into a Daytona SDK create call.
// Generated mappings cover field and policy conversion. This file chooses the
// snapshot or image source and applies the mapped options.
//
// Example:
//
//	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
//	instance, err := client.Create(ctx, nil)

package daytona

import (
	"context"
	"math"
	"strings"
	"time"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/options"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
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

// createPlan is the Daytona source and option list for one creation request.
type createPlan struct {
	// params is the snapshot or image argument accepted by the SDK create method.
	params any
	// options are the mapped creation settings, including policies and resources.
	options []func(*options.CreateSandbox)
}

// create calls the Daytona SDK and returns shared metadata.
// A nil request uses the default snapshot and provider defaults.
func (a *backend) create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	if request == nil {
		request = &sandbox.CreateOptions{}
	}
	plan, err := planCreate(request)
	if err != nil {
		return nil, err
	}
	instance, err := a.client.Create(ctx, plan.params, plan.options...)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, sandbox.NewError(sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "daytona",
			Operation: "create",
			Field:     sandbox.InfoFieldId,
			Message:   "sandbox-kit daytona: SDK returned no sandbox",
		}, nil)
	}
	// Use the same complete resolver as Get and lifecycle responses so a newly
	// created handle immediately has the class and metadata needed for accurate
	// capability decisions.
	info, err := daytonaSandboxInfo(instance)
	if err != nil {
		return nil, err
	}
	return &sandbox.CreateResult{Sandbox: &info}, nil
}

// RefreshSandbox retrieves current metadata for an existing Daytona sandbox.
func (a *backend) RefreshSandbox(ctx context.Context, id string) (*sandbox.SandboxInfo, error) {
	instance, err := a.client.Get(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, "refresh")
	}
	info, err := daytonaSandboxInfo(instance)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func daytonaSandboxInfo(instance *daytona.Sandbox) (sandbox.SandboxInfo, error) {
	info, err := mapSandboxInfo(instance)
	if err != nil {
		return info, err
	}
	if err := applyDaytonaResources(&info, instance); err != nil {
		return info, err
	}
	applyDaytonaState(&info, instance)
	applyDaytonaIsolation(&info, instance)
	if err := applyDaytonaProviderMetadata(&info, instance); err != nil {
		return info, err
	}
	return info, nil
}

func applyDaytonaResources(info *sandbox.SandboxInfo, instance *daytona.Sandbox) error {
	resources, err := mapAllocatedResources(instance)
	if err != nil {
		return err
	}
	if resources.CpuCores != nil {
		info.Resources = &resources
		info.Origins["resources"] = sandbox.ValueOriginProvider
	}
	return nil
}

func applyDaytonaIsolation(info *sandbox.SandboxInfo, instance *daytona.Sandbox) {
	if instance.SandboxClass == nil {
		return
	}
	isolation := &sandbox.IsolationConfig{}
	switch *instance.SandboxClass {
	case types.SandboxClassContainer:
		isolation.Kind = sandbox.IsolationKindContainer
	case types.SandboxClassLinuxVM:
		isolation.Kind = sandbox.IsolationKindLinuxVM
	default:
		isolation.Kind = sandbox.IsolationKindProviderSpecific
		isolation.ProviderKind = sandbox.Value(string(*instance.SandboxClass))
	}
	info.Isolation = isolation
	info.Origins["isolation"] = sandbox.ValueOriginProvider
}

func applyDaytonaProviderMetadata(info *sandbox.SandboxInfo, instance *daytona.Sandbox) error {
	metadata := map[string]any{}
	if instance.Gpu > 0 {
		metadata["gpu_count"] = int64(instance.Gpu)
	}
	if instance.GpuType != nil {
		metadata["gpu_type"] = string(*instance.GpuType)
	}
	if instance.DesiredState != nil {
		metadata["desired_state"] = string(*instance.DesiredState)
	}
	if instance.ErrorReason != nil {
		metadata["error_reason"] = *instance.ErrorReason
	}
	if instance.Recoverable != nil {
		metadata["recoverable"] = *instance.Recoverable
	}
	if instance.Spot {
		metadata["spot"] = true
	}
	if instance.SpotEvictedAt != nil {
		metadata["spot_evicted_at"] = *instance.SpotEvictedAt
	}
	if instance.BackupState != nil {
		metadata["backup_state"] = *instance.BackupState
	}
	if instance.BackupCreatedAt != nil {
		metadata["backup_created_at"] = *instance.BackupCreatedAt
	}
	if instance.UpdatedAt != nil {
		metadata["updated_at"] = *instance.UpdatedAt
	}
	if instance.LastActivityAt != nil {
		metadata["last_activity_at"] = *instance.LastActivityAt
	}
	if instance.AutoDestroyAt != nil {
		metadata["auto_destroy_at"] = *instance.AutoDestroyAt
	}
	if instance.QueueTimedOutAt != nil {
		metadata["queue_timed_out_at"] = *instance.QueueTimedOutAt
	}
	if instance.Snapshot != nil {
		metadata["snapshot"] = *instance.Snapshot
	}
	if instance.WarmPoolId != nil {
		metadata["warm_pool_id"] = *instance.WarmPoolId
	}
	metadata["auto_stop_interval_minutes"] = int64(instance.AutoStopInterval)
	metadata["auto_pause_interval_minutes"] = int64(instance.AutoPauseInterval)
	metadata["auto_archive_interval_minutes"] = int64(instance.AutoArchiveInterval)
	metadata["auto_delete_interval_minutes"] = int64(instance.AutoDeleteInterval)
	if len(metadata) == 0 {
		return nil
	}
	object, err := sandbox.MetadataObjectFromMap(metadata)
	if err != nil {
		return err
	}
	info.ProviderMetadata = object
	info.Origins[sandbox.InfoFieldProviderMetadata] = sandbox.ValueOriginProvider
	return nil
}

func applyDaytonaState(info *sandbox.SandboxInfo, instance *daytona.Sandbox) {
	switch string(instance.State) {
	case "creating", "restoring", "pending_build", "building_snapshot", "pulling_snapshot", "forking":
		info.State = sandbox.SandboxStateQueued
	case "starting", "resuming":
		info.State = sandbox.SandboxStateStarting
	case "started":
		info.State = sandbox.SandboxStateRunning
	case "stopping":
		info.State = sandbox.SandboxStateStopping
	case "stopped":
		info.State = sandbox.SandboxStateStopped
	case "pausing":
		info.State = sandbox.SandboxStatePausing
	case "paused":
		info.State = sandbox.SandboxStatePaused
	case "archiving":
		info.State = sandbox.SandboxStateArchiving
	case "archived":
		info.State = sandbox.SandboxStateArchived
	case "destroying":
		info.State = sandbox.SandboxStateDeleting
	case "destroyed":
		info.State = sandbox.SandboxStateDeleted
		info.Terminal = sandbox.Value(true)
	case "error", "build_failed":
		info.State = sandbox.SandboxStateError
	}
}

func (a *backend) daytonaSandbox(ctx context.Context, id string) (*daytona.Sandbox, error) {
	return a.client.Get(ctx, id)
}

func (a *backend) withSandbox(ctx context.Context, id, operation string, mutate func(context.Context, *daytona.Sandbox) error) (*sandbox.SandboxInfo, error) {
	instance, err := a.daytonaSandbox(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	if err := mutate(ctx, instance); err != nil {
		return nil, err
	}
	info, err := daytonaSandboxInfo(instance)
	return &info, err
}

func (a *backend) StartSandbox(ctx context.Context, id string, _ *sandbox.StartOptions) (*sandbox.SandboxInfo, error) {
	return a.withSandbox(ctx, id, "start", func(ctx context.Context, instance *daytona.Sandbox) error {
		if err := instance.StartWithTimeout(ctx, 0); err != nil {
			return mapProviderError(err, "start")
		}
		return nil
	})
}

func (a *backend) StopSandbox(ctx context.Context, id string, options *sandbox.StopOptions) (*sandbox.SandboxInfo, error) {
	return a.withSandbox(ctx, id, "stop", func(ctx context.Context, instance *daytona.Sandbox) error {
		force := options != nil && options.GetForce()
		if err := instance.StopWithTimeout(ctx, 0, force); err != nil {
			return mapProviderError(err, "stop")
		}
		return nil
	})
}

func (a *backend) PauseSandbox(ctx context.Context, id string, _ *sandbox.PauseOptions) (*sandbox.SandboxInfo, error) {
	return a.withSandbox(ctx, id, "pause", func(ctx context.Context, instance *daytona.Sandbox) error {
		if instance.SandboxClass == nil || (*instance.SandboxClass != types.SandboxClassLinuxVM && string(*instance.SandboxClass) != "windows") {
			return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "pause", Field: sandbox.InfoFieldIsolation, Message: "sandbox-kit daytona: pause requires a Linux VM or Windows sandbox"}, nil)
		}
		if err := instance.PauseWithTimeout(ctx, 0); err != nil {
			return mapProviderError(err, "pause")
		}
		return nil
	})
}

func (a *backend) ArchiveSandbox(ctx context.Context, id string, _ *sandbox.ArchiveOptions) (*sandbox.SandboxInfo, error) {
	return a.withSandbox(ctx, id, "archive", func(ctx context.Context, instance *daytona.Sandbox) error {
		if instance.SandboxClass == nil || *instance.SandboxClass != types.SandboxClassContainer {
			return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "archive", Field: sandbox.InfoFieldIsolation, Message: "sandbox-kit daytona: archive requires a container sandbox"}, nil)
		}
		if err := instance.Archive(ctx); err != nil {
			return mapProviderError(err, "archive")
		}
		return nil
	})
}

func (a *backend) DeleteSandbox(ctx context.Context, id string, options *sandbox.DeleteOptions) (*sandbox.SandboxInfo, error) {
	return a.withSandbox(ctx, id, "delete", func(ctx context.Context, instance *daytona.Sandbox) error {
		if options != nil && options.GetWait() {
			if err := instance.DeleteAndWait(ctx, 0); err != nil {
				return mapProviderError(err, "delete")
			}
			return nil
		}
		if err := instance.DeleteWithTimeout(ctx, 0); err != nil {
			return mapProviderError(err, "delete")
		}
		return nil
	})
}

// planCreate checks the request and builds SDK arguments without contacting Daytona.
func planCreate(request *sandbox.CreateOptions) (createPlan, error) {
	var plan createPlan
	if err := validateProviderCreate(request); err != nil {
		return plan, err
	}
	base, rest, err := mapCommonCreateFields(request)
	if err != nil {
		return plan, err
	}
	if err := planCreateRuntime(request, &base, &rest); err != nil {
		return plan, err
	}
	if err := planCreateIsolation(request, &base, &rest); err != nil {
		return plan, err
	}
	if err := planCreateSecurity(request, &base, &rest); err != nil {
		return plan, err
	}
	if err := planCreateNetwork(request, &base, &rest); err != nil {
		return plan, err
	}
	if err := planCreateLifetime(request, &base, &rest); err != nil {
		return plan, err
	}
	if err := planCreateProvisioning(request, &base, &rest, &plan); err != nil {
		return plan, err
	}
	if err := planCreateObservability(request, &base, &rest); err != nil {
		return plan, err
	}
	if err := planCreateSource(request, &base, &rest, &plan); err != nil {
		return plan, err
	}
	if err := sandbox.RejectUnmapped("daytona", &rest); err != nil {
		return plan, err
	}
	return plan, nil
}

func planCreateRuntime(request *sandbox.CreateOptions, base *types.SandboxBaseParams, rest *sandbox.CreateOptions) error {
	if request.Runtime == nil {
		return nil
	}
	config := *request.Runtime
	base.User = config.GetUser()
	base.Language = types.CodeLanguage(config.GetLanguage())
	config.User, config.Language = nil, nil
	if err := sandbox.RejectUnmapped("daytona runtime", &config); err != nil {
		return err
	}
	rest.Runtime = nil
	return nil
}

func planCreateIsolation(request *sandbox.CreateOptions, base *types.SandboxBaseParams, rest *sandbox.CreateOptions) error {
	if request.Isolation == nil {
		return nil
	}
	config := *request.Isolation
	base.Kvm = config.GetNestedVirtualization()
	config.NestedVirtualization = nil
	if err := sandbox.RejectUnmapped("daytona isolation (class comes from snapshot)", &config); err != nil {
		return err
	}
	rest.Isolation = nil
	return nil
}

func planCreateSecurity(request *sandbox.CreateOptions, base *types.SandboxBaseParams, rest *sandbox.CreateOptions) error {
	if request.Security == nil {
		return nil
	}
	config := *request.Security
	base.Public = config.GetPublicAccess()
	config.PublicAccess = nil
	if len(config.Secrets) > 0 {
		base.Secrets = map[string]string{}
	}
	for _, secret := range config.Secrets {
		if secret.GetEnvironmentVariable() == "" {
			return sandbox.NewError(sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create",
				Field: sandbox.CreateFieldSecuritySecretsEnvironmentVariable,
				Message: "sandbox-kit daytona: secret environment variable is required",
			}, nil)
		}
		if secret.Injection != sandbox.SecretInjectionDefault && secret.Injection != sandbox.SecretInjectionEgressPlaceholder {
			return sandbox.NewError(sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create",
				Field: sandbox.CreateFieldSecuritySecretsInjection,
				Message: "sandbox-kit daytona: secret injection must use egress placeholders",
			}, nil)
		}
		if _, exists := base.Secrets[secret.GetEnvironmentVariable()]; exists {
			return sandbox.NewError(sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create",
				Field: sandbox.CreateFieldSecuritySecretsEnvironmentVariable,
				Message: "sandbox-kit daytona: duplicate secret environment variable",
			}, nil)
		}
		base.Secrets[secret.GetEnvironmentVariable()] = secret.Reference
	}
	config.Secrets = nil
	if err := sandbox.RejectUnmapped("daytona security", &config); err != nil {
		return err
	}
	rest.Security = nil
	return nil
}

func planCreateNetwork(request *sandbox.CreateOptions, base *types.SandboxBaseParams, rest *sandbox.CreateOptions) error {
	if request.Network == nil {
		return nil
	}
	config := *request.Network
	if config.LinkedSandbox != nil && (request.Lifetime == nil || !request.Lifetime.GetEphemeral()) {
		return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create", Field: sandbox.CreateFieldNetworkLinkedSandbox, Message: "sandbox-kit daytona: linked sandboxes must be ephemeral"}, nil)
	}
	switch config.Egress {
	case sandbox.EgressModeDefault, sandbox.EgressModeAllowAll:
	case sandbox.EgressModeBlockAll:
		base.NetworkBlockAll = true
	case sandbox.EgressModeRestricted:
	default:
		return sandbox.NewError(sandbox.ErrorInfo{
			Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create",
			Field: sandbox.CreateFieldNetworkEgress,
			Message: "sandbox-kit daytona: unknown egress mode",
		}, nil)
	}
	config.Egress = 0
	if config.OutboundCidrs != nil {
		if len(config.OutboundCidrs.Entries) == 0 {
			return sandbox.NewError(sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create",
				Field: sandbox.CreateFieldNetworkOutboundCidrs,
				Message: "sandbox-kit daytona: empty CIDR allowlist mapping is not verified",
			}, nil)
		}
		value := strings.Join(config.OutboundCidrs.Entries, ",")
		base.NetworkAllowList = &value
	}
	config.OutboundCidrs = nil
	if config.OutboundDomains != nil {
		if len(config.OutboundDomains.Entries) == 0 {
			return sandbox.NewError(sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create",
				Field: sandbox.CreateFieldNetworkOutboundDomains,
				Message: "sandbox-kit daytona: empty domain allowlist mapping is not verified",
			}, nil)
		}
		value := strings.Join(config.OutboundDomains.Entries, ",")
		base.DomainAllowList = &value
	}
	config.OutboundDomains = nil
	base.LinkedSandbox = config.GetLinkedSandbox()
	base.OutboundProxyUrl = config.OutboundProxyUrl
	config.LinkedSandbox, config.OutboundProxyUrl = nil, nil
	if err := sandbox.RejectUnmapped("daytona network", &config); err != nil {
		return err
	}
	rest.Network = nil
	return nil
}

func planCreateLifetime(request *sandbox.CreateOptions, base *types.SandboxBaseParams, rest *sandbox.CreateOptions) error {
	if request.Lifetime == nil {
		return nil
	}
	config := *request.Lifetime
	base.Ephemeral = config.GetEphemeral()
	config.Ephemeral = nil
	if config.MaximumLifetime != nil {
		minutes, err := durationMinutesAllowZero(config.MaximumLifetime.AsDuration())
		if err != nil {
			return err
		}
		base.TtlMinutes = &minutes
	}
	config.MaximumLifetime = nil
	mapped, remaining, err := mapLifetimePolicies(&config)
	if err != nil {
		return err
	}
	base.AutoStopInterval, base.AutoPauseInterval = mapped.AutoStopInterval, mapped.AutoPauseInterval
	base.AutoArchiveInterval, base.AutoDeleteInterval = mapped.AutoArchiveInterval, mapped.AutoDeleteInterval
	config = remaining
	if err := sandbox.RejectUnmapped("daytona lifetime", &config); err != nil {
		return err
	}
	rest.Lifetime = nil
	return nil
}

func planCreateProvisioning(request *sandbox.CreateOptions, base *types.SandboxBaseParams, rest *sandbox.CreateOptions, plan *createPlan) error {
	if request.Provisioning == nil {
		return nil
	}
	config := *request.Provisioning
	if timeout := sandbox.OptionalDuration(config.Timeout); timeout != nil && *timeout > 0 {
		plan.options = append(plan.options, options.WithTimeout(*timeout))
	}
	config.Timeout = nil
	if config.QueueTimeout != nil {
		minutes, err := durationMinutesRequirePositive(config.QueueTimeout.AsDuration())
		if err != nil {
			return err
		}
		base.QueueTimeout = &minutes
	}
	config.QueueTimeout = nil
	switch config.WaitFor {
	case sandbox.WaitConditionDefault:
	case sandbox.WaitConditionSubmitted:
		plan.options = append(plan.options, options.WithWaitForStart(false))
	case sandbox.WaitConditionStarted:
		plan.options = append(plan.options, options.WithWaitForStart(true))
	default:
		return sandbox.NewError(sandbox.ErrorInfo{
			Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create",
			Field: sandbox.CreateFieldProvisioningWaitFor,
			Message: "sandbox-kit daytona: requested wait condition is unsupported",
		}, nil)
	}
	config.WaitFor = 0
	if err := sandbox.RejectUnmapped("daytona creation", &config); err != nil {
		return err
	}
	rest.Provisioning = nil
	return nil
}

func planCreateObservability(request *sandbox.CreateOptions, base *types.SandboxBaseParams, rest *sandbox.CreateOptions) error {
	if request.Observability == nil {
		return nil
	}
	config := *request.Observability
	base.OtelEndpointOverride = config.TelemetryEndpoint
	config.TelemetryEndpoint = nil
	if err := sandbox.RejectUnmapped("daytona observability", &config); err != nil {
		return err
	}
	rest.Observability = nil
	return nil
}

func planCreateSource(request *sandbox.CreateOptions, base *types.SandboxBaseParams, rest *sandbox.CreateOptions, plan *createPlan) error {
	source := request.GetSource()
	if source.GetWarmPool() != nil {
		return sandbox.NewError(sandbox.ErrorInfo{
			Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create",
			Field: sandbox.CreateFieldSourceWarmPool,
			Message: "sandbox-kit daytona: explicit warm-pool creation mapping is unavailable",
		}, nil)
	}
	if source.GetImage() != nil {
		params := types.ImageParams{SandboxBaseParams: *base, Image: source.GetImage().Reference}
		if request.Resources != nil {
			resources, remaining, err := mapResources(request.Resources)
			if err != nil {
				return err
			}
			if err := sandbox.RejectUnmapped("daytona resources", &remaining); err != nil {
				return err
			}
			params.Resources = &resources
			rest.Resources = nil
		}
		plan.params = params
	} else {
		params := types.SnapshotParams{SandboxBaseParams: *base}
		if snapshot := source.GetSnapshot(); snapshot != nil {
			if snapshot.Kind != sandbox.SnapshotKindDefault {
				return sandbox.NewError(sandbox.ErrorInfo{
					Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create",
					Field: sandbox.CreateFieldSourceSnapshotKind,
					Message: "sandbox-kit daytona: explicit snapshot-kind restore mapping is unavailable",
				}, nil)
			}
			if snapshot.ProviderKind != nil {
				return sandbox.NewError(sandbox.ErrorInfo{
					Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create",
					Field: sandbox.CreateFieldSourceSnapshotProviderKind,
					Message: "sandbox-kit daytona: provider snapshot kind is unsupported",
				}, nil)
			}
			params.Snapshot = snapshot.Reference
		}
		plan.params = params
	}
	rest.Source = nil
	return nil
}

// durationMinutesAllowZero converts a delay to whole minutes; zero is valid
// (Daytona uses it for immediate delete or disabled stop/pause).
func durationMinutesAllowZero(value time.Duration) (int, error) {
	return durationMinutes(value, true)
}

// durationMinutesRequirePositive converts a delay to whole minutes; zero is rejected.
func durationMinutesRequirePositive(value time.Duration) (int, error) {
	return durationMinutes(value, false)
}

func durationMinutes(value time.Duration, allowZero bool) (int, error) {
	if value < 0 || value%time.Minute != 0 || value/time.Minute > math.MaxInt32 || (!allowZero && value == 0) {
		return 0, sandbox.NewError(sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "daytona",
			Operation: "create",
			Field:     sandbox.CreateFieldLifetime,
			Message:   "sandbox-kit daytona: policy/queue durations require representable whole minutes",
		}, nil)
	}
	return int(value / time.Minute), nil
}

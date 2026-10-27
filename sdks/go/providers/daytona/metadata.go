package daytona

import (
	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type observedSandbox struct {
	native *daytona.Sandbox
	info   *sandbox.SandboxInfo
}

func sandboxInfo(native *daytona.Sandbox) (*sandbox.SandboxInfo, error) {
	info, err := mapSandboxInfo(native)
	if err != nil {
		return info, err
	}
	observed := &observedSandbox{native: native, info: info}
	if err := observed.resources(); err != nil {
		return info, err
	}
	observed.state()
	observed.isolation()
	observed.metadata()
	return info, nil
}

func (s *observedSandbox) resources() error {
	resources, err := mapAllocatedResources(s.native)
	if err != nil {
		return err
	}
	if resources.CpuCores != nil {
		s.info.Resources = resources
		s.info.Origins["resources"] = sandbox.ValueOriginProvider
	}
	return nil
}

func (s *observedSandbox) isolation() {
	if s.native.SandboxClass == nil {
		return
	}
	isolation := &sandbox.IsolationConfig{}
	switch *s.native.SandboxClass {
	case types.SandboxClassContainer:
		isolation.Kind = sandbox.IsolationKindContainer
	case types.SandboxClassLinuxVM:
		isolation.Kind = sandbox.IsolationKindLinuxVM
	default:
		isolation.Kind = sandbox.IsolationKindProviderSpecific
		isolation.ProviderKind = sandbox.Value(string(*s.native.SandboxClass))
	}
	s.info.Isolation = isolation
	s.info.Origins["isolation"] = sandbox.ValueOriginProvider
}

func (s *observedSandbox) metadata() {
	native := s.native
	values := make(map[string]*sandbox.MetadataValue, 18)
	if native.Gpu > 0 {
		values["gpu_count"] = metadataInt(int64(native.Gpu))
	}
	putText(values, "gpu_type", native.GpuType)
	putText(values, "desired_state", native.DesiredState)
	putText(values, "error_reason", native.ErrorReason)
	if native.Recoverable != nil {
		values["recoverable"] = metadataBool(*native.Recoverable)
	}
	if native.Spot {
		values["spot"] = metadataBool(true)
	}
	putText(values, "spot_evicted_at", native.SpotEvictedAt)
	putText(values, "backup_state", native.BackupState)
	putText(values, "backup_created_at", native.BackupCreatedAt)
	putText(values, "updated_at", native.UpdatedAt)
	putText(values, "last_activity_at", native.LastActivityAt)
	putText(values, "auto_destroy_at", native.AutoDestroyAt)
	putText(values, "queue_timed_out_at", native.QueueTimedOutAt)
	putText(values, "snapshot", native.Snapshot)
	putText(values, "warm_pool_id", native.WarmPoolId)
	putInt(values, "auto_stop_interval_minutes", native.AutoStopInterval)
	putInt(values, "auto_pause_interval_minutes", native.AutoPauseInterval)
	putInt(values, "auto_archive_interval_minutes", native.AutoArchiveInterval)
	putInt(values, "auto_delete_interval_minutes", native.AutoDeleteInterval)
	s.info.ProviderMetadata = &sandbox.MetadataObject{Values: values}
	s.info.Origins[sandbox.InfoFieldProviderMetadata] = sandbox.ValueOriginProvider
}

func putText[T ~string](values map[string]*sandbox.MetadataValue, key string, value *T) {
	if value == nil {
		return
	}
	values[key] = metadataText(string(*value))
}

func putInt(values map[string]*sandbox.MetadataValue, key string, value int) {
	values[key] = metadataInt(int64(value))
}

func metadataInt(value int64) *sandbox.MetadataValue {
	return &sandbox.MetadataValue{Value: &sandbox.MetadataSignedInteger{SignedInteger: value}}
}

func metadataText(value string) *sandbox.MetadataValue {
	return &sandbox.MetadataValue{Value: &sandbox.MetadataText{Text: value}}
}

func metadataBool(value bool) *sandbox.MetadataValue {
	return &sandbox.MetadataValue{Value: &sandbox.MetadataBoolean{Boolean: value}}
}

func (s *observedSandbox) state() {
	switch string(s.native.State) {
	case "creating", "restoring", "pending_build", "building_snapshot", "pulling_snapshot", "forking":
		s.info.State = sandbox.SandboxStateQueued
	case "starting", "resuming":
		s.info.State = sandbox.SandboxStateStarting
	case "started":
		s.info.State = sandbox.SandboxStateRunning
	case "stopping":
		s.info.State = sandbox.SandboxStateStopping
	case "stopped":
		s.info.State = sandbox.SandboxStateStopped
	case "pausing":
		s.info.State = sandbox.SandboxStatePausing
	case "paused":
		s.info.State = sandbox.SandboxStatePaused
	case "archiving":
		s.info.State = sandbox.SandboxStateArchiving
	case "archived":
		s.info.State = sandbox.SandboxStateArchived
	case "destroying":
		s.info.State = sandbox.SandboxStateDeleting
	case "destroyed":
		s.info.State = sandbox.SandboxStateDeleted
		s.info.Terminal = sandbox.Value(true)
	case "error", "build_failed":
		s.info.State = sandbox.SandboxStateError
	}
}

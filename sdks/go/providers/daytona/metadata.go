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
	values := make(map[string]*sandbox.MetadataValue, 18)
	if s.native.Gpu > 0 {
		values["gpu_count"] = metadataInt(int64(s.native.Gpu))
	}
	if s.native.GpuType != nil {
		values["gpu_type"] = metadataText(string(*s.native.GpuType))
	}
	if s.native.DesiredState != nil {
		values["desired_state"] = metadataText(string(*s.native.DesiredState))
	}
	if s.native.ErrorReason != nil {
		values["error_reason"] = metadataText(*s.native.ErrorReason)
	}
	if s.native.Recoverable != nil {
		values["recoverable"] = metadataBool(*s.native.Recoverable)
	}
	if s.native.Spot {
		values["spot"] = metadataBool(true)
	}
	if s.native.SpotEvictedAt != nil {
		values["spot_evicted_at"] = metadataText(*s.native.SpotEvictedAt)
	}
	if s.native.BackupState != nil {
		values["backup_state"] = metadataText(*s.native.BackupState)
	}
	if s.native.BackupCreatedAt != nil {
		values["backup_created_at"] = metadataText(*s.native.BackupCreatedAt)
	}
	if s.native.UpdatedAt != nil {
		values["updated_at"] = metadataText(*s.native.UpdatedAt)
	}
	if s.native.LastActivityAt != nil {
		values["last_activity_at"] = metadataText(*s.native.LastActivityAt)
	}
	if s.native.AutoDestroyAt != nil {
		values["auto_destroy_at"] = metadataText(*s.native.AutoDestroyAt)
	}
	if s.native.QueueTimedOutAt != nil {
		values["queue_timed_out_at"] = metadataText(*s.native.QueueTimedOutAt)
	}
	if s.native.Snapshot != nil {
		values["snapshot"] = metadataText(*s.native.Snapshot)
	}
	if s.native.WarmPoolId != nil {
		values["warm_pool_id"] = metadataText(*s.native.WarmPoolId)
	}
	values["auto_stop_interval_minutes"] = metadataInt(int64(s.native.AutoStopInterval))
	values["auto_pause_interval_minutes"] = metadataInt(int64(s.native.AutoPauseInterval))
	values["auto_archive_interval_minutes"] = metadataInt(int64(s.native.AutoArchiveInterval))
	values["auto_delete_interval_minutes"] = metadataInt(int64(s.native.AutoDeleteInterval))
	s.info.ProviderMetadata = &sandbox.MetadataObject{Values: values}
	s.info.Origins[sandbox.InfoFieldProviderMetadata] = sandbox.ValueOriginProvider
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

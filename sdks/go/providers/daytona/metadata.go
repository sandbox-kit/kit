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
	if err := observed.metadata(); err != nil {
		return info, err
	}
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

func (s *observedSandbox) metadata() error {
	metadata := map[string]any{}
	if s.native.Gpu > 0 {
		metadata["gpu_count"] = int64(s.native.Gpu)
	}
	if s.native.GpuType != nil {
		metadata["gpu_type"] = string(*s.native.GpuType)
	}
	if s.native.DesiredState != nil {
		metadata["desired_state"] = string(*s.native.DesiredState)
	}
	if s.native.ErrorReason != nil {
		metadata["error_reason"] = *s.native.ErrorReason
	}
	if s.native.Recoverable != nil {
		metadata["recoverable"] = *s.native.Recoverable
	}
	if s.native.Spot {
		metadata["spot"] = true
	}
	if s.native.SpotEvictedAt != nil {
		metadata["spot_evicted_at"] = *s.native.SpotEvictedAt
	}
	if s.native.BackupState != nil {
		metadata["backup_state"] = *s.native.BackupState
	}
	if s.native.BackupCreatedAt != nil {
		metadata["backup_created_at"] = *s.native.BackupCreatedAt
	}
	if s.native.UpdatedAt != nil {
		metadata["updated_at"] = *s.native.UpdatedAt
	}
	if s.native.LastActivityAt != nil {
		metadata["last_activity_at"] = *s.native.LastActivityAt
	}
	if s.native.AutoDestroyAt != nil {
		metadata["auto_destroy_at"] = *s.native.AutoDestroyAt
	}
	if s.native.QueueTimedOutAt != nil {
		metadata["queue_timed_out_at"] = *s.native.QueueTimedOutAt
	}
	if s.native.Snapshot != nil {
		metadata["snapshot"] = *s.native.Snapshot
	}
	if s.native.WarmPoolId != nil {
		metadata["warm_pool_id"] = *s.native.WarmPoolId
	}
	metadata["auto_stop_interval_minutes"] = int64(s.native.AutoStopInterval)
	metadata["auto_pause_interval_minutes"] = int64(s.native.AutoPauseInterval)
	metadata["auto_archive_interval_minutes"] = int64(s.native.AutoArchiveInterval)
	metadata["auto_delete_interval_minutes"] = int64(s.native.AutoDeleteInterval)
	if len(metadata) == 0 {
		return nil
	}
	object, err := sandbox.MetadataObjectFromMap(metadata)
	if err != nil {
		return err
	}
	s.info.ProviderMetadata = object
	s.info.Origins[sandbox.InfoFieldProviderMetadata] = sandbox.ValueOriginProvider
	return nil
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

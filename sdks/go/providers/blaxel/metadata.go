package blaxel

import (
	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type observedSandbox struct {
	native *blaxel.Sandbox
	info   *sandbox.SandboxInfo
}

func sandboxInfo(native *blaxel.Sandbox) (*sandbox.SandboxInfo, error) {
	info, err := mapSandboxInfo(native)
	if err != nil {
		return info, err
	}
	observed := &observedSandbox{native: native, info: info}
	if err := observed.resources(); err != nil {
		return info, err
	}
	observed.state()
	observed.metadata()
	return info, nil
}

func (s *observedSandbox) resources() error {
	resources, err := mapAllocatedResources(s.native)
	if err != nil {
		return err
	}
	if resources.MemoryMib != nil {
		s.info.Resources = resources
	}
	return nil
}

func (s *observedSandbox) state() {
	switch s.native.Status {
	case blaxel.StatusUploading, blaxel.StatusBuilding, blaxel.StatusDeploying, blaxel.StatusBuilt:
		s.info.State = sandbox.SandboxStateQueued
	case blaxel.StatusDeployed:
		s.info.State = runtimeState(s.native.State, sandbox.SandboxStateRunning)
	case blaxel.StatusDeleting, blaxel.StatusDeactivating:
		s.info.State = sandbox.SandboxStateDeleting
	case blaxel.StatusTerminated, blaxel.StatusDeactivated:
		s.info.State = sandbox.SandboxStateDeleted
		s.info.Terminal = sandbox.Value(true)
	case blaxel.StatusArchiving:
		s.info.State = sandbox.SandboxStateArchiving
	case blaxel.StatusArchived:
		s.info.State = sandbox.SandboxStateArchived
	case blaxel.StatusUnarchiving:
		s.info.State = sandbox.SandboxStateStarting
	case blaxel.StatusFailed:
		s.info.State = sandbox.SandboxStateError
	default:
		if mapped, ok := runtimeOnlyState(s.native.State); ok {
			s.info.State = mapped
		}
	}
}

func runtimeState(runtime blaxel.SandboxState, fallback sandbox.SandboxState) sandbox.SandboxState {
	if runtime == blaxel.SandboxStateStandby {
		return sandbox.SandboxStatePaused
	}
	return fallback
}

func runtimeOnlyState(runtime blaxel.SandboxState) (sandbox.SandboxState, bool) {
	switch runtime {
	case blaxel.SandboxStateRunning:
		return sandbox.SandboxStateRunning, true
	case blaxel.SandboxStateStandby:
		return sandbox.SandboxStatePaused, true
	default:
		return 0, false
	}
}

func (s *observedSandbox) metadata() {
	native := s.native
	values := make(map[string]*sandbox.MetadataValue, 12)
	putText(values, "display_name", native.Metadata.DisplayName)
	putText(values, "created_at", native.Metadata.CreatedAt)
	putText(values, "updated_at", native.Metadata.UpdatedAt)
	putText(values, "last_used_at", native.LastUsedAt)
	putText(values, "node_generation", native.NodeGeneration)
	if native.ExpiresIn > 0 {
		values["expires_in_seconds"] = metadataInt(native.ExpiresIn)
	}
	if native.State != "" {
		values["runtime_state"] = metadataText(string(native.State))
	}
	if len(values) > 0 {
		s.info.ProviderMetadata = &sandbox.MetadataObject{Values: values}
	}
	if native.Metadata.CreatedAt != "" {
		if parsed, err := parseBlaxelTime(native.Metadata.CreatedAt); err == nil {
			s.info.CreatedAt = timestamppb.New(parsed)
		}
	}
}

func putText(values map[string]*sandbox.MetadataValue, key, value string) {
	if value == "" {
		return
	}
	values[key] = metadataText(value)
}

func metadataInt(value int64) *sandbox.MetadataValue {
	return &sandbox.MetadataValue{Value: &sandbox.MetadataSignedInteger{SignedInteger: value}}
}

func metadataText(value string) *sandbox.MetadataValue {
	return &sandbox.MetadataValue{Value: &sandbox.MetadataText{Text: value}}
}

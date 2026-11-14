package blaxel

import (
	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
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
	case blaxel.StatusTerminated, blaxel.StatusDeactivated, blaxel.Status("DELETED"):
		// DELETED is not a typed Status constant in sdk-go v0.31.0; accept the
		// wire string so refresh matches TypeScript mapSandboxState.
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

package sandbox

import (
	"fmt"
	"sync"
	"time"
)

// Sandbox is the handle for one cloud sandbox.
// ID and Info describe it. Lifecycle methods change it. Close releases this
// handle and leaves the cloud sandbox running.
//
// Example:
//
//	instance, err := client.Create(ctx, nil)
//	fmt.Println(instance.ID(), instance.ProviderName())
//	info := instance.Info()
type Sandbox struct {
	mu            sync.RWMutex
	operations    Backend
	timeout       *time.Duration
	operationGate chan struct{}
	owner         *Client
	closed        bool
	// info is an owned metadata copy. Info returns another copy so callers cannot mutate this one.
	info *SandboxInfo
}

// sandboxFromResponse checks that creation returned an identifier
// and consistent provider metadata.
func sandboxFromResponse(backend Backend, response *CreateResult, operation string) (*Sandbox, error) {
	if response == nil || response.Sandbox == nil || response.Sandbox.Id == "" {
		return nil, NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  backend.Name(),
			Operation: operation,
			Field:     InfoFieldId,
			Message:   "sandbox-kit: provider returned no sandbox identity",
		}, nil)
	}
	info := cloneSandboxInfo(response.Sandbox)
	name := backend.Name()
	if info.Provider != "" && info.Provider != name {
		return nil, NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  name,
			Operation: operation,
			Field:     InfoFieldProvider,
			Message:   fmt.Sprintf("sandbox-kit: response provider %q differs from selected provider %q", info.Provider, name),
		}, nil)
	}
	info.Provider = name
	if err := validateResponseOrigins(info, operation); err != nil {
		return nil, err
	}
	return &Sandbox{operations: backend, info: info, operationGate: make(chan struct{}, 1)}, nil
}

// ID returns the provider sandbox identifier.
// A nil handle returns an empty string.
func (s *Sandbox) ID() string {
	if s == nil {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.info == nil {
		return ""
	}
	return s.info.Id
}

// ProviderName returns "daytona" or "modal".
// A nil handle returns an empty string.
func (s *Sandbox) ProviderName() string {
	if s == nil {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.info == nil {
		return ""
	}
	return s.info.Provider
}

// Info returns a copy of the sandbox metadata.
// Changing the result does not change the handle. A nil handle returns nil.
// State is a SandboxState. Origins maps a field path such as InfoFieldId to ValueOriginProvider or ValueOriginRequest.
//
// Example:
//
//	info := instance.Info()
//	fmt.Println(info.GetID(), info.Origins[InfoFieldId])
func (s *Sandbox) Info() *SandboxInfo {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.info == nil {
		return nil
	}
	info := cloneSandboxInfo(s.info)
	return info
}

// Capabilities lists the lifecycle operations this provider implements.
// Pass a value to Supports, or compare it with LifecycleCapabilityPause and the other capability constants.
// A closed client, a closed handle, or a terminal sandbox returns nil.
func (s *Sandbox) Capabilities() []LifecycleCapability {
	if s == nil || s.operations == nil {
		return nil
	}
	if s.owner != nil && s.owner.isClosed() {
		return nil
	}
	s.mu.RLock()
	blocked := s.closed || s.info == nil || s.info.GetTerminal()
	s.mu.RUnlock()
	if blocked {
		return nil
	}
	return staticLifecycleCapabilities(s.operations)
}

// operationView is a locked snapshot of the handle fields a lifecycle call reads.
// id is empty when metadata is missing. The snapshot does not clone that metadata.
type operationView struct {
	id       string
	closed   bool
	terminal bool
}

func (s *Sandbox) operationView() operationView {
	if s == nil {
		return operationView{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	view := operationView{closed: s.closed}
	if s.info != nil {
		view.id = s.info.Id
		view.terminal = s.info.GetTerminal()
	}
	return view
}

// Supports reports whether this provider implements capability.
// A terminal or closed handle supports nothing.
func (s *Sandbox) Supports(capability LifecycleCapability) bool {
	for _, candidate := range s.Capabilities() {
		if candidate == capability {
			return true
		}
	}
	return false
}

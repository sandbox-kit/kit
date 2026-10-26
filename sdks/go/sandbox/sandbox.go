package sandbox

import (
	"fmt"
	"sync"
	"time"
)

// Sandbox is the handle returned by Client.Create.
// It exposes identity and a copy of metadata. It does not expose the native SDK object.
// Lifecycle methods dispatch through private provider capabilities without exposing native SDK objects.
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

// ProviderName returns the selected provider, such as "daytona" or "modal".
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

// Info returns an owned metadata copy. Changing the result does not change the handle.
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

type sandboxCapabilityResolver interface {
	SandboxCapabilities(*SandboxInfo) []LifecycleCapability
}

// Capabilities returns the operations supported by this provider and resolved
// sandbox class. A terminal handle has no further lifecycle capabilities.
func (s *Sandbox) Capabilities() []LifecycleCapability {
	if s == nil || s.operations == nil {
		return nil
	}
	s.mu.RLock()
	closed := s.closed
	s.mu.RUnlock()
	if closed {
		return nil
	}
	if s.owner != nil && s.owner.isClosed() {
		return nil
	}
	info := s.Info()
	if info == nil || info.GetTerminal() {
		return nil
	}
	if resolver, ok := s.operations.(sandboxCapabilityResolver); ok {
		values := resolver.SandboxCapabilities(info)
		return append([]LifecycleCapability(nil), values...)
	}
	return staticLifecycleCapabilities(s.operations)
}

// Supports reports whether capability is available for this resolved handle.
func (s *Sandbox) Supports(capability LifecycleCapability) bool {
	for _, candidate := range s.Capabilities() {
		if candidate == capability {
			return true
		}
	}
	return false
}

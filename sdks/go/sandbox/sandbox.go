package sandbox

import (
	"fmt"
)

// Sandbox is the handle returned by Client.Create.
// It exposes identity and a copy of metadata. It does not expose the native SDK object.
// Lifecycle methods such as stop and delete are not implemented on this handle yet.
//
// Example:
//
//	instance, err := client.Create(ctx, nil)
//	fmt.Println(instance.ID(), instance.ProviderName())
//	info := instance.Info()
type Sandbox struct {
	// info is an owned metadata copy. Info returns another copy so callers cannot mutate this one.
	info *SandboxInfo
}

// sandboxFromResponse checks that creation returned an identifier
// and consistent provider metadata.
func sandboxFromResponse(backend Backend, response *CreateResult) (*Sandbox, error) {
	if response == nil || response.Sandbox == nil || response.Sandbox.ID == "" {
		return nil, NewError(ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  backend.Name(),
			Operation: "create",
			Field:     InfoFieldID,
			Message:   "sandbox-kit: creation binding returned no sandbox identity",
		}, nil)
	}
	info, err := response.Sandbox.Clone()
	if err != nil {
		return nil, NewError(ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  backend.Name(),
			Operation: "create",
			Field:     InfoFieldProviderMetadata,
			Message:   err.Error(),
		}, err)
	}
	name := backend.Name()
	if info.Provider != "" && info.Provider != name {
		return nil, NewError(ErrorInfo{
			Kind:      ErrorKindInvalidResponse,
			Provider:  name,
			Operation: "create",
			Field:     InfoFieldProvider,
			Message:   fmt.Sprintf("sandbox-kit: response provider %q differs from selected provider %q", info.Provider, name),
		}, nil)
	}
	info.Provider = name
	if err := validateResponseOrigins(info); err != nil {
		return nil, err
	}
	return &Sandbox{info: info}, nil
}

// ID returns the provider sandbox identifier.
func (s *Sandbox) ID() string { return s.info.ID }

// ProviderName returns the selected provider, such as "daytona" or "modal".
func (s *Sandbox) ProviderName() string { return s.info.Provider }

// Info returns an owned metadata copy. Changing the result does not change the handle.
//
// Example:
//
//	info := instance.Info()
//	fmt.Println(info.GetID(), info.Origins[InfoFieldID])
func (s *Sandbox) Info() *SandboxInfo {
	info, _ := s.info.Clone()
	return info
}

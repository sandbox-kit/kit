package sandbox

import (
	"fmt"
)

// Sandbox is the common handle returned by Client.Create. Lifecycle capabilities
// will be added to this handle separately; no native SDK object is exposed.
type Sandbox struct {
	info *SandboxInfo
}

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

func (s *Sandbox) ID() string           { return s.info.ID }
func (s *Sandbox) ProviderName() string { return s.info.Provider }

// Info returns a copy, preserving the handle's identity against caller mutation.
func (s *Sandbox) Info() *SandboxInfo {
	info, _ := s.info.Clone()
	return info
}

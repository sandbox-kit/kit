// Client runtime behavior is handwritten. Portable request and response data
// remains generated from protobuf.
// Client owns an initialized provider SDK. Create returns a sandbox handle.
// Close releases the SDK and leaves the cloud sandbox running.
//
// Example:
//
//	client, err := NewClient(provider, &Config{})
//	instance, err := client.Create(ctx, nil)
//	err = client.Close(ctx)
package sandbox

import (
	context "context"
	errors "errors"
	fmt "fmt"
	netip "net/netip"
	url "net/url"
	reflect "reflect"
	strings "strings"
	time "time"

	"buf.build/go/protovalidate"
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"
)

// Backend is the provider identity and local SDK cleanup contract. Operations
// are optional capabilities.
type Backend interface {
	// Name returns the provider identifier, such as "daytona" or "modal".
	Name() string
	// Close releases the native SDK. It does not stop or delete the cloud sandbox.
	//
	// Example:
	//
	//   err := backend.Close(ctx)
	//
	Close(context.Context) error
}

// Provider initializes the provider SDK from shared configuration.
type Provider interface {
	// Name returns the provider identifier selected by this factory.
	Name() string
	// NewClient builds a backend from shared configuration. Call it through
	// sandbox.NewClient,
	// which validates the configuration and owns cleanup.
	//
	// Example:
	//
	// 	backend, err := provider.NewClient(&config)
	//
	NewClient(*sandboxv1.Config) (Backend, error)
}

// Client owns its initialized provider SDK. Close it when finished.
type Client struct {
	// backend is the initialized provider SDK. Applications use Client methods
	// instead of this field.
	backend Backend
	// timeout is a copy of Config.Timeout. Nil adds no Kit deadline. An explicit
	// zero also adds none.
	timeout *time.Duration
	// lifecycle coordinates operations with idempotent client closure.
	lifecycle *clientLifecycle
}

// SandboxCreator is the optional create capability on a provider Backend.
type SandboxCreator interface {
	Create(context.Context, *CreateOptions) (*CreateResult, error)
}

// SandboxAttacher is the optional attach capability on a provider Backend.
type SandboxAttacher interface {
	AttachSandbox(context.Context, string) (*SandboxInfo, error)
}

// SandboxRefresher is the optional refresh capability on a provider Backend.
type SandboxRefresher interface {
	RefreshSandbox(context.Context, string) (*SandboxInfo, error)
}

// SandboxStarter is the optional start capability on a provider Backend.
type SandboxStarter interface {
	StartSandbox(context.Context, string, *StartOptions) (*SandboxInfo, error)
}

// SandboxStopper is the optional stop capability on a provider Backend.
type SandboxStopper interface {
	StopSandbox(context.Context, string, *StopOptions) (*SandboxInfo, error)
}

// SandboxPauser is the optional pause capability on a provider Backend.
type SandboxPauser interface {
	PauseSandbox(context.Context, string, *PauseOptions) (*SandboxInfo, error)
}

// SandboxTerminator is the optional terminate capability on a provider Backend.
type SandboxTerminator interface {
	TerminateSandbox(context.Context, string, *TerminateOptions) (*ExitResult, error)
}

// SandboxArchiver is the optional archive capability on a provider Backend.
type SandboxArchiver interface {
	ArchiveSandbox(context.Context, string, *ArchiveOptions) (*SandboxInfo, error)
}

// SandboxDeleter is the optional delete capability on a provider Backend.
type SandboxDeleter interface {
	DeleteSandbox(context.Context, string, *DeleteOptions) (*SandboxInfo, error)
}

// missingClientBinding reports whether a provider or backend value is missing.
// Nil interfaces, pointers, slices, and maps count as missing.
func missingClientBinding(binding any) bool {
	value := reflect.ValueOf(binding)
	if !value.IsValid() {
		return true
	}
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	}
	return false
}

// NewClient validates shared configuration and initializes the selected
// provider SDK.
// The returned client owns cleanup. Closing it does not stop the sandbox.
//
// Example:
//
//	client, err := NewClient(provider, &Config{
//	    Timeout: Duration(2 * time.Minute),
//	})
func NewClient(providerFactory Provider, config *sandboxv1.Config) (client *Client, result error) {
	provider := ""
	defer func() { result = WithErrorContext(result, provider, "initialize") }()
	if missingClientBinding(providerFactory) {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: "initialize", Field: "provider", Message: "sandbox-kit: provider is required"}, nil)
	}
	provider = providerFactory.Name()
	if err := validateProtoConfig(config); err != nil {
		return nil, err
	}
	expected := provider
	if strings.TrimSpace(expected) == "" {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: "initialize", Field: "provider", Message: "sandbox-kit: provider name is required"}, nil)
	}
	var timeout *time.Duration
	if config != nil && config.Timeout != nil {
		value := config.Timeout.AsDuration()
		timeout = &value
	}
	backend, err := providerFactory.NewClient(config)
	if err != nil {
		return nil, err
	}
	if missingClientBinding(backend) {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: expected, Operation: "initialize", Field: "backend", Message: "sandbox-kit: provider returned no initialized backend"}, nil)
	}
	if actual := backend.Name(); actual != expected {
		return nil, errors.Join(NewError(ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: expected, Operation: "initialize", Field: "provider", Message: fmt.Sprintf("sandbox-kit: initialized backend %q differs from selected provider %q", actual, expected)}, nil), backend.Close(context.Background()))
	}
	return &Client{backend: backend, timeout: timeout, lifecycle: &clientLifecycle{}}, nil
}

// validateProtoConfig validates the generated client contract before native
// provider initialization. Portable rules come from Protovalidate options on
// kit.sandbox.v1.Config; secure_endpoint stays handwritten. Provider-specific
// support remains provider-owned.
func validateProtoConfig(config *sandboxv1.Config) error {
	if config == nil {
		return nil
	}
	if err := protovalidate.Validate(config); err != nil {
		return WithErrorContext(mapProtovalidateError(err), "", "initialize")
	}
	if config.Endpoint != nil && !secureEndpoint(*config.Endpoint) {
		return NewError(ErrorInfo{
			Kind:      ErrorKindInvalidArgument,
			Operation: "initialize",
			Field:     "endpoint",
			Rule:      Value("secure_endpoint"),
			Message:   "sandbox-kit: endpoint must use HTTPS except for loopback HTTP development endpoints",
		}, nil)
	}
	return nil
}

func secureEndpoint(value string) bool {
	u, err := url.Parse(value)
	if err != nil || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
		return false
	}
	if u.Scheme == "https" {
		return true
	}
	if u.Scheme != "http" {
		return false
	}
	ip, err := netip.ParseAddr(u.Hostname())
	return u.Hostname() == "localhost" || (err == nil && ip.IsLoopback())
}

// ProviderName returns the selected provider identifier, such as "daytona" or
// "modal".
func (c *Client) ProviderName() string {
	if c == nil || c.backend == nil {
		return ""
	}
	return c.backend.Name()
}

// Close releases the native SDK. It does not stop or delete the cloud sandbox.
// A nil client is a no-op. A nil context returns an invalid-argument error.
//
// Example:
//
//	err := client.Close(ctx)
func (c *Client) Close(ctx context.Context) (result error) {
	provider := ""
	if c != nil && c.backend != nil {
		provider = c.backend.Name()
	}
	defer func() { result = WithErrorContext(result, provider, "close") }()
	if ctx == nil {
		return NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "close", Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	return c.closeBackend(ctx)
}

// Create creates one sandbox and returns a handle. A nil request uses provider
// defaults.
// Close the client when finished; that does not delete the sandbox.
//
// Example:
//
//	instance, err := client.Create(ctx, nil)
//	fmt.Println(instance.ID(), instance.ProviderName())
func (c *Client) Create(ctx context.Context, request *CreateOptions) (instance *Sandbox, result error) {
	provider := ""
	if c != nil && c.backend != nil {
		provider = c.backend.Name()
	}
	defer func() { result = WithErrorContext(result, provider, "create") }()
	if ctx == nil {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "create", Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c == nil || c.backend == nil {
		return nil, NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "create", Field: "client", Message: "sandbox-kit: client is not initialized"}, nil)
	}
	release, err := c.beginOperation("create")
	if err != nil {
		return nil, err
	}
	defer release()
	request, err = prepareCreateRequest(request)
	if err != nil {
		return nil, err
	}
	timeout := OptionalDuration(request.GetProvisioning().GetTimeout())
	if timeout == nil {
		timeout = c.timeout
	}
	if timeout != nil && *timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *timeout)
		defer cancel()
	}
	if timeout != nil {
		if request.Provisioning == nil {
			request.Provisioning = &ProvisioningOptions{}
		}
		request.Provisioning.Timeout = Duration(*timeout)
	}
	creator, ok := c.backend.(SandboxCreator)
	if !ok {
		return nil, NewError(ErrorInfo{Kind: ErrorKindUnsupported, Provider: provider, Operation: "create", Field: "", Message: "sandbox-kit: provider does not support sandbox creation"}, nil)
	}
	response, err := creator.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	instance, err = sandboxFromResponse(c.backend, response, "create")
	if err != nil {
		return nil, err
	}
	instance.timeout = c.timeout
	instance.owner = c
	return instance, nil
}

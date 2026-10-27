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
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
)

// Backend is the open provider connection held by Client.
// Close releases that connection. It does not stop or delete a cloud sandbox.
// CloseSandbox runs when the caller closes a sandbox handle. It does not stop or delete the cloud sandbox.
type Backend interface {
	// Name returns the provider name, "daytona" or "modal".
	Name() string
	// Close releases the provider connection. It does not stop or delete a cloud sandbox.
	Close(context.Context) error
	// CloseSandbox runs when the caller closes a sandbox handle.
	// It does not stop or delete the cloud sandbox.
	CloseSandbox(context.Context, string) error
}

// Provider is the value passed to NewClient.
// Use daytona.New or modal.New.
type Provider interface {
	// Name returns the provider name, "daytona" or "modal".
	Name() string
	// NewClient opens the provider from config.
	// Call the package NewClient. That checks config and owns cleanup.
	NewClient(*Config) (Backend, error)
}

// Client creates and looks up sandboxes.
// Close releases the local connection. It does not stop or delete a cloud sandbox.
type Client struct {
	// backend is the initialized provider SDK. Applications use Client methods
	// instead of this field.
	backend Backend
	// timeout is a copy of Config.Timeout. Nil adds no Kit deadline. An explicit
	// zero also adds none.
	timeout *time.Duration
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

// NewClient checks config and opens a client for the selected provider.
// Closing the client releases the provider SDK. It does not stop or delete a sandbox.
//
// Pass daytona.New or modal.New as providerFactory.
//
// Example:
//
//	client, err := NewClient(provider, &Config{
//	    Timeout: Duration(2 * time.Minute),
//	})
func NewClient(providerFactory Provider, config *Config) (client *Client, result error) {
	provider := ""
	defer func() { result = sdk.WithErrorContext(result, provider, "initialize") }()
	if missingClientBinding(providerFactory) {
		return nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: "initialize", Field: "provider", Message: "sandbox-kit: provider is required"}, nil)
	}
	provider = providerFactory.Name()
	if err := validateProtoConfig(config); err != nil {
		return nil, err
	}
	expected := provider
	if strings.TrimSpace(expected) == "" {
		return nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: "initialize", Field: "provider", Message: "sandbox-kit: provider name is required"}, nil)
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
		return nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: expected, Operation: "initialize", Field: "backend", Message: "sandbox-kit: provider returned no initialized backend"}, nil)
	}
	if actual := backend.Name(); actual != expected {
		return nil, errors.Join(NewError(&ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: expected, Operation: "initialize", Field: "provider", Message: fmt.Sprintf("sandbox-kit: initialized backend %q differs from selected provider %q", actual, expected)}, nil), backend.Close(context.Background()))
	}
	return &Client{backend: backend, timeout: timeout}, nil
}

// validateProtoConfig validates the generated client contract before native
// provider initialization. Portable rules come from Protovalidate options on
// kit.sandbox.v1.Config; secure_endpoint stays handwritten. Provider-specific
// support remains provider-owned.
func validateProtoConfig(config *Config) error {
	if config == nil {
		return nil
	}
	if err := protovalidate.Validate(config); err != nil {
		return sdk.WithErrorContext(mapProtovalidateError(err), "", "initialize")
	}
	if config.Endpoint != nil && !secureEndpoint(*config.Endpoint) {
		return NewError(&ErrorInfo{
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

// ProviderName returns "daytona" or "modal".
// A nil client returns an empty string.
func (c *Client) ProviderName() string {
	if c == nil || c.backend == nil {
		return ""
	}
	return c.backend.Name()
}

// Close releases the local provider connection. It does not stop or delete the cloud sandbox.
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
	defer func() { result = sdk.WithErrorContext(result, provider, "close") }()
	if ctx == nil {
		return NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "close", Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if c == nil || c.backend == nil {
		return nil
	}
	return c.backend.Close(ctx)
}

// Create creates one sandbox and returns its handle.
// A nil request uses the provider default. Modal requires an image in CreateOptions.Source.
// A nil client or a nil context returns ErrorKindInvalidArgument. Read a failure with errors.As and *Error.
// Close the client when finished. That does not delete the sandbox.
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
	defer func() { result = sdk.WithErrorContext(result, provider, "create") }()
	if ctx == nil {
		return nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "create", Field: "context", Message: "sandbox-kit: context is required"}, nil)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if c == nil || c.backend == nil {
		return nil, NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Provider: provider, Operation: "create", Field: "client", Message: "sandbox-kit: client is not initialized"}, nil)
	}
	request, err := prepareCreateRequest(request)
	if err != nil {
		return nil, err
	}
	timeout := sdk.OptionalDuration(request.GetProvisioning().GetTimeout())
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
		return nil, NewError(&ErrorInfo{Kind: ErrorKindUnsupported, Provider: provider, Operation: "create", Field: "", Message: "sandbox-kit: provider does not support sandbox creation"}, nil)
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
	return instance, nil
}

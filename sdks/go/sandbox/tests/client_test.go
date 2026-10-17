package sandbox_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/proto"
)

type testProvider struct {
	provider sandbox.Backend
	err      error
	calls    int
	config   *sandbox.Config
}

func (*testProvider) Name() string { return "test" }
func (f *testProvider) NewClient(config *sandbox.Config) (sandbox.Backend, error) {
	f.calls++
	f.config = config
	return f.provider, f.err
}
func (p *creationBackend) Close(context.Context) error { return nil }
func newTestClient(provider sandbox.Backend) (*sandbox.Client, error) {
	return sandbox.NewClient(&testProvider{provider: provider}, &sandbox.Config{})
}
func TestClientInitializesFromCommonConfiguration(t *testing.T) {
	factory := &testProvider{provider: &creationBackend{}}
	config := &sandbox.Config{
		Auth:     &sandbox.AuthConfig{Credentials: &sandbox.AuthConfig_ApiKey{ApiKey: &sandbox.APIKeyCredentials{Key: "test-key"}}},
		Endpoint: sandbox.Value("https://example.invalid"),
		Region:   sandbox.Value("region"),
		Timeout:  sandbox.Duration(time.Second),
	}
	client, err := sandbox.NewClient(factory, config)
	if err != nil {
		t.Fatal(err)
	}
	if factory.calls != 1 || factory.config.GetEndpoint() != config.GetEndpoint() || factory.config.GetRegion() != config.GetRegion() || client.ProviderName() != "test" {
		t.Fatal("initialization settings were lost")
	}
}
func TestClientRejectsInvalidConfigurationBeforeInitialization(t *testing.T) {
	for _, config := range []*sandbox.Config{
		{Timeout: sandbox.Duration(-time.Second)},
		{Endpoint: sandbox.Value("  ")},
		{Endpoint: sandbox.Value("ftp://example.invalid")},
		{Endpoint: sandbox.Value("not-a-url")},
		{Auth: &sandbox.AuthConfig{}},
		{Auth: &sandbox.AuthConfig{Credentials: &sandbox.AuthConfig_TokenPair{TokenPair: &sandbox.TokenPairCredentials{Id: "id"}}}},
		{Auth: &sandbox.AuthConfig{Credentials: &sandbox.AuthConfig_OauthRefresh{OauthRefresh: &sandbox.OAuthCredentials{RefreshToken: "refresh", ClientId: "id"}}}},
	} {
		factory := &testProvider{provider: &creationBackend{}}
		if client, err := sandbox.NewClient(factory, config); err == nil || client != nil {
			t.Fatal("invalid config accepted")
		}
		if factory.calls != 0 {
			t.Fatal("invalid config initialized SDK")
		}
	}
}
func TestClientRejectsNilFactoriesAndProviders(t *testing.T) {
	var typedNil *testProvider
	for _, factory := range []sandbox.Provider{nil, typedNil, &testProvider{}, &testProvider{provider: (*creationBackend)(nil)}} {
		if client, err := sandbox.NewClient(factory, &sandbox.Config{}); err == nil || client != nil {
			t.Fatal("nil factory/provider accepted")
		}
	}
}
func TestInitializationErrorsPassThrough(t *testing.T) {
	expected := errors.New("SDK failure")
	factory := &testProvider{err: expected}
	if _, err := sandbox.NewClient(factory, &sandbox.Config{}); !errors.Is(err, expected) {
		t.Fatal("SDK initialization error changed")
	}
}
func TestClientDefaultTimeoutAndPerOperationOverrides(t *testing.T) {
	for _, tc := range []struct {
		name                            string
		clientTimeout, operationTimeout *time.Duration
		hasDeadline                     bool
	}{
		{"client default", sandbox.Value(time.Second), nil, true},
		{"operation disables Kit timeout", sandbox.Value(time.Second), sandbox.Value(time.Duration(0)), false},
		{"operation timeout", nil, sandbox.Value(time.Second), true},
		{"no Kit timeout", nil, nil, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			provider := &creationBackend{call: func(ctx context.Context, r *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
				_, has := ctx.Deadline()
				if has != tc.hasDeadline {
					t.Fatalf("deadline=%v", has)
				}
				return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "id"}}, nil
			}}
			config := &sandbox.Config{Timeout: nil}
			if tc.clientTimeout != nil {
				config.Timeout = sandbox.Duration(*tc.clientTimeout)
			}
			client, err := sandbox.NewClient(&testProvider{provider: provider}, config)
			if err != nil {
				t.Fatal(err)
			}
			if config.Timeout != nil {
				config.Timeout.Seconds = 0
				config.Timeout.Nanos = 0
			}
			request := &sandbox.CreateOptions{}
			if tc.operationTimeout != nil {
				request.Provisioning = &sandbox.ProvisioningOptions{Timeout: tc.operationTimeout}
			}
			if _, err := client.Create(context.Background(), request); err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestConfigContractsAreProtobufMessages(t *testing.T) {
	config := &sandbox.Config{Auth: &sandbox.AuthConfig{Credentials: &sandbox.AuthConfig_TokenPair{TokenPair: &sandbox.TokenPairCredentials{Id: "secret-id", Secret: "secret-value"}}}}
	if _, ok := any(config).(proto.Message); !ok {
		t.Fatal("config is not a protobuf message")
	}
}

func TestProviderRequiresInitializationContract(t *testing.T) {
	// A name-only implementation no longer satisfies the initialization contract.
	typeNameOnly := reflect.TypeOf(nameOnly{})
	if typeNameOnly.Implements(reflect.TypeFor[sandbox.Provider]()) {
		t.Fatal("name-only provider satisfies factory contract")
	}
	if reflect.TypeOf(sandbox.NewClient).In(0) != reflect.TypeFor[sandbox.Provider]() {
		t.Fatal("constructor does not require Provider")
	}
	if reflect.TypeOf(sandbox.NewClient).In(1) != reflect.TypeFor[*sandbox.Config]() {
		t.Fatal("constructor does not require *Config")
	}
}

type nameOnly struct{}

func (nameOnly) Name() string { return "random" }

type cleanupBackend struct {
	name       string
	closed     int
	closeError error
}

func (p *cleanupBackend) Name() string { return p.name }
func (*cleanupBackend) Create(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	return nil, nil
}
func (p *cleanupBackend) Close(context.Context) error { p.closed++; return p.closeError }
func TestClientCleanupAndInvalidBackendIdentity(t *testing.T) {
	expected := errors.New("cleanup failure")
	provider := &cleanupBackend{name: "other", closeError: expected}
	if client, err := sandbox.NewClient(&testProvider{provider: provider}, &sandbox.Config{}); client != nil || err == nil || !errors.Is(err, expected) {
		t.Fatal("invalid backend was not rejected with cleanup error")
	}
	if provider.closed != 1 {
		t.Fatal("invalid initialized provider was not closed")
	}
	provider = &cleanupBackend{name: "test", closeError: expected}
	client, err := sandbox.NewClient(&testProvider{provider: provider}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Close(context.Background()); !errors.Is(err, expected) {
		t.Fatal("cleanup error changed")
	}
	if provider.closed != 1 {
		t.Fatal("owned SDK cleanup not delegated")
	}
	if err := client.Close(context.Background()); err != nil || provider.closed != 1 {
		t.Fatalf("close is not idempotent: %v, calls=%d", err, provider.closed)
	}
	_, err = client.Create(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindFailedPrecondition {
		t.Fatalf("closed client accepted an operation: %v", err)
	}
}

type managementBackend struct{}

func (*managementBackend) Name() string                { return "test" }
func (*managementBackend) Close(context.Context) error { return nil }

func TestCreationIsAnOptionalBackendCapability(t *testing.T) {
	client, err := sandbox.NewClient(&testProvider{provider: &managementBackend{}}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Create(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Operation != "create" {
		t.Fatalf("management-only backend returned %v", err)
	}
}
func TestCallerDeadlineSurvivesDisabledKitTimeout(t *testing.T) {
	provider := &creationBackend{call: func(ctx context.Context, r *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("caller deadline removed")
		}
		return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "id"}}, nil
	}}
	client, err := sandbox.NewClient(&testProvider{provider: provider}, &sandbox.Config{Timeout: sandbox.Duration(time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := client.Create(ctx, &sandbox.CreateOptions{Provisioning: &sandbox.ProvisioningOptions{Timeout: sandbox.Value(time.Duration(0))}}); err != nil {
		t.Fatal(err)
	}
}

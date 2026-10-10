package sandbox_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
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
	return sandbox.NewClient(sandbox.Config{Provider: &testProvider{provider: provider}})
}
func TestClientInitializesFromCommonConfiguration(t *testing.T) {
	factory := &testProvider{provider: &creationBackend{}}
	config := sandbox.Config{Provider: factory, Auth: &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: "test-key"}}, Endpoint: sandbox.Value("https://example.invalid"), Region: sandbox.Value("region"), Timeout: sandbox.Value(time.Second)}
	client, err := sandbox.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	if factory.calls != 1 || factory.config.GetEndpoint() != config.GetEndpoint() || factory.config.GetRegion() != config.GetRegion() || client.ProviderName() != "test" {
		t.Fatal("initialization settings were lost")
	}
}
func TestClientRejectsInvalidConfigurationBeforeInitialization(t *testing.T) {
	for _, config := range []sandbox.Config{
		{Timeout: sandbox.Value(-time.Second)}, {Endpoint: sandbox.Value("  ")}, {Endpoint: sandbox.Value("ftp://example.invalid")}, {Endpoint: sandbox.Value("not-a-url")},
		{Auth: &sandbox.AuthConfig{}},
		{Auth: &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: "key"}, BearerToken: &sandbox.BearerTokenCredentials{Token: "token"}}},
		{Auth: &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: "id"}}},
		{Auth: &sandbox.AuthConfig{OAuthRefresh: &sandbox.OAuthCredentials{RefreshToken: "refresh", ClientID: "id"}}},
	} {
		factory := &testProvider{provider: &creationBackend{}}
		config.Provider = factory
		if client, err := sandbox.NewClient(config); err == nil || client != nil {
			t.Fatal("invalid config accepted")
		}
		if factory.calls != 0 {
			t.Fatal("invalid config initialized SDK")
		}
	}
}
func TestClientRejectsNilFactoriesAndProviders(t *testing.T) {
	var typedNil *testProvider
	for _, config := range []sandbox.Config{{}, {Provider: typedNil}, {Provider: &testProvider{}}, {Provider: &testProvider{provider: (*creationBackend)(nil)}}} {
		if client, err := sandbox.NewClient(config); err == nil || client != nil {
			t.Fatal("nil factory/provider accepted")
		}
	}
}
func TestInitializationErrorsPassThrough(t *testing.T) {
	expected := errors.New("SDK failure")
	factory := &testProvider{err: expected}
	if _, err := sandbox.NewClient(sandbox.Config{Provider: factory}); !errors.Is(err, expected) {
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
			config := sandbox.Config{Provider: &testProvider{provider: provider}, Timeout: tc.clientTimeout}
			client, err := sandbox.NewClient(config)
			if err != nil {
				t.Fatal(err)
			}
			// Constructor snapshots defaults rather than retaining caller-owned pointers.
			if config.Timeout != nil {
				*config.Timeout = 0
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
func TestCredentialFieldsAreNotJSONSerialized(t *testing.T) {
	config := sandbox.Config{Auth: &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: "secret-id", Secret: "secret-value"}}}
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "secret-id") || strings.Contains(string(data), "secret-value") {
		t.Fatal("credentials serialized")
	}
}

func TestProviderRequiresInitializationContract(t *testing.T) {
	// A name-only implementation no longer satisfies the initialization contract.
	typeNameOnly := reflect.TypeOf(nameOnly{})
	if typeNameOnly.Implements(reflect.TypeFor[sandbox.Provider]()) {
		t.Fatal("name-only provider satisfies factory contract")
	}
	if reflect.TypeOf(sandbox.NewClient).In(0) != reflect.TypeFor[sandbox.Config]() {
		t.Fatal("constructor does not require Config")
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
	if client, err := sandbox.NewClient(sandbox.Config{Provider: &testProvider{provider: provider}}); client != nil || err == nil || !errors.Is(err, expected) {
		t.Fatal("invalid backend was not rejected with cleanup error")
	}
	if provider.closed != 1 {
		t.Fatal("invalid initialized provider was not closed")
	}
	provider = &cleanupBackend{name: "test", closeError: expected}
	client, err := sandbox.NewClient(sandbox.Config{Provider: &testProvider{provider: provider}})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Close(context.Background()); !errors.Is(err, expected) {
		t.Fatal("cleanup error changed")
	}
	if provider.closed != 1 {
		t.Fatal("owned SDK cleanup not delegated")
	}
}
func TestCallerDeadlineSurvivesDisabledKitTimeout(t *testing.T) {
	provider := &creationBackend{call: func(ctx context.Context, r *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("caller deadline removed")
		}
		return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "id"}}, nil
	}}
	client, err := sandbox.NewClient(sandbox.Config{Provider: &testProvider{provider: provider}, Timeout: sandbox.Value(time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := client.Create(ctx, &sandbox.CreateOptions{Provisioning: &sandbox.ProvisioningOptions{Timeout: sandbox.Value(time.Duration(0))}}); err != nil {
		t.Fatal(err)
	}
}

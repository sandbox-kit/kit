package modal

import (
	"context"
	"errors"
	"testing"
	"time"

	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type appsStub struct {
	sdk.AppService
	calls *int
}

func (s appsStub) FromName(context.Context, string, *sdk.AppFromNameParams) (*sdk.App, error) {
	*s.calls++
	return &sdk.App{}, nil
}

type imagesStub struct{ sdk.ImageService }

func (imagesStub) FromRegistry(string, *sdk.ImageFromRegistryParams) *sdk.Image { return &sdk.Image{} }

type proxiesStub struct {
	sdk.ProxyService
	calls *int
}

func (s proxiesStub) FromName(_ context.Context, name string, _ *sdk.ProxyFromNameParams) (*sdk.Proxy, error) {
	*s.calls++
	return &sdk.Proxy{ProxyID: "resolved-" + name}, nil
}

type sandboxesStub struct {
	sdk.SandboxService
	got **sdk.SandboxCreateParams
	err error
}

func (s sandboxesStub) FromID(_ context.Context, id string, _ *sdk.SandboxFromIDParams) (*sdk.Sandbox, error) {
	return &sdk.Sandbox{SandboxID: id}, nil
}

func (s sandboxesStub) Create(_ context.Context, _ *sdk.App, _ *sdk.Image, p *sdk.SandboxCreateParams) (*sdk.Sandbox, error) {
	*s.got = p
	if s.err != nil {
		return nil, s.err
	}
	return &sdk.Sandbox{SandboxID: "created"}, nil
}
func modalConfig() *sandbox.CreateOptions {
	return &sandbox.CreateOptions{Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "python:3.11"}}, Isolation: &sandbox.IsolationConfig{Kind: sandbox.IsolationKindContainer}, Resources: &sandbox.Resources{CPUCores: sandbox.Value(1.5), MemoryMiB: sandbox.Value(uint64(2048))}, Environment: map[string]string{"MODE": "test"}}
}
func TestTypedSDKMappingAndCommonHandle(t *testing.T) {
	calls := 0
	proxyCalls := 0
	var got *sdk.SandboxCreateParams
	native := &sdk.Client{Apps: appsStub{calls: &calls}, Images: imagesStub{}, Proxies: proxiesStub{calls: &proxyCalls}, Sandboxes: sandboxesStub{got: &got}}
	adapter := &backend{client: native, scope: &sandbox.Scope{AppName: sandbox.Value("my-app")}}
	if adapter.client != native {
		t.Fatal("SDK replaced")
	}
	client, _ := sandbox.NewClient(&operationProvider{adapter: adapter}, &sandbox.Config{})
	config := modalConfig()
	config.Network = &sandbox.NetworkConfig{ProxyReference: sandbox.Value("static-egress")}
	handle, err := client.Create(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	if handle.ID() != "created" || handle.ProviderName() != "modal" || handle.Info().GetIsolation().Kind != sandbox.IsolationKindContainer || calls != 1 || proxyCalls != 1 || got.Proxy == nil || got.Proxy.ProxyID != "resolved-static-egress" || got.CPU != 1.5 || got.MemoryMiB != 2048 || got.Env["MODE"] != "test" {
		t.Fatal("SDK mapping or response normalization failed")
	}
}

func TestAttachIsDistinctFromGet(t *testing.T) {
	var got *sdk.SandboxCreateParams
	adapter := &backend{client: &sdk.Client{Sandboxes: sandboxesStub{got: &got}}, scope: &sandbox.Scope{AppName: sandbox.Value("my-app")}}
	client, err := sandbox.NewClient(&operationProvider{adapter: adapter}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	handle, err := client.Attach(context.Background(), "sb-existing")
	if err != nil || handle.ID() != "sb-existing" {
		t.Fatalf("attach failed: %v", err)
	}
	if _, err := client.Get(context.Background(), "sb-existing"); err == nil {
		t.Fatal("Modal Get claimed provider retrieval")
	}
	if err := handle.Refresh(context.Background()); err == nil {
		t.Fatal("Modal Refresh claimed provider retrieval")
	}
}
func TestUnsupportedSettingsFailBeforeSDKCalls(t *testing.T) {
	calls := 0
	native := &sdk.Client{Apps: appsStub{calls: &calls}}
	adapter := &backend{client: native, scope: &sandbox.Scope{AppName: sandbox.Value("my-app")}}
	client, _ := sandbox.NewClient(&operationProvider{adapter: adapter}, &sandbox.Config{})
	config := modalConfig()
	config.Resources.DiskMiB = sandbox.Value(uint64(1024))
	if _, err := client.Create(context.Background(), config); err == nil {
		t.Fatal("unsupported disk ignored")
	}
	if calls != 0 {
		t.Fatal("SDK called before support validation")
	}
	config = modalConfig()
	config.Lifetime = &sandbox.LifetimePolicy{MaximumLifetime: sandbox.Value(time.Duration(0))}
	if _, err := client.Create(context.Background(), config); err == nil {
		t.Fatal("unlimited lifetime mapped to default")
	}
}
func TestSDKErrorPassThrough(t *testing.T) {
	calls := 0
	var got *sdk.SandboxCreateParams
	expected := errors.New("SDK failure")
	adapter := &backend{client: &sdk.Client{Apps: appsStub{calls: &calls}, Images: imagesStub{}, Sandboxes: sandboxesStub{got: &got, err: expected}}, scope: &sandbox.Scope{AppName: sandbox.Value("my-app")}}
	client, _ := sandbox.NewClient(&operationProvider{adapter: adapter}, &sandbox.Config{})
	if _, err := client.Create(context.Background(), modalConfig()); !errors.Is(err, expected) {
		t.Fatal("SDK error replaced")
	}
}

func TestGeneratedResourcePrecision(t *testing.T) {
	mapped, remaining, err := mapResources(&sandbox.Resources{CPUCores: sandbox.Value(0.5), MemoryMiB: sandbox.Value(uint64(512)), DiskMiB: sandbox.Value(uint64(1024))})
	if err != nil {
		t.Fatal(err)
	}
	if mapped.CPU != 0.5 || mapped.MemoryMiB != 512 {
		t.Fatal("resource mapping changed precision")
	}
	if remaining.CPUCores != nil || remaining.MemoryMiB != nil || remaining.DiskMiB == nil {
		t.Fatal("remaining unsupported intent was lost")
	}
	if _, _, err := mapResources(&sandbox.Resources{MemoryMiB: sandbox.Value(uint64(2147483648))}); err == nil {
		t.Fatal("overflow was accepted")
	}
}

// Operation tests inject SDK doubles behind a private factory,
// not a public SDK constructor.
type operationProvider struct{ adapter *backend }

func (f *operationProvider) Name() string { return f.adapter.Name() }
func (f *operationProvider) NewClient(*sandbox.Config) (sandbox.Backend, error) {
	return f.adapter, nil
}

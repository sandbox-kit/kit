package daytona

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	sdk "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestTypedSDKCreationAndResponseNormalization(t *testing.T) {
	var got map[string]any
	transport := transportFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != "POST" {
			t.Errorf("unexpected SDK request: %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Error(err)
		}
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"organizationId":"org","user":"daytona","env":{},"public":false,"networkBlockAll":false,"kvm":false,"gpu":0,"id":"created","name":"test","state":"started","target":"test-region","cpu":2,"memory":4,"disk":8,"labels":{},"toolboxProxyUrl":"http://localhost:9999"}`)), Request: r}, nil
	})
	native, err := sdk.NewClientWithConfig(&types.DaytonaConfig{APIKey: "test-only", APIUrl: "http://daytona.invalid", HTTPClient: &http.Client{Transport: transport}})
	if err != nil {
		t.Fatal(err)
	}
	defer native.Close(context.Background())
	adapter := &backend{client: native}
	if adapter.client != native {
		t.Fatal("SDK replaced")
	}
	client, _ := sandbox.NewClient(sandbox.Config{Provider: &operationProvider{adapter: adapter}})
	handle, err := client.Create(context.Background(), &sandbox.CreateOptions{Name: sandbox.Value("test"), Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "python:3.11"}}, Resources: &sandbox.Resources{CPUCores: sandbox.Value(2.0), MemoryMiB: sandbox.Value(uint64(4096)), DiskMiB: sandbox.Value(uint64(8192))}, Provisioning: &sandbox.ProvisioningOptions{WaitFor: sandbox.WaitConditionSubmitted}})
	if err != nil {
		t.Fatal(err)
	}
	if handle.ID() != "created" || handle.ProviderName() != "daytona" || handle.Info().GetResources().GetMemoryMiB() != 4096 {
		t.Fatal("response not normalized")
	}
	if got["memory"] != float64(4) || got["disk"] != float64(8) {
		t.Fatalf("incorrect unit conversion: %v", got)
	}
}
func TestUnrepresentableConfigRejectedBeforeCall(t *testing.T) {
	for _, config := range []*sandbox.CreateOptions{
		{Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "image"}}, Resources: &sandbox.Resources{CPUCores: sandbox.Value(1.5)}},
		{Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "image"}}, Resources: &sandbox.Resources{MemoryMiB: sandbox.Value(uint64(1000))}},
		{Resources: &sandbox.Resources{CPUCores: sandbox.Value(2.0)}},
		{Placement: &sandbox.Placement{Regions: []string{"region"}}},
		{Lifetime: &sandbox.LifetimePolicy{IdleStop: &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(time.Duration(0))}}},
	} {
		adapter := &backend{client: &sdk.Client{}}
		client, _ := sandbox.NewClient(sandbox.Config{Provider: &operationProvider{adapter: adapter}})
		if _, err := client.Create(context.Background(), config); err == nil {
			t.Fatal("unsupported config accepted")
		}
	}
}

func TestGeneratedResourceMappings(t *testing.T) {
	for _, tc := range []struct {
		name    string
		request *sandbox.Resources
		valid   bool
	}{
		{"whole GiB", &sandbox.Resources{CPUCores: sandbox.Value(2.0), MemoryMiB: sandbox.Value(uint64(4096))}, true},
		{"fractional CPU", &sandbox.Resources{CPUCores: sandbox.Value(0.5)}, false},
		{"fractional GiB", &sandbox.Resources{MemoryMiB: sandbox.Value(uint64(4097))}, false},
		{"maximum GiB", &sandbox.Resources{MemoryMiB: sandbox.Value(uint64(2147483647) * 1024)}, true},
		{"overflow GiB", &sandbox.Resources{MemoryMiB: sandbox.Value(uint64(2147483648) * 1024)}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mapped, remaining, err := mapResources(tc.request)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v err=%v", tc.valid, err)
			}
			if err == nil {
				if remaining.MemoryMiB != nil || remaining.CPUCores != nil {
					t.Fatal("mapped fields not consumed")
				}
				if tc.request.MemoryMiB != nil && uint64(mapped.Memory)*1024 != *tc.request.MemoryMiB {
					t.Fatal("incorrect conversion")
				}
			}
		})
	}
}

// Operation tests inject SDK doubles behind a private factory, not a public SDK constructor.
type operationProvider struct{ adapter *backend }

func (f *operationProvider) Name() string { return f.adapter.Name() }
func (f *operationProvider) NewClient(*sandbox.Config) (sandbox.Backend, error) {
	return f.adapter, nil
}

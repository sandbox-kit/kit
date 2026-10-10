package sandbox_test

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type creationBackend struct {
	call func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error)
}

func (*creationBackend) Name() string { return "test" }
func (p *creationBackend) Create(ctx context.Context, r *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	return p.call(ctx, r)
}

func TestCreateReturnsHandleAndProtectsCallerValues(t *testing.T) {
	request := &sandbox.CreateOptions{Environment: map[string]string{"MODE": "test"}, Provisioning: &sandbox.ProvisioningOptions{Timeout: sandbox.Value(time.Second)}}
	response := &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "sb-123", Labels: map[string]string{"team": "kit"}}}
	client, _ := newTestClient(&creationBackend{call: func(ctx context.Context, r *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("timeout not applied")
		}
		r.Environment["MODE"] = "changed"
		return response, nil
	}})
	instance, err := client.Create(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if instance.ID() != "sb-123" || instance.ProviderName() != "test" || request.Environment["MODE"] != "test" {
		t.Fatal("incorrect handle or mutated request")
	}
	response.Sandbox.ID = "changed"
	info := instance.Info()
	info.ID = "changed"
	info.Labels["team"] = "changed"
	if instance.ID() != "sb-123" || instance.Info().Labels["team"] != "kit" {
		t.Fatal("external mutation changed handle")
	}
}
func TestDefaultRequestAndNativeErrors(t *testing.T) {
	expected := errors.New("upstream failure")
	client, _ := newTestClient(&creationBackend{call: func(_ context.Context, r *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
		if !reflect.DeepEqual(r, &sandbox.CreateOptions{}) {
			t.Fatal("defaults replaced")
		}
		return nil, expected
	}})
	if _, err := client.Create(context.Background(), nil); !errors.Is(err, expected) {
		t.Fatal("native error replaced")
	}
}
func TestInvalidAndCancelledRequestsHaveNoProviderSideEffects(t *testing.T) {
	called := false
	client, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
		called = true
		return nil, nil
	}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Create(ctx, nil); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := client.Create(context.Background(), &sandbox.CreateOptions{Resources: &sandbox.Resources{CPUCores: sandbox.Value(math.NaN())}}); err == nil {
		t.Fatal("NaN accepted")
	}
	if called {
		t.Fatal("invalid request called provider")
	}
}
func TestPresenceAndExplicitZeroRoundTrip(t *testing.T) {
	request := &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{MaximumLifetime: sandbox.Value(time.Duration(0)), IdleStop: &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(time.Duration(0))}}, Network: &sandbox.NetworkConfig{Egress: sandbox.EgressModeRestricted, OutboundCIDRs: &sandbox.StringAllowlist{}}}
	if err := sandbox.ValidateCreateOptions(request); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	decoded := &sandbox.CreateOptions{}
	if err := json.Unmarshal(data, decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Lifetime.MaximumLifetime == nil || decoded.Lifetime.IdleStop.After == nil || decoded.Network.OutboundCIDRs == nil {
		t.Fatal("presence lost")
	}
}
func TestInvalidProviderResponses(t *testing.T) {
	for _, response := range []*sandbox.CreateResult{nil, {}, {Sandbox: &sandbox.SandboxInfo{}}, {Sandbox: &sandbox.SandboxInfo{ID: "id", Provider: "other"}}} {
		client, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
			return response, nil
		}})
		if _, err := client.Create(context.Background(), nil); err == nil {
			t.Fatal("invalid response accepted")
		}
	}

}

func TestCreatePreservesExplicitEmptyCollections(t *testing.T) {
	request := &sandbox.CreateOptions{Environment: map[string]string{}, Runtime: &sandbox.RuntimeConfig{Entrypoint: []string{}}, Placement: &sandbox.Placement{Regions: []string{}}}
	client, err := newTestClient(&creationBackend{call: func(_ context.Context, r *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
		if r.Environment == nil || r.Runtime.Entrypoint == nil || r.Placement.Regions == nil {
			t.Fatal("explicit empty collection became unspecified")
		}
		return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "empty-preserved"}}, nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Create(context.Background(), request); err != nil {
		t.Fatal(err)
	}
}

func TestRenamedFieldsPreserveJSONKeys(t *testing.T) {
	options := &sandbox.CreateOptions{Resources: &sandbox.Resources{CPUCores: sandbox.Value(2.0), MemoryMiB: sandbox.Value(uint64(4096))}, Provisioning: &sandbox.ProvisioningOptions{Timeout: sandbox.Value(time.Second)}}
	encoded, err := json.Marshal(options)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	if _, ok := fields["creation"]; !ok {
		t.Fatal("renaming Provisioning changed serialized creation key")
	}
	var resourceFields map[string]json.RawMessage
	if err := json.Unmarshal(fields["resources"], &resourceFields); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"cpu_cores", "memory_mib"} {
		if _, ok := resourceFields[key]; !ok {
			t.Fatalf("serialized key %s changed", key)
		}
	}
	config := sandbox.Config{Scope: &sandbox.Scope{AppName: sandbox.Value("app")}}
	encoded, err = json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	fields = nil
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatal(err)
	}
	if _, ok := fields["context"]; !ok {
		t.Fatal("renaming Scope changed serialized context key")
	}
}

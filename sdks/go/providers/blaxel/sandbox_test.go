package blaxel

import (
	"testing"
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestPlanCreateMapsRequiredFields(t *testing.T) {
	plan, err := planCreate(&sandbox.CreateOptions{
		Name: sandbox.Value("demo"),
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"},
		}},
		Labels:      map[string]string{"team": "kit"},
		Environment: map[string]string{"FOO": "bar"},
		Resources:   &sandbox.Resources{MemoryMib: sandbox.Value(uint64(4096))},
		Lifetime: &sandbox.LifetimePolicy{
			MaximumLifetime: sandbox.Duration(30 * time.Minute),
			IdleTerminate: &sandbox.AutomaticAction{
				Mode:  sandbox.PolicyModeAfter,
				After: sandbox.Duration(10 * time.Minute),
			},
		},
	}, "us-pdx-1")
	if err != nil {
		t.Fatal(err)
	}
	if plan.params.Sandbox.Metadata.Name != "demo" {
		t.Fatalf("name = %q", plan.params.Sandbox.Metadata.Name)
	}
	if !plan.params.Sandbox.Spec.Runtime.Image.Valid() || plan.params.Sandbox.Spec.Runtime.Image.Value != "blaxel/base-image:latest" {
		t.Fatalf("image = %#v", plan.params.Sandbox.Spec.Runtime.Image)
	}
	if !plan.params.Sandbox.Spec.Region.Valid() || plan.params.Sandbox.Spec.Region.Value != "us-pdx-1" {
		t.Fatalf("region = %#v", plan.params.Sandbox.Spec.Region)
	}
	if !plan.params.Sandbox.Spec.Runtime.Memory.Valid() || plan.params.Sandbox.Spec.Runtime.Memory.Value != 4295 {
		t.Fatalf("memory = %#v", plan.params.Sandbox.Spec.Runtime.Memory)
	}
	if !plan.params.Sandbox.Spec.Runtime.Ttl.Valid() || plan.params.Sandbox.Spec.Runtime.Ttl.Value != "30m" {
		t.Fatalf("ttl = %#v", plan.params.Sandbox.Spec.Runtime.Ttl)
	}
	if len(plan.params.Sandbox.Spec.Lifecycle.ExpirationPolicies) != 1 {
		t.Fatalf("policies = %#v", plan.params.Sandbox.Spec.Lifecycle.ExpirationPolicies)
	}
	policy := plan.params.Sandbox.Spec.Lifecycle.ExpirationPolicies[0]
	if policy.Type != blaxel.ExpirationPolicyTypeTtlIdle || !policy.Value.Valid() || policy.Value.Value != "10m" {
		t.Fatalf("idle policy = %#v", policy)
	}
	if plan.params.CreateIfNotExist.Valid() {
		t.Fatal("createIfNotExist must stay unset for native create")
	}
}

func TestBlaxelMemoryUnitConversion(t *testing.T) {
	mb, ok := mibToBlaxelMB(4096)
	if !ok || mb != 4295 {
		t.Fatalf("4096 MiB = (%d, %t) MB, want (4295, true)", mb, ok)
	}
	if mib := blaxelMBToMib(4295); mib != 4096 {
		t.Fatalf("4295 MB = %d MiB, want 4096", mib)
	}
}

func TestMapSandboxInfoUsesNameAsID(t *testing.T) {
	info, err := mapSandboxInfo(&blaxel.Sandbox{
		Metadata: blaxel.Metadata{
			Name:   "demo",
			Labels: map[string]string{"team": "kit"},
		},
		Spec:   blaxel.SandboxSpec{Region: "us-pdx-1", Runtime: blaxel.SandboxRuntime{Memory: 2048}},
		Status: blaxel.StatusDeployed,
		State:  blaxel.SandboxStateRunning,
	})
	if err != nil {
		t.Fatal(err)
	}
	if info.Id != "demo" || info.GetName() != "demo" || info.GetRegion() != "us-pdx-1" || info.GetProviderState() != "DEPLOYED" {
		t.Fatalf("info = %#v", info)
	}
	observed, err := sandboxInfo(&blaxel.Sandbox{
		Metadata: blaxel.Metadata{
			Name:        "demo",
			DisplayName: "Demo",
			CreatedAt:   "2026-01-02T03:04:05Z",
			UpdatedAt:   "2026-01-02T04:05:06Z",
		},
		Spec:       blaxel.SandboxSpec{Runtime: blaxel.SandboxRuntime{Memory: 2048}},
		Status:     blaxel.StatusDeployed,
		State:      blaxel.SandboxStateRunning,
		LastUsedAt: "2026-01-02T05:06:07Z",
		ExpiresIn:  3600,
	})
	if err != nil {
		t.Fatal(err)
	}
	if observed.State != sandbox.SandboxStateRunning || observed.Resources.GetMemoryMib() != 1953 {
		t.Fatalf("observed = %#v", observed)
	}
	// Blaxel keeps SandboxInfo to YAML mapSandboxInfo fields; no metadata bag or createdAt.
	if observed.ProviderMetadata != nil || observed.CreatedAt != nil {
		t.Fatalf("unexpected metadata = %#v createdAt = %#v", observed.ProviderMetadata, observed.CreatedAt)
	}

	standby, err := sandboxInfo(&blaxel.Sandbox{
		Metadata: blaxel.Metadata{Name: "demo"},
		Status:   blaxel.StatusDeployed,
		State:    blaxel.SandboxStateStandby,
	})
	if err != nil {
		t.Fatal(err)
	}
	if standby.State != sandbox.SandboxStatePaused {
		t.Fatalf("standby state = %#v", standby.State)
	}

	deleted, err := sandboxInfo(&blaxel.Sandbox{
		Metadata: blaxel.Metadata{Name: "demo"},
		Status:   blaxel.Status("DELETED"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if deleted.State != sandbox.SandboxStateDeleted || deleted.GetTerminal() != true {
		t.Fatalf("deleted = %#v", deleted)
	}
}

func TestBlaxelCreateTimeoutSeconds(t *testing.T) {
	for _, tc := range []struct {
		timeout *time.Duration
		want    int
	}{
		{nil, 0},
		{durationPtr(0), 0},
		{durationPtr(30 * time.Second), 30},
		{durationPtr(time.Second), 1},
		{durationPtr(50 * time.Second), 50},
		{durationPtr(1500 * time.Millisecond), 0},
		{durationPtr(120 * time.Second), 0},
	} {
		request := &sandbox.CreateOptions{}
		if tc.timeout != nil {
			request.Provisioning = &sandbox.ProvisioningOptions{
				Timeout: sandbox.Duration(*tc.timeout),
			}
		}
		if got := blaxelCreateTimeoutSeconds(request); got != tc.want {
			t.Fatalf("timeout %#v: got %d want %d", tc.timeout, got, tc.want)
		}
	}
}

func durationPtr(value time.Duration) *time.Duration {
	return &value
}

func TestFormatBlaxelDuration(t *testing.T) {
	for _, tc := range []struct {
		in      time.Duration
		minimum time.Duration
		want    string
		err     string
	}{
		{30 * time.Second, time.Second, "30s", ""},
		{5 * time.Minute, 5 * time.Minute, "5m", ""},
		{2 * time.Hour, 5 * time.Minute, "2h", ""},
		{48 * time.Hour, 5 * time.Minute, "2d", ""},
		{30 * time.Second, time.Minute, "", "duration must be at least 1m"},
		{4 * time.Minute, 5 * time.Minute, "", "duration must be at least 5m"},
		{1500 * time.Millisecond, time.Second, "", "duration must be a positive whole number of seconds"},
		{0, time.Second, "", "duration must be a positive whole number of seconds"},
	} {
		got, err := formatBlaxelDuration(tc.in, tc.minimum)
		if tc.err == "" {
			if err != nil || got != tc.want {
				t.Fatalf("%v (min %v): got %q err %v", tc.in, tc.minimum, got, err)
			}
			continue
		}
		if err == nil || err.Error() != tc.err {
			t.Fatalf("%v (min %v): got err %v want %q", tc.in, tc.minimum, err, tc.err)
		}
	}
}

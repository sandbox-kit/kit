package modal

import (
	"math"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestCPUNativePrecisionAndRange(t *testing.T) {
	for _, tc := range []struct {
		value float64
		valid bool
	}{
		{0.001, true}, {1.001, true}, {1.5, true}, {float64(math.MaxUint32) / 1000, true},
		{0.0001, false}, {1.0005, false}, {float64(math.MaxUint32+1) / 1000, false}, {math.NaN(), false}, {math.Inf(1), false},
	} {
		mapped, _, err := mapResources(&sandbox.Resources{CPUCores: sandbox.Value(tc.value)})
		if (err == nil) != tc.valid {
			t.Fatalf("CPU %v: err=%v", tc.value, err)
		}
		if tc.valid && uint32(mapped.CPU*1000) != uint32(math.Round(tc.value*1000)) {
			t.Fatalf("CPU %v truncated to %v", tc.value, mapped.CPU*1000)
		}
		_, _, err = mapResources(&sandbox.Resources{CPULimitCores: sandbox.Value(tc.value)})
		if (err == nil) != tc.valid {
			t.Fatalf("CPU limit %v: err=%v", tc.value, err)
		}
	}
}

func TestNativeDurationBounds(t *testing.T) {
	for _, tc := range []struct {
		value time.Duration
		valid bool
	}{
		{time.Second, true}, {time.Duration(math.MaxUint32) * time.Second, true},
		{24 * time.Hour, true}, {24*time.Hour + time.Second, true},
		{0, false}, {-time.Second, false}, {time.Millisecond, false}, {time.Duration(math.MaxUint32+1) * time.Second, false},
	} {
		mapped, _, err := mapLifetimeDuration(&sandbox.LifetimePolicy{MaximumLifetime: sandbox.Value(tc.value)})
		if (err == nil) != (tc.valid && tc.value <= 24*time.Hour) {
			t.Fatalf("lifetime %v: %v", tc.value, err)
		}
		if err == nil && mapped.Timeout != tc.value {
			t.Fatal("duration changed")
		}
		_, _, err = mapIdleDuration(&sandbox.AutomaticAction{After: sandbox.Value(tc.value)})
		if (err == nil) != tc.valid {
			t.Fatalf("idle %v: %v", tc.value, err)
		}
	}
	if mapped, _, err := mapLifetimeDuration(&sandbox.LifetimePolicy{}); err != nil || mapped.Timeout != 0 {
		t.Fatal("absent duration changed")
	}
}

func TestStateCapturePreservesPresenceAndOwnership(t *testing.T) {
	var state backend
	captureBackendState(&state, &sandbox.Config{})
	if state.scope != nil || state.region != nil {
		t.Fatal("absent values became supplied")
	}
	config := &sandbox.Config{Scope: &sandbox.Scope{AppName: sandbox.Value("app"), OrganizationID: sandbox.Value("org"), ProjectID: sandbox.Value("project")}, Region: sandbox.Value("")}
	captureBackendState(&state, config)
	if state.scope.Environment != nil || state.region == nil || *state.region != "" || state.scope.GetOrganizationID() != "org" || state.scope.GetProjectID() != "project" {
		t.Fatal("schema fields or presence lost")
	}
	*config.Scope.AppName = "changed"
	*config.Region = "changed"
	if state.scope.GetAppName() != "app" || *state.region != "" {
		t.Fatal("capture aliases caller data")
	}
}

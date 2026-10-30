package daytona

import (
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestRejectedCreateSettings(t *testing.T) {
	for _, tc := range []struct {
		name    string
		request *sandbox.CreateOptions
		kind    sandbox.ErrorKind
		field   string
	}{
		{
			name:    "placement",
			request: &sandbox.CreateOptions{Placement: &sandbox.Placement{Regions: []string{"eu"}}},
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldPlacement,
		},
		{
			name:    "idle terminate",
			request: &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{IdleTerminate: &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled}}},
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldLifetimeIdleTerminate,
		},
		{
			name:    "verbose",
			request: &sandbox.CreateOptions{Observability: &sandbox.ObservabilityConfig{Verbose: sandbox.Value(true)}},
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldObservabilityVerbose,
		},
		{
			name:    "workload identity",
			request: &sandbox.CreateOptions{Security: &sandbox.SecurityConfig{WorkloadIdentity: sandbox.Value(true)}},
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldSecurityWorkloadIdentity,
		},
		{
			name:    "runtime version",
			request: &sandbox.CreateOptions{Runtime: &sandbox.RuntimeConfig{Version: sandbox.Value("3.12")}},
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldRuntime,
		},
		{
			name:    "isolation class",
			request: &sandbox.CreateOptions{Isolation: &sandbox.IsolationConfig{Kind: sandbox.IsolationKindContainer}},
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldIsolation,
		},
		{
			name:    "private network",
			request: &sandbox.CreateOptions{Network: &sandbox.NetworkConfig{PrivateNetwork: sandbox.Value(true)}},
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldNetwork,
		},
		{
			name:    "maximum lifetime",
			request: &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{MaximumLifetime: sandbox.Duration(30 * time.Second)}},
			kind:    sandbox.ErrorKindInvalidArgument,
			field:   sandbox.CreateFieldLifetimeMaximumLifetime,
		},
		{
			name:    "fractional create timeout",
			request: &sandbox.CreateOptions{Provisioning: &sandbox.ProvisioningOptions{Timeout: sandbox.Duration(1500 * time.Millisecond)}},
			kind:    sandbox.ErrorKindInvalidArgument,
			field:   sandbox.CreateFieldProvisioningTimeout,
		},
		{
			name: "cpu limit",
			request: &sandbox.CreateOptions{
				Source:    alpineImage(),
				Resources: &sandbox.Resources{CpuLimitCores: sandbox.Value(2.0)},
			},
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldResources,
		},
		{
			name: "memory limit",
			request: &sandbox.CreateOptions{
				Source:    alpineImage(),
				Resources: &sandbox.Resources{MemoryLimitMib: sandbox.Value(uint64(2048))},
			},
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldResources,
		},
		{
			name:    "queue timeout",
			request: &sandbox.CreateOptions{Provisioning: &sandbox.ProvisioningOptions{QueueTimeout: sandbox.Duration(30 * time.Second)}},
			kind:    sandbox.ErrorKindInvalidArgument,
			field:   sandbox.CreateFieldProvisioningQueueTimeout,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := planCreate(tc.request)
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != tc.kind || detail.Field != tc.field {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
}

func alpineImage() *sandbox.SandboxSource {
	return &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}}}
}

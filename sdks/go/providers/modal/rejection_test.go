package modal

import (
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestRejectedCreateSettings(t *testing.T) {
	scope := &sandbox.Scope{AppName: sandbox.Value("app")}
	for _, tc := range []struct {
		name    string
		request *sandbox.CreateOptions
		kind    sandbox.ErrorKind
		field   string
	}{
		{
			name:    "disk",
			request: imaged(&sandbox.CreateOptions{Resources: &sandbox.Resources{DiskMib: sandbox.Value(uint64(1024))}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldResourcesDiskMib,
		},
		{
			name:    "placement target",
			request: imaged(&sandbox.CreateOptions{Placement: &sandbox.Placement{Target: sandbox.Value("us")}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldPlacementTarget,
		},
		{
			name:    "queue timeout",
			request: imaged(&sandbox.CreateOptions{Provisioning: &sandbox.ProvisioningOptions{QueueTimeout: sandbox.Duration(time.Minute)}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldProvisioningQueueTimeout,
		},
		{
			name:    "telemetry endpoint",
			request: imaged(&sandbox.CreateOptions{Observability: &sandbox.ObservabilityConfig{TelemetryEndpoint: sandbox.Value("https://otel.example")}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldObservabilityTelemetryEndpoint,
		},
		{
			name:    "public access",
			request: imaged(&sandbox.CreateOptions{Security: &sandbox.SecurityConfig{PublicAccess: sandbox.Value(true)}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldSecurityPublicAccess,
		},
		{
			name:    "linked sandbox",
			request: imaged(&sandbox.CreateOptions{Network: &sandbox.NetworkConfig{LinkedSandbox: sandbox.Value("other")}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldNetwork,
		},
		{
			name:    "ephemeral lifetime",
			request: imaged(&sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{Ephemeral: sandbox.Value(true)}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldLifetime,
		},
		{
			name:    "runtime language",
			request: imaged(&sandbox.CreateOptions{Runtime: &sandbox.RuntimeConfig{Language: sandbox.Value("python")}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldRuntime,
		},
		{
			name:    "nested virtualization",
			request: imaged(&sandbox.CreateOptions{Isolation: &sandbox.IsolationConfig{NestedVirtualization: sandbox.Value(true)}}),
			kind:    sandbox.ErrorKindUnsupported,
			field:   sandbox.CreateFieldIsolation,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := planCreate(tc.request, scope, "")
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != tc.kind || detail.Field != tc.field {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
}

func TestRejectedClientSettings(t *testing.T) {
	for _, tc := range []struct {
		name   string
		config *sandbox.Config
		field  string
	}{
		{"organization", &sandbox.Config{Scope: &sandbox.Scope{OrganizationId: sandbox.Value("org")}}, sandbox.ConfigFieldScopeOrganizationId},
		{"project", &sandbox.Config{Scope: &sandbox.Scope{ProjectId: sandbox.Value("project")}}, sandbox.ConfigFieldScopeProjectId},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mapClientConfig(tc.config)
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != tc.field {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
}

func imaged(request *sandbox.CreateOptions) *sandbox.CreateOptions {
	request.Source = &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}}}
	return request
}

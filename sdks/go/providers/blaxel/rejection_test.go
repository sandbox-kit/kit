package blaxel

import (
	"errors"
	"testing"
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
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
			name:    "missing name",
			request: imaged(&sandbox.CreateOptions{}),
			kind:    sandbox.ErrorKindInvalidArgument,
			field:   sandbox.CreateFieldName,
		},
		{
			name:    "missing image",
			request: &sandbox.CreateOptions{Name: sandbox.Value("sandbox")},
			kind:    sandbox.ErrorKindInvalidArgument,
			field:   sandbox.CreateFieldSourceImage,
		},
		{
			name: "cpu cores",
			request: namedImage(&sandbox.CreateOptions{
				Resources: &sandbox.Resources{CpuCores: sandbox.Value(2.0)},
			}),
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldResourcesCpuCores,
		},
		{
			name: "disk",
			request: namedImage(&sandbox.CreateOptions{
				Resources: &sandbox.Resources{DiskMib: sandbox.Value(uint64(1024))},
			}),
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldResourcesDiskMib,
		},
		{
			name: "placement",
			request: namedImage(&sandbox.CreateOptions{
				Placement: &sandbox.Placement{Cloud: sandbox.Value("aws")},
			}),
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldPlacement,
		},
		{
			name: "security",
			request: namedImage(&sandbox.CreateOptions{
				Security: &sandbox.SecurityConfig{PublicAccess: sandbox.Value(true)},
			}),
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldSecurity,
		},
		{
			name: "private network",
			request: namedImage(&sandbox.CreateOptions{
				Network: &sandbox.NetworkConfig{PrivateNetwork: sandbox.Value(true)},
			}),
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldNetwork,
		},
		{
			name: "egress mode",
			request: namedImage(&sandbox.CreateOptions{
				Network: &sandbox.NetworkConfig{Egress: sandbox.EgressModeRestricted},
			}),
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldNetwork,
		},
		{
			name: "ephemeral",
			request: namedImage(&sandbox.CreateOptions{
				Lifetime: &sandbox.LifetimePolicy{Ephemeral: sandbox.Value(true)},
			}),
			kind:  sandbox.ErrorKindUnsupported,
			field: sandbox.CreateFieldLifetime,
		},
		{
			name: "fractional maximum lifetime",
			request: namedImage(&sandbox.CreateOptions{
				Lifetime: &sandbox.LifetimePolicy{MaximumLifetime: sandbox.Duration(1500 * time.Millisecond)},
			}),
			kind:  sandbox.ErrorKindInvalidArgument,
			field: sandbox.CreateFieldLifetimeMaximumLifetime,
		},
		{
			name: "maximum lifetime below floor",
			request: namedImage(&sandbox.CreateOptions{
				Lifetime: &sandbox.LifetimePolicy{MaximumLifetime: sandbox.Duration(4 * time.Minute)},
			}),
			kind:  sandbox.ErrorKindInvalidArgument,
			field: sandbox.CreateFieldLifetimeMaximumLifetime,
		},
		{
			name: "idle terminate below floor",
			request: namedImage(&sandbox.CreateOptions{
				Lifetime: &sandbox.LifetimePolicy{
					IdleTerminate: &sandbox.AutomaticAction{
						Mode:  sandbox.PolicyModeAfter,
						After: sandbox.Duration(30 * time.Second),
					},
				},
			}),
			kind:  sandbox.ErrorKindInvalidArgument,
			field: sandbox.CreateFieldLifetimeIdleTerminate,
		},
		{
			name: "reserved port",
			request: namedImage(&sandbox.CreateOptions{
				Network: &sandbox.NetworkConfig{
					Ports: []*sandbox.PortBinding{
						{Port: 443, Transport: sandbox.PortTransportTLS},
					},
				},
			}),
			kind:  sandbox.ErrorKindInvalidArgument,
			field: sandbox.CreateFieldNetworkPortsPort,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := planCreate(tc.request, "us-pdx-1")
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
		kind   sandbox.ErrorKind
	}{
		{
			name:   "endpoint",
			config: &sandbox.Config{Endpoint: sandbox.Value("https://example.invalid")},
			field:  sandbox.ConfigFieldEndpoint,
			kind:   sandbox.ErrorKindUnsupported,
		},
		{
			name:   "app",
			config: &sandbox.Config{Scope: &sandbox.Scope{AppName: sandbox.Value("app")}},
			field:  sandbox.ConfigFieldScopeAppName,
			kind:   sandbox.ErrorKindUnsupported,
		},
		{
			name:   "environment",
			config: &sandbox.Config{Scope: &sandbox.Scope{Environment: sandbox.Value("env")}},
			field:  sandbox.ConfigFieldScopeEnvironment,
			kind:   sandbox.ErrorKindUnsupported,
		},
		{
			name:   "organization",
			config: &sandbox.Config{Scope: &sandbox.Scope{OrganizationId: sandbox.Value("org")}},
			field:  sandbox.ConfigFieldScopeOrganizationId,
			kind:   sandbox.ErrorKindUnsupported,
		},
		{
			name:   "project",
			config: &sandbox.Config{Scope: &sandbox.Scope{ProjectId: sandbox.Value("project")}},
			field:  sandbox.ConfigFieldScopeProjectId,
			kind:   sandbox.ErrorKindUnsupported,
		},
		{
			name:   "token pair",
			config: &sandbox.Config{Auth: &sandbox.AuthConfig{Credentials: &sandbox.TokenPair{TokenPair: &sandbox.TokenPairCredentials{Id: "id", Secret: "secret"}}}},
			field:  sandbox.ConfigFieldAuth,
			kind:   sandbox.ErrorKindUnsupported,
		},
		{
			name:   "bearer",
			config: &sandbox.Config{Auth: &sandbox.AuthConfig{Credentials: &sandbox.BearerToken{BearerToken: &sandbox.BearerTokenCredentials{Token: "jwt"}}}},
			field:  sandbox.ConfigFieldAuth,
			kind:   sandbox.ErrorKindUnsupported,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mapClientConfig(tc.config)
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != tc.kind || detail.Field != tc.field {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
}

func TestMissingRegionRejected(t *testing.T) {
	_, err := planCreate(namedImage(&sandbox.CreateOptions{}), "")
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.ConfigFieldRegion {
		t.Fatalf("error = %v detail = %#v", err, detail)
	}
}

func TestCreateTimeoutRejected(t *testing.T) {
	_, err := planCreate(namedImage(&sandbox.CreateOptions{
		Provisioning: &sandbox.ProvisioningOptions{Timeout: sandbox.Duration(30 * time.Second)},
	}), "us-pdx-1")
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != sandbox.CreateFieldProvisioningTimeout {
		t.Fatalf("error = %v detail = %#v", err, detail)
	}
	if detail.Message != "sandbox-kit blaxel: Blaxel create timeout is not mapped" {
		t.Fatalf("message = %q", detail.Message)
	}
}

func imaged(request *sandbox.CreateOptions) *sandbox.CreateOptions {
	request.Source = &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"}}}
	return request
}

func namedImage(request *sandbox.CreateOptions) *sandbox.CreateOptions {
	request.Name = sandbox.Value("sandbox")
	return imaged(request)
}

func TestMappedNetworkPortsAndDomains(t *testing.T) {
	plan, err := planCreate(namedImage(&sandbox.CreateOptions{
		Network: &sandbox.NetworkConfig{
			OutboundDomains: &sandbox.StringAllowlist{Entries: []string{"example.com"}},
			Ports: []*sandbox.PortBinding{
				{Port: 3000, Transport: sandbox.PortTransportTLS},
				{Port: 3001, Transport: sandbox.PortTransportPlain},
			},
		},
	}), "us-pdx-1")
	if err != nil {
		t.Fatal(err)
	}
	if got := plan.params.Sandbox.Spec.Network.AllowedDomains; len(got) != 1 || got[0] != "example.com" {
		t.Fatalf("AllowedDomains = %#v", got)
	}
	ports := plan.params.Sandbox.Spec.Runtime.Ports
	if len(ports) != 2 || ports[0].Target != 3000 || ports[0].Protocol != blaxel.PortProtocolTls {
		t.Fatalf("Ports = %#v", ports)
	}
	if ports[1].Target != 3001 || ports[1].Protocol != blaxel.PortProtocolTcp {
		t.Fatalf("Ports[1] = %#v", ports[1])
	}
}

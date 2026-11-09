package blaxel

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
			name: "network",
			request: namedImage(&sandbox.CreateOptions{
				Network: &sandbox.NetworkConfig{PrivateNetwork: sandbox.Value(true)},
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

func imaged(request *sandbox.CreateOptions) *sandbox.CreateOptions {
	request.Source = &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"}}}
	return request
}

func namedImage(request *sandbox.CreateOptions) *sandbox.CreateOptions {
	request.Name = sandbox.Value("sandbox")
	return imaged(request)
}

package modal

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestCommonInitializationMappings(t *testing.T) {
	config := &sandbox.Config{
		Auth:  &sandbox.AuthConfig{Credentials: &sandbox.TokenPair{TokenPair: &sandbox.TokenPairCredentials{Id: "token-id", Secret: "token-secret"}}},
		Scope: &sandbox.Scope{AppName: sandbox.Value("my-app"), Environment: sandbox.Value("dev")}, Region: sandbox.Value("us-east"),
	}
	params, err := mapClientConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if params.TokenID != "token-id" || params.TokenSecret != "token-secret" || params.Environment != "dev" {
		t.Fatal("common settings not mapped")
	}
	config.Auth = &sandbox.AuthConfig{Credentials: &sandbox.OAuthRefresh{OauthRefresh: &sandbox.OAuthCredentials{RefreshToken: "refresh", ClientId: "client", ClientSecret: sandbox.Value("secret")}}}
	params, err = mapClientConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if params.OAuthCredentials.RefreshToken != "refresh" || params.OAuthCredentials.ClientID != "client" || params.OAuthCredentials.ClientSecret != "secret" {
		t.Fatal("OAuth settings not mapped")
	}
	oauth := config.Auth.GetOauthRefresh()
	oauth.ClientSecret = nil
	oauth.JwtKey = sandbox.Value("key")
	params, err = mapClientConfig(config)
	if err != nil || params.OAuthCredentials.JWTKey != "key" {
		t.Fatal("JWT-key OAuth mapping failed")
	}
}
func TestUnsupportedClientSettingsRejectedBeforeSDKInitialization(t *testing.T) {
	for _, config := range []*sandbox.Config{
		{Endpoint: sandbox.Value("https://example.invalid")},
		{Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{ApiKey: &sandbox.APIKeyCredentials{Key: "test"}}}},
		{Auth: &sandbox.AuthConfig{Credentials: &sandbox.BearerToken{BearerToken: &sandbox.BearerTokenCredentials{Token: "test"}}}},
		{Scope: &sandbox.Scope{OrganizationId: sandbox.Value("org")}},
	} {
		if provider, err := New().NewClient(config); err == nil || provider != nil {
			t.Fatal("unsupported client setting initialized SDK")
		}
	}
	_, err := mapClientConfig(&sandbox.Config{Endpoint: sandbox.Value("https://example.invalid")})
	if err == nil || !strings.Contains(err.Error(), "Endpoint") {
		t.Fatal("endpoint exception not explained")
	}
}
func TestUnsupportedFieldsUseCanonicalPaths(t *testing.T) {
	scope := &sandbox.Scope{AppName: sandbox.Value("app")}
	image := &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "image"}}}
	for _, tc := range []struct {
		name    string
		request *sandbox.CreateOptions
		field   string
	}{
		{
			name:    "disk",
			request: &sandbox.CreateOptions{Source: image, Resources: &sandbox.Resources{DiskMib: sandbox.Value(uint64(1024))}},
			field:   sandbox.CreateFieldResourcesDiskMib,
		},
		{
			name:    "placement target",
			request: &sandbox.CreateOptions{Source: image, Placement: &sandbox.Placement{Target: sandbox.Value("us")}},
			field:   sandbox.CreateFieldPlacementTarget,
		},
		{
			name:    "queue timeout",
			request: &sandbox.CreateOptions{Source: image, Provisioning: &sandbox.ProvisioningOptions{QueueTimeout: sandbox.Duration(time.Minute)}},
			field:   sandbox.CreateFieldProvisioningQueueTimeout,
		},
		{
			name:    "telemetry endpoint",
			request: &sandbox.CreateOptions{Source: image, Observability: &sandbox.ObservabilityConfig{TelemetryEndpoint: sandbox.Value("https://otel.example")}},
			field:   sandbox.CreateFieldObservabilityTelemetryEndpoint,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := planCreate(tc.request, scope, "")
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != tc.field {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
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

func TestClientPlacementDefaultsAndOverrides(t *testing.T) {
	scope := &sandbox.Scope{AppName: sandbox.Value("app")}
	request := &sandbox.CreateOptions{Source: &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "image"}}}}
	plan, err := planCreate(request, scope, "default-region")
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.params.Regions) != 1 || plan.params.Regions[0] != "default-region" {
		t.Fatal("client region missing")
	}
	request.Placement = &sandbox.Placement{Regions: []string{"override"}}
	plan, err = planCreate(request, scope, "default-region")
	if err != nil || plan.params.Regions[0] != "override" {
		t.Fatal("per-create region not respected")
	}
	request.Placement.Regions = []string{}
	plan, err = planCreate(request, scope, "default-region")
	if err != nil || plan.params.Regions == nil || len(plan.params.Regions) != 0 {
		t.Fatal("explicit empty placement lost")
	}
}

func TestManagedProxyIsRetainedForNativeLookup(t *testing.T) {
	request := &sandbox.CreateOptions{
		Source:  &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "image"}}},
		Network: &sandbox.NetworkConfig{ProxyReference: sandbox.Value("static-egress")},
	}
	plan, err := planCreate(request, &sandbox.Scope{AppName: sandbox.Value("app")}, "")
	if err != nil {
		t.Fatal(err)
	}
	if plan.proxy != "static-egress" {
		t.Fatalf("proxy lookup name = %q", plan.proxy)
	}
}

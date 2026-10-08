package modal

import (
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestCommonInitializationMappings(t *testing.T) {
	config := &sandbox.Config{Auth: &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: "token-id", Secret: "token-secret"}}, Scope: &sandbox.Scope{AppName: sandbox.Value("my-app"), Environment: sandbox.Value("dev")}, Region: sandbox.Value("us-east")}
	params, err := mapClientConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if params.TokenID != "token-id" || params.TokenSecret != "token-secret" || params.Environment != "dev" {
		t.Fatal("common settings not mapped")
	}
	config.Auth = &sandbox.AuthConfig{OAuthRefresh: &sandbox.OAuthCredentials{RefreshToken: "refresh", ClientID: "client", ClientSecret: sandbox.Value("secret")}}
	params, err = mapClientConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if params.OAuthCredentials.RefreshToken != "refresh" || params.OAuthCredentials.ClientID != "client" || params.OAuthCredentials.ClientSecret != "secret" {
		t.Fatal("OAuth settings not mapped")
	}
	config.Auth.OAuthRefresh.ClientSecret = nil
	config.Auth.OAuthRefresh.JWTKey = sandbox.Value("key")
	params, err = mapClientConfig(config)
	if err != nil || params.OAuthCredentials.JWTKey != "key" {
		t.Fatal("JWT-key OAuth mapping failed")
	}
}
func TestUnsupportedClientSettingsRejectedBeforeSDKInitialization(t *testing.T) {
	for _, config := range []*sandbox.Config{
		{Endpoint: sandbox.Value("https://example.invalid")},
		{Auth: &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: "test"}}},
		{Auth: &sandbox.AuthConfig{BearerToken: &sandbox.BearerTokenCredentials{Token: "test"}}},
		{Scope: &sandbox.Scope{OrganizationID: sandbox.Value("org")}},
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
func TestClientPlacementDefaultsAndOverrides(t *testing.T) {
	scope := &sandbox.Scope{AppName: sandbox.Value("app")}
	request := &sandbox.CreateOptions{Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "image"}}}
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

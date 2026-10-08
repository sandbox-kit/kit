package daytona

import (
	"context"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestCommonInitializationMappings(t *testing.T) {
	config := &sandbox.Config{Auth: &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: "api-key"}}, Endpoint: sandbox.Value("https://example.invalid"), Region: sandbox.Value("region"), Timeout: sandbox.Value(time.Second)}
	params, err := mapClientConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if params.APIKey != "api-key" || params.APIUrl != "https://example.invalid" || params.Target != "region" || params.Timeout != nil {
		t.Fatal("client settings were not mapped with correct timeout semantics")
	}
	config.Auth = &sandbox.AuthConfig{BearerToken: &sandbox.BearerTokenCredentials{Token: "jwt"}}
	config.Scope = &sandbox.Scope{OrganizationID: sandbox.Value("org")}
	params, err = mapClientConfig(config)
	if err != nil || params.JWTToken != "jwt" || params.OrganizationID != "org" {
		t.Fatal("bearer authentication mapping failed")
	}
}
func TestUnsupportedClientSettingsRejectedBeforeSDKInitialization(t *testing.T) {
	for _, config := range []*sandbox.Config{
		{Auth: &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: "id", Secret: "secret"}}},
		{Scope: &sandbox.Scope{AppName: sandbox.Value("app")}},
		{Scope: &sandbox.Scope{Environment: sandbox.Value("env")}},
	} {
		if provider, err := New().NewClient(config); err == nil || provider != nil {
			t.Fatal("unsupported client setting initialized SDK")
		}
	}
}
func TestProviderInitializesAndOwnsSDKWithoutRequests(t *testing.T) {
	// The SDK constructor creates local clients; authentication is not probed here.
	t.Setenv("DAYTONA_JWT_TOKEN", "")
	t.Setenv("DAYTONA_API_KEY", "")
	client, err := sandbox.NewClient(sandbox.Config{Provider: New(), Auth: &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: "test-only"}}, Endpoint: sandbox.Value("https://example.invalid"), Region: sandbox.Value("region")})
	if err != nil {
		t.Fatal(err)
	}
	if client.ProviderName() != "daytona" {
		t.Fatal("wrong provider")
	}
	if err := client.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

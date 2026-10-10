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

func TestEndpointEnvironmentResolutionAndValidation(t *testing.T) {
	for _, tc := range []struct {
		name, api, server string
		explicit          *string
		want              string
		valid             bool
	}{
		{name: "default", want: "https://app.daytona.io/api", valid: true},
		{name: "API URL", api: "https://api.example.test", want: "https://api.example.test", valid: true},
		{name: "server fallback", server: "https://server.example.test", want: "https://server.example.test", valid: true},
		{name: "API precedence", api: "https://api.example.test", server: "http://bad.example.test", want: "https://api.example.test", valid: true},
		{name: "insecure API", api: "http://bad.example.test"},
		{name: "insecure server", server: "http://bad.example.test"},
		{name: "explicit precedence", api: "http://bad.example.test", explicit: sandbox.Value("https://explicit.example.test"), want: "https://explicit.example.test", valid: true},
		{name: "local testing", api: "http://127.0.0.1:1234", want: "http://127.0.0.1:1234", valid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("DAYTONA_API_URL", tc.api)
			t.Setenv("DAYTONA_SERVER_URL", tc.server)
			config := &sandbox.Config{Endpoint: tc.explicit}
			params, err := mapClientConfig(config)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v error=%v", tc.valid, err)
			}
			if tc.valid && params.APIUrl != tc.want {
				t.Fatalf("got %q want %q", params.APIUrl, tc.want)
			}
			if config.Endpoint != tc.explicit {
				t.Fatal("caller config was mutated")
			}
		})
	}
}

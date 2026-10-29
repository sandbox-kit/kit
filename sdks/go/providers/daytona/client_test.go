package daytona

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestCommonInitializationMappings(t *testing.T) {
	config := &sandbox.Config{
		Auth:     &sandbox.AuthConfig{Credentials: &sandbox.APIKey{ApiKey: &sandbox.APIKeyCredentials{Key: "api-key"}}},
		Endpoint: sandbox.Value("https://example.invalid"),
		Region:   sandbox.Value("region"),
		Timeout:  sandbox.Duration(time.Second),
	}
	params, err := mapClientConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if params.APIKey != "api-key" || params.APIUrl != "https://example.invalid" || params.Target != "region" || params.Timeout != nil {
		t.Fatal("client settings were not mapped with correct timeout semantics")
	}
	config.Auth = &sandbox.AuthConfig{Credentials: &sandbox.BearerToken{BearerToken: &sandbox.BearerTokenCredentials{Token: "jwt"}}}
	config.Scope = &sandbox.Scope{OrganizationId: sandbox.Value("org")}
	params, err = mapClientConfig(config)
	if err != nil || params.JWTToken != "jwt" || params.OrganizationID != "org" {
		t.Fatal("bearer authentication mapping failed")
	}
}
func TestUnsupportedClientSettingsRejectedBeforeSDKInitialization(t *testing.T) {
	for _, config := range []*sandbox.Config{
		{Auth: &sandbox.AuthConfig{Credentials: &sandbox.TokenPair{TokenPair: &sandbox.TokenPairCredentials{Id: "id", Secret: "secret"}}}},
		{Scope: &sandbox.Scope{AppName: sandbox.Value("app")}},
		{Scope: &sandbox.Scope{Environment: sandbox.Value("env")}},
	} {
		if provider, err := New().NewClient(config); err == nil || provider != nil {
			t.Fatal("unsupported client setting initialized SDK")
		}
	}
}
func TestUnsupportedScopeFieldsUseCanonicalPaths(t *testing.T) {
	for _, tc := range []struct {
		name  string
		scope *sandbox.Scope
		field string
	}{
		{"app", &sandbox.Scope{AppName: sandbox.Value("app")}, sandbox.ConfigFieldScopeAppName},
		{"environment", &sandbox.Scope{Environment: sandbox.Value("env")}, sandbox.ConfigFieldScopeEnvironment},
		{"project", &sandbox.Scope{ProjectId: sandbox.Value("project")}, sandbox.ConfigFieldScopeProjectId},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mapClientConfig(&sandbox.Config{Scope: tc.scope})
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != tc.field {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
}

func TestProviderInitializesAndOwnsSDKWithoutRequests(t *testing.T) {
	// The SDK constructor creates local clients; authentication is not probed here.
	t.Setenv("DAYTONA_JWT_TOKEN", "")
	t.Setenv("DAYTONA_API_KEY", "")
	client, err := sandbox.NewClient(New(), &sandbox.Config{
		Auth:     &sandbox.AuthConfig{Credentials: &sandbox.APIKey{ApiKey: &sandbox.APIKeyCredentials{Key: "test-only"}}},
		Endpoint: sandbox.Value("https://example.invalid"),
		Region:   sandbox.Value("region"),
	})
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

func TestEndpointPassThroughAndDefault(t *testing.T) {
	for _, tc := range []struct {
		name     string
		explicit *string
		want     string
	}{
		{name: "default", want: "https://app.daytona.io/api"},
		{name: "explicit", explicit: sandbox.Value("https://explicit.example.test"), want: "https://explicit.example.test"},
		{name: "loopback", explicit: sandbox.Value("http://127.0.0.1:1234"), want: "http://127.0.0.1:1234"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("DAYTONA_API_URL", "https://should-not-use.example/api")
			t.Setenv("DAYTONA_SERVER_URL", "https://should-not-use-either.example/api")
			config := &sandbox.Config{Endpoint: tc.explicit}
			params, err := mapClientConfig(config)
			if err != nil {
				t.Fatal(err)
			}
			if params.APIUrl != tc.want {
				t.Fatalf("got %q want %q", params.APIUrl, tc.want)
			}
			if config.Endpoint != tc.explicit {
				t.Fatal("caller config was mutated")
			}
		})
	}
}

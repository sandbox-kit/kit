package blaxel

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestCommonInitializationMappings(t *testing.T) {
	config := &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
			ApiKey: &sandbox.APIKeyCredentials{Key: "api-key"},
		}},
		Scope:   &sandbox.Scope{Workspace: sandbox.Value("workspace")},
		Region:  sandbox.Value("us-pdx-1"),
		Timeout: sandbox.Duration(time.Second),
	}
	opts, err := mapClientConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if len(opts) != 2 {
		t.Fatalf("expected api key and workspace options, got %d", len(opts))
	}
}

func TestUnsupportedClientSettingsRejectedBeforeSDKInitialization(t *testing.T) {
	for _, config := range []*sandbox.Config{
		{Auth: &sandbox.AuthConfig{Credentials: &sandbox.TokenPair{TokenPair: &sandbox.TokenPairCredentials{Id: "id", Secret: "secret"}}}},
		{Scope: &sandbox.Scope{AppName: sandbox.Value("app")}},
		{Endpoint: sandbox.Value("https://example.invalid")},
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
		{"organization", &sandbox.Scope{OrganizationId: sandbox.Value("org")}, sandbox.ConfigFieldScopeOrganizationId},
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
	t.Setenv("BL_API_KEY", "")
	t.Setenv("BL_WORKSPACE", "")
	client, err := sandbox.NewClient(New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
			ApiKey: &sandbox.APIKeyCredentials{Key: "test-only"},
		}},
		Scope:  &sandbox.Scope{Workspace: sandbox.Value("test-workspace")},
		Region: sandbox.Value("us-pdx-1"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.ProviderName() != "blaxel" {
		t.Fatal("wrong provider")
	}
	if err := client.Close(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestAuthReportsTheMissingCredential(t *testing.T) {
	cases := []struct {
		name   string
		config *sandbox.Config
		field  string
	}{
		{"missing key", &sandbox.Config{Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{}}}, sandbox.ConfigFieldAuthApiKey},
		{"blank key", &sandbox.Config{Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{ApiKey: &sandbox.APIKeyCredentials{}}}}, sandbox.ConfigFieldAuthApiKeyKey},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := mapClientConfig(tc.config)
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != tc.field || detail.Provider != "blaxel" {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
}

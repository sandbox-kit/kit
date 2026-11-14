package sandboxv1_test

import (
	"encoding/json"
	"strings"
	"testing"

	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/internal/contracts/kit/sandbox/v1"
)

func TestCredentialJSONRedactsSensitiveFields(t *testing.T) {
	tests := []struct {
		name        string
		credentials []any
		want        string
	}{
		{name: "api key", credentials: []any{&sandboxv1.APIKeyCredentials{Key: "api-key-secret"}}, want: `{}`},
		{name: "token pair", credentials: []any{&sandboxv1.TokenPairCredentials{Id: "token-id-secret", Secret: "token-secret"}}, want: `{}`},
		{name: "bearer token", credentials: []any{&sandboxv1.BearerTokenCredentials{Token: "bearer-secret"}}, want: `{}`},
		{name: "OAuth", credentials: []any{&sandboxv1.OAuthCredentials{RefreshToken: "refresh-secret", ClientId: "public-client", ClientSecret: sandboxv1String("client-secret"), JwtKey: sandboxv1String("jwt-secret")}}, want: `{"client_id":"public-client"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, credential := range test.credentials {
				encoded, err := json.Marshal(credential)
				if err != nil {
					t.Fatal(err)
				}
				if string(encoded) != test.want {
					t.Fatalf("JSON = %s, want %s", encoded, test.want)
				}
				for _, secret := range []string{"api-key-secret", "token-id-secret", "token-secret", "bearer-secret", "refresh-secret", "client-secret", "jwt-secret"} {
					if strings.Contains(string(encoded), secret) {
						t.Fatalf("JSON contains sensitive value %q", secret)
					}
				}
			}
		})
	}
}

func TestConfigJSONUsesCredentialRedaction(t *testing.T) {
	config := &sandboxv1.Config{Auth: &sandboxv1.AuthConfig{Credentials: &sandboxv1.AuthConfig_ApiKey{ApiKey: &sandboxv1.APIKeyCredentials{Key: "nested-secret"}}}}
	encoded, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "nested-secret") {
		t.Fatalf("config JSON contains credential: %s", encoded)
	}
}

func sandboxv1String(value string) *string { return &value }

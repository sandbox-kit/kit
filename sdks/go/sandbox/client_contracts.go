package sandbox

import (
	time "time"

	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/internal/contracts/kit/sandbox/v1"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

type (
	// Config is the configuration passed to NewClient.
	// Set Auth for the provider account. Modal also requires Scope.AppName.
	// An omitted Endpoint uses the provider default. Modal rejects Endpoint.
	Config = sandboxv1.Config

	// AuthConfig is the credential set on Config.Auth.
	// Set Credentials to one of APIKey, TokenPair, BearerToken, or OAuthRefresh.
	AuthConfig = sandboxv1.AuthConfig

	// APIKey is the API-key form of AuthConfig.Credentials.
	// Daytona accepts an API key.
	APIKey = sandboxv1.AuthConfig_ApiKey

	// TokenPair is the id-and-secret form of AuthConfig.Credentials.
	// Modal accepts a token pair or OAuth refresh. Omitting Auth lets the Modal SDK use credentials it finds in the environment. An API key is rejected.
	TokenPair = sandboxv1.AuthConfig_TokenPair

	// BearerToken is the bearer-token form of AuthConfig.Credentials.
	BearerToken = sandboxv1.AuthConfig_BearerToken

	// OAuthRefresh is the refresh-token form of AuthConfig.Credentials.
	OAuthRefresh = sandboxv1.AuthConfig_OauthRefresh

	// APIKeyCredentials holds the key for an APIKey credential.
	APIKeyCredentials = sandboxv1.APIKeyCredentials

	// TokenPairCredentials holds the id and secret for a TokenPair credential.
	TokenPairCredentials = sandboxv1.TokenPairCredentials

	// BearerTokenCredentials holds the token for a BearerToken credential.
	BearerTokenCredentials = sandboxv1.BearerTokenCredentials

	// OAuthCredentials holds the refresh token and client credentials for OAuthRefresh.
	OAuthCredentials = sandboxv1.OAuthCredentials

	// Scope selects the account boundary for the client.
	// Modal requires AppName. Environment is often "main".
	// Daytona accepts OrganizationId.
	Scope = sandboxv1.Scope
)

// Duration converts a Go duration into the timeout value used by Config and lifecycle options.
func Duration(value time.Duration) *durationpb.Duration {
	return durationpb.New(value)
}

// ValidateConfig checks client configuration before a provider is contacted.
// A nil config is valid. A remote endpoint must use HTTPS. Localhost and
// loopback may use HTTP.
func ValidateConfig(config *Config) error {
	return validateProtoConfig(config)
}

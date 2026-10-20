package sandbox

import (
	time "time"

	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Client configuration types are standard protobuf-generated contracts. These
// aliases preserve the SDK's concise package entrypoint while ensuring there is
// exactly one representation of each portable value.
type (
	Config                  = sandboxv1.Config
	AuthConfig              = sandboxv1.AuthConfig
	AuthConfig_ApiKey       = sandboxv1.AuthConfig_ApiKey
	AuthConfig_TokenPair    = sandboxv1.AuthConfig_TokenPair
	AuthConfig_BearerToken  = sandboxv1.AuthConfig_BearerToken
	AuthConfig_OauthRefresh = sandboxv1.AuthConfig_OauthRefresh
	APIKeyCredentials       = sandboxv1.APIKeyCredentials
	TokenPairCredentials    = sandboxv1.TokenPairCredentials
	BearerTokenCredentials  = sandboxv1.BearerTokenCredentials
	OAuthCredentials        = sandboxv1.OAuthCredentials
	Scope                   = sandboxv1.Scope
)

// Duration builds a protobuf duration for client and operation configuration.
func Duration(value time.Duration) *durationpb.Duration {
	return durationpb.New(value)
}

// ValidateConfig validates generated client configuration before provider
// initialization.
func ValidateConfig(config *Config) error {
	return validateProtoConfig(config)
}

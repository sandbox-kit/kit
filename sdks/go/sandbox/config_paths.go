// Config field-path constants identify client settings in errors. Values are canonical protobuf paths.
//
// Example:
//
//	fmt.Println(ConfigFieldEndpoint) // endpoint
package sandbox

// ConfigFieldAuth is the canonical path "auth".
const ConfigFieldAuth = "auth"

// ConfigFieldAuthApiKey is the canonical path "auth.api_key".
const ConfigFieldAuthApiKey = "auth.api_key"

// ConfigFieldAuthApiKeyKey is the canonical path "auth.api_key.key".
const ConfigFieldAuthApiKeyKey = "auth.api_key.key"

// ConfigFieldAuthTokenPair is the canonical path "auth.token_pair".
const ConfigFieldAuthTokenPair = "auth.token_pair"

// ConfigFieldAuthTokenPairId is the canonical path "auth.token_pair.id".
const ConfigFieldAuthTokenPairId = "auth.token_pair.id"

// ConfigFieldAuthTokenPairSecret is the canonical path
// "auth.token_pair.secret".
const ConfigFieldAuthTokenPairSecret = "auth.token_pair.secret"

// ConfigFieldAuthBearerToken is the canonical path "auth.bearer_token".
const ConfigFieldAuthBearerToken = "auth.bearer_token"

// ConfigFieldAuthBearerTokenToken is the canonical path
// "auth.bearer_token.token".
const ConfigFieldAuthBearerTokenToken = "auth.bearer_token.token"

// ConfigFieldAuthOauthRefresh is the canonical path "auth.oauth_refresh".
const ConfigFieldAuthOauthRefresh = "auth.oauth_refresh"

// ConfigFieldAuthOauthRefreshRefreshToken is the canonical path
// "auth.oauth_refresh.refresh_token".
const ConfigFieldAuthOauthRefreshRefreshToken = "auth.oauth_refresh.refresh_token"

// ConfigFieldAuthOauthRefreshClientId is the canonical path
// "auth.oauth_refresh.client_id".
const ConfigFieldAuthOauthRefreshClientId = "auth.oauth_refresh.client_id"

// ConfigFieldAuthOauthRefreshClientSecret is the canonical path
// "auth.oauth_refresh.client_secret".
const ConfigFieldAuthOauthRefreshClientSecret = "auth.oauth_refresh.client_secret"

// ConfigFieldAuthOauthRefreshJwtKey is the canonical path
// "auth.oauth_refresh.jwt_key".
const ConfigFieldAuthOauthRefreshJwtKey = "auth.oauth_refresh.jwt_key"

// ConfigFieldEndpoint is the canonical path "endpoint".
const ConfigFieldEndpoint = "endpoint"

// ConfigFieldRegion is the canonical path "region".
const ConfigFieldRegion = "region"

// ConfigFieldScope is the canonical path "scope".
const ConfigFieldScope = "scope"

// ConfigFieldScopeAppName is the canonical path "scope.app_name".
const ConfigFieldScopeAppName = "scope.app_name"

// ConfigFieldScopeEnvironment is the canonical path "scope.environment".
const ConfigFieldScopeEnvironment = "scope.environment"

// ConfigFieldScopeOrganizationId is the canonical path
// "scope.organization_id".
const ConfigFieldScopeOrganizationId = "scope.organization_id"

// ConfigFieldScopeProjectId is the canonical path "scope.project_id".
const ConfigFieldScopeProjectId = "scope.project_id"

// ConfigFieldTimeout is the canonical path "timeout".
const ConfigFieldTimeout = "timeout"

// ConfigFieldProvider is the canonical path "provider".
const ConfigFieldProvider = "provider"

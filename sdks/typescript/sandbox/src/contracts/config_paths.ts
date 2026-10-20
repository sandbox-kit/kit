// Config field-path constants identify client settings in errors.
// Values are canonical protobuf paths (snake_case), matching Go.
//
// Example:
//
//   console.log(ConfigFieldEndpoint); // endpoint

/** Canonical path "auth". */
export const ConfigFieldAuth = "auth" as const;

/** Canonical path "auth.api_key". */
export const ConfigFieldAuthApiKey = "auth.api_key" as const;

/** Canonical path "auth.api_key.key". */
export const ConfigFieldAuthApiKeyKey = "auth.api_key.key" as const;

/** Canonical path "auth.token_pair". */
export const ConfigFieldAuthTokenPair = "auth.token_pair" as const;

/** Canonical path "auth.token_pair.id". */
export const ConfigFieldAuthTokenPairId = "auth.token_pair.id" as const;

/** Canonical path "auth.token_pair.secret". */
export const ConfigFieldAuthTokenPairSecret = "auth.token_pair.secret" as const;

/** Canonical path "auth.bearer_token". */
export const ConfigFieldAuthBearerToken = "auth.bearer_token" as const;

/** Canonical path "auth.bearer_token.token". */
export const ConfigFieldAuthBearerTokenToken = "auth.bearer_token.token" as const;

/** Canonical path "auth.oauth_refresh". */
export const ConfigFieldAuthOauthRefresh = "auth.oauth_refresh" as const;

/** Canonical path "auth.oauth_refresh.refresh_token". */
export const ConfigFieldAuthOauthRefreshRefreshToken = "auth.oauth_refresh.refresh_token" as const;

/** Canonical path "auth.oauth_refresh.client_id". */
export const ConfigFieldAuthOauthRefreshClientId = "auth.oauth_refresh.client_id" as const;

/** Canonical path "auth.oauth_refresh.client_secret". */
export const ConfigFieldAuthOauthRefreshClientSecret = "auth.oauth_refresh.client_secret" as const;

/** Canonical path "auth.oauth_refresh.jwt_key". */
export const ConfigFieldAuthOauthRefreshJwtKey = "auth.oauth_refresh.jwt_key" as const;

/** Canonical path "endpoint". */
export const ConfigFieldEndpoint = "endpoint" as const;

/** Canonical path "region". */
export const ConfigFieldRegion = "region" as const;

/** Canonical path "scope". */
export const ConfigFieldScope = "scope" as const;

/** Canonical path "scope.app_name". */
export const ConfigFieldScopeAppName = "scope.app_name" as const;

/** Canonical path "scope.environment". */
export const ConfigFieldScopeEnvironment = "scope.environment" as const;

/** Canonical path "scope.organization_id". */
export const ConfigFieldScopeOrganizationId = "scope.organization_id" as const;

/** Canonical path "scope.project_id". */
export const ConfigFieldScopeProjectId = "scope.project_id" as const;

/** Canonical path "timeout". */
export const ConfigFieldTimeout = "timeout" as const;

/** Canonical path "provider". */
export const ConfigFieldProvider = "provider" as const;

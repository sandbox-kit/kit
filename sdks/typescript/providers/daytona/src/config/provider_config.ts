// Client configuration composition for the native SDK constructor.
// Kept in sync with specs/providers/daytona.yaml client: rules (not emitted).

import type { DaytonaConfig as NativeClientConfig } from "@daytona/sdk";
import {
  ErrorKind,
  sandboxError,
  validateConfig,
  type Config,
  ConfigFieldAuth,
  ConfigFieldScopeAppName,
  ConfigFieldScopeEnvironment,
  ConfigFieldScopeProjectId
} from "@sandbox-kit/sandbox";
import {
  mapAPIKeyCredentials,
  mapBearerAuth,
  mapClientSettings,
  mapScope,
} from "#mappings/provider_mappings.js";

export type Environment = Readonly<Record<string, string | undefined>>;

export function mapGeneratedClientConfig(
  config: Config,
  environment: Environment = environmentVariables(),
): NativeClientConfig {
  let source = config;
  if (source.endpoint === undefined)
    source = {
      ...source,
      endpoint:
        environment.DAYTONA_API_URL ??
        environment.DAYTONA_SERVER_URL ??
        "https://app.daytona.io/api",
    };
  validateConfig(source, "daytona", "initialize");
  const target: NativeClientConfig = mapClientSettings(source, "initialize");
  if (source.scope !== undefined) {
    const mapped = mapScope(source.scope, "initialize");
    Object.assign(target, mapped);
    if (source.scope.appName !== undefined)
      unsupported(ConfigFieldScopeAppName, "unsupported client configuration field");
    if (source.scope.environment !== undefined)
      unsupported(
        ConfigFieldScopeEnvironment,
        "unsupported client configuration field",
      );
    if (source.scope.projectId !== undefined)
      unsupported(ConfigFieldScopeProjectId, "unsupported client configuration field");
  }
  if (source.auth?.apiKey !== undefined) {
    Object.assign(
      target,
      mapAPIKeyCredentials(source.auth.apiKey, "initialize"),
    );
  } else if (source.auth?.bearerToken !== undefined) {
    Object.assign(
      target,
      mapBearerAuth(source.auth.bearerToken, "initialize"),
    );
  } else if (source.auth !== undefined) {
    unsupported(ConfigFieldAuth, "unsupported authentication mode");
  }
  return target;
}

function environmentVariables(): Environment {
  return (
    (globalThis as { readonly process?: { readonly env?: Environment } })
      .process?.env ?? {}
  );
}

function unsupported(field: string, message: string): never {
  throw sandboxError({
    kind: ErrorKind.Unsupported,
    provider: "daytona",
    operation: "initialize",
    field,
    message: "sandbox-kit daytona: " + message,
  });
}

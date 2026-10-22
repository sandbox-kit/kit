// Client configuration composition for the native SDK constructor.
// Kept in sync with specs/providers/daytona.yaml client: rules (not emitted).

import type { DaytonaConfig as NativeClientConfig } from "@daytona/sdk";
import {
  ErrorKind,
  type Config,
  ConfigFieldAuth,
  ConfigFieldScopeAppName,
  ConfigFieldScopeEnvironment,
  ConfigFieldScopeProjectId,
} from "@sandbox-kit/sandbox";
import {
  mapAPIKeyCredentials,
  mapBearerAuth,
  mapClientSettings,
  mapScope,
} from "#mappings/provider_mappings.js";
import { providerError } from "#errors.js";

const defaultEndpoint = "https://app.daytona.io/api";

export function mapGeneratedClientConfig(config: Config): NativeClientConfig {
  const source: Config =
    config.endpoint === undefined
      ? { ...config, endpoint: defaultEndpoint }
      : config;
  const target: NativeClientConfig = mapClientSettings(source, "initialize");
  if (source.scope !== undefined) {
    const mapped = mapScope(source.scope, "initialize");
    Object.assign(target, mapped);
    if (source.scope.appName !== undefined)
      unsupported(
        ConfigFieldScopeAppName,
        "unsupported client configuration field",
      );
    if (source.scope.environment !== undefined)
      unsupported(
        ConfigFieldScopeEnvironment,
        "unsupported client configuration field",
      );
    if (source.scope.projectId !== undefined)
      unsupported(
        ConfigFieldScopeProjectId,
        "unsupported client configuration field",
      );
  }
  if (source.auth?.apiKey !== undefined) {
    Object.assign(
      target,
      mapAPIKeyCredentials(source.auth.apiKey, "initialize"),
    );
  } else if (source.auth?.bearerToken !== undefined) {
    Object.assign(target, mapBearerAuth(source.auth.bearerToken, "initialize"));
  } else if (source.auth !== undefined) {
    unsupported(ConfigFieldAuth, "unsupported authentication mode");
  }
  return target;
}

function unsupported(field: string, message: string): never {
  providerError({
    kind: ErrorKind.Unsupported,
    operation: "initialize",
    field,
    message,
  });
}

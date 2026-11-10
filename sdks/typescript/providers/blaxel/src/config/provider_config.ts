// Client configuration composition for the native Blaxel settings object.
// Kept in sync with specs/providers/blaxel.yaml client: rules when present.

import type { Config as NativeClientConfig } from "@blaxel/core";
import {
  ErrorKind,
  type Config,
  type Scope,
  ConfigFieldAuth,
  ConfigFieldEndpoint,
  ConfigFieldScopeAppName,
  ConfigFieldScopeEnvironment,
  ConfigFieldScopeOrganizationId,
  ConfigFieldScopeProjectId,
} from "@sandbox-kit/sandbox";
import {
  mapAPIKeyCredentials,
  mapClientScope,
} from "#mappings/provider_mappings.js";
import { providerError } from "#errors.js";

export function mapGeneratedClientConfig(config: Config): NativeClientConfig {
  const source = config;
  const target: NativeClientConfig = {};
  if (source.endpoint !== undefined)
    unsupported(
      ConfigFieldEndpoint,
      "Endpoint cannot be configured through this SDK's public client constructor",
    );
  if (source.scope !== undefined) {
    Object.assign(target, mapClientScope(source.scope, "initialize"));
    for (const [set, field] of [
      [source.scope.appName !== undefined, ConfigFieldScopeAppName],
      [source.scope.environment !== undefined, ConfigFieldScopeEnvironment],
      [
        source.scope.organizationId !== undefined,
        ConfigFieldScopeOrganizationId,
      ],
      [source.scope.projectId !== undefined, ConfigFieldScopeProjectId],
    ] as const) {
      if (set) unsupported(field, "unsupported client configuration field");
    }
  }
  if (source.auth?.apiKey !== undefined) {
    Object.assign(
      target,
      mapAPIKeyCredentials(source.auth.apiKey, "initialize"),
    );
  } else if (source.auth !== undefined) {
    unsupported(ConfigFieldAuth, "unsupported authentication mode");
  }
  return target;
}

export interface ProviderState {
  readonly scope?: Scope;
  readonly region?: string;
}

export function captureProviderState(config: Config): ProviderState {
  return {
    ...(config.scope === undefined
      ? {}
      : { scope: structuredClone(config.scope) }),
    ...(config.region === undefined ? {} : { region: config.region }),
  };
}

function unsupported(field: string, message: string): never {
  providerError({
    kind: ErrorKind.Unsupported,
    operation: "initialize",
    field,
    message,
  });
}

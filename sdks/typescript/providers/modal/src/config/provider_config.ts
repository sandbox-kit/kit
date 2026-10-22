// Client configuration composition for the native SDK constructor.
// Kept in sync with specs/providers/modal.yaml client: rules (not emitted).

import type { ModalClientParams as NativeClientConfig } from "modal";
import {
  ErrorKind,
  type Config,
  type Scope,
  ConfigFieldAuth,
  ConfigFieldEndpoint,
  ConfigFieldScopeProjectId,
} from "@sandbox-kit/sandbox";
import {
  mapClientScope,
  mapOAuthAuth,
  mapTokenPairCredentials,
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
    const mapped = mapClientScope(source.scope, "initialize");
    Object.assign(target, mapped);
    if (source.scope.organizationId !== undefined)
      unsupported(
        "scope.organizationId",
        "unsupported client configuration field",
      );
    if (source.scope.projectId !== undefined)
      unsupported(
        ConfigFieldScopeProjectId,
        "unsupported client configuration field",
      );
  }
  if (source.auth?.tokenPair !== undefined) {
    Object.assign(
      target,
      mapTokenPairCredentials(source.auth.tokenPair, "initialize"),
    );
  } else if (source.auth?.oauthRefresh !== undefined) {
    Object.assign(target, mapOAuthAuth(source.auth.oauthRefresh, "initialize"));
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

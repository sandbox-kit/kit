// Client configuration composition for the native SDK constructor.
// Kept in sync with specs/providers/modal.yaml client: rules (not emitted).

import type { ModalClientParams as NativeClientConfig } from "modal";
import {
  ErrorKind,
  resolveClientConfig,
  sandboxError,
  validateConfig,
  type Config,
  type ConfigInput,
  type Scope,
} from "@sandbox-kit/sandbox";
import {
  mapClientScope,
  mapOAuthAuth,
  mapTokenPairCredentials,
} from "./provider_mappings.js";

export function mapGeneratedClientConfig(config: Config): NativeClientConfig {
  let source = config;
  validateConfig(source, "modal", "initialize");
  const target: NativeClientConfig = {};
  if (source.endpoint !== undefined)
    unsupported(
      "endpoint",
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
      unsupported("scope.projectId", "unsupported client configuration field");
  }
  const credentials = source.auth?.credentials;
  if (credentials?.case === "tokenPair") {
    const mapped = mapTokenPairCredentials(credentials.value, "initialize");
    Object.assign(target, mapped);
  } else if (credentials?.case === "oauthRefresh") {
    const mapped = mapOAuthAuth(credentials.value, "initialize");
    Object.assign(target, mapped);
  } else if (credentials?.case !== undefined) {
    unsupported("auth", "unsupported authentication mode");
  }
  return target;
}

export interface ProviderState {
  readonly scope?: ConfigInput["scope"] | Scope;
  readonly region?: string;
}

export function captureProviderState(
  config: Config | ConfigInput,
): ProviderState {
  const source = resolveClientConfig(config);
  return {
    ...(source.scope === undefined
      ? {}
      : { scope: structuredClone(source.scope) }),
    ...(source.region === undefined ? {} : { region: source.region }),
  };
}

function unsupported(field: string, message: string): never {
  throw sandboxError({
    kind: ErrorKind.Unsupported,
    provider: "modal",
    operation: "initialize",
    field,
    message: "sandbox-kit modal: " + message,
  });
}

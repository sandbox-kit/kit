// Provider field mappings between shared CreateOptions/Config and the native SDK.
// Inventoried in specs/providers/modal.yaml groups; kept in sync by hand (not emitted).

import { durationMs } from "@bufbuild/protobuf/wkt";
import { ErrorKind, PolicyMode, ValueOrigin, sandboxError, type AutomaticAction, type CreateOptions, type LifetimePolicy, type OAuthCredentials, type Resources, type SandboxInfo, type Scope, type TokenPairCredentials } from "@sandbox-kit/sandbox";
import type { ModalClientParams as NativeModalClientParams, Sandbox as NativeSandbox, SandboxCreateParams as NativeSandboxCreateParams } from "modal";

type Mutable<T> = { -readonly [K in keyof T]: T[K] };

export function mapSandboxIdentity(source: NativeSandbox | undefined, operation: string = "create"): Partial<SandboxInfo> {
  const target: Mutable<Partial<SandboxInfo>> = {};
  if (source === undefined) return target;
  const origins: Record<string, ValueOrigin> = {};
  target.origins = origins;
  {
    if (source.sandboxId !== undefined) {
      const value = source.sandboxId;
      target.id = value;
      origins["id"] = ValueOrigin.Provider;
    }
  }
  target.provider = "modal";
  return target;
}

export function mapClientScope(source: Scope | undefined, operation: string = "initialize"): NativeModalClientParams {
  const target: Mutable<NativeModalClientParams> = {};
  if (source === undefined) return target;
  {
    if (source.environment !== undefined) {
      const value = source.environment;
      target.environment = value;
    }
  }
  return target;
}

export function mapCommonCreateFields(source: CreateOptions | undefined, operation: string = "create"): NativeSandboxCreateParams {
  const target: Mutable<NativeSandboxCreateParams> = {};
  if (source === undefined) return target;
  {
    if (source.name !== undefined) {
      const value = source.name;
      target.name = value;
    }
  }
  {
    if (source.environment !== undefined) {
      const value = source.environment;
      target.env = { ...value };
    }
  }
  {
    if (source.labels !== undefined) {
      const value = source.labels;
      target.tags = { ...value };
    }
  }
  return target;
}

export function mapResources(source: Resources | undefined, operation: string = "create"): NativeSandboxCreateParams {
  const target: Mutable<NativeSandboxCreateParams> = {};
  if (source === undefined) return target;
  {
    if (source.cpuCores !== undefined) {
      const value = source.cpuCores;
      target.cpu = value;
    }
  }
  {
    if (source.cpuLimitCores !== undefined) {
      const value = source.cpuLimitCores;
      target.cpuLimit = value;
    }
  }
  {
    if (source.memoryMib !== undefined) {
      const value = source.memoryMib;
      if (value > 2147483647n) mappingError(ErrorKind.InvalidArgument, "resources.memory_mib", "memory_mib cannot be represented", operation);
      target.memoryMiB = Number(value);
    }
  }
  {
    if (source.memoryLimitMib !== undefined) {
      const value = source.memoryLimitMib;
      if (value > 2147483647n) mappingError(ErrorKind.InvalidArgument, "resources.memory_limit_mib", "memory_limit_mib cannot be represented", operation);
      target.memoryLimitMiB = Number(value);
    }
  }
  return target;
}

export function mapLifetimeDuration(source: LifetimePolicy | undefined, operation: string = "create"): NativeSandboxCreateParams {
  const target: Mutable<NativeSandboxCreateParams> = {};
  if (source === undefined) return target;
  {
    if (source.maximumLifetime !== undefined) {
      const value = durationMs(source.maximumLifetime);
      if (value > 86400000) mappingError(ErrorKind.InvalidArgument, "lifetime.maximum_lifetime", "maximum_lifetime cannot be represented", operation);
      target.timeoutMs = value;
    }
  }
  return target;
}

export function mapIdleDuration(source: AutomaticAction | undefined, operation: string = "create"): NativeSandboxCreateParams {
  const target: Mutable<NativeSandboxCreateParams> = {};
  if (source === undefined) return target;
  {
    if (source.after !== undefined) {
      const value = durationMs(source.after);
      if (value > 4294967295000) mappingError(ErrorKind.InvalidArgument, "lifetime.idle_terminate.after", "after cannot be represented", operation);
      target.idleTimeoutMs = value;
    }
  }
  return target;
}

export function mapTokenPairCredentials(source: TokenPairCredentials | undefined, operation: string = "initialize"): NativeModalClientParams {
  const target: Mutable<NativeModalClientParams> = {};
  if (source === undefined) return target;
  {
    const value = source.id;
    target.tokenId = value;
  }
  {
    const value = source.secret;
    target.tokenSecret = value;
  }
  return target;
}

export function mapOAuthAuth(source: OAuthCredentials | undefined, operation: string = "create"): NativeModalClientParams {
  const target: Mutable<NativeModalClientParams> = {};
  if (source === undefined) return target;
  {
    const value = source.refreshToken;
    target.oauthRefreshToken = value;
  }
  {
    const value = source.clientId;
    target.oauthClientId = value;
  }
  {
    if (source.clientSecret !== undefined) {
      const value = source.clientSecret;
      target.oauthClientSecret = value;
    }
  }
  {
    if (source.jwtKey !== undefined) {
      const value = source.jwtKey;
      target.oauthJwtKey = value;
    }
  }
  return target;
}

function mappingError(kind: ErrorKind, field: string, message: string, operation: string): never {
  throw sandboxError({ kind, provider: "modal", operation, field, message: "sandbox-kit modal: " + message });
}

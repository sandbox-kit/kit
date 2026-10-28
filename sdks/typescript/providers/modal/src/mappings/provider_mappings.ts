// Provider field mappings between shared CreateOptions/Config and the native SDK.
// Inventoried in specs/providers/modal.yaml groups; kept in sync by hand (not emitted).

import {
  ErrorKind,
  type SandboxInfo,
  type Scope,
  CreateFieldResourcesMemoryLimitMib,
  CreateFieldResourcesMemoryMib,
  CreateFieldLifetimeIdleTerminateAfter,
  CreateFieldLifetimeMaximumLifetime,
  CreateFieldResourcesCpuCores,
  CreateFieldResourcesCpuLimitCores,
} from "@sandbox-kit/sandbox";
import type {
  CreateOptions,
  OAuthCredentials,
  Resources,
  TokenPairCredentials,
} from "@sandbox-kit/sandbox";
import type {
  ModalClientParams as NativeModalClientParams,
  Sandbox as NativeSandbox,
  SandboxCreateParams as NativeSandboxCreateParams,
} from "modal";
import { providerError, type ProviderErrorInput } from "#errors.js";

type Mutable<T> = { -readonly [K in keyof T]: T[K] };
type PublicLifetime = NonNullable<CreateOptions["lifetime"]>;
type PublicIdleTerminate = NonNullable<PublicLifetime["idleTerminate"]>;

export function mapSandboxIdentity(
  source: NativeSandbox | undefined,
  operation: string = "create",
): Partial<SandboxInfo> {
  const target: Mutable<Partial<SandboxInfo>> = {};
  if (source === undefined) return target;
  {
    if (source.sandboxId !== undefined) {
      const value = source.sandboxId;
      target.id = value;
    }
  }
  target.provider = "modal";
  return target;
}

export function mapClientScope(
  source: Scope | undefined,
  operation: string = "initialize",
): NativeModalClientParams {
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

export function mapCommonCreateFields(
  source: CreateOptions | undefined,
  operation: string = "create",
): NativeSandboxCreateParams {
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

export function mapResources(
  source: Resources | undefined,
  operation: string = "create",
): NativeSandboxCreateParams {
  const target: Mutable<NativeSandboxCreateParams> = {};
  if (source === undefined) return target;
  {
    if (source.cpuCores !== undefined) {
      target.cpu = scaledCpuCores(
        source.cpuCores,
        CreateFieldResourcesCpuCores,
        operation,
      );
    }
  }
  {
    if (source.cpuLimitCores !== undefined) {
      target.cpuLimit = scaledCpuCores(
        source.cpuLimitCores,
        CreateFieldResourcesCpuLimitCores,
        operation,
      );
    }
  }
  {
    if (source.memoryMib !== undefined) {
      const value = source.memoryMib;
      if (value > 2147483647n)
        mappingError({
          kind: ErrorKind.InvalidArgument,
          field: CreateFieldResourcesMemoryMib,
          message: "memory_mib cannot be represented",
          operation,
        });
      target.memoryMiB = Number(value);
    }
  }
  {
    if (source.memoryLimitMib !== undefined) {
      const value = source.memoryLimitMib;
      if (value > 2147483647n)
        mappingError({
          kind: ErrorKind.InvalidArgument,
          field: CreateFieldResourcesMemoryLimitMib,
          message: "memory_limit_mib cannot be represented",
          operation,
        });
      target.memoryLimitMiB = Number(value);
    }
  }
  return target;
}

export function mapLifetimeDuration(
  source: PublicLifetime | undefined,
  operation: string = "create",
): NativeSandboxCreateParams {
  const target: Mutable<NativeSandboxCreateParams> = {};
  if (source === undefined) return target;
  {
    if (source.maximumLifetime !== undefined) {
      target.timeoutMs = durationSecondsMs(
        source.maximumLifetime,
        CreateFieldLifetimeMaximumLifetime,
        86_400,
        operation,
      );
    }
  }
  return target;
}

export function mapIdleDuration(
  source: PublicIdleTerminate | undefined,
  operation: string = "create",
): NativeSandboxCreateParams {
  const target: Mutable<NativeSandboxCreateParams> = {};
  if (source === undefined) return target;
  {
    if (source.after !== undefined) {
      target.idleTimeoutMs = durationSecondsMs(
        source.after,
        CreateFieldLifetimeIdleTerminateAfter,
        4_294_967_295,
        operation,
      );
    }
  }
  return target;
}

/** Modal CPU fields are millis-scaled: value must equal round(value*1000)/1000. */
function scaledCpuCores(
  value: number,
  field: string,
  operation: string,
): number {
  const scaled = Math.round(value * 1000);
  if (
    !Number.isFinite(value) ||
    value < 0 ||
    value !== scaled / 1000 ||
    scaled > 4_294_967_295
  ) {
    mappingError({
      kind: ErrorKind.InvalidArgument,
      field,
      message:
        field.slice(field.lastIndexOf(".") + 1) + " cannot be represented",
      operation,
    });
  }
  return value * 1000 < scaled ? scaled / 1000 : value;
}

/** Modal lifetime/idle timeouts are whole seconds expressed as milliseconds. */
function durationSecondsMs(
  value: number,
  field: string,
  maximumSeconds: number,
  operation: string,
): number {
  if (
    !Number.isFinite(value) ||
    value <= 0 ||
    value % 1000 !== 0 ||
    value / 1000 > maximumSeconds
  ) {
    mappingError({
      kind: ErrorKind.InvalidArgument,
      field,
      message:
        field.slice(field.lastIndexOf(".") + 1) + " cannot be represented",
      operation,
    });
  }
  return value;
}

export function mapTokenPairCredentials(
  source: TokenPairCredentials | undefined,
  operation: string = "initialize",
): NativeModalClientParams {
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

export function mapOAuthAuth(
  source: OAuthCredentials | undefined,
  operation: string = "create",
): NativeModalClientParams {
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

function mappingError(input: ProviderErrorInput): never {
  providerError(input);
}

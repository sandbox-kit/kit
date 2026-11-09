// Provider field mappings between shared CreateOptions/Config and the native SDK.
// Inventoried in specs/providers/blaxel.yaml groups when present; kept in sync by hand.

import {
  ErrorKind,
  SandboxState,
  type SandboxInfo,
  type Scope,
  CreateFieldResourcesMemoryMib,
} from "@sandbox-kit/sandbox";
import type {
  APIKeyCredentials,
  CreateOptions,
  Resources,
} from "@sandbox-kit/sandbox";
import type {
  Config as NativeClientConfig,
  SandboxCreateConfiguration,
} from "@blaxel/core";
import { providerError, type ProviderErrorInput } from "#errors.js";

type Mutable<T> = { -readonly [K in keyof T]: T[K] };

/** Minimal native sandbox shape used for identity mapping. */
export interface BlaxelSandboxIdentity {
  readonly metadata?: {
    readonly name?: string;
    readonly labels?: Readonly<Record<string, string>>;
  };
  readonly status?: unknown;
  readonly spec?: { readonly region?: string };
}

export function mapSandboxIdentity(
  source: BlaxelSandboxIdentity | undefined,
  _operation: string = "create",
): Partial<SandboxInfo> {
  const target: Mutable<Partial<SandboxInfo>> = {};
  if (source === undefined) return target;
  const name = source.metadata?.name;
  if (name !== undefined && name !== "") {
    target.id = name;
    target.name = name;
  }
  if (source.status !== undefined) {
    target.providerState = String(source.status);
    const state = mapSandboxState(String(source.status));
    if (state !== undefined) target.state = state;
  }
  const region = source.spec?.region;
  if (region !== undefined && region !== "") target.region = region;
  const labels = source.metadata?.labels;
  if (labels !== undefined) target.labels = { ...labels };
  target.provider = "blaxel";
  return target;
}

export function mapClientScope(
  source: Scope | undefined,
  _operation: string = "initialize",
): NativeClientConfig {
  const target: Mutable<NativeClientConfig> = {};
  if (source === undefined) return target;
  if (source.workspace !== undefined) target.workspace = source.workspace;
  return target;
}

export function mapCommonCreateFields(
  source: CreateOptions | undefined,
  _operation: string = "create",
): SandboxCreateConfiguration {
  const target: Mutable<SandboxCreateConfiguration> = {};
  if (source === undefined) return target;
  if (source.name !== undefined) target.name = source.name;
  if (source.environment !== undefined) {
    target.envs = Object.entries(source.environment).map(([name, value]) => ({
      name,
      value: String(value),
    }));
  }
  if (source.labels !== undefined) target.labels = { ...source.labels };
  return target;
}

export function mapResources(
  source: Resources | undefined,
  operation: string = "create",
): Pick<SandboxCreateConfiguration, "memory"> {
  const target: Mutable<Pick<SandboxCreateConfiguration, "memory">> = {};
  if (source === undefined) return target;
  if (source.memoryMib !== undefined) {
    const value = source.memoryMib;
    if (value > BigInt(Number.MAX_SAFE_INTEGER))
      mappingError({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldResourcesMemoryMib,
        message: "memory_mib cannot be represented",
        operation,
      });
    target.memory = Number(value);
  }
  return target;
}

export function mapAPIKeyCredentials(
  source: APIKeyCredentials | undefined,
  _operation: string = "initialize",
): NativeClientConfig {
  const target: Mutable<NativeClientConfig> = {};
  if (source === undefined) return target;
  target.apiKey = source.key;
  return target;
}

function mapSandboxState(status: string): SandboxState | undefined {
  const normalized = status.toUpperCase();
  switch (normalized) {
    case "RUNNING":
    case "READY":
    case "DEPLOYED":
      return SandboxState.Running;
    case "PENDING":
    case "CREATING":
    case "DEPLOYING":
    case "STARTING":
      return SandboxState.Starting;
    case "STANDBY":
    case "SUSPENDED":
    case "SLEEPING":
      return SandboxState.Paused;
    case "ARCHIVED":
      return SandboxState.Archived;
    case "FAILED":
      return SandboxState.Error;
    case "TERMINATED":
    case "DELETED":
      return SandboxState.Terminated;
    default:
      return undefined;
  }
}

function mappingError(input: ProviderErrorInput): never {
  providerError(input);
}

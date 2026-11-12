// Provider field mappings between shared CreateOptions/Config and the native SDK.
// Inventoried in specs/providers/blaxel.yaml groups when present; kept in sync by hand.

import {
  ConfigFieldScopeWorkspace,
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
import { providerError } from "#errors.js";

type Mutable<T> = { -readonly [K in keyof T]: T[K] };

/** Minimal native sandbox shape used for identity mapping. */
export interface BlaxelSandboxIdentity {
  readonly metadata?: {
    readonly name?: string;
    readonly labels?: Readonly<Record<string, string>>;
  };
  /** Deployment status (DEPLOYED, FAILED, …). */
  readonly status?: unknown;
  readonly spec?: {
    readonly region?: string;
    readonly runtime?: { readonly memory?: number };
  };
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
  if (source.status !== undefined && String(source.status) !== "") {
    target.providerState = String(source.status);
  }
  const state = mapSandboxState(source.status);
  if (state !== undefined) target.state = state;
  const region = source.spec?.region;
  if (region !== undefined && region !== "") target.region = region;
  const labels = source.metadata?.labels;
  if (labels !== undefined) target.labels = { ...labels };
  target.provider = "blaxel";
  return target;
}

export function mapClientScope(
  source: Scope | undefined,
  operation: string = "initialize",
): NativeClientConfig {
  const target: Mutable<NativeClientConfig> = {};
  if (source === undefined) return target;
  if (source.workspace !== undefined) {
    if (source.workspace.trim() === "") {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation,
        field: ConfigFieldScopeWorkspace,
        message: "workspace requires a non-blank value",
      });
    }
    target.workspace = source.workspace;
  }
  return target;
}

export function mapCommonCreateFields(
  source: CreateOptions | undefined,
  _operation: string = "create",
): SandboxCreateConfiguration {
  const target: Mutable<SandboxCreateConfiguration> = {};
  if (source === undefined) return target;
  if (source.name !== undefined) target.name = source.name;
  // Match Go: omit envs when the map is empty so explicit {} equals omit.
  if (source.environment !== undefined) {
    const entries = Object.entries(source.environment);
    if (entries.length > 0) {
      target.envs = entries.map(([name, value]) => ({
        name,
        value: String(value),
      }));
    }
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
    if (value === 0n || value > BigInt(Number.MAX_SAFE_INTEGER))
      providerError({
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

/** Maps native Blaxel memory (MB) onto Kit `resources.memoryMib` 1:1. */
export function mapAllocatedResources(
  source: BlaxelSandboxIdentity | undefined,
): Resources | undefined {
  const memory = source?.spec?.runtime?.memory;
  if (memory === undefined || !Number.isFinite(memory) || memory <= 0) {
    return undefined;
  }
  return { memoryMib: BigInt(Math.trunc(memory)) };
}

/** Maps the public deployment status exposed by @blaxel/core. */
function mapSandboxState(status: unknown): SandboxState | undefined {
  const normalized =
    status === undefined || String(status) === ""
      ? ""
      : String(status).toUpperCase();
  switch (normalized) {
    case "UPLOADING":
    case "BUILDING":
    case "DEPLOYING":
    case "BUILT":
      return SandboxState.Queued;
    case "DEPLOYED":
      return SandboxState.Running;
    case "DELETING":
    case "DEACTIVATING":
      return SandboxState.Deleting;
    case "TERMINATED":
    case "DEACTIVATED":
    case "DELETED":
      return SandboxState.Deleted;
    case "ARCHIVING":
      return SandboxState.Archiving;
    case "ARCHIVED":
      return SandboxState.Archived;
    case "UNARCHIVING":
      return SandboxState.Starting;
    case "FAILED":
      return SandboxState.Error;
    default:
      return undefined;
  }
}

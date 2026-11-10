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
  /** Deployment status (DEPLOYED, FAILED, …). */
  readonly status?: unknown;
  /** Runtime state (RUNNING / STANDBY). */
  readonly state?: unknown;
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
  const runtime = runtimeState(source);
  if (source.status !== undefined && String(source.status) !== "") {
    target.providerState = String(source.status);
  } else if (runtime !== undefined) {
    target.providerState = runtime;
  }
  const state = mapSandboxState(source.status, runtime);
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

function runtimeState(source: BlaxelSandboxIdentity): string | undefined {
  if (source.state !== undefined && String(source.state) !== "")
    return String(source.state);
  // SandboxInstance keeps the wire model on a private field; read it when present.
  const nested = (source as { sandbox?: { state?: unknown } }).sandbox?.state;
  if (nested !== undefined && String(nested) !== "") return String(nested);
  return undefined;
}

/** Mirrors Go metadata.go: deployment status, then runtime standby. */
function mapSandboxState(
  status: unknown,
  runtime: string | undefined,
): SandboxState | undefined {
  const normalized =
    status === undefined || String(status) === ""
      ? ""
      : String(status).toUpperCase();
  const runtimeNormalized = runtime?.toUpperCase();
  switch (normalized) {
    case "UPLOADING":
    case "BUILDING":
    case "DEPLOYING":
    case "BUILT":
      return SandboxState.Queued;
    case "DEPLOYED":
      return runtimeNormalized === "STANDBY"
        ? SandboxState.Paused
        : SandboxState.Running;
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
      switch (runtimeNormalized) {
        case "RUNNING":
          return SandboxState.Running;
        case "STANDBY":
          return SandboxState.Paused;
        default:
          return undefined;
      }
  }
}

function mappingError(input: ProviderErrorInput): never {
  providerError(input);
}

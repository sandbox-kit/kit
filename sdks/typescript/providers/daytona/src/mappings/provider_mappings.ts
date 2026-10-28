// Provider field mappings between shared CreateOptions/Config and the native SDK.
// Inventoried in specs/providers/daytona.yaml groups; kept in sync by hand (not emitted).

import {
  ErrorKind,
  PolicyMode,
  type SandboxInfo,
  type Scope,
  CreateFieldLifetimeIdlePause,
  CreateFieldLifetimeIdleStop,
  CreateFieldLifetimeStoppedArchive,
  CreateFieldLifetimeStoppedDelete,
  CreateFieldResourcesCpuCores,
  CreateFieldResourcesDiskMib,
  CreateFieldResourcesMemoryMib,
} from "@sandbox-kit/sandbox";
import type {
  APIKeyCredentials,
  BearerTokenCredentials,
  Config,
  CreateOptions,
  Resources,
} from "@sandbox-kit/sandbox";
import type {
  CreateSandboxBaseParams as NativeCreateSandboxBaseParams,
  DaytonaConfig as NativeDaytonaConfig,
  Resources as NativeResources,
  Sandbox as NativeSandbox,
} from "@daytona/sdk";
import { providerError, type ProviderErrorInput } from "#errors.js";

type Mutable<T> = { -readonly [K in keyof T]: T[K] };
type PublicLifetime = NonNullable<CreateOptions["lifetime"]>;
type PublicAction = NonNullable<PublicLifetime["idleStop"]>;

/** Field-specific options for mapping one Daytona minute-based lifetime policy. */
interface MinutePolicySpec {
  field: string;
  operation: string;
  allowDisable: boolean;
  forbidImmediate: boolean;
  maxMinutes?: number;
}

export function mapLifetimePolicies(
  source: PublicLifetime | undefined,
  operation: string = "create",
): NativeCreateSandboxBaseParams {
  const target: Mutable<NativeCreateSandboxBaseParams> = {};
  if (source === undefined) return target;
  const idleStop = mapMinutePolicy(source.idleStop, {
    field: CreateFieldLifetimeIdleStop,
    operation,
    allowDisable: true,
    forbidImmediate: true,
  });
  if (idleStop !== undefined) target.autoStopInterval = idleStop;
  const idlePause = mapMinutePolicy(source.idlePause, {
    field: CreateFieldLifetimeIdlePause,
    operation,
    allowDisable: true,
    forbidImmediate: true,
  });
  if (idlePause !== undefined) target.autoPauseInterval = idlePause;
  const stoppedArchive = mapMinutePolicy(source.stoppedArchive, {
    field: CreateFieldLifetimeStoppedArchive,
    operation,
    allowDisable: false,
    forbidImmediate: true,
    maxMinutes: 43200,
  });
  if (stoppedArchive !== undefined) target.autoArchiveInterval = stoppedArchive;
  const stoppedDelete = mapMinutePolicy(source.stoppedDelete, {
    field: CreateFieldLifetimeStoppedDelete,
    operation,
    allowDisable: false,
    forbidImmediate: false,
  });
  if (stoppedDelete !== undefined) target.autoDeleteInterval = stoppedDelete;
  return target;
}

function mapMinutePolicy(
  policy: PublicAction | undefined,
  spec: MinutePolicySpec,
): number | undefined {
  if (policy === undefined) return undefined;
  const mode = policy.mode ?? PolicyMode.Default;
  if (mode === PolicyMode.Disabled) {
    if (!spec.allowDisable)
      mappingError({
        kind: ErrorKind.Unsupported,
        field: spec.field,
        message: "explicit disabling is not representable",
        operation: spec.operation,
      });
    return 0;
  }
  if (mode === PolicyMode.After) {
    if (policy.after === undefined)
      mappingError({
        kind: ErrorKind.InvalidArgument,
        field: spec.field,
        message: "duration is required",
        operation: spec.operation,
      });
    const afterMs = policy.after;
    if (spec.forbidImmediate && afterMs === 0)
      mappingError({
        kind: ErrorKind.Unsupported,
        field: spec.field,
        message: "immediate action is not representable",
        operation: spec.operation,
      });
    if (!Number.isFinite(afterMs) || afterMs < 0 || afterMs % 60_000 !== 0)
      mappingError({
        kind: ErrorKind.InvalidArgument,
        field: spec.field,
        message: "delay cannot be represented in whole minutes",
        operation: spec.operation,
      });
    const minutes = afterMs / 60_000;
    if (spec.maxMinutes !== undefined && minutes > spec.maxMinutes)
      mappingError({
        kind: ErrorKind.InvalidArgument,
        field: spec.field,
        message: "delay exceeds documented maximum",
        operation: spec.operation,
      });
    return minutes;
  }
  if (mode !== PolicyMode.Default)
    mappingError({
      kind: ErrorKind.InvalidArgument,
      field: spec.field,
      message: "unknown policy mode",
      operation: spec.operation,
    });
  return undefined;
}

export function mapCommonCreateFields(
  source: CreateOptions | undefined,
  operation: string = "create",
): NativeCreateSandboxBaseParams {
  const target: Mutable<NativeCreateSandboxBaseParams> = {};
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
      target.envVars = { ...value };
    }
  }
  {
    if (source.labels !== undefined) {
      const value = source.labels;
      target.labels = { ...value };
    }
  }
  return target;
}

export function mapResources(
  source: Resources | undefined,
  operation: string = "create",
): NativeResources {
  const target: Mutable<NativeResources> = {};
  if (source === undefined) return target;
  {
    if (source.cpuCores !== undefined) {
      const value = source.cpuCores;
      if (!Number.isInteger(value) || value > 2147483647)
        mappingError({
          kind: ErrorKind.InvalidArgument,
          field: CreateFieldResourcesCpuCores,
          message: "cpu_cores cannot be represented",
          operation,
        });
      target.cpu = Number(value);
    }
  }
  {
    if (source.memoryMib !== undefined) {
      const value = source.memoryMib;
      if (value % 1024n !== 0n || value / 1024n > 2147483647n)
        mappingError({
          kind: ErrorKind.InvalidArgument,
          field: CreateFieldResourcesMemoryMib,
          message: "memory_mib cannot be represented",
          operation,
        });
      target.memory = Number(value / 1024n);
    }
  }
  {
    if (source.diskMib !== undefined) {
      const value = source.diskMib;
      if (value % 1024n !== 0n || value / 1024n > 2147483647n)
        mappingError({
          kind: ErrorKind.InvalidArgument,
          field: CreateFieldResourcesDiskMib,
          message: "disk_mib cannot be represented",
          operation,
        });
      target.disk = Number(value / 1024n);
    }
  }
  return target;
}

export function mapSandboxInfo(
  source: NativeSandbox | undefined,
  operation: string = "create",
): Partial<SandboxInfo> {
  const target: Mutable<Partial<SandboxInfo>> = {};
  if (source === undefined) return target;
  {
    if (source.id !== undefined) {
      const value = source.id;
      target.id = value;
    }
  }
  {
    if (source.name !== undefined) {
      const value = source.name;
      target.name = value;
    }
  }
  {
    if (source.state !== undefined) {
      const value = source.state;
      target.providerState = String(value);
    }
  }
  {
    if (source.target !== undefined) {
      const value = source.target;
      target.region = value;
    }
  }
  {
    if (source.labels !== undefined) {
      const value = source.labels;
      target.labels = { ...value };
    }
  }
  target.provider = "daytona";
  return target;
}

export function mapAllocatedResources(
  source: NativeSandbox | undefined,
  operation: string = "create",
): Partial<Resources> {
  const target: Mutable<Partial<Resources>> = {};
  if (source === undefined) return target;
  if (
    !(
      (source.cpu ?? 0) > 0 &&
      (source.memory ?? 0) > 0 &&
      (source.disk ?? 0) > 0
    )
  )
    return target;
  {
    if (source.cpu !== undefined) {
      const value = source.cpu;
      target.cpuCores = Number(value);
    }
  }
  {
    if (source.memory !== undefined) {
      const value = source.memory;
      if (!Number.isSafeInteger(value) || value < 0)
        mappingError({
          kind: ErrorKind.InvalidResponse,
          field: "memory",
          message: "memory cannot be represented",
          operation,
        });
      target.memoryMib = BigInt(value) * 1024n;
    }
  }
  {
    if (source.disk !== undefined) {
      const value = source.disk;
      if (!Number.isSafeInteger(value) || value < 0)
        mappingError({
          kind: ErrorKind.InvalidResponse,
          field: "disk",
          message: "disk cannot be represented",
          operation,
        });
      target.diskMib = BigInt(value) * 1024n;
    }
  }
  return target;
}

export function mapClientSettings(
  source: Config | undefined,
  operation: string = "initialize",
): NativeDaytonaConfig {
  const target: Mutable<NativeDaytonaConfig> = {};
  if (source === undefined) return target;
  {
    if (source.endpoint !== undefined) {
      const value = source.endpoint;
      target.apiUrl = value;
    }
  }
  {
    if (source.region !== undefined) {
      const value = source.region;
      target.target = value;
    }
  }
  return target;
}

export function mapAPIKeyCredentials(
  source: APIKeyCredentials | undefined,
  operation: string = "initialize",
): NativeDaytonaConfig {
  const target: Mutable<NativeDaytonaConfig> = {};
  if (source === undefined) return target;
  {
    const value = source.key;
    target.apiKey = value;
  }
  return target;
}

export function mapBearerAuth(
  source: BearerTokenCredentials | undefined,
  operation: string = "initialize",
): NativeDaytonaConfig {
  const target: Mutable<NativeDaytonaConfig> = {};
  if (source === undefined) return target;
  {
    const value = source.token;
    target.jwtToken = value;
  }
  return target;
}

export function mapScope(
  source: Scope | undefined,
  operation: string = "initialize",
): NativeDaytonaConfig {
  const target: Mutable<NativeDaytonaConfig> = {};
  if (source === undefined) return target;
  {
    if (source.organizationId !== undefined) {
      const value = source.organizationId;
      target.organizationId = value;
    }
  }
  return target;
}

function mappingError(input: ProviderErrorInput): never {
  providerError(input);
}

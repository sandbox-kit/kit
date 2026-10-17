// Provider-specific orchestration around generated portable field mappings.
// Duration/oneof-sensitive mapping stays handwritten until area-5 regenerates providers.
import { durationMs, timestampFromDate } from "@bufbuild/protobuf/wkt";
import {
  EgressMode,
  ErrorKind,
  IsolationKind,
  SandboxState,
  SecretInjection,
  SnapshotKind,
  ValueOrigin,
  WaitCondition,
  materializeClientConfig,
  metadataObjectFromMap,
  sandboxError,
  type ConfigInput,
  type CreateOptions,
  type SandboxInfo,
} from "@sandbox-kit/sandbox";
import {
  type CreateSandboxBaseParams,
  type CreateSandboxFromImageParams,
  type CreateSandboxFromSnapshotParams,
  type DaytonaConfig,
  type Sandbox as DaytonaSandbox,
} from "@daytona/sdk";
import {
  mapAllocatedResources as generatedMapAllocatedResources,
  mapCommonCreateFields as generatedMapCommonCreateFields,
  mapLifetimePolicies as generatedMapLifetimePolicies,
  mapResources as generatedMapResources,
  mapSandboxInfo as generatedMapSandboxInfo,
} from "./provider.mappings.gen.js";
import { validateProviderCreate } from "./provider.checks.gen.js";
import { mapGeneratedClientConfig, type Environment } from "./provider.config.gen.js";

export type DaytonaCreateParams = CreateSandboxFromImageParams | CreateSandboxFromSnapshotParams;

export function mapClientConfig(config: ConfigInput, environment?: Environment): DaytonaConfig {
  const source = materializeClientConfig(config);
  const result = mapGeneratedClientConfig(source, environment);
  if (source.auth?.credentials.case === "bearerToken" && source.scope?.organizationId === undefined)
    configInvalid("scope.organizationId", "organization ID is required with bearer authentication");
  return result;
}

export function mapCreateRequest(request: CreateOptions): DaytonaCreateParams {
  validateProviderCreate(request);
  const base: CreateSandboxBaseParams = generatedMapCommonCreateFields(request);
  mapRuntime(request, base);
  mapIsolation(request, base);
  mapSecurity(request, base);
  mapNetwork(request, base);
  mapLifetime(request, base);
  mapProvisioning(request, base);
  mapObservability(request, base);
  if (request.placement !== undefined)
    unsupported("placement", "per-request Daytona placement is unsupported; use Config.region");
  if (request.readiness !== undefined) unsupported("readiness", "Daytona readiness probes are not mapped");
  if (request.providerOptions !== undefined) unsupported("providerOptions", "Daytona provider options are not mapped");
  const source = request.source;
  if (source?.kind?.case === "warmPool")
    unsupported("source.warmPool", "explicit warm-pool creation mapping is unavailable");
  if (source?.kind?.case === "image")
    return {
      ...base,
      image: nonblank(source.kind.value.reference, "source.image.reference"),
      ...(request.resources === undefined ? {} : { resources: mapResources(request.resources) }),
    };
  if (request.resources !== undefined) unsupported("resources", "Daytona resources require an image source");
  if (source?.kind?.case === "snapshot") {
    const snapshot = source.kind.value;
    if ((snapshot.kind ?? SnapshotKind.Default) !== SnapshotKind.Default)
      unsupported("source.snapshot.kind", "explicit snapshot-kind restore mapping is unavailable");
    if (snapshot.providerKind !== undefined)
      unsupported("source.snapshot.providerKind", "provider snapshot kind is unsupported");
    return { ...base, snapshot: nonblank(snapshot.reference, "source.snapshot.reference") };
  }
  return { ...base };
}

function mapRuntime(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.runtime;
  if (value === undefined) return;
  if (value.language !== undefined) target.language = value.language;
  if (value.user !== undefined) target.user = nonblank(value.user, "runtime.user");
  if (
    value.version !== undefined ||
    value.workingDirectory !== undefined ||
    value.entrypoint !== undefined ||
    value.pty !== undefined
  )
    unsupported("runtime", "runtime version, working directory, entrypoint, and PTY are not mapped by Daytona");
}

function mapIsolation(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.isolation;
  if (value === undefined) return;
  if ((value.kind ?? IsolationKind.Default) !== IsolationKind.Default || value.providerKind !== undefined)
    unsupported("isolation", "Daytona isolation class comes from the selected source");
  if (value.nestedVirtualization !== undefined) target.kvm = value.nestedVirtualization;
}

function mapSecurity(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.security;
  if (value === undefined) return;
  if (value.publicAccess !== undefined) target.public = value.publicAccess;
  if (value.workloadIdentity !== undefined)
    unsupported("security.workloadIdentity", "Daytona workload identity is not mapped");
  if (value.secrets !== undefined) {
    const secrets: Record<string, string> = {};
    for (const secret of value.secrets) {
      const variable =
        secret.environmentVariable === undefined
          ? ""
          : nonblank(secret.environmentVariable, "security.secrets.environmentVariable");
      if (variable === "") invalid("security.secrets.environmentVariable", "secret environment variable is required");
      const injection = secret.injection ?? SecretInjection.Default;
      if (injection !== SecretInjection.Default && injection !== SecretInjection.EgressPlaceholder)
        unsupported("security.secrets.injection", "Daytona secrets require egress-placeholder injection");
      if (Object.hasOwn(secrets, variable))
        invalid("security.secrets.environmentVariable", "duplicate secret environment variable");
      secrets[variable] = nonblank(secret.reference, "security.secrets.reference");
    }
    target.secrets = secrets;
  }
}

function mapNetwork(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.network;
  if (value === undefined) return;
  if (value.linkedSandbox !== undefined && request.lifetime?.ephemeral !== true)
    invalid("network.linkedSandbox", "linked sandboxes must be ephemeral");
  const egress = value.egress ?? EgressMode.Default;
  if (egress === EgressMode.BlockAll) target.networkBlockAll = true;
  else if (egress !== EgressMode.Default && egress !== EgressMode.AllowAll && egress !== EgressMode.Restricted)
    invalid("network.egress", "unknown egress mode");
  if (value.outboundCidrs !== undefined) {
    if (!value.outboundCidrs.entries?.length)
      unsupported("network.outboundCidrs", "empty CIDR allowlist mapping is not verified");
    target.networkAllowList = value.outboundCidrs.entries.join(",");
  }
  if (value.outboundDomains !== undefined) {
    if (!value.outboundDomains.entries?.length)
      unsupported("network.outboundDomains", "empty domain allowlist mapping is not verified");
    target.domainAllowList = value.outboundDomains.entries.join(",");
  }
  if (value.linkedSandbox !== undefined) target.linkedSandbox = nonblank(value.linkedSandbox, "network.linkedSandbox");
  if (value.outboundProxyUrl !== undefined) target.outboundProxyUrl = value.outboundProxyUrl;
  if (
    value.inboundCidrs !== undefined ||
    value.privateNetwork !== undefined ||
    value.proxyReference !== undefined ||
    value.customDomain !== undefined ||
    value.ports !== undefined
  )
    unsupported("network", "one or more Daytona network fields are not mapped");
}

function mapLifetime(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.lifetime;
  if (value === undefined) return;
  Object.assign(target, generatedMapLifetimePolicies(value));
  if (value.ephemeral !== undefined) target.ephemeral = value.ephemeral;
  if (value.maximumLifetime !== undefined)
    target.ttlMinutes = wholeMinutes(durationMs(value.maximumLifetime), "lifetime.maximumLifetime", true);
  if (value.idleTerminate !== undefined)
    unsupported("lifetime.idleTerminate", "Daytona idle termination is not mapped");
}

function mapProvisioning(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.provisioning;
  if (value === undefined) return;
  if (value.queueTimeout !== undefined)
    target.queueTimeout = wholeMinutes(durationMs(value.queueTimeout), "provisioning.queueTimeout", false);
  const waitFor = value.waitFor ?? WaitCondition.Default;
  if (waitFor !== WaitCondition.Default && waitFor !== WaitCondition.Started)
    unsupported("provisioning.waitFor", "the Daytona TypeScript SDK only returns after startup");
}

function mapObservability(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.observability;
  if (value === undefined) return;
  if (value.telemetryEndpoint !== undefined)
    target.otelEndpointOverride = nonblank(value.telemetryEndpoint, "observability.telemetryEndpoint");
  if (value.verbose !== undefined) unsupported("observability.verbose", "Daytona verbose creation is not mapped");
}

function mapResources(value: NonNullable<CreateOptions["resources"]>): ReturnType<typeof generatedMapResources> {
  if (value.cpuCores !== undefined && value.cpuCores <= 0) invalid("resources.cpuCores", "CPU must be positive");
  if (value.cpuLimitCores !== undefined || value.memoryLimitMib !== undefined)
    unsupported("resources", "Daytona resource limits are not mapped");
  return generatedMapResources(value);
}

function wholeMinutes(value: number, field: string, zeroAllowed: boolean): number {
  if (
    !Number.isFinite(value) ||
    value < 0 ||
    value % 60_000 !== 0 ||
    value / 60_000 > 2_147_483_647 ||
    (!zeroAllowed && value === 0)
  )
    invalid(field, "duration must be representable in whole minutes");
  return value / 60_000;
}

export function mapSandboxInfo(value: DaytonaSandbox, operation = "create"): SandboxInfo {
  const mapped = generatedMapSandboxInfo(value, operation);
  const origins: Record<string, ValueOrigin> = { ...(mapped.origins ?? {}) };
  const state = mapState(value.state === undefined ? undefined : String(value.state));
  const result: SandboxInfo = {
    ...mapped,
    id: mapped.id ?? value.id,
    provider: "daytona",
    state,
    labels: mapped.labels ?? {},
    origins,
  };
  const resources = generatedMapAllocatedResources(value, operation);
  if (Object.keys(resources).length > 0) {
    resultMutable(result).resources = resources;
    origins.resources = ValueOrigin.Provider;
  }
  if (value.sandboxClass !== undefined) {
    const native = String(value.sandboxClass);
    resultMutable(result).isolation =
      native === "container"
        ? { kind: IsolationKind.Container }
        : native === "linux-vm"
          ? { kind: IsolationKind.LinuxVm }
          : { kind: IsolationKind.ProviderSpecific, providerKind: native };
    origins.isolation = ValueOrigin.Provider;
  }
  if (state === SandboxState.Deleted) {
    resultMutable(result).terminal = true;
    origins.terminal = ValueOrigin.Provider;
  }
  const metadata: Record<string, string | number | boolean> = {};
  if (value.gpu > 0) metadata.gpu_count = value.gpu;
  if (value.gpuType !== undefined) metadata.gpu_type = String(value.gpuType);
  if (value.desiredState !== undefined) metadata.desired_state = String(value.desiredState);
  if (value.errorReason !== undefined) metadata.error_reason = value.errorReason;
  if (value.recoverable !== undefined) metadata.recoverable = value.recoverable;
  if (value.spot === true) metadata.spot = true;
  for (const [key, entry] of [
    ["spot_evicted_at", value.spotEvictedAt],
    ["backup_state", value.backupState],
    ["backup_created_at", value.backupCreatedAt],
    ["updated_at", value.updatedAt],
    ["last_activity_at", value.lastActivityAt],
    ["auto_destroy_at", value.autoDestroyAt],
    ["queue_timed_out_at", value.queueTimedOutAt],
    ["snapshot", value.snapshot],
    ["warm_pool_id", value.warmPoolId],
  ] as const)
    if (entry !== undefined) metadata[key] = String(entry);
  for (const [key, entry] of [
    ["auto_stop_interval_minutes", value.autoStopInterval],
    ["auto_pause_interval_minutes", value.autoPauseInterval],
    ["auto_archive_interval_minutes", value.autoArchiveInterval],
    ["auto_delete_interval_minutes", value.autoDeleteInterval],
  ] as const)
    if (entry !== undefined) metadata[key] = entry;
  if (Object.keys(metadata).length > 0) {
    const providerMetadata = metadataObjectFromMap(metadata);
    if (providerMetadata !== undefined) {
      resultMutable(result).providerMetadata = providerMetadata;
      origins.provider_metadata = ValueOrigin.Provider;
    }
  }
  if (value.createdAt !== undefined) {
    const date = new Date(value.createdAt);
    if (!Number.isNaN(date.valueOf())) {
      resultMutable(result).createdAt = timestampFromDate(date);
      origins.created_at = ValueOrigin.Provider;
    }
  }
  return result;
}

function resultMutable(value: SandboxInfo): { -readonly [K in keyof SandboxInfo]: SandboxInfo[K] } {
  return value;
}
function mapState(value: string | undefined): SandboxState {
  switch (value) {
    case "creating":
    case "restoring":
    case "pending_build":
    case "building_snapshot":
    case "pulling_snapshot":
    case "forking":
      return SandboxState.Queued;
    case "starting":
    case "resuming":
      return SandboxState.Starting;
    case "started":
      return SandboxState.Running;
    case "stopping":
      return SandboxState.Stopping;
    case "stopped":
      return SandboxState.Stopped;
    case "pausing":
      return SandboxState.Pausing;
    case "paused":
      return SandboxState.Paused;
    case "archiving":
      return SandboxState.Archiving;
    case "archived":
      return SandboxState.Archived;
    case "destroying":
      return SandboxState.Deleting;
    case "destroyed":
      return SandboxState.Deleted;
    case "error":
    case "build_failed":
      return SandboxState.Error;
    default:
      return SandboxState.Unknown;
  }
}

function nonblank(value: string, field: string): string {
  if (value.trim() === "") invalid(field, "value must not be blank");
  return value;
}
function invalid(field: string, message: string): never {
  throw sandboxError({
    kind: ErrorKind.InvalidArgument,
    provider: "daytona",
    operation: "create",
    field,
    message: "sandbox-kit daytona: " + message,
  });
}
function unsupported(field: string, message: string): never {
  throw sandboxError({
    kind: ErrorKind.Unsupported,
    provider: "daytona",
    operation: "create",
    field,
    message: "sandbox-kit daytona: " + message,
  });
}
function configInvalid(field: string, message: string): never {
  throw sandboxError({
    kind: ErrorKind.InvalidArgument,
    provider: "daytona",
    operation: "initialize",
    field,
    message: "sandbox-kit daytona: " + message,
  });
}

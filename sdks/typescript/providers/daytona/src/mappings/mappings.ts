// Provider-specific orchestration around generated portable field mappings.
import {
  EgressMode,
  ErrorKind,
  IsolationKind,
  SandboxState,
  SecretInjection,
  SnapshotKind,
  WaitCondition,
  type Config,
  type CreateOptions,
  type MetadataValue,
  type SandboxInfo,
  CreateFieldIsolation,
  CreateFieldLifetimeMaximumLifetime,
  CreateFieldNetwork,
  CreateFieldNetworkEgress,
  CreateFieldNetworkLinkedSandbox,
  CreateFieldObservabilityVerbose,
  CreateFieldProvisioningQueueTimeout,
  CreateFieldProvisioningTimeout,
  CreateFieldProvisioningWaitFor,
  CreateFieldReadiness,
  CreateFieldResources,
  CreateFieldResourcesCpuCores,
  CreateFieldRuntime,
  CreateFieldSecuritySecretsEnvironmentVariable,
  CreateFieldSecuritySecretsInjection,
  CreateFieldSourceSnapshotKind,
  CreateFieldRuntimeUser,
  CreateFieldSecuritySecretsReference,
  CreateFieldSecurityWorkloadIdentity,
  CreateFieldSourceImageReference,
  CreateFieldSourceSnapshotProviderKind,
  CreateFieldSourceSnapshotReference,
  CreateFieldSourceWarmPool,
  CreateFieldNetworkOutboundCidrs,
  CreateFieldNetworkOutboundDomains,
  CreateFieldLifetimeIdleTerminate,
  CreateFieldObservabilityTelemetryEndpoint,
  ConfigFieldScopeOrganizationId,
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
} from "./provider_mappings.js";
import { validateProviderCreate } from "#checks/provider_checks.js";
import { mapGeneratedClientConfig } from "#config/provider_config.js";
import { providerError } from "#errors.js";

export type DaytonaCreateParams =
  | CreateSandboxFromImageParams
  | CreateSandboxFromSnapshotParams;

export function mapClientConfig(config: Config): DaytonaConfig {
  const result = mapGeneratedClientConfig(config);
  if (
    config.auth?.bearerToken !== undefined &&
    config.scope?.organizationId === undefined
  )
    configInvalid(
      ConfigFieldScopeOrganizationId,
      "organization ID is required with bearer authentication",
    );
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
  if (request.readiness !== undefined)
    unsupported(
      CreateFieldReadiness,
      "Daytona readiness probes are not mapped",
    );
  const source = request.source;
  if (source?.warmPool !== undefined)
    unsupported(
      CreateFieldSourceWarmPool,
      "explicit warm-pool creation mapping is unavailable",
    );
  if (source?.image !== undefined)
    return {
      ...base,
      image: nonblank(source.image.reference, CreateFieldSourceImageReference),
      ...(request.resources === undefined
        ? {}
        : { resources: mapResources(request.resources) }),
    };
  if (request.resources !== undefined)
    unsupported(
      CreateFieldResources,
      "Daytona resources require an image source",
    );
  if (source?.snapshot !== undefined) {
    const snapshot = source.snapshot;
    if ((snapshot.kind ?? SnapshotKind.Default) !== SnapshotKind.Default)
      unsupported(
        CreateFieldSourceSnapshotKind,
        "explicit snapshot-kind restore mapping is unavailable",
      );
    if (snapshot.providerKind !== undefined)
      unsupported(
        CreateFieldSourceSnapshotProviderKind,
        "provider snapshot kind is unsupported",
      );
    return {
      ...base,
      snapshot: nonblank(
        snapshot.reference,
        CreateFieldSourceSnapshotReference,
      ),
    };
  }
  return { ...base };
}

function mapRuntime(
  request: CreateOptions,
  target: CreateSandboxBaseParams,
): void {
  const value = request.runtime;
  if (value === undefined) return;
  if (value.language !== undefined) target.language = value.language;
  if (value.user !== undefined)
    target.user = nonblank(value.user, CreateFieldRuntimeUser);
  if (
    value.version !== undefined ||
    value.workingDirectory !== undefined ||
    value.entrypoint !== undefined ||
    value.pty !== undefined
  )
    unsupported(
      CreateFieldRuntime,
      "runtime version, working directory, entrypoint, and PTY are not mapped by Daytona",
    );
}

function mapIsolation(
  request: CreateOptions,
  target: CreateSandboxBaseParams,
): void {
  const value = request.isolation;
  if (value === undefined) return;
  if (
    (value.kind ?? IsolationKind.Default) !== IsolationKind.Default ||
    value.providerKind !== undefined
  )
    unsupported(
      CreateFieldIsolation,
      "Daytona isolation class comes from the selected source",
    );
  if (value.nestedVirtualization !== undefined)
    target.kvm = value.nestedVirtualization;
}

function mapSecurity(
  request: CreateOptions,
  target: CreateSandboxBaseParams,
): void {
  const value = request.security;
  if (value === undefined) return;
  if (value.publicAccess !== undefined) target.public = value.publicAccess;
  if (value.workloadIdentity !== undefined)
    unsupported(
      CreateFieldSecurityWorkloadIdentity,
      "Daytona workload identity is not mapped",
    );
  if (value.secrets !== undefined) {
    const secrets: Record<string, string> = {};
    for (const secret of value.secrets) {
      const variable =
        secret.environmentVariable === undefined
          ? ""
          : nonblank(
              secret.environmentVariable,
              CreateFieldSecuritySecretsEnvironmentVariable,
            );
      if (variable === "")
        invalid(
          CreateFieldSecuritySecretsEnvironmentVariable,
          "secret environment variable is required",
        );
      const injection = secret.injection ?? SecretInjection.Default;
      if (
        injection !== SecretInjection.Default &&
        injection !== SecretInjection.EgressPlaceholder
      )
        unsupported(
          CreateFieldSecuritySecretsInjection,
          "Daytona secrets require egress-placeholder injection",
        );
      if (Object.hasOwn(secrets, variable))
        invalid(
          CreateFieldSecuritySecretsEnvironmentVariable,
          "duplicate secret environment variable",
        );
      secrets[variable] = nonblank(
        secret.reference,
        CreateFieldSecuritySecretsReference,
      );
    }
    target.secrets = secrets;
  }
}

function mapNetwork(
  request: CreateOptions,
  target: CreateSandboxBaseParams,
): void {
  const value = request.network;
  if (value === undefined) return;
  if (value.linkedSandbox !== undefined && request.lifetime?.ephemeral !== true)
    invalid(
      CreateFieldNetworkLinkedSandbox,
      "linked sandboxes must be ephemeral",
    );
  const egress = value.egress ?? EgressMode.Default;
  if (egress === EgressMode.BlockAll) target.networkBlockAll = true;
  else if (
    egress !== EgressMode.Default &&
    egress !== EgressMode.AllowAll &&
    egress !== EgressMode.Restricted
  )
    invalid(CreateFieldNetworkEgress, "unknown egress mode");
  if (value.outboundCidrs !== undefined) {
    if (!value.outboundCidrs.length)
      unsupported(
        CreateFieldNetworkOutboundCidrs,
        "empty CIDR allowlist mapping is not verified",
      );
    target.networkAllowList = value.outboundCidrs.join(",");
  }
  if (value.outboundDomains !== undefined) {
    if (!value.outboundDomains.length)
      unsupported(
        CreateFieldNetworkOutboundDomains,
        "empty domain allowlist mapping is not verified",
      );
    target.domainAllowList = value.outboundDomains.join(",");
  }
  if (value.linkedSandbox !== undefined)
    target.linkedSandbox = nonblank(
      value.linkedSandbox,
      CreateFieldNetworkLinkedSandbox,
    );
  if (value.outboundProxyUrl !== undefined)
    target.outboundProxyUrl = value.outboundProxyUrl;
  if (
    value.inboundCidrs !== undefined ||
    value.privateNetwork !== undefined ||
    value.proxyReference !== undefined ||
    value.customDomain !== undefined ||
    value.ports !== undefined
  )
    unsupported(
      CreateFieldNetwork,
      "one or more Daytona network fields are not mapped",
    );
}

function mapLifetime(
  request: CreateOptions,
  target: CreateSandboxBaseParams,
): void {
  const value = request.lifetime;
  if (value === undefined) return;
  Object.assign(target, generatedMapLifetimePolicies(value));
  if (value.ephemeral !== undefined) target.ephemeral = value.ephemeral;
  if (value.maximumLifetime !== undefined)
    target.ttlMinutes = wholeMinutes(
      value.maximumLifetime,
      CreateFieldLifetimeMaximumLifetime,
      true,
    );
  if (value.idleTerminate !== undefined)
    unsupported(
      CreateFieldLifetimeIdleTerminate,
      "Daytona idle termination is not mapped",
    );
}

function mapProvisioning(
  request: CreateOptions,
  target: CreateSandboxBaseParams,
): void {
  const value = request.provisioning;
  if (value === undefined) return;
  if (value.timeout !== undefined)
    requireWholeSeconds(value.timeout, CreateFieldProvisioningTimeout);
  if (value.queueTimeout !== undefined)
    target.queueTimeout = wholeMinutes(
      value.queueTimeout,
      CreateFieldProvisioningQueueTimeout,
      false,
    );
  const waitFor = value.waitFor ?? WaitCondition.Default;
  if (waitFor !== WaitCondition.Default && waitFor !== WaitCondition.Started)
    unsupported(
      CreateFieldProvisioningWaitFor,
      "the Daytona TypeScript SDK only returns after startup",
    );
}

function mapObservability(
  request: CreateOptions,
  target: CreateSandboxBaseParams,
): void {
  const value = request.observability;
  if (value === undefined) return;
  if (value.telemetryEndpoint !== undefined)
    target.otelEndpointOverride = nonblank(
      value.telemetryEndpoint,
      CreateFieldObservabilityTelemetryEndpoint,
    );
  if (value.verbose !== undefined)
    unsupported(
      CreateFieldObservabilityVerbose,
      "Daytona verbose creation is not mapped",
    );
}

function mapResources(
  value: NonNullable<CreateOptions["resources"]>,
): ReturnType<typeof generatedMapResources> {
  if (value.cpuCores !== undefined && value.cpuCores <= 0)
    invalid(CreateFieldResourcesCpuCores, "CPU must be positive");
  if (value.cpuLimitCores !== undefined || value.memoryLimitMib !== undefined)
    unsupported(CreateFieldResources, "Daytona resource limits are not mapped");
  return generatedMapResources(value);
}

function requireWholeSeconds(value: number, field: string): void {
  if (!Number.isFinite(value) || value < 0 || value % 1_000 !== 0)
    invalid(field, "timeout must be a whole number of seconds");
}

function wholeMinutes(
  value: number,
  field: string,
  zeroAllowed: boolean,
): number {
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

export function mapSandboxInfo(
  value: DaytonaSandbox,
  operation = "create",
): SandboxInfo {
  const mapped = generatedMapSandboxInfo(value, operation);
  const state = mapState(
    value.state === undefined ? undefined : String(value.state),
  );
  const result: SandboxInfo = {
    ...mapped,
    id: mapped.id ?? value.id,
    provider: "daytona",
    state,
    labels: mapped.labels ?? {},
  };
  const resources = generatedMapAllocatedResources(value, operation);
  if (Object.keys(resources).length > 0) {
    resultMutable(result).resources = resources;
  }
  if (value.sandboxClass !== undefined) {
    const native = String(value.sandboxClass);
    resultMutable(result).isolation =
      native === "container"
        ? { kind: IsolationKind.Container }
        : native === "linux-vm"
          ? { kind: IsolationKind.LinuxVm }
          : { kind: IsolationKind.ProviderSpecific, providerKind: native };
  }
  if (state === SandboxState.Deleted) {
    resultMutable(result).terminal = true;
  }
  const metadata: { [key: string]: MetadataValue } = {};
  if (value.gpu > 0) metadata.gpu_count = scalarMetadata(value.gpu);
  assignText(metadata, "gpu_type", value.gpuType);
  assignText(metadata, "desired_state", value.desiredState);
  assignText(metadata, "error_reason", value.errorReason);
  if (value.recoverable !== undefined)
    metadata.recoverable = scalarMetadata(value.recoverable);
  if (value.spot === true) metadata.spot = scalarMetadata(true);
  assignText(metadata, "spot_evicted_at", value.spotEvictedAt);
  assignText(metadata, "backup_state", value.backupState);
  assignText(metadata, "backup_created_at", value.backupCreatedAt);
  assignText(metadata, "updated_at", value.updatedAt);
  assignText(metadata, "last_activity_at", value.lastActivityAt);
  assignText(metadata, "auto_destroy_at", value.autoDestroyAt);
  assignText(metadata, "queue_timed_out_at", value.queueTimedOutAt);
  assignText(metadata, "snapshot", value.snapshot);
  assignText(metadata, "warm_pool_id", value.warmPoolId);
  assignNumber(metadata, "auto_stop_interval_minutes", value.autoStopInterval);
  assignNumber(
    metadata,
    "auto_pause_interval_minutes",
    value.autoPauseInterval,
  );
  assignNumber(
    metadata,
    "auto_archive_interval_minutes",
    value.autoArchiveInterval,
  );
  assignNumber(
    metadata,
    "auto_delete_interval_minutes",
    value.autoDeleteInterval,
  );
  if (Object.keys(metadata).length > 0) {
    resultMutable(result).providerMetadata = { values: metadata };
  }
  if (value.createdAt !== undefined) {
    const date = new Date(value.createdAt);
    if (!Number.isNaN(date.valueOf())) {
      resultMutable(result).createdAt = date;
    }
  }
  return result;
}

function resultMutable(value: SandboxInfo): {
  -readonly [K in keyof SandboxInfo]: SandboxInfo[K];
} {
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

function assignText(
  metadata: { [key: string]: MetadataValue },
  key: string,
  value: unknown,
): void {
  if (value !== undefined) metadata[key] = scalarMetadata(String(value));
}

function assignNumber(
  metadata: { [key: string]: MetadataValue },
  key: string,
  value: number | undefined,
): void {
  if (value !== undefined) metadata[key] = scalarMetadata(value);
}

/** Scalar metadata for fields this provider just read. Matches the public map conversion. */
function scalarMetadata(value: string | number | boolean): MetadataValue {
  if (typeof value === "boolean") return { value: { case: "boolean", value } };
  if (typeof value === "string") return { value: { case: "text", value } };
  if (!Number.isFinite(value)) {
    throw new Error("sandbox-kit: metadata numbers must be finite");
  }
  if (Number.isInteger(value)) {
    const asBigInt = BigInt(value);
    if (value < 0) return { value: { case: "signedInteger", value: asBigInt } };
    return { value: { case: "unsignedInteger", value: asBigInt } };
  }
  return { value: { case: "number", value } };
}

function nonblank(value: string, field: string): string {
  if (value.trim() === "") invalid(field, "value must not be blank");
  return value;
}
function invalid(field: string, message: string): never {
  providerError({
    kind: ErrorKind.InvalidArgument,
    operation: "create",
    field,
    message,
  });
}
function unsupported(field: string, message: string): never {
  providerError({
    kind: ErrorKind.Unsupported,
    operation: "create",
    field,
    message,
  });
}
function configInvalid(field: string, message: string): never {
  providerError({
    kind: ErrorKind.InvalidArgument,
    operation: "initialize",
    field,
    message,
  });
}

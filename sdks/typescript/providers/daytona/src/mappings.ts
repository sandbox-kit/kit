// Provider-specific orchestration around generated portable field mappings.
import {
  EgressMode,
  ErrorKind,
  IsolationKind,
  SandboxState,
  SecretInjection,
  SnapshotKind,
  ValueOrigin,
  WaitCondition,
  sandboxError,
  type Config,
  type CreateOptions,
  type CreateResult,
  type SandboxError,
  type SandboxInfo,
} from "@sandbox-kit/sandbox";
import {
  DaytonaError,
  type CreateSandboxBaseParams,
  type CreateSandboxFromImageParams,
  type CreateSandboxFromSnapshotParams,
  type DaytonaConfig,
  type Sandbox as DaytonaSandbox,
} from "@daytona/sdk";
import {
  mapAPIKeyCredentials as generatedMapAPIKeyCredentials,
  mapBearerAuth as generatedMapBearerAuth,
  mapClientSettings as generatedMapClientSettings,
  mapCommonCreateFields as generatedMapCommonCreateFields,
  mapLifetimePolicies as generatedMapLifetimePolicies,
  mapResources as generatedMapResources,
  mapScope as generatedMapScope,
} from "./provider.mappings.gen.js";

export type DaytonaCreateParams = CreateSandboxFromImageParams | CreateSandboxFromSnapshotParams;

type Environment = Readonly<Record<string, string | undefined>>;

export function mapClientConfig(config: Config, environment: Environment = environmentVariables()): DaytonaConfig {
  const result: DaytonaConfig = {
    ...generatedMapClientSettings(config),
    apiUrl: config.endpoint ?? environment.DAYTONA_API_URL ?? environment.DAYTONA_SERVER_URL ?? "https://app.daytona.io/api",
    ...(config.region === undefined ? {} : { target: configNonblank(config.region, "region") }),
  };
  if (config.endpoint !== undefined) secureEndpoint(config.endpoint);
  if (config.timeout !== undefined && (!Number.isFinite(config.timeout) || config.timeout < 0)) configInvalid("timeout", "timeout must be a finite non-negative duration");
  if (config.scope?.appName !== undefined || config.scope?.environment !== undefined || config.scope?.projectId !== undefined) configUnsupported("scope", "unsupported Daytona client scope");
  if (config.scope?.organizationId !== undefined) { configNonblank(config.scope.organizationId, "scope.organizationId"); Object.assign(result, generatedMapScope(config.scope)); }
  if (config.auth?.credentials === "apiKey") { configNonblank(config.auth.apiKey.key, "auth.apiKey.key"); Object.assign(result, generatedMapAPIKeyCredentials(config.auth.apiKey)); }
  else if (config.auth?.credentials === "bearerToken") {
    configNonblank(config.auth.bearerToken.token, "auth.bearerToken.token");
    Object.assign(result, generatedMapBearerAuth(config.auth.bearerToken));
    if (config.scope?.organizationId === undefined) configInvalid("scope.organizationId", "organization ID is required with bearer authentication");
    configNonblank(config.scope.organizationId, "scope.organizationId");
    Object.assign(result, generatedMapScope(config.scope));
  } else if (config.auth?.credentials !== undefined) configUnsupported("auth", "unsupported Daytona authentication mode");
  return result;
}

function environmentVariables(): Environment {
  return (globalThis as { readonly process?: { readonly env?: Environment } }).process?.env ?? {};
}

export function mapCreateRequest(request: CreateOptions): DaytonaCreateParams {
  const base: CreateSandboxBaseParams = generatedMapCommonCreateFields(request);
  mapRuntime(request, base); mapIsolation(request, base); mapSecurity(request, base); mapNetwork(request, base); mapLifetime(request, base); mapProvisioning(request, base); mapObservability(request, base);
  if (request.placement !== undefined) unsupported("placement", "per-request Daytona placement is unsupported; use Config.region");
  if (request.readiness !== undefined) unsupported("readiness", "Daytona readiness probes are not mapped");
  if (request.providerOptions !== undefined) unsupported("providerOptions", "Daytona provider options are not mapped");
  const source = request.source;
  if (source?.kind === "warmPool") unsupported("source.warmPool", "explicit warm-pool creation mapping is unavailable");
  if (source?.kind === "image") return { ...base, image: nonblank(source.image.reference, "source.image.reference"), ...(request.resources === undefined ? {} : { resources: mapResources(request.resources) }) };
  if (request.resources !== undefined) unsupported("resources", "Daytona resources require an image source");
  if (source?.kind === "snapshot") {
    if (source.snapshot.kind !== SnapshotKind.Default) unsupported("source.snapshot.kind", "explicit snapshot-kind restore mapping is unavailable");
    if (source.snapshot.providerKind !== undefined) unsupported("source.snapshot.providerKind", "provider snapshot kind is unsupported");
    return { ...base, snapshot: nonblank(source.snapshot.reference, "source.snapshot.reference") };
  }
  return { ...base };
}

function mapRuntime(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.runtime; if (value === undefined) return;
  if (value.language !== undefined) { const language = nonblank(value.language, "runtime.language"); if (!["python", "javascript", "typescript"].includes(language)) invalid("runtime.language", "unsupported Daytona language"); target.language = language; }
  if (value.user !== undefined) target.user = nonblank(value.user, "runtime.user");
  if (value.version !== undefined || value.workingDirectory !== undefined || value.entrypoint !== undefined || value.pty !== undefined) unsupported("runtime", "runtime version, working directory, entrypoint, and PTY are not mapped by Daytona");
}

function mapIsolation(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.isolation; if (value === undefined) return;
  if (value.kind !== IsolationKind.Default || value.providerKind !== undefined) unsupported("isolation", "Daytona isolation class comes from the selected source");
  if (value.nestedVirtualization !== undefined) target.kvm = value.nestedVirtualization;
}

function mapSecurity(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.security; if (value === undefined) return;
  if (value.publicAccess !== undefined) target.public = value.publicAccess;
  if (value.workloadIdentity !== undefined) unsupported("security.workloadIdentity", "Daytona workload identity is not mapped");
  if (value.secrets !== undefined) { const secrets: Record<string, string> = {}; for (const secret of value.secrets) { const variable = secret.environmentVariable === undefined ? "" : nonblank(secret.environmentVariable, "security.secrets.environmentVariable"); if (variable === "") invalid("security.secrets.environmentVariable", "secret environment variable is required"); if (secret.injection !== SecretInjection.Default && secret.injection !== SecretInjection.EgressPlaceholder) unsupported("security.secrets.injection", "Daytona secrets require egress-placeholder injection"); if (Object.hasOwn(secrets, variable)) invalid("security.secrets.environmentVariable", "duplicate secret environment variable"); secrets[variable] = nonblank(secret.reference, "security.secrets.reference"); } target.secrets = secrets; }
}

function mapNetwork(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.network; if (value === undefined) return;
  if (value.egress === EgressMode.BlockAll) target.networkBlockAll = true;
  else if (value.egress !== EgressMode.Default && value.egress !== EgressMode.AllowAll && value.egress !== EgressMode.Restricted) invalid("network.egress", "unknown egress mode");
  if (value.outboundCidrs !== undefined) { if (!value.outboundCidrs.entries?.length) unsupported("network.outboundCidrs", "empty CIDR allowlist mapping is not verified"); target.networkAllowList = value.outboundCidrs.entries.join(","); }
  if (value.outboundDomains !== undefined) { if (!value.outboundDomains.entries?.length) unsupported("network.outboundDomains", "empty domain allowlist mapping is not verified"); target.domainAllowList = value.outboundDomains.entries.join(","); }
  if (value.linkedSandbox !== undefined) target.linkedSandbox = nonblank(value.linkedSandbox, "network.linkedSandbox");
  if (value.outboundProxyUrl !== undefined) { httpUrl(value.outboundProxyUrl, "network.outboundProxyUrl"); target.outboundProxyUrl = value.outboundProxyUrl; }
  if (value.inboundCidrs !== undefined || value.privateNetwork !== undefined || value.proxyReference !== undefined || value.customDomain !== undefined || value.ports !== undefined) unsupported("network", "one or more Daytona network fields are not mapped");
}

function mapLifetime(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.lifetime; if (value === undefined) return;
	Object.assign(target, generatedMapLifetimePolicies(value));
  if (value.ephemeral !== undefined) target.ephemeral = value.ephemeral;
  if (value.maximumLifetime !== undefined) target.ttlMinutes = wholeMinutes(value.maximumLifetime, "lifetime.maximumLifetime", true);
  if (value.idleTerminate !== undefined) unsupported("lifetime.idleTerminate", "Daytona idle termination is not mapped");
  if ((target.autoStopInterval ?? 0) > 0 && (target.autoPauseInterval ?? 0) > 0) invalid("lifetime", "idle stop and idle pause are mutually exclusive");
  if ((target.autoPauseInterval ?? 0) > 0 && (value.ephemeral === true || target.autoDeleteInterval === 0)) invalid("lifetime.idlePause", "ephemeral sandboxes cannot auto-pause");
  if (value.ephemeral === true && value.stoppedArchive !== undefined) invalid("lifetime.stoppedArchive", "ephemeral sandboxes cannot auto-archive");
}

function mapProvisioning(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.provisioning; if (value === undefined) return;
  if (value.queueTimeout !== undefined) target.queueTimeout = wholeMinutes(value.queueTimeout, "provisioning.queueTimeout", false);
  if (value.waitFor !== WaitCondition.Default && value.waitFor !== WaitCondition.Started) unsupported("provisioning.waitFor", "the Daytona TypeScript SDK only returns after startup");
}

function mapObservability(request: CreateOptions, target: CreateSandboxBaseParams): void {
  const value = request.observability; if (value === undefined) return;
  if (value.telemetryEndpoint !== undefined) target.otelEndpointOverride = nonblank(value.telemetryEndpoint, "observability.telemetryEndpoint");
  if (value.verbose !== undefined) unsupported("observability.verbose", "Daytona verbose creation is not mapped");
}

function mapResources(value: NonNullable<CreateOptions["resources"]>): ReturnType<typeof generatedMapResources> {
  if (value.cpuCores !== undefined && value.cpuCores <= 0) invalid("resources.cpuCores", "CPU must be positive");
  if (value.cpuLimitCores !== undefined || value.memoryLimitMib !== undefined) unsupported("resources", "Daytona resource limits are not mapped");
  return generatedMapResources(value);
}

function wholeMinutes(value: number, field: string, zeroAllowed: boolean): number { if (!Number.isFinite(value) || value < 0 || value % 60_000 !== 0 || value / 60_000 > 2_147_483_647 || (!zeroAllowed && value === 0)) invalid(field, "duration must be representable in whole minutes"); return value / 60_000; }

export function mapSandboxInfo(value: DaytonaSandbox): SandboxInfo {
  const origins: Record<string, ValueOrigin> = { id: ValueOrigin.Provider };
  const state = mapState(value.state === undefined ? undefined : String(value.state));
  const result: SandboxInfo = { id: value.id, provider: "daytona", state, origins };
  if (value.name !== "") { resultMutable(result).name = value.name; origins.name = ValueOrigin.Provider; }
  if (value.state !== undefined) { resultMutable(result).providerState = String(value.state); origins.provider_state = ValueOrigin.Provider; }
  if (value.target !== "") { resultMutable(result).region = value.target; origins.region = ValueOrigin.Provider; }
  if (value.labels !== undefined) { resultMutable(result).labels = { ...value.labels }; origins.labels = ValueOrigin.Provider; }
  if (value.cpu > 0 && value.memory > 0 && value.disk > 0) { resultMutable(result).resources = { cpuCores: value.cpu, memoryMib: BigInt(value.memory) * 1024n, diskMib: BigInt(value.disk) * 1024n }; origins.resources = ValueOrigin.Provider; }
  if (value.sandboxClass !== undefined) { const native = String(value.sandboxClass); resultMutable(result).isolation = native === "container" ? { kind: IsolationKind.Container } : native === "linux-vm" ? { kind: IsolationKind.LinuxVm } : { kind: IsolationKind.ProviderSpecific, providerKind: native }; origins.isolation = ValueOrigin.Provider; }
  if (state === SandboxState.Deleted) { resultMutable(result).terminal = true; origins.terminal = ValueOrigin.Provider; }
  const metadata: Record<string, string | number | boolean> = {};
  if (value.gpu > 0) metadata.gpu_count = value.gpu;
  if (value.gpuType !== undefined) metadata.gpu_type = String(value.gpuType);
  if (value.desiredState !== undefined) metadata.desired_state = String(value.desiredState);
  if (value.errorReason !== undefined) metadata.error_reason = value.errorReason;
  if (value.recoverable !== undefined) metadata.recoverable = value.recoverable;
  if (value.spot === true) metadata.spot = true;
  for (const [key, entry] of [["spot_evicted_at", value.spotEvictedAt], ["backup_state", value.backupState], ["backup_created_at", value.backupCreatedAt], ["updated_at", value.updatedAt], ["last_activity_at", value.lastActivityAt], ["auto_destroy_at", value.autoDestroyAt], ["queue_timed_out_at", value.queueTimedOutAt], ["snapshot", value.snapshot], ["warm_pool_id", value.warmPoolId]] as const) if (entry !== undefined) metadata[key] = String(entry);
  for (const [key, entry] of [["auto_stop_interval_minutes", value.autoStopInterval], ["auto_pause_interval_minutes", value.autoPauseInterval], ["auto_archive_interval_minutes", value.autoArchiveInterval], ["auto_delete_interval_minutes", value.autoDeleteInterval]] as const) if (entry !== undefined) metadata[key] = entry;
  if (Object.keys(metadata).length > 0) { resultMutable(result).providerMetadata = metadata; origins.provider_metadata = ValueOrigin.Provider; }
  if (value.createdAt !== undefined) { const date = new Date(value.createdAt); if (!Number.isNaN(date.valueOf())) { resultMutable(result).createdAt = date; origins.created_at = ValueOrigin.Provider; } }
  return result;
}

function resultMutable(value: SandboxInfo): { -readonly [K in keyof SandboxInfo]: SandboxInfo[K] } { return value; }
function mapState(value: string | undefined): SandboxState { switch (value) { case "creating": case "restoring": case "pending_build": case "building_snapshot": case "pulling_snapshot": case "forking": return SandboxState.Queued; case "starting": case "resuming": return SandboxState.Starting; case "started": return SandboxState.Running; case "stopping": return SandboxState.Stopping; case "stopped": return SandboxState.Stopped; case "pausing": return SandboxState.Pausing; case "paused": return SandboxState.Paused; case "archiving": return SandboxState.Archiving; case "archived": return SandboxState.Archived; case "destroying": return SandboxState.Deleting; case "destroyed": return SandboxState.Deleted; case "error": case "build_failed": return SandboxState.Error; default: return SandboxState.Unknown; } }

export function mapDaytonaError(error: unknown, operation: string): SandboxError {
	if (error instanceof DOMException && error.name === "AbortError") return sandboxError({ kind: ErrorKind.Canceled, provider: "daytona", operation, field: "", message: error.message }, error);
	if (error instanceof DOMException && error.name === "TimeoutError") return sandboxError({ kind: ErrorKind.Timeout, provider: "daytona", operation, field: "", message: error.message }, error);
  if (!(error instanceof DaytonaError)) return sandboxError({ kind: ErrorKind.Unknown, provider: "daytona", operation, field: "", message: error instanceof Error ? error.message : "sandbox-kit daytona: provider operation failed" }, error);
  const byStatus: Readonly<Record<number, ErrorKind>> = { 400: ErrorKind.InvalidArgument, 401: ErrorKind.Authentication, 403: ErrorKind.PermissionDenied, 404: ErrorKind.NotFound, 408: ErrorKind.Timeout, 409: ErrorKind.Conflict, 410: ErrorKind.NotFound, 412: ErrorKind.FailedPrecondition, 422: ErrorKind.InvalidArgument, 429: ErrorKind.RateLimited, 500: ErrorKind.Internal, 502: ErrorKind.Unavailable, 503: ErrorKind.Unavailable, 504: ErrorKind.Timeout };
  const kind = error.code === "SANDBOX_SPOT_EVICTED" ? ErrorKind.Unavailable : error.statusCode === 0 ? ErrorKind.Unavailable : error.statusCode === undefined ? ErrorKind.Unknown : (byStatus[error.statusCode] ?? ErrorKind.Unknown);
  return sandboxError({ kind, provider: "daytona", operation, field: "", message: error.message, ...(error.statusCode === undefined ? {} : { statusCode: error.statusCode }), ...(error.code === undefined ? {} : { providerCode: error.code }), ...(error.source === undefined ? {} : { providerSource: error.source }) }, error);
}

function nonblank(value: string, field: string): string { if (value.trim() === "") invalid(field, "value must not be blank"); return value; }
function configNonblank(value: string, field: string): string { if (value.trim() === "") configInvalid(field, "value must not be blank"); return value; }
function secureEndpoint(value: string): void { let url: URL; try { url = new URL(value); } catch { configInvalid("endpoint", "endpoint must be a valid URL"); } if (url!.protocol !== "https:" && url!.hostname !== "localhost" && url!.hostname !== "127.0.0.1") configInvalid("endpoint", "endpoint must use HTTPS except for loopback development"); }
function httpUrl(value: string, field: string): void { let url: URL; try { url = new URL(value); } catch { invalid(field, "value must be a valid HTTP URL"); } if (url!.protocol !== "http:" && url!.protocol !== "https:") invalid(field, "value must be an HTTP URL"); }
function invalid(field: string, message: string): never { throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: "daytona", operation: "create", field, message: "sandbox-kit daytona: " + message }); }
function unsupported(field: string, message: string): never { throw sandboxError({ kind: ErrorKind.Unsupported, provider: "daytona", operation: "create", field, message: "sandbox-kit daytona: " + message }); }
function configInvalid(field: string, message: string): never { throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: "daytona", operation: "initialize", field, message: "sandbox-kit daytona: " + message }); }
function configUnsupported(field: string, message: string): never { throw sandboxError({ kind: ErrorKind.Unsupported, provider: "daytona", operation: "initialize", field, message: "sandbox-kit daytona: " + message }); }


// Provider-specific orchestration around generated portable field mappings.
import type {
  Config as NativeClientConfig,
  Port,
  SandboxCreateConfiguration,
  SandboxLifecycle,
} from "@blaxel/core";
import {
  EgressMode,
  ErrorKind,
  IsolationKind,
  PolicyMode,
  PortTransport,
  SandboxState,
  WaitCondition,
  type Config,
  type CreateOptions,
  type SandboxInfo,
  ConfigFieldRegion,
  CreateFieldIsolation,
  CreateFieldLifetime,
  CreateFieldLifetimeIdleTerminate,
  CreateFieldLifetimeIdleTerminateAfter,
  CreateFieldLifetimeIdleTerminateMode,
  CreateFieldLifetimeMaximumLifetime,
  CreateFieldName,
  CreateFieldNetwork,
  CreateFieldNetworkPortsPort,
  CreateFieldNetworkPortsTransport,
  CreateFieldObservability,
  CreateFieldProvisioningQueueTimeout,
  CreateFieldProvisioningWaitFor,
  CreateFieldReadiness,
  CreateFieldRuntime,
  CreateFieldRuntimeWorkingDirectory,
  CreateFieldSecurity,
  CreateFieldSourceImageReference,
} from "@sandbox-kit/sandbox";
import {
  mapAllocatedResources,
  mapCommonCreateFields as generatedMapCommonCreateFields,
  mapResources as generatedMapResources,
  mapSandboxIdentity as generatedMapSandboxIdentity,
  type BlaxelSandboxIdentity,
} from "./provider_mappings.js";
import { validateProviderCreate } from "#checks/provider_checks.js";
import { mapGeneratedClientConfig } from "#config/provider_config.js";
import type { ProviderState } from "#config/provider_config.js";
import { providerError } from "#errors.js";

export interface BlaxelCreatePlan {
  readonly params: SandboxCreateConfiguration;
}

export function mapClientConfig(config: Config): NativeClientConfig {
  return mapGeneratedClientConfig(config);
}

export function mapCreateRequest(
  request: CreateOptions,
  config: ProviderState,
): BlaxelCreatePlan {
  validateProviderCreate(request);
  if (config.region === undefined || config.region.trim() === "")
    invalid(ConfigFieldRegion, "Config.region is required");

  const params: SandboxCreateConfiguration =
    generatedMapCommonCreateFields(request);
  params.name = nonblank(request.name!, CreateFieldName);
  params.image = nonblank(
    request.source!.image!.reference,
    CreateFieldSourceImageReference,
  );
  params.region = nonblank(config.region, ConfigFieldRegion);

  mapRuntime(request);
  mapIsolation(request);
  mapResources(request, params);
  mapNetwork(request, params);
  mapSecurity(request);
  mapLifetime(request, params);
  mapReadiness(request);
  mapObservability(request);
  mapProvisioning(request);

  return { params };
}

function mapRuntime(r: CreateOptions): void {
  const v = r.runtime;
  if (!v) return;
  if (v.workingDirectory !== undefined)
    unsupported(
      CreateFieldRuntimeWorkingDirectory,
      "Blaxel create-time working directory is not mapped",
    );
  if (hasUnsupportedRuntimeFields(v))
    unsupported(
      CreateFieldRuntime,
      "runtime user, language, version, entrypoint, and PTY are not mapped by Blaxel",
    );
}

function mapIsolation(r: CreateOptions): void {
  const v = r.isolation;
  if (!v) return;
  if (hasUnsupportedIsolation(v))
    unsupported(CreateFieldIsolation, "Blaxel isolation class is not mapped");
}

function mapResources(r: CreateOptions, p: SandboxCreateConfiguration): void {
  const v = r.resources;
  if (!v) return;
  Object.assign(p, generatedMapResources(v));
}

function mapNetwork(r: CreateOptions, p: SandboxCreateConfiguration): void {
  const v = r.network;
  if (!v) return;
  rejectUnsupportedNetwork(v);

  if (v.outboundDomains !== undefined) {
    p.network = {
      ...(p.network ?? {}),
      allowedDomains: [...v.outboundDomains],
    };
  }

  if (v.ports) p.ports = mapPorts(v.ports);
}

const reservedBlaxelPorts = new Set([80, 443, 8080]);

function mapPorts(
  bindings: NonNullable<NonNullable<CreateOptions["network"]>["ports"]>,
): Port[] {
  const ports: Port[] = [];
  for (const x of bindings) {
    const port = x.port;
    if (port === undefined)
      invalid(CreateFieldNetworkPortsPort, "port is required");
    if (reservedBlaxelPorts.has(port))
      invalid(
        CreateFieldNetworkPortsPort,
        "ports 80, 443, and 8080 are reserved by Blaxel",
      );
    const protocol = blaxelPortProtocol(x.transport ?? PortTransport.Default);
    ports.push(
      protocol === undefined ? { target: port } : { target: port, protocol },
    );
  }
  return ports;
}

/** Default egress is a no-op (matches Go). Everything else here is unsupported. */
function rejectUnsupportedNetwork(
  network: NonNullable<CreateOptions["network"]>,
): void {
  if (hasUnsupportedNetwork(network)) {
    unsupported(
      CreateFieldNetwork,
      "requested network settings are unsupported by Blaxel",
    );
  }
}

function blaxelPortProtocol(
  transport: PortTransport,
): Port["protocol"] | undefined {
  // Kit Default leaves the protocol to the provider. Blaxel defaults to HTTP
  // when protocol is omitted (normalizePorts).
  if (transport === PortTransport.Default) return undefined;
  if (transport === PortTransport.Plain) return "TCP";
  if (transport === PortTransport.Tls) return "TLS";
  // Blaxel has no HTTP/2 port protocol (Kit Http2Tls included).
  invalid(
    CreateFieldNetworkPortsTransport,
    "requested port transport is unsupported",
  );
}

function mapSecurity(r: CreateOptions): void {
  rejectIfSet(
    r.security !== undefined,
    CreateFieldSecurity,
    "Blaxel security fields are not mapped",
  );
}

function mapLifetime(r: CreateOptions, p: SandboxCreateConfiguration): void {
  const v = r.lifetime;
  if (!v) return;
  // Idle first (matches Go) so idle_terminate errors win over unmapped sibling fields.
  if (v.maximumLifetime !== undefined) {
    p.ttl = durationString({
      ms: v.maximumLifetime,
      field: CreateFieldLifetimeMaximumLifetime,
      minimumMs: 5 * 60_000,
    });
  }

  if (v.idleTerminate) {
    const lifecycle = mapIdleTerminate(v.idleTerminate);
    if (lifecycle !== undefined) p.lifecycle = lifecycle;
  }

  if (hasUnsupportedLifetimeFields(v))
    unsupported(
      CreateFieldLifetime,
      "one or more Blaxel lifetime fields are not mapped",
    );
}

function mapIdleTerminate(
  action: NonNullable<NonNullable<CreateOptions["lifetime"]>["idleTerminate"]>,
): SandboxLifecycle | undefined {
  const mode = action.mode ?? PolicyMode.Default;
  if (mode === PolicyMode.Disabled) {
    unsupported(
      CreateFieldLifetimeIdleTerminate,
      "disabling idle terminate is not representable",
    );
  }
  if (mode === PolicyMode.After) {
    if (action.after === undefined)
      invalid(CreateFieldLifetimeIdleTerminateAfter, "duration is required");
    const value = durationString({
      ms: action.after,
      field: CreateFieldLifetimeIdleTerminateAfter,
      minimumMs: 60_000,
    });
    return {
      expirationPolicies: [{ type: "ttl-idle", action: "delete", value }],
    };
  }
  if (mode !== PolicyMode.Default)
    invalid(CreateFieldLifetimeIdleTerminateMode, "unknown policy mode");
  return undefined;
}

function mapReadiness(r: CreateOptions): void {
  rejectIfSet(
    r.readiness !== undefined,
    CreateFieldReadiness,
    "Blaxel readiness probes are not mapped",
  );
}

function mapObservability(r: CreateOptions): void {
  rejectIfSet(
    r.observability !== undefined,
    CreateFieldObservability,
    "Blaxel observability fields are not mapped",
  );
}

// Kit Client.create applies provisioning.timeout as the abort deadline.
// This provider does not forward a native create timeout — use the Kit deadline.
function mapProvisioning(r: CreateOptions): void {
  const v = r.provisioning;
  if (!v) return;
  if (v.queueTimeout !== undefined)
    unsupported(
      CreateFieldProvisioningQueueTimeout,
      "Blaxel queue timeout is unsupported",
    );

  const waitFor = v.waitFor ?? WaitCondition.Default;
  if (waitFor === WaitCondition.Default || waitFor === WaitCondition.Scheduled)
    return;
  unsupported(
    CreateFieldProvisioningWaitFor,
    "requested wait condition is unsupported",
  );
}

export function mapSandboxInfo(
  value: BlaxelSandboxIdentity,
  request?: CreateOptions,
): SandboxInfo {
  const mapped = generatedMapSandboxIdentity(value);
  const id = mapped.id ?? request?.name ?? value.metadata?.name ?? "";
  const result: SandboxInfo = {
    id,
    provider: "blaxel",
    state: mapped.state ?? SandboxState.Unknown,
    labels: mapped.labels ?? {},
  };
  if (mapped.name !== undefined) mutable(result).name = mapped.name;
  else if (request?.name !== undefined) mutable(result).name = request.name;
  if (request?.labels !== undefined && mapped.labels === undefined)
    mutable(result).labels = { ...request.labels };
  if (mapped.region !== undefined) mutable(result).region = mapped.region;
  if (mapped.providerState !== undefined)
    mutable(result).providerState = mapped.providerState;
  if (result.state === SandboxState.Deleted) mutable(result).terminal = true;
  const resources = mapAllocatedResources(value);
  if (resources !== undefined) mutable(result).resources = resources;
  return result;
}

function mutable(v: SandboxInfo): {
  -readonly [K in keyof SandboxInfo]: SandboxInfo[K];
} {
  return v;
}

function nonblank(v: string, f: string): string {
  if (!v.trim()) invalid(f, "must not be blank");
  return v;
}

interface DurationStringArgs {
  ms: number;
  field: string;
  minimumMs: number;
}

/** Converts Kit millisecond durations into Blaxel duration strings (s/m/h/d). */
function durationString(args: DurationStringArgs): string {
  const { ms, field, minimumMs } = args;
  if (!Number.isFinite(ms) || ms <= 0 || ms % 1_000 !== 0)
    invalid(field, "must be a whole number of seconds");
  if (ms < minimumMs) invalid(field, "is below the Blaxel minimum duration");
  const seconds = ms / 1_000;
  if (seconds % 86_400 === 0) return `${seconds / 86_400}d`;
  if (seconds % 3_600 === 0) return `${seconds / 3_600}h`;
  if (seconds % 60 === 0) return `${seconds / 60}m`;
  return `${seconds}s`;
}

function rejectIfSet(set: boolean, field: string, message: string): void {
  if (set) unsupported(field, message);
}

function isDefined<T>(value: T | undefined): value is T {
  return value !== undefined;
}

function hasUnsupportedRuntimeFields(
  runtime: NonNullable<CreateOptions["runtime"]>,
): boolean {
  return [
    runtime.entrypoint,
    runtime.pty,
    runtime.user,
    runtime.language,
    runtime.version,
  ].some(isDefined);
}

function hasUnsupportedIsolation(
  isolation: NonNullable<CreateOptions["isolation"]>,
): boolean {
  return (
    (isolation.kind !== undefined &&
      isolation.kind !== IsolationKind.Default) ||
    isolation.providerKind !== undefined ||
    isolation.nestedVirtualization !== undefined
  );
}

function hasUnsupportedLifetimeFields(
  lifetime: NonNullable<CreateOptions["lifetime"]>,
): boolean {
  return [
    lifetime.ephemeral,
    lifetime.idleStop,
    lifetime.idlePause,
    lifetime.stoppedArchive,
    lifetime.stoppedDelete,
  ].some(isDefined);
}

function hasUnsupportedNetwork(
  network: NonNullable<CreateOptions["network"]>,
): boolean {
  const nonDefaultEgress =
    network.egress !== undefined && network.egress !== EgressMode.Default;
  return (
    nonDefaultEgress ||
    [
      network.outboundCidrs,
      network.inboundCidrs,
      network.privateNetwork,
      network.linkedSandbox,
      network.proxyReference,
      network.outboundProxyUrl,
      network.customDomain,
    ].some(isDefined)
  );
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

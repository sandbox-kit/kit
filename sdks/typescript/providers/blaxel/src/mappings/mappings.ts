// Provider-specific orchestration around generated portable field mappings.
import type {
  Config as NativeClientConfig,
  Port,
  SandboxCreateConfiguration,
  SandboxCreateOptions,
  SandboxLifecycle,
} from "@blaxel/core";
import {
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
  readonly options?: SandboxCreateOptions;
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
  if (
    v.entrypoint !== undefined ||
    v.pty !== undefined ||
    v.user !== undefined ||
    v.language !== undefined ||
    v.version !== undefined
  )
    unsupported(
      CreateFieldRuntime,
      "runtime user, language, version, entrypoint, and PTY are not mapped by Blaxel",
    );
}

function mapIsolation(r: CreateOptions): void {
  const v = r.isolation;
  if (!v) return;
  if (
    (v.kind !== undefined && v.kind !== IsolationKind.Default) ||
    v.providerKind !== undefined ||
    v.nestedVirtualization !== undefined
  )
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
  if (
    v.egress !== undefined ||
    v.outboundCidrs !== undefined ||
    v.inboundCidrs !== undefined ||
    v.privateNetwork !== undefined ||
    v.linkedSandbox !== undefined ||
    v.proxyReference !== undefined ||
    v.outboundProxyUrl !== undefined ||
    v.customDomain !== undefined
  )
    unsupported(
      CreateFieldNetwork,
      "requested network settings are unsupported by Blaxel",
    );

  if (v.outboundDomains !== undefined) {
    p.network = {
      ...(p.network ?? {}),
      allowedDomains: [...v.outboundDomains],
    };
  }

  if (v.ports) {
    const ports: Port[] = [];
    for (const x of v.ports) {
      const port = x.port;
      if (port === undefined)
        invalid(CreateFieldNetworkPortsPort, "port is required");
      if (port === 80 || port === 443 || port === 8080)
        invalid(
          CreateFieldNetworkPortsPort,
          "ports 80, 443, and 8080 are reserved by Blaxel",
        );
      const protocol = blaxelPortProtocol(x.transport ?? PortTransport.Default);
      ports.push(
        protocol === undefined ? { target: port } : { target: port, protocol },
      );
    }
    p.ports = ports;
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
  // Blaxel exposes HTTP (not H2). Kit HTTP/2+TLS maps to that HTTP protocol.
  if (transport === PortTransport.Http2Tls) return "HTTP";
  invalid(
    CreateFieldNetworkPortsTransport,
    "requested port transport is unsupported",
  );
}

function mapSecurity(r: CreateOptions): void {
  if (r.security !== undefined)
    unsupported(CreateFieldSecurity, "Blaxel security settings are not mapped");
}

function mapLifetime(r: CreateOptions, p: SandboxCreateConfiguration): void {
  const v = r.lifetime;
  if (!v) return;
  if (
    v.ephemeral !== undefined ||
    v.idleStop !== undefined ||
    v.idlePause !== undefined ||
    v.stoppedArchive !== undefined ||
    v.stoppedDelete !== undefined
  )
    unsupported(
      CreateFieldLifetime,
      "requested lifecycle policy is unsupported by Blaxel",
    );

  if (v.maximumLifetime !== undefined) {
    p.ttl = durationString(
      v.maximumLifetime,
      CreateFieldLifetimeMaximumLifetime,
      5 * 60_000,
    );
  }

  if (v.idleTerminate) {
    const mode = v.idleTerminate.mode ?? PolicyMode.Default;
    if (mode === PolicyMode.Disabled) {
      unsupported(
        CreateFieldLifetimeIdleTerminate,
        "disabling idle terminate is not representable",
      );
    } else if (mode === PolicyMode.After) {
      if (v.idleTerminate.after === undefined)
        invalid(CreateFieldLifetimeIdleTerminateAfter, "duration is required");
      const value = durationString(
        v.idleTerminate.after,
        CreateFieldLifetimeIdleTerminateAfter,
        60_000,
      );
      const lifecycle: SandboxLifecycle = {
        expirationPolicies: [{ type: "ttl-idle", action: "delete", value }],
      };
      p.lifecycle = lifecycle;
    } else if (mode !== PolicyMode.Default)
      invalid(CreateFieldLifetimeIdleTerminateMode, "unknown policy mode");
  }
}

function mapReadiness(r: CreateOptions): void {
  if (r.readiness !== undefined)
    unsupported(CreateFieldReadiness, "Blaxel readiness probes are not mapped");
}

function mapObservability(r: CreateOptions): void {
  if (r.observability !== undefined)
    unsupported(
      CreateFieldObservability,
      "Blaxel observability settings are not mapped",
    );
}

// Kit Client.create applies provisioning.timeout as the abort deadline.
// @blaxel/core also accepts SandboxCreateOptions.timeout (1–50s); this
// provider does not forward it — use the Kit deadline instead.
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

/** Converts Kit millisecond durations into Blaxel duration strings (s/m/h/d). */
function durationString(ms: number, field: string, minimumMs: number): string {
  if (!Number.isFinite(ms) || ms <= 0 || ms % 1_000 !== 0)
    invalid(field, "must be a whole number of seconds");
  if (ms < minimumMs) invalid(field, "is below the Blaxel minimum duration");
  const seconds = ms / 1_000;
  if (seconds % 86_400 === 0) return `${seconds / 86_400}d`;
  if (seconds % 3_600 === 0) return `${seconds / 3_600}h`;
  if (seconds % 60 === 0) return `${seconds / 60}m`;
  return `${seconds}s`;
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

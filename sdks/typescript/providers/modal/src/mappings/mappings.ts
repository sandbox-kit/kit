// Provider-specific orchestration around generated portable field mappings.
import {
  Probe,
  type ModalClientParams,
  type Sandbox,
  type SandboxCreateParams,
} from "modal";
import {
  EgressMode,
  ErrorKind,
  IsolationKind,
  PolicyMode,
  PortTransport,
  SecretInjection,
  SandboxState,
  WaitCondition,
  type Config,
  type CreateOptions,
  type SandboxInfo,
  CreateFieldIsolation,
  CreateFieldIsolationKind,
  CreateFieldLifetime,
  CreateFieldLifetimeIdleTerminateAfter,
  CreateFieldLifetimeIdleTerminateMode,
  CreateFieldNetwork,
  CreateFieldNetworkEgress,
  CreateFieldNetworkPortsPort,
  CreateFieldNetworkPortsTransport,
  CreateFieldPlacementTarget,
  CreateFieldProvisioningTimeout,
  CreateFieldProvisioningWaitFor,
  CreateFieldReadiness,
  CreateFieldRuntime,
  CreateFieldSecuritySecretsEnvironmentVariable,
  CreateFieldSecuritySecretsInjection,
  CreateFieldSource,
  ConfigFieldScopeAppName,
  CreateFieldName,
  CreateFieldNetworkInboundCidrs,
  CreateFieldNetworkProxyReference,
  CreateFieldObservabilityTelemetryEndpoint,
  CreateFieldProvisioningQueueTimeout,
  CreateFieldResourcesCpuLimitCores,
  CreateFieldResourcesDiskMib,
  CreateFieldResourcesMemoryLimitMib,
  CreateFieldResourcesMemoryMib,
  CreateFieldSecurityPublicAccess,
  CreateFieldSecuritySecretsReference,
  CreateFieldSourceImageReference,
} from "@sandbox-kit/sandbox";
import {
  mapCommonCreateFields as generatedMapCommonCreateFields,
  mapIdleDuration as generatedMapIdleDuration,
  mapLifetimeDuration as generatedMapLifetimeDuration,
  mapResources as generatedMapResources,
  mapSandboxIdentity as generatedMapSandboxIdentity,
} from "./provider_mappings.js";
import { validateProviderCreate } from "#checks/provider_checks.js";
import { mapGeneratedClientConfig } from "#config/provider_config.js";
import type { ProviderState } from "#config/provider_config.js";
import { providerError } from "#errors.js";

export interface ModalCreatePlan {
  readonly app: string;
  readonly environment?: string;
  readonly image: string;
  readonly params: SandboxCreateParams;
  readonly secretNames: readonly string[];
  readonly proxyName?: string;
  readonly waitUntilReady: boolean;
  readonly readyTimeoutMs?: number;
}

export function mapClientConfig(config: Config): ModalClientParams {
  const result = mapGeneratedClientConfig(config);
  if (config.timeout !== undefined)
    result.timeoutMs = finiteDurationMs(config.timeout, "timeout", true);
  return result;
}

export function mapCreateRequest(
  request: CreateOptions,
  config: ProviderState,
): ModalCreatePlan {
  validateProviderCreate(request);
  const app =
    config.scope?.appName === undefined
      ? invalid(ConfigFieldScopeAppName, "Modal app name is required")
      : nonblank(config.scope.appName, ConfigFieldScopeAppName);
  if (request.source?.image === undefined)
    unsupported(CreateFieldSource, "Modal requires an explicit image source");
  const params: SandboxCreateParams = generatedMapCommonCreateFields(request);
  if (request.name !== undefined)
    params.name = nonblank(request.name, CreateFieldName);
  if (config.region !== undefined)
    params.regions = [nonblankConfig(config.region, "region")];
  mapRuntime(request, params);
  mapIsolation(request, params);
  mapResources(request, params);
  mapPlacement(request, params);
  mapNetwork(request, params);
  const secretNames = mapSecurity(request, params);
  mapLifetime(request, params);
  mapReadiness(request, params);
  const wait = mapProvisioning(request, params);
  mapObservability(request, params);
  return {
    app,
    ...(config.scope?.environment === undefined
      ? {}
      : { environment: config.scope.environment }),
    image: nonblank(
      request.source.image.reference,
      CreateFieldSourceImageReference,
    ),
    params,
    secretNames,
    ...(request.network?.proxyReference === undefined
      ? {}
      : {
          proxyName: nonblank(
            request.network.proxyReference,
            CreateFieldNetworkProxyReference,
          ),
        }),
    ...wait,
  };
}
function mapRuntime(r: CreateOptions, p: SandboxCreateParams): void {
  const v = r.runtime;
  if (!v) return;
  if (v.workingDirectory !== undefined) p.workdir = v.workingDirectory;
  if (v.entrypoint !== undefined) p.command = [...v.entrypoint];
  if (v.pty !== undefined) p.pty = v.pty;
  if (
    v.language !== undefined ||
    v.version !== undefined ||
    v.user !== undefined
  )
    unsupported(
      CreateFieldRuntime,
      "language, version, and user must come from the Modal image",
    );
}
function mapIsolation(r: CreateOptions, p: SandboxCreateParams): void {
  const v = r.isolation;
  if (!v) return;
  const kind = v.kind ?? IsolationKind.Default;
  if (kind === IsolationKind.LinuxVm) p.runtime = "vm";
  else if (kind === IsolationKind.Container) p.runtime = "gvisor";
  else if (kind !== IsolationKind.Default)
    unsupported(
      CreateFieldIsolationKind,
      "requested Modal isolation is unsupported",
    );
  if (v.providerKind !== undefined || v.nestedVirtualization !== undefined)
    unsupported(
      CreateFieldIsolation,
      "provider kind and nested virtualization are unsupported",
    );
}
function mapResources(r: CreateOptions, p: SandboxCreateParams): void {
  const v = r.resources;
  if (!v) return;
  if (v.cpuLimitCores !== undefined)
    positive(v.cpuLimitCores, CreateFieldResourcesCpuLimitCores);
  if (v.memoryMib !== undefined)
    integer(v.memoryMib, CreateFieldResourcesMemoryMib, 128);
  if (v.memoryLimitMib !== undefined)
    integer(v.memoryLimitMib, CreateFieldResourcesMemoryLimitMib, 1);
  if (v.diskMib !== undefined)
    unsupported(
      CreateFieldResourcesDiskMib,
      "Modal disk override is unsupported",
    );
  Object.assign(p, generatedMapResources(v));
}
function mapPlacement(r: CreateOptions, p: SandboxCreateParams): void {
  const v = r.placement;
  if (!v) return;
  if (v.cloud !== undefined) p.cloud = v.cloud;
  if (v.regions !== undefined) p.regions = [...v.regions];
  if (v.target !== undefined)
    unsupported(
      CreateFieldPlacementTarget,
      "Modal placement target is unsupported",
    );
}
function mapNetwork(r: CreateOptions, p: SandboxCreateParams): void {
  const v = r.network;
  if (!v) return;
  const egress = v.egress ?? EgressMode.Default;
  if (egress === EgressMode.BlockAll) p.blockNetwork = true;
  else if (
    egress !== EgressMode.Default &&
    egress !== EgressMode.AllowAll &&
    egress !== EgressMode.Restricted
  )
    invalid(CreateFieldNetworkEgress, "unknown egress mode");
  if (v.outboundCidrs !== undefined)
    p.outboundCidrAllowlist = [...v.outboundCidrs];
  if (v.outboundDomains !== undefined)
    p.outboundDomainAllowlist = [...v.outboundDomains];
  if (v.inboundCidrs !== undefined) {
    if (!v.inboundCidrs.length)
      unsupported(
        CreateFieldNetworkInboundCidrs,
        "empty inbound deny-all is not representable",
      );
    p.inboundCidrAllowlist = [...v.inboundCidrs];
  }
  if (v.privateNetwork !== undefined) p.i6pn = v.privateNetwork;
  if (v.customDomain !== undefined) p.customDomain = v.customDomain;
  if (v.ports) {
    p.encryptedPorts = [];
    p.h2Ports = [];
    p.unencryptedPorts = [];
    for (const x of v.ports) {
      const port = x.port;
      if (port === undefined)
        invalid(CreateFieldNetworkPortsPort, "port is required");
      const transport = x.transport ?? PortTransport.Default;
      if (transport === PortTransport.Tls) p.encryptedPorts.push(port);
      else if (transport === PortTransport.Http2Tls) p.h2Ports.push(port);
      else if (transport === PortTransport.Plain) p.unencryptedPorts.push(port);
      else
        invalid(
          CreateFieldNetworkPortsTransport,
          "explicit transport is required",
        );
    }
  }
  if (v.linkedSandbox !== undefined || v.outboundProxyUrl !== undefined)
    unsupported(
      CreateFieldNetwork,
      "linked sandbox and outbound proxy URL mappings are unsupported",
    );
}
function mapSecurity(
  r: CreateOptions,
  p: SandboxCreateParams,
): readonly string[] {
  const v = r.security;
  if (!v) return [];
  if (v.workloadIdentity !== undefined)
    p.includeOidcIdentityToken = v.workloadIdentity;
  if (v.publicAccess !== undefined)
    unsupported(
      CreateFieldSecurityPublicAccess,
      "Modal public access is represented by ports",
    );
  const names: string[] = [];
  for (const s of v.secrets ?? []) {
    if (s.environmentVariable !== undefined)
      unsupported(
        CreateFieldSecuritySecretsEnvironmentVariable,
        "per-variable secret renaming is unsupported",
      );
    const injection = s.injection ?? SecretInjection.Default;
    if (
      injection !== SecretInjection.Default &&
      injection !== SecretInjection.Environment
    )
      unsupported(
        CreateFieldSecuritySecretsInjection,
        "egress-placeholder injection is unsupported",
      );
    names.push(nonblank(s.reference, CreateFieldSecuritySecretsReference));
  }
  return names;
}
function mapLifetime(r: CreateOptions, p: SandboxCreateParams): void {
  const v = r.lifetime;
  if (!v) return;
  if (v.maximumLifetime !== undefined)
    finiteDurationMs(v.maximumLifetime, "lifetime.maximumLifetime", false);
  Object.assign(p, generatedMapLifetimeDuration(v));
  if (v.idleTerminate) {
    const mode = v.idleTerminate.mode ?? PolicyMode.Default;
    if (mode === PolicyMode.Disabled) p.idleTimeoutMs = 0;
    else if (mode === PolicyMode.After) {
      if (v.idleTerminate.after === undefined)
        invalid(CreateFieldLifetimeIdleTerminateAfter, "duration is required");
      finiteDurationMs(
        v.idleTerminate.after,
        "lifetime.idleTerminate.after",
        false,
      );
      Object.assign(p, generatedMapIdleDuration(v.idleTerminate));
    } else if (mode !== PolicyMode.Default)
      invalid(CreateFieldLifetimeIdleTerminateMode, "unknown policy mode");
  }
  if (
    v.ephemeral !== undefined ||
    v.idleStop !== undefined ||
    v.idlePause !== undefined ||
    v.stoppedArchive !== undefined ||
    v.stoppedDelete !== undefined
  )
    unsupported(
      CreateFieldLifetime,
      "requested lifecycle policy is unsupported by Modal",
    );
}
function mapReadiness(r: CreateOptions, p: SandboxCreateParams): void {
  const v = r.readiness;
  if (!v) return;
  if (v.tcpPort !== undefined) {
    p.readinessProbe = Probe.withTcp(v.tcpPort);
  } else if (v.command !== undefined)
    p.readinessProbe = Probe.withExec([...(v.command.argv ?? [])]);
  else invalid(CreateFieldReadiness, "probe is required");
}
function mapProvisioning(
  r: CreateOptions,
  _p: SandboxCreateParams,
): { waitUntilReady: boolean; readyTimeoutMs?: number } {
  const v = r.provisioning;
  if (!v) return { waitUntilReady: false };
  if (v.queueTimeout !== undefined)
    unsupported(
      CreateFieldProvisioningQueueTimeout,
      "Modal queue timeout is unsupported",
    );
  const waitFor = v.waitFor ?? WaitCondition.Default;
  if (waitFor === WaitCondition.Default || waitFor === WaitCondition.Scheduled)
    return { waitUntilReady: false };
  if (waitFor === WaitCondition.Ready) {
    if (v.timeout === undefined || v.timeout <= 0)
      invalid(
        CreateFieldProvisioningTimeout,
        "ready wait requires a positive timeout",
      );
    if (!r.readiness)
      invalid(CreateFieldReadiness, "ready wait requires a probe");
    return { waitUntilReady: true, readyTimeoutMs: v.timeout };
  }
  unsupported(
    CreateFieldProvisioningWaitFor,
    "requested wait condition is unsupported",
  );
}
function mapObservability(r: CreateOptions, p: SandboxCreateParams): void {
  const v = r.observability;
  if (!v) return;
  if (v.verbose !== undefined) p.verbose = v.verbose;
  if (v.telemetryEndpoint !== undefined)
    unsupported(
      CreateFieldObservabilityTelemetryEndpoint,
      "per-sandbox telemetry endpoint is unsupported",
    );
}

export function mapSandboxInfo(
  value: Sandbox,
  request?: CreateOptions,
): SandboxInfo {
  const mapped = generatedMapSandboxIdentity(value);
  const result: SandboxInfo = {
    id: mapped.id ?? value.sandboxId,
    provider: "modal",
    state: SandboxState.Running,
    labels: {},
  };
  if (request?.name !== undefined) {
    mutable(result).name = request.name;
  }
  if (request?.labels !== undefined) {
    mutable(result).labels = { ...request.labels };
  }
  if (request?.isolation !== undefined) {
    mutable(result).isolation = { ...request.isolation };
  }
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
function nonblankConfig(v: string, f: string): string {
  if (!v.trim()) configInvalid(f, "must not be blank");
  return v;
}
function positive(v: number, f: string): number {
  if (!Number.isFinite(v) || v <= 0) invalid(f, "must be positive");
  return v;
}
function integer(v: bigint, f: string, min: number): number {
  if (v < BigInt(min) || v > 2_147_483_647n)
    invalid(f, "is outside the supported range");
  return Number(v);
}
function finiteDurationMs(v: number, f: string, zero: boolean): number {
  if (!Number.isFinite(v) || v < 0 || (!zero && v === 0))
    invalid(f, "must be a finite positive duration");
  return v;
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

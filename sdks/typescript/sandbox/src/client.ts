// Client and sandbox runtime behavior is handwritten. Portable client data
// contracts come from standard protobuf generation.

import { create, type MessageInitShape } from "@bufbuild/protobuf";
import { durationFromMs, durationMs, type Duration } from "@bufbuild/protobuf/wkt";
import {
  ConfigSchema,
  type APIKeyCredentials,
  type ArchiveSandboxRequest,
  type AttachSandboxRequest,
  type AuthConfig,
  type BearerTokenCredentials,
  type Config,
  type DeleteSandboxRequest,
  type GetSandboxRequest,
  type OAuthCredentials,
  type PauseSandboxRequest,
  type RefreshSandboxRequest,
  type SandboxReference,
  type Scope,
  type StartSandboxRequest,
  type StopSandboxRequest,
  type TerminateSandboxRequest,
  type TokenPairCredentials,
} from "@sandbox-kit/contracts/kit/sandbox/v1/client_pb.js";
import { ErrorKind as ProtoErrorKind } from "@sandbox-kit/contracts/kit/sandbox/v1/sandbox_pb.js";
import { EgressModeNumber, IsolationKindNumber, PolicyModeNumber, PortTransportNumber, SandboxState, SecretInjectionNumber, SnapshotKindNumber, WaitConditionNumber, type ArchiveOptions, type AutomaticAction, type CommandProbe, type CreateOptions, type CreateResult, type DeleteOptions, type ExitResult, type ImageSource, type IsolationConfig, type LifecycleCapability, type LifetimePolicy, type NetworkConfig, type ObservabilityConfig, type PauseOptions, type Placement, type PortBinding, type ProviderDefaultSource, type ProvisioningOptions, type ReadinessProbe, type Resources, type RuntimeConfig, type SandboxInfo, type SandboxSource, type SecretReference, type SecurityConfig, type SnapshotSource, type StartOptions, type StopOptions, type StringAllowlist, type TerminateOptions, type WarmPoolSource } from "./sandbox.gen.js";
import type { ErrorInfo } from "./sandbox.gen.js";

export type {
  APIKeyCredentials,
  ArchiveSandboxRequest,
  AttachSandboxRequest,
  AuthConfig,
  BearerTokenCredentials,
  Config,
  DeleteSandboxRequest,
  GetSandboxRequest,
  OAuthCredentials,
  PauseSandboxRequest,
  RefreshSandboxRequest,
  SandboxReference,
  Scope,
  StartSandboxRequest,
  StopSandboxRequest,
  TerminateSandboxRequest,
  TokenPairCredentials,
};

/** Plain-object client configuration accepted at construction time. */
export type ConfigInput = MessageInitShape<typeof ConfigSchema>;

// These ergonomic names reference the protobuf-generated enum values; they do
// not define a second error classification contract.
export const ErrorKind = {
  Unknown: ProtoErrorKind.UNKNOWN,
  InvalidArgument: ProtoErrorKind.INVALID_ARGUMENT,
  Unsupported: ProtoErrorKind.UNSUPPORTED,
  Authentication: ProtoErrorKind.AUTHENTICATION,
  PermissionDenied: ProtoErrorKind.PERMISSION_DENIED,
  NotFound: ProtoErrorKind.NOT_FOUND,
  AlreadyExists: ProtoErrorKind.ALREADY_EXISTS,
  Conflict: ProtoErrorKind.CONFLICT,
  RateLimited: ProtoErrorKind.RATE_LIMITED,
  ResourceExhausted: ProtoErrorKind.RESOURCE_EXHAUSTED,
  Timeout: ProtoErrorKind.TIMEOUT,
  Canceled: ProtoErrorKind.CANCELED,
  Unavailable: ProtoErrorKind.UNAVAILABLE,
  InvalidResponse: ProtoErrorKind.INVALID_RESPONSE,
  Internal: ProtoErrorKind.INTERNAL,
  FailedPrecondition: ProtoErrorKind.FAILED_PRECONDITION,
} as const;

export type ErrorKind = (typeof ErrorKind)[keyof typeof ErrorKind];

/** Builds a protobuf duration for client configuration. */
export function duration(ms: number): Duration {
  return durationFromMs(ms);
}

export interface OperationOptions {
  readonly signal?: AbortSignal;
  readonly timeoutMs?: number;
}

export interface Backend {
	readonly name: string;
  create?(request: CreateOptions, options: OperationOptions): Promise<CreateResult>;
  attach?(id: string, options: OperationOptions): Promise<SandboxInfo>;
  refresh?(id: string, options: OperationOptions): Promise<SandboxInfo>;
  start?(id: string, request: StartOptions, options: OperationOptions): Promise<SandboxInfo>;
  stop?(id: string, request: StopOptions, options: OperationOptions): Promise<SandboxInfo>;
  pause?(id: string, request: PauseOptions, options: OperationOptions): Promise<SandboxInfo>;
  terminate?(id: string, request: TerminateOptions, options: OperationOptions): Promise<ExitResult>;
  archive?(id: string, request: ArchiveOptions, options: OperationOptions): Promise<SandboxInfo>;
  delete?(id: string, request: DeleteOptions, options: OperationOptions): Promise<SandboxInfo>;
  close(options: OperationOptions): Promise<void>;
  closeSandbox?(id: string, options: OperationOptions): Promise<void>;
  capabilities?(info: SandboxInfo): readonly LifecycleCapability[];
}

export interface Provider {
  readonly name: string;
  initialize(config: Config): Backend;
}

/** Provider selection stays handwritten and is never part of the portable Config message. */
export type ClientConfig = ConfigInput & {
  readonly provider: Provider;
};

export class SandboxError extends Error {
  readonly kind: ErrorInfo["kind"];
  readonly provider: string;
  readonly operation: string;
  readonly field: string;
  readonly rule: string | undefined;
  readonly statusCode: number | undefined;
  readonly providerCode: string | undefined;
  readonly providerSource: string | undefined;

  constructor(info: ErrorInfo, options?: ErrorOptions) {
    super(info.message, options);
    this.name = "SandboxError";
    this.kind = info.kind;
    this.provider = info.provider;
    this.operation = info.operation;
    this.field = info.field;
    this.rule = info.rule;
    this.statusCode = info.statusCode;
    this.providerCode = info.providerCode;
    this.providerSource = info.providerSource;
  }
}

export function sandboxError(info: ErrorInfo, cause?: unknown): SandboxError {
  return new SandboxError(info, cause === undefined ? undefined : { cause });
}
function validationError(provider: string, operation: string, field: string, rule: string, message: string): never { throw sandboxError({ kind: ErrorKind.InvalidArgument, provider, operation, field, rule, message: "sandbox-kit: " + message }); }
function joinValidationPath(parent: string, field: string): string { return parent ? parent + "." + field : field; }
function validSecureEndpoint(value: string): boolean { try { const url = new URL(value); if (!url.hostname || url.username || url.password || url.hash) return false; if (url.protocol === "https:") return true; if (url.protocol !== "http:") return false; const host = url.hostname.toLowerCase(); return host === "localhost" || host === "[::1]" || host === "::1" || /^127(?:\.[0-9]{1,3}){3}$/.test(host); } catch { return false; } }

function validateAutomaticActionValue(value: AutomaticAction | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.after !== undefined && !Number.isFinite(value.after)) validationError(provider, operation, joinValidationPath(path, "after"), "finite", "must be finite");
  if (value.after !== undefined && !(value.after >= 0)) validationError(provider, operation, joinValidationPath(path, "after"), "minimum", "must be at least 0");
  if (value.mode !== undefined && !(PolicyModeNumber[value.mode] >= 0)) validationError(provider, operation, joinValidationPath(path, "mode"), "minimum", "must be at least 0");
  if (value.mode !== undefined && !(PolicyModeNumber[value.mode] <= 2)) validationError(provider, operation, joinValidationPath(path, "mode"), "maximum", "must be at most 2");
  if (value.mode === "after") { const count0 = [value.after !== undefined].filter(Boolean).length;
    if (count0 === 0) validationError(provider, operation, joinValidationPath(path, "after"), "requires_any", "requires a value"); }
  if (value.mode !== "after") { const count1 = [value.after !== undefined].filter(Boolean).length;
    if (count1 > 0) validationError(provider, operation, joinValidationPath(path, "after"), "forbids", "forbids these values"); }
  if (value.mode !== undefined && !Object.hasOwn(PolicyModeNumber, value.mode)) validationError(provider, operation, joinValidationPath(path, "mode"), "known_enum", "must be a known PolicyMode value");
}

function validateCommandProbeValue(value: CommandProbe | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.argv !== undefined && value.argv.length < 1) validationError(provider, operation, joinValidationPath(path, "argv"), "min_items", "must contain at least 1 item(s)");
}

function validateCreateOptionsValue(value: CreateOptions | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  validateSandboxSourceValue(value.source, provider, operation, joinValidationPath(path, "source"));
  validateRuntimeConfigValue(value.runtime, provider, operation, joinValidationPath(path, "runtime"));
  validateIsolationConfigValue(value.isolation, provider, operation, joinValidationPath(path, "isolation"));
  validateResourcesValue(value.resources, provider, operation, joinValidationPath(path, "resources"));
  validatePlacementValue(value.placement, provider, operation, joinValidationPath(path, "placement"));
  validateNetworkConfigValue(value.network, provider, operation, joinValidationPath(path, "network"));
  validateSecurityConfigValue(value.security, provider, operation, joinValidationPath(path, "security"));
  validateLifetimePolicyValue(value.lifetime, provider, operation, joinValidationPath(path, "lifetime"));
  validateProvisioningOptionsValue(value.provisioning, provider, operation, joinValidationPath(path, "provisioning"));
  validateReadinessProbeValue(value.readiness, provider, operation, joinValidationPath(path, "readiness"));
  validateObservabilityConfigValue(value.observability, provider, operation, joinValidationPath(path, "observability"));
}

function validateImageSourceValue(value: ImageSource | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.reference !== undefined && !value.reference.trim()) validationError(provider, operation, joinValidationPath(path, "reference"), "nonblank", "must not be blank");
  if (value.reference !== undefined && (value.reference.length > 512 || !/^(?:[A-Za-z0-9](?:[A-Za-z0-9.-]*[A-Za-z0-9])?(?::[0-9]+)?\/)?[a-z0-9]+(?:(?:[._]|__|-+)[a-z0-9]+)*(?:\/[a-z0-9]+(?:(?:[._]|__|-+)[a-z0-9]+)*)*(?::[A-Za-z0-9_][A-Za-z0-9_.-]{0,127})?(?:@sha256:[a-fA-F0-9]{64})?$/.test(value.reference))) validationError(provider, operation, joinValidationPath(path, "reference"), "image_reference", "must be a valid image reference");
}

function validateIsolationConfigValue(value: IsolationConfig | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.kind !== undefined && !Object.hasOwn(IsolationKindNumber, value.kind)) validationError(provider, operation, joinValidationPath(path, "kind"), "known_enum", "must be a known IsolationKind value");
}

function validateLifetimePolicyValue(value: LifetimePolicy | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.maximumLifetime !== undefined && !Number.isFinite(value.maximumLifetime)) validationError(provider, operation, joinValidationPath(path, "maximumLifetime"), "finite", "must be finite");
  if (value.maximumLifetime !== undefined && !(value.maximumLifetime >= 0)) validationError(provider, operation, joinValidationPath(path, "maximumLifetime"), "minimum", "must be at least 0");
  validateAutomaticActionValue(value.idleTerminate, provider, operation, joinValidationPath(path, "idleTerminate"));
  validateAutomaticActionValue(value.idleStop, provider, operation, joinValidationPath(path, "idleStop"));
  validateAutomaticActionValue(value.idlePause, provider, operation, joinValidationPath(path, "idlePause"));
  validateAutomaticActionValue(value.stoppedArchive, provider, operation, joinValidationPath(path, "stoppedArchive"));
  validateAutomaticActionValue(value.stoppedDelete, provider, operation, joinValidationPath(path, "stoppedDelete"));
}

function validateNetworkConfigValue(value: NetworkConfig | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.linkedSandbox !== undefined && !value.linkedSandbox.trim()) validationError(provider, operation, joinValidationPath(path, "linkedSandbox"), "nonblank", "must not be blank");
  if (value.proxyReference !== undefined && !value.proxyReference.trim()) validationError(provider, operation, joinValidationPath(path, "proxyReference"), "nonblank", "must not be blank");
  if (value.egress === "restricted") { const count0 = [value.outboundCidrs !== undefined, value.outboundDomains !== undefined].filter(Boolean).length;
    if (count0 === 0) validationError(provider, operation, joinValidationPath(path, "outboundCidrs,outboundDomains"), "requires_any", "requires a value"); }
  if (value.egress === "block_all") { const count1 = [value.outboundCidrs !== undefined, value.outboundDomains !== undefined, value.privateNetwork === true, value.linkedSandbox !== undefined, value.inboundCidrs !== undefined, value.ports !== undefined && value.ports.length > 0].filter(Boolean).length;
    if (count1 > 0) validationError(provider, operation, joinValidationPath(path, "outboundCidrs,outboundDomains,privateNetwork,linkedSandbox,inboundCidrs,ports"), "forbids", "forbids these values"); }
  if (value.egress === "allow_all") { const count2 = [value.outboundCidrs !== undefined, value.outboundDomains !== undefined].filter(Boolean).length;
    if (count2 > 0) validationError(provider, operation, joinValidationPath(path, "outboundCidrs,outboundDomains"), "forbids", "forbids these values"); }
  if (value.egress !== undefined && !Object.hasOwn(EgressModeNumber, value.egress)) validationError(provider, operation, joinValidationPath(path, "egress"), "known_enum", "must be a known EgressMode value");
  validateStringAllowlistValue(value.outboundCidrs, provider, operation, joinValidationPath(path, "outboundCidrs"));
  validateStringAllowlistValue(value.outboundDomains, provider, operation, joinValidationPath(path, "outboundDomains"));
  validateStringAllowlistValue(value.inboundCidrs, provider, operation, joinValidationPath(path, "inboundCidrs"));
  for (const item of value.ports ?? []) validatePortBindingValue(item, provider, operation, joinValidationPath(path, "ports") + "[]");
}

function validateObservabilityConfigValue(value: ObservabilityConfig | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.telemetryEndpoint !== undefined && !value.telemetryEndpoint.trim()) validationError(provider, operation, joinValidationPath(path, "telemetryEndpoint"), "nonblank", "must not be blank");
}

function validatePlacementValue(value: Placement | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  for (const item of value.regions ?? []) {
  if (item !== undefined && !item.trim()) validationError(provider, operation, joinValidationPath(path, "regions") + "[]", "nonblank", "must not be blank");
  }
}

function validatePortBindingValue(value: PortBinding | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.port !== undefined && !(value.port >= 1)) validationError(provider, operation, joinValidationPath(path, "port"), "minimum", "must be at least 1");
  if (value.port !== undefined && !(value.port <= 65535)) validationError(provider, operation, joinValidationPath(path, "port"), "maximum", "must be at most 65535");
  if (value.transport !== undefined && !Object.hasOwn(PortTransportNumber, value.transport)) validationError(provider, operation, joinValidationPath(path, "transport"), "known_enum", "must be a known PortTransport value");
}

function validateProviderDefaultSourceValue(value: ProviderDefaultSource | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
}

function validateProvisioningOptionsValue(value: ProvisioningOptions | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.queueTimeout !== undefined && !Number.isFinite(value.queueTimeout)) validationError(provider, operation, joinValidationPath(path, "queueTimeout"), "finite", "must be finite");
  if (value.queueTimeout !== undefined && !(value.queueTimeout >= 0)) validationError(provider, operation, joinValidationPath(path, "queueTimeout"), "minimum", "must be at least 0");
  if (value.timeout !== undefined && !Number.isFinite(value.timeout)) validationError(provider, operation, joinValidationPath(path, "timeout"), "finite", "must be finite");
  if (value.timeout !== undefined && !(value.timeout >= 0)) validationError(provider, operation, joinValidationPath(path, "timeout"), "minimum", "must be at least 0");
  if (value.waitFor !== undefined && !Object.hasOwn(WaitConditionNumber, value.waitFor)) validationError(provider, operation, joinValidationPath(path, "waitFor"), "known_enum", "must be a known WaitCondition value");
}

function validateReadinessProbeValue(value: ReadinessProbe | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.check === "tcpPort") {
  if (value.tcpPort !== undefined && !(value.tcpPort >= 1)) validationError(provider, operation, joinValidationPath(path, "tcpPort"), "minimum", "must be at least 1");
  if (value.tcpPort !== undefined && !(value.tcpPort <= 65535)) validationError(provider, operation, joinValidationPath(path, "tcpPort"), "maximum", "must be at most 65535");
  }
  if (true) { const count0 = [value.check === "tcpPort", value.check === "command"].filter(Boolean).length;
    if (count0 !== 1) validationError(provider, operation, joinValidationPath(path, "tcpPort,command"), "exactly_one", "requires exactly one value"); }
  if (value.check === "command") validateCommandProbeValue(value.command, provider, operation, joinValidationPath(path, "command"));
}

function validateResourcesValue(value: Resources | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.cpuCores !== undefined && !Number.isFinite(value.cpuCores)) validationError(provider, operation, joinValidationPath(path, "cpuCores"), "finite", "must be finite");
  if (value.cpuCores !== undefined && !(value.cpuCores > 0)) validationError(provider, operation, joinValidationPath(path, "cpuCores"), "exclusive_minimum", "must be greater than 0");
  if (value.cpuLimitCores !== undefined && !Number.isFinite(value.cpuLimitCores)) validationError(provider, operation, joinValidationPath(path, "cpuLimitCores"), "finite", "must be finite");
  if (value.cpuLimitCores !== undefined && !(value.cpuLimitCores >= 0)) validationError(provider, operation, joinValidationPath(path, "cpuLimitCores"), "minimum", "must be at least 0");
  if (value.diskMib !== undefined && !(value.diskMib > 0n)) validationError(provider, operation, joinValidationPath(path, "diskMib"), "exclusive_minimum", "must be greater than 0");
  if (value.memoryMib !== undefined && !(value.memoryMib > 0n)) validationError(provider, operation, joinValidationPath(path, "memoryMib"), "exclusive_minimum", "must be greater than 0");
  if (true && value.cpuLimitCores !== undefined && value.cpuCores !== undefined && value.cpuLimitCores > 0 && value.cpuLimitCores < value.cpuCores) validationError(provider, operation, joinValidationPath(path, "cpuLimitCores"), "limit", "must not be smaller than cpuCores");
  if (true && value.memoryLimitMib !== undefined && value.memoryMib !== undefined && value.memoryLimitMib > 0 && value.memoryLimitMib < value.memoryMib) validationError(provider, operation, joinValidationPath(path, "memoryLimitMib"), "limit", "must not be smaller than memoryMib");
  if (true && value.cpuLimitCores !== undefined && value.cpuLimitCores > 0 && !(value.cpuCores !== undefined && value.cpuCores > 0)) validationError(provider, operation, joinValidationPath(path, "cpuLimitCores"), "requires_positive", "requires a positive cpuCores");
  if (true && value.memoryLimitMib !== undefined && value.memoryLimitMib > 0 && !(value.memoryMib !== undefined && value.memoryMib > 0)) validationError(provider, operation, joinValidationPath(path, "memoryLimitMib"), "requires_positive", "requires a positive memoryMib");
}

function validateRuntimeConfigValue(value: RuntimeConfig | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  for (const item of value.entrypoint ?? []) {
  if (item !== undefined && !item.trim()) validationError(provider, operation, joinValidationPath(path, "entrypoint") + "[]", "nonblank", "must not be blank");
  }
  if (value.language !== undefined && !value.language.trim()) validationError(provider, operation, joinValidationPath(path, "language"), "nonblank", "must not be blank");
  if (value.user !== undefined && !value.user.trim()) validationError(provider, operation, joinValidationPath(path, "user"), "nonblank", "must not be blank");
  if (value.version !== undefined && !value.version.trim()) validationError(provider, operation, joinValidationPath(path, "version"), "nonblank", "must not be blank");
  if (value.workingDirectory !== undefined && !value.workingDirectory.trim()) validationError(provider, operation, joinValidationPath(path, "workingDirectory"), "nonblank", "must not be blank");
}

function validateSandboxSourceValue(value: SandboxSource | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (true) { const count0 = [value.kind === "providerDefault", value.kind === "image", value.kind === "snapshot", value.kind === "warmPool"].filter(Boolean).length;
    if (count0 > 1) validationError(provider, operation, joinValidationPath(path, "providerDefault,image,snapshot,warmPool"), "at_most_one", "allows at most one value"); }
  if (value.kind === "providerDefault") validateProviderDefaultSourceValue(value.providerDefault, provider, operation, joinValidationPath(path, "providerDefault"));
  if (value.kind === "image") validateImageSourceValue(value.image, provider, operation, joinValidationPath(path, "image"));
  if (value.kind === "snapshot") validateSnapshotSourceValue(value.snapshot, provider, operation, joinValidationPath(path, "snapshot"));
  if (value.kind === "warmPool") validateWarmPoolSourceValue(value.warmPool, provider, operation, joinValidationPath(path, "warmPool"));
}

function validateSecretReferenceValue(value: SecretReference | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.environmentVariable !== undefined && !value.environmentVariable.trim()) validationError(provider, operation, joinValidationPath(path, "environmentVariable"), "nonblank", "must not be blank");
  if (value.reference !== undefined && !value.reference.trim()) validationError(provider, operation, joinValidationPath(path, "reference"), "nonblank", "must not be blank");
  if (value.injection !== undefined && !Object.hasOwn(SecretInjectionNumber, value.injection)) validationError(provider, operation, joinValidationPath(path, "injection"), "known_enum", "must be a known SecretInjection value");
}

function validateSecurityConfigValue(value: SecurityConfig | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  for (const item of value.secrets ?? []) validateSecretReferenceValue(item, provider, operation, joinValidationPath(path, "secrets") + "[]");
}

function validateSnapshotSourceValue(value: SnapshotSource | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.reference !== undefined && !value.reference.trim()) validationError(provider, operation, joinValidationPath(path, "reference"), "nonblank", "must not be blank");
  if (value.kind !== undefined && !Object.hasOwn(SnapshotKindNumber, value.kind)) validationError(provider, operation, joinValidationPath(path, "kind"), "known_enum", "must be a known SnapshotKind value");
}

function validateStringAllowlistValue(value: StringAllowlist | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
}

function validateWarmPoolSourceValue(value: WarmPoolSource | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.reference !== undefined && !value.reference.trim()) validationError(provider, operation, joinValidationPath(path, "reference"), "nonblank", "must not be blank");
}

export function validateCreateOptions(value: CreateOptions, provider = "", operation = "create"): void { validateCreateOptionsValue(value, provider, operation, ""); }

function validateAPIKeyCredentialsValue(value: APIKeyCredentials | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.key !== undefined && !value.key.trim()) validationError(provider, operation, joinValidationPath(path, "key"), "nonblank", "must not be blank");
}

function validateAuthConfigValue(value: AuthConfig | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  const credentials = value.credentials;
  if (credentials === undefined || credentials.case === undefined) {
    validationError(provider, operation, joinValidationPath(path, "apiKey,tokenPair,bearerToken,oauthRefresh"), "exactly_one", "requires exactly one value");
  }
  if (credentials.case === "apiKey") validateAPIKeyCredentialsValue(credentials.value, provider, operation, joinValidationPath(path, "apiKey"));
  if (credentials.case === "tokenPair") validateTokenPairCredentialsValue(credentials.value, provider, operation, joinValidationPath(path, "tokenPair"));
  if (credentials.case === "bearerToken") validateBearerTokenCredentialsValue(credentials.value, provider, operation, joinValidationPath(path, "bearerToken"));
  if (credentials.case === "oauthRefresh") validateOAuthCredentialsValue(credentials.value, provider, operation, joinValidationPath(path, "oauthRefresh"));
}

function validateBearerTokenCredentialsValue(value: BearerTokenCredentials | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.token !== undefined && !value.token.trim()) validationError(provider, operation, joinValidationPath(path, "token"), "nonblank", "must not be blank");
}

function validateConfigValue(value: Config | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.endpoint !== undefined && !value.endpoint.trim()) validationError(provider, operation, joinValidationPath(path, "endpoint"), "nonblank", "must not be blank");
  if (value.endpoint !== undefined && !validSecureEndpoint(value.endpoint)) validationError(provider, operation, joinValidationPath(path, "endpoint"), "secure_endpoint", "must use HTTPS except for loopback HTTP development endpoints");
  if (value.region !== undefined && !value.region.trim()) validationError(provider, operation, joinValidationPath(path, "region"), "nonblank", "must not be blank");
  if (value.timeout !== undefined) {
    let timeoutMs: number;
    try {
      timeoutMs = durationMs(value.timeout);
    } catch {
      validationError(provider, operation, joinValidationPath(path, "timeout"), "finite", "must be a non-negative protobuf duration");
    }
    if (!Number.isFinite(timeoutMs) || timeoutMs < 0) validationError(provider, operation, joinValidationPath(path, "timeout"), "minimum", "must be at least 0");
  }
  validateAuthConfigValue(value.auth, provider, operation, joinValidationPath(path, "auth"));
  validateScopeValue(value.scope, provider, operation, joinValidationPath(path, "scope"));
}

function validateOAuthCredentialsValue(value: OAuthCredentials | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.clientId !== undefined && !value.clientId.trim()) validationError(provider, operation, joinValidationPath(path, "clientId"), "nonblank", "must not be blank");
  if (value.clientSecret !== undefined && !value.clientSecret.trim()) validationError(provider, operation, joinValidationPath(path, "clientSecret"), "nonblank", "must not be blank");
  if (value.jwtKey !== undefined && !value.jwtKey.trim()) validationError(provider, operation, joinValidationPath(path, "jwtKey"), "nonblank", "must not be blank");
  if (value.refreshToken !== undefined && !value.refreshToken.trim()) validationError(provider, operation, joinValidationPath(path, "refreshToken"), "nonblank", "must not be blank");
  if (true) { const count0 = [value.clientSecret !== undefined, value.jwtKey !== undefined].filter(Boolean).length;
    if (count0 !== 1) validationError(provider, operation, joinValidationPath(path, "clientSecret,jwtKey"), "exactly_one", "requires exactly one value"); }
}

function validateScopeValue(value: Scope | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.appName !== undefined && !value.appName.trim()) validationError(provider, operation, joinValidationPath(path, "appName"), "nonblank", "must not be blank");
  if (value.environment !== undefined && !value.environment.trim()) validationError(provider, operation, joinValidationPath(path, "environment"), "nonblank", "must not be blank");
  if (value.organizationId !== undefined && !value.organizationId.trim()) validationError(provider, operation, joinValidationPath(path, "organizationId"), "nonblank", "must not be blank");
  if (value.projectId !== undefined && !value.projectId.trim()) validationError(provider, operation, joinValidationPath(path, "projectId"), "nonblank", "must not be blank");
}

function validateTokenPairCredentialsValue(value: TokenPairCredentials | undefined, provider: string, operation: string, path: string): void {
  if (value === undefined) return;
  if (value.id !== undefined && !value.id.trim()) validationError(provider, operation, joinValidationPath(path, "id"), "nonblank", "must not be blank");
  if (value.secret !== undefined && !value.secret.trim()) validationError(provider, operation, joinValidationPath(path, "secret"), "nonblank", "must not be blank");
}

export function validateConfig(value: Config | ConfigInput, provider = "", operation = "create"): void { validateConfigValue(value as Config, provider, operation, ""); }

function clone<T>(value: T): T {
  return structuredClone(value);
}

function operationSignal(options: OperationOptions, inheritedTimeoutMs?: number): { signal?: AbortSignal; dispose(): void } {
  const timeoutMs = options.timeoutMs ?? inheritedTimeoutMs;
  if (timeoutMs !== undefined && (!Number.isFinite(timeoutMs) || timeoutMs < 0)) {
    throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: "", operation: "", field: "timeoutMs", message: "sandbox-kit: timeoutMs must be a finite non-negative number" });
  }
  if ((timeoutMs === undefined || timeoutMs === 0) && options.signal === undefined) return { dispose() {} };
  const controller = new AbortController();
  const abort = () => controller.abort(options.signal?.reason);
  if (options.signal?.aborted) abort();
  else options.signal?.addEventListener("abort", abort, { once: true });
  const timer = timeoutMs !== undefined && timeoutMs > 0
    ? setTimeout(() => controller.abort(new DOMException("The operation timed out", "TimeoutError")), timeoutMs)
    : undefined;
  return {
    signal: controller.signal,
    dispose() {
      if (timer !== undefined) clearTimeout(timer);
      options.signal?.removeEventListener("abort", abort);
    },
  };
}

function operationOptions(options: OperationOptions, signal: AbortSignal | undefined): OperationOptions {
  return {
    ...(options.timeoutMs === undefined ? {} : { timeoutMs: options.timeoutMs }),
    ...(signal === undefined ? {} : { signal }),
  };
}

function abortReason(signal: AbortSignal): DOMException {
  const reason = signal.reason;
  return reason instanceof DOMException && (reason.name === "AbortError" || reason.name === "TimeoutError")
    ? reason
    : new DOMException("The operation was aborted", "AbortError");
}

async function awaitWithSignal<T>(operation: () => Promise<T>, signal?: AbortSignal): Promise<T> {
  if (signal?.aborted) throw abortReason(signal);
  const pending = operation();
  if (signal === undefined) return pending;
  return new Promise<T>((resolve, reject) => {
    const aborted = () => reject(abortReason(signal));
    const cleanup = () => signal.removeEventListener("abort", aborted);
    signal.addEventListener("abort", aborted, { once: true });
    pending.then(
      (value) => { cleanup(); resolve(value); },
      (error: unknown) => { cleanup(); reject(error); },
    );
  });
}
interface ClientLifecycle { begin(operation: string): () => void; isClosed(): boolean; }
const sandboxConstructorToken: unique symbol = Symbol("sandbox-kit.Sandbox");
export class Sandbox implements AsyncDisposable {
  readonly #backend: Backend; readonly #owner: ClientLifecycle; readonly #timeoutMs: number | undefined; #info: SandboxInfo; #closed = false; #tail: Promise<void> = Promise.resolve();
  constructor(token: typeof sandboxConstructorToken, backend: Backend, owner: ClientLifecycle, info: SandboxInfo, timeoutMs?: number) { if (token !== sandboxConstructorToken) throw new TypeError("Sandbox handles are created by Client"); this.#backend = backend; this.#owner = owner; this.#info = clone(info); this.#timeoutMs = timeoutMs; }
  get id(): string { return this.#info.id; } get providerName(): string { return this.#info.provider; } info(): SandboxInfo { return clone(this.#info); }
  capabilities(): readonly LifecycleCapability[] { if (this.#closed || this.#owner.isClosed() || this.#info.terminal) return []; return this.#backend.capabilities?.(clone(this.#info)) ?? []; }
  supports(capability: LifecycleCapability): boolean { return this.capabilities().includes(capability); }
  async refresh(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#serialized("refresh", options, this.#timeoutMs, async (signal) => {
      const method = this.#backend.refresh;
      if (!method) unsupported(this.providerName, "refresh");
      const info = await awaitWithSignal(() => method.call(this.#backend, this.id, operationOptions(options, signal)), signal);
      this.#info = clone({ ...info, provider: this.providerName });
      return this.info();
    });
  }
  async start(request: StartOptions = {}, options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update("start", this.#backend.start, request, options);
  }
  async stop(request: StopOptions = {}, options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update("stop", this.#backend.stop, request, options);
  }
  async pause(request: PauseOptions = {}, options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update("pause", this.#backend.pause, request, options);
  }
  async terminate(request: TerminateOptions = {}, options: OperationOptions = {}): Promise<ExitResult> {
    return this.#serialized("terminate", options, request.timeout ?? this.#timeoutMs, async (signal) => {
      const method = this.#backend.terminate;
      if (!method) unsupported(this.providerName, "terminate");
      const result = await awaitWithSignal(() => method.call(this.#backend, this.id, clone(request), operationOptions(options, signal)), signal);
      this.#info = { ...this.#info, state: SandboxState.Terminated, terminal: true };
      return result;
    });
  }
  async archive(request: ArchiveOptions = {}, options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update("archive", this.#backend.archive, request, options);
  }
  async delete(request: DeleteOptions = {}, options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update("delete", this.#backend.delete, request, options);
  }
  async close(options: OperationOptions = {}): Promise<void> { if (this.#closed) return; await this.#serialized("close", options, this.#timeoutMs, async (signal) => { const close = this.#backend.closeSandbox; if (close) await awaitWithSignal(() => close.call(this.#backend, this.id, operationOptions(options, signal)), signal); this.#closed = true; }); }
  async [Symbol.asyncDispose](): Promise<void> { await this.close(); }
  async #update<T extends StartOptions | StopOptions | PauseOptions | ArchiveOptions | DeleteOptions>(
    operation: string,
    method: ((id: string, request: T, options: OperationOptions) => Promise<SandboxInfo>) | undefined,
    request: T,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    return this.#serialized(operation, options, request.timeout ?? this.#timeoutMs, async (signal) => {
      if (!method) unsupported(this.providerName, operation);
      const info = await awaitWithSignal(
        () => method.call(this.#backend, this.id, clone(request), operationOptions(options, signal)),
        signal,
      );
      this.#info = clone({ ...info, provider: this.providerName });
      return this.info();
    });
  }
  async #serialized<T>(operation: string, options: OperationOptions, timeoutMs: number | undefined, body: (signal: AbortSignal | undefined) => Promise<T>): Promise<T> { this.#ensureOpen(operation); const releaseClient = this.#owner.begin(operation); const previous = this.#tail; let unlock!: () => void; this.#tail = new Promise<void>((resolve) => { unlock = resolve; }); let acquired = false; let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, timeoutMs); await awaitWithSignal(() => previous, scoped.signal); acquired = true; this.#ensureOpen(operation); return await body(scoped.signal); } catch (error) { throw withContext(error, this.providerName, operation); } finally { scoped?.dispose(); if (acquired) unlock(); else void previous.then(unlock, unlock); releaseClient(); } }
  #ensureOpen(operation: string): void { if (this.#closed || this.#owner.isClosed()) throw sandboxError({ kind: ErrorKind.FailedPrecondition, provider: this.providerName, operation, field: "sandbox", message: "sandbox-kit: sandbox handle is closed" }); if (operation !== "close" && this.#info.terminal) throw sandboxError({ kind: ErrorKind.FailedPrecondition, provider: this.providerName, operation, field: "sandbox", message: "sandbox-kit: sandbox lifecycle is terminal" }); }
}

export class Client implements AsyncDisposable {
  readonly #backend: Backend; readonly #providerName: string; readonly #timeoutMs: number | undefined; readonly #lifecycle: ClientLifecycle; #active = 0; #closing = false; #closed = false; #closePromise: Promise<void> | undefined; #drain: Promise<void> | undefined; #resolveDrain: (() => void) | undefined;
  constructor(config: ClientConfig) { if (!config.provider?.name?.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: "", operation: "initialize", field: "provider", message: "sandbox-kit: provider is required" }); validateConfig(config, config.provider.name, "initialize"); this.#providerName = config.provider.name; this.#timeoutMs = config.timeout === undefined ? undefined : durationMs(config.timeout as Duration); this.#lifecycle = { begin: (operation) => this.#begin(operation), isClosed: () => this.#closing || this.#closed }; try { this.#backend = config.provider.initialize(materializeConfig(config)); } catch (error) { throw withContext(error, config.provider.name, "initialize"); } if (!this.#backend) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: config.provider.name, operation: "initialize", field: "provider", message: "sandbox-kit: provider returned no backend" }); if (this.#backend.name !== this.#providerName) { void this.#backend.close({}).catch(() => {}); throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: this.#providerName, operation: "initialize", field: "provider", message: "sandbox-kit: initialized backend identity differs from selected provider" }); } }
  get providerName(): string { return this.#providerName; }
  async attach(id: string, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("attach"); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (!id.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: this.providerName, operation: "attach", field: "id", message: "sandbox-kit: sandbox ID must not be blank" }); if (!this.#backend.attach) unsupported(this.providerName, "attach"); const info = await awaitWithSignal(() => this.#backend.attach!(id, operationOptions(options, scoped?.signal)), scoped.signal); return new Sandbox(sandboxConstructorToken, this.#backend, this.#lifecycle, { ...info, id, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "attach"); } finally { scoped?.dispose(); release(); } }
  async get(id: string, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("get"); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (!id.trim()) throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: this.providerName, operation: "get", field: "id", message: "sandbox-kit: sandbox ID must not be blank" }); if (!this.#backend.refresh) unsupported(this.providerName, "get"); const info = await awaitWithSignal(() => this.#backend.refresh!(id, operationOptions(options, scoped?.signal)), scoped.signal); return new Sandbox(sandboxConstructorToken, this.#backend, this.#lifecycle, { ...info, id, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "get"); } finally { scoped?.dispose(); release(); } }
  async create(request: CreateOptions = {}, options: OperationOptions = {}): Promise<Sandbox> { const release = this.#begin("create"); const timeoutMs = options.timeoutMs ?? request.provisioning?.timeout ?? this.#timeoutMs; let scoped: ReturnType<typeof operationSignal> | undefined; try { validateCreateOptions(request, this.providerName, "create"); scoped = operationSignal(options, timeoutMs); const create = this.#backend.create; if (!create) unsupported(this.providerName, "create"); const result = await awaitWithSignal(() => create.call(this.#backend, clone(request), operationOptions({ ...(timeoutMs === undefined ? {} : { timeoutMs }) }, scoped?.signal)), scoped.signal); if (!result.sandbox?.id) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: this.providerName, operation: "create", field: "id", message: "sandbox-kit: provider returned no sandbox identity" }); if (result.sandbox.provider && result.sandbox.provider !== this.providerName) throw sandboxError({ kind: ErrorKind.InvalidResponse, provider: this.providerName, operation: "create", field: "provider", message: "sandbox-kit: response provider differs from selected provider" }); return new Sandbox(sandboxConstructorToken, this.#backend, this.#lifecycle, { ...result.sandbox, provider: this.providerName }, this.#timeoutMs); } catch (error) { throw withContext(error, this.providerName, "create"); } finally { scoped?.dispose(); release(); } }
  close(options: OperationOptions = {}): Promise<void> { if (this.#closePromise) return this.#closePromise; this.#closing = true; this.#closePromise = this.#close(options); return this.#closePromise; }
  async #close(options: OperationOptions): Promise<void> { let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); if (this.#active > 0) { this.#drain ??= new Promise<void>((resolve) => { this.#resolveDrain = resolve; }); await awaitWithSignal(() => this.#drain!, scoped.signal); } await awaitWithSignal(() => this.#backend.close(operationOptions(options, scoped?.signal)), scoped.signal); this.#closed = true; } catch (error) { this.#closePromise = undefined; this.#closing = false; throw withContext(error, this.providerName, "close"); } finally { scoped?.dispose(); } }
  async [Symbol.asyncDispose](): Promise<void> { await this.close(); }
  #begin(operation: string): () => void { if (this.#closing || this.#closed) throw sandboxError({ kind: ErrorKind.FailedPrecondition, provider: this.providerName, operation, field: "client", message: "sandbox-kit: client is closed" }); this.#active++; let released = false; return () => { if (released) return; released = true; this.#active--; if (this.#active === 0) { this.#resolveDrain?.(); this.#resolveDrain = undefined; this.#drain = undefined; } }; }
}

/** Materializes a plain client config into the generated protobuf message. */
export function materializeClientConfig(input: ConfigInput): Config {
  return create(ConfigSchema, input);
}

function materializeConfig(config: ClientConfig): Config {
  const { provider: _provider, ...shared } = config;
  return materializeClientConfig(shared);
}
function unsupported(provider: string, operation: string): never { throw sandboxError({ kind: ErrorKind.Unsupported, provider, operation, field: "", message: "sandbox-kit: provider does not support " + operation }); }
function withContext(error: unknown, provider: string, operation: string): SandboxError { if (error instanceof SandboxError) { if (error.provider === provider && error.operation === operation) return error; return sandboxError({ kind: error.kind, provider: error.provider || provider, operation: error.operation || operation, field: error.field, message: error.message, ...(error.rule === undefined ? {} : { rule: error.rule }), ...(error.statusCode === undefined ? {} : { statusCode: error.statusCode }), ...(error.providerCode === undefined ? {} : { providerCode: error.providerCode }), ...(error.providerSource === undefined ? {} : { providerSource: error.providerSource }) }, error.cause); } if (error instanceof DOMException && error.name === "AbortError") return sandboxError({ kind: ErrorKind.Canceled, provider, operation, field: "", message: error.message }, error); if (error instanceof DOMException && error.name === "TimeoutError") return sandboxError({ kind: ErrorKind.Timeout, provider, operation, field: "", message: error.message }, error); return sandboxError({ kind: ErrorKind.Unknown, provider, operation, field: "", message: error instanceof Error ? error.message : "sandbox-kit: provider operation failed" }, error); }

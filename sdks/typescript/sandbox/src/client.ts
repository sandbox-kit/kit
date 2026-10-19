// Client and sandbox runtime behavior is handwritten. Portable client data
// contracts come from standard protobuf generation.

import {
  create,
  type MessageInitShape,
  type MessageShape,
} from "@bufbuild/protobuf";
import { pathToString } from "@bufbuild/protobuf/reflect";
import {
  durationFromMs,
  durationMs,
  type Duration,
} from "@bufbuild/protobuf/wkt";
import { createValidator } from "@bufbuild/protovalidate";
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
import {
  CreateOptionsSchema,
  SandboxState,
  type ArchiveOptions,
  type CreateOptions,
  type CreateOptionsInput,
  type CreateResult,
  type DeleteOptions,
  type ErrorInfo,
  type ExitResult,
  type LifecycleCapability,
  type PauseOptions,
  type SandboxInfo,
  type StartOptions,
  type StopOptions,
  type TerminateOptions,
} from "./sandbox_data_contracts.js";

const createOptionsValidator = createValidator({ failFast: true });

const stableRuleIDs: Record<string, string> = {
  "enum.defined_only": "known_enum",
  "duration.gte": "minimum",
  "oneof.required": "exactly_one",
  cpu_limit: "limit",
  memory_limit: "limit",
  cpu_requires_positive: "requires_positive",
  memory_requires_positive: "requires_positive",
  restricted_requires_allowlist: "requires_any",
  block_all_forbids_network: "forbids",
  allow_all_forbids_allowlists: "forbids",
  after_requires_duration: "requires_any",
  after_forbids_duration: "forbids",
};

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
  create?(
    request: CreateOptions,
    options: OperationOptions,
  ): Promise<CreateResult>;
  attach?(id: string, options: OperationOptions): Promise<SandboxInfo>;
  refresh?(id: string, options: OperationOptions): Promise<SandboxInfo>;
  start?(
    id: string,
    request: StartOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  stop?(
    id: string,
    request: StopOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  pause?(
    id: string,
    request: PauseOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  terminate?(
    id: string,
    request: TerminateOptions,
    options: OperationOptions,
  ): Promise<ExitResult>;
  archive?(
    id: string,
    request: ArchiveOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  delete?(
    id: string,
    request: DeleteOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
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
function validationError(
  provider: string,
  operation: string,
  field: string,
  rule: string,
  message: string,
): never {
  throw sandboxError({
    kind: ErrorKind.InvalidArgument,
    provider,
    operation,
    field,
    rule,
    message: "sandbox-kit: " + message,
  });
}
function joinValidationPath(parent: string, field: string): string {
  return parent ? parent + "." + field : field;
}
function validSecureEndpoint(value: string): boolean {
  try {
    const url = new URL(value);
    if (!url.hostname || url.username || url.password || url.hash) return false;
    if (url.protocol === "https:") return true;
    if (url.protocol !== "http:") return false;
    const host = url.hostname.toLowerCase();
    return (
      host === "localhost" ||
      host === "[::1]" ||
      host === "::1" ||
      /^127(?:\.[0-9]{1,3}){3}$/.test(host)
    );
  } catch {
    return false;
  }
}

export function validateCreateOptions(
  value: CreateOptions,
  provider = "",
  operation = "create",
): void {
  const result = createOptionsValidator.validate(
    CreateOptionsSchema,
    value as MessageShape<typeof CreateOptionsSchema>,
  );
  if (result.kind === "valid") {
    return;
  }
  if (result.kind === "error") {
    throw sandboxError(
      {
        kind: ErrorKind.InvalidArgument,
        provider,
        operation,
        field: "",
        message: "sandbox-kit: invalid configuration",
      },
      result.error,
    );
  }
  const violation = result.violations[0]!;
  const ruleID = violation.ruleId;
  let field = pathToString(violation.field);
  let rule = stableRuleIDs[ruleID] ?? ruleID;
  if (!field) {
    field = messageFieldHint(ruleID);
  }
  validationError(
    provider,
    operation,
    field,
    rule,
    violation.message || "invalid configuration",
  );
}

function messageFieldHint(ruleID: string): string {
  switch (ruleID) {
    case "cpu_limit":
    case "cpu_requires_positive":
      return "resources.cpu_limit_cores";
    case "memory_limit":
    case "memory_requires_positive":
      return "resources.memory_limit_mib";
    case "restricted_requires_allowlist":
      return "network.outbound_cidrs,outbound_domains";
    case "block_all_forbids_network":
      return "network.outbound_cidrs,outbound_domains,private_network,linked_sandbox,inbound_cidrs,ports";
    case "allow_all_forbids_allowlists":
      return "network.outbound_cidrs,outbound_domains";
    case "after_requires_duration":
    case "after_forbids_duration":
      return "after";
    case "oneof.required":
      return "tcp_port,command";
    default:
      return "";
  }
}

function validateAPIKeyCredentialsValue(
  value: APIKeyCredentials | undefined,
  provider: string,
  operation: string,
  path: string,
): void {
  if (value === undefined) return;
  if (value.key !== undefined && !value.key.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "key"),
      "nonblank",
      "must not be blank",
    );
}

function validateAuthConfigValue(
  value: AuthConfig | undefined,
  provider: string,
  operation: string,
  path: string,
): void {
  if (value === undefined) return;
  const credentials = value.credentials;
  if (credentials === undefined || credentials.case === undefined) {
    validationError(
      provider,
      operation,
      joinValidationPath(path, "apiKey,tokenPair,bearerToken,oauthRefresh"),
      "exactly_one",
      "requires exactly one value",
    );
  }
  if (credentials.case === "apiKey")
    validateAPIKeyCredentialsValue(
      credentials.value,
      provider,
      operation,
      joinValidationPath(path, "apiKey"),
    );
  if (credentials.case === "tokenPair")
    validateTokenPairCredentialsValue(
      credentials.value,
      provider,
      operation,
      joinValidationPath(path, "tokenPair"),
    );
  if (credentials.case === "bearerToken")
    validateBearerTokenCredentialsValue(
      credentials.value,
      provider,
      operation,
      joinValidationPath(path, "bearerToken"),
    );
  if (credentials.case === "oauthRefresh")
    validateOAuthCredentialsValue(
      credentials.value,
      provider,
      operation,
      joinValidationPath(path, "oauthRefresh"),
    );
}

function validateBearerTokenCredentialsValue(
  value: BearerTokenCredentials | undefined,
  provider: string,
  operation: string,
  path: string,
): void {
  if (value === undefined) return;
  if (value.token !== undefined && !value.token.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "token"),
      "nonblank",
      "must not be blank",
    );
}

function validateConfigValue(
  value: Config | undefined,
  provider: string,
  operation: string,
  path: string,
): void {
  if (value === undefined) return;
  if (value.endpoint !== undefined && !value.endpoint.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "endpoint"),
      "nonblank",
      "must not be blank",
    );
  if (value.endpoint !== undefined && !validSecureEndpoint(value.endpoint))
    validationError(
      provider,
      operation,
      joinValidationPath(path, "endpoint"),
      "secure_endpoint",
      "must use HTTPS except for loopback HTTP development endpoints",
    );
  if (value.region !== undefined && !value.region.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "region"),
      "nonblank",
      "must not be blank",
    );
  if (value.timeout !== undefined) {
    let timeoutMs: number;
    try {
      timeoutMs = durationMs(value.timeout);
    } catch {
      validationError(
        provider,
        operation,
        joinValidationPath(path, "timeout"),
        "finite",
        "must be a non-negative protobuf duration",
      );
    }
    if (!Number.isFinite(timeoutMs) || timeoutMs < 0)
      validationError(
        provider,
        operation,
        joinValidationPath(path, "timeout"),
        "minimum",
        "must be at least 0",
      );
  }
  validateAuthConfigValue(
    value.auth,
    provider,
    operation,
    joinValidationPath(path, "auth"),
  );
  validateScopeValue(
    value.scope,
    provider,
    operation,
    joinValidationPath(path, "scope"),
  );
}

function validateOAuthCredentialsValue(
  value: OAuthCredentials | undefined,
  provider: string,
  operation: string,
  path: string,
): void {
  if (value === undefined) return;
  if (value.clientId !== undefined && !value.clientId.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "clientId"),
      "nonblank",
      "must not be blank",
    );
  if (value.clientSecret !== undefined && !value.clientSecret.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "clientSecret"),
      "nonblank",
      "must not be blank",
    );
  if (value.jwtKey !== undefined && !value.jwtKey.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "jwtKey"),
      "nonblank",
      "must not be blank",
    );
  if (value.refreshToken !== undefined && !value.refreshToken.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "refreshToken"),
      "nonblank",
      "must not be blank",
    );
  if (true) {
    const count0 = [
      value.clientSecret !== undefined,
      value.jwtKey !== undefined,
    ].filter(Boolean).length;
    if (count0 !== 1)
      validationError(
        provider,
        operation,
        joinValidationPath(path, "clientSecret,jwtKey"),
        "exactly_one",
        "requires exactly one value",
      );
  }
}

function validateScopeValue(
  value: Scope | undefined,
  provider: string,
  operation: string,
  path: string,
): void {
  if (value === undefined) return;
  if (value.appName !== undefined && !value.appName.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "appName"),
      "nonblank",
      "must not be blank",
    );
  if (value.environment !== undefined && !value.environment.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "environment"),
      "nonblank",
      "must not be blank",
    );
  if (value.organizationId !== undefined && !value.organizationId.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "organizationId"),
      "nonblank",
      "must not be blank",
    );
  if (value.projectId !== undefined && !value.projectId.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "projectId"),
      "nonblank",
      "must not be blank",
    );
}

function validateTokenPairCredentialsValue(
  value: TokenPairCredentials | undefined,
  provider: string,
  operation: string,
  path: string,
): void {
  if (value === undefined) return;
  if (value.id !== undefined && !value.id.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "id"),
      "nonblank",
      "must not be blank",
    );
  if (value.secret !== undefined && !value.secret.trim())
    validationError(
      provider,
      operation,
      joinValidationPath(path, "secret"),
      "nonblank",
      "must not be blank",
    );
}

export function validateConfig(
  value: Config | ConfigInput,
  provider = "",
  operation = "create",
): void {
  validateConfigValue(value as Config, provider, operation, "");
}

function clone<T>(value: T): T {
  return structuredClone(value);
}

function operationSignal(
  options: OperationOptions,
  inheritedTimeoutMs?: number,
): { signal?: AbortSignal; dispose(): void } {
  const timeoutMs = options.timeoutMs ?? inheritedTimeoutMs;
  if (
    timeoutMs !== undefined &&
    (!Number.isFinite(timeoutMs) || timeoutMs < 0)
  ) {
    throw sandboxError({
      kind: ErrorKind.InvalidArgument,
      provider: "",
      operation: "",
      field: "timeoutMs",
      message: "sandbox-kit: timeoutMs must be a finite non-negative number",
    });
  }
  if (
    (timeoutMs === undefined || timeoutMs === 0) &&
    options.signal === undefined
  )
    return { dispose() {} };
  const controller = new AbortController();
  const abort = () => controller.abort(options.signal?.reason);
  if (options.signal?.aborted) abort();
  else options.signal?.addEventListener("abort", abort, { once: true });
  const timer =
    timeoutMs !== undefined && timeoutMs > 0
      ? setTimeout(
          () =>
            controller.abort(
              new DOMException("The operation timed out", "TimeoutError"),
            ),
          timeoutMs,
        )
      : undefined;
  return {
    signal: controller.signal,
    dispose() {
      if (timer !== undefined) clearTimeout(timer);
      options.signal?.removeEventListener("abort", abort);
    },
  };
}

function operationOptions(
  options: OperationOptions,
  signal: AbortSignal | undefined,
): OperationOptions {
  return {
    ...(options.timeoutMs === undefined
      ? {}
      : { timeoutMs: options.timeoutMs }),
    ...(signal === undefined ? {} : { signal }),
  };
}

function abortReason(signal: AbortSignal): DOMException {
  const reason = signal.reason;
  return reason instanceof DOMException &&
    (reason.name === "AbortError" || reason.name === "TimeoutError")
    ? reason
    : new DOMException("The operation was aborted", "AbortError");
}

async function awaitWithSignal<T>(
  operation: () => Promise<T>,
  signal?: AbortSignal,
): Promise<T> {
  if (signal?.aborted) throw abortReason(signal);
  const pending = operation();
  if (signal === undefined) return pending;
  return new Promise<T>((resolve, reject) => {
    const aborted = () => reject(abortReason(signal));
    const cleanup = () => signal.removeEventListener("abort", aborted);
    signal.addEventListener("abort", aborted, { once: true });
    pending.then(
      (value) => {
        cleanup();
        resolve(value);
      },
      (error: unknown) => {
        cleanup();
        reject(error);
      },
    );
  });
}
interface ClientLifecycle {
  begin(operation: string): () => void;
  isClosed(): boolean;
}
const sandboxConstructorToken: unique symbol = Symbol("sandbox-kit.Sandbox");
export class Sandbox implements AsyncDisposable {
  readonly #backend: Backend;
  readonly #owner: ClientLifecycle;
  readonly #timeoutMs: number | undefined;
  #info: SandboxInfo;
  #closed = false;
  #tail: Promise<void> = Promise.resolve();
  constructor(
    token: typeof sandboxConstructorToken,
    backend: Backend,
    owner: ClientLifecycle,
    info: SandboxInfo,
    timeoutMs?: number,
  ) {
    if (token !== sandboxConstructorToken)
      throw new TypeError("Sandbox handles are created by Client");
    this.#backend = backend;
    this.#owner = owner;
    this.#info = clone(info);
    this.#timeoutMs = timeoutMs;
  }
  get id(): string {
    return this.#info.id;
  }
  get providerName(): string {
    return this.#info.provider;
  }
  info(): SandboxInfo {
    return clone(this.#info);
  }
  capabilities(): readonly LifecycleCapability[] {
    if (this.#closed || this.#owner.isClosed() || this.#info.terminal)
      return [];
    return this.#backend.capabilities?.(clone(this.#info)) ?? [];
  }
  supports(capability: LifecycleCapability): boolean {
    return this.capabilities().includes(capability);
  }
  async refresh(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#serialized(
      "refresh",
      options,
      this.#timeoutMs,
      async (signal) => {
        const method = this.#backend.refresh;
        if (!method) unsupported(this.providerName, "refresh");
        const info = await awaitWithSignal(
          () =>
            method.call(
              this.#backend,
              this.id,
              operationOptions(options, signal),
            ),
          signal,
        );
        this.#info = clone({ ...info, provider: this.providerName });
        return this.info();
      },
    );
  }
  async start(
    request: StartOptions = {},
    options: OperationOptions = {},
  ): Promise<SandboxInfo> {
    return this.#update("start", this.#backend.start, request, options);
  }
  async stop(
    request: StopOptions = {},
    options: OperationOptions = {},
  ): Promise<SandboxInfo> {
    return this.#update("stop", this.#backend.stop, request, options);
  }
  async pause(
    request: PauseOptions = {},
    options: OperationOptions = {},
  ): Promise<SandboxInfo> {
    return this.#update("pause", this.#backend.pause, request, options);
  }
  async terminate(
    request: TerminateOptions = {},
    options: OperationOptions = {},
  ): Promise<ExitResult> {
    return this.#serialized(
      "terminate",
      options,
      optionalTimeoutMs(request.timeout) ?? this.#timeoutMs,
      async (signal) => {
        const method = this.#backend.terminate;
        if (!method) unsupported(this.providerName, "terminate");
        const result = await awaitWithSignal(
          () =>
            method.call(
              this.#backend,
              this.id,
              clone(request),
              operationOptions(options, signal),
            ),
          signal,
        );
        this.#info = {
          ...this.#info,
          state: SandboxState.Terminated,
          terminal: true,
        };
        return result;
      },
    );
  }
  async archive(
    request: ArchiveOptions = {},
    options: OperationOptions = {},
  ): Promise<SandboxInfo> {
    return this.#update("archive", this.#backend.archive, request, options);
  }
  async delete(
    request: DeleteOptions = {},
    options: OperationOptions = {},
  ): Promise<SandboxInfo> {
    return this.#update("delete", this.#backend.delete, request, options);
  }
  async close(options: OperationOptions = {}): Promise<void> {
    if (this.#closed) return;
    await this.#serialized(
      "close",
      options,
      this.#timeoutMs,
      async (signal) => {
        const close = this.#backend.closeSandbox;
        if (close)
          await awaitWithSignal(
            () =>
              close.call(
                this.#backend,
                this.id,
                operationOptions(options, signal),
              ),
            signal,
          );
        this.#closed = true;
      },
    );
  }
  async [Symbol.asyncDispose](): Promise<void> {
    await this.close();
  }
  async #update<
    T extends
      | StartOptions
      | StopOptions
      | PauseOptions
      | ArchiveOptions
      | DeleteOptions,
  >(
    operation: string,
    method:
      | ((
          id: string,
          request: T,
          options: OperationOptions,
        ) => Promise<SandboxInfo>)
      | undefined,
    request: T,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    return this.#serialized(
      operation,
      options,
      optionalTimeoutMs(request.timeout) ?? this.#timeoutMs,
      async (signal) => {
        if (!method) unsupported(this.providerName, operation);
        const info = await awaitWithSignal(
          () =>
            method.call(
              this.#backend,
              this.id,
              clone(request),
              operationOptions(options, signal),
            ),
          signal,
        );
        this.#info = clone({ ...info, provider: this.providerName });
        return this.info();
      },
    );
  }
  async #serialized<T>(
    operation: string,
    options: OperationOptions,
    timeoutMs: number | undefined,
    body: (signal: AbortSignal | undefined) => Promise<T>,
  ): Promise<T> {
    this.#ensureOpen(operation);
    const releaseClient = this.#owner.begin(operation);
    const previous = this.#tail;
    let unlock!: () => void;
    this.#tail = new Promise<void>((resolve) => {
      unlock = resolve;
    });
    let acquired = false;
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      scoped = operationSignal(options, timeoutMs);
      await awaitWithSignal(() => previous, scoped.signal);
      acquired = true;
      this.#ensureOpen(operation);
      return await body(scoped.signal);
    } catch (error) {
      throw withContext(error, this.providerName, operation);
    } finally {
      scoped?.dispose();
      if (acquired) unlock();
      else void previous.then(unlock, unlock);
      releaseClient();
    }
  }
  #ensureOpen(operation: string): void {
    if (this.#closed || this.#owner.isClosed())
      throw sandboxError({
        kind: ErrorKind.FailedPrecondition,
        provider: this.providerName,
        operation,
        field: "sandbox",
        message: "sandbox-kit: sandbox handle is closed",
      });
    if (operation !== "close" && this.#info.terminal)
      throw sandboxError({
        kind: ErrorKind.FailedPrecondition,
        provider: this.providerName,
        operation,
        field: "sandbox",
        message: "sandbox-kit: sandbox lifecycle is terminal",
      });
  }
}

export class Client implements AsyncDisposable {
  readonly #backend: Backend;
  readonly #providerName: string;
  readonly #timeoutMs: number | undefined;
  readonly #lifecycle: ClientLifecycle;
  #active = 0;
  #closing = false;
  #closed = false;
  #closePromise: Promise<void> | undefined;
  #drain: Promise<void> | undefined;
  #resolveDrain: (() => void) | undefined;
  constructor(config: ClientConfig) {
    if (!config.provider?.name?.trim())
      throw sandboxError({
        kind: ErrorKind.InvalidArgument,
        provider: "",
        operation: "initialize",
        field: "provider",
        message: "sandbox-kit: provider is required",
      });
    validateConfig(config, config.provider.name, "initialize");
    this.#providerName = config.provider.name;
    this.#timeoutMs =
      config.timeout === undefined
        ? undefined
        : durationMs(config.timeout as Duration);
    this.#lifecycle = {
      begin: (operation) => this.#begin(operation),
      isClosed: () => this.#closing || this.#closed,
    };
    try {
      this.#backend = config.provider.initialize(materializeConfig(config));
    } catch (error) {
      throw withContext(error, config.provider.name, "initialize");
    }
    if (!this.#backend)
      throw sandboxError({
        kind: ErrorKind.InvalidResponse,
        provider: config.provider.name,
        operation: "initialize",
        field: "provider",
        message: "sandbox-kit: provider returned no backend",
      });
    if (this.#backend.name !== this.#providerName) {
      void this.#backend.close({}).catch(() => {});
      throw sandboxError({
        kind: ErrorKind.InvalidResponse,
        provider: this.#providerName,
        operation: "initialize",
        field: "provider",
        message:
          "sandbox-kit: initialized backend identity differs from selected provider",
      });
    }
  }
  get providerName(): string {
    return this.#providerName;
  }
  async attach(id: string, options: OperationOptions = {}): Promise<Sandbox> {
    const release = this.#begin("attach");
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      scoped = operationSignal(options, this.#timeoutMs);
      if (!id.trim())
        throw sandboxError({
          kind: ErrorKind.InvalidArgument,
          provider: this.providerName,
          operation: "attach",
          field: "id",
          message: "sandbox-kit: sandbox ID must not be blank",
        });
      if (!this.#backend.attach) unsupported(this.providerName, "attach");
      const info = await awaitWithSignal(
        () =>
          this.#backend.attach!(id, operationOptions(options, scoped?.signal)),
        scoped.signal,
      );
      return new Sandbox(
        sandboxConstructorToken,
        this.#backend,
        this.#lifecycle,
        { ...info, id, provider: this.providerName },
        this.#timeoutMs,
      );
    } catch (error) {
      throw withContext(error, this.providerName, "attach");
    } finally {
      scoped?.dispose();
      release();
    }
  }
  async get(id: string, options: OperationOptions = {}): Promise<Sandbox> {
    const release = this.#begin("get");
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      scoped = operationSignal(options, this.#timeoutMs);
      if (!id.trim())
        throw sandboxError({
          kind: ErrorKind.InvalidArgument,
          provider: this.providerName,
          operation: "get",
          field: "id",
          message: "sandbox-kit: sandbox ID must not be blank",
        });
      if (!this.#backend.refresh) unsupported(this.providerName, "get");
      const info = await awaitWithSignal(
        () =>
          this.#backend.refresh!(id, operationOptions(options, scoped?.signal)),
        scoped.signal,
      );
      return new Sandbox(
        sandboxConstructorToken,
        this.#backend,
        this.#lifecycle,
        { ...info, id, provider: this.providerName },
        this.#timeoutMs,
      );
    } catch (error) {
      throw withContext(error, this.providerName, "get");
    } finally {
      scoped?.dispose();
      release();
    }
  }
  async create(
    request: CreateOptionsInput = {},
    options: OperationOptions = {},
  ): Promise<Sandbox> {
    const release = this.#begin("create");
    const materialized = materializeCreateOptions(request);
    const timeoutMs =
      options.timeoutMs ??
      optionalTimeoutMs(materialized.provisioning?.timeout) ??
      this.#timeoutMs;
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      validateCreateOptions(materialized, this.providerName, "create");
      scoped = operationSignal(options, timeoutMs);
      const create = this.#backend.create;
      if (!create) unsupported(this.providerName, "create");
      const result = await awaitWithSignal(
        () =>
          create.call(
            this.#backend,
            clone(materialized),
            operationOptions(
              { ...(timeoutMs === undefined ? {} : { timeoutMs }) },
              scoped?.signal,
            ),
          ),
        scoped.signal,
      );
      if (!result.sandbox?.id)
        throw sandboxError({
          kind: ErrorKind.InvalidResponse,
          provider: this.providerName,
          operation: "create",
          field: "id",
          message: "sandbox-kit: provider returned no sandbox identity",
        });
      if (
        result.sandbox.provider &&
        result.sandbox.provider !== this.providerName
      )
        throw sandboxError({
          kind: ErrorKind.InvalidResponse,
          provider: this.providerName,
          operation: "create",
          field: "provider",
          message:
            "sandbox-kit: response provider differs from selected provider",
        });
      return new Sandbox(
        sandboxConstructorToken,
        this.#backend,
        this.#lifecycle,
        { ...result.sandbox, provider: this.providerName },
        this.#timeoutMs,
      );
    } catch (error) {
      throw withContext(error, this.providerName, "create");
    } finally {
      scoped?.dispose();
      release();
    }
  }
  close(options: OperationOptions = {}): Promise<void> {
    if (this.#closePromise) return this.#closePromise;
    this.#closing = true;
    this.#closePromise = this.#close(options);
    return this.#closePromise;
  }
  async #close(options: OperationOptions): Promise<void> {
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      scoped = operationSignal(options, this.#timeoutMs);
      if (this.#active > 0) {
        this.#drain ??= new Promise<void>((resolve) => {
          this.#resolveDrain = resolve;
        });
        await awaitWithSignal(() => this.#drain!, scoped.signal);
      }
      await awaitWithSignal(
        () => this.#backend.close(operationOptions(options, scoped?.signal)),
        scoped.signal,
      );
      this.#closed = true;
    } catch (error) {
      this.#closePromise = undefined;
      this.#closing = false;
      throw withContext(error, this.providerName, "close");
    } finally {
      scoped?.dispose();
    }
  }
  async [Symbol.asyncDispose](): Promise<void> {
    await this.close();
  }
  #begin(operation: string): () => void {
    if (this.#closing || this.#closed)
      throw sandboxError({
        kind: ErrorKind.FailedPrecondition,
        provider: this.providerName,
        operation,
        field: "client",
        message: "sandbox-kit: client is closed",
      });
    this.#active++;
    let released = false;
    return () => {
      if (released) return;
      released = true;
      this.#active--;
      if (this.#active === 0) {
        this.#resolveDrain?.();
        this.#resolveDrain = undefined;
        this.#drain = undefined;
      }
    };
  }
}

/** Materializes a plain client config into the generated protobuf message. */
export function materializeClientConfig(input: ConfigInput): Config {
  return create(ConfigSchema, input);
}

/** Materializes plain create options into the generated protobuf message. */
export function materializeCreateOptions(
  input: CreateOptionsInput,
): CreateOptions {
  return create(CreateOptionsSchema, input) as CreateOptions;
}

function optionalTimeoutMs(value: Duration | undefined): number | undefined {
  if (value === undefined) return undefined;
  return durationMs(value);
}

function materializeConfig(config: ClientConfig): Config {
  const { provider: _provider, ...shared } = config;
  return materializeClientConfig(shared);
}
function unsupported(provider: string, operation: string): never {
  throw sandboxError({
    kind: ErrorKind.Unsupported,
    provider,
    operation,
    field: "",
    message: "sandbox-kit: provider does not support " + operation,
  });
}
function withContext(
  error: unknown,
  provider: string,
  operation: string,
): SandboxError {
  if (error instanceof SandboxError) {
    if (error.provider === provider && error.operation === operation)
      return error;
    return sandboxError(
      {
        kind: error.kind,
        provider: error.provider || provider,
        operation: error.operation || operation,
        field: error.field,
        message: error.message,
        ...(error.rule === undefined ? {} : { rule: error.rule }),
        ...(error.statusCode === undefined
          ? {}
          : { statusCode: error.statusCode }),
        ...(error.providerCode === undefined
          ? {}
          : { providerCode: error.providerCode }),
        ...(error.providerSource === undefined
          ? {}
          : { providerSource: error.providerSource }),
      },
      error.cause,
    );
  }
  if (error instanceof DOMException && error.name === "AbortError")
    return sandboxError(
      {
        kind: ErrorKind.Canceled,
        provider,
        operation,
        field: "",
        message: error.message,
      },
      error,
    );
  if (error instanceof DOMException && error.name === "TimeoutError")
    return sandboxError(
      {
        kind: ErrorKind.Timeout,
        provider,
        operation,
        field: "",
        message: error.message,
      },
      error,
    );
  return sandboxError(
    {
      kind: ErrorKind.Unknown,
      provider,
      operation,
      field: "",
      message:
        error instanceof Error
          ? error.message
          : "sandbox-kit: provider operation failed",
    },
    error,
  );
}

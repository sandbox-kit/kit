// Client and sandbox runtime behavior is handwritten. Portable client data
// contracts come from standard protobuf generation.

import {
  create,
  type MessageInitShape,
  type MessageShape,
} from "@bufbuild/protobuf";
import { pathToString } from "@bufbuild/protobuf/reflect";
import { durationFromMs, type Duration } from "@bufbuild/protobuf/wkt";
import { createValidator } from "@bufbuild/protovalidate";
import {
  ConfigSchema,
  type APIKeyCredentials as WireAPIKeyCredentials,
  type AuthConfig,
  type BearerTokenCredentials as WireBearerTokenCredentials,
  type Config as WireConfig,
  type OAuthCredentials as WireOAuthCredentials,
  type TokenPairCredentials as WireTokenPairCredentials,
} from "@sandbox-kit/contracts/kit/sandbox/v1/client_pb.js";
import {
  CreateOptionsSchema,
  ErrorKind as ProtoErrorKind,
  type CommandProbe,
  type CreateOptions as ProtoCreateOptions,
  type ImageSource,
  type ProviderDefaultSource,
  type ReadinessProbe,
  type SandboxSource,
  type SnapshotSource,
  type WarmPoolSource,
} from "@sandbox-kit/contracts/kit/sandbox/v1/sandbox_pb.js";
import {
  SandboxState,
  type ArchiveOptions,
  type CreateOptions,
  type CreateResult,
  type DeleteOptions,
  type ErrorInfo,
  type ExitResult,
  type LifecycleCapability,
  type NetworkConfigInput,
  type PauseOptions,
  type ReadinessProbeInput,
  type SandboxInfo,
  type SandboxSourceInput,
  type StartOptions,
  type StopOptions,
  type TerminateOptions,
} from "./sandbox_data_contracts.js";
import {
  materializeDurationFields,
  type DurationFieldsToMs,
} from "./duration_ms.js";
import type { RemapFields } from "./input_remap.js";
import { metadataObjectFromMap } from "./metadata_object.js";
import type { FlattenOneof } from "./oneof_input.js";
import { protoEnum } from "./proto_enum.js";
import type { StripProto } from "./strip_proto.js";

const contractValidator = createValidator({ failFast: true });

const stableRuleIDs: Record<string, string> = {
  "enum.defined_only": "known_enum",
  "duration.gte": "minimum",
  "oneof.required": "exactly_one",
  auth_credentials_required: "exactly_one",
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

/** Public credential shapes (no protobuf `$typeName` / getters). */
export type APIKeyCredentials = StripProto<WireAPIKeyCredentials>;
export type TokenPairCredentials = StripProto<WireTokenPairCredentials>;
export type BearerTokenCredentials = StripProto<WireBearerTokenCredentials>;
export type OAuthCredentials = StripProto<WireOAuthCredentials>;
export type Scope = StripProto<
  import("@sandbox-kit/contracts/kit/sandbox/v1/client_pb.js").Scope
>;

/** Public auth: set exactly one credential variant (no protobuf `{ case, value }`). */
export type AuthInput = FlattenOneof<NonNullable<AuthConfig["credentials"]>>;

/**
 * Public client settings (ms timeout, flat auth).
 * Prefer {@link ClientConfig} at construction (`Config` plus `provider`).
 */
export type Config = DurationFieldsToMs<
  RemapFields<
    StripProto<MessageInitShape<typeof ConfigSchema>>,
    { auth: AuthInput }
  >
>;

/** Wire config after materialize; package-private for Protovalidate. */
type ConfigWire = StripProto<WireConfig>;
/** Wire create options after materialize; package-private for Protovalidate. */
type CreateOptionsWire = StripProto<ProtoCreateOptions>;

// Ergonomic PascalCase names over protobuf enum values; not a second contract.
export const ErrorKind = protoEnum(ProtoErrorKind);
export type ErrorKind = (typeof ErrorKind)[keyof typeof ErrorKind];

export interface OperationOptions {
  readonly signal?: AbortSignal;
  /** Operation deadline in milliseconds. */
  readonly timeout?: number;
}

/** Public stop options: Kit call options plus optional force. */
export type StopCallOptions = OperationOptions & { readonly force?: boolean };
/** Public terminate options: Kit call options plus optional wait. */
export type TerminateCallOptions = OperationOptions & {
  readonly wait?: boolean;
};
/** Public delete options: Kit call options plus optional wait. */
export type DeleteCallOptions = OperationOptions & { readonly wait?: boolean };

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
  /** Receives public {@link Config}. */
  initialize(config: Config): Backend;
}

/** Provider selection stays handwritten and is never part of the portable Config message. */
export type ClientConfig = Config & {
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

function validateCreateOptions(
  value: CreateOptionsWire,
  provider = "",
  operation = "create",
): void {
  const result = contractValidator.validate(
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
    case "auth_credentials_required":
      return "api_key,token_pair,bearer_token,oauth_refresh";
    case "exactly_one":
      return "client_secret,jwt_key";
    default:
      return "";
  }
}

export function validateConfig(
  value: Config,
  provider = "",
  operation = "create",
): void {
  const message = materializeClientConfig(value);
  const result = contractValidator.validate(
    ConfigSchema,
    message as MessageShape<typeof ConfigSchema>,
  );
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
  if (result.kind === "invalid") {
    const violation = result.violations[0]!;
    const ruleID = violation.ruleId;
    let field = pathToString(violation.field);
    const rule = stableRuleIDs[ruleID] ?? ruleID;
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
  if (
    message.endpoint !== undefined &&
    !validSecureEndpoint(message.endpoint)
  ) {
    validationError(
      provider,
      operation,
      "endpoint",
      "secure_endpoint",
      "must use HTTPS except for loopback HTTP development endpoints",
    );
  }
}

function clone<T>(value: T): T {
  return structuredClone(value);
}

function operationSignal(
  options: OperationOptions,
  inheritedTimeoutMs?: number,
): { signal?: AbortSignal; dispose(): void } {
  const timeout = options.timeout ?? inheritedTimeoutMs;
  if (timeout !== undefined && (!Number.isFinite(timeout) || timeout < 0)) {
    throw sandboxError({
      kind: ErrorKind.InvalidArgument,
      provider: "",
      operation: "",
      field: "timeout",
      message: "sandbox-kit: timeout must be a finite non-negative number",
    });
  }
  if ((timeout === undefined || timeout === 0) && options.signal === undefined)
    return { dispose() {} };
  const controller = new AbortController();
  const abort = () => controller.abort(options.signal?.reason);
  if (options.signal?.aborted) abort();
  else options.signal?.addEventListener("abort", abort, { once: true });
  const timer =
    timeout !== undefined && timeout > 0
      ? setTimeout(
          () =>
            controller.abort(
              new DOMException("The operation timed out", "TimeoutError"),
            ),
          timeout,
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
    ...(options.timeout === undefined ? {} : { timeout: options.timeout }),
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
  readonly #timeout: number | undefined;
  #info: SandboxInfo;
  #closed = false;
  #tail: Promise<void> = Promise.resolve();
  constructor(
    token: typeof sandboxConstructorToken,
    backend: Backend,
    owner: ClientLifecycle,
    info: SandboxInfo,
    timeout?: number,
  ) {
    if (token !== sandboxConstructorToken)
      throw new TypeError("Sandbox handles are created by Client");
    this.#backend = backend;
    this.#owner = owner;
    this.#info = clone(info);
    this.#timeout = timeout;
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
      this.#timeout,
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
  async start(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "start",
      this.#backend.start,
      wireTimeoutOptions(options),
      options,
    );
  }
  async stop(options: StopCallOptions = {}): Promise<SandboxInfo> {
    const { force, ...call } = options;
    return this.#update(
      "stop",
      this.#backend.stop,
      {
        ...wireTimeoutOptions(call),
        ...(force === undefined ? {} : { force }),
      },
      call,
    );
  }
  async pause(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "pause",
      this.#backend.pause,
      wireTimeoutOptions(options),
      options,
    );
  }
  async terminate(options: TerminateCallOptions = {}): Promise<ExitResult> {
    const { wait, ...call } = options;
    const wire: TerminateOptions = {
      ...wireTimeoutOptions(call),
      ...(wait === undefined ? {} : { wait }),
    };
    return this.#serialized(
      "terminate",
      call,
      call.timeout ?? this.#timeout,
      async (signal) => {
        const method = this.#backend.terminate;
        if (!method) unsupported(this.providerName, "terminate");
        const result = await awaitWithSignal(
          () =>
            method.call(
              this.#backend,
              this.id,
              clone(wire),
              operationOptions(call, signal),
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
  async archive(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "archive",
      this.#backend.archive,
      wireTimeoutOptions(options),
      options,
    );
  }
  async delete(options: DeleteCallOptions = {}): Promise<SandboxInfo> {
    const { wait, ...call } = options;
    return this.#update(
      "delete",
      this.#backend.delete,
      { ...wireTimeoutOptions(call), ...(wait === undefined ? {} : { wait }) },
      call,
    );
  }
  async close(options: OperationOptions = {}): Promise<void> {
    if (this.#closed) return;
    await this.#serialized("close", options, this.#timeout, async (signal) => {
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
    });
  }
  async [Symbol.asyncDispose](): Promise<void> {
    await this.close();
  }
  async #update(
    operation: string,
    method:
      | ((
          id: string,
          request:
            | StartOptions
            | StopOptions
            | PauseOptions
            | ArchiveOptions
            | DeleteOptions,
          options: OperationOptions,
        ) => Promise<SandboxInfo>)
      | undefined,
    wire:
      | StartOptions
      | StopOptions
      | PauseOptions
      | ArchiveOptions
      | DeleteOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    return this.#serialized(
      operation,
      options,
      options.timeout ?? this.#timeout,
      async (signal) => {
        if (!method) unsupported(this.providerName, operation);
        const info = await awaitWithSignal(
          () =>
            method.call(
              this.#backend,
              this.id,
              clone(wire),
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
    timeout: number | undefined,
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
      scoped = operationSignal(options, timeout);
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
  readonly #timeout: number | undefined;
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
    this.#timeout = config.timeout;
    this.#lifecycle = {
      begin: (operation) => this.#begin(operation),
      isClosed: () => this.#closing || this.#closed,
    };
    try {
      const { provider: _provider, ...shared } = config;
      this.#backend = config.provider.initialize(shared);
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
      scoped = operationSignal(options, this.#timeout);
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
        this.#timeout,
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
      scoped = operationSignal(options, this.#timeout);
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
        this.#timeout,
      );
    } catch (error) {
      throw withContext(error, this.providerName, "get");
    } finally {
      scoped?.dispose();
      release();
    }
  }
  async create(
    request: CreateOptions = {},
    options: OperationOptions = {},
  ): Promise<Sandbox> {
    const release = this.#begin("create");
    const materialized = materializeCreateOptions(request);
    const timeout =
      options.timeout ?? request.provisioning?.timeout ?? this.#timeout;
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      // Validate the wire projection first; providers receive the public request.
      validateCreateOptions(materialized, this.providerName, "create");
      scoped = operationSignal(options, timeout);
      const create = this.#backend.create;
      if (!create) unsupported(this.providerName, "create");
      const providerRequest: CreateOptions = {
        ...request,
        ...(timeout === undefined
          ? {}
          : {
              provisioning: {
                ...request.provisioning,
                timeout,
              },
            }),
      };
      const result = await awaitWithSignal(
        () =>
          create.call(
            this.#backend,
            providerRequest,
            operationOptions(
              { ...(timeout === undefined ? {} : { timeout: timeout }) },
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
        this.#timeout,
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
      scoped = operationSignal(options, this.#timeout);
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
function materializeClientConfig(input: Config): ConfigWire {
  const { auth, timeout, ...rest } = input;
  return create(ConfigSchema, {
    ...rest,
    ...(timeout === undefined ? {} : { timeout: durationFromMs(timeout) }),
    ...(auth === undefined
      ? {}
      : { auth: { credentials: authCredentials(auth) } }),
  } as MessageInitShape<typeof ConfigSchema>) as ConfigWire;
}

function wireTimeoutOptions(options: OperationOptions): { timeout?: Duration } {
  return options.timeout === undefined
    ? {}
    : { timeout: durationFromMs(options.timeout) };
}

/** Materializes plain create options into the generated protobuf message. */
function materializeCreateOptions(input: CreateOptions): CreateOptionsWire {
  const { source, readiness, providerOptions, network, ...rest } = input;
  const withDurations = materializeDurationFields(rest);
  return create(CreateOptionsSchema, {
    ...withDurations,
    ...(network === undefined
      ? {}
      : { network: materializeNetworkConfig(network) }),
    ...(source === undefined ? {} : { source: { kind: sourceKind(source) } }),
    ...(readiness === undefined
      ? {}
      : { readiness: { check: readinessCheck(readiness) } }),
    ...(providerOptions === undefined
      ? {}
      : { providerOptions: metadataObjectFromMap(providerOptions) }),
  } as unknown as MessageInitShape<
    typeof CreateOptionsSchema
  >) as CreateOptionsWire;
}

/** Wrap public string-array allowlists into wire StringAllowlist messages. */
function materializeNetworkConfig(network: NetworkConfigInput) {
  const { outboundCidrs, outboundDomains, inboundCidrs, ...rest } = network;
  return {
    ...rest,
    ...(outboundCidrs === undefined
      ? {}
      : { outboundCidrs: { entries: outboundCidrs } }),
    ...(outboundDomains === undefined
      ? {}
      : { outboundDomains: { entries: outboundDomains } }),
    ...(inboundCidrs === undefined
      ? {}
      : { inboundCidrs: { entries: inboundCidrs } }),
  };
}

function authCredentials(
  auth: AuthInput,
): NonNullable<AuthConfig["credentials"]> {
  if (auth.apiKey !== undefined)
    return { case: "apiKey", value: auth.apiKey as WireAPIKeyCredentials };
  if (auth.tokenPair !== undefined)
    return {
      case: "tokenPair",
      value: auth.tokenPair as WireTokenPairCredentials,
    };
  if (auth.bearerToken !== undefined)
    return {
      case: "bearerToken",
      value: auth.bearerToken as WireBearerTokenCredentials,
    };
  if (auth.oauthRefresh !== undefined)
    return {
      case: "oauthRefresh",
      value: auth.oauthRefresh as WireOAuthCredentials,
    };
  return { case: undefined };
}

function sourceKind(
  source: SandboxSourceInput,
): NonNullable<SandboxSource["kind"]> {
  if (source.providerDefault !== undefined)
    return {
      case: "providerDefault",
      value: {} as ProviderDefaultSource,
    };
  if (source.image !== undefined)
    return { case: "image", value: source.image as ImageSource };
  if (source.snapshot !== undefined)
    return { case: "snapshot", value: source.snapshot as SnapshotSource };
  if (source.warmPool !== undefined)
    return { case: "warmPool", value: source.warmPool as WarmPoolSource };
  return { case: undefined } as NonNullable<SandboxSource["kind"]>;
}

function readinessCheck(
  readiness: ReadinessProbeInput,
): NonNullable<ReadinessProbe["check"]> {
  if (readiness.tcpPort !== undefined)
    return { case: "tcpPort", value: readiness.tcpPort };
  if (readiness.command !== undefined)
    return { case: "command", value: readiness.command as CommandProbe };
  return { case: undefined } as NonNullable<ReadinessProbe["check"]>;
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
        kind: error.kind as ErrorInfo["kind"],
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

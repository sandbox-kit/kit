// Client runtime behavior is handwritten. Portable request and response data
// remain generated from protobuf. This module owns construction and create /
// attach / get / close; validation, errors, and the Sandbox handle live beside it.

import type { CreateOptions } from "#contracts/sandbox_data_contracts.js";
import { InfoFieldId } from "#contracts/field_paths.js";
import { ConfigFieldProvider } from "#contracts/config_paths.js";
import {
  ErrorKind,
  sandboxError,
  unsupported,
  withContext,
} from "#client/errors.js";
import {
  awaitWithSignal,
  operationOptions,
  operationSignal,
} from "#client/operation.js";
import {
  Sandbox,
  sandboxConstructorToken,
  type ClientLifecycle,
} from "#client/sandbox.js";
import type {
  Backend,
  ClientConfig,
  OperationOptions,
} from "#client/types.js";
import {
  validateConfig,
  validateCreateOptions,
} from "#client/validation.js";

export type {
  APIKeyCredentials,
  AuthInput,
  Backend,
  BearerTokenCredentials,
  ClientConfig,
  Config,
  DeleteCallOptions,
  OAuthCredentials,
  OperationOptions,
  Provider,
  Scope,
  StopCallOptions,
  TerminateCallOptions,
  TokenPairCredentials,
} from "#client/types.js";
export {
  ErrorKind,
  SandboxError,
  sandboxError,
} from "#client/errors.js";
export { Sandbox } from "#client/sandbox.js";
export {
  validateConfig,
  validateCreateOptions,
} from "#client/validation.js";

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
        field: ConfigFieldProvider,
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
        field: ConfigFieldProvider,
        message: "sandbox-kit: provider returned no backend",
      });
    if (this.#backend.name !== this.#providerName) {
      void this.#backend.close({}).catch(() => {});
      throw sandboxError({
        kind: ErrorKind.InvalidResponse,
        provider: this.#providerName,
        operation: "initialize",
        field: ConfigFieldProvider,
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
          field: InfoFieldId,
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
          field: InfoFieldId,
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
    const timeout =
      options.timeout ?? request.provisioning?.timeout ?? this.#timeout;
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      // Validate the wire projection first; providers receive the public request.
      validateCreateOptions(request, this.providerName, "create");
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
          field: InfoFieldId,
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
          field: ConfigFieldProvider,
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

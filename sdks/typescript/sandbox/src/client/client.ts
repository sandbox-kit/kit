// Client runtime behavior is handwritten. Portable request and response data
// remain generated from protobuf. This module owns construction and create /
// attach / get / close; validation, errors, and the Sandbox handle live beside it.

import type {
  CreateOptions,
  SandboxInfo,
} from "#contracts/sandbox_data_contracts.js";
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
import { Sandbox, sandboxConstructorToken } from "#client/sandbox.js";
import type { Backend, ClientConfig, OperationOptions } from "#client/types.js";
import { validateConfig, validateCreateOptions } from "#client/validation.js";

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
export { ErrorKind, SandboxError, sandboxError } from "#client/errors.js";
export { Sandbox } from "#client/sandbox.js";
export { validateConfig, validateCreateOptions } from "#client/validation.js";

/**
 * Creates and looks up sandboxes for one provider.
 * `close` closes your connection to that provider. The sandboxes keep running.
 * A later call on a handle returns the provider's result.
 */
export class Client implements AsyncDisposable {
  readonly #backend: Backend;
  readonly #providerName: string;
  readonly #timeout: number | undefined;
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
  /** `"daytona"` or `"modal"`. */
  get providerName(): string {
    return this.#providerName;
  }
  /**
   * Returns a handle for an existing sandbox ID without looking it up.
   * Modal supports attach. Daytona throws `ErrorKind.Unsupported`. Use `get`.
   */
  async attach(id: string, options: OperationOptions = {}): Promise<Sandbox> {
    return this.#openById("attach", id, options, (fetchOptions) => {
      if (!this.#backend.attach) unsupported(this.providerName, "attach");
      return this.#backend.attach(id, fetchOptions);
    });
  }
  /**
   * Looks up a sandbox by ID and returns a handle the provider has verified.
   * Daytona supports get. Modal throws `ErrorKind.Unsupported`. Use `attach`.
   */
  async get(id: string, options: OperationOptions = {}): Promise<Sandbox> {
    return this.#openById("get", id, options, (fetchOptions) => {
      if (!this.#backend.refresh) unsupported(this.providerName, "get");
      return this.#backend.refresh(id, fetchOptions);
    });
  }
  async #openById(
    operation: "attach" | "get",
    id: string,
    options: OperationOptions,
    fetch: (options: OperationOptions) => Promise<SandboxInfo>,
  ): Promise<Sandbox> {
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      scoped = operationSignal(options, this.#timeout);
      if (!id.trim())
        throw sandboxError({
          kind: ErrorKind.InvalidArgument,
          provider: this.providerName,
          operation,
          field: InfoFieldId,
          message: "sandbox-kit: sandbox ID must not be blank",
        });
      const info = await awaitWithSignal(
        () => fetch(operationOptions(options, scoped?.signal)),
        scoped.signal,
      );
      return new Sandbox(sandboxConstructorToken, {
        backend: this.#backend,
        info: { ...info, id, provider: this.providerName },
        ...(this.#timeout === undefined ? {} : { timeout: this.#timeout }),
      });
    } catch (error) {
      throw withContext(error, this.providerName, operation);
    } finally {
      scoped?.dispose();
    }
  }
  /**
   * Creates one sandbox and returns its handle.
   * An empty request uses the provider default. Daytona uses its default
   * snapshot. Modal requires an image.
   */
  async create(
    request: CreateOptions = {},
    options: OperationOptions = {},
  ): Promise<Sandbox> {
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
      return new Sandbox(sandboxConstructorToken, {
        backend: this.#backend,
        info: { ...result.sandbox, provider: this.providerName },
        ...(this.#timeout === undefined ? {} : { timeout: this.#timeout }),
      });
    } catch (error) {
      throw withContext(error, this.providerName, "create");
    } finally {
      scoped?.dispose();
    }
  }
  /**
   * Closes your connection to the provider.
   * The sandboxes keep running. You can keep using a handle.
   */
  async close(options: OperationOptions = {}): Promise<void> {
    let scoped: ReturnType<typeof operationSignal> | undefined;
    try {
      scoped = operationSignal(options, this.#timeout);
      await awaitWithSignal(
        () => this.#backend.close(operationOptions(options, scoped?.signal)),
        scoped.signal,
      );
    } catch (error) {
      throw withContext(error, this.providerName, "close");
    } finally {
      scoped?.dispose();
    }
  }
  async [Symbol.asyncDispose](): Promise<void> {
    await this.close();
  }
}

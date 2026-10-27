// Sandbox handle runtime behavior is handwritten. Lifecycle methods dispatch
// through optional Backend methods.

import {
  SandboxState,
  type ArchiveOptions,
  type DeleteOptions,
  type ExitResult,
  LifecycleCapability,
  type PauseOptions,
  type SandboxInfo,
  type StartOptions,
  type StopOptions,
  type TerminateOptions,
} from "#contracts/sandbox_data_contracts.js";
import {
  ErrorKind,
  sandboxError,
  unsupported,
  withContext,
} from "#client/errors.js";
import {
  awaitWithSignal,
  clone,
  operationOptions,
  operationSignal,
  publicTimeoutOptions,
} from "#client/operation.js";
import type {
  Backend,
  DeleteCallOptions,
  OperationOptions,
  StopCallOptions,
  TerminateCallOptions,
} from "#client/types.js";

/** Coordinates sandbox operations with client closure. */
export interface ClientLifecycle {
  begin(operation: string): () => void;
  isClosed(): boolean;
}

export const sandboxConstructorToken: unique symbol = Symbol(
  "sandbox-kit.Sandbox",
);

/** Internal construction inputs for {@link Sandbox}; only Client may call. */
export interface SandboxInit {
  backend: Backend;
  owner: ClientLifecycle;
  info: SandboxInfo;
  timeout?: number;
}

export class Sandbox implements AsyncDisposable {
  readonly #backend: Backend;
  readonly #owner: ClientLifecycle;
  readonly #timeout: number | undefined;
  #info: SandboxInfo;
  #closed = false;
  #tail: Promise<void> = Promise.resolve();
  constructor(token: typeof sandboxConstructorToken, init: SandboxInit) {
    if (token !== sandboxConstructorToken)
      throw new TypeError("Sandbox handles are created by Client");
    this.#backend = init.backend;
    this.#owner = init.owner;
    this.#info = clone(init.info);
    this.#timeout = init.timeout;
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
    return providerCapabilities(this.#backend);
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
        return this.#storeInfo(info);
      },
    );
  }
  async start(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "start",
      this.#backend.start,
      publicTimeoutOptions(options),
      options,
    );
  }
  async stop(options: StopCallOptions = {}): Promise<SandboxInfo> {
    const { force, ...call } = options;
    return this.#update(
      "stop",
      this.#backend.stop,
      {
        ...publicTimeoutOptions(call),
        ...(force === undefined ? {} : { force }),
      },
      call,
    );
  }
  async pause(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "pause",
      this.#backend.pause,
      publicTimeoutOptions(options),
      options,
    );
  }
  async terminate(options: TerminateCallOptions = {}): Promise<ExitResult> {
    const { wait, ...call } = options;
    const request: TerminateOptions = {
      ...publicTimeoutOptions(call),
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
              { ...request },
              operationOptions(call, signal),
            ),
          signal,
        );
        this.#assignInfo({
          ...this.#info,
          state: SandboxState.Terminated,
          terminal: true,
        });
        return result;
      },
    );
  }
  async archive(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "archive",
      this.#backend.archive,
      publicTimeoutOptions(options),
      options,
    );
  }
  async delete(options: DeleteCallOptions = {}): Promise<SandboxInfo> {
    const { wait, ...call } = options;
    return this.#update(
      "delete",
      this.#backend.delete,
      {
        ...publicTimeoutOptions(call),
        ...(wait === undefined ? {} : { wait }),
      },
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
    request:
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
              { ...request },
              operationOptions(options, signal),
            ),
          signal,
        );
        return this.#storeInfo(info);
      },
    );
  }
  #storeInfo(info: SandboxInfo): SandboxInfo {
    const stored = clone(info);
    stored.provider = this.providerName;
    this.#assignInfo(stored);
    return clone(stored);
  }
  #assignInfo(info: SandboxInfo): void {
    this.#info = info;
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

function providerCapabilities(backend: Backend): LifecycleCapability[] {
  const values: LifecycleCapability[] = [];
  if (backend.refresh) values.push(LifecycleCapability.Refresh);
  if (backend.start) values.push(LifecycleCapability.Start);
  if (backend.stop) values.push(LifecycleCapability.Stop);
  if (backend.pause) values.push(LifecycleCapability.Pause);
  if (backend.terminate) values.push(LifecycleCapability.Terminate);
  if (backend.archive) values.push(LifecycleCapability.Archive);
  if (backend.delete) values.push(LifecycleCapability.Delete);
  return values;
}

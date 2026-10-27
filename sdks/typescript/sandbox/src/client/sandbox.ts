// Sandbox handle runtime behavior is handwritten. Lifecycle methods dispatch
// through optional Backend methods.

import {
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
import { unsupported, withContext } from "#client/errors.js";
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

export const sandboxConstructorToken: unique symbol = Symbol(
  "sandbox-kit.Sandbox",
);

/** Internal construction inputs for {@link Sandbox}; only Client may call. */
export interface SandboxInit {
  backend: Backend;
  info: SandboxInfo;
  timeout?: number;
}

export class Sandbox implements AsyncDisposable {
  readonly #backend: Backend;
  readonly #timeout: number | undefined;
  #info: SandboxInfo;
  #tail: Promise<void> = Promise.resolve();
  constructor(token: typeof sandboxConstructorToken, init: SandboxInit) {
    if (token !== sandboxConstructorToken)
      throw new TypeError("Sandbox handles are created by Client");
    this.#backend = init.backend;
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
    await this.#serialized("close", options, this.#timeout, async (signal) => {
      await awaitWithSignal(
        () =>
          this.#backend.closeSandbox(
            this.id,
            operationOptions(options, signal),
          ),
        signal,
      );
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
    this.#info = stored;
    return clone(stored);
  }
  async #serialized<T>(
    operation: string,
    options: OperationOptions,
    timeout: number | undefined,
    body: (signal: AbortSignal | undefined) => Promise<T>,
  ): Promise<T> {
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
      return await body(scoped.signal);
    } catch (error) {
      throw withContext(error, this.providerName, operation);
    } finally {
      scoped?.dispose();
      if (acquired) unlock();
      else void previous.then(unlock, unlock);
    }
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

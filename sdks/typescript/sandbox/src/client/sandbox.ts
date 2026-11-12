// Sandbox handle runtime behavior is handwritten. Lifecycle, filesystem, and
// process methods dispatch through optional Backend methods.

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
import { InfoFieldId, InfoFieldProvider } from "#contracts/field_paths.js";
import { ErrorKind, sandboxError, unsupported, withContext } from "#client/errors.js";
import {
  awaitWithSignal,
  clone,
  operationOptions,
  operationSignal,
  publicTimeoutOptions,
} from "#client/operation.js";
import { Filesystem } from "#client/filesystem.js";
import { Process } from "#client/process.js";
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

/**
 * Handle for one cloud sandbox.
 * Lifecycle methods change it. `filesystem` exposes the shared file
 * operations. `process` runs commands. `close` leaves the cloud sandbox as
 * it is, and you can keep using the handle.
 */
export class Sandbox implements AsyncDisposable {
  readonly #backend: Backend;
  readonly #timeout: number | undefined;
  /**
   * File operations for this sandbox: read, write, copy, make directory, list,
   * stat, and remove. Paths inside the sandbox are absolute.
   */
  readonly filesystem: Filesystem;
  /**
   * Command runner for this sandbox. Set exactly one of `argv` or `shell`.
   * A working directory, when set, must be absolute.
   */
  readonly process: Process;
  #info: SandboxInfo;
  #tail: Promise<void> = Promise.resolve();
  constructor(token: typeof sandboxConstructorToken, init: SandboxInit) {
    if (token !== sandboxConstructorToken)
      throw new TypeError("Sandbox handles are created by Client");
    this.#backend = init.backend;
    this.#info = clone(init.info);
    this.#timeout = init.timeout;
    const self = this;
    const host = {
      get providerName() {
        return self.providerName;
      },
      get id() {
        return self.id;
      },
      get backend() {
        return self.#backend;
      },
      run: <T>(
        operation: string,
        options: OperationOptions,
        body: (signal: AbortSignal | undefined) => Promise<T>,
      ): Promise<T> =>
        self.#serialized(
          operation,
          options,
          options.timeout ?? self.#timeout,
          body,
        ),
    };
    this.filesystem = new Filesystem(host);
    this.process = new Process(host);
  }
  /** Sandbox identifier. */
  get id(): string {
    return this.#info.id;
  }
  /** `"daytona"`, `"modal"`, or `"blaxel"`. */
  get providerName(): string {
    return this.#info.provider;
  }
  /** A copy of the sandbox metadata. Editing the result leaves the handle unchanged. */
  info(): SandboxInfo {
    return clone(this.#info);
  }
  /**
   * Lifecycle operations this provider implements.
   * File and process operations are not lifecycle capabilities; use
   * {@link filesystem} or {@link process}.
   */
  capabilities(): readonly LifecycleCapability[] {
    return providerCapabilities(this.#backend);
  }
  /**
   * Reports whether this provider implements a lifecycle `capability`.
   * File and process operations are checked by calling {@link filesystem} or
   * {@link process} methods, not `supports`.
   */
  supports(capability: LifecycleCapability): boolean {
    return this.capabilities().includes(capability);
  }
  /**
   * Loads current metadata from the provider and stores it on this handle.
   * Daytona and Blaxel support refresh. Modal throws `ErrorKind.Unsupported`.
   */
  async refresh(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#serialized(
      "refresh",
      options,
      this.#timeout,
      async (signal) => {
        const method = this.#backend.refresh;
        if (!method) unsupported(this.providerName, "refresh");
        return this.#storeInfo(
          "refresh",
          await method.call(
            this.#backend,
            this.id,
            operationOptions(options, signal),
          ),
        );
      },
    );
  }
  /**
   * Starts a stopped sandbox and stores the updated metadata on this handle.
   * An omitted `timeout` uses the client timeout as this call's deadline.
   * Daytona supports start. Modal and Blaxel throw `ErrorKind.Unsupported`.
   */
  async start(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "start",
      this.#backend.start,
      publicTimeoutOptions(options),
      options,
    );
  }
  /**
   * Stops a sandbox and stores the updated metadata on this handle.
   * Set `force` to request an immediate stop.
   * Daytona supports stop. Modal and Blaxel throw `ErrorKind.Unsupported`.
   */
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
  /**
   * Pauses a sandbox and stores the updated metadata on this handle.
   * Daytona accepts or rejects the sandbox. Modal and Blaxel throw
   * `ErrorKind.Unsupported`.
   */
  async pause(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "pause",
      this.#backend.pause,
      publicTimeoutOptions(options),
      options,
    );
  }
  /**
   * Ends a sandbox. The result includes `exitCode` when `wait` is set.
   * An omitted `timeout` uses the client timeout as this call's deadline.
   * Modal waits on that deadline and does not take a separate terminate duration.
   * Modal supports terminate. Daytona and Blaxel throw `ErrorKind.Unsupported`.
   * Use `delete` or `stop` on Daytona; `delete` on Blaxel.
   */
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
        return method.call(
          this.#backend,
          this.id,
          { ...request },
          operationOptions(call, signal),
        );
      },
    );
  }
  /**
   * Archives a sandbox and stores the updated metadata on this handle.
   * Daytona archive does not take a separate duration. The call still uses
   * `timeout` as its deadline.
   * Daytona accepts or rejects the sandbox. Modal and Blaxel throw
   * `ErrorKind.Unsupported`.
   */
  async archive(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(
      "archive",
      this.#backend.archive,
      publicTimeoutOptions(options),
      options,
    );
  }
  /**
   * Deletes the cloud sandbox and stores the metadata the provider returns.
   * Daytona and Blaxel honor `wait` to block until deletion finishes. Blaxel
   * without `wait` returns `deleted` / `terminal: true` as soon as delete is
   * accepted. `close` leaves the sandbox in place. Modal throws
   * `ErrorKind.Unsupported` — use `terminate`.
   */
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
  /**
   * Closes this handle. The sandbox stays as it is, and you can keep using the handle.
   * Closing a Daytona or Blaxel handle does not contact the provider. Closing a
   * Modal handle does not terminate the sandbox.
   */
  async close(options: OperationOptions = {}): Promise<void> {
    await this.#serialized("close", options, this.#timeout, (signal) =>
      this.#backend.closeSandbox(this.id, operationOptions(options, signal)),
    );
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
        return this.#storeInfo(
          operation,
          await method.call(
            this.#backend,
            this.id,
            { ...request },
            operationOptions(options, signal),
          ),
        );
      },
    );
  }
  #storeInfo(operation: string, info: SandboxInfo): SandboxInfo {
    if (!info?.id || info.id !== this.id)
      throw sandboxError({
        kind: ErrorKind.InvalidResponse,
        provider: this.providerName,
        operation,
        field: InfoFieldId,
        message: "sandbox-kit: provider returned a different sandbox identity",
      });
    if (info.provider && info.provider !== this.providerName)
      throw sandboxError({
        kind: ErrorKind.InvalidResponse,
        provider: this.providerName,
        operation,
        field: InfoFieldProvider,
        message: "sandbox-kit: provider returned inconsistent sandbox metadata",
      });
    const stored = clone(info);
    stored.provider = this.providerName;
    this.#info = stored;
    return clone(stored);
  }
  // One active provider call per handle. Abort/timeout still rejects the
  // caller promptly (providers often ignore AbortSignal). The gate stays held
  // until provider work settles so a follow-up call cannot overlap it.
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
    let work: Promise<T> | undefined;
    try {
      scoped = operationSignal(options, timeout);
      await awaitWithSignal(previous, scoped.signal);
      acquired = true;
      work = body(scoped.signal);
      // Race the already-started provider promise so abort rejects promptly
      // without dropping that promise (unlock still waits on it below).
      return await awaitWithSignal(work, scoped.signal);
    } catch (error) {
      throw withContext(error, this.providerName, operation);
    } finally {
      scoped?.dispose();
      if (!acquired) {
        void previous.then(unlock, unlock);
      } else {
        // Do not await work here: that would delay abort/timeout rejection when
        // the provider ignores the signal. Unlock when the native call settles.
        void Promise.resolve(work).then(unlock, unlock);
      }
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

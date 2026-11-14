import {
  ErrorKind,
  FileFieldParents,
  SandboxState,
  type FileInfo,
  type OperationOptions,
  type Provider,
  type SandboxInfo,
} from "@sandbox-kit/sandbox";
import type {
  Backend,
  Config,
  CreateOptions,
  CreateResult,
  DeleteOptions,
  ExecOptions,
  ExecResult,
  FileCopyOptions,
  FileReadOptions,
  FileWriteOptions,
  ListDirectoryOptions,
  MakeDirectoryOptions,
  RemoveOptions,
  StatOptions,
} from "@sandbox-kit/sandbox";
import {
  blaxelCreateTimeoutSeconds,
  mapClientConfig,
  mapCreateRequest,
  mapSandboxInfo,
} from "#mappings/mappings.js";
import { mapBlaxelError, providerError } from "#errors.js";
import { BackendBase, ProviderBase } from "./provider.js";
import {
  captureProviderState,
  type ProviderState,
} from "#config/provider_config.js";
import { blaxelStat, directoryInfo, fileInfo } from "./filesystem.js";
import {
  ensureLocalParent,
  mapCallerPathError,
  readLocalFile,
  requireLocalDestination,
} from "./local.js";
import {
  defaultNativeFactory,
  type NativeBlaxel,
  type NativeFactory,
  type NativeSandboxHandle,
} from "./native.js";
import {
  execTimeoutSeconds,
  finishProcess,
  processRequest,
} from "./process.js";

/** @internal Exported for adapter boundary tests, but not from the package root. */
export class BlaxelProvider extends ProviderBase {
  constructor(
    private readonly nativeFactory: NativeFactory = defaultNativeFactory,
  ) {
    super();
  }

  initialize(config: Config): Backend {
    try {
      const nativeConfig = mapClientConfig(config);
      const native = this.nativeFactory(nativeConfig);
      return new BlaxelBackend(native, captureProviderState(config));
    } catch (error) {
      throw mapBlaxelError(error, "initialize");
    }
  }
}

interface WithSandboxArgs<T> {
  operation: string;
  id: string;
  options: OperationOptions;
  body: (sandbox: NativeSandboxHandle) => Promise<T>;
}

class BlaxelBackend extends BackendBase implements Backend {
  constructor(
    private readonly native: NativeBlaxel,
    private readonly config: ProviderState,
  ) {
    super();
  }

  async create(
    request: CreateOptions,
    options: OperationOptions,
  ): Promise<CreateResult> {
    try {
      throwIfAborted(options.signal);
      const plan = mapCreateRequest(request, this.config);
      throwIfAborted(options.signal);
      const timeoutSeconds = blaxelCreateTimeoutSeconds(
        request.provisioning?.timeout,
      );
      const sandbox = await this.native.create(
        plan.params,
        timeoutSeconds === undefined ? undefined : { timeout: timeoutSeconds },
      );
      await deleteCreatedSandboxAfterAbort(
        this.native,
        sandbox,
        request.name,
        options.signal,
      );
      return { sandbox: mapSandboxInfo(sandbox) };
    } catch (error) {
      throw mapBlaxelError(error, "create");
    }
  }

  async refresh(id: string, options: OperationOptions): Promise<SandboxInfo> {
    return this.#withSandbox({
      operation: "refresh",
      id,
      options,
      body: async (sandbox) => mapSandboxInfo(sandbox, "refresh"),
    });
  }

  async delete(
    id: string,
    request: DeleteOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    try {
      throwIfAborted(options.signal);
      const deleted = await this.native.delete(id);
      throwIfAborted(options.signal);
      if (request.wait === true) {
        // Blaxel has no native wait-for-delete; Kit polls until gone.
        return await this.#waitUntilDeleted(id, options.signal);
      }
      // DELETE accepted; Blaxel may still finish asynchronously.
      return {
        ...mapSandboxInfo(deleted, "delete"),
        state: SandboxState.Deleted,
        terminal: true,
        providerState: "TERMINATED",
      };
    } catch (error) {
      throw mapBlaxelError(error, "delete");
    }
  }

  async readFile(
    id: string,
    path: string,
    _request: FileReadOptions,
    options: OperationOptions,
  ): Promise<Uint8Array> {
    return this.#withSandbox({
      operation: "read_file",
      id,
      options,
      body: async (sandbox) => {
        const blob = await sandbox.fs.readBinary(path);
        return new Uint8Array(await blob.arrayBuffer());
      },
    });
  }

  async writeFile(
    id: string,
    path: string,
    data: Uint8Array,
    _request: FileWriteOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#withSandbox({
      operation: "write_file",
      id,
      options,
      body: async (sandbox) => {
        await sandbox.fs.writeBinary(path, data);
      },
    });
  }

  async copyFromLocal(
    id: string,
    localPath: string,
    remotePath: string,
    _request: FileCopyOptions,
    options: OperationOptions,
  ): Promise<void> {
    const data = await readLocalFile(localPath, "copy_from_local");
    await this.#withSandbox({
      operation: "copy_from_local",
      id,
      options,
      body: async (sandbox) => {
        await sandbox.fs.writeBinary(remotePath, data);
      },
    });
  }

  async copyToLocal(
    id: string,
    remotePath: string,
    localPath: string,
    _request: FileCopyOptions,
    options: OperationOptions,
  ): Promise<void> {
    await requireLocalDestination(localPath, "copy_to_local");
    await ensureLocalParent(localPath, "copy_to_local");
    await this.#withSandbox({
      operation: "copy_to_local",
      id,
      options,
      body: async (sandbox) => {
        try {
          await sandbox.fs.download(remotePath, localPath);
        } catch (error) {
          mapCallerPathError(error, "copy_to_local", localPath);
        }
      },
    });
  }

  async makeDirectory(
    id: string,
    path: string,
    request: MakeDirectoryOptions,
    options: OperationOptions,
  ): Promise<void> {
    if (request.parents === false) {
      providerError({
        kind: ErrorKind.Unsupported,
        operation: "make_directory",
        field: FileFieldParents,
        message: "parents false is not available",
      });
    }
    await this.#withSandbox({
      operation: "make_directory",
      id,
      options,
      body: async (sandbox) => {
        await sandbox.fs.mkdir(path);
      },
    });
  }

  async listDirectory(
    id: string,
    path: string,
    _request: ListDirectoryOptions,
    options: OperationOptions,
  ): Promise<FileInfo[]> {
    return this.#withSandbox({
      operation: "list_directory",
      id,
      options,
      body: async (sandbox) => {
        const listing = await sandbox.fs.ls(path);
        return [
          ...listing.subdirectories.map((entry) => directoryInfo(entry)),
          ...listing.files.map((entry) => fileInfo(entry, "list_directory")),
        ];
      },
    });
  }

  async stat(
    id: string,
    path: string,
    _request: StatOptions,
    options: OperationOptions,
  ): Promise<FileInfo> {
    return this.#withSandbox({
      operation: "stat",
      id,
      options,
      body: async (sandbox) => blaxelStat(sandbox.fs, path),
    });
  }

  async remove(
    id: string,
    path: string,
    request: RemoveOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#withSandbox({
      operation: "remove",
      id,
      options,
      body: async (sandbox) => {
        await sandbox.fs.rm(path, request.recursive === true);
      },
    });
  }

  async run(
    id: string,
    request: ExecOptions,
    options: OperationOptions,
  ): Promise<ExecResult> {
    return this.#withSandbox({
      operation: "run",
      id,
      options,
      body: async (sandbox) => {
        const timeout = execTimeoutSeconds(request.timeout);
        const started = await sandbox.process.exec(
          processRequest(request, timeout),
        );
        return finishProcess({
          process: sandbox.process,
          started,
          timeout,
          signal: options.signal,
        });
      },
    });
  }

  async #withSandbox<T>(args: WithSandboxArgs<T>): Promise<T> {
    let sandbox: NativeSandboxHandle;
    try {
      throwIfAborted(args.options.signal);
      sandbox = await this.native.get(args.id);
      throwIfAborted(args.options.signal);
    } catch (error) {
      // A 404 here identifies the sandbox, not the requested filesystem path.
      throw mapBlaxelError(error, args.operation, false);
    }
    try {
      return await args.body(sandbox);
    } catch (error) {
      throw mapBlaxelError(error, args.operation);
    }
  }

  /** Poll get until NotFound or Deleted. Do not use the DELETE body as done. */
  async #waitUntilDeleted(
    id: string,
    signal: AbortSignal | undefined,
  ): Promise<SandboxInfo> {
    for (;;) {
      throwIfAborted(signal);
      await delay(1_000, signal);
      try {
        const info = mapSandboxInfo(await this.native.get(id), "delete");
        if (info.state === SandboxState.Deleted) {
          return { ...info, terminal: true, providerState: "TERMINATED" };
        }
      } catch (error) {
        const mapped = mapBlaxelError(error, "delete");
        if (mapped.kind === ErrorKind.NotFound) {
          return {
            id,
            provider: "blaxel",
            state: SandboxState.Deleted,
            terminal: true,
            providerState: "TERMINATED",
          };
        }
        throw mapped;
      }
    }
  }

  async closeSandbox(_id: string, _options: OperationOptions): Promise<void> {}

  async close(options: OperationOptions): Promise<void> {
    throwIfAborted(options.signal);
  }
}

function throwIfAborted(signal: AbortSignal | undefined): void {
  signal?.throwIfAborted();
}

// @blaxel/core create has no AbortSignal. If Kit's outer race has already
// returned cancellation, delete the eventually-created named sandbox so it
// does not become an orphaned resource.
async function deleteCreatedSandboxAfterAbort(
  native: NativeBlaxel,
  sandbox: NativeSandboxHandle,
  requestedName: string | undefined,
  signal: AbortSignal | undefined,
): Promise<void> {
  if (!signal?.aborted) return;
  const name = sandbox.metadata?.name ?? requestedName;
  if (name) {
    try {
      await native.delete(name);
    } catch {
      // Preserve the caller's cancellation; a cleanup failure cannot safely
      // replace it with a successful create result.
    }
  }
  throwIfAborted(signal);
}

function delay(ms: number, signal: AbortSignal | undefined): Promise<void> {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) {
      reject(
        signal.reason ??
          new DOMException("The operation was aborted", "AbortError"),
      );
      return;
    }
    const timer = setTimeout(() => {
      signal?.removeEventListener("abort", onAbort);
      resolve();
    }, ms);
    const onAbort = () => {
      clearTimeout(timer);
      reject(
        signal?.reason ??
          new DOMException("The operation was aborted", "AbortError"),
      );
    };
    signal?.addEventListener("abort", onAbort, { once: true });
  });
}

/**
 * Selects Blaxel as the provider for the shared Sandbox Kit client.
 *
 * Create needs `name`, an image source, and client `region`. Set
 * `auth.apiKey` and `scope.workspace`, or omit them to use `BL_API_KEY` /
 * `BL_WORKSPACE` from the environment. Set `memoryMib` for memory size; Kit
 * converts MiB to Blaxel decimal MB, rounding up so the request is not
 * under-provisioned. Volumes are not part of create. In TypeScript, one process
 * accepts one Blaxel configuration; a conflicting client configuration is
 * rejected.
 */
export function blaxel(): Provider {
  return new BlaxelProvider();
}

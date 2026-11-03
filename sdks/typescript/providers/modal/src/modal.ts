import {
  ModalClient,
  type FileInfo as ModalFileInfo,
  type ModalClientParams,
  type Sandbox as ModalSandbox,
} from "modal";
import {
  ErrorKind,
  ExecFieldTimeout,
  FileType,
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
  ExecOptions,
  ExecResult,
  ExitResult,
  FileCopyOptions,
  FileReadOptions,
  FileWriteOptions,
  ListDirectoryOptions,
  MakeDirectoryOptions,
  RemoveOptions,
  StatOptions,
  TerminateOptions,
} from "@sandbox-kit/sandbox";
import {
  mapClientConfig,
  mapCreateRequest,
  mapSandboxInfo,
} from "#mappings/mappings.js";
import {
  mapModalError,
  providerError,
  requireLocalDestination,
  requireLocalFile,
  translateLocalCopyError,
} from "#errors.js";
import { BackendBase, ProviderBase } from "./provider.js";
import {
  captureProviderState,
  type ProviderState,
} from "#config/provider_config.js";

interface NativeModal {
  readonly apps: Pick<ModalClient["apps"], "fromName">;
  readonly images: Pick<ModalClient["images"], "fromRegistry">;
  readonly proxies: Pick<ModalClient["proxies"], "fromName">;
  readonly sandboxes: Pick<ModalClient["sandboxes"], "create" | "fromId">;
  readonly secrets: Pick<ModalClient["secrets"], "fromName">;
  close(): void;
}
type NativeFactory = (config: ModalClientParams) => NativeModal;

/** @internal Exported for adapter boundary tests, but not from the package root. */
export class ModalProvider extends ProviderBase {
  constructor(
    private readonly nativeFactory: NativeFactory = (config) =>
      new ModalClient(config),
  ) {
    super();
  }
  initialize(config: Config): Backend {
    try {
      return new ModalBackend(
        this.nativeFactory(mapClientConfig(config)),
        captureProviderState(config),
      );
    } catch (error) {
      throw mapModalError(error, "initialize");
    }
  }
}

class ModalBackend extends BackendBase implements Backend {
  constructor(
    private readonly native: NativeModal,
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
      const app = await this.native.apps.fromName(plan.app, {
        ...(plan.environment === undefined
          ? {}
          : { environment: plan.environment }),
      });
      throwIfAborted(options.signal);
      const secrets = await Promise.all(
        plan.secretNames.map((name) =>
          this.native.secrets.fromName(name, {
            ...(plan.environment === undefined
              ? {}
              : { environment: plan.environment }),
          }),
        ),
      );
      throwIfAborted(options.signal);
      const proxy =
        plan.proxyName === undefined
          ? undefined
          : await this.native.proxies.fromName(plan.proxyName, {
              ...(plan.environment === undefined
                ? {}
                : { environment: plan.environment }),
            });
      throwIfAborted(options.signal);
      const image = this.native.images.fromRegistry(plan.image);
      const sandbox = await this.native.sandboxes.create(app, image, {
        ...plan.params,
        ...(secrets.length === 0 ? {} : { secrets }),
        ...(proxy === undefined ? {} : { proxy }),
      });
      if (plan.waitUntilReady) {
        throwIfAborted(options.signal);
        const readyTimeoutMs = plan.readyTimeoutMs ?? options.timeout;
        await sandbox.waitUntilReady(readyTimeoutMs);
      }
      return { sandbox: mapSandboxInfo(sandbox, request) };
    } catch (error) {
      throw mapModalError(error, "create");
    }
  }

  async attach(id: string, options: OperationOptions): Promise<SandboxInfo> {
    try {
      throwIfAborted(options.signal);
      return mapSandboxInfo(await this.native.sandboxes.fromId(id));
    } catch (error) {
      throw mapModalError(error, "attach");
    }
  }

  async terminate(
    id: string,
    request: TerminateOptions,
    options: OperationOptions,
  ): Promise<ExitResult> {
    try {
      throwIfAborted(options.signal);
      const sandbox = await this.native.sandboxes.fromId(id);
      throwIfAborted(options.signal);
      await sandbox.terminate();
      if (request.wait !== true) return {};
      return { exitCode: await waitForExit(sandbox, options.signal) };
    } catch (error) {
      throw mapModalError(error, "terminate");
    }
  }

  async readFile(
    id: string,
    path: string,
    _request: FileReadOptions,
    options: OperationOptions,
  ): Promise<Uint8Array> {
    return this.#file("read_file", id, options, async (sandbox) =>
      sandbox.filesystem.readBytes(path),
    );
  }

  async writeFile(
    id: string,
    path: string,
    data: Uint8Array,
    _request: FileWriteOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#file("write_file", id, options, async (sandbox) => {
      await sandbox.filesystem.writeBytes(data, path);
    });
  }

  async copyFromLocal(
    id: string,
    localPath: string,
    remotePath: string,
    _request: FileCopyOptions,
    options: OperationOptions,
  ): Promise<void> {
    await requireLocalFile(localPath, "copy_from_local");
    await this.#file("copy_from_local", id, options, async (sandbox) => {
      try {
        await sandbox.filesystem.copyFromLocal(localPath, remotePath);
      } catch (error) {
        throw translateLocalCopyError(error, "copy_from_local");
      }
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
    await this.#file("copy_to_local", id, options, async (sandbox) => {
      try {
        await sandbox.filesystem.copyToLocal(remotePath, localPath);
      } catch (error) {
        throw translateLocalCopyError(error, "copy_to_local");
      }
    });
  }

  async makeDirectory(
    id: string,
    path: string,
    request: MakeDirectoryOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#file("make_directory", id, options, async (sandbox) => {
      await sandbox.filesystem.makeDirectory(path, {
        createParents: request.parents ?? true,
      });
    });
  }

  async listDirectory(
    id: string,
    path: string,
    _request: ListDirectoryOptions,
    options: OperationOptions,
  ): Promise<FileInfo[]> {
    return this.#file("list_directory", id, options, async (sandbox) => {
      const entries = await sandbox.filesystem.listFiles(path);
      return entries.map((entry) => modalFileInfo(entry));
    });
  }

  async stat(
    id: string,
    path: string,
    _request: StatOptions,
    options: OperationOptions,
  ): Promise<FileInfo> {
    return this.#file("stat", id, options, async (sandbox) =>
      modalFileInfo(await sandbox.filesystem.stat(path)),
    );
  }

  async remove(
    id: string,
    path: string,
    request: RemoveOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#file("remove", id, options, async (sandbox) => {
      await sandbox.filesystem.remove(path, {
        recursive: request.recursive === true,
      });
    });
  }

  async run(
    id: string,
    request: ExecOptions,
    options: OperationOptions,
  ): Promise<ExecResult> {
    return this.#file("run", id, options, async (sandbox) => {
      const command = modalCommand(request);
      const timeoutMs = execTimeoutMs(request.timeout);
      const process = await sandbox.exec(command, {
        mode: "binary",
        stdout: "pipe",
        stderr: "pipe",
        // Modal rejects an explicit timeoutMs of 0; omit the field for no limit.
        ...(timeoutMs === undefined ? {} : { timeoutMs }),
        ...(request.workingDirectory === undefined
          ? {}
          : { workdir: request.workingDirectory }),
        ...(request.environment === undefined
          ? {}
          : { env: request.environment }),
      });
      await process.closeStdin();
      const [stdout, stderr, exitCode] = await Promise.all([
        process.stdout.readBytes(),
        process.stderr.readBytes(),
        process.wait(),
      ]);
      return { exitCode, stdout, stderr };
    });
  }

  async #file<T>(
    operation: string,
    id: string,
    options: OperationOptions,
    body: (sandbox: ModalSandbox) => Promise<T>,
  ): Promise<T> {
    try {
      throwIfAborted(options.signal);
      const sandbox = await this.native.sandboxes.fromId(id);
      throwIfAborted(options.signal);
      return await body(sandbox);
    } catch (error) {
      throw mapModalError(error, operation);
    }
  }

  async closeSandbox(id: string, options: OperationOptions): Promise<void> {
    try {
      throwIfAborted(options.signal);
      (await this.native.sandboxes.fromId(id)).detach();
    } catch (error) {
      throw mapModalError(error, "close");
    }
  }

  async close(options: OperationOptions): Promise<void> {
    try {
      throwIfAborted(options.signal);
      this.native.close();
    } catch (error) {
      throw mapModalError(error, "close");
    }
  }
}

function modalFileInfo(info: ModalFileInfo): FileInfo {
  const type =
    info.type === "directory"
      ? FileType.Directory
      : info.type === "symlink"
        ? FileType.Symlink
        : FileType.File;
  return {
    name: info.name,
    path: info.path,
    type,
    size: BigInt(Math.trunc(info.size)),
    mode: (info.mode & 0o7777).toString(8).padStart(4, "0"),
    ...(info.permissions ? { permissions: info.permissions } : {}),
    ...(info.owner ? { owner: info.owner } : {}),
    ...(info.group ? { group: info.group } : {}),
    ...(info.symlinkTarget ? { symlinkTarget: info.symlinkTarget } : {}),
    ...(Number.isFinite(info.modifiedTime)
      ? { modifiedAt: new Date(info.modifiedTime * 1000) }
      : {}),
  };
}

function throwIfAborted(signal: AbortSignal | undefined): void {
  signal?.throwIfAborted();
}

/** Poll until the sandbox exits. The pinned Modal wait() loop does not take a signal. */
async function waitForExit(
  sandbox: ModalSandbox,
  signal: AbortSignal | undefined,
): Promise<number> {
  for (;;) {
    throwIfAborted(signal);
    const code = await sandbox.poll();
    if (code !== null) return code;
    await abortableDelay(signal, 1_000);
  }
}

function abortableDelay(
  signal: AbortSignal | undefined,
  ms: number,
): Promise<void> {
  throwIfAborted(signal);
  return new Promise((resolve, reject) => {
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

function modalCommand(request: ExecOptions): string[] {
  if (request.shell !== undefined) return ["bash", "-c", request.shell];
  return [...(request.argv ?? [])];
}

/**
 * Maps Kit ExecOptions.timeout (ms) to Modal's timeoutMs.
 * Omitted and explicit zero become undefined so Kit omits timeoutMs (Modal
 * rejects an explicit 0). Non-zero values must be whole seconds.
 */
function execTimeoutMs(timeout: number | undefined): number | undefined {
  if (timeout === undefined || timeout === 0) return undefined;
  if (!Number.isFinite(timeout) || timeout < 0 || timeout % 1_000 !== 0) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "run",
      field: ExecFieldTimeout,
      message: "timeout must be a whole number of seconds",
    });
  }
  return timeout;
}

/** Selects Modal as the provider for the shared Sandbox Kit client. */
export function modal(): Provider {
  return new ModalProvider();
}

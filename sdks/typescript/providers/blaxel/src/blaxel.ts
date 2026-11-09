import { basename, dirname } from "node:path";
import {
  SandboxInstance,
  settings,
  type Config as NativeClientConfig,
  type Directory,
  type File as BlaxelFile,
  type SandboxCreateConfiguration,
  type SandboxCreateOptions,
} from "@blaxel/core";

type BlaxelSubdirectory = Directory["subdirectories"][number];
import {
  ErrorKind,
  ExecFieldTimeout,
  FileFieldParents,
  FileType,
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
  mapClientConfig,
  mapCreateRequest,
  mapSandboxInfo,
} from "#mappings/mappings.js";
import {
  ensureLocalParent,
  mapBlaxelError,
  providerError,
  readLocalFile,
  requireLocalDestination,
} from "#errors.js";
import { BackendBase, ProviderBase } from "./provider.js";
import {
  captureProviderState,
  type ProviderState,
} from "#config/provider_config.js";

interface NativeSandboxHandle {
  readonly metadata?: {
    readonly name?: string;
    readonly labels?: Readonly<Record<string, string>>;
  };
  readonly status?: unknown;
  readonly spec?: { readonly region?: string };
  readonly fs: {
    readBinary(path: string): Promise<Blob>;
    writeBinary(path: string, content: Uint8Array): Promise<unknown>;
    mkdir(path: string, permissions?: string): Promise<unknown>;
    ls(path: string): Promise<Directory>;
    rm(path: string, recursive?: boolean): Promise<unknown>;
    download(src: string, destinationPath: string): Promise<void>;
  };
  readonly process: {
    exec(request: {
      command: string;
      waitForCompletion?: boolean;
      timeout?: number;
      workingDir?: string;
      env?: Record<string, string>;
    }): Promise<{
      exitCode?: number;
      stdout?: string;
      stderr?: string;
    }>;
  };
  delete(): Promise<unknown>;
}

interface NativeBlaxel {
  configure(config: NativeClientConfig): void;
  create(
    params: SandboxCreateConfiguration,
    options?: SandboxCreateOptions,
  ): Promise<NativeSandboxHandle>;
  get(name: string): Promise<NativeSandboxHandle>;
  delete(name: string): Promise<unknown>;
}

type NativeFactory = (config: NativeClientConfig) => NativeBlaxel;

/** @internal Exported for adapter boundary tests, but not from the package root. */
export class BlaxelProvider extends ProviderBase {
  constructor(
    private readonly nativeFactory: NativeFactory = (config) => {
      settings.setConfig(config);
      return {
        configure(next: NativeClientConfig): void {
          settings.setConfig(next);
        },
        create(params, options) {
          return SandboxInstance.create(params, options);
        },
        get(name) {
          return SandboxInstance.get(name);
        },
        delete(name) {
          return SandboxInstance.delete(name);
        },
      };
    },
  ) {
    super();
  }

  initialize(config: Config): Backend {
    try {
      const nativeConfig = mapClientConfig(config);
      const native = this.nativeFactory(nativeConfig);
      native.configure(nativeConfig);
      return new BlaxelBackend(native, captureProviderState(config));
    } catch (error) {
      throw mapBlaxelError(error, "initialize");
    }
  }
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
      const sandbox = await this.native.create(plan.params, plan.options);
      return { sandbox: mapSandboxInfo(sandbox, request) };
    } catch (error) {
      throw mapBlaxelError(error, "create");
    }
  }

  async refresh(id: string, options: OperationOptions): Promise<SandboxInfo> {
    try {
      throwIfAborted(options.signal);
      const sandbox = await this.native.get(id);
      throwIfAborted(options.signal);
      return mapSandboxInfo(sandbox);
    } catch (error) {
      throw mapBlaxelError(error, "refresh");
    }
  }

  async delete(
    id: string,
    _request: DeleteOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    try {
      throwIfAborted(options.signal);
      const sandbox = await this.native.get(id);
      throwIfAborted(options.signal);
      const info = mapSandboxInfo(sandbox);
      await sandbox.delete();
      return {
        ...info,
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
    return this.#withSandbox("read_file", id, options, async (sandbox) => {
      const blob = await sandbox.fs.readBinary(path);
      return new Uint8Array(await blob.arrayBuffer());
    });
  }

  async writeFile(
    id: string,
    path: string,
    data: Uint8Array,
    _request: FileWriteOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#withSandbox("write_file", id, options, async (sandbox) => {
      await sandbox.fs.writeBinary(path, data);
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
    await this.#withSandbox("copy_from_local", id, options, async (sandbox) => {
      await sandbox.fs.writeBinary(remotePath, data);
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
    await this.#withSandbox("copy_to_local", id, options, async (sandbox) => {
      await sandbox.fs.download(remotePath, localPath);
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
    await this.#withSandbox("make_directory", id, options, async (sandbox) => {
      await sandbox.fs.mkdir(path);
    });
  }

  async listDirectory(
    id: string,
    path: string,
    _request: ListDirectoryOptions,
    options: OperationOptions,
  ): Promise<FileInfo[]> {
    return this.#withSandbox("list_directory", id, options, async (sandbox) => {
      const listing = await sandbox.fs.ls(path);
      return [
        ...listing.subdirectories.map((entry) => directoryInfo(entry)),
        ...listing.files.map((entry) => fileInfo(entry)),
      ];
    });
  }

  async stat(
    id: string,
    path: string,
    _request: StatOptions,
    options: OperationOptions,
  ): Promise<FileInfo> {
    return this.#withSandbox("stat", id, options, async (sandbox) => {
      const parent = dirname(path);
      const name = basename(path);
      if (name === "" || name === "/" || path === "/") {
        const listing = await sandbox.fs.ls(path);
        return {
          name: listing.name || "/",
          path: listing.path || path,
          type: FileType.Directory,
          size: 0n,
        };
      }
      try {
        const listing = await sandbox.fs.ls(parent === "." ? "/" : parent);
        const file = listing.files.find((entry) => entry.name === name);
        if (file) return fileInfo(file);
        const directory = listing.subdirectories.find(
          (entry) => entry.name === name,
        );
        if (directory) return directoryInfo(directory);
      } catch {
        // Fall through to treating the path as a directory.
      }
      const self = await sandbox.fs.ls(path);
      return {
        name: self.name || name,
        path: self.path || path,
        type: FileType.Directory,
        size: 0n,
      };
    });
  }

  async remove(
    id: string,
    path: string,
    request: RemoveOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#withSandbox("remove", id, options, async (sandbox) => {
      await sandbox.fs.rm(path, request.recursive === true);
    });
  }

  async run(
    id: string,
    request: ExecOptions,
    options: OperationOptions,
  ): Promise<ExecResult> {
    return this.#withSandbox("run", id, options, async (sandbox) => {
      const timeout = execTimeoutSeconds(request.timeout);
      const result = await sandbox.process.exec({
        command: blaxelCommand(request),
        waitForCompletion: true,
        // Blaxel: omit/0 means no auto-kill. Kit omit/0 map to 0.
        timeout: timeout ?? 0,
        ...(request.workingDirectory === undefined
          ? {}
          : { workingDir: request.workingDirectory }),
        ...(request.environment === undefined
          ? {}
          : { env: { ...request.environment } }),
      });
      return {
        exitCode: result.exitCode ?? 0,
        stdout: new TextEncoder().encode(result.stdout ?? ""),
        stderr: new TextEncoder().encode(result.stderr ?? ""),
      };
    });
  }

  async #withSandbox<T>(
    operation: string,
    id: string,
    options: OperationOptions,
    body: (sandbox: NativeSandboxHandle) => Promise<T>,
  ): Promise<T> {
    try {
      throwIfAborted(options.signal);
      const sandbox = await this.native.get(id);
      throwIfAborted(options.signal);
      return await body(sandbox);
    } catch (error) {
      throw mapBlaxelError(error, operation);
    }
  }

  async closeSandbox(_id: string, _options: OperationOptions): Promise<void> {}

  async close(options: OperationOptions): Promise<void> {
    throwIfAborted(options.signal);
  }
}

function fileInfo(info: BlaxelFile): FileInfo {
  const modified = Date.parse(info.lastModified);
  return {
    name: info.name,
    path: info.path,
    type: FileType.File,
    size: BigInt(Math.trunc(info.size)),
    ...(info.permissions ? { permissions: info.permissions } : {}),
    ...(info.owner ? { owner: info.owner } : {}),
    ...(info.group ? { group: info.group } : {}),
    ...(Number.isNaN(modified) ? {} : { modifiedAt: new Date(modified) }),
  };
}

function directoryInfo(info: BlaxelSubdirectory): FileInfo {
  return {
    name: info.name,
    path: info.path,
    type: FileType.Directory,
    size: 0n,
  };
}

function throwIfAborted(signal: AbortSignal | undefined): void {
  signal?.throwIfAborted();
}

function blaxelCommand(request: ExecOptions): string {
  if (request.shell !== undefined) return request.shell;
  return posixShellJoin(request.argv ?? []);
}

function posixShellJoin(argv: readonly string[]): string {
  return argv.map(posixQuote).join(" ");
}

function posixQuote(value: string): string {
  if (value === "") return "''";
  return "'" + value.replaceAll("'", `'"'"'`) + "'";
}

/**
 * Maps Kit ExecOptions.timeout (ms) to Blaxel's process timeout (seconds).
 * Omitted and explicit zero become undefined so the caller sends 0 (no limit).
 */
function execTimeoutSeconds(timeout: number | undefined): number | undefined {
  if (timeout === undefined || timeout === 0) return undefined;
  if (!Number.isFinite(timeout) || timeout < 0 || timeout % 1_000 !== 0) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "run",
      field: ExecFieldTimeout,
      message: "timeout must be a whole number of seconds",
    });
  }
  return timeout / 1_000;
}

/** Selects Blaxel as the provider for the shared Sandbox Kit client. */
export function blaxel(): Provider {
  return new BlaxelProvider();
}

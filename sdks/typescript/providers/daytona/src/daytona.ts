import { Buffer } from "node:buffer";
import {
  Daytona,
  DaytonaError,
  type DaytonaConfig,
  type FileDownloadResponse,
  type Sandbox as DaytonaSandbox,
} from "@daytona/sdk";
import {
  ConfigFieldTimeout,
  CreateFieldProvisioningTimeout,
  ErrorKind,
  ExecFieldTimeout,
  FileFieldParents,
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
  ArchiveOptions,
  DeleteOptions,
  ExecOptions,
  ExecResult,
  FileCopyOptions,
  FileReadOptions,
  FileWriteOptions,
  ListDirectoryOptions,
  MakeDirectoryOptions,
  PauseOptions,
  RemoveOptions,
  StartOptions,
  StatOptions,
  StopOptions,
} from "@sandbox-kit/sandbox";
import type { FileInfo as DaytonaFileInfo } from "@daytona/sdk";

import {
  mapClientConfig,
  mapCreateRequest,
  mapSandboxInfo,
} from "#mappings/mappings.js";
import {
  ensureLocalParent,
  mapDaytonaError,
  providerError,
  requireLocalDestination,
  requireLocalFile,
  translateLocalCopyError,
} from "#errors.js";
import { BackendBase, ProviderBase } from "./provider.js";

interface NativeDaytona {
  create(
    params: ReturnType<typeof mapCreateRequest>,
    options?: { readonly timeout?: number },
  ): Promise<DaytonaSandbox>;
  get(id: string): Promise<DaytonaSandbox>;
  [Symbol.asyncDispose](): Promise<void>;
}

type NativeFactory = (config: DaytonaConfig) => NativeDaytona;

/** @internal Exported from this module for boundary testing, but not from the package root. */
export class DaytonaProvider extends ProviderBase {
  constructor(
    private readonly nativeFactory: NativeFactory = (config) =>
      new Daytona(config),
  ) {
    super();
  }

  initialize(config: Config): Backend {
    try {
      return new DaytonaBackend(this.nativeFactory(mapClientConfig(config)));
    } catch (error) {
      throw mapDaytonaError(error, "initialize");
    }
  }
}

class DaytonaBackend extends BackendBase implements Backend {
  constructor(private readonly native: NativeDaytona) {
    super();
  }

  async create(
    request: CreateOptions,
    options: OperationOptions,
  ): Promise<CreateResult> {
    try {
      throwIfAborted(options.signal);
      const params = mapCreateRequest(request);
      throwIfAborted(options.signal);
      const timeout = wholeSeconds(
        options.timeout,
        "create",
        request.provisioning?.timeout !== undefined
          ? CreateFieldProvisioningTimeout
          : ConfigFieldTimeout,
      );
      const sandbox = await this.native.create(
        params,
        // Create treats an omitted timeout as 60 seconds. Zero is omitted too:
        // the Go SDK's create option expires immediately when the timeout is 0.
        timeout === undefined || timeout === 0 ? undefined : { timeout },
      );
      return { sandbox: mapSandboxInfo(sandbox, "create") };
    } catch (error) {
      throw mapDaytonaError(error, "create");
    }
  }

  async refresh(id: string, options: OperationOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "refresh", options, async (value) => {
      await value.refreshData();
    });
  }
  async start(
    id: string,
    request: StartOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    const timeout = wholeSeconds(request.timeout, "start", ConfigFieldTimeout);
    return this.withSandbox(id, "start", options, (value) =>
      value.start(timeout),
    );
  }
  async stop(
    id: string,
    request: StopOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    const timeout = wholeSeconds(request.timeout, "stop", ConfigFieldTimeout);
    return this.withSandbox(id, "stop", options, (value) =>
      value.stop(timeout, request.force),
    );
  }
  async pause(
    id: string,
    request: PauseOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    const timeout = wholeSeconds(request.timeout, "pause", ConfigFieldTimeout);
    return this.withSandbox(id, "pause", options, (value) =>
      value.pause(timeout),
    );
  }
  async archive(
    id: string,
    _request: ArchiveOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    return this.withSandbox(id, "archive", options, (value) => value.archive());
  }
  async delete(
    id: string,
    request: DeleteOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    const timeout = wholeSeconds(request.timeout, "delete", ConfigFieldTimeout);
    return this.withSandbox(id, "delete", options, (value) =>
      value.delete(timeout, request.wait),
    );
  }

  private async withSandbox(
    id: string,
    operation: string,
    options: OperationOptions,
    action: (sandbox: DaytonaSandbox) => Promise<void>,
  ): Promise<SandboxInfo> {
    try {
      throwIfAborted(options.signal);
      const sandbox = await this.native.get(id);
      throwIfAborted(options.signal);
      await action(sandbox);
      return mapSandboxInfo(sandbox, operation);
    } catch (error) {
      throw mapDaytonaError(error, operation);
    }
  }

  async readFile(
    id: string,
    path: string,
    request: FileReadOptions,
    options: OperationOptions,
  ): Promise<Uint8Array> {
    return this.#file(
      "read_file",
      id,
      options,
      request.timeout,
      async (sandbox, timeout) => {
        const seconds = fileSeconds(timeout);
        // downloadFile(path, 0) drops the 0, because 0 is not a string path,
        // and Daytona then applies its 30-minute default.
        if (!seconds) {
          const [file] = await sandbox.fs.downloadFiles([{ source: path }], 0);
          const body = file?.result;
          if (!file || file.error || !body) throw downloadFailure(file);
          return new Uint8Array(
            typeof body === "string" ? Buffer.from(body) : body,
          );
        }
        const data = await sandbox.fs.downloadFile(path, seconds);
        return new Uint8Array(data);
      },
    );
  }

  async writeFile(
    id: string,
    path: string,
    data: Uint8Array,
    request: FileWriteOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#file(
      "write_file",
      id,
      options,
      request.timeout,
      async (sandbox, timeout) => {
        await sandbox.fs.uploadFile(
          Buffer.from(data),
          path,
          fileSeconds(timeout),
        );
      },
    );
  }

  async copyFromLocal(
    id: string,
    localPath: string,
    remotePath: string,
    request: FileCopyOptions,
    options: OperationOptions,
  ): Promise<void> {
    await requireLocalFile(localPath, "copy_from_local");
    await this.#file(
      "copy_from_local",
      id,
      options,
      request.timeout,
      async (sandbox, timeout) => {
        await sandbox.fs.uploadFile(
          localPath,
          remotePath,
          fileSeconds(timeout),
        );
      },
    );
  }

  async copyToLocal(
    id: string,
    remotePath: string,
    localPath: string,
    request: FileCopyOptions,
    options: OperationOptions,
  ): Promise<void> {
    await requireLocalDestination(localPath, "copy_to_local");
    await ensureLocalParent(localPath, "copy_to_local");
    await this.#file(
      "copy_to_local",
      id,
      options,
      request.timeout,
      async (sandbox, timeout) => {
        try {
          await sandbox.fs.downloadFile(
            remotePath,
            localPath,
            fileSeconds(timeout),
          );
        } catch (error) {
          throw translateLocalCopyError(error, "copy_to_local");
        }
      },
    );
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
    await this.#file(
      "make_directory",
      id,
      options,
      undefined,
      async (sandbox) => {
        await sandbox.fs.createFolder(path, "0755");
      },
    );
  }

  async listDirectory(
    id: string,
    path: string,
    _request: ListDirectoryOptions,
    options: OperationOptions,
  ): Promise<FileInfo[]> {
    return this.#file(
      "list_directory",
      id,
      options,
      undefined,
      async (sandbox) => {
        const entries = await sandbox.fs.listFiles(path);
        return entries.map((entry) => daytonaFileInfo(entry));
      },
    );
  }

  async stat(
    id: string,
    path: string,
    _request: StatOptions,
    options: OperationOptions,
  ): Promise<FileInfo> {
    return this.#file("stat", id, options, undefined, async (sandbox) =>
      daytonaFileInfo(await sandbox.fs.getFileDetails(path)),
    );
  }

  async remove(
    id: string,
    path: string,
    request: RemoveOptions,
    options: OperationOptions,
  ): Promise<void> {
    await this.#file("remove", id, options, undefined, async (sandbox) => {
      await sandbox.fs.deleteFile(path, request.recursive === true);
    });
  }

  async run(
    id: string,
    request: ExecOptions,
    options: OperationOptions,
  ): Promise<ExecResult> {
    try {
      throwIfAborted(options.signal);
      const timeout = wholeSeconds(request.timeout, "run", ExecFieldTimeout);
      const sandbox = await this.native.get(id);
      throwIfAborted(options.signal);
      const command = daytonaCommand(request);
      // Omitted and explicit zero become 0 so Daytona adds no command limit
      // (its omitted default is 10 seconds).
      const seconds = timeout ?? 0;
      const response = await sandbox.process.executeCommand(
        command,
        request.workingDirectory,
        request.environment,
        seconds,
      );
      return {
        exitCode: response.exitCode,
        stdout: new TextEncoder().encode(response.result ?? ""),
        stderr: new Uint8Array(),
      };
    } catch (error) {
      throw mapDaytonaError(error, "run");
    }
  }

  async #file<T>(
    operation: string,
    id: string,
    options: OperationOptions,
    timeoutMs: number | undefined,
    body: (sandbox: DaytonaSandbox, timeout: number | undefined) => Promise<T>,
  ): Promise<T> {
    try {
      throwIfAborted(options.signal);
      const timeout = wholeSeconds(timeoutMs, operation, ConfigFieldTimeout);
      const sandbox = await this.native.get(id);
      throwIfAborted(options.signal);
      return await body(sandbox, timeout);
    } catch (error) {
      throw mapDaytonaError(error, operation);
    }
  }

  async closeSandbox(_id: string, _options: OperationOptions): Promise<void> {}

  async close(options: OperationOptions): Promise<void> {
    try {
      throwIfAborted(options.signal);
      await this.native[Symbol.asyncDispose]();
    } catch (error) {
      throw mapDaytonaError(error, "close");
    }
  }
}

function daytonaFileInfo(info: DaytonaFileInfo): FileInfo {
  const modified = fileDate(info);
  return {
    name: info.name,
    path: info.path || "",
    type: info.isDir ? FileType.Directory : FileType.File,
    size: BigInt(info.size),
    ...(info.mode ? { mode: info.mode } : {}),
    ...(info.permissions ? { permissions: info.permissions } : {}),
    ...(info.owner ? { owner: info.owner } : {}),
    ...(info.group ? { group: info.group } : {}),
    ...(modified ? { modifiedAt: modified } : {}),
  };
}

function fileDate(info: DaytonaFileInfo): Date | undefined {
  const raw = info.modifiedAt || info.modTime;
  if (!raw) return undefined;
  const parsed = Date.parse(raw);
  if (Number.isNaN(parsed)) return undefined;
  return new Date(parsed);
}

function throwIfAborted(signal: AbortSignal | undefined): void {
  signal?.throwIfAborted();
}

/**
 * Converts a Kit file timeout (ms) for Daytona upload/download.
 * Omitted and explicit zero become `0` so Daytona adds no transfer limit.
 * Daytona's native default for an omitted argument is 30 minutes.
 */
function fileSeconds(timeout: number | undefined): number {
  return timeout ?? 0;
}

function downloadFailure(file: FileDownloadResponse | undefined): DaytonaError {
  const details = file?.errorDetails;
  if (details) {
    return new DaytonaError(
      details.message,
      details.statusCode,
      undefined,
      details.code ?? details.errorCode,
      details.source,
    );
  }
  return new DaytonaError(file?.error ?? "download failed");
}

/** Native Daytona APIs take whole seconds. Lifecycle zero means no operation timeout. */
function wholeSeconds(
  value: number | undefined,
  operation: string,
  field: string,
): number | undefined {
  if (value === undefined) return undefined;
  if (!Number.isFinite(value) || value < 0 || value % 1_000 !== 0)
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation,
      field,
      message: "timeout must be a whole number of seconds",
    });
  return value / 1_000;
}

function daytonaCommand(request: ExecOptions): string {
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

/** Selects Daytona as the provider for the shared Sandbox Kit client. */
export function daytona(): Provider {
  return new DaytonaProvider();
}

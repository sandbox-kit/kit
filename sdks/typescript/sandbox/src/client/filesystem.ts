// Filesystem calls are handwritten. They validate paths, then dispatch
// through optional Backend methods.

import { FileFieldLocalPath, FileFieldPath } from "#contracts/field_paths.js";
import type {
  FileCopyOptions,
  FileInfo,
  FileReadOptions,
  FileWriteOptions,
  ListDirectoryOptions,
  MakeDirectoryOptions,
  RemoveOptions,
  StatOptions,
} from "#contracts/sandbox_data_contracts.js";
import { ErrorKind, sandboxError, unsupported } from "#client/errors.js";
import { awaitWithSignal, operationOptions } from "#client/operation.js";
import type { Backend, OperationOptions } from "#client/types.js";

type FileCall<T> = T & OperationOptions;

interface FilesystemHost {
  readonly providerName: string;
  readonly id: string;
  readonly backend: Backend;
  run<T>(
    operation: string,
    options: OperationOptions,
    body: (signal: AbortSignal | undefined) => Promise<T>,
  ): Promise<T>;
}

/** File operations for one sandbox. Paths inside the sandbox are absolute. */
export class Filesystem {
  constructor(private readonly host: FilesystemHost) {}

  /** Returns the bytes of the file at `path`. */
  readFile(
    path: string,
    options: FileCall<FileReadOptions> = {},
  ): Promise<Uint8Array> {
    return this.#call("read_file", path, options, (signal) => {
      const method = this.host.backend.readFile;
      if (!method) unsupported(this.host.providerName, "read_file");
      return awaitWithSignal(
        () =>
          method.call(
            this.host.backend,
            this.host.id,
            path,
            requestFields(options),
            operationOptions(options, signal),
          ),
        signal,
      );
    });
  }

  /** Replaces the file at `path` and creates missing parent directories. */
  writeFile(
    path: string,
    data: Uint8Array,
    options: FileCall<FileWriteOptions> = {},
  ): Promise<void> {
    return this.#call("write_file", path, options, (signal) => {
      const method = this.host.backend.writeFile;
      if (!method) unsupported(this.host.providerName, "write_file");
      return awaitWithSignal(
        () =>
          method.call(
            this.host.backend,
            this.host.id,
            path,
            data,
            requestFields(options),
            operationOptions(options, signal),
          ),
        signal,
      );
    });
  }

  /** Copies a file from the caller into the sandbox. */
  copyFromLocal(
    localPath: string,
    remotePath: string,
    options: FileCall<FileCopyOptions> = {},
  ): Promise<void> {
    return this.#copy(
      "copy_from_local",
      localPath,
      remotePath,
      options,
      (signal) => {
        const method = this.host.backend.copyFromLocal;
        if (!method) unsupported(this.host.providerName, "copy_from_local");
        return awaitWithSignal(
          () =>
            method.call(
              this.host.backend,
              this.host.id,
              localPath,
              remotePath,
              requestFields(options),
              operationOptions(options, signal),
            ),
          signal,
        );
      },
    );
  }

  /** Copies a sandbox file onto the caller's disk. */
  copyToLocal(
    remotePath: string,
    localPath: string,
    options: FileCall<FileCopyOptions> = {},
  ): Promise<void> {
    return this.#copy(
      "copy_to_local",
      localPath,
      remotePath,
      options,
      (signal) => {
        const method = this.host.backend.copyToLocal;
        if (!method) unsupported(this.host.providerName, "copy_to_local");
        return awaitWithSignal(
          () =>
            method.call(
              this.host.backend,
              this.host.id,
              remotePath,
              localPath,
              requestFields(options),
              operationOptions(options, signal),
            ),
          signal,
        );
      },
    );
  }

  /**
   * Creates the directory at `path`.
   * An omitted `parents` flag is true: missing parents are created.
   * Daytona rejects `false`.
   */
  makeDirectory(
    path: string,
    options: FileCall<MakeDirectoryOptions> = {},
  ): Promise<void> {
    return this.#call("make_directory", path, options, (signal) => {
      const method = this.host.backend.makeDirectory;
      if (!method) unsupported(this.host.providerName, "make_directory");
      return awaitWithSignal(
        () =>
          method.call(
            this.host.backend,
            this.host.id,
            path,
            directoryRequest(options),
            operationOptions(options, signal),
          ),
        signal,
      );
    });
  }

  /** Lists the entries in one directory. */
  listDirectory(
    path: string,
    options: FileCall<ListDirectoryOptions> = {},
  ): Promise<FileInfo[]> {
    return this.#call("list_directory", path, options, (signal) => {
      const method = this.host.backend.listDirectory;
      if (!method) unsupported(this.host.providerName, "list_directory");
      return awaitWithSignal(
        () =>
          method.call(
            this.host.backend,
            this.host.id,
            path,
            requestFields(options),
            operationOptions(options, signal),
          ),
        signal,
      );
    });
  }

  /** Returns metadata for one path. A symlink result describes the symlink. */
  stat(path: string, options: FileCall<StatOptions> = {}): Promise<FileInfo> {
    return this.#call("stat", path, options, (signal) => {
      const method = this.host.backend.stat;
      if (!method) unsupported(this.host.providerName, "stat");
      return awaitWithSignal(
        () =>
          method.call(
            this.host.backend,
            this.host.id,
            path,
            requestFields(options),
            operationOptions(options, signal),
          ),
        signal,
      );
    });
  }

  /**
   * Removes a file, or a directory when `recursive` is set.
   * An omitted `recursive` flag removes a file or an empty directory.
   */
  remove(path: string, options: FileCall<RemoveOptions> = {}): Promise<void> {
    return this.#call("remove", path, options, (signal) => {
      const method = this.host.backend.remove;
      if (!method) unsupported(this.host.providerName, "remove");
      return awaitWithSignal(
        () =>
          method.call(
            this.host.backend,
            this.host.id,
            path,
            removeRequest(options),
            operationOptions(options, signal),
          ),
        signal,
      );
    });
  }

  async #call<T>(
    operation: string,
    path: string,
    options: OperationOptions,
    body: (signal: AbortSignal | undefined) => Promise<T>,
  ): Promise<T> {
    absolutePath(this.host.providerName, operation, path);
    return this.host.run(operation, options, body);
  }

  async #copy<T>(
    operation: string,
    localPath: string,
    remotePath: string,
    options: OperationOptions,
    body: (signal: AbortSignal | undefined) => Promise<T>,
  ): Promise<T> {
    absolutePath(this.host.providerName, operation, remotePath);
    callerPath(this.host.providerName, operation, localPath);
    return this.host.run(operation, options, body);
  }
}

function absolutePath(provider: string, operation: string, path: string): void {
  if (path.length === 0 || !path.startsWith("/") || path.includes("\0")) {
    throw sandboxError({
      kind: ErrorKind.InvalidArgument,
      provider,
      operation,
      field: FileFieldPath,
      message: "sandbox-kit: sandbox path must be absolute",
    });
  }
}

function callerPath(provider: string, operation: string, path: string): void {
  if (path.length === 0 || path.includes("\0")) {
    throw sandboxError({
      kind: ErrorKind.InvalidArgument,
      provider,
      operation,
      field: FileFieldLocalPath,
      message: "sandbox-kit: local path is required",
    });
  }
}

function requestFields<T extends { timeout?: number }>(
  options: T & OperationOptions,
): T {
  const request = {} as T;
  if (options.timeout !== undefined) request.timeout = options.timeout;
  return request;
}

function directoryRequest(
  options: FileCall<MakeDirectoryOptions>,
): MakeDirectoryOptions {
  return {
    ...requestFields(options),
    ...(options.parents === undefined ? {} : { parents: options.parents }),
  };
}

function removeRequest(options: FileCall<RemoveOptions>): RemoveOptions {
  return {
    ...requestFields(options),
    ...(options.recursive === undefined
      ? {}
      : { recursive: options.recursive }),
  };
}

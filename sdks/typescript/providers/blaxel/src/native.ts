import {
  getFilesystemByPath,
  ResponseError,
  SandboxInstance,
  settings,
  type Config as NativeClientConfig,
  type Directory,
  type File as BlaxelFile,
  type SandboxCreateConfiguration,
  type SandboxFileSystem,
} from "@blaxel/core";
import { ErrorKind, FileFieldPath } from "@sandbox-kit/sandbox";
import type { BlaxelSandboxIdentity } from "#mappings/provider_mappings.js";
import { providerError } from "./errors.js";

export interface BlaxelProcessResult {
  readonly pid?: string;
  readonly name?: string;
  readonly status?: string;
  readonly exitCode?: number;
  readonly stdout?: string;
  readonly stderr?: string;
}

export interface ProcessExecRequest {
  command: string;
  waitForCompletion?: boolean;
  timeout?: number;
  workingDir?: string;
  env?: Record<string, string>;
}

export interface NativeSandboxHandle extends BlaxelSandboxIdentity {
  readonly fs: {
    readBinary(path: string): Promise<Blob>;
    writeBinary(path: string, content: Uint8Array): Promise<unknown>;
    mkdir(path: string, permissions?: string): Promise<unknown>;
    ls(path: string): Promise<Directory>;
    /** GET path: directory listing or file metadata (size, permissions, …). */
    get(path: string): Promise<Directory | BlaxelFile>;
    rm(path: string, recursive?: boolean): Promise<unknown>;
    download(src: string, destinationPath: string): Promise<void>;
  };
  readonly process: {
    exec(request: ProcessExecRequest): Promise<BlaxelProcessResult>;
    wait(
      identifier: string,
      options?: {
        maxWait?: number;
        signal?: AbortSignal;
      },
    ): Promise<BlaxelProcessResult>;
  };
}

export interface NativeBlaxel {
  configure(config: NativeClientConfig): void;
  create(
    params: SandboxCreateConfiguration,
    /** Blaxel create timeout in whole seconds (1..50). */
    options?: { readonly timeout?: number },
  ): Promise<NativeSandboxHandle>;
  get(name: string): Promise<NativeSandboxHandle>;
  delete(name: string): Promise<BlaxelSandboxIdentity>;
}

export type NativeFactory = (config: NativeClientConfig) => NativeBlaxel;

export function defaultNativeFactory(config: NativeClientConfig): NativeBlaxel {
  settings.setConfig(config);
  return {
    configure(next: NativeClientConfig): void {
      settings.setConfig(next);
    },
    async create(params, options) {
      return projectHandle(
        await SandboxInstance.create(
          params,
          options?.timeout === undefined
            ? undefined
            : { timeout: options.timeout },
        ),
      );
    },
    async get(name) {
      return projectHandle(await SandboxInstance.get(name));
    },
    async delete(name) {
      return (await SandboxInstance.delete(name)) as BlaxelSandboxIdentity;
    },
  };
}

/**
 * `@blaxel/core` types `SandboxInstance.sandbox` private and has no `state`
 * getter. Project a plain Kit handle (identity + wrapped fs/process) so
 * mapping never mutates the native instance or reaches into its private wire
 * model after projection. Wrap `fs` with `get` so Stat uses the same
 * filesystem GET as Go.
 */
export function projectHandle(instance: SandboxInstance): NativeSandboxHandle {
  const state = readWireRuntimeState(instance);
  return {
    metadata: instance.metadata,
    status: instance.status,
    ...(state !== undefined ? { state } : {}),
    spec: instance.spec,
    fs: wrapFilesystem(instance.fs),
    process: {
      exec: (request) =>
        instance.process.exec(request) as Promise<BlaxelProcessResult>,
      wait: (identifier, options) =>
        instance.process.wait(
          identifier,
          options,
        ) as Promise<BlaxelProcessResult>,
    },
  };
}

function wrapFilesystem(fs: SandboxFileSystem): NativeSandboxHandle["fs"] {
  return {
    readBinary: (path) => fs.readBinary(path),
    writeBinary: (path, content) => fs.writeBinary(path, content),
    mkdir: (path, permissions) => fs.mkdir(path, permissions),
    ls: (path) => fs.ls(path),
    rm: (path, recursive) => fs.rm(path, recursive),
    download: (src, destinationPath) => fs.download(src, destinationPath),
    get: (path) => filesystemGet(fs, path),
  };
}

async function filesystemGet(
  fs: SandboxFileSystem,
  path: string,
): Promise<Directory | BlaxelFile> {
  const { response, data, error } = await getFilesystemByPath({
    path: { path },
    baseUrl: fs.url,
    client: fs.client,
  });
  if (!response.ok || data == null) {
    throw new ResponseError(response, data, error);
  }
  if (isDirectoryListing(data)) return data;
  if (isFileMetadata(data)) {
    const { content: _content, ...meta } = data as BlaxelFile & {
      content?: string;
    };
    return meta;
  }
  // Match Go Stat: a non-listing GET body is treated as file metadata; when it
  // is not usable Kit reports InvalidResponse (not Unknown from a plain Error).
  providerError({
    kind: ErrorKind.InvalidResponse,
    operation: "stat",
    field: FileFieldPath,
    message: "SDK returned no file metadata",
  });
}

/** True when a filesystem GET body is a directory listing (shared with Stat). */
export function isDirectoryListing(data: unknown): data is Directory {
  return (
    !!data &&
    typeof data === "object" &&
    ("files" in data || "subdirectories" in data)
  );
}

function isFileMetadata(data: unknown): data is BlaxelFile {
  return (
    !!data &&
    typeof data === "object" &&
    "size" in data &&
    !("files" in data) &&
    !("subdirectories" in data)
  );
}

/**
 * Runtime state lives on Blaxel's private wire model (`sandbox.state`). The
 * published `SandboxInstance` typings expose status/spec getters but not state,
 * so this is the single boundary that reads it for Kit's handle projection.
 */
type BlaxelWireRuntime = {
  readonly sandbox?: { readonly state?: unknown };
};

function readWireRuntimeState(instance: SandboxInstance): string | undefined {
  const value = (instance as unknown as BlaxelWireRuntime).sandbox?.state;
  if (value === undefined || value === null || value === "") return undefined;
  const text = String(value);
  return text === "" ? undefined : text;
}

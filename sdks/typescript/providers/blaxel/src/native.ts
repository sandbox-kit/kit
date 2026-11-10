import {
  SandboxInstance,
  settings,
  type Config as NativeClientConfig,
  type Directory,
  type SandboxCreateConfiguration,
  type SandboxCreateOptions,
} from "@blaxel/core";
import type { BlaxelSandboxIdentity } from "#mappings/provider_mappings.js";

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
    rm(path: string, recursive?: boolean): Promise<unknown>;
    download(src: string, destinationPath: string): Promise<void>;
  };
  readonly process: {
    exec(request: ProcessExecRequest): Promise<BlaxelProcessResult>;
    wait(
      identifier: string,
      options?: {
        maxWait?: number;
        interval?: number;
        signal?: AbortSignal;
      },
    ): Promise<BlaxelProcessResult>;
  };
  delete(): Promise<unknown>;
}

export interface NativeBlaxel {
  configure(config: NativeClientConfig): void;
  create(
    params: SandboxCreateConfiguration,
    options?: SandboxCreateOptions,
  ): Promise<NativeSandboxHandle>;
  get(name: string): Promise<NativeSandboxHandle>;
  delete(name: string): Promise<unknown>;
}

export type NativeFactory = (config: NativeClientConfig) => NativeBlaxel;

export function defaultNativeFactory(config: NativeClientConfig): NativeBlaxel {
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
}

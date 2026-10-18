import { Daytona, type DaytonaConfig, type Sandbox as DaytonaSandbox } from "@daytona/sdk";
import { durationMs } from "@bufbuild/protobuf/wkt";
import {
  type Backend,
  type Config,
  type CreateOptions,
  type CreateResult,
  IsolationKind,
  LifecycleCapability,
  type OperationOptions,
  type ArchiveOptions,
  type DeleteOptions,
  type PauseOptions,
  type Provider,
  type SandboxInfo,
  type StartOptions,
  type StopOptions,
  metadataObjectAsMap,
} from "@sandbox-kit/sandbox";
import type { Duration } from "@sandbox-kit/sandbox";

import { mapClientConfig, mapCreateRequest, mapSandboxInfo } from "./mappings.js";
import { mapDaytonaError } from "./errors.js";
import { BackendBase, ProviderBase } from "./provider.js";

interface NativeDaytona {
  create(params: ReturnType<typeof mapCreateRequest>, options?: { readonly timeout?: number }): Promise<DaytonaSandbox>;
  get(id: string): Promise<DaytonaSandbox>;
  [Symbol.asyncDispose](): Promise<void>;
}

type NativeFactory = (config: DaytonaConfig) => NativeDaytona;

/** @internal Exported from this module for boundary testing, but not from the package root. */
export class DaytonaProvider extends ProviderBase {
  constructor(private readonly nativeFactory: NativeFactory = (config) => new Daytona(config)) {
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

  async create(request: CreateOptions, options: OperationOptions): Promise<CreateResult> {
    try {
      options.signal?.throwIfAborted();
      const params = mapCreateRequest(request);
      const timeout = options.timeoutMs === undefined ? undefined : options.timeoutMs / 1_000;
      const sandbox = await this.native.create(params, timeout === undefined ? undefined : { timeout });
      return { sandbox: mapSandboxInfo(sandbox, "create") };
    } catch (error) {
      throw mapDaytonaError(error, "create");
    }
  }

  capabilities(info: SandboxInfo): readonly LifecycleCapability[] {
    if (info.terminal) return [];
    const result: LifecycleCapability[] = [
      LifecycleCapability.Refresh,
      LifecycleCapability.Start,
      LifecycleCapability.Stop,
      LifecycleCapability.Delete,
    ];
    const metadata = metadataObjectAsMap(info.providerMetadata);
    const gpuCount = metadata?.gpu_count;
    const hasGpu =
      (typeof gpuCount === "number" && gpuCount > 0) ||
      (typeof gpuCount === "bigint" && gpuCount > 0n);
    if (
      !hasGpu &&
      (info.isolation?.kind === IsolationKind.LinuxVm || info.isolation?.providerKind === "windows")
    )
      result.push(LifecycleCapability.Pause);
    if (!hasGpu && info.isolation?.kind === IsolationKind.Container)
      result.push(LifecycleCapability.Archive);
    return result;
  }

  async refresh(id: string): Promise<SandboxInfo> {
    return this.withSandbox(id, "refresh", async (value) => {
      await value.refreshData();
    });
  }
  async start(id: string, request: StartOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "start", (value) => value.start(seconds(request.timeout)));
  }
  async stop(id: string, request: StopOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "stop", (value) => value.stop(seconds(request.timeout), request.force));
  }
  async pause(id: string, request: PauseOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "pause", (value) => value.pause(seconds(request.timeout)));
  }
  async archive(id: string, _request: ArchiveOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "archive", (value) => value.archive());
  }
  async delete(id: string, request: DeleteOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "delete", (value) => value.delete(seconds(request.timeout), request.wait));
  }
  async closeSandbox(_id: string): Promise<void> {}

  private async withSandbox(
    id: string,
    operation: string,
    action: (sandbox: DaytonaSandbox) => Promise<void>,
  ): Promise<SandboxInfo> {
    try {
      const sandbox = await this.native.get(id);
      await action(sandbox);
      return mapSandboxInfo(sandbox, operation);
    } catch (error) {
      throw mapDaytonaError(error, operation);
    }
  }

  async close(_options: OperationOptions): Promise<void> {
    try {
      await this.native[Symbol.asyncDispose]();
    } catch (error) {
      throw mapDaytonaError(error, "close");
    }
  }
}

function seconds(value: Duration | undefined): number | undefined {
  return value === undefined ? undefined : durationMs(value) / 1_000;
}

/** Selects Daytona as the provider for the shared Sandbox Kit client. */
export function daytona(): Provider {
  return new DaytonaProvider();
}

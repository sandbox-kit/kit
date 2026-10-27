import {
  Daytona,
  type DaytonaConfig,
  type Sandbox as DaytonaSandbox,
} from "@daytona/sdk";
import {
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
  PauseOptions,
  StartOptions,
  StopOptions,
} from "@sandbox-kit/sandbox";

import {
  mapClientConfig,
  mapCreateRequest,
  mapSandboxInfo,
} from "#mappings/mappings.js";
import { mapDaytonaError } from "#errors.js";
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
      options.signal?.throwIfAborted();
      const params = mapCreateRequest(request);
      const timeout = secondsFromMs(options.timeout);
      const sandbox = await this.native.create(
        params,
        timeout === undefined ? undefined : { timeout },
      );
      return { sandbox: mapSandboxInfo(sandbox, "create") };
    } catch (error) {
      throw mapDaytonaError(error, "create");
    }
  }

  async refresh(id: string): Promise<SandboxInfo> {
    return this.withSandbox(id, "refresh", async (value) => {
      await value.refreshData();
    });
  }
  async start(id: string, request: StartOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "start", (value) =>
      value.start(secondsFromMs(request.timeout)),
    );
  }
  async stop(id: string, request: StopOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "stop", (value) =>
      value.stop(secondsFromMs(request.timeout), request.force),
    );
  }
  async pause(id: string, request: PauseOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "pause", (value) =>
      value.pause(secondsFromMs(request.timeout)),
    );
  }
  async archive(id: string, _request: ArchiveOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "archive", (value) => value.archive());
  }
  async delete(id: string, request: DeleteOptions): Promise<SandboxInfo> {
    return this.withSandbox(id, "delete", (value) =>
      value.delete(secondsFromMs(request.timeout), request.wait),
    );
  }

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

  async closeSandbox(_id: string, _options: OperationOptions): Promise<void> {}

  async close(_options: OperationOptions): Promise<void> {
    try {
      await this.native[Symbol.asyncDispose]();
    } catch (error) {
      throw mapDaytonaError(error, "close");
    }
  }
}

/** Native Daytona lifecycle APIs take whole seconds; Kit timeouts are ms. */
function secondsFromMs(value: number | undefined): number | undefined {
  return value === undefined ? undefined : value / 1_000;
}

/** Selects Daytona as the provider for the shared Sandbox Kit client. */
export function daytona(): Provider {
  return new DaytonaProvider();
}

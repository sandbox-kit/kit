import {
  Daytona,
  type DaytonaConfig,
  type Sandbox as DaytonaSandbox,
} from "@daytona/sdk";
import {
  ConfigFieldTimeout,
  CreateFieldProvisioningTimeout,
  ErrorKind,
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
import { mapDaytonaError, providerError } from "#errors.js";
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

function throwIfAborted(signal: AbortSignal | undefined): void {
  signal?.throwIfAborted();
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

/** Selects Daytona as the provider for the shared Sandbox Kit client. */
export function daytona(): Provider {
  return new DaytonaProvider();
}

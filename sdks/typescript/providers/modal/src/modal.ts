import {
  ModalClient,
  type ModalClientParams,
  type Sandbox as ModalSandbox,
} from "modal";
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
  ExitResult,
  TerminateOptions,
} from "@sandbox-kit/sandbox";
import {
  mapClientConfig,
  mapCreateRequest,
  mapSandboxInfo,
} from "#mappings/mappings.js";
import { mapModalError } from "#errors.js";
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

/** Selects Modal as the provider for the shared Sandbox Kit client. */
export function modal(): Provider {
  return new ModalProvider();
}

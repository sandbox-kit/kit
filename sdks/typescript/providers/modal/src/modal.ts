import {
  ModalClient,
  type ModalClientParams,
  type Sandbox as ModalSandbox,
} from "modal";
import {
  LifecycleCapability,
  type Backend,
  type Config,
  type CreateOptions,
  type CreateResult,
  type ExitResult,
  type OperationOptions,
  type Provider,
  type SandboxInfo,
  type TerminateOptions,
} from "@sandbox-kit/sandbox";
import {
  mapClientConfig,
  mapCreateRequest,
  mapSandboxInfo,
} from "./mappings.js";
import { mapModalError } from "./errors.js";
import { BackendBase, ProviderBase } from "./provider.js";
import { captureProviderState, type ProviderState } from "./provider_config.js";

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
      options.signal?.throwIfAborted();
      const plan = mapCreateRequest(request, this.config);
      const app = await this.native.apps.fromName(plan.app, {
        ...(plan.environment === undefined
          ? {}
          : { environment: plan.environment }),
      });
      const secrets = await Promise.all(
        plan.secretNames.map((name) =>
          this.native.secrets.fromName(name, {
            ...(plan.environment === undefined
              ? {}
              : { environment: plan.environment }),
          }),
        ),
      );
      const proxy =
        plan.proxyName === undefined
          ? undefined
          : await this.native.proxies.fromName(plan.proxyName, {
              ...(plan.environment === undefined
                ? {}
                : { environment: plan.environment }),
            });
      const image = this.native.images.fromRegistry(plan.image);
      const sandbox = await this.native.sandboxes.create(app, image, {
        ...plan.params,
        ...(secrets.length === 0 ? {} : { secrets }),
        ...(proxy === undefined ? {} : { proxy }),
      });
      if (plan.waitUntilReady) {
        const readyTimeoutMs = plan.readyTimeoutMs ?? options.timeout;
        await sandbox.waitUntilReady(readyTimeoutMs);
      }
      return { sandbox: mapSandboxInfo(sandbox, request) };
    } catch (error) {
      throw mapModalError(error, "create");
    }
  }

  async attach(id: string): Promise<SandboxInfo> {
    try {
      return mapSandboxInfo(await this.native.sandboxes.fromId(id));
    } catch (error) {
      throw mapModalError(error, "attach");
    }
  }

  capabilities(): readonly LifecycleCapability[] {
    return [LifecycleCapability.Terminate];
  }

  async terminate(id: string, request: TerminateOptions): Promise<ExitResult> {
    try {
      const sandbox = await this.native.sandboxes.fromId(id);
      if (request.wait === true)
        return { exitCode: await sandbox.terminate({ wait: true }) };
      await sandbox.terminate();
      return {};
    } catch (error) {
      throw mapModalError(error, "terminate");
    }
  }

  async closeSandbox(id: string): Promise<void> {
    try {
      (await this.native.sandboxes.fromId(id)).detach();
    } catch (error) {
      throw mapModalError(error, "close");
    }
  }

  async close(): Promise<void> {
    this.native.close();
  }
}

/** Selects Modal as the provider for the shared Sandbox Kit client. */
export function modal(): Provider {
  return new ModalProvider();
}

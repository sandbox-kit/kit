import { ModalClient, type ModalClientParams, type Sandbox as ModalSandbox } from "modal";
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
import { mapClientConfig, mapCreateRequest, mapModalError, mapSandboxInfo } from "./mappings.js";
import { providerName } from "./provider.gen.js";

interface NativeModal {
  readonly apps: ModalClient["apps"];
  readonly images: ModalClient["images"];
  readonly sandboxes: ModalClient["sandboxes"];
  readonly secrets: ModalClient["secrets"];
  close(): void;
}
type NativeFactory = (config: ModalClientParams) => NativeModal;

/** @internal Exported for adapter boundary tests, but not from the package root. */
export class ModalProvider implements Provider {
  readonly name = providerName;
  constructor(private readonly nativeFactory: NativeFactory = (config) => new ModalClient(config)) {}
  initialize(config: Config): Backend {
    try { return new ModalBackend(this.nativeFactory(mapClientConfig(config)), structuredClone(config)); }
    catch (error) { throw mapModalError(error, "initialize"); }
  }
}

class ModalBackend implements Backend {
  readonly name = providerName;
  constructor(private readonly native: NativeModal, private readonly config: Config) {}

  async create(request: CreateOptions, _options: OperationOptions): Promise<CreateResult> {
    try {
      const plan = mapCreateRequest(request, this.config);
      const app = await this.native.apps.fromName(plan.app, { ...(plan.environment === undefined ? {} : { environment: plan.environment }) });
      const secrets = await Promise.all(plan.secretNames.map((name) => this.native.secrets.fromName(name, { ...(plan.environment === undefined ? {} : { environment: plan.environment }) })));
      const image = this.native.images.fromRegistry(plan.image);
      const sandbox = await this.native.sandboxes.create(app, image, { ...plan.params, ...(secrets.length === 0 ? {} : { secrets }) });
      if (plan.waitUntilReady) await sandbox.waitUntilReady(plan.readyTimeoutMs);
      return { sandbox: mapSandboxInfo(sandbox, request) };
    } catch (error) { throw mapModalError(error, "create"); }
  }

  async attach(id: string): Promise<SandboxInfo> {
    try { return mapSandboxInfo(await this.native.sandboxes.fromId(id)); }
    catch (error) { throw mapModalError(error, "attach"); }
  }

  capabilities(): readonly LifecycleCapability[] { return [LifecycleCapability.Terminate]; }

  async terminate(id: string, request: TerminateOptions): Promise<ExitResult> {
    try {
      const sandbox = await this.native.sandboxes.fromId(id);
      if (request.wait === true) return { exitCode: await sandbox.terminate({ wait: true }) };
      await sandbox.terminate();
      return {};
    } catch (error) { throw mapModalError(error, "terminate"); }
  }

  async closeSandbox(id: string): Promise<void> {
    try { (await this.native.sandboxes.fromId(id)).detach(); }
    catch (error) { throw mapModalError(error, "close"); }
  }

  async close(): Promise<void> { this.native.close(); }
}

/** Selects Modal as the provider for the shared Sandbox Kit client. */
export function modal(): Provider { return new ModalProvider(); }

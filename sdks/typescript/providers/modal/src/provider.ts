// Modal provider construction and native SDK orchestration are handwritten.

import type { Backend, Config, CreateOptions, CreateResult, OperationOptions, Provider } from "@sandbox-kit/sandbox";

export const providerName = "modal" as const;

/** @internal Shared generated provider identity wiring. */
export abstract class ProviderBase implements Provider {
  readonly name = providerName;
  abstract initialize(config: Config): Backend;
}

/** @internal Shared generated backend identity and required operation surface. */
export abstract class BackendBase implements Backend {
  readonly name = providerName;
  abstract create(request: CreateOptions, options: OperationOptions): Promise<CreateResult>;
  abstract close(options: OperationOptions): Promise<void>;
}

// Daytona provider construction and native SDK orchestration are handwritten.

import type {
  Backend,
  Config,
  CreateOptions,
  CreateResult,
  OperationOptions,
  Provider,
} from "@sandbox-kit/sandbox";

export const providerName = "daytona" as const;

/** @internal Shared provider identity wiring. */
export abstract class ProviderBase implements Provider {
  readonly name = providerName;
  abstract initialize(config: Config): Backend;
}

/** @internal Shared backend identity and required create/close surface. */
export abstract class BackendBase implements Backend {
  readonly name = providerName;
  abstract create(
    request: CreateOptions,
    options: OperationOptions,
  ): Promise<CreateResult>;
  abstract close(options: OperationOptions): Promise<void>;
  abstract closeSandbox(id: string, options: OperationOptions): Promise<void>;
}

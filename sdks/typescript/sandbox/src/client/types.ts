// Public client / provider / backend contracts are handwritten facades over
// portable protobuf config and create/lifecycle option types.

import type { MessageInitShape } from "@bufbuild/protobuf";
import {
  ConfigSchema,
  type APIKeyCredentials as WireAPIKeyCredentials,
  type AuthConfig,
  type BearerTokenCredentials as WireBearerTokenCredentials,
  type OAuthCredentials as WireOAuthCredentials,
  type TokenPairCredentials as WireTokenPairCredentials,
} from "@sandbox-kit/contracts/client";
import type {
  ArchiveOptions,
  CreateOptions,
  CreateResult,
  DeleteOptions,
  ExitResult,
  PauseOptions,
  SandboxInfo,
  StartOptions,
  StopOptions,
  TerminateOptions,
} from "#contracts/sandbox_data_contracts.js";
import type { DurationFieldsToMs } from "#internal/duration_ms.js";
import type { RemapFields } from "#internal/input_remap.js";
import type { FlattenOneof } from "#internal/oneof_input.js";
import type { StripProto } from "#internal/strip_proto.js";

/** Public credential shapes (no protobuf `$typeName` / getters). */
export type APIKeyCredentials = StripProto<WireAPIKeyCredentials>;
export type TokenPairCredentials = StripProto<WireTokenPairCredentials>;
export type BearerTokenCredentials = StripProto<WireBearerTokenCredentials>;
export type OAuthCredentials = StripProto<WireOAuthCredentials>;
export type Scope = StripProto<import("@sandbox-kit/contracts/client").Scope>;

/** Public auth: set exactly one credential variant (no protobuf `{ case, value }`). */
export type AuthInput = FlattenOneof<NonNullable<AuthConfig["credentials"]>>;

/**
 * Public client settings (ms timeout, flat auth).
 * Prefer {@link ClientConfig} at construction (`Config` plus `provider`).
 */
export type Config = DurationFieldsToMs<
  RemapFields<
    StripProto<MessageInitShape<typeof ConfigSchema>>,
    { auth: AuthInput }
  >
>;

export interface OperationOptions {
  readonly signal?: AbortSignal;
  /** Operation deadline in milliseconds. */
  readonly timeout?: number;
}

/** Public stop options: Kit call options plus optional force. */
export type StopCallOptions = OperationOptions & { readonly force?: boolean };
/** Public terminate options: Kit call options plus optional wait. */
export type TerminateCallOptions = OperationOptions & {
  readonly wait?: boolean;
};
/** Public delete options: Kit call options plus optional wait. */
export type DeleteCallOptions = OperationOptions & { readonly wait?: boolean };

export interface Backend {
  readonly name: string;
  create?(
    request: CreateOptions,
    options: OperationOptions,
  ): Promise<CreateResult>;
  attach?(id: string, options: OperationOptions): Promise<SandboxInfo>;
  refresh?(id: string, options: OperationOptions): Promise<SandboxInfo>;
  start?(
    id: string,
    request: StartOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  stop?(
    id: string,
    request: StopOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  pause?(
    id: string,
    request: PauseOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  terminate?(
    id: string,
    request: TerminateOptions,
    options: OperationOptions,
  ): Promise<ExitResult>;
  archive?(
    id: string,
    request: ArchiveOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  delete?(
    id: string,
    request: DeleteOptions,
    options: OperationOptions,
  ): Promise<SandboxInfo>;
  close(options: OperationOptions): Promise<void>;
  closeSandbox?(id: string, options: OperationOptions): Promise<void>;
}

export interface Provider {
  readonly name: string;
  /** Receives public {@link Config}. */
  initialize(config: Config): Backend;
}

/** Provider selection stays handwritten and is never part of the portable Config message. */
export type ClientConfig = Config & {
  readonly provider: Provider;
};

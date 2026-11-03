// Public client / provider / backend contracts are handwritten facades over
// portable protobuf config, create/lifecycle, and filesystem option types.

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
  ExecOptions,
  ExecResult,
  FileCopyOptions,
  FileInfo,
  FileReadOptions,
  FileWriteOptions,
  ListDirectoryOptions,
  MakeDirectoryOptions,
  PauseOptions,
  RemoveOptions,
  StatOptions,
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

/** Set exactly one of `apiKey`, `tokenPair`, `bearerToken`, or `oauthRefresh`. */
export type AuthInput = FlattenOneof<NonNullable<AuthConfig["credentials"]>>;

/**
 * Shared client settings. Timeouts are milliseconds.
 * Pass these fields to {@link ClientConfig} together with `provider`.
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
  readFile?(
    id: string,
    path: string,
    request: FileReadOptions,
    options: OperationOptions,
  ): Promise<Uint8Array>;
  writeFile?(
    id: string,
    path: string,
    data: Uint8Array,
    request: FileWriteOptions,
    options: OperationOptions,
  ): Promise<void>;
  copyFromLocal?(
    id: string,
    localPath: string,
    remotePath: string,
    request: FileCopyOptions,
    options: OperationOptions,
  ): Promise<void>;
  copyToLocal?(
    id: string,
    remotePath: string,
    localPath: string,
    request: FileCopyOptions,
    options: OperationOptions,
  ): Promise<void>;
  makeDirectory?(
    id: string,
    path: string,
    request: MakeDirectoryOptions,
    options: OperationOptions,
  ): Promise<void>;
  listDirectory?(
    id: string,
    path: string,
    request: ListDirectoryOptions,
    options: OperationOptions,
  ): Promise<FileInfo[]>;
  stat?(
    id: string,
    path: string,
    request: StatOptions,
    options: OperationOptions,
  ): Promise<FileInfo>;
  remove?(
    id: string,
    path: string,
    request: RemoveOptions,
    options: OperationOptions,
  ): Promise<void>;
  run?(
    id: string,
    request: ExecOptions,
    options: OperationOptions,
  ): Promise<ExecResult>;
  close(options: OperationOptions): Promise<void>;
  closeSandbox(id: string, options: OperationOptions): Promise<void>;
}

export interface Provider {
  readonly name: string;
  /** Receives public {@link Config}. */
  initialize(config: Config): Backend;
}

/** Settings for `new Client`, including the provider (`daytona()` or `modal()`). */
export type ClientConfig = Config & {
  readonly provider: Provider;
};

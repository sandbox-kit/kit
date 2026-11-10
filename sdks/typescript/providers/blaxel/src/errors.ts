// Native Blaxel error translation is handwritten.
// Classification order: existing SandboxError, abort/timeout DOMException,
// then Blaxel helpers / ResponseError status, then unknown.
// Cause stays on SandboxError.cause.

import { mkdir, readFile, stat, writeFile } from "node:fs/promises";
import { dirname } from "node:path";

import {
  CredentialsError,
  getBlaxelErrorCode,
  getBlaxelErrorStatus,
  isCreationTimeoutError,
  isGatewayError,
  isGatewayTimeout,
  ResponseError,
  SandboxCreationTimeoutError,
  SandboxGatewayError,
} from "@blaxel/core";
import {
  ErrorKind,
  FileFieldLocalPath,
  SandboxError,
  sandboxError,
  type ErrorInfo,
} from "@sandbox-kit/sandbox";

/** Inputs for translating a native/provider failure into {@link SandboxError}. */
interface MappedErrorInput {
  kind: ErrorInfo["kind"];
  error: unknown;
  operation: string;
  statusCode?: number;
  providerCode?: string;
}

/** Inputs for throwing a Blaxel-scoped validation or mapping error. */
export interface ProviderErrorInput {
  kind: ErrorInfo["kind"];
  operation: string;
  field: string;
  message: string;
}

export function mapBlaxelError(
  error: unknown,
  operation: string,
): SandboxError {
  if (error instanceof SandboxError) return error;
  if (error instanceof DOMException && error.name === "AbortError")
    return create({ kind: ErrorKind.Canceled, error, operation });
  if (error instanceof DOMException && error.name === "TimeoutError")
    return create({ kind: ErrorKind.Timeout, error, operation });

  const cause = causeOf(error);
  let kind: ErrorInfo["kind"] = ErrorKind.Unknown;
  const statusCode =
    getBlaxelErrorStatus(error) ?? getBlaxelErrorStatus(cause);
  const providerCode = stringifyCode(getBlaxelErrorCode(error));

  if (error instanceof CredentialsError) kind = ErrorKind.Authentication;
  if (
    kind === ErrorKind.Unknown &&
    (error instanceof SandboxCreationTimeoutError ||
      isCreationTimeoutError(error))
  )
    kind = ErrorKind.Timeout;
  // process.wait reports its deadline as a plain Error. A 504 on one poll is
  // retried inside wait and may be that error's cause. A bare edge 504
  // (completion hold or origin read) is a timeout too, matching Go.
  if (
    kind === ErrorKind.Unknown &&
    (isProcessWaitTimeout(error) ||
      isGatewayTimeout(error) ||
      isGatewayTimeout(cause))
  )
    kind = ErrorKind.Timeout;
  if (
    kind === ErrorKind.Unknown &&
    (error instanceof SandboxGatewayError || isGatewayError(error))
  )
    kind = ErrorKind.Unavailable;

  if (kind === ErrorKind.Unknown && providerCode !== undefined)
    kind = kindFromCode(providerCode);

  if (kind === ErrorKind.Unknown && statusCode !== undefined)
    kind = kindFromStatus(statusCode);

  if (
    kind === ErrorKind.Unknown &&
    error instanceof ResponseError &&
    error.status !== undefined
  )
    kind = kindFromStatus(error.status);

  return create({
    kind,
    error,
    operation,
    ...(statusCode === undefined ? {} : { statusCode }),
    ...(providerCode === undefined ? {} : { providerCode }),
  });
}

function kindFromCode(code: string): ErrorInfo["kind"] {
  switch (code) {
    case "AUTHENTICATION_REQUIRED":
    case "AUTHENTICATION_FAILED":
    case "TOKEN_REVOKED":
      return ErrorKind.Authentication;
    case "FORBIDDEN":
    case "UNAUTHORIZED":
    case "UNAUTHORIZED_WORKSPACE":
      return ErrorKind.PermissionDenied;
    case "ROUTE_NOT_FOUND":
    case "WORKLOAD_NOT_FOUND":
    case "WORKSPACE_NOT_FOUND":
    case "IMAGE_NOT_FOUND":
    case "VOLUME_NOT_FOUND":
    case "ACCOUNT_NOT_FOUND":
    case "SANDBOX_DELETED":
    case "VOLUME_DELETED":
      return ErrorKind.NotFound;
    case "SANDBOX_ALREADY_EXISTS":
    case "VOLUME_ALREADY_ATTACHED":
      return ErrorKind.AlreadyExists;
    case "BAD_REQUEST":
    case "POLICY_VIOLATION": // reserved gateway code; matches Go ErrPolicyViolation
    case "INVALID_INPUT":
    case "INVALID_IMAGE":
    case "INVALID_PORTS":
    case "INVALID_REGION":
    case "INVALID_VOLUME":
    case "INVALID_VOLUME_COUNT":
    case "INVALID_VOLUME_REGION":
    case "VALIDATION_FAILED":
    case "VALIDATION_ERROR":
    case "WORKSPACE_REQUIRED":
      return ErrorKind.InvalidArgument;
    case "CREATION_TIMEOUT":
    case "UPSTREAM_CONNECT_TIMEOUT":
    case "UPSTREAM_TIMEOUT":
      return ErrorKind.Timeout;
    case "QUOTA_EXCEEDED":
    case "USAGE_LIMIT_EXCEEDED":
      return ErrorKind.ResourceExhausted;
    case "RATE_LIMIT_EXCEEDED":
    case "WORKSPACE_RATE_LIMIT_EXCEEDED":
      return ErrorKind.RateLimited;
    case "SANDBOX_DELETION_IN_PROGRESS":
      return ErrorKind.Conflict;
    case "WORKLOAD_UNAVAILABLE":
    case "ROUTING_UNAVAILABLE":
    case "SERVICE_UNAVAILABLE":
    case "AUTHENTICATION_UNAVAILABLE":
    case "UPSTREAM_CONNECT_FAILED":
    case "UPSTREAM_ERROR":
    case "GATEWAY_INTERNAL_ERROR":
      return ErrorKind.Unavailable;
    case "WORKLOAD_FAILED":
    case "INTERNAL_ERROR":
    case "HANDLER_ERROR":
    case "DATABASE_ERROR":
      return ErrorKind.Internal;
    default:
      return ErrorKind.Unknown;
  }
}

function kindFromStatus(status: number): ErrorInfo["kind"] {
  switch (status) {
    case 400:
    case 422:
      return ErrorKind.InvalidArgument;
    case 401:
      return ErrorKind.Authentication;
    case 402:
      return ErrorKind.ResourceExhausted;
    case 403:
      return ErrorKind.PermissionDenied;
    case 404:
    case 410:
      return ErrorKind.NotFound;
    case 408:
      return ErrorKind.Timeout;
    case 409:
      return ErrorKind.Conflict;
    case 412:
      return ErrorKind.FailedPrecondition;
    case 429:
      return ErrorKind.RateLimited;
    case 500:
      return ErrorKind.Internal;
    case 502:
    case 503:
      return ErrorKind.Unavailable;
    case 504:
      return ErrorKind.Timeout;
    default:
      return ErrorKind.Unknown;
  }
}

function isProcessWaitTimeout(error: unknown): boolean {
  return (
    error instanceof Error &&
    error.message.startsWith("Process did not finish in time")
  );
}

function causeOf(error: unknown): unknown {
  return error instanceof Error ? error.cause : undefined;
}

function stringifyCode(code: unknown): string | undefined {
  if (typeof code === "string" && code !== "") return code;
  if (typeof code === "number") return String(code);
  return undefined;
}

function create(input: MappedErrorInput): SandboxError {
  const { kind, error, operation, statusCode, providerCode } = input;
  return sandboxError(
    {
      kind,
      provider: "blaxel",
      operation,
      field: "",
      message:
        error instanceof Error
          ? error.message
          : "sandbox-kit blaxel: provider operation failed",
      ...(statusCode === undefined ? {} : { statusCode }),
      ...(providerCode === undefined ? {} : { providerCode }),
    },
    error,
  );
}

/** Rejects a missing local file or a directory before Blaxel reads it. */
export async function requireLocalFile(
  localPath: string,
  operation: string,
): Promise<void> {
  let info;
  try {
    info = await stat(localPath);
  } catch (error) {
    if (nodeCode(error)) throw localFileError(error, operation);
    throw error;
  }
  if (info.isDirectory()) localDirectory(operation);
}

/** Rejects a local destination that is already a directory. */
export async function requireLocalDestination(
  localPath: string,
  operation: string,
): Promise<void> {
  try {
    const info = await stat(localPath);
    if (info.isDirectory()) localDirectory(operation);
  } catch (error) {
    if (nodeCode(error) === "ENOENT") return;
    if (nodeCode(error)) throw localFileError(error, operation);
    throw error;
  }
}

/** Creates the parent directory for a local destination path. */
export async function ensureLocalParent(
  localPath: string,
  operation: string,
): Promise<void> {
  try {
    await mkdir(dirname(localPath), { recursive: true });
  } catch (error) {
    throw localFileError(error, operation);
  }
}

/** Reads a local file for upload into the sandbox. */
export async function readLocalFile(
  localPath: string,
  operation: string,
): Promise<Uint8Array> {
  await requireLocalFile(localPath, operation);
  try {
    return await readFile(localPath);
  } catch (error) {
    throw localFileError(error, operation);
  }
}

/** Writes downloaded sandbox bytes to a local path. */
export async function writeLocalFile(
  localPath: string,
  data: Uint8Array,
  operation: string,
): Promise<void> {
  await requireLocalDestination(localPath, operation);
  await ensureLocalParent(localPath, operation);
  try {
    await writeFile(localPath, data);
  } catch (error) {
    throw localFileError(error, operation);
  }
}

/** Classifies a local filesystem failure. Other errors stay with provider mapping. */
export function translateLocalCopyError(
  error: unknown,
  operation: string,
): unknown {
  if (error instanceof SandboxError) return error;
  const code = nodeCode(error);
  if (!code || !localCodes.has(code)) return error;
  return localFileError(error, operation);
}

function localDirectory(operation: string): never {
  providerError({
    kind: ErrorKind.InvalidArgument,
    operation,
    field: FileFieldLocalPath,
    message: "local path is a directory",
  });
}

function localFileError(error: unknown, operation: string): SandboxError {
  return sandboxError(
    {
      kind:
        nodeCode(error) === "ENOENT"
          ? ErrorKind.NotFound
          : ErrorKind.InvalidArgument,
      provider: "blaxel",
      operation,
      field: FileFieldLocalPath,
      message: error instanceof Error ? error.message : "local path failed",
    },
    error,
  );
}

function nodeCode(error: unknown): string | undefined {
  if (typeof error !== "object" || !error || !("code" in error))
    return undefined;
  return typeof error.code === "string" ? error.code : undefined;
}

const localCodes = new Set([
  "ENOENT",
  "EISDIR",
  "ENOTDIR",
  "EACCES",
  "EPERM",
  "ENOSPC",
]);

/** Local throw helper for mappings, checks, and config. */
export function providerError(input: ProviderErrorInput): never {
  throw sandboxError({
    kind: input.kind,
    provider: "blaxel",
    operation: input.operation,
    field: input.field,
    message: "sandbox-kit blaxel: " + input.message,
  });
}

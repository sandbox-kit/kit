// Native Modal error translation is handwritten.
// Classification order: existing SandboxError, abort/timeout DOMException,
// then native Modal types, then unknown. Cause stays on SandboxError.cause.

import { stat } from "node:fs/promises";

import {
  AlreadyExistsError,
  ClientClosedError,
  ConflictError,
  InvalidError,
  NotFoundError,
  ResourceExhaustedError,
  SandboxFilesystemDirectoryNotEmptyError,
  SandboxFilesystemFileTooLargeError,
  SandboxFilesystemIsADirectoryError,
  SandboxFilesystemNotADirectoryError,
  SandboxFilesystemNotFoundError,
  SandboxFilesystemPathAlreadyExistsError,
  SandboxFilesystemPermissionError,
  SandboxTimeoutError,
  TimeoutError,
} from "modal";
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
  providerSource?: string;
}

/** Inputs for throwing a Modal-scoped validation or mapping error. */
export interface ProviderErrorInput {
  kind: ErrorInfo["kind"];
  operation: string;
  field: string;
  message: string;
}

export function mapModalError(error: unknown, operation: string): SandboxError {
  if (error instanceof SandboxError) return error;
  if (error instanceof DOMException && error.name === "AbortError")
    return create({ kind: ErrorKind.Canceled, error, operation });
  if (error instanceof DOMException && error.name === "TimeoutError")
    return create({ kind: ErrorKind.Timeout, error, operation });
  let kind: ErrorInfo["kind"] = ErrorKind.Unknown;
  if (kind === ErrorKind.Unknown && error instanceof InvalidError)
    kind = ErrorKind.InvalidArgument;
  if (kind === ErrorKind.Unknown && error instanceof NotFoundError)
    kind = ErrorKind.NotFound;
  if (kind === ErrorKind.Unknown && error instanceof AlreadyExistsError)
    kind = ErrorKind.AlreadyExists;
  if (kind === ErrorKind.Unknown && error instanceof ConflictError)
    kind = ErrorKind.Conflict;
  if (kind === ErrorKind.Unknown && error instanceof ResourceExhaustedError)
    kind = ErrorKind.ResourceExhausted;
  if (kind === ErrorKind.Unknown && error instanceof TimeoutError)
    kind = ErrorKind.Timeout;
  if (kind === ErrorKind.Unknown && error instanceof SandboxTimeoutError)
    kind = ErrorKind.Timeout;
  if (kind === ErrorKind.Unknown && error instanceof ClientClosedError)
    kind = ErrorKind.FailedPrecondition;
  if (
    kind === ErrorKind.Unknown &&
    error instanceof SandboxFilesystemNotFoundError
  )
    kind = ErrorKind.NotFound;
  if (
    kind === ErrorKind.Unknown &&
    error instanceof SandboxFilesystemPermissionError
  )
    kind = ErrorKind.PermissionDenied;
  if (
    kind === ErrorKind.Unknown &&
    error instanceof SandboxFilesystemPathAlreadyExistsError
  )
    kind = ErrorKind.AlreadyExists;
  if (
    kind === ErrorKind.Unknown &&
    error instanceof SandboxFilesystemDirectoryNotEmptyError
  )
    kind = ErrorKind.FailedPrecondition;
  if (
    kind === ErrorKind.Unknown &&
    error instanceof SandboxFilesystemFileTooLargeError
  )
    kind = ErrorKind.ResourceExhausted;
  if (
    kind === ErrorKind.Unknown &&
    (error instanceof SandboxFilesystemNotADirectoryError ||
      error instanceof SandboxFilesystemIsADirectoryError)
  )
    kind = ErrorKind.InvalidArgument;
  return create({ kind, error, operation });
}

function create(input: MappedErrorInput): SandboxError {
  const { kind, error, operation, statusCode, providerCode, providerSource } =
    input;
  return sandboxError(
    {
      kind,
      provider: "modal",
      operation,
      field: "",
      message:
        error instanceof Error
          ? error.message
          : "sandbox-kit modal: provider operation failed",
      ...(statusCode === undefined ? {} : { statusCode }),
      ...(providerCode === undefined ? {} : { providerCode }),
      ...(providerSource === undefined ? {} : { providerSource }),
    },
    error,
  );
}

/** Rejects a missing local file or a directory before Modal reads it. */
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

/** Rejects a local destination that is already a directory. A missing path is left for Modal to create. */
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
      provider: "modal",
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
    provider: "modal",
    operation: input.operation,
    field: input.field,
    message: "sandbox-kit modal: " + input.message,
  });
}

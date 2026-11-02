// Native Daytona error translation is handwritten.
// Classification order: existing SandboxError, abort/timeout DOMException,
// then DaytonaError code/status maps, then unknown. Cause stays on SandboxError.cause.

import { stat } from "node:fs/promises";
import { DaytonaError } from "@daytona/sdk";
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

/** Inputs for throwing a Daytona-scoped validation or mapping error. */
export interface ProviderErrorInput {
  kind: ErrorInfo["kind"];
  operation: string;
  field: string;
  message: string;
}

export function mapDaytonaError(
  error: unknown,
  operation: string,
): SandboxError {
  if (error instanceof SandboxError) return error;
  if (error instanceof DOMException && error.name === "AbortError")
    return create({ kind: ErrorKind.Canceled, error, operation });
  if (error instanceof DOMException && error.name === "TimeoutError")
    return create({ kind: ErrorKind.Timeout, error, operation });
  let kind: ErrorInfo["kind"] = ErrorKind.Unknown;
  let statusCode: number | undefined;
  let providerCode: string | undefined;
  let providerSource: string | undefined;
  if (error instanceof DaytonaError) {
    statusCode = error.statusCode;
    providerCode = error.code;
    providerSource = error.source;
    if (kind === ErrorKind.Unknown)
      switch (error.code) {
        case "SANDBOX_SPOT_EVICTED":
          kind = ErrorKind.Unavailable;
          break;
      }
    if (kind === ErrorKind.Unknown)
      switch (error.statusCode) {
        case 400:
          kind = ErrorKind.InvalidArgument;
          break;
        case 401:
          kind = ErrorKind.Authentication;
          break;
        case 403:
          kind = ErrorKind.PermissionDenied;
          break;
        case 404:
          kind = ErrorKind.NotFound;
          break;
        case 408:
          kind = ErrorKind.Timeout;
          break;
        case 409:
          kind = ErrorKind.Conflict;
          break;
        case 410:
          kind = ErrorKind.NotFound;
          break;
        case 412:
          kind = ErrorKind.FailedPrecondition;
          break;
        case 422:
          kind = ErrorKind.InvalidArgument;
          break;
        case 429:
          kind = ErrorKind.RateLimited;
          break;
        case 500:
          kind = ErrorKind.Internal;
          break;
        case 502:
          kind = ErrorKind.Unavailable;
          break;
        case 503:
          kind = ErrorKind.Unavailable;
          break;
        case 504:
          kind = ErrorKind.Timeout;
          break;
      }
    if (kind === ErrorKind.Unknown && error.statusCode === 0 && !error.code)
      kind = ErrorKind.Unavailable;
  }
  return create({
    kind,
    error,
    operation,
    ...(statusCode === undefined ? {} : { statusCode }),
    ...(providerCode === undefined ? {} : { providerCode }),
    ...(providerSource === undefined ? {} : { providerSource }),
  });
}

function create(input: MappedErrorInput): SandboxError {
  const { kind, error, operation, statusCode, providerCode, providerSource } =
    input;
  return sandboxError(
    {
      kind,
      provider: "daytona",
      operation,
      field: "",
      message:
        error instanceof Error
          ? error.message
          : "sandbox-kit daytona: provider operation failed",
      ...(statusCode === undefined ? {} : { statusCode }),
      ...(providerCode === undefined ? {} : { providerCode }),
      ...(providerSource === undefined ? {} : { providerSource }),
    },
    error,
  );
}

/** Rejects a missing local file or a directory before Daytona reads it. */
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

/** Rejects a local destination that is already a directory. A missing path is left for the write. */
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

/**
 * Classifies a local write failure from copyToLocal.
 * Other errors are returned unchanged for provider mapping.
 */
export function translateLocalCopyError(
  error: unknown,
  operation: string,
): unknown {
  if (error instanceof SandboxError) return error;
  if (nodeCode(error)) return localFileError(error, operation);
  if (
    error instanceof DaytonaError &&
    error.message.includes("Write stream failed")
  ) {
    return sandboxError(
      {
        kind: error.message.includes("ENOENT")
          ? ErrorKind.NotFound
          : ErrorKind.InvalidArgument,
        provider: "daytona",
        operation,
        field: FileFieldLocalPath,
        message: error.message,
      },
      error,
    );
  }
  return error;
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
      provider: "daytona",
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

/** Local throw helper for mappings, checks, and config. */
export function providerError(input: ProviderErrorInput): never {
  throw sandboxError({
    kind: input.kind,
    provider: "daytona",
    operation: input.operation,
    field: input.field,
    message: "sandbox-kit daytona: " + input.message,
  });
}

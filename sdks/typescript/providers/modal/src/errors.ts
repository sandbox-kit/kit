// Native Modal error translation is handwritten.
// Classification order: existing SandboxError, abort/timeout DOMException,
// then native Modal types, then unknown. Cause stays on SandboxError.cause.

import {
  AlreadyExistsError,
  ClientClosedError,
  ConflictError,
  InvalidError,
  NotFoundError,
  ResourceExhaustedError,
  SandboxTimeoutError,
  TimeoutError,
} from "modal";
import {
  ErrorKind,
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

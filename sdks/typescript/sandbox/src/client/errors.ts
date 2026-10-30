// Portable error wrappers are handwritten. ErrorInfo / ErrorKind remain
// protobuf-owned contract data.

import { ErrorKind as ProtoErrorKind } from "@sandbox-kit/contracts/sandbox";
import type { ErrorInfo } from "#contracts/sandbox_data_contracts.js";
import { protoEnum } from "#internal/proto_enum.js";

// Ergonomic PascalCase names over protobuf enum values; not a second contract.
export const ErrorKind = protoEnum(ProtoErrorKind);
export type ErrorKind = (typeof ErrorKind)[keyof typeof ErrorKind];

/**
 * Failure from the client or a sandbox handle.
 * `kind` is the category, `field` is the setting to change, and `cause` is the
 * provider error when there is one.
 */
export class SandboxError extends Error {
  readonly kind: ErrorInfo["kind"];
  readonly provider: string;
  readonly operation: string;
  readonly field: string;
  readonly rule: string | undefined;
  readonly statusCode: number | undefined;
  readonly providerCode: string | undefined;
  readonly providerSource: string | undefined;

  constructor(info: ErrorInfo, options?: ErrorOptions) {
    super(info.message, options);
    this.name = "SandboxError";
    this.kind = info.kind;
    this.provider = info.provider;
    this.operation = info.operation;
    this.field = info.field;
    this.rule = info.rule;
    this.statusCode = info.statusCode;
    this.providerCode = info.providerCode;
    this.providerSource = info.providerSource;
  }
}

export function sandboxError(info: ErrorInfo, cause?: unknown): SandboxError {
  return new SandboxError(info, cause === undefined ? undefined : { cause });
}

export function unsupported(provider: string, operation: string): never {
  throw sandboxError({
    kind: ErrorKind.Unsupported,
    provider,
    operation,
    field: "",
    message: "sandbox-kit: provider does not support " + operation,
  });
}

export function withContext(
  error: unknown,
  provider: string,
  operation: string,
): SandboxError {
  if (error instanceof SandboxError) {
    if (error.provider === provider && error.operation === operation)
      return error;
    return sandboxError(
      {
        kind: error.kind as ErrorInfo["kind"],
        provider: error.provider || provider,
        operation: error.operation || operation,
        field: error.field,
        message: error.message,
        ...(error.rule === undefined ? {} : { rule: error.rule }),
        ...(error.statusCode === undefined
          ? {}
          : { statusCode: error.statusCode }),
        ...(error.providerCode === undefined
          ? {}
          : { providerCode: error.providerCode }),
        ...(error.providerSource === undefined
          ? {}
          : { providerSource: error.providerSource }),
      },
      error.cause,
    );
  }
  if (error instanceof DOMException && error.name === "AbortError")
    return sandboxError(
      {
        kind: ErrorKind.Canceled,
        provider,
        operation,
        field: "",
        message: error.message,
      },
      error,
    );
  if (error instanceof DOMException && error.name === "TimeoutError")
    return sandboxError(
      {
        kind: ErrorKind.Timeout,
        provider,
        operation,
        field: "",
        message: error.message,
      },
      error,
    );
  return sandboxError(
    {
      kind: ErrorKind.Unknown,
      provider,
      operation,
      field: "",
      message:
        error instanceof Error
          ? error.message
          : "sandbox-kit: provider operation failed",
    },
    error,
  );
}

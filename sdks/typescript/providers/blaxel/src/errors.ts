// Native Blaxel error translation is handwritten.
// Classification order: existing SandboxError, abort/timeout DOMException,
// then Blaxel helpers / ResponseError status, then unknown.
// Cause stays on SandboxError.cause.

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

const CODE_KINDS: Record<string, ErrorInfo["kind"]> = {
  AUTHENTICATION_REQUIRED: ErrorKind.Authentication,
  AUTHENTICATION_FAILED: ErrorKind.Authentication,
  TOKEN_REVOKED: ErrorKind.Authentication,
  FORBIDDEN: ErrorKind.PermissionDenied,
  UNAUTHORIZED: ErrorKind.PermissionDenied,
  UNAUTHORIZED_WORKSPACE: ErrorKind.PermissionDenied,
  ROUTE_NOT_FOUND: ErrorKind.NotFound,
  WORKLOAD_NOT_FOUND: ErrorKind.NotFound,
  WORKSPACE_NOT_FOUND: ErrorKind.NotFound,
  IMAGE_NOT_FOUND: ErrorKind.NotFound,
  VOLUME_NOT_FOUND: ErrorKind.NotFound,
  ACCOUNT_NOT_FOUND: ErrorKind.NotFound,
  SANDBOX_DELETED: ErrorKind.NotFound,
  VOLUME_DELETED: ErrorKind.NotFound,
  SANDBOX_ALREADY_EXISTS: ErrorKind.AlreadyExists,
  VOLUME_ALREADY_ATTACHED: ErrorKind.AlreadyExists,
  BAD_REQUEST: ErrorKind.InvalidArgument,
  POLICY_VIOLATION: ErrorKind.InvalidArgument,
  INVALID_INPUT: ErrorKind.InvalidArgument,
  INVALID_IMAGE: ErrorKind.InvalidArgument,
  INVALID_PORTS: ErrorKind.InvalidArgument,
  INVALID_REGION: ErrorKind.InvalidArgument,
  INVALID_VOLUME: ErrorKind.InvalidArgument,
  INVALID_VOLUME_COUNT: ErrorKind.InvalidArgument,
  INVALID_VOLUME_REGION: ErrorKind.InvalidArgument,
  VALIDATION_FAILED: ErrorKind.InvalidArgument,
  VALIDATION_ERROR: ErrorKind.InvalidArgument,
  WORKSPACE_REQUIRED: ErrorKind.InvalidArgument,
  CREATION_TIMEOUT: ErrorKind.Timeout,
  UPSTREAM_CONNECT_TIMEOUT: ErrorKind.Timeout,
  UPSTREAM_TIMEOUT: ErrorKind.Timeout,
  QUOTA_EXCEEDED: ErrorKind.ResourceExhausted,
  USAGE_LIMIT_EXCEEDED: ErrorKind.ResourceExhausted,
  RATE_LIMIT_EXCEEDED: ErrorKind.RateLimited,
  WORKSPACE_RATE_LIMIT_EXCEEDED: ErrorKind.RateLimited,
  SANDBOX_DELETION_IN_PROGRESS: ErrorKind.Conflict,
  WORKLOAD_UNAVAILABLE: ErrorKind.Unavailable,
  ROUTING_UNAVAILABLE: ErrorKind.Unavailable,
  SERVICE_UNAVAILABLE: ErrorKind.Unavailable,
  AUTHENTICATION_UNAVAILABLE: ErrorKind.Unavailable,
  UPSTREAM_CONNECT_FAILED: ErrorKind.Unavailable,
  UPSTREAM_ERROR: ErrorKind.Unavailable,
  GATEWAY_INTERNAL_ERROR: ErrorKind.Unavailable,
  WORKLOAD_FAILED: ErrorKind.Internal,
  INTERNAL_ERROR: ErrorKind.Internal,
  HANDLER_ERROR: ErrorKind.Internal,
  DATABASE_ERROR: ErrorKind.Internal,
};

const STATUS_KINDS: Record<number, ErrorInfo["kind"]> = {
  400: ErrorKind.InvalidArgument,
  422: ErrorKind.InvalidArgument,
  401: ErrorKind.Authentication,
  402: ErrorKind.ResourceExhausted,
  403: ErrorKind.PermissionDenied,
  404: ErrorKind.NotFound,
  410: ErrorKind.NotFound,
  408: ErrorKind.Timeout,
  // Bare 409 is usually name collision; coded conflicts (e.g. deletion in
  // progress) still win through CODE_KINDS above.
  409: ErrorKind.AlreadyExists,
  412: ErrorKind.FailedPrecondition,
  429: ErrorKind.RateLimited,
  500: ErrorKind.Internal,
  502: ErrorKind.Unavailable,
  503: ErrorKind.Unavailable,
  504: ErrorKind.Timeout,
};

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
  const statusCode = getBlaxelErrorStatus(error) ?? getBlaxelErrorStatus(cause);
  const providerCode = stringifyCode(
    getBlaxelErrorCode(error) ?? getBlaxelErrorCode(cause),
  );
  return create({
    kind: classifyKind({ error, cause, providerCode, statusCode }),
    error,
    operation,
    ...(statusCode === undefined ? {} : { statusCode }),
    ...(providerCode === undefined ? {} : { providerCode }),
  });
}

interface ClassifyArgs {
  error: unknown;
  cause: unknown;
  providerCode: string | undefined;
  statusCode: number | undefined;
}

function classifyKind(args: ClassifyArgs): ErrorInfo["kind"] {
  const { error, cause, providerCode, statusCode } = args;
  if (error instanceof CredentialsError) return ErrorKind.Authentication;
  if (
    error instanceof SandboxCreationTimeoutError ||
    isCreationTimeoutError(error)
  )
    return ErrorKind.Timeout;
  // process.wait deadline is a plain Error; a 504 on one poll may be its cause.
  // A bare edge 504 is a timeout too, matching Go.
  if (
    isProcessWaitTimeout(error) ||
    isGatewayTimeout(error) ||
    isGatewayTimeout(cause)
  )
    return ErrorKind.Timeout;
  if (error instanceof SandboxGatewayError || isGatewayError(error))
    return ErrorKind.Unavailable;

  // Unknown codes fall through so status can still classify the failure.
  if (providerCode !== undefined && CODE_KINDS[providerCode] !== undefined)
    return CODE_KINDS[providerCode]!;
  if (statusCode !== undefined && STATUS_KINDS[statusCode] !== undefined)
    return STATUS_KINDS[statusCode]!;
  if (error instanceof ResponseError && error.status !== undefined)
    return STATUS_KINDS[error.status] ?? ErrorKind.Unknown;
  return ErrorKind.Unknown;
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

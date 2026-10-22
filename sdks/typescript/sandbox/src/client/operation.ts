// Operation deadline / AbortSignal helpers shared by Client and Sandbox.

import { ConfigFieldTimeout } from "#contracts/config_paths.js";
import { ErrorKind, sandboxError } from "#client/errors.js";
import type { OperationOptions } from "#client/types.js";

export function clone<T>(value: T): T {
  return structuredClone(value);
}

export function operationSignal(
  options: OperationOptions,
  inheritedTimeoutMs?: number,
): { signal?: AbortSignal; dispose(): void } {
  const timeout = options.timeout ?? inheritedTimeoutMs;
  if (timeout !== undefined && (!Number.isFinite(timeout) || timeout < 0)) {
    throw sandboxError({
      kind: ErrorKind.InvalidArgument,
      provider: "",
      operation: "",
      field: ConfigFieldTimeout,
      message: "sandbox-kit: timeout must be a finite non-negative number",
    });
  }
  if ((timeout === undefined || timeout === 0) && options.signal === undefined)
    return { dispose() {} };
  const controller = new AbortController();
  const abort = () => controller.abort(options.signal?.reason);
  if (options.signal?.aborted) abort();
  else options.signal?.addEventListener("abort", abort, { once: true });
  const timer =
    timeout !== undefined && timeout > 0
      ? setTimeout(
          () =>
            controller.abort(
              new DOMException("The operation timed out", "TimeoutError"),
            ),
          timeout,
        )
      : undefined;
  return {
    signal: controller.signal,
    dispose() {
      if (timer !== undefined) clearTimeout(timer);
      options.signal?.removeEventListener("abort", abort);
    },
  };
}

export function operationOptions(
  options: OperationOptions,
  signal: AbortSignal | undefined,
): OperationOptions {
  return {
    ...(options.timeout === undefined ? {} : { timeout: options.timeout }),
    ...(signal === undefined ? {} : { signal }),
  };
}

function abortReason(signal: AbortSignal): DOMException {
  const reason = signal.reason;
  return reason instanceof DOMException &&
    (reason.name === "AbortError" || reason.name === "TimeoutError")
    ? reason
    : new DOMException("The operation was aborted", "AbortError");
}

export async function awaitWithSignal<T>(
  operation: () => Promise<T>,
  signal?: AbortSignal,
): Promise<T> {
  if (signal?.aborted) throw abortReason(signal);
  const pending = operation();
  if (signal === undefined) return pending;
  return new Promise<T>((resolve, reject) => {
    const aborted = () => reject(abortReason(signal));
    const cleanup = () => signal.removeEventListener("abort", aborted);
    signal.addEventListener("abort", aborted, { once: true });
    pending.then(
      (value) => {
        cleanup();
        resolve(value);
      },
      (error: unknown) => {
        cleanup();
        reject(error);
      },
    );
  });
}

/** Lifecycle request timeouts stay milliseconds for Backend / providers. */
export function publicTimeoutOptions(options: OperationOptions): {
  timeout?: number;
} {
  return options.timeout === undefined ? {} : { timeout: options.timeout };
}

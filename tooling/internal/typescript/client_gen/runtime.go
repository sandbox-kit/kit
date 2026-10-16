package clientgen

// runtimeHelpers is handwritten orchestration that cannot be represented by
// the portable API declarations. Keeping it as TypeScript source makes the
// runtime reviewable; semantic surfaces and lifecycle methods stay generated
// from the compiled model.
const runtimeHelpers = `
function clone<T>(value: T): T {
  return structuredClone(value);
}

function operationSignal(options: OperationOptions, inheritedTimeoutMs?: number): { signal?: AbortSignal; dispose(): void } {
  const timeoutMs = options.timeoutMs ?? inheritedTimeoutMs;
  if (timeoutMs !== undefined && (!Number.isFinite(timeoutMs) || timeoutMs < 0)) {
    throw sandboxError({ kind: ErrorKind.InvalidArgument, provider: "", operation: "", field: "timeoutMs", message: "sandbox-kit: timeoutMs must be a finite non-negative number" });
  }
  if ((timeoutMs === undefined || timeoutMs === 0) && options.signal === undefined) return { dispose() {} };
  const controller = new AbortController();
  const abort = () => controller.abort(options.signal?.reason);
  if (options.signal?.aborted) abort();
  else options.signal?.addEventListener("abort", abort, { once: true });
  const timer = timeoutMs !== undefined && timeoutMs > 0
    ? setTimeout(() => controller.abort(new DOMException("The operation timed out", "TimeoutError")), timeoutMs)
    : undefined;
  return {
    signal: controller.signal,
    dispose() {
      if (timer !== undefined) clearTimeout(timer);
      options.signal?.removeEventListener("abort", abort);
    },
  };
}

function operationOptions(options: OperationOptions, signal: AbortSignal | undefined): OperationOptions {
  return {
    ...(options.timeoutMs === undefined ? {} : { timeoutMs: options.timeoutMs }),
    ...(signal === undefined ? {} : { signal }),
  };
}

function abortReason(signal: AbortSignal): DOMException {
  const reason = signal.reason;
  return reason instanceof DOMException && (reason.name === "AbortError" || reason.name === "TimeoutError")
    ? reason
    : new DOMException("The operation was aborted", "AbortError");
}

async function awaitWithSignal<T>(operation: () => Promise<T>, signal?: AbortSignal): Promise<T> {
  if (signal?.aborted) throw abortReason(signal);
  const pending = operation();
  if (signal === undefined) return pending;
  return new Promise<T>((resolve, reject) => {
    const aborted = () => reject(abortReason(signal));
    const cleanup = () => signal.removeEventListener("abort", aborted);
    signal.addEventListener("abort", aborted, { once: true });
    pending.then(
      (value) => { cleanup(); resolve(value); },
      (error: unknown) => { cleanup(); reject(error); },
    );
  });
}
`

const sandboxUpdateTemplate = `
  async #update<T extends {{.Options}}>(
    operation: string,
    method: ((id: string, request: T, options: OperationOptions) => Promise<SandboxInfo>) | undefined,
    request: T,
    options: OperationOptions,
  ): Promise<SandboxInfo> {
    return this.#serialized(operation, options, request.timeout ?? this.#timeoutMs, async (signal) => {
      if (!method) unsupported(this.providerName, operation);
      const info = await awaitWithSignal(
        () => method.call(this.#backend, this.id, clone(request), operationOptions(options, signal)),
        signal,
      );
      this.#info = clone({ ...info, provider: this.providerName });
      return this.info();
    });
  }
`

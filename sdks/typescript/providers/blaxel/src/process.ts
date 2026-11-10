import {
  ErrorKind,
  ExecFieldTimeout,
  sandboxError,
  type ExecOptions,
  type ExecResult,
} from "@sandbox-kit/sandbox";
import { providerError } from "./errors.js";
import type {
  BlaxelProcessResult,
  NativeSandboxHandle,
  ProcessExecRequest,
} from "./native.js";

/** Blaxel's documented limit for a single waitForCompletion request. */
const blaxelCompletionHoldSeconds = 60;

const terminalProcessStatus = new Set([
  "completed",
  "failed",
  "killed",
  "stopped",
]);

export function processRequest(
  request: ExecOptions,
  timeout: number | undefined,
): ProcessExecRequest {
  const body: ProcessExecRequest = {
    command: blaxelCommand(request),
    // 1–60s can use Blaxel's completion hold; omit/0/longer timeouts poll.
    waitForCompletion: holdsCompletion(timeout),
    // omit/0 → 0 disables Blaxel's default process auto-kill.
    timeout: timeout ?? 0,
  };
  if (request.workingDirectory !== undefined) {
    body.workingDir = request.workingDirectory;
  }
  if (request.environment !== undefined) {
    body.env = { ...request.environment };
  }
  return body;
}

/** Internal wait orchestration after Process.exec. */
export interface FinishProcessArgs {
  process: NativeSandboxHandle["process"];
  started: BlaxelProcessResult;
  timeout: number | undefined;
  signal: AbortSignal | undefined;
}

export async function finishProcess(
  args: FinishProcessArgs,
): Promise<ExecResult> {
  if (processFinished(args.started.status)) return execResult(args.started);
  const identifier = args.started.pid || args.started.name;
  if (!identifier) {
    throw sandboxError({
      kind: ErrorKind.InvalidResponse,
      provider: "blaxel",
      operation: "run",
      field: "",
      message: "sandbox-kit blaxel: SDK returned no process identifier",
    });
  }
  const waitOptions: {
    maxWait: number;
    interval: number;
    signal?: AbortSignal;
  } = {
    maxWait: args.timeout === undefined ? -1 : args.timeout * 1_000,
    interval: 1_000,
  };
  if (args.signal !== undefined) waitOptions.signal = args.signal;
  return execResult(await args.process.wait(identifier, waitOptions));
}

/**
 * Maps Kit ExecOptions.timeout (ms) to Blaxel seconds.
 * Omitted and explicit zero become undefined so the caller sends 0 (no limit).
 */
export function execTimeoutSeconds(
  timeout: number | undefined,
): number | undefined {
  if (timeout === undefined || timeout === 0) return undefined;
  if (!isWholeSecondTimeout(timeout)) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "run",
      field: ExecFieldTimeout,
      message: "timeout must be a whole number of seconds",
    });
  }
  return timeout / 1_000;
}

function isWholeSecondTimeout(timeout: number): boolean {
  return Number.isFinite(timeout) && timeout >= 0 && timeout % 1_000 === 0;
}

function holdsCompletion(timeoutSeconds: number | undefined): boolean {
  return (
    timeoutSeconds !== undefined &&
    timeoutSeconds > 0 &&
    timeoutSeconds <= blaxelCompletionHoldSeconds
  );
}

function processFinished(status: string | undefined): boolean {
  return status !== undefined && terminalProcessStatus.has(status);
}

function execResult(result: BlaxelProcessResult): ExecResult {
  return {
    exitCode: result.exitCode ?? 0,
    stdout: new TextEncoder().encode(result.stdout ?? ""),
    stderr: new TextEncoder().encode(result.stderr ?? ""),
  };
}

function blaxelCommand(request: ExecOptions): string {
  if (request.shell !== undefined) return request.shell;
  return posixShellJoin(request.argv ?? []);
}

function posixShellJoin(argv: readonly string[]): string {
  return argv.map(posixQuote).join(" ");
}

function posixQuote(value: string): string {
  if (value === "") return "''";
  return "'" + value.replaceAll("'", `'"'"'`) + "'";
}

// Process calls are handwritten. They validate argv/shell and working
// directory, then dispatch through optional Backend.run.

import {
  ExecFieldArgv,
  ExecFieldShell,
  ExecFieldWorkingDirectory,
} from "#contracts/field_paths.js";
import type {
  ExecOptions,
  ExecResult,
} from "#contracts/sandbox_data_contracts.js";
import { ErrorKind, sandboxError, unsupported } from "#client/errors.js";
import { operationOptions } from "#client/operation.js";
import type { Backend, OperationOptions } from "#client/types.js";

type ProcessCall = ExecOptions & OperationOptions;

interface ProcessHost {
  readonly providerName: string;
  readonly id: string;
  readonly backend: Backend;
  run<T>(
    operation: string,
    options: OperationOptions,
    body: (signal: AbortSignal | undefined) => Promise<T>,
  ): Promise<T>;
}

/**
 * Runs commands inside one sandbox and waits for each to finish. Set exactly
 * one of `argv` or `shell`. A working directory, when set, must be absolute.
 * An omitted operation timeout uses the client timeout; explicit zero adds no
 * Kit deadline.
 */
export class Process {
  constructor(private readonly host: ProcessHost) {}

  /**
   * Runs one command and returns its exit code and collected output.
   * A non-zero exit code is not an error.
   */
  async run(options: ProcessCall): Promise<ExecResult> {
    validateExecOptions(this.host.providerName, "run", options);
    return this.host.run("run", options, (signal) => {
      const method = this.host.backend.run;
      if (!method) unsupported(this.host.providerName, "run");
      return method.call(
        this.host.backend,
        this.host.id,
        requestFields(options),
        operationOptions(options, signal),
      );
    });
  }
}

function validateExecOptions(
  provider: string,
  operation: string,
  options: ExecOptions,
): void {
  const hasArgv = Array.isArray(options.argv) && options.argv.length > 0;
  const hasShell = options.shell !== undefined;
  if (hasArgv && hasShell) {
    throw sandboxError({
      kind: ErrorKind.InvalidArgument,
      provider,
      operation,
      field: ExecFieldArgv,
      message: "sandbox-kit: set exactly one of argv or shell",
    });
  }
  if (!hasArgv && !hasShell) {
    throw sandboxError({
      kind: ErrorKind.InvalidArgument,
      provider,
      operation,
      field: ExecFieldArgv,
      message: "sandbox-kit: set argv or shell",
    });
  }
  if (hasShell && options.shell!.trim() === "") {
    throw sandboxError({
      kind: ErrorKind.InvalidArgument,
      provider,
      operation,
      field: ExecFieldShell,
      message: "sandbox-kit: shell must not be blank",
    });
  }
  if (options.workingDirectory !== undefined) {
    const path = options.workingDirectory;
    if (!path || !path.startsWith("/") || path.includes("\0")) {
      throw sandboxError({
        kind: ErrorKind.InvalidArgument,
        provider,
        operation,
        field: ExecFieldWorkingDirectory,
        message: "sandbox-kit: working directory must be absolute",
      });
    }
  }
}

function requestFields(options: ProcessCall): ExecOptions {
  const request: ExecOptions = {};
  if (options.argv !== undefined) request.argv = options.argv;
  if (options.shell !== undefined) request.shell = options.shell;
  if (options.workingDirectory !== undefined) {
    request.workingDirectory = options.workingDirectory;
  }
  if (options.environment !== undefined) {
    request.environment = options.environment;
  }
  if (options.timeout !== undefined) request.timeout = options.timeout;
  return request;
}

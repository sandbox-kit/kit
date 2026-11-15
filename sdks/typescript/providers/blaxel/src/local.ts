import { mkdir, readFile, stat } from "node:fs/promises";
import { dirname, resolve } from "node:path";

import {
  ErrorKind,
  FileFieldLocalPath,
  SandboxError,
  sandboxError,
} from "@sandbox-kit/sandbox";
import { providerError } from "./errors.js";

/** Rejects a missing local file or a directory before Blaxel reads it. */
async function requireLocalFile(
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

/**
 * Classifies a local destination path failure from Blaxel download as
 * `local_path`. Matches Go: only when the error path equals `localPath`.
 * Other errors are rethrown for `mapBlaxelError`.
 */
export function mapCallerPathError(
  error: unknown,
  operation: string,
  localPath: string,
): never {
  if (error instanceof SandboxError) throw error;
  const pathErr = nodePathError(error) ?? nodePathError(causeOf(error));
  if (pathErr !== undefined && resolve(pathErr.path) === resolve(localPath)) {
    throw localFileError(error, operation);
  }
  throw error;
}

/** Node `ErrnoException`-style path error (Go `os.PathError` counterpart). */
function nodePathError(error: unknown): { path: string } | undefined {
  if (typeof error !== "object" || !error || !("path" in error))
    return undefined;
  const path = error.path;
  return typeof path === "string" && path !== "" ? { path } : undefined;
}

function causeOf(error: unknown): unknown {
  return error instanceof Error ? error.cause : undefined;
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
  const code = nodeCode(error) ?? nodeCode(causeOf(error));
  return sandboxError(
    {
      kind: code === "ENOENT" ? ErrorKind.NotFound : ErrorKind.InvalidArgument,
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

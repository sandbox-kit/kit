import { mkdir, readFile, stat } from "node:fs/promises";
import { dirname } from "node:path";

import {
  ErrorKind,
  FileFieldLocalPath,
  SandboxError,
  sandboxError,
} from "@sandbox-kit/sandbox";
import { providerError } from "./errors.js";

/** Rejects a missing local file or a directory before Blaxel reads it. */
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

import { basename, dirname } from "node:path";
import type { Directory, File as BlaxelFile } from "@blaxel/core";
import {
  ErrorKind,
  FileFieldPath,
  FileType,
  sandboxError,
  type FileInfo,
} from "@sandbox-kit/sandbox";
import type { NativeSandboxHandle } from "./native.js";

type BlaxelSubdirectory = Directory["subdirectories"][number];

export function fileInfo(info: BlaxelFile): FileInfo {
  const modified = Date.parse(info.lastModified);
  return {
    name: info.name,
    path: info.path,
    type: FileType.File,
    size: BigInt(Math.trunc(info.size)),
    ...(info.permissions ? { permissions: info.permissions } : {}),
    ...(info.owner ? { owner: info.owner } : {}),
    ...(info.group ? { group: info.group } : {}),
    ...(Number.isNaN(modified) ? {} : { modifiedAt: new Date(modified) }),
  };
}

export function directoryInfo(info: BlaxelSubdirectory): FileInfo {
  return {
    name: info.name,
    path: info.path,
    type: FileType.Directory,
    size: 0n,
  };
}

/** Blaxel has no Stat; ls wraps GET and rejects file paths with this error. */
export async function blaxelStat(
  fs: NativeSandboxHandle["fs"],
  path: string,
): Promise<FileInfo> {
  try {
    const listing = await fs.ls(path);
    return {
      name: listing.name || basename(path) || "/",
      path: listing.path || path,
      type: FileType.Directory,
      size: 0n,
    };
  } catch (error) {
    if (!isBlaxelFilePath(error)) throw error;
  }
  const parent = dirname(path);
  const name = basename(path);
  const listing = await fs.ls(parent === "." ? "/" : parent);
  const file = listing.files.find((entry) => entry.name === name);
  if (file) return fileInfo(file);
  throw sandboxError({
    kind: ErrorKind.NotFound,
    provider: "blaxel",
    operation: "stat",
    field: FileFieldPath,
    message: `sandbox-kit blaxel: path not found: ${path}`,
  });
}

/** Blaxel `fs.ls` throws this when GET returns a file instead of a directory. */
function isBlaxelFilePath(error: unknown): boolean {
  if (!(error instanceof Error)) return false;
  try {
    const body = JSON.parse(error.message) as { error?: unknown };
    return body.error === "Directory not found";
  } catch {
    return false;
  }
}

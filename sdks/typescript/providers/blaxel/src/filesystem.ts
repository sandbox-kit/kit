import { basename } from "node:path";
import type { Directory, File as BlaxelFile } from "@blaxel/core";
import { FileType, type FileInfo } from "@sandbox-kit/sandbox";
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

/**
 * Stat one path. Uses the same filesystem GET as Go: a directory listing or
 * file metadata (size, permissions when Blaxel supplies them).
 */
export async function blaxelStat(
  fs: NativeSandboxHandle["fs"],
  path: string,
): Promise<FileInfo> {
  const data = await fs.get(path);
  if (isDirectoryListing(data)) {
    return {
      name: data.name || basename(path) || "/",
      path: data.path || path,
      type: FileType.Directory,
      size: 0n,
    };
  }
  return fileInfo({
    name: data.name || basename(path) || path,
    path: data.path || path,
    size: data.size,
    permissions: data.permissions ?? "",
    owner: data.owner ?? "",
    group: data.group ?? "",
    lastModified: data.lastModified ?? "",
  });
}

function isDirectoryListing(data: Directory | BlaxelFile): data is Directory {
  return "files" in data || "subdirectories" in data;
}

import type { Directory, File as BlaxelFile } from "@blaxel/core";
import {
  ErrorKind,
  FileFieldPath,
  FileType,
  type FileInfo,
} from "@sandbox-kit/sandbox";
import { providerError } from "./errors.js";
import type { NativeSandboxHandle } from "./native.js";

type BlaxelSubdirectory = Directory["subdirectories"][number];

export function fileInfo(info: BlaxelFile, operation = "stat"): FileInfo {
  requireFileIdentity(info, operation);
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
 * file metadata (size, permissions when Blaxel supplies them). Does not invent
 * name or path when Blaxel omits them.
 */
export async function blaxelStat(
  fs: NativeSandboxHandle["fs"],
  path: string,
): Promise<FileInfo> {
  const data = await fs.get(path);
  if (isDirectoryListing(data)) {
    return {
      name: data.name ?? "",
      path: data.path ?? "",
      type: FileType.Directory,
      size: 0n,
    };
  }
  return fileInfo(
    {
      name: data.name ?? "",
      path: data.path ?? "",
      size: data.size,
      permissions: data.permissions ?? "",
      owner: data.owner ?? "",
      group: data.group ?? "",
      lastModified: data.lastModified ?? "",
    },
    "stat",
  );
}

/** Match Go blaxelFileInfo: empty name and path is not usable file metadata. */
function requireFileIdentity(
  info: Pick<BlaxelFile, "name" | "path">,
  operation: string,
): void {
  if (!info.name && !info.path) {
    providerError({
      kind: ErrorKind.InvalidResponse,
      operation,
      field: FileFieldPath,
      message: "SDK returned no file metadata",
    });
  }
}

function isDirectoryListing(data: Directory | BlaxelFile): data is Directory {
  return "files" in data || "subdirectories" in data;
}

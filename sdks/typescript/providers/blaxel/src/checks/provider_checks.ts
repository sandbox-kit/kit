import {
  ErrorKind,
  CreateFieldEnvironment,
  CreateFieldName,
  CreateFieldPlacement,
  CreateFieldResourcesCpuCores,
  CreateFieldResourcesCpuLimitCores,
  CreateFieldResourcesDiskMib,
  CreateFieldResourcesMemoryLimitMib,
  CreateFieldSource,
  CreateFieldSourceImage,
  validEnvName,
} from "@sandbox-kit/sandbox";
import type { CreateOptions } from "@sandbox-kit/sandbox";
import { providerError } from "#errors.js";

// validateProviderCreate rejects creation settings this provider cannot honor.
// Rules are inventoried in specs/providers/blaxel.yaml checks when present.
export function validateProviderCreate(request: CreateOptions): void {
  if (request.name === undefined || request.name.trim() === "") {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldName,
      message: "name is required",
    });
  }
  if (request.source?.image === undefined) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldSourceImage,
      message: "an explicit image source is required",
    });
  }
  if (
    request.source?.snapshot !== undefined ||
    request.source?.warmPool !== undefined
  ) {
    providerError({
      kind: ErrorKind.Unsupported,
      operation: "create",
      field: CreateFieldSource,
      message: "only an image source is supported",
    });
  }
  if (request.placement !== undefined) {
    providerError({
      kind: ErrorKind.Unsupported,
      operation: "create",
      field: CreateFieldPlacement,
      message:
        "per-request Blaxel placement is unsupported; use Config.region",
    });
  }
  const resources = request.resources;
  if (resources !== undefined) {
    if (resources.cpuCores !== undefined) {
      providerError({
        kind: ErrorKind.Unsupported,
        operation: "create",
        field: CreateFieldResourcesCpuCores,
        message: "Blaxel derives CPU from memory; cpu_cores is not mapped",
      });
    }
    if (resources.cpuLimitCores !== undefined) {
      providerError({
        kind: ErrorKind.Unsupported,
        operation: "create",
        field: CreateFieldResourcesCpuLimitCores,
        message: "Blaxel resource limits are not mapped",
      });
    }
    if (resources.memoryLimitMib !== undefined) {
      providerError({
        kind: ErrorKind.Unsupported,
        operation: "create",
        field: CreateFieldResourcesMemoryLimitMib,
        message: "Blaxel resource limits are not mapped",
      });
    }
    if (resources.diskMib !== undefined) {
      providerError({
        kind: ErrorKind.Unsupported,
        operation: "create",
        field: CreateFieldResourcesDiskMib,
        message: "Blaxel disk size is not mapped as MiB",
      });
    }
  }
  for (const value of Object.keys(request.environment ?? {})) {
    if (!validEnvName(value)) {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation: "create",
        field: CreateFieldEnvironment,
        message: "environment is outside supported values",
      });
    }
  }
}

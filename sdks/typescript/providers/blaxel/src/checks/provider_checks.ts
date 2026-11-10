import {
  ErrorKind,
  CreateFieldEnvironment,
  CreateFieldName,
  CreateFieldNetworkOutboundDomainsEntries,
  CreateFieldPlacement,
  CreateFieldResourcesCpuCores,
  CreateFieldResourcesCpuLimitCores,
  CreateFieldResourcesDiskMib,
  CreateFieldResourcesMemoryLimitMib,
  CreateFieldSource,
  CreateFieldSourceImage,
  validDomain,
  validEnvName,
  type ErrorInfo,
} from "@sandbox-kit/sandbox";
import type { CreateOptions, Resources } from "@sandbox-kit/sandbox";
import { providerError } from "#errors.js";

// Blaxel Metadata.Name: lowercase alphanumeric with single hyphens, max 49.
const PROVIDER_SANDBOX_NAME = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;
const PROVIDER_SANDBOX_NAME_MAX = 49;

// validateProviderCreate rejects creation settings this provider cannot honor.
// Inventoried in specs/providers/blaxel.yaml checks when present.
export function validateProviderCreate(request: CreateOptions): void {
  if (request.name === undefined || request.name.trim() === "") {
    reject(ErrorKind.InvalidArgument, CreateFieldName, "name is required");
  }
  if (
    request.name.length > PROVIDER_SANDBOX_NAME_MAX ||
    !PROVIDER_SANDBOX_NAME.test(request.name)
  ) {
    reject(
      ErrorKind.InvalidArgument,
      CreateFieldName,
      "name must be lowercase alphanumeric with single hyphens, at most 49 characters",
    );
  }
  if (
    request.source?.image === undefined ||
    !request.source.image.reference?.trim()
  ) {
    reject(
      ErrorKind.InvalidArgument,
      CreateFieldSourceImage,
      "an explicit image source is required",
    );
  }
  if (
    request.source?.snapshot !== undefined ||
    request.source?.warmPool !== undefined
  ) {
    reject(
      ErrorKind.Unsupported,
      CreateFieldSource,
      "only an image source is supported",
    );
  }
  if (request.placement !== undefined) {
    reject(
      ErrorKind.Unsupported,
      CreateFieldPlacement,
      "per-request Blaxel placement is unsupported; use Config.region",
    );
  }
  rejectUnsupportedResources(request.resources);
  for (const value of Object.keys(request.environment ?? {})) {
    if (!validEnvName(value)) {
      reject(
        ErrorKind.InvalidArgument,
        CreateFieldEnvironment,
        "environment is outside supported values",
      );
    }
  }
  for (const value of request.network?.outboundDomains ?? []) {
    if (!validDomain(value)) {
      reject(
        ErrorKind.InvalidArgument,
        CreateFieldNetworkOutboundDomainsEntries,
        "network.outbound_domains.entries is outside supported values",
      );
    }
  }
}

function rejectUnsupportedResources(resources: Resources | undefined): void {
  if (resources === undefined) return;
  const checks: ReadonlyArray<{
    set: boolean;
    field: string;
    message: string;
  }> = [
    {
      set: resources.cpuCores !== undefined,
      field: CreateFieldResourcesCpuCores,
      message: "Blaxel derives CPU from memory; cpu_cores is not mapped",
    },
    {
      set: resources.cpuLimitCores !== undefined,
      field: CreateFieldResourcesCpuLimitCores,
      message: "Blaxel resource limits are not mapped",
    },
    {
      set: resources.memoryLimitMib !== undefined,
      field: CreateFieldResourcesMemoryLimitMib,
      message: "Blaxel resource limits are not mapped",
    },
    {
      set: resources.diskMib !== undefined,
      field: CreateFieldResourcesDiskMib,
      message: "Blaxel disk size is not mapped as MiB",
    },
  ];
  for (const check of checks) {
    if (check.set) reject(ErrorKind.Unsupported, check.field, check.message);
  }
}

function reject(
  kind: ErrorInfo["kind"],
  field: string,
  message: string,
): never {
  providerError({ kind, operation: "create", field, message });
}

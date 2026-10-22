import {
  ErrorKind,
  CreateFieldEnvironment,
  CreateFieldNetworkInboundCidrsEntries,
  CreateFieldNetworkOutboundCidrsEntries,
  CreateFieldNetworkOutboundDomainsEntries,
  CreateFieldPlacementCloud,
  CreateFieldResourcesCpuCores,
  CreateFieldResourcesMemoryMib,
  CreateFieldRuntimeWorkingDirectory,
  validCidr,
  validDomain,
  validEnvName,
} from "@sandbox-kit/sandbox";
import type { CreateOptions } from "@sandbox-kit/sandbox";
import { providerError } from "#errors.js";

// validateProviderCreate rejects creation settings this provider cannot honor.
// Rules are inventoried in specs/providers/modal.yaml checks.
export function validateProviderCreate(request: CreateOptions): void {
  const cpuCores = request.resources?.cpuCores;
  if (cpuCores !== undefined && cpuCores < 0.125) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldResourcesCpuCores,
      message: "resources.cpu_cores is outside supported values",
    });
  }
  const memoryMib = request.resources?.memoryMib;
  if (memoryMib !== undefined && memoryMib < BigInt("128")) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldResourcesMemoryMib,
      message: "resources.memory_mib is outside supported values",
    });
  }
  const workingDirectory = request.runtime?.workingDirectory;
  if (workingDirectory !== undefined && !workingDirectory.startsWith("/")) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldRuntimeWorkingDirectory,
      message: "runtime.working_directory is outside supported values",
    });
  }
  const cloud = request.placement?.cloud;
  if (
    cloud !== undefined &&
    cloud !== "aws" &&
    cloud !== "gcp" &&
    cloud !== "oci" &&
    cloud !== "auto"
  ) {
    providerError({
      kind: ErrorKind.Unsupported,
      operation: "create",
      field: CreateFieldPlacementCloud,
      message: "placement.cloud is outside supported values",
    });
  }
  for (const value of request.network?.outboundCidrs ?? []) {
    if (value !== "" && !validCidr(value)) {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation: "create",
        field: CreateFieldNetworkOutboundCidrsEntries,
        message: "network.outbound_cidrs.entries is outside supported values",
      });
    }
  }
  for (const value of request.network?.inboundCidrs ?? []) {
    if (!validCidr(value)) {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation: "create",
        field: CreateFieldNetworkInboundCidrsEntries,
        message: "network.inbound_cidrs.entries is outside supported values",
      });
    }
  }
  for (const value of request.network?.outboundDomains ?? []) {
    if (value !== "" && !validDomain(value)) {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation: "create",
        field: CreateFieldNetworkOutboundDomainsEntries,
        message: "network.outbound_domains.entries is outside supported values",
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

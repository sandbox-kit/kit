import {
  ErrorKind,
  sandboxError,
  type CreateOptions,
} from "@sandbox-kit/sandbox";

// validateProviderCreate rejects creation settings this provider cannot honor.
// Rules are inventoried in specs/providers/modal.yaml checks.
export function validateProviderCreate(request: CreateOptions): void {
  const cpuCores = request.resources?.cpuCores;
  if (cpuCores !== undefined && cpuCores < 0.125) {
    fail(
      ErrorKind.InvalidArgument,
      "resources.cpuCores",
      "resources.cpu_cores is outside supported values",
    );
  }
  const memoryMib = request.resources?.memoryMib;
  if (memoryMib !== undefined && memoryMib < BigInt("128")) {
    fail(
      ErrorKind.InvalidArgument,
      "resources.memoryMib",
      "resources.memory_mib is outside supported values",
    );
  }
  const workingDirectory = request.runtime?.workingDirectory;
  if (workingDirectory !== undefined && !workingDirectory.startsWith("/")) {
    fail(
      ErrorKind.InvalidArgument,
      "runtime.workingDirectory",
      "runtime.working_directory is outside supported values",
    );
  }
  const cloud = request.placement?.cloud;
  if (
    cloud !== undefined &&
    cloud !== "aws" &&
    cloud !== "gcp" &&
    cloud !== "oci" &&
    cloud !== "auto"
  ) {
    fail(
      ErrorKind.Unsupported,
      "placement.cloud",
      "placement.cloud is outside supported values",
    );
  }
  for (const value of request.network?.outboundCidrs?.entries ?? []) {
    if (value !== "" && !validCidr(value)) {
      fail(
        ErrorKind.InvalidArgument,
        "network.outboundCidrs.entries",
        "network.outbound_cidrs.entries is outside supported values",
      );
    }
  }
  for (const value of request.network?.inboundCidrs?.entries ?? []) {
    if (!validCidr(value)) {
      fail(
        ErrorKind.InvalidArgument,
        "network.inboundCidrs.entries",
        "network.inbound_cidrs.entries is outside supported values",
      );
    }
  }
  for (const value of request.network?.outboundDomains?.entries ?? []) {
    if (value !== "" && !validDomain(value)) {
      fail(
        ErrorKind.InvalidArgument,
        "network.outboundDomains.entries",
        "network.outbound_domains.entries is outside supported values",
      );
    }
  }
  for (const value of Object.keys(request.environment ?? {})) {
    if (!validEnvName(value)) {
      fail(
        ErrorKind.InvalidArgument,
        "environment",
        "environment is outside supported values",
      );
    }
  }
}

function validEnvName(value: string): boolean {
  return /^[A-Za-z_][A-Za-z0-9_]*$/.test(value);
}

function validDomain(value: string): boolean {
  const domain = value.startsWith("*.") ? value.slice(2) : value;
  return (
    domain.length <= 253 &&
    /^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.?$/i.test(
      domain,
    )
  );
}

function validCidr(value: string): boolean {
  const separator = value.lastIndexOf("/");
  if (separator <= 0) return false;
  const address = value.slice(0, separator);
  const prefix = Number(value.slice(separator + 1));
  if (!Number.isInteger(prefix)) return false;
  return address.includes(":")
    ? prefix >= 0 && prefix <= 128 && validIpv6(address)
    : prefix >= 0 && prefix <= 32 && validIpv4(address);
}

function validIpv4(value: string): boolean {
  const parts = value.split(".");
  return (
    parts.length === 4 &&
    parts.every(
      (part) => /^(0|[1-9][0-9]{0,2})$/.test(part) && Number(part) <= 255,
    )
  );
}

function validIpv6(value: string): boolean {
  if ((value.match(/::/g)?.length ?? 0) > 1) return false;
  const halves = value.split("::");
  if (halves.length > 2) return false;
  const count = (part: string): number | undefined => {
    if (part === "") return 0;
    const groups = part.split(":");
    let total = 0;
    for (const group of groups) {
      if (/^[0-9a-f]{1,4}$/i.test(group)) total++;
      else if (validIpv4(group)) total += 2;
      else return undefined;
    }
    return total;
  };
  const left = count(halves[0] ?? "");
  const right = count(halves[1] ?? "");
  return (
    left !== undefined &&
    right !== undefined &&
    (halves.length === 2 ? left + right < 8 : left === 8)
  );
}

function fail(kind: ErrorKind, field: string, message: string): never {
  throw sandboxError({
    kind,
    provider: "modal",
    operation: "create",
    field,
    message: "sandbox-kit modal: " + message,
  });
}

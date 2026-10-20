import {
  ErrorKind,
  PolicyMode,
  sandboxError
} from "@sandbox-kit/sandbox";
import type {
  CreateOptions
} from "@sandbox-kit/sandbox";

type PublicAction = NonNullable<
  NonNullable<CreateOptions["lifetime"]>["idleStop"]
>;

// validateProviderCreate rejects creation settings this provider cannot honor.
// Rules are inventoried in specs/providers/daytona.yaml checks.
export function validateProviderCreate(request: CreateOptions): void {
  const language = request.runtime?.language;
  if (
    language !== undefined &&
    language !== "python" &&
    language !== "javascript" &&
    language !== "typescript"
  ) {
    fail(
      ErrorKind.Unsupported,
      "runtime.language",
      "runtime.language is outside supported values",
    );
  }
  const outboundProxyUrl = request.network?.outboundProxyUrl;
  if (outboundProxyUrl !== undefined && !validHttpUrl(outboundProxyUrl)) {
    fail(
      ErrorKind.InvalidArgument,
      "network.outboundProxyUrl",
      "network.outbound_proxy_url is outside supported values",
    );
  }
  for (const value of request.network?.outboundCidrs ?? []) {
    if (!validCidr(value)) {
      fail(
        ErrorKind.InvalidArgument,
        "network.outboundCidrs.entries",
        "network.outbound_cidrs.entries is outside supported values",
      );
    }
  }
  for (const value of request.network?.outboundDomains ?? []) {
    if (!validDomain(value)) {
      fail(
        ErrorKind.InvalidArgument,
        "network.outboundDomains.entries",
        "network.outbound_domains.entries is outside supported values",
      );
    }
  }
  for (const item of request.security?.secrets ?? []) {
    const environmentVariable = item.environmentVariable;
    if (
      environmentVariable !== undefined &&
      !validEnvName(environmentVariable)
    ) {
      fail(
        ErrorKind.InvalidArgument,
        "security.secrets.environmentVariable",
        "security.secrets.environment_variable is outside supported values",
      );
    }
  }
  if (
    delayed(request.lifetime?.idleStop) &&
    delayed(request.lifetime?.idlePause)
  ) {
    invalid("lifetime", "incompatible lifetime policies");
  }
  if (
    delayed(request.lifetime?.idlePause) &&
    request.lifetime?.ephemeral === true
  ) {
    invalid("lifetime", "incompatible lifetime policies");
  }
  if (
    delayed(request.lifetime?.idlePause) &&
    immediate(request.lifetime?.stoppedDelete)
  ) {
    invalid("lifetime", "incompatible lifetime policies");
  }
  if (
    delayed(request.lifetime?.stoppedDelete) &&
    request.lifetime?.ephemeral === true
  ) {
    invalid("lifetime", "incompatible lifetime policies");
  }
  if (
    delayed(request.lifetime?.stoppedArchive) &&
    request.lifetime?.ephemeral === true
  ) {
    invalid("lifetime", "incompatible lifetime policies");
  }
}

function delayed(policy: PublicAction | undefined): boolean {
  return (
    (policy?.mode ?? PolicyMode.Default) === PolicyMode.After &&
    policy?.after !== undefined &&
    policy.after > 0
  );
}

function immediate(policy: PublicAction | undefined): boolean {
  return (
    (policy?.mode ?? PolicyMode.Default) === PolicyMode.After &&
    policy?.after !== undefined &&
    policy.after === 0
  );
}

function validEnvName(value: string): boolean {
  return /^[A-Za-z_][A-Za-z0-9_]*$/.test(value);
}

function validHttpUrl(value: string): boolean {
  try {
    const parsed = new URL(value);
    return (
      parsed.hostname !== "" &&
      (parsed.protocol === "http:" || parsed.protocol === "https:")
    );
  } catch {
    return false;
  }
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

function invalid(field: string, message: string): never {
  fail(ErrorKind.InvalidArgument, field, message);
}

function fail(kind: ErrorKind, field: string, message: string): never {
  throw sandboxError({
    kind,
    provider: "daytona",
    operation: "create",
    field,
    message: "sandbox-kit daytona: " + message,
  });
}

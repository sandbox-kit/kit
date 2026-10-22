import {
  ErrorKind,
  PolicyMode,
  CreateFieldLifetime,
  CreateFieldNetworkOutboundCidrsEntries,
  CreateFieldNetworkOutboundDomainsEntries,
  CreateFieldNetworkOutboundProxyUrl,
  CreateFieldRuntimeLanguage,
  CreateFieldSecuritySecretsEnvironmentVariable,
  validCidr,
  validDomain,
  validEnvName,
} from "@sandbox-kit/sandbox";
import type { CreateOptions } from "@sandbox-kit/sandbox";
import { providerError } from "#errors.js";

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
    providerError({
      kind: ErrorKind.Unsupported,
      operation: "create",
      field: CreateFieldRuntimeLanguage,
      message: "runtime.language is outside supported values",
    });
  }
  const outboundProxyUrl = request.network?.outboundProxyUrl;
  if (outboundProxyUrl !== undefined && !validHttpUrl(outboundProxyUrl)) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldNetworkOutboundProxyUrl,
      message: "network.outbound_proxy_url is outside supported values",
    });
  }
  for (const value of request.network?.outboundCidrs ?? []) {
    if (!validCidr(value)) {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation: "create",
        field: CreateFieldNetworkOutboundCidrsEntries,
        message: "network.outbound_cidrs.entries is outside supported values",
      });
    }
  }
  for (const value of request.network?.outboundDomains ?? []) {
    if (!validDomain(value)) {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation: "create",
        field: CreateFieldNetworkOutboundDomainsEntries,
        message: "network.outbound_domains.entries is outside supported values",
      });
    }
  }
  for (const item of request.security?.secrets ?? []) {
    const environmentVariable = item.environmentVariable;
    if (
      environmentVariable !== undefined &&
      !validEnvName(environmentVariable)
    ) {
      providerError({
        kind: ErrorKind.InvalidArgument,
        operation: "create",
        field: CreateFieldSecuritySecretsEnvironmentVariable,
        message:
          "security.secrets.environment_variable is outside supported values",
      });
    }
  }
  if (
    delayed(request.lifetime?.idleStop) &&
    delayed(request.lifetime?.idlePause)
  ) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldLifetime,
      message: "incompatible lifetime policies",
    });
  }
  if (
    delayed(request.lifetime?.idlePause) &&
    request.lifetime?.ephemeral === true
  ) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldLifetime,
      message: "incompatible lifetime policies",
    });
  }
  if (
    delayed(request.lifetime?.idlePause) &&
    immediate(request.lifetime?.stoppedDelete)
  ) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldLifetime,
      message: "incompatible lifetime policies",
    });
  }
  if (
    delayed(request.lifetime?.stoppedDelete) &&
    request.lifetime?.ephemeral === true
  ) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldLifetime,
      message: "incompatible lifetime policies",
    });
  }
  if (
    delayed(request.lifetime?.stoppedArchive) &&
    request.lifetime?.ephemeral === true
  ) {
    providerError({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldLifetime,
      message: "incompatible lifetime policies",
    });
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

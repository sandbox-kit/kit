// Shared create/config validation is Protovalidate on materialized wire
// messages; secure_endpoint stays handwritten.

import {
  create,
  type MessageInitShape,
  type MessageShape,
} from "@bufbuild/protobuf";
import { pathToString } from "@bufbuild/protobuf/reflect";
import { durationFromMs } from "@bufbuild/protobuf/wkt";
import { createValidator } from "@bufbuild/protovalidate";
import {
  ConfigSchema,
  type APIKeyCredentials as WireAPIKeyCredentials,
  type AuthConfig,
  type BearerTokenCredentials as WireBearerTokenCredentials,
  type Config as WireConfig,
  type OAuthCredentials as WireOAuthCredentials,
  type TokenPairCredentials as WireTokenPairCredentials,
} from "@sandbox-kit/contracts/client";
import {
  CreateOptionsSchema,
  type CommandProbe,
  type CreateOptions as ProtoCreateOptions,
  type ImageSource,
  type ProviderDefaultSource,
  type ReadinessProbe,
  type SandboxSource,
  type SnapshotSource,
  type WarmPoolSource,
} from "@sandbox-kit/contracts/sandbox";
import type {
  CreateOptions,
  NetworkConfigInput,
  ReadinessProbeInput,
  SandboxSourceInput,
} from "#contracts/sandbox_data_contracts.js";
import { materializeDurationFields } from "#internal/duration_ms.js";
import type { StripProto } from "#internal/strip_proto.js";
import { metadataObjectFromMap } from "#metadata/metadata_object.js";
import { ErrorKind, sandboxError } from "#client/errors.js";
import type { AuthInput, Config } from "#client/types.js";

const contractValidator = createValidator({ failFast: true });

const stableRuleIDs: Record<string, string> = {
  "enum.defined_only": "known_enum",
  "duration.gte": "minimum",
  "oneof.required": "exactly_one",
  auth_credentials_required: "exactly_one",
  cpu_limit: "limit",
  memory_limit: "limit",
  cpu_requires_positive: "requires_positive",
  memory_requires_positive: "requires_positive",
  restricted_requires_allowlist: "requires_any",
  block_all_forbids_network: "forbids",
  allow_all_forbids_allowlists: "forbids",
  after_requires_duration: "requires_any",
  after_forbids_duration: "forbids",
};

/** Wire config after materialize; package-private for Protovalidate. */
type ConfigWire = StripProto<WireConfig>;
/** Wire create options after materialize; package-private for Protovalidate. */
type CreateOptionsWire = StripProto<ProtoCreateOptions>;

function validationError(
  provider: string,
  operation: string,
  field: string,
  rule: string,
  message: string,
): never {
  throw sandboxError({
    kind: ErrorKind.InvalidArgument,
    provider,
    operation,
    field,
    rule,
    message: "sandbox-kit: " + message,
  });
}

function validSecureEndpoint(value: string): boolean {
  try {
    const url = new URL(value);
    if (!url.hostname || url.username || url.password || url.hash) return false;
    if (url.protocol === "https:") return true;
    if (url.protocol !== "http:") return false;
    const host = url.hostname.toLowerCase();
    return (
      host === "localhost" ||
      host === "[::1]" ||
      host === "::1" ||
      /^127(?:\.[0-9]{1,3}){3}$/.test(host)
    );
  } catch {
    return false;
  }
}

/**
 * Rejects public {@link CreateOptions} the shared contract cannot represent.
 * Materializes wire form privately for Protovalidate. Provider-specific checks
 * remain in each provider adapter.
 */
export function validateCreateOptions(
  value: CreateOptions,
  provider = "",
  operation = "create",
): void {
  const message = materializeCreateOptions(value);
  const result = contractValidator.validate(
    CreateOptionsSchema,
    message as MessageShape<typeof CreateOptionsSchema>,
  );
  if (result.kind === "valid") {
    return;
  }
  if (result.kind === "error") {
    throw sandboxError(
      {
        kind: ErrorKind.InvalidArgument,
        provider,
        operation,
        field: "",
        message: "sandbox-kit: invalid configuration",
      },
      result.error,
    );
  }
  const violation = result.violations[0]!;
  const ruleID = violation.ruleId;
  let field = pathToString(violation.field);
  const rule = stableRuleIDs[ruleID] ?? ruleID;
  if (!field) {
    field = messageFieldHint(ruleID);
  }
  validationError(
    provider,
    operation,
    field,
    rule,
    violation.message || "invalid configuration",
  );
}

function messageFieldHint(ruleID: string): string {
  switch (ruleID) {
    case "cpu_limit":
    case "cpu_requires_positive":
      return "resources.cpu_limit_cores";
    case "memory_limit":
    case "memory_requires_positive":
      return "resources.memory_limit_mib";
    case "restricted_requires_allowlist":
      return "network.outbound_cidrs,outbound_domains";
    case "block_all_forbids_network":
      return "network.outbound_cidrs,outbound_domains,private_network,linked_sandbox,inbound_cidrs,ports";
    case "allow_all_forbids_allowlists":
      return "network.outbound_cidrs,outbound_domains";
    case "after_requires_duration":
    case "after_forbids_duration":
      return "after";
    case "oneof.required":
      return "tcp_port,command";
    case "auth_credentials_required":
      return "api_key,token_pair,bearer_token,oauth_refresh";
    case "exactly_one":
      return "client_secret,jwt_key";
    default:
      return "";
  }
}

export function validateConfig(
  value: Config,
  provider = "",
  operation = "create",
): void {
  const message = materializeClientConfig(value);
  const result = contractValidator.validate(
    ConfigSchema,
    message as MessageShape<typeof ConfigSchema>,
  );
  if (result.kind === "error") {
    throw sandboxError(
      {
        kind: ErrorKind.InvalidArgument,
        provider,
        operation,
        field: "",
        message: "sandbox-kit: invalid configuration",
      },
      result.error,
    );
  }
  if (result.kind === "invalid") {
    const violation = result.violations[0]!;
    const ruleID = violation.ruleId;
    let field = pathToString(violation.field);
    const rule = stableRuleIDs[ruleID] ?? ruleID;
    if (!field) {
      field = messageFieldHint(ruleID);
    }
    validationError(
      provider,
      operation,
      field,
      rule,
      violation.message || "invalid configuration",
    );
  }
  if (
    message.endpoint !== undefined &&
    !validSecureEndpoint(message.endpoint)
  ) {
    validationError(
      provider,
      operation,
      "endpoint",
      "secure_endpoint",
      "must use HTTPS except for loopback HTTP development endpoints",
    );
  }
}

/** Materializes a plain client config into the generated protobuf message. */
function materializeClientConfig(input: Config): ConfigWire {
  const { auth, timeout, ...rest } = input;
  return create(ConfigSchema, {
    ...rest,
    ...(timeout === undefined ? {} : { timeout: durationFromMs(timeout) }),
    ...(auth === undefined
      ? {}
      : { auth: { credentials: authCredentials(auth) } }),
  } as MessageInitShape<typeof ConfigSchema>) as ConfigWire;
}

/** Materializes plain create options into the generated protobuf message. */
function materializeCreateOptions(input: CreateOptions): CreateOptionsWire {
  const { source, readiness, providerOptions, network, ...rest } = input;
  const withDurations = materializeDurationFields(rest);
  return create(CreateOptionsSchema, {
    ...withDurations,
    ...(network === undefined
      ? {}
      : { network: materializeNetworkConfig(network) }),
    ...(source === undefined ? {} : { source: { kind: sourceKind(source) } }),
    ...(readiness === undefined
      ? {}
      : { readiness: { check: readinessCheck(readiness) } }),
    ...(providerOptions === undefined
      ? {}
      : { providerOptions: metadataObjectFromMap(providerOptions) }),
  } as unknown as MessageInitShape<
    typeof CreateOptionsSchema
  >) as CreateOptionsWire;
}

/** Wrap public string-array allowlists into wire StringAllowlist messages. */
function materializeNetworkConfig(network: NetworkConfigInput) {
  const { outboundCidrs, outboundDomains, inboundCidrs, ...rest } = network;
  return {
    ...rest,
    ...(outboundCidrs === undefined
      ? {}
      : { outboundCidrs: { entries: outboundCidrs } }),
    ...(outboundDomains === undefined
      ? {}
      : { outboundDomains: { entries: outboundDomains } }),
    ...(inboundCidrs === undefined
      ? {}
      : { inboundCidrs: { entries: inboundCidrs } }),
  };
}

function authCredentials(
  auth: AuthInput,
): NonNullable<AuthConfig["credentials"]> {
  if (auth.apiKey !== undefined)
    return { case: "apiKey", value: auth.apiKey as WireAPIKeyCredentials };
  if (auth.tokenPair !== undefined)
    return {
      case: "tokenPair",
      value: auth.tokenPair as WireTokenPairCredentials,
    };
  if (auth.bearerToken !== undefined)
    return {
      case: "bearerToken",
      value: auth.bearerToken as WireBearerTokenCredentials,
    };
  if (auth.oauthRefresh !== undefined)
    return {
      case: "oauthRefresh",
      value: auth.oauthRefresh as WireOAuthCredentials,
    };
  return { case: undefined };
}

function sourceKind(
  source: SandboxSourceInput,
): NonNullable<SandboxSource["kind"]> {
  if (source.providerDefault !== undefined)
    return {
      case: "providerDefault",
      value: {} as ProviderDefaultSource,
    };
  if (source.image !== undefined)
    return { case: "image", value: source.image as ImageSource };
  if (source.snapshot !== undefined)
    return { case: "snapshot", value: source.snapshot as SnapshotSource };
  if (source.warmPool !== undefined)
    return { case: "warmPool", value: source.warmPool as WarmPoolSource };
  return { case: undefined } as NonNullable<SandboxSource["kind"]>;
}

function readinessCheck(
  readiness: ReadinessProbeInput,
): NonNullable<ReadinessProbe["check"]> {
  if (readiness.tcpPort !== undefined)
    return { case: "tcpPort", value: readiness.tcpPort };
  if (readiness.command !== undefined)
    return { case: "command", value: readiness.command as CommandProbe };
  return { case: undefined } as NonNullable<ReadinessProbe["check"]>;
}

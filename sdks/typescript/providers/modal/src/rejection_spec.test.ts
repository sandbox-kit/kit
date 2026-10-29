import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import {
  ErrorKind,
  IsolationKind,
  PolicyMode,
  type Config,
  type CreateOptions,
} from "@sandbox-kit/sandbox";
import { describe, expect, it } from "vitest";
import { parse } from "yaml";

import { mapClientConfig, mapCreateRequest } from "#mappings/mappings.js";

const modalScope = { scope: { appName: "app" } };

describe("shared rejection spec", () => {
  const cases = loadCases().filter((item) => item.provider === "modal");

  it("covers the Modal cases", () => {
    expect(cases.length).toBeGreaterThan(0);
  });

  for (const item of cases) {
    it(item.name, () => {
      const run = () => {
        if (item.operation === "create")
          return mapCreateRequest(createRequest(item.create), modalScope);
        if (item.operation === "initialize")
          return mapClientConfig(clientConfig(item.config));
        throw new Error(`operation ${item.operation}`);
      };
      expect(run).toThrowError(
        expect.objectContaining({
          kind: kind(item.kind),
          field: item.field,
        }),
      );
    });
  }
});

interface SpecCase {
  name: string;
  provider: string;
  operation: string;
  kind: string;
  field: string;
  create?: Raw;
  config?: Raw;
}

type Raw = { [key: string]: unknown };

function loadCases(): SpecCase[] {
  const path = join(
    dirname(fileURLToPath(import.meta.url)),
    "../../../../../specs/rejections.yaml",
  );
  const parsed: unknown = parse(readFileSync(path, "utf8"));
  if (!isRecord(parsed)) throw new Error("rejection spec is not a document");
  only(parsed, ["version", "cases"], "spec");
  if (parsed.version !== 1) throw new Error("rejection spec version");
  if (!Array.isArray(parsed.cases)) throw new Error("rejection spec cases");
  return parsed.cases.map(specCase);
}

function specCase(value: unknown): SpecCase {
  if (!isRecord(value)) throw new Error("rejection case");
  only(
    value,
    ["name", "provider", "operation", "kind", "field", "create", "config"],
    "case",
  );
  if (typeof value.name !== "string" || typeof value.provider !== "string")
    throw new Error("rejection case identity");
  if (typeof value.operation !== "string" || typeof value.kind !== "string")
    throw new Error("rejection case result");
  if (typeof value.field !== "string") throw new Error("rejection case field");
  return {
    name: value.name,
    provider: value.provider,
    operation: value.operation,
    kind: value.kind,
    field: value.field,
    ...(isRecord(value.create) ? { create: value.create } : {}),
    ...(isRecord(value.config) ? { config: value.config } : {}),
  };
}

function createRequest(value: Raw | undefined): CreateOptions {
  if (value === undefined) throw new Error("create request is missing");
  only(
    value,
    [
      "source",
      "runtime",
      "isolation",
      "resources",
      "placement",
      "security",
      "network",
      "lifetime",
      "provisioning",
      "observability",
    ],
    "create",
  );
  const request: {
    -readonly [K in keyof CreateOptions]: CreateOptions[K];
  } = {};
  const source = record(value.source, "source");
  if (source !== undefined) {
    only(source, ["image"], "source");
    const image = record(source.image, "image");
    if (image !== undefined) {
      only(image, ["reference"], "image");
      request.source = { image: { reference: text(image.reference, "reference") } };
    }
  }
  const runtime = record(value.runtime, "runtime");
  if (runtime !== undefined) {
    only(runtime, ["language", "version"], "runtime");
    request.runtime = {
      ...(runtime.language === undefined
        ? {}
        : { language: text(runtime.language, "language") }),
      ...(runtime.version === undefined
        ? {}
        : { version: text(runtime.version, "version") }),
    };
  }
  const isolation = record(value.isolation, "isolation");
  if (isolation !== undefined) {
    only(isolation, ["kind", "nested_virtualization"], "isolation");
    request.isolation = {
      ...(isolation.kind === undefined
        ? {}
        : { kind: isolationKind(text(isolation.kind, "kind")) }),
      ...(isolation.nested_virtualization === undefined
        ? {}
        : { nestedVirtualization: flag(isolation.nested_virtualization) }),
    };
  }
  const resources = record(value.resources, "resources");
  if (resources !== undefined) {
    only(
      resources,
      ["cpu_limit_cores", "memory_limit_mib", "disk_mib"],
      "resources",
    );
    request.resources = {
      ...(resources.cpu_limit_cores === undefined
        ? {}
        : { cpuLimitCores: number(resources.cpu_limit_cores) }),
      ...(resources.memory_limit_mib === undefined
        ? {}
        : { memoryLimitMib: BigInt(number(resources.memory_limit_mib)) }),
      ...(resources.disk_mib === undefined
        ? {}
        : { diskMib: BigInt(number(resources.disk_mib)) }),
    };
  }
  const placement = record(value.placement, "placement");
  if (placement !== undefined) {
    only(placement, ["regions", "target"], "placement");
    request.placement = {
      ...(placement.regions === undefined
        ? {}
        : { regions: stringList(placement.regions) }),
      ...(placement.target === undefined
        ? {}
        : { target: text(placement.target, "target") }),
    };
  }
  const security = record(value.security, "security");
  if (security !== undefined) {
    only(security, ["public_access", "workload_identity"], "security");
    request.security = {
      ...(security.public_access === undefined
        ? {}
        : { publicAccess: flag(security.public_access) }),
      ...(security.workload_identity === undefined
        ? {}
        : { workloadIdentity: flag(security.workload_identity) }),
    };
  }
  const network = record(value.network, "network");
  if (network !== undefined) {
    only(network, ["private_network", "linked_sandbox"], "network");
    request.network = {
      ...(network.private_network === undefined
        ? {}
        : { privateNetwork: flag(network.private_network) }),
      ...(network.linked_sandbox === undefined
        ? {}
        : { linkedSandbox: text(network.linked_sandbox, "linked_sandbox") }),
    };
  }
  const lifetime = record(value.lifetime, "lifetime");
  if (lifetime !== undefined) {
    only(
      lifetime,
      ["maximum_lifetime_ms", "ephemeral", "idle_terminate"],
      "lifetime",
    );
    const idle = record(lifetime.idle_terminate, "idle_terminate");
    if (idle !== undefined) only(idle, ["mode"], "idle_terminate");
    request.lifetime = {
      ...(lifetime.maximum_lifetime_ms === undefined
        ? {}
        : { maximumLifetime: number(lifetime.maximum_lifetime_ms) }),
      ...(lifetime.ephemeral === undefined
        ? {}
        : { ephemeral: flag(lifetime.ephemeral) }),
      ...(idle === undefined
        ? {}
        : {
            idleTerminate: {
              mode:
                text(idle.mode, "mode") === "disabled"
                  ? PolicyMode.Disabled
                  : unexpected(text(idle.mode, "mode")),
            },
          }),
    };
  }
  const provisioning = record(value.provisioning, "provisioning");
  if (provisioning !== undefined) {
    only(provisioning, ["timeout_ms", "queue_timeout_ms"], "provisioning");
    request.provisioning = {
      ...(provisioning.timeout_ms === undefined
        ? {}
        : { timeout: number(provisioning.timeout_ms) }),
      ...(provisioning.queue_timeout_ms === undefined
        ? {}
        : { queueTimeout: number(provisioning.queue_timeout_ms) }),
    };
  }
  const observability = record(value.observability, "observability");
  if (observability !== undefined) {
    only(observability, ["verbose", "telemetry_endpoint"], "observability");
    request.observability = {
      ...(observability.verbose === undefined
        ? {}
        : { verbose: flag(observability.verbose) }),
      ...(observability.telemetry_endpoint === undefined
        ? {}
        : {
            telemetryEndpoint: text(
              observability.telemetry_endpoint,
              "telemetry_endpoint",
            ),
          }),
    };
  }
  return request;
}

function clientConfig(value: Raw | undefined): Config {
  if (value === undefined) throw new Error("config is missing");
  only(value, ["scope"], "config");
  const scope = record(value.scope, "scope");
  if (scope === undefined) throw new Error("config scope is missing");
  only(scope, ["organization_id", "project_id"], "scope");
  return {
    scope: {
      ...(scope.organization_id === undefined
        ? {}
        : { organizationId: text(scope.organization_id, "organization_id") }),
      ...(scope.project_id === undefined
        ? {}
        : { projectId: text(scope.project_id, "project_id") }),
    },
  };
}

function kind(value: string): (typeof ErrorKind)[keyof typeof ErrorKind] {
  if (value === "unsupported") return ErrorKind.Unsupported;
  if (value === "invalid_argument") return ErrorKind.InvalidArgument;
  return unexpected(value);
}

function isolationKind(value: string): (typeof IsolationKind)[keyof typeof IsolationKind] {
  if (value === "container") return IsolationKind.Container;
  return unexpected(value);
}

function unexpected(value: string): never {
  throw new Error(`unexpected spec value ${value}`);
}

function only(value: Raw, keys: readonly string[], label: string): void {
  for (const key of Object.keys(value)) {
    if (!keys.includes(key)) throw new Error(`${label} has unknown key ${key}`);
  }
}

function record(value: unknown, label: string): Raw | undefined {
  if (value === undefined) return undefined;
  if (!isRecord(value)) throw new Error(`${label} is not a record`);
  return value;
}

function isRecord(value: unknown): value is Raw {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function text(value: unknown, label: string): string {
  if (typeof value !== "string" || value === "")
    throw new Error(`${label} must be text`);
  return value;
}

function number(value: unknown): number {
  if (typeof value !== "number" || !Number.isFinite(value))
    throw new Error("spec number");
  return value;
}

function flag(value: unknown): boolean {
  if (typeof value !== "boolean") throw new Error("spec flag");
  return value;
}

function stringList(value: unknown): string[] {
  if (!Array.isArray(value) || value.some((item) => typeof item !== "string"))
    throw new Error("spec string list");
  return value;
}

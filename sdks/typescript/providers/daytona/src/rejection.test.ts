import {
  CreateFieldIsolation,
  CreateFieldLifetimeIdleTerminate,
  CreateFieldLifetimeMaximumLifetime,
  CreateFieldNetwork,
  CreateFieldObservabilityVerbose,
  CreateFieldPlacement,
  CreateFieldProvisioningQueueTimeout,
  CreateFieldProvisioningTimeout,
  CreateFieldResources,
  CreateFieldRuntime,
  CreateFieldSecurityWorkloadIdentity,
  ErrorKind,
  IsolationKind,
  PolicyMode,
  type CreateOptions,
} from "@sandbox-kit/sandbox";
import { describe, expect, it } from "vitest";

import { mapCreateRequest } from "#mappings/mappings.js";

describe("rejected Daytona create settings", () => {
  const cases: {
    name: string;
    request: CreateOptions;
    kind: (typeof ErrorKind)[keyof typeof ErrorKind];
    field: string;
  }[] = [
    {
      name: "placement",
      request: { placement: { regions: ["eu"] } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldPlacement,
    },
    {
      name: "idle terminate",
      request: {
        lifetime: { idleTerminate: { mode: PolicyMode.Disabled } },
      },
      kind: ErrorKind.Unsupported,
      field: CreateFieldLifetimeIdleTerminate,
    },
    {
      name: "verbose",
      request: { observability: { verbose: true } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldObservabilityVerbose,
    },
    {
      name: "workload identity",
      request: { security: { workloadIdentity: true } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldSecurityWorkloadIdentity,
    },
    {
      name: "runtime version",
      request: { runtime: { version: "3.12" } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldRuntime,
    },
    {
      name: "isolation class",
      request: { isolation: { kind: IsolationKind.Container } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldIsolation,
    },
    {
      name: "private network",
      request: { network: { privateNetwork: true } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldNetwork,
    },
    {
      name: "maximum lifetime",
      request: { lifetime: { maximumLifetime: 30_000 } },
      kind: ErrorKind.InvalidArgument,
      field: CreateFieldLifetimeMaximumLifetime,
    },
    {
      name: "fractional create timeout",
      request: { provisioning: { timeout: 1_500 } },
      kind: ErrorKind.InvalidArgument,
      field: CreateFieldProvisioningTimeout,
    },
    {
      name: "cpu limit",
      request: {
        source: { image: { reference: "alpine:3.21" } },
        resources: { cpuLimitCores: 2 },
      },
      kind: ErrorKind.Unsupported,
      field: CreateFieldResources,
    },
    {
      name: "memory limit",
      request: {
        source: { image: { reference: "alpine:3.21" } },
        resources: { memoryLimitMib: 2048n },
      },
      kind: ErrorKind.Unsupported,
      field: CreateFieldResources,
    },
    {
      name: "queue timeout",
      request: { provisioning: { queueTimeout: 30_000 } },
      kind: ErrorKind.InvalidArgument,
      field: CreateFieldProvisioningQueueTimeout,
    },
  ];

  for (const item of cases) {
    it(item.name, () => {
      expect(() => mapCreateRequest(item.request)).toThrowError(
        expect.objectContaining({ kind: item.kind, field: item.field }),
      );
    });
  }
});

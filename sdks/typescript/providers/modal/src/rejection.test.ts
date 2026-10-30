import {
  ConfigFieldScopeOrganizationId,
  ConfigFieldScopeProjectId,
  CreateFieldIsolation,
  CreateFieldLifetime,
  CreateFieldNetwork,
  CreateFieldObservabilityTelemetryEndpoint,
  CreateFieldPlacementTarget,
  CreateFieldProvisioningQueueTimeout,
  CreateFieldResourcesDiskMib,
  CreateFieldRuntime,
  CreateFieldSecurityPublicAccess,
  ErrorKind,
  type Config,
  type CreateOptions,
} from "@sandbox-kit/sandbox";
import { describe, expect, it } from "vitest";

import { mapClientConfig, mapCreateRequest } from "#mappings/mappings.js";

const app = { scope: { appName: "app" } };

describe("rejected Modal create settings", () => {
  const cases: {
    name: string;
    request: CreateOptions;
    kind: (typeof ErrorKind)[keyof typeof ErrorKind];
    field: string;
  }[] = [
    {
      name: "disk",
      request: imaged({ resources: { diskMib: 1024n } }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldResourcesDiskMib,
    },
    {
      name: "placement target",
      request: imaged({ placement: { target: "us" } }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldPlacementTarget,
    },
    {
      name: "queue timeout",
      request: imaged({ provisioning: { queueTimeout: 60_000 } }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldProvisioningQueueTimeout,
    },
    {
      name: "telemetry endpoint",
      request: imaged({
        observability: { telemetryEndpoint: "https://otel.example" },
      }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldObservabilityTelemetryEndpoint,
    },
    {
      name: "public access",
      request: imaged({ security: { publicAccess: true } }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldSecurityPublicAccess,
    },
    {
      name: "linked sandbox",
      request: imaged({ network: { linkedSandbox: "other" } }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldNetwork,
    },
    {
      name: "ephemeral lifetime",
      request: imaged({ lifetime: { ephemeral: true } }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldLifetime,
    },
    {
      name: "runtime language",
      request: imaged({ runtime: { language: "python" } }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldRuntime,
    },
    {
      name: "nested virtualization",
      request: imaged({ isolation: { nestedVirtualization: true } }),
      kind: ErrorKind.Unsupported,
      field: CreateFieldIsolation,
    },
  ];

  for (const item of cases) {
    it(item.name, () => {
      expect(() => mapCreateRequest(item.request, app)).toThrowError(
        expect.objectContaining({ kind: item.kind, field: item.field }),
      );
    });
  }
});

describe("rejected Modal client settings", () => {
  const cases: { name: string; config: Config; field: string }[] = [
    {
      name: "organization",
      config: { scope: { organizationId: "org" } },
      field: ConfigFieldScopeOrganizationId,
    },
    {
      name: "project",
      config: { scope: { projectId: "project" } },
      field: ConfigFieldScopeProjectId,
    },
  ];

  for (const item of cases) {
    it(item.name, () => {
      expect(() => mapClientConfig(item.config)).toThrowError(
        expect.objectContaining({
          kind: ErrorKind.Unsupported,
          field: item.field,
        }),
      );
    });
  }
});

function imaged(request: CreateOptions): CreateOptions {
  return {
    ...request,
    source: { image: { reference: "alpine:3.21" } },
  };
}

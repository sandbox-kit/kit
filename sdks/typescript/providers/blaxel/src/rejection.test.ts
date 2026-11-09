import {
  ConfigFieldAuth,
  ConfigFieldEndpoint,
  ConfigFieldScopeAppName,
  ConfigFieldScopeOrganizationId,
  CreateFieldIsolation,
  CreateFieldLifetime,
  CreateFieldName,
  CreateFieldNetwork,
  CreateFieldObservability,
  CreateFieldPlacement,
  CreateFieldProvisioningQueueTimeout,
  CreateFieldReadiness,
  CreateFieldResources,
  CreateFieldRuntime,
  CreateFieldSecurity,
  CreateFieldSource,
  ErrorKind,
  IsolationKind,
  type Config,
  type CreateOptions,
} from "@sandbox-kit/sandbox";
import { describe, expect, it } from "vitest";

import { mapClientConfig, mapCreateRequest } from "#mappings/mappings.js";

const workspace = { scope: { workspace: "ws" } };

describe("rejected Blaxel create settings", () => {
  const cases: {
    name: string;
    request: CreateOptions;
    kind: (typeof ErrorKind)[keyof typeof ErrorKind];
    field: string;
  }[] = [
    {
      name: "cpu resources",
      request: { resources: { cpuCores: 1 } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldResources,
    },
    {
      name: "disk",
      request: { resources: { diskMib: 1024n } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldResources,
    },
    {
      name: "placement",
      request: { placement: { regions: ["us-pdx-1"] } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldPlacement,
    },
    {
      name: "queue timeout",
      request: { provisioning: { queueTimeout: 60_000 } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldProvisioningQueueTimeout,
    },
    {
      name: "readiness",
      request: { readiness: { tcpPort: 8080 } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldReadiness,
    },
    {
      name: "observability",
      request: {
        observability: { telemetryEndpoint: "https://otel.example" },
      },
      kind: ErrorKind.Unsupported,
      field: CreateFieldObservability,
    },
    {
      name: "security",
      request: { security: { publicAccess: true } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldSecurity,
    },
    {
      name: "private network",
      request: { network: { privateNetwork: true } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldNetwork,
    },
    {
      name: "ephemeral lifetime",
      request: { lifetime: { ephemeral: true } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldLifetime,
    },
    {
      name: "runtime",
      request: { runtime: { language: "python" } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldRuntime,
    },
    {
      name: "isolation",
      request: { isolation: { kind: IsolationKind.Container } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldIsolation,
    },
  ];

  it("missing name", () => {
    expect(() =>
      mapCreateRequest(
        { source: { image: { reference: "blaxel/base-image:latest" } } },
        workspace,
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldName,
      }),
    );
  });

  it("missing image", () => {
    expect(() => mapCreateRequest({ name: "box" }, workspace)).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.Unsupported,
        field: CreateFieldSource,
      }),
    );
  });

  for (const item of cases) {
    it(item.name, () => {
      expect(() =>
        mapCreateRequest(namedImage(item.request), workspace),
      ).toThrowError(
        expect.objectContaining({ kind: item.kind, field: item.field }),
      );
    });
  }
});

describe("rejected Blaxel client settings", () => {
  const cases: { name: string; config: Config; field: string }[] = [
    {
      name: "endpoint",
      config: { endpoint: "https://unsupported.example" },
      field: ConfigFieldEndpoint,
    },
    {
      name: "organization",
      config: { scope: { organizationId: "org" } },
      field: ConfigFieldScopeOrganizationId,
    },
    {
      name: "app name",
      config: { scope: { appName: "app" } },
      field: ConfigFieldScopeAppName,
    },
    {
      name: "token pair auth",
      config: {
        auth: { tokenPair: { id: "id", secret: "secret" } },
      },
      field: ConfigFieldAuth,
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

function namedImage(request: CreateOptions): CreateOptions {
  return {
    name: "sandbox",
    ...request,
    source: { image: { reference: "blaxel/base-image:latest" } },
  };
}

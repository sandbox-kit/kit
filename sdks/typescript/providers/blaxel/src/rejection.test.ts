import {
  ConfigFieldAuth,
  ConfigFieldEndpoint,
  ConfigFieldRegion,
  ConfigFieldScopeAppName,
  ConfigFieldScopeOrganizationId,
  CreateFieldIsolation,
  CreateFieldLifetime,
  CreateFieldName,
  CreateFieldNetwork,
  CreateFieldObservability,
  CreateFieldPlacement,
  CreateFieldProvisioningQueueTimeout,
  CreateFieldProvisioningWaitFor,
  CreateFieldReadiness,
  CreateFieldResourcesCpuCores,
  CreateFieldResourcesDiskMib,
  CreateFieldRuntime,
  CreateFieldSecurity,
  CreateFieldSourceImage,
  ErrorKind,
  IsolationKind,
  WaitCondition,
  type Config,
  type CreateOptions,
} from "@sandbox-kit/sandbox";
import { describe, expect, it } from "vitest";

import { mapClientConfig, mapCreateRequest } from "#mappings/mappings.js";

const provider = { scope: { workspace: "ws" }, region: "us-pdx-1" };

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
      field: CreateFieldResourcesCpuCores,
    },
    {
      name: "disk",
      request: { resources: { diskMib: 1024n } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldResourcesDiskMib,
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
      name: "wait for ready",
      request: { provisioning: { waitFor: WaitCondition.Ready } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldProvisioningWaitFor,
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
      name: "runtime language",
      request: { runtime: { language: "python" } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldRuntime,
    },
    {
      name: "isolation container",
      request: { isolation: { kind: IsolationKind.Container } },
      kind: ErrorKind.Unsupported,
      field: CreateFieldIsolation,
    },
  ];

  it("missing name", () => {
    expect(() =>
      mapCreateRequest(
        { source: { image: { reference: "blaxel/base-image:latest" } } },
        provider,
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldName,
      }),
    );
  });

  it("missing image", () => {
    expect(() => mapCreateRequest({ name: "box" }, provider)).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldSourceImage,
      }),
    );
  });

  it("missing region", () => {
    expect(() =>
      mapCreateRequest(namedImage({}), { scope: { workspace: "ws" } }),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: ConfigFieldRegion,
      }),
    );
  });

  it("accepts kit create timeout as deadline", () => {
    expect(() =>
      mapCreateRequest(
        namedImage({ provisioning: { timeout: 30_000 } }),
        provider,
      ),
    ).not.toThrow();
  });

  it("accepts default isolation", () => {
    expect(() =>
      mapCreateRequest(
        namedImage({ isolation: { kind: IsolationKind.Default } }),
        provider,
      ),
    ).not.toThrow();
  });

  for (const item of cases) {
    it(item.name, () => {
      expect(() =>
        mapCreateRequest(namedImage(item.request), provider),
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

import { DaytonaError, type DaytonaConfig, type Sandbox as DaytonaSandbox } from "@daytona/sdk";
import {
  Client,
  EgressMode,
  ErrorKind,
  IsolationKind,
  PolicyMode,
  SecretInjection,
  WaitCondition,
  type CreateOptions,
} from "@sandbox-kit/sandbox";
import { describe, expect, it } from "vitest";

import { DaytonaProvider } from "./daytona.js";
import {
  mapClientConfig,
  mapCreateRequest,
  mapDaytonaError,
  type DaytonaCreateParams,
} from "./mappings.js";

describe("Daytona mappings", () => {
  it("maps API-key client configuration without exposing the provider", () => {
    expect(
      mapClientConfig({
        endpoint: "https://daytona.example/api",
        region: "us",
        auth: { credentials: "apiKey", apiKey: { key: "secret" } },
      }),
    ).toEqual({
      apiKey: "secret",
      apiUrl: "https://daytona.example/api",
      target: "us",
    });
  });

  it("maps an image request with exact resource and policy units", () => {
    const request: CreateOptions = {
      source: { kind: "image", image: { reference: "node:24" } },
      runtime: { language: "typescript" },
      resources: { cpuCores: 2, memoryMib: 2_048n, diskMib: 10_240n },
      lifetime: {
        idleStop: { mode: PolicyMode.After, after: 120_000 },
      },
    };

    expect(mapCreateRequest(request)).toEqual({
      image: "node:24",
      language: "typescript",
      resources: { cpu: 2, memory: 2, disk: 10 },
      autoStopInterval: 2,
    });
  });

  it("maps supported isolation, security, network, and lifetime settings", () => {
    expect(
      mapCreateRequest({
        isolation: {
          kind: IsolationKind.Default,
          nestedVirtualization: true,
        },
        security: {
          publicAccess: false,
          secrets: [
            {
              reference: "github-token",
              environmentVariable: "GITHUB_TOKEN",
              injection: SecretInjection.EgressPlaceholder,
            },
          ],
        },
        network: {
          egress: EgressMode.Restricted,
          outboundCidrs: { entries: ["192.0.2.0/24"] },
          outboundDomains: { entries: ["example.com"] },
          linkedSandbox: "sandbox-parent",
          outboundProxyUrl: "https://proxy.example.com",
        },
        lifetime: {
          maximumLifetime: 3_600_000,
          ephemeral: false,
          idleStop: { mode: PolicyMode.After, after: 900_000 },
          stoppedArchive: { mode: PolicyMode.After, after: 86_400_000 },
          stoppedDelete: { mode: PolicyMode.After, after: 172_800_000 },
        },
        provisioning: {
          queueTimeout: 300_000,
          waitFor: WaitCondition.Started,
        },
        observability: {
          telemetryEndpoint: "https://otel.example.com",
        },
      }),
    ).toEqual({
      kvm: true,
      public: false,
      secrets: { GITHUB_TOKEN: "github-token" },
      networkAllowList: "192.0.2.0/24",
      domainAllowList: "example.com",
      linkedSandbox: "sandbox-parent",
      outboundProxyUrl: "https://proxy.example.com",
      ttlMinutes: 60,
      ephemeral: false,
      autoStopInterval: 15,
      autoArchiveInterval: 1_440,
      autoDeleteInterval: 2_880,
      queueTimeout: 5,
      otelEndpointOverride: "https://otel.example.com",
    });
  });

  it("maps bearer authentication with its required organization scope", () => {
    expect(
      mapClientConfig({
        auth: { credentials: "bearerToken", bearerToken: { token: "jwt" } },
        scope: { organizationId: "organization-1" },
      }),
    ).toMatchObject({ jwtToken: "jwt", organizationId: "organization-1" });
  });

  it("rejects wait intent the Daytona TypeScript SDK cannot represent", () => {
    expect(() =>
      mapCreateRequest({ provisioning: { waitFor: WaitCondition.Submitted } }),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.Unsupported,
        field: "provisioning.waitFor",
      }),
    );
  });

  it("preserves Daytona error details and cause", () => {
    const native = new DaytonaError("limited", 429, undefined, "RATE", "DAYTONA_API");
    expect(mapDaytonaError(native, "create")).toMatchObject({
      kind: ErrorKind.RateLimited,
      statusCode: 429,
      providerCode: "RATE",
      providerSource: "DAYTONA_API",
      cause: native,
    });
  });
});

describe("Daytona provider orchestration", () => {
  it("creates through the native SDK and only disposes local client resources", async () => {
    let nativeConfig: DaytonaConfig | undefined;
    let nativeRequest: DaytonaCreateParams | undefined;
    let closeCalls = 0;
    const nativeSandbox = {
      id: "sandbox-1",
      name: "example",
      target: "us",
      labels: {},
      cpu: 2,
      memory: 4,
      disk: 10,
      state: "started",
      toolboxProxyUrl: "",
    } as DaytonaSandbox;
    const native = {
      async create(params: DaytonaCreateParams): Promise<DaytonaSandbox> {
        nativeRequest = params;
        return nativeSandbox;
      },
      async get(): Promise<DaytonaSandbox> { return nativeSandbox; },
      async [Symbol.asyncDispose](): Promise<void> {
        closeCalls++;
      },
    };
    const provider = new DaytonaProvider((config) => {
      nativeConfig = config;
      return native;
    });
    const client = new Client({ provider, auth: { credentials: "apiKey", apiKey: { key: "secret" } } });

    const sandbox = await client.create({ source: { kind: "image", image: { reference: "node:24" } } });
    await client.close();

    expect(nativeConfig?.apiKey).toBe("secret");
    expect(nativeRequest).toMatchObject({ image: "node:24" });
    expect(sandbox.id).toBe("sandbox-1");
    expect(closeCalls).toBe(1);
  });
});

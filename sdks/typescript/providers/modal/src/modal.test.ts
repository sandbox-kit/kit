import { describe, expect, it } from "vitest";
import {
  Client,
  EgressMode,
  ErrorKind,
  IsolationKind,
  PolicyMode,
  PortTransport,
  SecretInjection,
  WaitCondition,
  CreateFieldEnvironment,
  CreateFieldLifetimeIdleTerminateAfter,
  CreateFieldLifetimeMaximumLifetime,
  CreateFieldNetworkInboundCidrsEntries,
  CreateFieldResourcesCpuCores,
  CreateFieldResourcesCpuLimitCores,
  ConfigFieldAuthTokenPairId,
  ConfigFieldEndpoint,
} from "@sandbox-kit/sandbox";
import {
  App,
  InvalidError,
  Proxy,
  type Image,
  type Sandbox as ModalSandbox,
  type SandboxCreateParams,
} from "modal";
import { mapClientConfig, mapCreateRequest } from "#mappings/mappings.js";
import { ModalProvider } from "./modal.js";
import {
  mapIdleDuration,
  mapLifetimeDuration,
  mapResources,
} from "#mappings/provider_mappings.js";
import { captureProviderState } from "#config/provider_config.js";
import { mapModalError } from "#errors.js";

describe("Modal mappings", () => {
  it("maps token authentication and environment scope", () => {
    expect(
      mapClientConfig({
        auth: {
          tokenPair: { id: "id", secret: "secret" },
        },
        scope: { environment: "dev" },
        timeout: 30_000,
      }),
    ).toEqual({
      tokenId: "id",
      tokenSecret: "secret",
      environment: "dev",
      timeoutMs: 30_000,
    });
  });

  it("rejects unsupported client endpoint before SDK construction", () => {
    expect(() =>
      mapClientConfig({ endpoint: "https://unsupported.example" }),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.Unsupported,
        field: ConfigFieldEndpoint,
        operation: "initialize",
      }),
    );
  });

  it("reports canonical fields for invalid durations", () => {
    const config = { scope: { appName: "app" } };
    const image = { source: { image: { reference: "alpine:3.21" } } };
    expect(() =>
      mapCreateRequest({ ...image, lifetime: { maximumLifetime: 0 } }, config),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldLifetimeMaximumLifetime,
      }),
    );
    expect(() =>
      mapCreateRequest(
        {
          ...image,
          lifetime: {
            idleTerminate: { mode: PolicyMode.After, after: 0 },
          },
        },
        config,
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldLifetimeIdleTerminateAfter,
      }),
    );
  });

  it("rejects blank auth credentials in Client construction", () => {
    expect(
      () =>
        new Client({
          provider: new ModalProvider(() => {
            throw new Error("native SDK must not construct");
          }),
          auth: {
            tokenPair: { id: "", secret: "secret" },
          },
        }),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: ConfigFieldAuthTokenPairId,
        operation: "initialize",
      }),
    );
  });

  it("captures retained configuration without aliasing caller state", () => {
    const config = {
      scope: { appName: "original", environment: "dev" },
      region: "us-east-1",
    };
    const state = captureProviderState(config);
    config.scope.appName = "changed";

    expect(state).toMatchObject({
      scope: { appName: "original", environment: "dev" },
      region: "us-east-1",
    });
  });

  it("maps the Go creation surface into the native Modal plan", () => {
    const plan = mapCreateRequest(
      {
        name: "example",
        source: {
          image: { reference: "alpine:3.21" },
        },
        runtime: {
          workingDirectory: "/tmp",
          entrypoint: ["sleep", "infinity"],
          pty: true,
        },
        isolation: { kind: IsolationKind.LinuxVm },
        resources: {
          cpuCores: 0.5,
          cpuLimitCores: 1,
          memoryMib: 512n,
          memoryLimitMib: 1024n,
        },
        placement: { cloud: "aws", regions: ["us-east-1"] },
        network: {
          egress: EgressMode.Restricted,
          outboundCidrs: ["192.0.2.0/24"],
          outboundDomains: ["example.com"],
          inboundCidrs: ["198.51.100.0/24"],
          privateNetwork: true,
          proxyReference: "static-egress",
          ports: [{ port: 443, transport: PortTransport.Tls }],
        },
        security: {
          workloadIdentity: true,
          secrets: [
            { reference: "credential", injection: SecretInjection.Environment },
          ],
        },
        lifetime: {
          maximumLifetime: 300_000,
          idleTerminate: {
            mode: PolicyMode.After,
            after: 60_000,
          },
        },
        readiness: { command: { argv: ["true"] } },
        provisioning: {
          timeout: 120_000,
          waitFor: WaitCondition.Ready,
        },
        observability: { verbose: true },
        labels: { purpose: "test" },
        environment: { NODE_ENV: "test" },
      },
      { scope: { appName: "app", environment: "dev" }, region: "us-east-1" },
    );
    expect(plan).toMatchObject({
      app: "app",
      environment: "dev",
      image: "alpine:3.21",
      secretNames: ["credential"],
      proxyName: "static-egress",
      waitUntilReady: true,
      readyTimeoutMs: 120_000,
      params: {
        name: "example",
        runtime: "vm",
        cpu: 0.5,
        cpuLimit: 1,
        memoryMiB: 512,
        memoryLimitMiB: 1024,
        cloud: "aws",
        regions: ["us-east-1"],
        i6pn: true,
        encryptedPorts: [443],
        includeOidcIdentityToken: true,
        timeoutMs: 300_000,
        idleTimeoutMs: 60_000,
        verbose: true,
      },
    });
  });

  it("rejects CPU values that are not millis-representable", () => {
    expect(mapResources({ cpuCores: 0.5 })).toMatchObject({ cpu: 0.5 });
    expect(() => mapResources({ cpuCores: 1.0005 })).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldResourcesCpuCores,
      }),
    );
    expect(() => mapResources({ cpuLimitCores: 0.0001 })).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldResourcesCpuLimitCores,
      }),
    );
  });

  it("rejects lifetime and idle durations that are not whole seconds", () => {
    expect(mapLifetimeDuration({ maximumLifetime: 60_000 })).toMatchObject({
      timeoutMs: 60_000,
    });
    expect(() => mapLifetimeDuration({ maximumLifetime: 1_500 })).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldLifetimeMaximumLifetime,
      }),
    );
    expect(() => mapIdleDuration({ after: 500 })).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldLifetimeIdleTerminateAfter,
      }),
    );
  });

  it("classifies native errors from the provider specification", () => {
    const native = new InvalidError("invalid request");
    expect(mapModalError(native, "create")).toMatchObject({
      kind: ErrorKind.InvalidArgument,
      provider: "modal",
      operation: "create",
      cause: native,
    });
  });

  it("applies generated environment and network checks before mapping", () => {
    const config = { scope: { appName: "app" } };
    expect(() =>
      mapCreateRequest(
        {
          source: {
            image: { reference: "alpine:3.21" },
          },
          environment: { "NOT VALID": "value" },
        },
        config,
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldEnvironment,
      }),
    );
    expect(() =>
      mapCreateRequest(
        {
          source: {
            image: { reference: "alpine:3.21" },
          },
          network: { inboundCidrs: ["300.1.1.1/24"] },
        },
        config,
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldNetworkInboundCidrsEntries,
      }),
    );
  });
});

describe("Modal provider orchestration", () => {
  it("resolves a managed proxy before native sandbox creation", async () => {
    let proxyLookup:
      | { readonly name: string; readonly environment?: string }
      | undefined;
    let nativeCreate: SandboxCreateParams | undefined;
    let closeCalls = 0;
    const nativeSandbox = { sandboxId: "sandbox-1" } as ModalSandbox;
    const native = {
      apps: {
        async fromName(name: string): Promise<App> {
          return new App("app-id", name);
        },
      },
      images: {
        fromRegistry(): Image {
          return {} as Image;
        },
      },
      proxies: {
        async fromName(
          name: string,
          params?: { readonly environment?: string },
        ): Promise<Proxy> {
          proxyLookup = {
            name,
            ...(params?.environment === undefined
              ? {}
              : { environment: params.environment }),
          };
          return new Proxy("proxy-id");
        },
      },
      secrets: {
        async fromName(): Promise<never> {
          throw new Error("unexpected secret lookup");
        },
      },
      sandboxes: {
        async create(
          _app: App,
          _image: Image,
          params?: SandboxCreateParams,
        ): Promise<ModalSandbox> {
          nativeCreate = params;
          return nativeSandbox;
        },
        async fromId(): Promise<ModalSandbox> {
          return nativeSandbox;
        },
      },
      close(): void {
        closeCalls++;
      },
    };
    const client = new Client({
      provider: new ModalProvider(() => native),
      scope: { appName: "app", environment: "dev" },
    });

    const sandbox = await client.create({
      source: { image: { reference: "alpine:3.21" } },
      network: { proxyReference: "static-egress" },
    });
    await client.close();

    expect(proxyLookup).toEqual({ name: "static-egress", environment: "dev" });
    expect(nativeCreate?.proxy?.proxyId).toBe("proxy-id");
    expect(sandbox.id).toBe("sandbox-1");
    expect(closeCalls).toBe(1);
  });

  it("does not create a sandbox after the signal aborts during lookup", async () => {
    const controller = new AbortController();
    let created = false;
    const native = {
      apps: {
        async fromName(name: string): Promise<App> {
          controller.abort();
          return new App("app-id", name);
        },
      },
      images: {
        fromRegistry(): Image {
          return {} as Image;
        },
      },
      proxies: {
        async fromName(): Promise<Proxy> {
          throw new Error("unexpected proxy lookup");
        },
      },
      secrets: {
        async fromName(): Promise<never> {
          throw new Error("unexpected secret lookup");
        },
      },
      sandboxes: {
        async create(): Promise<ModalSandbox> {
          created = true;
          throw new Error("native create");
        },
        async fromId(): Promise<ModalSandbox> {
          throw new Error("unexpected attach");
        },
      },
      close(): void {},
    };
    const backend = new ModalProvider(() => native).initialize({
      scope: { appName: "app" },
    });
    if (backend.create === undefined) throw new Error("create is missing");
    await expect(
      backend.create(
        { source: { image: { reference: "alpine:3.21" } } },
        { signal: controller.signal },
      ),
    ).rejects.toMatchObject({
      kind: ErrorKind.Canceled,
      operation: "create",
    });
    expect(created).toBe(false);
  });

  it("stops waiting for termination when the signal aborts", async () => {
    const controller = new AbortController();
    let polls = 0;
    const nativeSandbox = {
      async terminate(): Promise<void> {},
      async poll(): Promise<number | null> {
        polls++;
        controller.abort();
        return null;
      },
    } as ModalSandbox;
    const backend = new ModalProvider(() => ({
      apps: {
        async fromName(): Promise<App> {
          throw new Error("unexpected app lookup");
        },
      },
      images: {
        fromRegistry(): Image {
          return {} as Image;
        },
      },
      proxies: {
        async fromName(): Promise<Proxy> {
          throw new Error("unexpected proxy lookup");
        },
      },
      secrets: {
        async fromName(): Promise<never> {
          throw new Error("unexpected secret lookup");
        },
      },
      sandboxes: {
        async create(): Promise<ModalSandbox> {
          throw new Error("unexpected create");
        },
        async fromId(): Promise<ModalSandbox> {
          return nativeSandbox;
        },
      },
      close(): void {},
    })).initialize({ scope: { appName: "app" } });
    if (backend.terminate === undefined) throw new Error("terminate is missing");
    await expect(
      backend.terminate("sandbox-1", { wait: true }, { signal: controller.signal }),
    ).rejects.toMatchObject({
      kind: ErrorKind.Canceled,
      operation: "terminate",
    });
    expect(polls).toBe(1);
  });
});

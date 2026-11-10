import { tmpdir } from "node:os";
import { join } from "node:path";
import { describe, expect, it } from "vitest";
import {
  Client,
  ErrorKind,
  PolicyMode,
  PortTransport,
  SandboxState,
  WaitCondition,
  ConfigFieldAuthApiKeyKey,
  ConfigFieldEndpoint,
  CreateFieldEnvironment,
  CreateFieldLifetimeIdleTerminate,
  CreateFieldLifetimeMaximumLifetime,
  CreateFieldName,
  CreateFieldNetworkPortsPort,
  CreateFieldProvisioningWaitFor,
  FileFieldLocalPath,
} from "@sandbox-kit/sandbox";
import { mapClientConfig, mapCreateRequest } from "#mappings/mappings.js";
import { BlaxelProvider } from "./blaxel.js";
import { mapResources } from "#mappings/provider_mappings.js";
import { captureProviderState } from "#config/provider_config.js";
import { mapBlaxelError } from "#errors.js";
import { ResponseError } from "@blaxel/core";

describe("Blaxel mappings", () => {
  it("maps API key authentication and workspace scope", () => {
    expect(
      mapClientConfig({
        auth: { apiKey: { key: "bl_key" } },
        scope: { workspace: "my-workspace" },
      }),
    ).toEqual({
      apiKey: "bl_key",
      workspace: "my-workspace",
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

  it("requires name and image on create", () => {
    const config = { scope: { workspace: "ws" } };
    expect(() =>
      mapCreateRequest(
        { source: { image: { reference: "blaxel/base-image:latest" } } },
        config,
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldName,
      }),
    );
  });

  it("rejects blank auth credentials in Client construction", () => {
    expect(
      () =>
        new Client({
          provider: new BlaxelProvider(() => {
            throw new Error("native SDK must not construct");
          }),
          auth: { apiKey: { key: "" } },
          scope: { workspace: "ws" },
        }),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: ConfigFieldAuthApiKeyKey,
        operation: "initialize",
      }),
    );
  });

  it("captures retained configuration without aliasing caller state", () => {
    const config = {
      scope: { workspace: "original" },
      region: "us-pdx-1",
    };
    const state = captureProviderState(config);
    config.scope.workspace = "changed";

    expect(state).toMatchObject({
      scope: { workspace: "original" },
      region: "us-pdx-1",
    });
  });

  it("maps the creation surface into the native Blaxel plan", () => {
    const plan = mapCreateRequest(
      {
        name: "example",
        source: { image: { reference: "blaxel/base-image:latest" } },
        resources: { memoryMib: 4096n },
        network: {
          outboundDomains: ["example.com"],
          ports: [{ port: 3000, transport: PortTransport.Http2Tls }],
        },
        lifetime: {
          maximumLifetime: 24 * 60 * 60_000,
          idleTerminate: {
            mode: PolicyMode.After,
            after: 30 * 60_000,
          },
        },
        provisioning: {
          timeout: 30_000,
        },
        labels: { purpose: "test" },
        environment: { NODE_ENV: "test" },
      },
      { scope: { workspace: "ws" }, region: "us-pdx-1" },
    );
    expect(plan).toMatchObject({
      options: { timeout: 30 },
      params: {
        name: "example",
        image: "blaxel/base-image:latest",
        memory: 4096,
        region: "us-pdx-1",
        ttl: "1d",
        labels: { purpose: "test" },
        envs: [{ name: "NODE_ENV", value: "test" }],
        ports: [{ target: 3000, protocol: "HTTP" }],
        network: { allowedDomains: ["example.com"] },
        lifecycle: {
          expirationPolicies: [
            { type: "ttl-idle", action: "delete", value: "30m" },
          ],
        },
      },
    });
  });

  it("maps memory without a Kit-side floor", () => {
    expect(
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          resources: { memoryMib: 128n },
        },
        { scope: { workspace: "ws" } },
      ).params.memory,
    ).toBe(128);
    expect(mapResources({ memoryMib: 4096n })).toEqual({ memory: 4096 });
  });

  it("rejects WaitCondition.Ready", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          provisioning: {
            timeout: 30_000,
            waitFor: WaitCondition.Ready,
          },
        },
        { scope: { workspace: "ws" } },
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.Unsupported,
        field: CreateFieldProvisioningWaitFor,
      }),
    );
  });

  it("rejects Blaxel-reserved ports", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          network: { ports: [{ port: 443, transport: PortTransport.Tls }] },
        },
        { scope: { workspace: "ws" } },
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldNetworkPortsPort,
      }),
    );
  });

  it("rejects lifetime durations that are not whole seconds", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          lifetime: { maximumLifetime: 1_500 },
        },
        { scope: { workspace: "ws" } },
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldLifetimeMaximumLifetime,
      }),
    );
  });

  it("rejects disabling idle terminate", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          lifetime: { idleTerminate: { mode: PolicyMode.Disabled } },
        },
        { scope: { workspace: "ws" } },
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.Unsupported,
        field: CreateFieldLifetimeIdleTerminate,
      }),
    );
  });

  it("classifies native ResponseError status codes", () => {
    const native = new ResponseError(
      new Response(null, { status: 404, statusText: "Not Found" }),
      { error: { code: "WORKLOAD_NOT_FOUND", message: "missing" } },
      new Error("missing"),
    );
    expect(mapBlaxelError(native, "refresh")).toMatchObject({
      kind: ErrorKind.NotFound,
      provider: "blaxel",
      operation: "refresh",
      cause: native,
    });
  });

  it("applies environment checks before mapping", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          environment: { "NOT VALID": "value" },
        },
        { scope: { workspace: "ws" } },
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldEnvironment,
      }),
    );
  });
});

describe("Blaxel provider orchestration", () => {
  it("creates with SandboxInstance.create and never createIfNotExists", async () => {
    let created:
      | { params: unknown; options: unknown }
      | undefined;
    let configured: unknown;
    const nativeSandbox = {
      metadata: { name: "example" },
      status: "RUNNING",
      spec: { region: "us-pdx-1" },
      fs: {} as never,
      process: {} as never,
      async wait() {
        return this;
      },
      async delete() {},
    };
    const client = new Client({
      provider: new BlaxelProvider(() => ({
        configure(config) {
          configured = config;
        },
        async create(params, options) {
          created = { params, options };
          return nativeSandbox;
        },
        async get() {
          throw new Error("unexpected get");
        },
        async delete() {
          throw new Error("unexpected delete");
        },
      })),
      auth: { apiKey: { key: "bl_key" } },
      scope: { workspace: "ws" },
      region: "us-pdx-1",
    });

    const sandbox = await client.create({
      name: "example",
      source: { image: { reference: "blaxel/base-image:latest" } },
      resources: { memoryMib: 2048n },
    });
    await client.close();

    expect(configured).toEqual({ apiKey: "bl_key", workspace: "ws" });
    expect(created?.params).toMatchObject({
      name: "example",
      image: "blaxel/base-image:latest",
      memory: 2048,
      region: "us-pdx-1",
    });
    expect(created?.options).toBeUndefined();
    expect(sandbox.id).toBe("example");
    expect(sandbox.info().name).toBe("example");
  });

  it("looks up sandboxes by name for get/refresh", async () => {
    let lookedUp: string | undefined;
    const backend = new BlaxelProvider(() => ({
      configure() {},
      async create() {
        throw new Error("unexpected create");
      },
      async get(name) {
        lookedUp = name;
        return {
          metadata: { name },
          status: "RUNNING",
          fs: {} as never,
          process: {} as never,
          async wait() {
            return this;
          },
          async delete() {},
        };
      },
      async delete() {
        throw new Error("unexpected delete");
      },
    })).initialize({
      auth: { apiKey: { key: "bl_key" } },
      scope: { workspace: "ws" },
    });
    if (backend.refresh === undefined) throw new Error("refresh is missing");
    const info = await backend.refresh("existing", {});
    expect(lookedUp).toBe("existing");
    expect(info).toMatchObject({ id: "existing", provider: "blaxel" });
  });

  it("deletes by name after a provider-verified lookup", async () => {
    const deleted: string[] = [];
    const backend = new BlaxelProvider(() => ({
      configure() {},
      async create() {
        throw new Error("unexpected create");
      },
      async get(name) {
        return {
          metadata: { name },
          status: "RUNNING",
          fs: {} as never,
          process: {} as never,
          async wait() {
            return this;
          },
          async delete() {
            deleted.push(name);
          },
        };
      },
      async delete(name) {
        deleted.push(`static:${name}`);
      },
    })).initialize({
      auth: { apiKey: { key: "bl_key" } },
      scope: { workspace: "ws" },
    });
    if (backend.delete === undefined) throw new Error("delete is missing");
    const info = await backend.delete("doomed", {}, {});
    expect(deleted).toEqual(["doomed"]);
    expect(info).toMatchObject({
      id: "doomed",
      state: SandboxState.Deleted,
      terminal: true,
      providerState: "TERMINATED",
    });
  });

  it("does not create a sandbox after the signal aborts during mapping", async () => {
    const controller = new AbortController();
    let created = false;
    const backend = new BlaxelProvider(() => ({
      configure() {},
      async create() {
        created = true;
        throw new Error("native create");
      },
      async get() {
        throw new Error("unexpected get");
      },
      async delete() {
        throw new Error("unexpected delete");
      },
    })).initialize({
      auth: { apiKey: { key: "bl_key" } },
      scope: { workspace: "ws" },
    });
    controller.abort();
    if (backend.create === undefined) throw new Error("create is missing");
    await expect(
      backend.create(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
        },
        { signal: controller.signal },
      ),
    ).rejects.toMatchObject({
      kind: ErrorKind.Canceled,
      operation: "create",
    });
    expect(created).toBe(false);
  });

  it("rejects a missing or directory local file before Blaxel is called", async () => {
    let lookedUp = false;
    const backend = new BlaxelProvider(() => ({
      configure() {},
      async create() {
        throw new Error("unexpected create");
      },
      async get() {
        lookedUp = true;
        throw new Error("native get must not run");
      },
      async delete() {
        throw new Error("unexpected delete");
      },
    })).initialize({
      auth: { apiKey: { key: "bl_key" } },
      scope: { workspace: "ws" },
    });
    const copyFromLocal = backend.copyFromLocal;
    const copyToLocal = backend.copyToLocal;
    if (!copyFromLocal || !copyToLocal)
      throw new Error("copy methods are missing");
    await expect(
      copyFromLocal.call(
        backend,
        "sandbox-1",
        join(tmpdir(), "kit-missing-source.txt"),
        "/tmp/note.txt",
        {},
        {},
      ),
    ).rejects.toMatchObject({
      kind: ErrorKind.NotFound,
      field: FileFieldLocalPath,
      operation: "copy_from_local",
    });
    await expect(
      copyFromLocal.call(
        backend,
        "sandbox-1",
        tmpdir(),
        "/tmp/note.txt",
        {},
        {},
      ),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: FileFieldLocalPath,
      operation: "copy_from_local",
    });
    await expect(
      copyToLocal.call(backend, "sandbox-1", "/tmp/note.txt", tmpdir(), {}, {}),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: FileFieldLocalPath,
      operation: "copy_to_local",
    });
    expect(lookedUp).toBe(false);
  });

  it("runs commands with waitForCompletion and second-based timeouts", async () => {
    const seen: Array<{
      command: string;
      waitForCompletion?: boolean;
      timeout?: number;
    }> = [];
    const backend = new BlaxelProvider(() => ({
      configure() {},
      async create() {
        throw new Error("unexpected create");
      },
      async get() {
        return {
          metadata: { name: "sandbox-1" },
          fs: {} as never,
          process: {
            async exec(request) {
              const entry: (typeof seen)[number] = {
                command: request.command,
              };
              if (request.waitForCompletion !== undefined)
                entry.waitForCompletion = request.waitForCompletion;
              if (request.timeout !== undefined) entry.timeout = request.timeout;
              seen.push(entry);
              return { exitCode: 0, stdout: "ok", stderr: "" };
            },
          },
          async wait() {
            return this;
          },
          async delete() {},
        };
      },
      async delete() {
        throw new Error("unexpected delete");
      },
    })).initialize({
      auth: { apiKey: { key: "bl_key" } },
      scope: { workspace: "ws" },
    });
    if (backend.run === undefined) throw new Error("run is missing");

    await backend.run("sandbox-1", { argv: ["echo", "hi"] }, {});
    await backend.run("sandbox-1", { argv: ["echo", "hi"], timeout: 0 }, {});
    await backend.run(
      "sandbox-1",
      { argv: ["echo", "hi"], timeout: 5_000 },
      {},
    );
    await expect(
      backend.run("sandbox-1", { argv: ["echo", "hi"], timeout: 1_500 }, {}),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: "timeout",
      operation: "run",
    });

    expect(seen).toEqual([
      { command: "'echo' 'hi'", waitForCompletion: true, timeout: 0 },
      { command: "'echo' 'hi'", waitForCompletion: true, timeout: 0 },
      { command: "'echo' 'hi'", waitForCompletion: true, timeout: 5 },
    ]);
  });
});

import { tmpdir } from "node:os";
import { join } from "node:path";
import { describe, expect, it, vi } from "vitest";
import {
  Client,
  EgressMode,
  ErrorKind,
  FileType,
  PolicyMode,
  PortTransport,
  SandboxState,
  WaitCondition,
  ConfigFieldAuthApiKeyKey,
  ConfigFieldEndpoint,
  ConfigFieldRegion,
  ConfigFieldScopeWorkspace,
  CreateFieldEnvironment,
  CreateFieldLifetimeIdleTerminate,
  CreateFieldLifetimeMaximumLifetime,
  CreateFieldName,
  CreateFieldNetworkPortsPort,
  CreateFieldNetworkPortsTransport,
  CreateFieldProvisioningWaitFor,
  FileFieldLocalPath,
} from "@sandbox-kit/sandbox";
import {
  mapClientConfig,
  mapCreateRequest,
  mapSandboxInfo,
} from "#mappings/mappings.js";
import { BlaxelProvider } from "./blaxel.js";
import {
  mapAllocatedResources,
  mapResources,
  mapSandboxIdentity,
} from "#mappings/provider_mappings.js";
import { captureProviderState } from "#config/provider_config.js";
import { mapBlaxelError } from "#errors.js";
import { projectHandle } from "./native.js";
import { ResponseError, SandboxInstance } from "@blaxel/core";

const provider = { region: "us-pdx-1" };

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

  it("rejects blank workspace scope", () => {
    expect(() => mapClientConfig({ scope: { workspace: "  " } })).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: ConfigFieldScopeWorkspace,
        operation: "initialize",
      }),
    );
  });

  it("requires name and image on create", () => {
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

  it("requires client region on create", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
        },
        {},
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: ConfigFieldRegion,
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

  it("captures create region from client config", () => {
    expect(
      captureProviderState({
        scope: { workspace: "original" },
        region: "us-pdx-1",
      }),
    ).toEqual({ region: "us-pdx-1" });
    expect(captureProviderState({ scope: { workspace: "ws" } })).toEqual({});
  });

  it("maps the creation surface into the native Blaxel plan", () => {
    const plan = mapCreateRequest(
      {
        name: "example",
        source: { image: { reference: "blaxel/base-image:latest" } },
        resources: { memoryMib: 4096n },
        network: {
          outboundDomains: ["example.com"],
          ports: [{ port: 3000, transport: PortTransport.Tls }],
        },
        lifetime: {
          maximumLifetime: 24 * 60 * 60_000,
          idleTerminate: {
            mode: PolicyMode.After,
            after: 30 * 60_000,
          },
        },
        labels: { purpose: "test" },
        environment: { NODE_ENV: "test" },
      },
      provider,
    );
    expect(plan).toMatchObject({
      params: {
        name: "example",
        image: "blaxel/base-image:latest",
        memory: 4096,
        region: "us-pdx-1",
        ttl: "1d",
        labels: { purpose: "test" },
        envs: [{ name: "NODE_ENV", value: "test" }],
        ports: [{ target: 3000, protocol: "TLS" }],
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
        provider,
      ).params.memory,
    ).toBe(128);
    expect(mapResources({ memoryMib: 4096n })).toEqual({ memory: 4096 });
  });

  it("accepts kit create timeout as the client deadline", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          provisioning: { timeout: 30_000 },
        },
        provider,
      ),
    ).not.toThrow();
  });

  it("rejects WaitCondition.Ready", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          provisioning: {
            waitFor: WaitCondition.Ready,
          },
        },
        provider,
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
        provider,
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
        provider,
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
        provider,
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

  it("classifies POLICY_VIOLATION as invalid argument", () => {
    const native = new ResponseError(
      new Response(null, { status: 403, statusText: "Forbidden" }),
      {
        error: {
          code: "POLICY_VIOLATION",
          message: "policy",
          status: 403,
        },
      },
      new Error("policy"),
    );
    expect(mapBlaxelError(native, "create")).toMatchObject({
      kind: ErrorKind.InvalidArgument,
      provider: "blaxel",
      operation: "create",
      statusCode: 403,
      providerCode: "POLICY_VIOLATION",
      cause: native,
    });
  });

  it("classifies bare HTTP 409 as already exists", () => {
    const native = new ResponseError(
      new Response(null, { status: 409, statusText: "Conflict" }),
      { error: { message: "conflict" } },
      new Error("conflict"),
    );
    expect(mapBlaxelError(native, "create")).toMatchObject({
      kind: ErrorKind.AlreadyExists,
      provider: "blaxel",
      operation: "create",
      statusCode: 409,
      cause: native,
    });
  });

  it("classifies bare HTTP 409 outside create as conflict", () => {
    const native = new ResponseError(
      new Response(null, { status: 409, statusText: "Conflict" }),
      { error: { message: "conflict" } },
      new Error("conflict"),
    );
    expect(mapBlaxelError(native, "delete")).toMatchObject({
      kind: ErrorKind.Conflict,
      provider: "blaxel",
      operation: "delete",
      statusCode: 409,
      cause: native,
    });
  });

  it("keeps SANDBOX_DELETION_IN_PROGRESS as conflict", () => {
    const native = new ResponseError(
      new Response(null, { status: 409, statusText: "Conflict" }),
      {
        error: {
          code: "SANDBOX_DELETION_IN_PROGRESS",
          message: "deleting",
          status: 409,
        },
      },
      new Error("deleting"),
    );
    expect(mapBlaxelError(native, "delete")).toMatchObject({
      kind: ErrorKind.Conflict,
      providerCode: "SANDBOX_DELETION_IN_PROGRESS",
      statusCode: 409,
      cause: native,
    });
  });

  it("classifies edge gateway 504 as timeout", () => {
    const native = new ResponseError(
      new Response(null, { status: 504, statusText: "Gateway Timeout" }),
      { error: { message: "gateway timeout" } },
      new Error("gateway timeout"),
    );
    expect(mapBlaxelError(native, "run")).toMatchObject({
      kind: ErrorKind.Timeout,
      provider: "blaxel",
      operation: "run",
      cause: native,
    });
    const waited = new Error(
      "Process did not finish in time (proc); it may still be running",
      { cause: native },
    );
    expect(mapBlaxelError(waited, "run")).toMatchObject({
      kind: ErrorKind.Timeout,
      operation: "run",
      statusCode: 504,
      cause: waited,
    });
  });

  it("reads provider codes from error.cause", () => {
    const native = new ResponseError(
      new Response(null, { status: 403, statusText: "Forbidden" }),
      {
        error: {
          code: "POLICY_VIOLATION",
          message: "policy",
          status: 403,
        },
      },
      new Error("policy"),
    );
    const wrapped = new Error("wrapped", { cause: native });
    expect(mapBlaxelError(wrapped, "create")).toMatchObject({
      kind: ErrorKind.InvalidArgument,
      providerCode: "POLICY_VIOLATION",
      statusCode: 403,
      cause: wrapped,
    });
  });

  it("accepts default egress with outbound domains", () => {
    const plan = mapCreateRequest(
      {
        name: "example",
        source: { image: { reference: "blaxel/base-image:latest" } },
        network: {
          egress: EgressMode.Default,
          outboundDomains: ["example.com"],
        },
      },
      provider,
    );
    expect(plan.params.network).toEqual({
      allowedDomains: ["example.com"],
    });
  });

  it("maps allocated memory onto SandboxInfo resources", () => {
    expect(
      mapAllocatedResources({
        metadata: { name: "demo" },
        spec: { runtime: { memory: 2048 } },
      }),
    ).toEqual({ memoryMib: 2048n });
    expect(
      mapSandboxInfo({
        metadata: { name: "demo" },
        status: "DEPLOYED",
        state: "RUNNING",
        spec: { region: "us-pdx-1", runtime: { memory: 4096 } },
      }),
    ).toMatchObject({
      id: "demo",
      resources: { memoryMib: 4096n },
    });
  });

  it("marks DELETED and TERMINATED refresh as terminal", () => {
    for (const status of ["DELETED", "TERMINATED", "DEACTIVATED"] as const) {
      expect(
        mapSandboxInfo({
          metadata: { name: "gone" },
          status,
        }),
      ).toMatchObject({
        id: "gone",
        state: SandboxState.Deleted,
        terminal: true,
        providerState: status,
      });
    }
  });

  it("omits protocol for default port transport", () => {
    const plan = mapCreateRequest(
      {
        name: "example",
        source: { image: { reference: "blaxel/base-image:latest" } },
        network: { ports: [{ port: 3000 }] },
      },
      provider,
    );
    expect(plan.params.ports).toEqual([{ target: 3000 }]);
  });

  it("rejects HTTP/2 TLS port transport", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          network: {
            ports: [{ port: 3000, transport: PortTransport.Http2Tls }],
          },
        },
        provider,
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldNetworkPortsTransport,
      }),
    );
  });

  it("applies environment checks before mapping", () => {
    expect(() =>
      mapCreateRequest(
        {
          name: "example",
          source: { image: { reference: "blaxel/base-image:latest" } },
          environment: { "NOT VALID": "value" },
        },
        provider,
      ),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldEnvironment,
      }),
    );
  });

  it("maps DEPLOYED plus runtime STANDBY to paused", () => {
    expect(
      mapSandboxIdentity({
        metadata: { name: "demo" },
        status: "DEPLOYED",
        state: "STANDBY",
        spec: { region: "us-pdx-1" },
      }),
    ).toMatchObject({
      id: "demo",
      providerState: "DEPLOYED",
      state: SandboxState.Paused,
      region: "us-pdx-1",
    });
  });

  it("maps DEPLOYED plus runtime RUNNING to running", () => {
    expect(
      mapSandboxIdentity({
        metadata: { name: "demo" },
        status: "DEPLOYED",
        state: "RUNNING",
      }),
    ).toMatchObject({
      providerState: "DEPLOYED",
      state: SandboxState.Running,
    });
  });

  it("maps projected runtime state from the native adapter", () => {
    expect(
      mapSandboxIdentity({
        metadata: { name: "demo" },
        status: "DEPLOYED",
        // defaultNativeFactory.projectHandle copies wire sandbox.state here.
        state: "STANDBY",
      }).state,
    ).toBe(SandboxState.Paused);
  });

  it("projectHandle copies private wire sandbox.state onto state", () => {
    const instance = new SandboxInstance({
      metadata: { name: "demo" },
      status: "DEPLOYED",
      spec: { region: "us-pdx-1" },
    } as never);
    // @blaxel/core types sandbox private; seed the wire model for the adapter.
    Object.assign(instance as object, {
      sandbox: {
        metadata: { name: "demo" },
        status: "DEPLOYED",
        state: "STANDBY",
        spec: { region: "us-pdx-1" },
      },
    });
    const handle = projectHandle(instance);
    expect(handle.state).toBe("STANDBY");
    expect(
      mapSandboxIdentity({
        metadata: { name: "demo" },
        status: "DEPLOYED",
        state: handle.state,
        spec: { region: "us-pdx-1" },
      }).state,
    ).toBe(SandboxState.Paused);
  });
});

describe("Blaxel provider orchestration", () => {
  it("creates with SandboxInstance.create and never createIfNotExists", async () => {
    let created: unknown;
    let configured: unknown;
    const nativeSandbox = {
      metadata: { name: "example" },
      status: "DEPLOYED",
      state: "RUNNING",
      spec: { region: "us-pdx-1" },
      fs: {} as never,
      process: {} as never,
    };
    const client = new Client({
      provider: new BlaxelProvider(() => ({
        configure(config) {
          configured = config;
        },
        async create(params) {
          created = params;
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
    expect(created).toMatchObject({
      name: "example",
      image: "blaxel/base-image:latest",
      memory: 2048,
      region: "us-pdx-1",
    });
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
          status: "DEPLOYED",
          state: "RUNNING",
          fs: {} as never,
          process: {} as never,
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

  it("unifies a successful delete as Deleted", async () => {
    const deleted: string[] = [];
    const backend = new BlaxelProvider(() => ({
      configure() {},
      async create() {
        throw new Error("unexpected create");
      },
      async get() {
        throw new Error("unexpected get");
      },
      async delete(name) {
        deleted.push(name);
        // Wire may still say DELETING; Kit reports Deleted.
        return { metadata: { name }, status: "DELETING" };
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

  it("does not poll get when wait is omitted", async () => {
    let gets = 0;
    const backend = new BlaxelProvider(() => ({
      configure() {},
      async create() {
        throw new Error("unexpected create");
      },
      async get() {
        gets += 1;
        throw new Error("unexpected get");
      },
      async delete(name) {
        return { metadata: { name }, status: "DELETING" };
      },
    })).initialize({
      auth: { apiKey: { key: "bl_key" } },
      scope: { workspace: "ws" },
    });
    if (backend.delete === undefined) throw new Error("delete is missing");
    await backend.delete("doomed", { wait: false }, {});
    expect(gets).toBe(0);
  });

  it("waits until get returns NotFound when wait is true", async () => {
    vi.useFakeTimers();
    try {
      let gets = 0;
      const backend = new BlaxelProvider(() => ({
        configure() {},
        async create() {
          throw new Error("unexpected create");
        },
        async get(name) {
          gets += 1;
          if (gets < 2) {
            return {
              metadata: { name },
              status: "DELETING",
              fs: {} as never,
              process: {} as never,
            };
          }
          throw new ResponseError(
            new Response(null, { status: 404, statusText: "Not Found" }),
            undefined,
            undefined,
          );
        },
        async delete(name) {
          return { metadata: { name }, status: "DELETING" };
        },
      })).initialize({
        auth: { apiKey: { key: "bl_key" } },
        scope: { workspace: "ws" },
      });
      if (backend.delete === undefined) throw new Error("delete is missing");
      const pending = backend.delete("doomed", { wait: true }, {});
      await vi.advanceTimersByTimeAsync(1_000);
      await vi.advanceTimersByTimeAsync(1_000);
      const info = await pending;
      expect(gets).toBe(2);
      expect(info).toMatchObject({
        id: "doomed",
        state: SandboxState.Deleted,
        terminal: true,
        providerState: "TERMINATED",
      });
    } finally {
      vi.useRealTimers();
    }
  });

  it("stops waiting for delete when the signal aborts", async () => {
    vi.useFakeTimers();
    try {
      const controller = new AbortController();
      const backend = new BlaxelProvider(() => ({
        configure() {},
        async create() {
          throw new Error("unexpected create");
        },
        async get(name) {
          return {
            metadata: { name },
            status: "DELETING",
            fs: {} as never,
            process: {} as never,
          };
        },
        async delete(name) {
          return { metadata: { name }, status: "DELETING" };
        },
      })).initialize({
        auth: { apiKey: { key: "bl_key" } },
        scope: { workspace: "ws" },
      });
      if (backend.delete === undefined) throw new Error("delete is missing");
      const pending = backend.delete("doomed", { wait: true }, {
        signal: controller.signal,
      });
      const expectation = expect(pending).rejects.toMatchObject({
        kind: ErrorKind.Canceled,
        operation: "delete",
      });
      await vi.advanceTimersByTimeAsync(500);
      controller.abort();
      await expectation;
    } finally {
      vi.useRealTimers();
    }
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

  it("starts without waitForCompletion and waits when still running", async () => {
    const seen: Array<{
      command: string;
      waitForCompletion?: boolean;
      timeout?: number;
      workingDir?: string;
      env?: Record<string, string>;
    }> = [];
    const waits: Array<{
      identifier: string;
      maxWait?: number;
      signal?: AbortSignal;
    }> = [];
    const signal = new AbortController().signal;
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
              if (request.timeout !== undefined)
                entry.timeout = request.timeout;
              if (request.workingDir !== undefined)
                entry.workingDir = request.workingDir;
              if (request.env !== undefined) entry.env = request.env;
              seen.push(entry);
              return {
                exitCode: 0,
                stdout: "",
                stderr: "",
                status: "running",
                pid: "long",
                name: "proc",
              };
            },
            async wait(identifier, options) {
              waits.push({
                identifier,
                ...(options?.maxWait === undefined
                  ? {}
                  : { maxWait: options.maxWait }),
                ...(options?.signal === undefined
                  ? {}
                  : { signal: options.signal }),
              });
              return {
                exitCode: 7,
                stdout: "out",
                stderr: "err",
                status: "failed",
              };
            },
          },
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

    const finished = await backend.run(
      "sandbox-1",
      { argv: ["echo", "hi"] },
      {},
    );
    await backend.run("sandbox-1", { argv: ["echo", "hi"], timeout: 0 }, {});
    await backend.run(
      "sandbox-1",
      {
        argv: ["echo", "hi"],
        timeout: 5_000,
        workingDirectory: "/workspace",
        environment: { MODE: "test" },
      },
      {},
    );
    await backend.run(
      "sandbox-1",
      { argv: ["echo", "hi"], timeout: 60_000 },
      {},
    );
    await backend.run(
      "sandbox-1",
      { argv: ["echo", "hi"], timeout: 61_000 },
      { signal },
    );
    await expect(
      backend.run("sandbox-1", { argv: ["echo", "hi"], timeout: 1_500 }, {}),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: "timeout",
      operation: "run",
    });

    expect(finished).toMatchObject({ exitCode: 7 });
    expect(finished.stdout).toEqual(new TextEncoder().encode("out"));
    expect(finished.stderr).toEqual(new TextEncoder().encode("err"));
    expect(seen).toEqual([
      { command: "'echo' 'hi'", waitForCompletion: false, timeout: 0 },
      { command: "'echo' 'hi'", waitForCompletion: false, timeout: 0 },
      {
        command: "'echo' 'hi'",
        waitForCompletion: false,
        timeout: 5,
        workingDir: "/workspace",
        env: { MODE: "test" },
      },
      { command: "'echo' 'hi'", waitForCompletion: false, timeout: 60 },
      { command: "'echo' 'hi'", waitForCompletion: false, timeout: 61 },
    ]);
    expect(waits).toEqual([
      { identifier: "long", maxWait: -1 },
      { identifier: "long", maxWait: -1 },
      { identifier: "long", maxWait: 5_000 },
      { identifier: "long", maxWait: 60_000 },
      { identifier: "long", maxWait: 61_000, signal },
    ]);
  });

  it("returns a finished process without waiting", async () => {
    const backend = runBackend({
      async exec() {
        return {
          status: "killed",
          exitCode: 9,
          stdout: "partial",
          stderr: "stopped",
        };
      },
      async wait() {
        throw new Error("wait must not run");
      },
    });
    const result = await backend.run!(
      "sandbox-1",
      { argv: ["sleep", "1"] },
      {},
    );
    expect(result.exitCode).toBe(9);
    expect(result.stdout).toEqual(new TextEncoder().encode("partial"));
    expect(result.stderr).toEqual(new TextEncoder().encode("stopped"));
  });

  it("waits when start still returns running", async () => {
    const waits: string[] = [];
    const backend = runBackend({
      async exec(request) {
        expect(request.waitForCompletion).toBe(false);
        expect(request.timeout).toBe(5);
        return { status: "running", pid: "held", stdout: "", stderr: "" };
      },
      async wait(identifier) {
        waits.push(identifier);
        return {
          status: "completed",
          exitCode: 0,
          stdout: "late",
          stderr: "",
        };
      },
    });
    const result = await backend.run!(
      "sandbox-1",
      { argv: ["echo", "hi"], timeout: 5_000 },
      {},
    );
    expect(waits).toEqual(["held"]);
    expect(result.stdout).toEqual(new TextEncoder().encode("late"));
  });

  it("rejects a running process with no identifier", async () => {
    const backend = runBackend({
      async exec() {
        return { status: "running", exitCode: 0, stdout: "", stderr: "" };
      },
      async wait() {
        throw new Error("wait must not run");
      },
    });
    await expect(
      backend.run!("sandbox-1", { argv: ["echo", "hi"] }, {}),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidResponse,
      operation: "run",
    });
  });

  it("maps an aborted process wait and a wait deadline", async () => {
    const controller = new AbortController();
    const aborted = runBackend({
      async exec() {
        return { status: "running", pid: "long", stdout: "", stderr: "" };
      },
      async wait(_identifier, options) {
        controller.abort();
        options?.signal?.throwIfAborted();
        return { status: "completed", exitCode: 0, stdout: "", stderr: "" };
      },
    });
    await expect(
      aborted.run!(
        "sandbox-1",
        { argv: ["sleep", "30"], timeout: 0 },
        { signal: controller.signal },
      ),
    ).rejects.toMatchObject({ kind: ErrorKind.Canceled, operation: "run" });

    const deadline = runBackend({
      async exec() {
        return { status: "running", pid: "long", stdout: "", stderr: "" };
      },
      async wait() {
        throw new Error(
          "Process did not finish in time (long); it may still be running",
        );
      },
    });
    await expect(
      deadline.run!(
        "sandbox-1",
        { argv: ["sleep", "90"], timeout: 61_000 },
        {},
      ),
    ).rejects.toMatchObject({ kind: ErrorKind.Timeout, operation: "run" });
  });
});

describe("Blaxel filesystem.stat", () => {
  it("returns directory metadata from get on the path", async () => {
    const backend = fsBackend({
      async get(path) {
        expect(path).toBe("/app");
        return {
          name: "app",
          path: "/app",
          files: [],
          subdirectories: [],
        };
      },
    });
    await expect(backend.stat!("sandbox-1", "/app", {}, {})).resolves.toEqual({
      name: "app",
      path: "/app",
      type: FileType.Directory,
      size: 0n,
    });
  });

  it("returns file metadata including size from get", async () => {
    const backend = fsBackend({
      async get(path) {
        expect(path).toBe("/app/readme.txt");
        return {
          name: "readme.txt",
          path: "/app/readme.txt",
          size: 12,
          permissions: "0644",
          owner: "user",
          group: "user",
          lastModified: "2024-01-02T03:04:05Z",
        };
      },
    });
    await expect(
      backend.stat!("sandbox-1", "/app/readme.txt", {}, {}),
    ).resolves.toMatchObject({
      name: "readme.txt",
      path: "/app/readme.txt",
      type: FileType.File,
      size: 12n,
      permissions: "0644",
      owner: "user",
      group: "user",
    });
  });

  it("propagates get errors", async () => {
    const missing = fsBackend({
      async get() {
        throw new Error("boom");
      },
    });
    await expect(
      missing.stat!("sandbox-1", "/missing", {}, {}),
    ).rejects.toMatchObject({ kind: ErrorKind.Unknown, operation: "stat" });
  });
});

function fsBackend(fs: {
  get(path: string): Promise<
    | {
        name: string;
        path: string;
        files: Array<{
          name: string;
          path: string;
          size: number;
          permissions?: string;
          owner?: string;
          group?: string;
          lastModified?: string;
        }>;
        subdirectories: Array<{ name: string; path: string }>;
      }
    | {
        name: string;
        path: string;
        size: number;
        permissions?: string;
        owner?: string;
        group?: string;
        lastModified?: string;
      }
  >;
}) {
  return new BlaxelProvider(() => ({
    configure() {},
    async create() {
      throw new Error("unexpected create");
    },
    async get() {
      return {
        metadata: { name: "sandbox-1" },
        fs: fs as never,
        process: {} as never,
      };
    },
    async delete() {
      throw new Error("unexpected delete");
    },
  })).initialize({
    auth: { apiKey: { key: "bl_key" } },
    scope: { workspace: "ws" },
  });
}

function runBackend(process: {
  exec(request: {
    command: string;
    waitForCompletion?: boolean;
    timeout?: number;
  }): Promise<{
    pid?: string;
    status?: string;
    exitCode?: number;
    stdout?: string;
    stderr?: string;
  }>;
  wait(
    identifier: string,
    options?: { signal?: AbortSignal },
  ): Promise<{
    exitCode?: number;
    stdout?: string;
    stderr?: string;
    status?: string;
  }>;
}) {
  return new BlaxelProvider(() => ({
    configure() {},
    async create() {
      throw new Error("unexpected create");
    },
    async get() {
      return {
        metadata: { name: "sandbox-1" },
        fs: {} as never,
        process,
      };
    },
    async delete() {
      throw new Error("unexpected delete");
    },
  })).initialize({
    auth: { apiKey: { key: "bl_key" } },
    scope: { workspace: "ws" },
  });
}

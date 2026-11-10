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
  ConfigFieldRegion,
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
import {
  mapResources,
  mapSandboxIdentity,
} from "#mappings/provider_mappings.js";
import { captureProviderState } from "#config/provider_config.js";
import { mapBlaxelError } from "#errors.js";
import { ResponseError } from "@blaxel/core";

const provider = { scope: { workspace: "ws" }, region: "us-pdx-1" };

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
        { scope: { workspace: "ws" } },
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
        labels: { purpose: "test" },
        environment: { NODE_ENV: "test" },
      },
      provider,
    );
    expect(plan.options).toBeUndefined();
    expect(plan).toMatchObject({
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

  it("reads runtime state from SandboxInstance wire model", () => {
    expect(
      mapSandboxIdentity({
        metadata: { name: "demo" },
        status: "DEPLOYED",
        // SandboxInstance exposes runtime state only on its private wire model.
        ...({ sandbox: { state: "STANDBY" } } as object),
      }).state,
    ).toBe(SandboxState.Paused);
  });
});

describe("Blaxel provider orchestration", () => {
  it("creates with SandboxInstance.create and never createIfNotExists", async () => {
    let created: { params: unknown; options: unknown } | undefined;
    let configured: unknown;
    const nativeSandbox = {
      metadata: { name: "example" },
      status: "DEPLOYED",
      state: "RUNNING",
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
          status: "DEPLOYED",
          state: "RUNNING",
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

  it("keeps short runs on waitForCompletion and polls the rest", async () => {
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
      interval?: number;
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
              if (request.waitForCompletion) {
                return {
                  exitCode: 0,
                  stdout: "ok",
                  stderr: "",
                  status: "completed",
                  pid: "short",
                };
              }
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
                ...(options?.interval === undefined
                  ? {}
                  : { interval: options.interval }),
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

    const polled = await backend.run(
      "sandbox-1",
      { argv: ["echo", "hi"] },
      {},
    );
    await backend.run("sandbox-1", { argv: ["echo", "hi"], timeout: 0 }, {});
    const held = await backend.run(
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

    expect(polled).toMatchObject({ exitCode: 7 });
    expect(polled.stdout).toEqual(new TextEncoder().encode("out"));
    expect(polled.stderr).toEqual(new TextEncoder().encode("err"));
    expect(held).toMatchObject({ exitCode: 0 });
    expect(held.stdout).toEqual(new TextEncoder().encode("ok"));
    expect(seen).toEqual([
      { command: "'echo' 'hi'", waitForCompletion: false, timeout: 0 },
      { command: "'echo' 'hi'", waitForCompletion: false, timeout: 0 },
      {
        command: "'echo' 'hi'",
        waitForCompletion: true,
        timeout: 5,
        workingDir: "/workspace",
        env: { MODE: "test" },
      },
      { command: "'echo' 'hi'", waitForCompletion: true, timeout: 60 },
      { command: "'echo' 'hi'", waitForCompletion: false, timeout: 61 },
    ]);
    expect(waits).toEqual([
      { identifier: "long", maxWait: -1, interval: 1_000 },
      { identifier: "long", maxWait: -1, interval: 1_000 },
      { identifier: "long", maxWait: 61_000, interval: 1_000, signal },
    ]);
  });

  it("returns a finished process without polling", async () => {
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
}

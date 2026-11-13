import { describe, expect, it } from "vitest";
import { Buffer } from "node:buffer";
import {
  ExecFieldArgv,
  ExecFieldShell,
  ExecFieldWorkingDirectory,
  FileFieldLocalPath,
  FileFieldPath,
  InfoFieldId,
  CreateFieldLifetimeIdleStopAfter,
  CreateFieldLifetimeIdleStopMode,
  CreateFieldNetwork,
  CreateFieldProvisioningTimeout,
  CreateFieldReadiness,
  CreateFieldResourcesCpuCores,
  CreateFieldResourcesCpuLimitCores,
  CreateFieldSecuritySecretsReference,
  CreateFieldSourceImageReference,
} from "#contracts/field_paths.js";
import {
  ConfigFieldEndpoint,
  ConfigFieldProvider,
  ConfigFieldTimeout,
} from "#contracts/config_paths.js";

import {
  Client,
  EgressMode,
  ErrorKind,
  SandboxError,
  sandboxError,
  PolicyMode,
  validateCreateOptions,
  type Backend,
  type CreateOptions,
  type CreateResult,
  type ExitResult,
  type OperationOptions,
  type Provider,
  LifecycleCapability,
  SandboxState,
  type SandboxInfo,
  type StartOptions,
  type StopOptions,
  type TerminateOptions,
} from "../index.js";

class RecordingBackend implements Backend {
  readonly name = "test";
  closeCalls = 0;
  closeSandboxCalls = 0;
  createCalls = 0;
  terminated = false;
  lastCreateRequest: CreateOptions | undefined;
  // Optional Backend methods tests attach at runtime.
  run?: Backend["run"];
  readFile?: Backend["readFile"];

  async create(
    request: CreateOptions,
    _options: OperationOptions,
  ): Promise<CreateResult> {
    this.createCalls++;
    this.lastCreateRequest = request;
    return { sandbox: { id: "sandbox-1", provider: this.name } };
  }

  async close(_options: OperationOptions): Promise<void> {
    this.closeCalls++;
  }

  async closeSandbox(_id: string, _options: OperationOptions): Promise<void> {
    this.closeSandboxCalls++;
  }

  async attach(id: string): Promise<SandboxInfo> {
    return { id, provider: this.name, state: SandboxState.Running };
  }

  async terminate(_id: string, options: TerminateOptions): Promise<ExitResult> {
    this.terminated = true;
    return options.wait ? { exitCode: 0 } : {};
  }
}

function fixture(): {
  readonly client: Client;
  readonly backend: RecordingBackend;
} {
  const backend = new RecordingBackend();
  const provider: Provider = { name: backend.name, initialize: () => backend };
  return { client: new Client({ provider }), backend };
}

class HangingBackend extends RecordingBackend {
  override async create(
    _request: CreateOptions,
    _options: OperationOptions,
  ): Promise<CreateResult> {
    this.createCalls++;
    return new Promise<CreateResult>(() => {});
  }

  override async terminate(
    _id: string,
    _options: TerminateOptions,
  ): Promise<ExitResult> {
    return new Promise<ExitResult>(() => {});
  }
}

function hangingFixture(): {
  readonly client: Client;
  readonly backend: HangingBackend;
} {
  const backend = new HangingBackend();
  const provider: Provider = { name: backend.name, initialize: () => backend };
  return { client: new Client({ provider }), backend };
}

describe("Client cleanup", () => {
  it("starts cleanup when a provider returns the wrong backend identity", async () => {
    const backend = new RecordingBackend();
    expect(
      () =>
        new Client({
          provider: { name: "different", initialize: () => backend },
        }),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidResponse,
        field: ConfigFieldProvider,
      }),
    );
    await Promise.resolve();
    expect(backend.closeCalls).toBe(1);
  });

  it("sends each close to the provider", async () => {
    const { client, backend } = fixture();
    await Promise.all([client.close(), client.close()]);
    expect(backend.closeCalls).toBe(2);
  });

  it("supports explicit async disposal", async () => {
    const { client, backend } = fixture();
    await client[Symbol.asyncDispose]();
    expect(backend.closeCalls).toBe(1);
  });

  it("supports await using as an optional convenience", async () => {
    const { client, backend } = fixture();
    {
      await using scopedClient = client;
      expect(scopedClient.providerName).toBe("test");
    }
    expect(backend.closeCalls).toBe(1);
  });

  it("calls the provider again after close", async () => {
    const { client, backend } = fixture();
    await client.close();
    await expect(client.create()).resolves.toBeDefined();
    expect(backend.closeCalls).toBe(1);
    expect(backend.createCalls).toBe(1);
    await client.close();
    expect(backend.closeCalls).toBe(2);
  });

  it("leaves a handle usable after its client closes", async () => {
    const { client, backend } = fixture();
    const sandbox = await client.attach("sandbox-existing");
    await client.close();

    await expect(sandbox.terminate()).resolves.toEqual({});
    expect(backend.terminated).toBe(true);
    expect(sandbox.supports(LifecycleCapability.Terminate)).toBe(true);
  });

  it("supports management-only providers without a create capability", async () => {
    const backend: Backend = {
      name: "manage",
      close: async () => {},
      closeSandbox: async () => {},
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });

    await expect(client.create()).rejects.toMatchObject({
      kind: ErrorKind.Unsupported,
      operation: "create",
    });
  });
});

describe("Provider response validation", () => {
  it("reports the public get operation and preserves the provider error", async () => {
    const native = new Error("native failure");
    const providerError = sandboxError(
      {
        kind: ErrorKind.Unknown,
        provider: "wrong-provider",
        operation: "refresh",
        field: "",
        message: "lookup failed",
      },
      native,
    );
    const backend: Backend = {
      name: "test",
      close: async () => {},
      closeSandbox: async () => {},
      refresh: async () => {
        throw providerError;
      },
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });

    const error = await client.get("requested").catch((reason: unknown) => reason);
    expect(error).toBeInstanceOf(SandboxError);
    expect(error).toMatchObject({
      provider: "test",
      operation: "get",
      message: "lookup failed",
    });
    expect((error as Error).cause).toBe(providerError);
    expect(providerError.cause).toBe(native);
  });

  it("rejects lookup responses with a different identity or provider", async () => {
    const attachBackend: Backend = {
      name: "test",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async () => ({ id: "different", provider: "test" }),
    };
    const attachClient = new Client({
      provider: { name: attachBackend.name, initialize: () => attachBackend },
    });
    await expect(attachClient.attach("requested")).rejects.toMatchObject({
      kind: ErrorKind.InvalidResponse,
      operation: "attach",
      field: InfoFieldId,
    });

    const getBackend: Backend = {
      name: "test",
      close: async () => {},
      closeSandbox: async () => {},
      refresh: async (id) => ({ id, provider: "different" }),
    };
    const getClient = new Client({
      provider: { name: getBackend.name, initialize: () => getBackend },
    });
    await expect(getClient.get("requested")).rejects.toMatchObject({
      kind: ErrorKind.InvalidResponse,
      operation: "get",
      field: ConfigFieldProvider,
    });
  });

  it("preserves a handle when a lifecycle response changes its identity", async () => {
    const backend: Backend = {
      name: "test",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({ id, provider: "test" }),
      refresh: async () => ({ id: "different", provider: "test" }),
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });
    const sandbox = await client.attach("requested");

    await expect(sandbox.refresh()).rejects.toMatchObject({
      kind: ErrorKind.InvalidResponse,
      operation: "refresh",
      field: InfoFieldId,
    });
    expect(sandbox.id).toBe("requested");
  });
});

describe("Shared request validation", () => {
  it("exposes validateCreateOptions for public CreateOptions", () => {
    expect(() => validateCreateOptions({})).not.toThrow();
    expect(() => validateCreateOptions({ resources: { cpuCores: 0 } })).toThrow(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldResourcesCpuCores,
        rule: "exclusive_minimum",
      }),
    );
    expect(() =>
      validateCreateOptions({
        lifetime: { idleStop: { mode: PolicyMode.After } },
      }),
    ).toThrow(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldLifetimeIdleStopAfter,
        rule: "requires_any",
      }),
    );
    expect(() =>
      validateCreateOptions({
        resources: { cpuCores: 2, cpuLimitCores: 1 },
      }),
    ).toThrow(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldResourcesCpuLimitCores,
        rule: "limit",
      }),
    );
    expect(() =>
      validateCreateOptions({
        network: { egress: EgressMode.Restricted },
      }),
    ).toThrow(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldNetwork,
      }),
    );
    expect(() =>
      validateCreateOptions({ readiness: {} } as CreateOptions),
    ).toThrow(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: CreateFieldReadiness,
      }),
    );
  });

  it("rejects invalid portable intent before the selected provider is called", async () => {
    const { client, backend } = fixture();

    await expect(
      client.create({ resources: { cpuCores: 0 } }),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: CreateFieldResourcesCpuCores,
      rule: "exclusive_minimum",
    });
    expect(backend.createCalls).toBe(0);
  });

  it("validates nested source choices using the shared schema", async () => {
    const { client, backend } = fixture();

    await expect(
      client.create({
        source: {
          image: { reference: "not an image reference!" },
        },
      }),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: CreateFieldSourceImageReference,
    });
    expect(backend.createCalls).toBe(0);
  });

  it("validates descendants through containers without direct rules", async () => {
    const { client, backend } = fixture();

    await expect(
      client.create({ security: { secrets: [{ reference: "" }] } }),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: CreateFieldSecuritySecretsReference,
      message: "sandbox-kit: must not be blank",
    });
    expect(backend.createCalls).toBe(0);
  });

  it("rejects unknown runtime enum values from JavaScript callers", async () => {
    const { client, backend } = fixture();
    const request = {
      lifetime: { idleStop: { mode: 99 } },
    } as unknown as CreateOptions;

    await expect(client.create(request)).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: CreateFieldLifetimeIdleStopMode,
    });
    expect(backend.createCalls).toBe(0);
  });

  it("allows callers to omit proto3 default enum fields", async () => {
    const { client, backend } = fixture();

    await expect(
      client.create({
        isolation: {},
        network: { ports: [{ port: 8080 }] },
        provisioning: { timeout: 1_000 },
        security: { secrets: [{ reference: "service-secret" }] },
      }),
    ).resolves.toBeDefined();
    expect(backend.createCalls).toBe(1);
  });

  it("rejects non-finite duration values", async () => {
    const { client, backend } = fixture();

    await expect(
      client.create({ provisioning: { timeout: -1 } }),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: CreateFieldProvisioningTimeout,
    });
    expect(backend.createCalls).toBe(0);
  });
});

describe("Public create options for providers", () => {
  it("keeps nested create input owned by the caller", async () => {
    const request: CreateOptions = {
      environment: { ORIGINAL: "value" },
      runtime: { entrypoint: ["echo"] },
    };
    const backend: Backend = {
      name: "test",
      close: async () => {},
      closeSandbox: async () => {},
      create: async (providerRequest) => {
        providerRequest.environment!.ORIGINAL = "mutated";
        providerRequest.runtime!.entrypoint!.push("mutated");
        return { sandbox: { id: "sandbox-1", provider: "test" } };
      },
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });

    await client.create(request);

    expect(request).toEqual({
      environment: { ORIGINAL: "value" },
      runtime: { entrypoint: ["echo"] },
    });
  });

  it("keeps nested client configuration owned by the caller", () => {
    const auth = { apiKey: { key: "secret" } } as const;
    const provider: Provider = {
      name: "test",
      initialize: (config) => {
        config.auth!.apiKey!.key = "mutated";
        return new RecordingBackend();
      },
    };

    new Client({ provider, auth });

    expect(auth.apiKey.key).toBe("secret");
  });

  it("hands the provider public create options after a create call", async () => {
    const { client, backend } = fixture();

    await client.create({
      source: { image: { reference: "node:24" } },
      provisioning: { timeout: 5_000 },
    });

    expect(backend.lastCreateRequest?.source).toMatchObject({
      image: { reference: "node:24" },
    });
    expect(backend.lastCreateRequest?.provisioning?.timeout).toBe(5_000);
  });

  it("writes the resolved create timeout into provisioning for providers", async () => {
    const { client, backend } = fixture();

    await client.create(
      { source: { image: { reference: "node:24" } } },
      { timeout: 15_000 },
    );

    expect(backend.lastCreateRequest?.provisioning?.timeout).toBe(15_000);
  });
});

describe("Provider input ownership", () => {
  it("keeps process arrays and maps owned by the caller", async () => {
    const argv = ["echo", "hello"];
    const environment = { ORIGINAL: "value" };
    const backend: Backend = {
      name: "test",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({ id, provider: "test" }),
      run: async (_id, request) => {
        request.argv!.push("mutated");
        request.environment!.ORIGINAL = "mutated";
        return { exitCode: 0 };
      },
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });
    const sandbox = await client.attach("sandbox-1");

    await sandbox.process.run({ argv, environment });

    expect(argv).toEqual(["echo", "hello"]);
    expect(environment).toEqual({ ORIGINAL: "value" });
  });

  it("keeps write buffers owned by the caller", async () => {
    const data = Buffer.from([1, 2, 3]);
    const backend: Backend = {
      name: "test",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({ id, provider: "test" }),
      writeFile: async (_id, _path, providerData) => {
        providerData[0] = 9;
      },
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });
    const sandbox = await client.attach("sandbox-1");

    await sandbox.filesystem.writeFile("/tmp/file", data);

    expect([...data]).toEqual([1, 2, 3]);
  });

  it("captures process and file inputs before queued operations run", async () => {
    let release!: () => void;
    const blocked = new Promise<void>((resolve) => { release = resolve; });
    let observedArgv: string[] | undefined;
    let observedEnvironment: Record<string, string> | undefined;
    let observedData: number[] | undefined;
    const backend: Backend = {
      name: "test",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({ id, provider: "test" }),
      start: async (id) => { await blocked; return { id, provider: "test" }; },
      run: async (_id, request) => {
        observedArgv = request.argv;
        observedEnvironment = request.environment;
        return { exitCode: 0 };
      },
      writeFile: async (_id, _path, data) => { observedData = [...data]; },
    };
    const client = new Client({ provider: { name: "test", initialize: () => backend } });
    const sandbox = await client.attach("sandbox-1");
    const first = sandbox.start();
    const argv = ["echo", "hello"];
    const environment = { KEY: "original" };
    const data = Buffer.from([1, 2, 3]);
    const command = sandbox.process.run({ argv, environment });
    const write = sandbox.filesystem.writeFile("/tmp/file", data);
    argv.length = 0;
    environment.KEY = "changed";
    data[0] = 9;
    release();
    await Promise.all([first, command, write]);
    expect(observedArgv).toEqual(["echo", "hello"]);
    expect(observedEnvironment).toEqual({ KEY: "original" });
    expect(observedData).toEqual([1, 2, 3]);
  });
});

describe("Shared client configuration validation", () => {
  it.each([
    ["https://api.example.test/api", true],
    ["http://127.0.0.1:1234", true],
    ["http://[::1]:1234", true],
    ["http://[0:0:0:0:0:0:0:1]:1234", true],
    ["http://localhost:1234", true],
    ["http://LOCALHOST:1234", true],
    ["http://api.example.test", false],
    ["http://127.0.0.1.evil.test", false],
    ["http://0.0.0.0:1234", false],
    ["https://user:secret@api.example.test", false],
    ["https://api.example.test/#fragment", false],
    ["http://@localhost", false],
    ["http://127.1", false],
  ])("validates secure endpoint %s as %s", (endpoint, valid) => {
    const provider: Provider = {
      name: "test",
      initialize: () => new RecordingBackend(),
    };
    const construct = () => new Client({ provider, endpoint });

    if (valid) expect(construct).not.toThrow();
    else
      expect(construct).toThrowError(
        expect.objectContaining({
          kind: ErrorKind.InvalidArgument,
          field: ConfigFieldEndpoint,
        }),
      );
  });

  it("rejects insecure remote endpoints before initializing the provider", () => {
    const backend = new RecordingBackend();
    let initialized = false;
    const provider: Provider = {
      name: backend.name,
      initialize: () => {
        initialized = true;
        return backend;
      },
    };

    expect(
      () => new Client({ provider, endpoint: "http://remote.example/api" }),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: ConfigFieldEndpoint,
        operation: "initialize",
      }),
    );
    expect(initialized).toBe(false);
  });

  it("rejects a non-finite default timeout before initialization", () => {
    const backend = new RecordingBackend();
    let initialized = false;
    const provider: Provider = {
      name: backend.name,
      initialize: () => {
        initialized = true;
        return backend;
      },
    };

    expect(
      () =>
        new Client({
          provider,
          timeout: -1,
        }),
    ).toThrowError(
      expect.objectContaining({
        kind: ErrorKind.InvalidArgument,
        field: ConfigFieldTimeout,
        operation: "initialize",
      }),
    );
    expect(initialized).toBe(false);
  });
});

describe("Operation cancellation", () => {
  it("does not invoke the provider for an already-aborted operation", async () => {
    const { client, backend } = fixture();
    const controller = new AbortController();
    controller.abort(new Error("caller stopped work"));

    await expect(
      client.create({}, { signal: controller.signal }),
    ).rejects.toMatchObject({
      kind: ErrorKind.Canceled,
      operation: "create",
    });
    expect(backend.createCalls).toBe(0);
  });

  it("closes the provider after a rejected create", async () => {
    const { client, backend } = fixture();

    await expect(client.create({}, { timeout: -1 })).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: ConfigFieldTimeout,
    });
    await expect(client.close()).resolves.toBeUndefined();
    expect(backend.closeCalls).toBe(1);
  });

  it("enforces a caller timeout when a provider operation does not settle", async () => {
    const { client, backend } = hangingFixture();

    await expect(client.create({}, { timeout: 10 })).rejects.toMatchObject({
      kind: ErrorKind.Timeout,
      operation: "create",
    });
    expect(backend.createCalls).toBe(1);
  });

  it("maps an external abort signal to a canceled error", async () => {
    const { client } = hangingFixture();
    const controller = new AbortController();
    const operation = client.create({}, { signal: controller.signal });
    controller.abort(new Error("caller stopped work"));

    await expect(operation).rejects.toMatchObject({
      kind: ErrorKind.Canceled,
      operation: "create",
    });
  });

  it("applies request timeouts to lifecycle operations", async () => {
    const { client } = hangingFixture();
    const sandbox = await client.attach("sandbox-existing");

    await expect(sandbox.terminate({ timeout: 10 })).rejects.toMatchObject({
      kind: ErrorKind.Timeout,
      operation: "terminate",
    });
  });
});

describe("Sandbox lifecycle", () => {
  it("calls handle close on every provider", async () => {
    const { client, backend } = fixture();
    const sandbox = await client.attach("sandbox-existing");
    await sandbox.close();
    expect(backend.closeSandboxCalls).toBe(1);
    expect(sandbox.supports(LifecycleCapability.Terminate)).toBe(true);
  });

  it("attaches, reports capabilities, and terminates through the provider", async () => {
    const { client, backend } = fixture();
    const sandbox = await client.attach("sandbox-existing");
    expect(sandbox.supports(LifecycleCapability.Terminate)).toBe(true);
    await expect(sandbox.terminate({ wait: true })).resolves.toEqual({
      exitCode: 0,
    });
    expect(backend.terminated).toBe(true);
    expect(sandbox.supports(LifecycleCapability.Terminate)).toBe(true);
    expect(sandbox.info().terminal).toBeUndefined();
    await expect(sandbox.terminate()).resolves.toEqual({});
  });

  it("keeps an omitted lifecycle timeout off the provider request", async () => {
    const seen: Array<StartOptions | StopOptions> = [];
    const backend: Backend = {
      name: "omit-timeout",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({
        id,
        provider: "omit-timeout",
        state: SandboxState.Running,
      }),
      start: async (id, request) => {
        seen.push(request);
        return { id, provider: "omit-timeout", state: SandboxState.Running };
      },
      stop: async (id, request) => {
        seen.push(request);
        return { id, provider: "omit-timeout", state: SandboxState.Stopped };
      },
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
      timeout: 60_000,
    });
    const sandbox = await client.attach("sandbox-existing");

    await sandbox.start();
    await sandbox.stop({ force: true });
    await sandbox.start({ timeout: 0 });

    expect(seen).toEqual([{}, { force: true }, { timeout: 0 }]);
  });

  it("uses the client timeout as the deadline when a lifecycle call omits one", async () => {
    const backend: Backend = {
      name: "inherit-deadline",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({ id, provider: "inherit-deadline" }),
      start: async () => new Promise(() => {}),
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
      timeout: 10,
    });
    const sandbox = await client.attach("sandbox-existing");

    await expect(sandbox.start()).rejects.toMatchObject({
      kind: ErrorKind.Timeout,
      operation: "start",
    });
  });

  it("rejects lifecycle operations absent from the selected provider", async () => {
    const { client } = fixture();
    const sandbox = await client.attach("sandbox-existing");
    await expect(sandbox.pause()).rejects.toMatchObject({
      kind: ErrorKind.Unsupported,
      operation: "pause",
    });
  });

  it("serializes lifecycle transitions on one handle", async () => {
    let active = 0;
    let maximumActive = 0;
    let releaseFirst!: () => void;
    const firstBlocked = new Promise<void>((resolve) => {
      releaseFirst = resolve;
    });
    let calls = 0;
    const backend: Backend = {
      name: "serial",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({
        id,
        provider: "serial",
        state: SandboxState.Running,
      }),
      start: async (id) => {
        active++;
        maximumActive = Math.max(maximumActive, active);
        calls++;
        if (calls === 1) await firstBlocked;
        active--;
        return { id, provider: "serial", state: SandboxState.Running };
      },
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });
    const sandbox = await client.attach("sandbox-existing");

    const first = sandbox.start();
    const second = sandbox.start();
    await Promise.resolve();
    releaseFirst();
    await Promise.all([first, second]);

    expect(maximumActive).toBe(1);
  });

  it("preserves serialization when a queued transition is canceled", async () => {
    let active = 0;
    let maximumActive = 0;
    let releaseFirst!: () => void;
    const firstBlocked = new Promise<void>((resolve) => {
      releaseFirst = resolve;
    });
    let calls = 0;
    const backend: Backend = {
      name: "serial-cancel",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({
        id,
        provider: "serial-cancel",
        state: SandboxState.Running,
      }),
      start: async (id) => {
        active++;
        maximumActive = Math.max(maximumActive, active);
        calls++;
        if (calls === 1) await firstBlocked;
        active--;
        return { id, provider: "serial-cancel", state: SandboxState.Running };
      },
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });
    const sandbox = await client.attach("sandbox-existing");
    const controller = new AbortController();

    const first = sandbox.start();
    const canceled = sandbox.start({ signal: controller.signal });
    controller.abort();
    await expect(canceled).rejects.toMatchObject({ kind: ErrorKind.Canceled });

    const third = sandbox.start();
    await Promise.resolve();
    expect(calls).toBe(1);

    releaseFirst();
    await Promise.all([first, third]);
    expect(maximumActive).toBe(1);
  });

  it("holds serialization until an aborted in-flight provider call settles", async () => {
    let active = 0;
    let maximumActive = 0;
    let releaseFirst!: () => void;
    const firstBlocked = new Promise<void>((resolve) => {
      releaseFirst = resolve;
    });
    const backend: Backend = {
      name: "serial-inflight-abort",
      close: async () => {},
      closeSandbox: async () => {},
      attach: async (id) => ({
        id,
        provider: "serial-inflight-abort",
        state: SandboxState.Running,
      }),
      start: async (id) => {
        active++;
        maximumActive = Math.max(maximumActive, active);
        await firstBlocked;
        active--;
        return {
          id,
          provider: "serial-inflight-abort",
          state: SandboxState.Running,
        };
      },
    };
    const client = new Client({
      provider: { name: backend.name, initialize: () => backend },
    });
    const sandbox = await client.attach("sandbox-existing");
    const controller = new AbortController();

    const first = sandbox.start({ signal: controller.signal });
    await Promise.resolve();
    controller.abort();
    await expect(first).rejects.toMatchObject({ kind: ErrorKind.Canceled });

    const second = sandbox.start();
    await Promise.resolve();
    expect(active).toBe(1);
    expect(maximumActive).toBe(1);

    releaseFirst();
    await second;
    expect(maximumActive).toBe(1);
  });

  it("holds process serialization until an aborted in-flight run settles", async () => {
    let active = 0;
    let maximumActive = 0;
    let releaseFirst!: () => void;
    const firstBlocked = new Promise<void>((resolve) => {
      releaseFirst = resolve;
    });
    const { client, backend } = fixture();
    backend.run = async () => {
      active++;
      maximumActive = Math.max(maximumActive, active);
      await firstBlocked;
      active--;
      return {
        exitCode: 0,
        stdout: new Uint8Array(),
        stderr: new Uint8Array(),
      };
    };
    const sandbox = await client.create();
    const controller = new AbortController();

    const first = sandbox.process.run({
      argv: ["sleep", "1"],
      signal: controller.signal,
    });
    await Promise.resolve();
    controller.abort();
    await expect(first).rejects.toMatchObject({ kind: ErrorKind.Canceled });

    const second = sandbox.process.run({ argv: ["true"] });
    await Promise.resolve();
    expect(active).toBe(1);
    expect(maximumActive).toBe(1);

    releaseFirst();
    await second;
    expect(maximumActive).toBe(1);
  });

  it("holds filesystem serialization until an aborted in-flight read settles", async () => {
    let active = 0;
    let maximumActive = 0;
    let releaseFirst!: () => void;
    const firstBlocked = new Promise<void>((resolve) => {
      releaseFirst = resolve;
    });
    const { client, backend } = fixture();
    backend.readFile = async () => {
      active++;
      maximumActive = Math.max(maximumActive, active);
      await firstBlocked;
      active--;
      return new TextEncoder().encode("ok");
    };
    const sandbox = await client.create();
    const controller = new AbortController();

    const first = sandbox.filesystem.readFile("/tmp/a.txt", {
      signal: controller.signal,
    });
    await Promise.resolve();
    controller.abort();
    await expect(first).rejects.toMatchObject({ kind: ErrorKind.Canceled });

    const second = sandbox.filesystem.readFile("/tmp/b.txt");
    await Promise.resolve();
    expect(active).toBe(1);
    expect(maximumActive).toBe(1);

    releaseFirst();
    await second;
    expect(maximumActive).toBe(1);
  });

  it("rejects a relative sandbox path before the provider call", async () => {
    const { client, backend } = fixture();
    const files: Backend = backend;
    let calls = 0;
    files.readFile = async () => {
      calls++;
      return new Uint8Array();
    };
    const sandbox = await client.create();
    await expect(
      sandbox.filesystem.readFile("relative.txt"),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      operation: "read_file",
      field: FileFieldPath,
    });
    expect(calls).toBe(0);
  });

  it("reads an absolute path through the provider", async () => {
    const { client, backend } = fixture();
    const files: Backend = backend;
    files.readFile = async (_id, path) => {
      expect(path).toBe("/workspace/note.txt");
      return new TextEncoder().encode("hello");
    };
    const sandbox = await client.create();
    const data = await sandbox.filesystem.readFile("/workspace/note.txt");
    expect(new TextDecoder().decode(data)).toBe("hello");
  });

  it("requires a local path for copy", async () => {
    const { client } = fixture();
    const sandbox = await client.create();
    await expect(
      sandbox.filesystem.copyFromLocal("", "/workspace/note.txt"),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      operation: "copy_from_local",
      field: FileFieldLocalPath,
    });
  });

  it("reports filesystem methods the provider does not implement", async () => {
    const { client } = fixture();
    const sandbox = await client.create();
    await expect(sandbox.filesystem.stat("/workspace")).rejects.toMatchObject({
      kind: ErrorKind.Unsupported,
      operation: "stat",
    });
  });

  it("rejects a process call without argv or shell before the provider", async () => {
    const { client, backend } = fixture();
    let calls = 0;
    backend.run = async () => {
      calls++;
      return {
        exitCode: 0,
        stdout: new Uint8Array(),
        stderr: new Uint8Array(),
      };
    };
    const sandbox = await client.create();
    await expect(sandbox.process.run({})).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      operation: "run",
      field: ExecFieldArgv,
    });
    expect(calls).toBe(0);
  });

  it("rejects a blank shell before the provider", async () => {
    const { client, backend } = fixture();
    backend.run = async () => {
      throw new Error("provider was called");
    };
    const sandbox = await client.create();
    for (const shell of ["", "   "]) {
      await expect(sandbox.process.run({ shell })).rejects.toMatchObject({
        kind: ErrorKind.InvalidArgument,
        operation: "run",
        field: ExecFieldShell,
      });
    }
  });

  it("rejects a relative working directory before the provider", async () => {
    const { client, backend } = fixture();
    backend.run = async () => {
      throw new Error("provider was called");
    };
    const sandbox = await client.create();
    await expect(
      sandbox.process.run({ argv: ["pwd"], workingDirectory: "relative" }),
    ).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      operation: "run",
      field: ExecFieldWorkingDirectory,
    });
  });

  it("runs argv through the provider and keeps a non-zero exit code", async () => {
    const { client, backend } = fixture();
    backend.run = async (_id, request) => {
      expect(request.argv).toEqual(["echo", "hello"]);
      return {
        exitCode: 7,
        stdout: new TextEncoder().encode("hello\n"),
        stderr: new Uint8Array(),
      };
    };
    const sandbox = await client.create();
    const result = await sandbox.process.run({ argv: ["echo", "hello"] });
    expect(result.exitCode).toBe(7);
    expect(new TextDecoder().decode(result.stdout)).toBe("hello\n");
  });

  it("reports process methods the provider does not implement", async () => {
    const { client } = fixture();
    const sandbox = await client.create();
    await expect(sandbox.process.run({ argv: ["true"] })).rejects.toMatchObject(
      {
        kind: ErrorKind.Unsupported,
        operation: "run",
      },
    );
  });
});

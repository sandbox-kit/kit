import { describe, expect, it } from "vitest";
import {
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

describe("Shared client configuration validation", () => {
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
});

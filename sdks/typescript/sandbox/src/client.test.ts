import { describe, expect, it } from "vitest";

import {
  Client,
  ErrorKind,
  type Backend,
  type CreateOptions,
  type CreateResult,
  type OperationOptions,
  type Provider,
  LifecycleCapability,
  type SandboxInfo,
  type TerminateOptions,
  type ExitResult,
} from "./index.js";

class RecordingBackend implements Backend {
  readonly name = "test";
  closeCalls = 0;
  createCalls = 0;
  terminated = false;

  async create(_request: CreateOptions, _options: OperationOptions): Promise<CreateResult> {
    this.createCalls++;
    return { sandbox: { id: "sandbox-1", provider: this.name, state: "unknown" } };
  }

  async close(_options: OperationOptions): Promise<void> {
    this.closeCalls++;
  }

  async attach(id: string): Promise<SandboxInfo> {
    return { id, provider: this.name, state: "running" };
  }

  capabilities(): readonly LifecycleCapability[] { return [LifecycleCapability.Terminate]; }

  async terminate(_id: string, options: TerminateOptions): Promise<ExitResult> {
    this.terminated = true;
    return options.wait ? { exitCode: 0 } : {};
  }
}

function fixture(): { readonly client: Client; readonly backend: RecordingBackend } {
  const backend = new RecordingBackend();
  const provider: Provider = { name: backend.name, initialize: () => backend };
  return { client: new Client({ provider }), backend };
}

class HangingBackend extends RecordingBackend {
  override async create(_request: CreateOptions, _options: OperationOptions): Promise<CreateResult> {
    this.createCalls++;
    return new Promise<CreateResult>(() => {});
  }

  override async terminate(_id: string, _options: TerminateOptions): Promise<ExitResult> {
    return new Promise<ExitResult>(() => {});
  }
}

function hangingFixture(): { readonly client: Client; readonly backend: HangingBackend } {
  const backend = new HangingBackend();
  const provider: Provider = { name: backend.name, initialize: () => backend };
  return { client: new Client({ provider }), backend };
}

describe("Client cleanup", () => {
  it("makes close idempotent", async () => {
    const { client, backend } = fixture();
    await Promise.all([client.close(), client.close()]);
    expect(backend.closeCalls).toBe(1);
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

  it("rejects operations after close without deleting a sandbox", async () => {
    const { client, backend } = fixture();
    await client.close();
    await expect(client.create()).rejects.toMatchObject({
      kind: ErrorKind.FailedPrecondition,
      operation: "create",
    });
    expect(backend.createCalls).toBe(0);
  });
});

describe("Shared request validation", () => {
  it("rejects invalid portable intent before the selected provider is called", async () => {
    const { client, backend } = fixture();

    await expect(client.create({ resources: { cpuCores: 0 } })).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: "resources.cpuCores",
    });
    expect(backend.createCalls).toBe(0);
  });

  it("validates nested source choices using the shared schema", async () => {
    const { client, backend } = fixture();

    await expect(client.create({ source: { kind: "image", image: { reference: "not an image reference!" } } })).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: "source.image.reference",
    });
    expect(backend.createCalls).toBe(0);
  });

  it("validates descendants through containers without direct rules", async () => {
    const { client, backend } = fixture();

    await expect(client.create({ security: { secrets: [{ reference: "", injection: "default" }] } })).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: "security.secrets[].reference",
      message: "sandbox-kit: must not be blank",
    });
    expect(backend.createCalls).toBe(0);
  });

  it("rejects unknown runtime enum values from JavaScript callers", async () => {
    const { client, backend } = fixture();
    const request = { lifetime: { idleStop: { mode: "unknown-mode" } } } as unknown as CreateOptions;

    await expect(client.create(request)).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      field: "lifetime.idleStop.mode",
    });
    expect(backend.createCalls).toBe(0);
  });

  it("allows callers to omit proto3 default enum fields", async () => {
    const { client, backend } = fixture();

    await expect(client.create({
      isolation: {},
      network: { ports: [{ port: 8080 }] },
      provisioning: { timeout: 1_000 },
      security: { secrets: [{ reference: "service-secret" }] },
    })).resolves.toBeDefined();
    expect(backend.createCalls).toBe(1);
  });
});

describe("Shared client configuration validation", () => {
  it("rejects insecure remote endpoints before initializing the provider", () => {
    const backend = new RecordingBackend();
    let initialized = false;
    const provider: Provider = { name: backend.name, initialize: () => { initialized = true; return backend; } };

    expect(() => new Client({ provider, endpoint: "http://remote.example/api" })).toThrowError(
      expect.objectContaining({ kind: ErrorKind.InvalidArgument, field: "endpoint", operation: "initialize" }),
    );
    expect(initialized).toBe(false);
  });
});

describe("Operation cancellation", () => {
  it("does not invoke the provider for an already-aborted operation", async () => {
    const { client, backend } = fixture();
    const controller = new AbortController();
    controller.abort(new Error("caller stopped work"));

    await expect(client.create({}, { signal: controller.signal })).rejects.toMatchObject({
      kind: ErrorKind.Canceled,
      operation: "create",
    });
    expect(backend.createCalls).toBe(0);
  });

  it("releases client operation bookkeeping after invalid timeout input", async () => {
    const { client, backend } = fixture();

    await expect(client.create({}, { timeoutMs: -1 })).rejects.toMatchObject({
      kind: ErrorKind.InvalidArgument,
      operation: "create",
      field: "timeoutMs",
    });
    await expect(client.close()).resolves.toBeUndefined();
    expect(backend.closeCalls).toBe(1);
  });

  it("enforces a caller timeout when a provider operation does not settle", async () => {
    const { client, backend } = hangingFixture();

    await expect(client.create({}, { timeoutMs: 10 })).rejects.toMatchObject({
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
  it("attaches, reports capabilities, and terminates through the provider", async () => {
    const { client, backend } = fixture();
    const sandbox = await client.attach("sandbox-existing");
    expect(sandbox.supports(LifecycleCapability.Terminate)).toBe(true);
    await expect(sandbox.terminate({ wait: true })).resolves.toEqual({ exitCode: 0 });
    expect(backend.terminated).toBe(true);
  });

  it("rejects lifecycle operations absent from the selected provider", async () => {
    const { client } = fixture();
    const sandbox = await client.attach("sandbox-existing");
    await expect(sandbox.pause()).rejects.toMatchObject({ kind: ErrorKind.Unsupported, operation: "pause" });
  });
});

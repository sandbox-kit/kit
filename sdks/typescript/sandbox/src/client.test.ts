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

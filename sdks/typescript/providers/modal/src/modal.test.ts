import { describe, expect, it } from "vitest";
import {
  EgressMode, ErrorKind, IsolationKind, PolicyMode, PortTransport, SecretInjection, WaitCondition,
} from "@sandbox-kit/sandbox";
import { InvalidError } from "modal";
import { mapClientConfig, mapCreateRequest } from "./mappings.js";
import { captureProviderState } from "./provider.config.gen.js";
import { mapModalError } from "./provider.errors.gen.js";

describe("Modal mappings", () => {
  it("maps token authentication and environment scope", () => {
    expect(mapClientConfig({ auth: { credentials: "tokenPair", tokenPair: { id: "id", secret: "secret" } }, scope: { environment: "dev" }, timeout: 30_000 })).toEqual({ tokenId: "id", tokenSecret: "secret", environment: "dev", timeoutMs: 30_000 });
  });

  it("rejects generated unsupported and malformed client configuration", () => {
    expect(() => mapClientConfig({ endpoint: "https://unsupported.example" })).toThrowError(
      expect.objectContaining({ kind: ErrorKind.Unsupported, field: "endpoint", operation: "initialize" }),
    );
    expect(() => mapClientConfig({
      auth: { credentials: "tokenPair", tokenPair: { id: "", secret: "secret" } },
    })).toThrowError(expect.objectContaining({ kind: ErrorKind.InvalidArgument, field: "auth.tokenPair.id" }));
  });

  it("captures retained configuration without aliasing caller state", () => {
    const config = { scope: { appName: "original", environment: "dev" }, region: "us-east-1" };
    const state = captureProviderState(config);
    config.scope.appName = "changed";

    expect(state).toEqual({ scope: { appName: "original", environment: "dev" }, region: "us-east-1" });
  });

  it("maps the Go creation surface into the native Modal plan", () => {
    const plan = mapCreateRequest({
      name: "example", source: { kind: "image", image: { reference: "alpine:3.21" } },
      runtime: { workingDirectory: "/tmp", entrypoint: ["sleep", "infinity"], pty: true },
      isolation: { kind: IsolationKind.LinuxVm },
      resources: { cpuCores: 0.5, cpuLimitCores: 1, memoryMib: 512n, memoryLimitMib: 1024n },
      placement: { cloud: "aws", regions: ["us-east-1"] },
      network: { egress: EgressMode.Restricted, outboundCidrs: { entries: ["192.0.2.0/24"] }, outboundDomains: { entries: ["example.com"] }, inboundCidrs: { entries: ["198.51.100.0/24"] }, privateNetwork: true, ports: [{ port: 443, transport: PortTransport.Tls }] },
      security: { workloadIdentity: true, secrets: [{ reference: "credential", injection: SecretInjection.Environment }] },
      lifetime: { maximumLifetime: 300_000, idleTerminate: { mode: PolicyMode.After, after: 60_000 } },
      readiness: { check: "command", command: { argv: ["true"] } },
      provisioning: { timeout: 120_000, waitFor: WaitCondition.Ready },
      observability: { verbose: true }, labels: { purpose: "test" }, environment: { NODE_ENV: "test" },
    }, { scope: { appName: "app", environment: "dev" }, region: "us-east-1" });
    expect(plan).toMatchObject({ app: "app", environment: "dev", image: "alpine:3.21", secretNames: ["credential"], waitUntilReady: true, readyTimeoutMs: 120_000, params: { name: "example", runtime: "vm", cpu: 0.5, cpuLimit: 1, memoryMiB: 512, memoryLimitMiB: 1024, cloud: "aws", regions: ["us-east-1"], i6pn: true, encryptedPorts: [443], includeOidcIdentityToken: true, timeoutMs: 300_000, idleTimeoutMs: 60_000, verbose: true } });
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
    expect(() => mapCreateRequest({
      source: { kind: "image", image: { reference: "alpine:3.21" } },
      environment: { "NOT VALID": "value" },
    }, config)).toThrowError(expect.objectContaining({ kind: ErrorKind.InvalidArgument, field: "environment" }));
    expect(() => mapCreateRequest({
      source: { kind: "image", image: { reference: "alpine:3.21" } },
      network: { inboundCidrs: { entries: ["300.1.1.1/24"] } },
    }, config)).toThrowError(expect.objectContaining({ kind: ErrorKind.InvalidArgument, field: "network.inboundCidrs.entries" }));
  });
});

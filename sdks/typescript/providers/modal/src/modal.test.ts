import { describe, expect, it } from "vitest";
import {
  EgressMode, IsolationKind, PolicyMode, PortTransport, SecretInjection, WaitCondition,
} from "@sandbox-kit/sandbox";
import { mapClientConfig, mapCreateRequest } from "./mappings.js";

describe("Modal mappings", () => {
  it("maps token authentication and environment scope", () => {
    expect(mapClientConfig({ auth: { credentials: "tokenPair", tokenPair: { id: "id", secret: "secret" } }, scope: { environment: "dev" }, timeout: 30_000 })).toEqual({ tokenId: "id", tokenSecret: "secret", environment: "dev", timeoutMs: 30_000 });
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
});

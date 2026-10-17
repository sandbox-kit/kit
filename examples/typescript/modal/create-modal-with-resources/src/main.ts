import { modal } from "@sandbox-kit/provider-modal";
import {
  Client,
  EgressMode,
  IsolationKind,
  LifecycleCapability,
  PolicyMode,
  SandboxError,
  WaitCondition,
  type ClientConfig,
  type Sandbox,
  type SandboxSource,
  duration,
} from "@sandbox-kit/sandbox";

function requireEnv(name: string): string {
  const value = process.env[name];
  if (!value) throw new Error(`${name} is required`);
  return value;
}

function clientConfig(): ClientConfig {
  return {
    provider: modal(),
    auth: {
      credentials: { case: "tokenPair", value: {
        id: requireEnv("MODAL_TOKEN_ID"),
        secret: requireEnv("MODAL_TOKEN_SECRET"),
      } },
    },
    scope: {
      appName: requireEnv("MODAL_APP_NAME"),
      environment: process.env.MODAL_ENVIRONMENT || "main",
    },
    timeout: duration(120_000),
  };
}

function assertCapability(
  sandbox: Sandbox,
  capability: (typeof LifecycleCapability)[keyof typeof LifecycleCapability],
): void {
  if (!sandbox.supports(capability)) {
    throw new Error(`sandbox ${sandbox.id} does not support ${capability}`);
  }
}

const imageSource: SandboxSource = {
  kind: "image",
  image: { reference: "alpine:3.21" },
};

await using client = new Client(clientConfig());

const sandbox = await client.create({
  source: imageSource,
  resources: { cpuCores: 0.5, memoryMib: 512n },
  lifetime: { maximumLifetime: 5 * 60_000 },
});
const info = sandbox.info();
console.log(`Modal sandbox created: ${sandbox.id} (requested 0.5 CPU cores, 512 MiB, 5-minute lifetime)`);
console.log(
  `ID origin: ${info.origins?.id ?? "unknown"}; allocated resources reported: ${info.resources !== undefined}`,
);

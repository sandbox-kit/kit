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
      credentials: "tokenPair",
      tokenPair: {
        id: requireEnv("MODAL_TOKEN_ID"),
        secret: requireEnv("MODAL_TOKEN_SECRET"),
      },
    },
    scope: {
      appName: requireEnv("MODAL_APP_NAME"),
      environment: process.env.MODAL_ENVIRONMENT || "main",
    },
    timeout: 120_000,
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
  lifetime: { idleTerminate: { mode: PolicyMode.After, after: 60_000 } },
});
console.log(`Modal sandbox created: ${sandbox.id}`);

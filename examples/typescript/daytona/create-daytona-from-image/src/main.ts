import { daytona } from "@sandbox-kit/provider-daytona";
import {
  Client,
  LifecycleCapability,
  PolicyMode,
  SandboxError,
  type ClientConfig,
  type Sandbox,
} from "@sandbox-kit/sandbox";

function requireEnv(name: string): string {
  const value = process.env[name];
  if (!value) throw new Error(`${name} is required`);
  return value;
}

function clientConfig(): ClientConfig {
  const endpoint = process.env.DAYTONA_API_URL;
  const region = process.env.DAYTONA_TARGET;
  return {
    provider: daytona(),
    auth: { credentials: "apiKey", apiKey: { key: requireEnv("DAYTONA_API_KEY") } },
    timeout: 120_000,
    ...(endpoint ? { endpoint } : {}),
    ...(region ? { region } : {}),
  };
}

function assertCapability(sandbox: Sandbox, capability: (typeof LifecycleCapability)[keyof typeof LifecycleCapability]): void {
  if (!sandbox.supports(capability)) {
    throw new Error(`sandbox ${sandbox.id} does not support ${capability}`);
  }
}

await using client = new Client(clientConfig());

const sandbox = await client.create({
  source: { kind: "image", image: { reference: "ubuntu:22.04" } },
  resources: { cpuCores: 1, memoryMib: 1_024n, diskMib: 3_072n },
});
console.log(`Daytona sandbox created from ubuntu:22.04: ${sandbox.id}`);


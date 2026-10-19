import { daytona } from "@sandbox-kit/daytona";
import {
  Client,
  LifecycleCapability,
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
    auth: { apiKey: { key: requireEnv("DAYTONA_API_KEY") } },
    timeout: 120_000,
    ...(endpoint ? { endpoint } : {}),
    ...(region ? { region } : {}),
  };
}

function assertCapability(
  sandbox: Sandbox,
  capability: LifecycleCapability,
): void {
  if (!sandbox.supports(capability)) {
    throw new Error(`sandbox ${sandbox.id} does not support ${capability}`);
  }
}

await using client = new Client(clientConfig());

const sandbox = await client.get(requireEnv("DAYTONA_SANDBOX_ID"));
assertCapability(sandbox, LifecycleCapability.Start);
// Single options bag: Kit abort deadline + native start timeout (same ms value).
const info = await sandbox.start({ timeout: 60_000 });
console.log(
  `Daytona sandbox started with 60s timeout: ${sandbox.id} (${info.state})`,
);

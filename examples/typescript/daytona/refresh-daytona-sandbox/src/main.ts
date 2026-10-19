import { daytona } from "@sandbox-kit/provider-daytona";
import {
  Client,
  LifecycleCapability,
  duration,
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
    auth: { credentials: { case: "apiKey", value: { key: requireEnv("DAYTONA_API_KEY") } } },
    timeout: duration(120_000),
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
assertCapability(sandbox, LifecycleCapability.Refresh);
const info = await sandbox.refresh();
console.log(`Daytona metadata refreshed: ${sandbox.id} (${info.state})`);

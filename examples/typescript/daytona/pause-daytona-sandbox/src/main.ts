import { daytona } from "@sandbox-kit/daytona";
import {
  Client,
  LifecycleCapability,
  type Sandbox,
} from "@sandbox-kit/sandbox";

function assertCapability(
  sandbox: Sandbox,
  capability: LifecycleCapability,
): void {
  if (!sandbox.supports(capability)) {
    throw new Error(`sandbox ${sandbox.id} does not support ${capability}`);
  }
}

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
  ...(process.env.DAYTONA_API_URL ? { endpoint: process.env.DAYTONA_API_URL } : {}),
  ...(process.env.DAYTONA_TARGET ? { region: process.env.DAYTONA_TARGET } : {}),
});

const sandbox = await client.get(process.env.DAYTONA_SANDBOX_ID!);
assertCapability(sandbox, LifecycleCapability.Pause);
const info = await sandbox.pause();
console.log(`Daytona pause completed: ${sandbox.id} (${info.state}, native ${info.providerState ?? "unknown"})`);

import { blaxel } from "@sandbox-kit/blaxel";
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

const region = process.env.BL_REGION || "us-pdx-1";
const workspace = process.env.BL_WORKSPACE?.trim();

await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: process.env.BL_API_KEY! } },
  ...(workspace ? { scope: { workspace } } : {}),
  region,
  timeout: 120_000,
});

const sandbox = await client.get(process.env.BL_SANDBOX_NAME!);
assertCapability(sandbox, LifecycleCapability.Refresh);
const info = await sandbox.refresh();
console.log(`Blaxel metadata refreshed: ${sandbox.id} (${info.state})`);

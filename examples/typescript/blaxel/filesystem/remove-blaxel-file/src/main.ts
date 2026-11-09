import { blaxel } from "@sandbox-kit/blaxel";
import { Client } from "@sandbox-kit/sandbox";

const region = process.env.BL_REGION || "us-pdx-1";

await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: process.env.BL_API_KEY! } },
  scope: { workspace: process.env.BL_WORKSPACE! },
  region,
  timeout: 120_000,
});

const sandbox = await client.get(process.env.BL_SANDBOX_ID!);
const remote = process.env.BL_REMOTE_PATH!;
const recursive = process.env.BL_RECURSIVE === "true";
await sandbox.filesystem.remove(remote, recursive ? { recursive: true } : {});
console.log(`removed ${remote}`);

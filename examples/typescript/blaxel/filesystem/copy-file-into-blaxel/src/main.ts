import { blaxel } from "@sandbox-kit/blaxel";
import { Client } from "@sandbox-kit/sandbox";

const region = process.env.BL_REGION || "us-pdx-1";
const workspace = process.env.BL_WORKSPACE?.trim();

await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: process.env.BL_API_KEY! } },
  ...(workspace ? { scope: { workspace } } : {}),
  region,
  timeout: 120000,
});

const sandbox = await client.get(process.env.BL_SANDBOX_NAME!);
const remote = process.env.BL_REMOTE_PATH!;
const local = process.env.BL_LOCAL_PATH!;
await sandbox.filesystem.copyFromLocal(local, remote);
console.log(`copied ${local} to ${remote}`);

import { blaxel } from "@sandbox-kit/blaxel";
import { Client } from "@sandbox-kit/sandbox";

const region = process.env.BL_REGION || "us-pdx-1";

await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: process.env.BL_API_KEY! } },
  scope: { workspace: process.env.BL_WORKSPACE! },
  region,
  timeout: 120000,
});

const sandbox = await client.get(process.env.BL_SANDBOX_ID!);
const remote = process.env.BL_REMOTE_PATH!;
const body = new TextEncoder().encode("hello from sandbox-kit\n");
await sandbox.filesystem.writeFile(remote, body);
console.log(`wrote ${remote} (${body.byteLength} bytes)`);

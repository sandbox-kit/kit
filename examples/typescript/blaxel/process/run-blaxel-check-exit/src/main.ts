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

const sandbox = await client.get(process.env.BL_SANDBOX_NAME!);
const result = await sandbox.process.run({ shell: "exit 7" });
if (result.exitCode !== 7) {
  throw new Error(`want exit 7, got ${result.exitCode}`);
}
console.log(`command finished with exit=${result.exitCode}`);

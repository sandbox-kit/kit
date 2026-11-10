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
const result = await sandbox.process.run({
  argv: ["echo", "hello from kit"],
});
console.log(
  `exit=${result.exitCode} stdout=${JSON.stringify(new TextDecoder().decode(result.stdout))} stderr=${JSON.stringify(new TextDecoder().decode(result.stderr))}`,
);

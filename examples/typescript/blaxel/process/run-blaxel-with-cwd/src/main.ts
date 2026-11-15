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
const cwd = process.env.BL_WORKING_DIRECTORY || "/tmp";
const result = await sandbox.process.run({
  argv: ["pwd"],
  workingDirectory: cwd,
});
console.log(
  `exit=${result.exitCode} cwd_out=${JSON.stringify(new TextDecoder().decode(result.stdout))}`,
);

import { blaxel } from "@sandbox-kit/blaxel";
import { Client, PolicyMode } from "@sandbox-kit/sandbox";

const region = process.env.BL_REGION || "us-pdx-1";
const workspace = process.env.BL_WORKSPACE?.trim();

await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: process.env.BL_API_KEY! } },
  ...(workspace ? { scope: { workspace } } : {}),
  region,
  timeout: 120000,
});

const sandbox = await client.create({
  name: "kit-create-blaxel-idle",
  source: { image: { reference: "blaxel/base-image:latest" } },
  lifetime: { idleTerminate: { mode: PolicyMode.After, after: 60_000 } },
});
console.log(`Blaxel sandbox created with 60s idle delete: ${sandbox.id}`);

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

const sandbox = await client.create({
  name: "kit-create-blaxel-ttl",
  source: { image: { reference: "blaxel/base-image:latest" } },
  lifetime: { maximumLifetime: 60 * 60_000 },
});
console.log(`Blaxel sandbox created with 60-minute maximum lifetime: ${sandbox.id}`);

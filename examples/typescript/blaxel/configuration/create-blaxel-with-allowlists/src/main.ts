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

const sandbox = await client.create({
  name: "kit-create-blaxel-allowlists",
  source: { image: { reference: "blaxel/base-image:latest" } },
  network: { outboundDomains: ["example.com"] },
});
console.log(`Blaxel sandbox created with outbound domains: ${sandbox.id}`);

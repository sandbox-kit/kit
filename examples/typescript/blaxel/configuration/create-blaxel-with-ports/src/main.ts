import { blaxel } from "@sandbox-kit/blaxel";
import { Client, PortTransport } from "@sandbox-kit/sandbox";

const region = process.env.BL_REGION || "us-pdx-1";

await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: process.env.BL_API_KEY! } },
  scope: { workspace: process.env.BL_WORKSPACE! },
  region,
  timeout: 120000,
});

const sandbox = await client.create({
  name: "kit-create-blaxel-ports",
  source: { image: { reference: "blaxel/base-image:latest" } },
  network: {
    ports: [
      { port: 3000, transport: PortTransport.Tls },
      { port: 3001, transport: PortTransport.Plain },
    ],
  },
});
console.log(`Blaxel sandbox created with ports: ${sandbox.id}`);

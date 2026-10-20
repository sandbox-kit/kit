import { daytona } from "@sandbox-kit/daytona";
import {
  Client,
  EgressMode,
} from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
  ...(process.env.DAYTONA_API_URL ? { endpoint: process.env.DAYTONA_API_URL } : {}),
  ...(process.env.DAYTONA_TARGET ? { region: process.env.DAYTONA_TARGET } : {}),
});

const sandbox = await client.create({
  source: { image: { reference: "ubuntu:22.04" } },
  network: {
    egress: EgressMode.Restricted,
    outboundCidrs: ["192.0.2.0/24"],
    outboundDomains: ["example.com"],
  },
});
console.log(`Daytona sandbox created with restricted egress: ${sandbox.id}`);

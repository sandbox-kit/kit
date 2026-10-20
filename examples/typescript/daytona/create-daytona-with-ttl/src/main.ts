import { daytona } from "@sandbox-kit/daytona";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
  ...(process.env.DAYTONA_API_URL ? { endpoint: process.env.DAYTONA_API_URL } : {}),
  ...(process.env.DAYTONA_TARGET ? { region: process.env.DAYTONA_TARGET } : {}),
});

const sandbox = await client.create({
  source: { image: { reference: "ubuntu:22.04" } },
  lifetime: { maximumLifetime: 60 * 60_000 },
});
console.log(
  `Daytona sandbox created with 60-minute maximum lifetime: ${sandbox.id}`,
);

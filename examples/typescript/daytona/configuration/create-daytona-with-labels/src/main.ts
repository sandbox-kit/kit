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
  labels: { team: "platform", purpose: "kit-example" },
  environment: { NODE_ENV: "development" },
});
console.log(
  `Daytona sandbox created with labels and environment: ${sandbox.id}`,
);

import { daytona } from "@sandbox-kit/daytona";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
  ...(process.env.DAYTONA_API_URL ? { endpoint: process.env.DAYTONA_API_URL } : {}),
  ...(process.env.DAYTONA_TARGET ? { region: process.env.DAYTONA_TARGET } : {}),
});

const parent = await client.create();
const child = await client.create({ network: { linkedSandbox: parent.id }, lifetime: { ephemeral: true } });
console.log(`Daytona parent ${parent.id}; linked ephemeral child ${child.id}`);

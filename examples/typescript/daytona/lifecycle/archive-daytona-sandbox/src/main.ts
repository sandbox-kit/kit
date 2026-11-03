import { daytona } from "@sandbox-kit/daytona";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
  ...(process.env.DAYTONA_API_URL ? { endpoint: process.env.DAYTONA_API_URL } : {}),
  ...(process.env.DAYTONA_TARGET ? { region: process.env.DAYTONA_TARGET } : {}),
});

const sandbox = await client.get(process.env.DAYTONA_SANDBOX_ID!);
const info = await sandbox.archive();
console.log(`Daytona archive requested: ${sandbox.id} (${info.state})`);


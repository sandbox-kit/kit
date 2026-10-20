import { daytona } from "@sandbox-kit/daytona";
import {
  Client,
  PolicyMode,
} from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
  ...(process.env.DAYTONA_API_URL ? { endpoint: process.env.DAYTONA_API_URL } : {}),
  ...(process.env.DAYTONA_TARGET ? { region: process.env.DAYTONA_TARGET } : {}),
});

const sandbox = await client.create({
  lifetime: {
    idlePause: { mode: PolicyMode.Disabled },
    stoppedDelete: { mode: PolicyMode.After, after: 10 * 60_000 },
  },
});
const info = sandbox.info();
console.log(`Daytona sandbox created: ${sandbox.id} (auto-pause disabled, deletion 10 minutes after stopping)`);
console.log(`ID origin: ${info.origins?.id ?? "unknown"}; resource origin: ${info.origins?.resources ?? "unknown"}`);

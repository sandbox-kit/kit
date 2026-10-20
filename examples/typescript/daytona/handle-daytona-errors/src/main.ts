import { daytona } from "@sandbox-kit/daytona";
import {
  Client,
  SandboxError,
} from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
  ...(process.env.DAYTONA_API_URL ? { endpoint: process.env.DAYTONA_API_URL } : {}),
  ...(process.env.DAYTONA_TARGET ? { region: process.env.DAYTONA_TARGET } : {}),
});

try {
  await client.create({ resources: { cpuCores: -1 } });
  console.log("Sandbox created.");
} catch (error) {
  if (error instanceof SandboxError) {
    console.log(`${error.kind}: ${error.message} (field=${error.field || "none"}, operation=${error.operation})`);
  } else {
    throw error;
  }
}

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
const cwd = process.env.DAYTONA_WORKING_DIRECTORY || "/tmp";
const result = await sandbox.process.run({
  argv: ["pwd"],
  workingDirectory: cwd,
});
console.log(`exit=${result.exitCode} cwd_out=${JSON.stringify(new TextDecoder().decode(result.stdout))}`);

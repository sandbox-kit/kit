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
const result = await sandbox.process.run({
  argv: ["echo", "hello from kit"],
});
console.log(
  `exit=${result.exitCode} stdout=${JSON.stringify(new TextDecoder().decode(result.stdout))} stderr=${JSON.stringify(new TextDecoder().decode(result.stderr))}`,
);

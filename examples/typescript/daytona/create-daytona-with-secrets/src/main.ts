import { daytona } from "@sandbox-kit/daytona";
import {
  Client,
  SecretInjection,
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
  security: {
    secrets: [
      {
        reference: process.env.DAYTONA_SECRET_NAME!,
        environmentVariable: "GITHUB_TOKEN",
        injection: SecretInjection.EgressPlaceholder,
      },
    ],
  },
});
console.log(`Daytona sandbox created with injected secret: ${sandbox.id}`);

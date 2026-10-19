import { daytona } from "@sandbox-kit/provider-daytona";
import {
  Client,
  duration,
  type ClientConfig,
  type CreateOptionsInput,
} from "@sandbox-kit/sandbox";

function requireEnv(name: string): string {
  const value = process.env[name];
  if (!value) throw new Error(`${name} is required`);
  return value;
}

function clientConfig(): ClientConfig {
  const endpoint = process.env.DAYTONA_API_URL;
  const region = process.env.DAYTONA_TARGET;
  return {
    provider: daytona(),
    auth: { credentials: { case: "apiKey", value: { key: requireEnv("DAYTONA_API_KEY") } } },
    timeout: duration(120_000),
    ...(endpoint ? { endpoint } : {}),
    ...(region ? { region } : {}),
  };
}

await using client = new Client(clientConfig());

const createOptions = {
  source: { kind: { case: "snapshot", value: { reference: "daytona-small" } } },
  runtime: { language: "python", user: "daytona" },
} satisfies CreateOptionsInput;

const sandbox = await client.create(createOptions);
console.log(`Daytona Python sandbox created: ${sandbox.id} (${sandbox.info().isolation?.kind ?? "default"})`);

import { daytona } from "@sandbox-kit/daytona";
import {
  Client,
  type ClientConfig,
  type CreateOptions,
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
    auth: { apiKey: { key: requireEnv("DAYTONA_API_KEY") } },
    timeout: 120_000,
    ...(endpoint ? { endpoint } : {}),
    ...(region ? { region } : {}),
  };
}

await using client = new Client(clientConfig());

const createOptions = {
  source: { image: { reference: "ubuntu:22.04" } },
  isolation: { nestedVirtualization: true},
} satisfies CreateOptions;

const sandbox = await client.create(createOptions);
console.log(`Daytona sandbox created with nested virtualization: ${sandbox.id}`);

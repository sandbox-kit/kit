import { daytona } from "@sandbox-kit/daytona";
import {
  Client,
  PolicyMode,
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
    auth: { apiKey: { key: requireEnv("DAYTONA_API_KEY")  } },
    timeout: 120_000,
    ...(endpoint ? { endpoint } : {}),
    ...(region ? { region } : {}),
  };
}

await using client = new Client(clientConfig());

const createOptions = {
  lifetime: {
    idlePause: { mode: PolicyMode.Disabled },
    stoppedDelete: { mode: PolicyMode.After, after: 10 * 60_000 },
  },
} satisfies CreateOptions;

const sandbox = await client.create(createOptions);
const info = sandbox.info();
console.log(`Daytona sandbox created: ${sandbox.id} (auto-pause disabled, deletion 10 minutes after stopping)`);
console.log(`ID origin: ${info.origins?.id ?? "unknown"}; resource origin: ${info.origins?.resources ?? "unknown"}`);

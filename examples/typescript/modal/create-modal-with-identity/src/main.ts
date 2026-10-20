import { modal } from "@sandbox-kit/modal";
import {
  Client,
  type ClientConfig,
  type SandboxSourceInput,
  type CreateOptions,
} from "@sandbox-kit/sandbox";

function requireEnv(name: string): string {
  const value = process.env[name];
  if (!value) throw new Error(`${name} is required`);
  return value;
}

function clientConfig(): ClientConfig {
  return {
    provider: modal(),
    auth: {
      tokenPair: {
        id: requireEnv("MODAL_TOKEN_ID"),
        secret: requireEnv("MODAL_TOKEN_SECRET"),
      },
    },
    scope: {
      appName: requireEnv("MODAL_APP_NAME"),
      environment: process.env.MODAL_ENVIRONMENT || "main",
    },
    timeout: 120_000,
  };
}

const imageSource: SandboxSourceInput = {
  image: { reference: "alpine:3.21" },
};

await using client = new Client(clientConfig());

const createOptions = {
  source: imageSource,
  security: { workloadIdentity: true },
} satisfies CreateOptions;

const sandbox = await client.create(createOptions);
console.log(`Modal sandbox created: ${sandbox.id}`);

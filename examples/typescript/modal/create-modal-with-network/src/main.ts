import { modal } from "@sandbox-kit/provider-modal";
import {
  Client,
  duration,
  type ClientConfig,
  type SandboxSourceInput,
  type CreateOptionsInput,
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
      credentials: { case: "tokenPair", value: {
        id: requireEnv("MODAL_TOKEN_ID"),
        secret: requireEnv("MODAL_TOKEN_SECRET"),
      } },
    },
    scope: {
      appName: requireEnv("MODAL_APP_NAME"),
      environment: process.env.MODAL_ENVIRONMENT || "main",
    },
    timeout: duration(120_000),
  };
}

const imageSource: SandboxSourceInput = {
  kind: { case: "image", value: { reference: "alpine:3.21" } },
};

await using client = new Client(clientConfig());

const createOptions = {
  source: imageSource,
  network: { proxyReference: requireEnv("MODAL_PROXY_NAME") },
} satisfies CreateOptionsInput;

const sandbox = await client.create(createOptions);
console.log(`Modal sandbox created: ${sandbox.id}`);

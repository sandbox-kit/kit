import { modal } from "@sandbox-kit/provider-modal";
import { Client, type ClientConfig } from "@sandbox-kit/sandbox";

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

await using client = new Client(clientConfig());

const sandbox = await client.attach(requireEnv("MODAL_SANDBOX_ID"));
console.log(
  `Modal handle attached: id=${sandbox.id} capabilities=${sandbox.capabilities().join(",")}`,
);

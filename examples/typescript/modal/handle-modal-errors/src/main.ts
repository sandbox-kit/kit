import { modal } from "@sandbox-kit/provider-modal";
import {
  Client,
  SandboxError,
  type ClientConfig,
  type SandboxSourceInput,
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

try {
  await client.create({ source: imageSource, resources: { diskMib: 1_024n } });
  console.log("Sandbox created.");
} catch (error) {
  if (error instanceof SandboxError) {
    console.log(`${error.kind}: ${error.message} (field=${error.field || "none"}, operation=${error.operation})`);
  } else {
    throw error;
  }
}

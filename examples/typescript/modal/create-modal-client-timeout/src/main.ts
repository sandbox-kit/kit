import { modal } from "@sandbox-kit/modal";
import {
  Client,
  WaitCondition,
  type ClientConfig,
  type CreateOptions,
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
    // Kit injects this into provisioning.timeout for Modal ready-wait.
    timeout: 90_000,
  };
}

const imageSource: SandboxSourceInput = {
  image: { reference: "alpine:3.21" },
};

await using client = new Client(clientConfig());

const createOptions = {
  source: imageSource,
  runtime: { entrypoint: ["sleep", "300"] },
  readiness: { command: { argv: ["true"] } },
  provisioning: { waitFor: WaitCondition.Ready },
} satisfies CreateOptions;

const sandbox = await client.create(createOptions);
console.log(`Modal sandbox ready using client timeout for wait: ${sandbox.id}`);

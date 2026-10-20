import { modal } from "@sandbox-kit/modal";
import {
  Client,
  EgressMode,
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
    timeout: 120_000,
  };
}

const imageSource: SandboxSourceInput = {
  image: { reference: "alpine:3.21" },
};

await using client = new Client(clientConfig());

const createOptions = {
  source: imageSource,
  network: {
    egress: EgressMode.Restricted,
    outboundCidrs: ["192.0.2.0/24"],
    outboundDomains: ["example.com"],
    privateNetwork: true,
  },
} satisfies CreateOptions;

const sandbox = await client.create(createOptions);
console.log(
  `Modal sandbox created with restricted outbound allowlists: ${sandbox.id}`,
);

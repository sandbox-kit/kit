import { modal } from "@sandbox-kit/modal";
import {
  Client,
  SecretInjection,
} from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: modal(),
  auth: {
    tokenPair: {
      id: process.env.MODAL_TOKEN_ID!,
      secret: process.env.MODAL_TOKEN_SECRET!,
    },
  },
  scope: {
    appName: process.env.MODAL_APP_NAME!,
    environment: process.env.MODAL_ENVIRONMENT || "main",
  },
  timeout: 120_000,
});

const sandbox = await client.create({
  source: { image: { reference: "alpine:3.21" } },
  security: {
    secrets: [
      {
        reference: process.env.MODAL_SECRET_NAME!,
        injection: SecretInjection.Environment,
      },
    ],
  },
});
console.log(`Modal sandbox created with secret: ${sandbox.id}`);

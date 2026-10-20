import { modal } from "@sandbox-kit/modal";
import { Client } from "@sandbox-kit/sandbox";

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

const sandbox = await client.attach(process.env.MODAL_SANDBOX_ID!);
console.log(
  `Modal handle attached: id=${sandbox.id} capabilities=${sandbox.capabilities().join(",")}`,
);

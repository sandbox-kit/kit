import { modal } from "@sandbox-kit/modal";
import {
  Client,
  WaitCondition,
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
  // Kit injects this into provisioning.timeout for Modal ready-wait.
  timeout: 90_000,
});

const sandbox = await client.create({
  source: { image: { reference: "alpine:3.21" } },
  runtime: { entrypoint: ["sleep", "300"] },
  readiness: { command: { argv: ["true"] } },
  provisioning: { waitFor: WaitCondition.Ready },
});
console.log(`Modal sandbox ready using client timeout for wait: ${sandbox.id}`);

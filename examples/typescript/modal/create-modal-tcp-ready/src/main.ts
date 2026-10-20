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
  timeout: 120_000,
});

const sandbox = await client.create({
  source: { image: { reference: "alpine:3.21" } },
  runtime: {
    entrypoint: ["sh", "-c", "nc -l -p 8080"],
  },
  readiness: { tcpPort: 8080 },
  provisioning: { timeout: 120_000, waitFor: WaitCondition.Ready },
});
console.log(`Modal sandbox ready on TCP 8080: ${sandbox.id}`);

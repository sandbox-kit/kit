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
const remote = process.env.MODAL_REMOTE_PATH!;
const body = new TextEncoder().encode("hello from sandbox-kit\n");
await sandbox.filesystem.writeFile(remote, body);
console.log(`wrote ${remote} (${body.byteLength} bytes)`);

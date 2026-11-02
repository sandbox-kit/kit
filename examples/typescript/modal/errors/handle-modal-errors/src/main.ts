import { modal } from "@sandbox-kit/modal";
import {
  Client,
  SandboxError,
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

try {
  await client.create({ source: { image: { reference: "alpine:3.21" } }, resources: { diskMib: 1_024n } });
  console.log("Sandbox created.");
} catch (error) {
  if (error instanceof SandboxError) {
    console.log(`${error.kind}: ${error.message} (field=${error.field || "none"}, operation=${error.operation})`);
  } else {
    throw error;
  }
}

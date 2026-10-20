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

const sandbox = await client.create({
  source: { image: { reference: "alpine:3.21" } },
  resources: { cpuCores: 0.5, memoryMib: 512n },
  lifetime: { maximumLifetime: 5 * 60_000 },
});
const info = sandbox.info();
console.log(`Modal sandbox created: ${sandbox.id} (requested 0.5 CPU cores, 512 MiB, 5-minute lifetime)`);
console.log(
  `ID origin: ${info.origins?.id ?? "unknown"}; allocated resources reported: ${info.resources !== undefined}`,
);

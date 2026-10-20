import { modal } from "@sandbox-kit/modal";
import {
  Client,
  EgressMode,
  IsolationKind,
  PolicyMode,
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
    environment: process.env.MODAL_ENVIRONMENT ?? "main",
  },
  timeout: 120_000,
});

const sandbox = await client.create({
  source: { image: { reference: "alpine:3.21" } },
  runtime: {
    workingDirectory: "/tmp",
    entrypoint: ["sleep", "infinity"],
    pty: true,
  },
  isolation: { kind: IsolationKind.Container },
  resources: { cpuCores: 0.5, memoryMib: 512n },
  network: { egress: EgressMode.AllowAll },
  lifetime: {
    maximumLifetime: 5 * 60_000,
    idleTerminate: { mode: PolicyMode.After, after: 60_000 },
  },
  labels: { example: "create-modal-sandbox" },
});

console.log(`Sandbox created: ${sandbox.id}`);

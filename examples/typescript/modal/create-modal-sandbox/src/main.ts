import { modal } from "@sandbox-kit/provider-modal";
import { Client, EgressMode, IsolationKind, PolicyMode, duration } from "@sandbox-kit/sandbox";

const tokenId = process.env.MODAL_TOKEN_ID;
const tokenSecret = process.env.MODAL_TOKEN_SECRET;
const appName = process.env.MODAL_APP_NAME;
if (!tokenId || !tokenSecret || !appName)
  throw new Error("MODAL_TOKEN_ID, MODAL_TOKEN_SECRET, and MODAL_APP_NAME are required");

await using client = new Client({
  provider: modal(),
  auth: { credentials: { case: "tokenPair", value: { id: tokenId, secret: tokenSecret } } },
  scope: { appName, environment: process.env.MODAL_ENVIRONMENT ?? "main" },
  timeout: duration(120_000),
});

const sandbox = await client.create({
  source: { kind: "image", image: { reference: "alpine:3.21" } },
  runtime: { workingDirectory: "/tmp", entrypoint: ["sleep", "infinity"], pty: true },
  isolation: { kind: IsolationKind.Container },
  resources: { cpuCores: 0.5, memoryMib: 512n },
  network: { egress: EgressMode.AllowAll },
  lifetime: { maximumLifetime: 5 * 60_000, idleTerminate: { mode: PolicyMode.After, after: 60_000 } },
  labels: { example: "create-modal-sandbox" },
});

console.log(`Sandbox created: ${sandbox.id}`);

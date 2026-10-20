import { modal } from "@sandbox-kit/modal";
import {
  Client,
  LifecycleCapability,
  type Sandbox,
} from "@sandbox-kit/sandbox";

function assertCapability(
  sandbox: Sandbox,
  capability: LifecycleCapability,
): void {
  if (!sandbox.supports(capability)) {
    throw new Error(`sandbox ${sandbox.id} does not support ${capability}`);
  }
}

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
assertCapability(sandbox, LifecycleCapability.Terminate);
const exit = await sandbox.terminate({ wait: true, timeout: 60_000 });
console.log(
  `Modal sandbox terminated with 60s wait timeout: ${sandbox.id} (exit ${exit.exitCode ?? 0})`,
);

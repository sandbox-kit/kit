import { daytona } from "@sandbox-kit/provider-daytona";
import {
  Client,
  EgressMode,
  IsolationKind,
  PolicyMode,
} from "@sandbox-kit/sandbox";

async function main(): Promise<void> {
  await using client = new Client({
    provider: daytona(),
    timeout: 120_000,
  });

  const sandbox = await client.create({
    source: {
      kind: "image",
      image: { reference: "debian:12.9" },
    },
    runtime: { language: "typescript" },
    isolation: {
      kind: IsolationKind.Default,
      nestedVirtualization: false,
    },
    resources: {
      cpuCores: 2,
      memoryMib: 4_096n,
      diskMib: 10_240n,
    },
    environment: {
      NODE_ENV: "development",
    },
    labels: {
      example: "create-daytona-sandbox",
    },
    network: {
      egress: EgressMode.AllowAll,
    },
    security: {
      publicAccess: false,
    },
    lifetime: {
      idleStop: {
        mode: PolicyMode.After,
        after: 15 * 60_000,
      },
      stoppedDelete: {
        mode: PolicyMode.After,
        after: 24 * 60 * 60_000,
      },
    },
  });

  console.log(`Sandbox created: ${sandbox.id}`);
  console.log(`Provider: ${sandbox.providerName}`);
}

await main();

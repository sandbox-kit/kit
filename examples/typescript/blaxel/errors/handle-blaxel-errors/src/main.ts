import { blaxel } from "@sandbox-kit/blaxel";
import { Client, SandboxError } from "@sandbox-kit/sandbox";

const region = process.env.BL_REGION || "us-pdx-1";

await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: process.env.BL_API_KEY! } },
  scope: { workspace: process.env.BL_WORKSPACE! },
  region,
  timeout: 120000,
});

try {
  await client.create({
    name: "kit-handle-blaxel-errors",
    source: { image: { reference: "blaxel/base-image:latest" } },
    resources: { cpuCores: 2 },
  });
  console.log("Sandbox created.");
} catch (error) {
  if (error instanceof SandboxError) {
    console.log(
      `${error.kind}: ${error.message} (field=${error.field || "none"}, operation=${error.operation})`,
    );
  } else {
    throw error;
  }
}

import { modal } from "@sandbox-kit/modal";
import {
  Client,
  PortTransport,
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
  network: {
    ports: [
      { port: 443, transport: PortTransport.Tls },
      { port: 8080, transport: PortTransport.Plain },
    ],
  },
});
console.log(`Modal sandbox created with exposed ports: ${sandbox.id}`);

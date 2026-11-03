import { modal } from "@sandbox-kit/modal";
import {
  Client,
  EgressMode,
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
    egress: EgressMode.Restricted,
    outboundCidrs: ["192.0.2.0/24"],
    outboundDomains: ["example.com"],
    privateNetwork: true,
  },
});
console.log(
  `Modal sandbox created with restricted outbound allowlists: ${sandbox.id}`,
);

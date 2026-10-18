# Modal provider for Sandbox Kit

The provider adapts the official `modal@0.11.0` SDK to the Sandbox Kit
TypeScript contracts. Checks and portable field mappings are handwritten in
`provider_checks.ts`, `provider_mappings.ts`, and `provider_config.ts`, kept in
sync with `specs/providers/modal.yaml`. `src/modal.ts` owns SDK call ordering;
`src/mappings.ts` owns remaining composition that those helpers do not cover.

```ts
import { modal } from "@sandbox-kit/provider-modal";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: modal(),
  auth: {
    credentials: { case: "tokenPair", value: { id: "...", secret: "..." } },
  },
  scope: { appName: "my-app", environment: "main" },
});
```

Modal supports image creation, attach, terminate, and local handle detachment.
The creation adapter covers runtime command/workdir/PTY, gVisor or VM isolation,
CPU and memory reservations and limits, placement, network policy and ports,
secrets, OIDC workload identity, readiness, timeouts, labels, and verbose mode.
Unsupported shared intent is rejected before contacting Modal.

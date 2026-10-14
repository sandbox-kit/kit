# Modal provider for Sandbox Kit

The provider adapts the official `modal@0.11.0` SDK to the generated Sandbox Kit
TypeScript contracts. Configuration and portable field mappings are generated
from protobuf plus `specs/providers/modal.yaml`; `src/modal.ts` contains SDK call
ordering and `src/mappings.ts` contains only the remaining composition that
cannot be expressed as portable field mappings.

```ts
import { modal } from "@sandbox-kit/provider-modal";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: modal(),
  auth: { credentials: "tokenPair", tokenPair: { id: "...", secret: "..." } },
  scope: { appName: "my-app", environment: "main" },
});
```

Modal supports image creation, attach, terminate, and local handle detachment.
The creation adapter covers runtime command/workdir/PTY, gVisor or VM isolation,
CPU and memory reservations and limits, placement, network policy and ports,
secrets, OIDC workload identity, readiness, timeouts, labels, and verbose mode.
Unsupported shared intent is rejected before contacting Modal.

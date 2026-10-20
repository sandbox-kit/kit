# Modal provider for Sandbox Kit

Independent package over `modal@0.11.0`, published separately from the core
`@sandbox-kit/sandbox` client. Use it through that shared API.

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/modal
```

## Usage

```ts
import { modal } from "@sandbox-kit/modal";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: modal(),
  auth: {
    tokenPair: { id: "...", secret: "..." },
  },
  scope: { appName: "my-app", environment: "main" },
});

const sandbox = await client.create({
  source: { image: { reference: "alpine:3.21" } },
  provisioning: { timeout: 120_000 },
});
console.log(sandbox.id);
```

Closing the client releases local SDK resources only. The Modal app must already
exist.

## Examples

| Example    | Project                                                                                                    |
| ---------- | ---------------------------------------------------------------------------------------------------------- |
| Creation   | [create-modal-sandbox](../../../../examples/typescript/modal/create-modal-sandbox/README.md)               |
| Resources  | [create-modal-with-resources](../../../../examples/typescript/modal/create-modal-with-resources/README.md) |
| Attach     | [attach-modal-sandbox](../../../../examples/typescript/modal/attach-modal-sandbox/README.md)               |
| Full index | [TypeScript examples](../../../../examples/typescript/README.md)                                           |

## Capabilities

Modal supports image creation, attach, terminate, and local handle detachment.
The creation adapter covers runtime command/workdir/PTY, gVisor or VM isolation,
CPU and memory reservations and limits, placement, network policy and ports,
secrets, OIDC workload identity, readiness, timeouts, labels, and verbose mode.
Unsupported shared intent is rejected before contacting Modal.

## Implementation notes

Checks and portable field mappings are handwritten in `checks/provider_checks.ts`,
`mappings/provider_mappings.ts`, and `config/provider_config.ts`, kept in sync with
`specs/providers/modal.yaml`. `src/modal.ts` owns SDK call ordering;
`src/mappings.ts` owns remaining composition those helpers do not cover.

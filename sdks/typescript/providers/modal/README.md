# `@sandbox-kit/modal`

Modal provider for [`@sandbox-kit/sandbox`](https://www.npmjs.com/package/@sandbox-kit/sandbox).
Use the shared client API; this package supplies the `modal()` provider.

## Install

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/modal
```

Requires Node.js 20, 22, or 24.15+, and a Modal token pair.

## Quick start

```ts
import { modal } from "@sandbox-kit/modal";
import { Client } from "@sandbox-kit/sandbox";

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
  resources: { cpuCores: 0.5, memoryMib: 512n },
  lifetime: { maximumLifetime: 5 * 60_000 },
});
console.log(sandbox.id);
```

Closing the client releases local SDK resources only. Reclaim sandboxes with
idle/max lifetime policies or `terminate`.

## Configuration

| Setting        | Notes                                                      |
| -------------- | ---------------------------------------------------------- |
| Token pair     | `auth: { tokenPair: { id, secret } }`                      |
| OAuth refresh  | Supported with client ID and exactly one secret or JWT key |
| Omitted `auth` | Modal SDK environment / profile credentials                |
| App            | `scope.appName` — required                                 |
| Environment    | `scope.environment` (default often `main`)                 |
| Region         | `region` — default placement preference                    |
| Endpoint       | Not supported                                              |

## What you can do

**Create** from a registry image with runtime (command, workdir, PTY), gVisor or
Linux VM isolation, CPU/memory, placement, network policy and ports, secrets,
OIDC workload identity, readiness probes, timeouts, labels, and verbose mode.

**Lifecycle** on the handle: `attach`, `terminate`, and local `close`.
Provider-verified `get` / `refresh`, start, stop, pause, archive, and delete are
unsupported.

When you set resources, Modal requires at least 0.125 CPU cores and 128 MiB
memory. Omitted values keep Modal defaults. Disk size cannot be raised at create.

## Examples

| Example    | Project                                                                                                                           |
| ---------- | --------------------------------------------------------------------------------------------------------------------------------- |
| Creation   | [create-modal-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/create-modal-sandbox)               |
| Resources  | [create-modal-with-resources](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/create-modal-with-resources) |
| Attach     | [attach-modal-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/attach-modal-sandbox)               |
| Full index | [TypeScript examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript)                                           |

## License

MIT

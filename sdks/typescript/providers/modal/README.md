# `@sandbox-kit/modal`

Modal provider for [`@sandbox-kit/sandbox`](https://www.npmjs.com/package/@sandbox-kit/sandbox).
Use the shared client interface; this package supplies the `modal()` provider.

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

**Lifecycle** on the handle: `attach`, `terminate`, and `close`. `close` does
not terminate the sandbox.
Provider-verified `get` / `refresh`, start, stop, pause, archive, and delete are
unsupported.

**Files** through `sandbox.filesystem`: `readFile`, `writeFile`,
`copyFromLocal`, `copyToLocal`, `makeDirectory`, `listDirectory`, `stat`, and
`remove`. Paths in the sandbox are absolute. A file timeout is the abort
signal. `1_500` ends the call at 1.5 seconds. `parents` defaults to true:
missing parents are created, and an existing directory is accepted.
`parents: false` requires the parent to exist. A missing file passed to
`copyFromLocal` is not found before Modal is called, and the error field is
`local_path`. A local directory is an invalid argument with the same field.
`copyToLocal` creates a missing local parent directory.

**Commands** through `sandbox.process.run`. Set exactly one of `argv` or
`shell`. Kit maps `shell` to `bash -c` for Modal. A non-zero exit code is not
an SDK error—check `exitCode`. Modal collects both stdout and stderr. An
omitted or zero command timeout adds no Modal command limit: Kit omits
`timeoutMs` and must not send `timeoutMs: 0` (Modal rejects it). Provider
command timeouts must be whole seconds (`2000` is two seconds; `1500` is
rejected).

When you set resources, Modal requires at least 0.125 CPU cores and 128 MiB
memory. Omitted values keep Modal defaults. Disk size cannot be raised at create.

## Examples

| Example       | Project                                                                                                                                         |
| ------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| Creation      | [create-modal-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/configuration/create-modal-sandbox)               |
| Run a command | [run-modal-command](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/process/run-modal-command)                           |
| Process       | [process examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/process)                                              |
| Read a file   | [read-modal-file](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/filesystem/read-modal-file)                            |
| Filesystem    | [filesystem examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/filesystem)                                        |
| Resources     | [create-modal-with-resources](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/configuration/create-modal-with-resources) |
| Attach        | [attach-modal-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/modal/lifecycle/attach-modal-sandbox)                   |
| Full index    | [TypeScript examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript)                                                         |

## License

MIT

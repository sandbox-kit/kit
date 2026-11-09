# `@sandbox-kit/blaxel`

Blaxel provider for [`@sandbox-kit/sandbox`](https://www.npmjs.com/package/@sandbox-kit/sandbox).
Use the shared client interface; this package supplies the `blaxel()` provider.

## Install

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/blaxel
```

Requires Node.js 20, 22, or 24.15+, a Blaxel API key, and a workspace name.

## Quick start

```ts
import { blaxel } from "@sandbox-kit/blaxel";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: process.env.BL_API_KEY! } },
  scope: { workspace: process.env.BL_WORKSPACE! },
  region: process.env.BL_REGION ?? "us-pdx-1",
  timeout: 120_000,
});

const sandbox = await client.create({
  name: "kit-example",
  source: { image: { reference: "blaxel/base-image:latest" } },
  resources: { memoryMib: 4096n },
  lifetime: { maximumLifetime: 24 * 60 * 60_000 },
});
console.log(sandbox.id);
```

Closing the client releases local SDK resources only. Reclaim sandboxes with
TTL / idle policies or `delete`.

## Configuration

| Setting        | Notes                                                         |
| -------------- | ------------------------------------------------------------- |
| API key        | `auth: { apiKey: { key } }`                                   |
| Omitted `auth` | Blaxel SDK environment / CLI credentials                      |
| Workspace      | `scope.workspace` — Blaxel workspace name                     |
| Region         | `region` — default create placement                           |
| Endpoint       | Not supported                                                 |

Token-pair, bearer, and OAuth client modes are rejected. Scope fields other
than `workspace` are rejected.

## What you can do

**Create** requires a nonblank `name` and an image source. Optional memory
(MiB; Blaxel CPU follows memory), labels, environment variables, ports, outbound
domains, max-age TTL, idle-delete policy, and a create timeout (1–50 whole
seconds). Creation calls native `SandboxInstance.create`, not
`createIfNotExists`.

**Lifecycle** on the handle: `get` / `refresh` (lookup by name), `delete`, and
`close`. `close` does not delete the sandbox. Start, stop, pause, archive,
terminate, and attach are unsupported.

**Files** through `sandbox.filesystem`: `readFile`, `writeFile`,
`copyFromLocal`, `copyToLocal`, `makeDirectory`, `listDirectory`, `stat`, and
`remove`. Paths in the sandbox are absolute. `parents` defaults to true;
`parents: false` is rejected. A missing file passed to `copyFromLocal` is not
found before Blaxel is called, and the error field is `local_path`. A local
directory is an invalid argument with the same field. `copyToLocal` creates a
missing local parent directory.

**Commands** through `sandbox.process.run`. Set exactly one of `argv` or
`shell`. `argv` is sent as a POSIX-quoted shell string. Runs are blocking
(`waitForCompletion: true`). A non-zero exit code is not an SDK error—check
`exitCode`. Blaxel returns stdout and stderr. Omit or zero `timeout` for no
command limit. Command timeouts must be whole seconds (`2000` is two seconds;
`1500` is rejected). Aborting the signal rejects the `await`; the remote
command may keep running.

When you set memory, Blaxel requires at least 256 MiB. Omitted memory keeps
the Blaxel default. CPU, CPU limit, memory limit, and disk overrides are
rejected.

## License

MIT

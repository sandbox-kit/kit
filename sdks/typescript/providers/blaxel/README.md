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
console.log(sandbox.id); // same as create name; use with get / BL_SANDBOX_NAME
```

Closing the client releases local SDK resources only. Reclaim sandboxes with
TTL / idle policies or `delete`.

## Configuration

| Setting        | Notes                                     |
| -------------- | ----------------------------------------- |
| API key        | `auth: { apiKey: { key } }`               |
| Omitted `auth` | Blaxel SDK environment / CLI credentials  |
| Workspace      | `scope.workspace` — Blaxel workspace name |
| Region         | `region` — required for create            |
| Endpoint       | Not supported                             |

Token-pair, bearer, and OAuth client modes are rejected. Scope fields other
than `workspace` are rejected.

## What you can do

**Create** requires a Blaxel-valid `name` (lowercase alphanumeric with single
hyphens, at most 49 characters), an image source, and client `region`.
Optional `memoryMib` maps 1:1 onto Blaxel's memory field (Blaxel labels that
integer MB; Kit does not convert; Blaxel CPU follows memory), labels,
environment variables, ports, outbound domains, max-age TTL (at least 5m), and
idle-delete policy (at least 1m; explicit disable is rejected). Client /
create `timeout` is the Kit deadline. `@blaxel/core`'s native creation timeout
(1–50s) is not forwarded. `provisioning.waitFor` accepts only `Default` and
`Scheduled`; Ready (and Started / Submitted) are unsupported. Creation calls
native `SandboxInstance.create`, not `createIfNotExists`. Ports
80, 443, and 8080 are rejected (Blaxel-reserved). Default port transport leaves
the protocol to Blaxel (HTTP); plain maps to TCP; TLS maps to TLS. HTTP/2 over
TLS is rejected. Kit has no portable volume fields; Blaxel volumes are outside
create.

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
`shell`. `argv` is sent as a POSIX-quoted shell string. `run` waits until the
command ends. A non-zero exit code is not an SDK error—check `exitCode`.
Blaxel returns stdout and stderr. Omit or zero `timeout` for no process
auto-kill; Kit then polls until the command ends or the signal aborts. A
positive timeout of at most 60 seconds uses Blaxel's blocking completion
request. A longer timeout sets the process auto-kill and polls. Command
timeouts must be whole seconds (`2000` is two seconds; `1500` is rejected).
Aborting the signal rejects the `await`; the remote command may keep running.
An HTTP 504 during the wait is a timeout error.

Omitted memory keeps the Blaxel default. CPU, CPU limit, memory limit, and
disk overrides are rejected.

`@blaxel/core` applies credentials through process-global `settings`. Use one
Blaxel Kit client, or one credential set, per process; a second client
overwrites the first.

## Examples

| Example       | Project                                                                                                                              |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| Creation      | [create-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/configuration/create-blaxel-sandbox) |
| Run a command | [run-blaxel-command](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/process/run-blaxel-command)             |
| Process       | [process examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/process)                                  |
| Read a file   | [read-blaxel-file](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/filesystem/read-blaxel-file)              |
| Filesystem    | [filesystem examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/filesystem)                            |
| Retrieval     | [get-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/lifecycle/get-blaxel-sandbox)           |
| Full index    | [TypeScript examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript)                                              |

## License

MIT

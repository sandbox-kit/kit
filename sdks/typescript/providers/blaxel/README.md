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

Closing the client releases local resources only. Reclaim sandboxes with TTL /
idle policies or `delete`.

## Configuration

| Setting        | Notes                                                   |
| -------------- | ------------------------------------------------------- |
| API key        | `auth: { apiKey: { key } }`                             |
| Omitted `auth` | Uses `BL_API_KEY` / `BL_WORKSPACE` from the environment |
| Workspace      | `scope.workspace` — Blaxel workspace name               |
| Region         | `region` — required for create                          |
| Endpoint       | Not supported                                           |

Token-pair, bearer, and OAuth are rejected. Only `scope.workspace` is accepted
on scope.

In TypeScript, use one Blaxel client (or one API key) per process. A second
client with different credentials replaces the first.

## What you can do

### Create

You need a valid `name` (lowercase alphanumeric with single hyphens, max 49
characters), an image, and `region`. You can also set `memoryMib`, labels,
environment variables, ports, outbound domains, max-age TTL (≥5 minutes), and
idle-delete (≥1 minute). You cannot disable idle terminate. If you omit
memory, Blaxel applies its default. CPU and disk overrides are rejected —
Blaxel derives CPU from memory. Set `memoryMib` for memory size; Blaxel
documents that value as MB and Kit does not convert units.

Your client/create `timeout` is how long create may take. Ready wait is not
supported (`waitFor` only Default or Scheduled). Ports 80, 443, and 8080 are
reserved. Default port transport is HTTP; plain is TCP; TLS is TLS; HTTP/2
over TLS is not supported. Volumes are not part of create.

If the name already exists, you get `AlreadyExists`.

### Lifecycle

You can `get` / `refresh` (by name), `delete`, and `close`. `close` does not
delete the sandbox. Start, stop, pause, archive, terminate, and attach are not
available.

After a successful `delete`, you get `state: "deleted"` and `terminal: true`.
Blaxel can still be cleaning up — recreating the same name right away may
return `Conflict`.

On `get` / `refresh` / `info`, an idle sandbox shows as `paused`; an active one
shows as `running`.

### Filesystem

You get the usual file APIs: `readFile`, `writeFile`, `copyFromLocal`,
`copyToLocal`, `makeDirectory`, `listDirectory`, `stat`, and `remove`. Paths
must be absolute. `parents` defaults to true; `parents: false` is rejected.

For directories, `stat` returns type directory. For files, `stat` returns type
file with size and permission fields when Blaxel provides them.
Missing local files for `copyFromLocal` return not found with field
`local_path`.

### Process

`sandbox.process.run` waits until the command finishes. Set exactly one of
`argv` or `shell`. `argv` is sent as one quoted shell command string. A non-zero exit code is not
an error — check `exitCode`. You get both stdout and stderr.

Omit or zero `timeout` to wait without a Blaxel process time limit (your
signal can still abort). A positive timeout must be whole seconds (`2000` is
two seconds; `1500` is rejected) and caps how long the command may run.
Aborting the signal stops waiting; the remote command may keep running. If the
wait times out, you get a timeout error.

### Errors you will see

| Situation                                  | Kind            |
| ------------------------------------------ | --------------- |
| Create when the name already exists        | `AlreadyExists` |
| Conflict while a delete is still finishing | `Conflict`      |

## Examples

| Example       | Project                                                                                                                              |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| Creation      | [create-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/configuration/create-blaxel-sandbox) |
| Run a command | [run-blaxel-command](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/process/run-blaxel-command)             |
| Process       | [process examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/process)                                  |
| Read a file   | [read-blaxel-file](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/filesystem/read-blaxel-file)              |
| Filesystem    | [filesystem examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/filesystem)                            |
| Retrieval     | [get-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/lifecycle/get-blaxel-sandbox)           |
| Delete        | [delete-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/blaxel/lifecycle/delete-blaxel-sandbox)     |
| Full index    | [TypeScript examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript)                                              |

## License

MIT

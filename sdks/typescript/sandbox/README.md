# `@sandbox-kit/sandbox`

Create and manage cloud sandboxes with a consistent TypeScript client interface.
Pick a provider package (`@sandbox-kit/daytona`, `@sandbox-kit/modal`, or
`@sandbox-kit/blaxel`), construct a client, then create or open sandboxes
through that client.

Closing a client or handle releases **local** SDK resources only. It does not
stop, terminate, or delete the remote sandbox.

## Install

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/daytona
# or
npm install @sandbox-kit/sandbox @sandbox-kit/modal
# or
npm install @sandbox-kit/sandbox @sandbox-kit/blaxel
```

Requires Node.js 20, 22, or 24.15+.

## Quick start

Daytona sample (`create()` with no body uses Daytona's default snapshot):

```ts
import { daytona } from "@sandbox-kit/daytona";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
});

const sandbox = await client.create();
console.log(sandbox.id, client.providerName);
```

Prefer `await using` when available. Otherwise call `await client.close()` in a
`finally` block.

## Client configuration

Pass the provider on `ClientConfig`, then shared settings:

| Field      | Purpose                                                                    |
| ---------- | -------------------------------------------------------------------------- |
| `provider` | Provider factory (`daytona()`, `modal()`, `blaxel()`)                      |
| `auth`     | One credential shape (`apiKey`, `bearerToken`, `tokenPair`, …)             |
| `scope`    | App, environment, organization, project, or workspace (provider-dependent) |
| `endpoint` | API URL when the provider supports it                                      |
| `region`   | Target or placement preference                                             |
| `timeout`  | Default operation deadline in **milliseconds**                             |

Omit `auth` when the provider SDK should resolve credentials from its own
environment. Provider-specific auth and scope rules live in each provider
package README.

## Create a sandbox

Daytona example:

```ts
const sandbox = await client.create({
  source: { image: { reference: "daytonaio/playground:latest" } },
  provisioning: { timeout: 120_000 },
});
```

`create()` with no options uses the provider default when that is valid (for
example Daytona's default snapshot). Modal requires an image. Blaxel requires
`name`, an image, and client `region`. Timeouts and delays are millisecond
numbers. Nested exclusive fields use shapes like `source: { image: … }` or
`auth: { apiKey: … }`.

## Open an existing sandbox

```ts
const attached = await client.attach(id); // Modal: id only, no existence check
const fetched = await client.get(id); // Daytona: sandbox id; Blaxel: create name
```

## Run a command

```ts
const result = await sandbox.process.run({
  argv: ["echo", "hello from kit"],
});
console.log(result.exitCode, new TextDecoder().decode(result.stdout));
```

Set exactly one of `argv` or `shell`. A non-zero exit code is not an SDK
error—check `exitCode`. Optional `workingDirectory` must be absolute.
`environment` applies to that call only. Omit `timeout` to use the client
deadline; use `0` for no extra deadline on this call. Aborting the signal
rejects the `await`; the remote command may keep running. See the
[process guide](https://github.com/sandbox-kit/kit/blob/main/docs/process.md).

## Read and write files

```ts
await sandbox.filesystem.writeFile(
  "/tmp/note.txt",
  new TextEncoder().encode("hello\n"),
);
const data = await sandbox.filesystem.readFile("/tmp/note.txt");
console.log(new TextDecoder().decode(data));
```

Sandbox paths are absolute. The same handle also supports `copyFromLocal`,
`copyToLocal`, `makeDirectory`, `listDirectory`, `stat`, and `remove`. See the
[filesystem guide](https://github.com/sandbox-kit/kit/blob/main/docs/filesystem.md).

## Sandbox handle

| Method           | Purpose                                       |
| ---------------- | --------------------------------------------- |
| `id`             | Sandbox identifier                            |
| `info()`         | Copied metadata                               |
| `capabilities()` | Lifecycle operations this handle can run      |
| `supports(...)`  | Check one capability                          |
| `refresh`        | Refresh metadata from the provider            |
| `start`          | Start                                         |
| `stop`           | Stop (`force` when supported)                 |
| `pause`          | Pause                                         |
| `archive`        | Archive                                       |
| `delete`         | Delete (`wait` when supported)                |
| `terminate`      | Terminate (`wait` when supported)             |
| `close`          | Close the handle. The sandbox keeps running   |
| `filesystem`     | Read, write, copy, list, mkdir, stat, remove  |
| `process`        | Run a command and wait for exit code + output |

Check `supports` before a lifecycle call that might be unavailable. On Blaxel,
you can `get`, `refresh`, `delete`, and `close`. Other lifecycle methods are
not available. If the provider rejects a supported call, you get that
provider's error.

## Errors

Failures throw `SandboxError` with `kind`, `provider`, `operation`, `field`,
and `message`:

```ts
import { Client, SandboxError } from "@sandbox-kit/sandbox";

try {
  await client.create({ resources: { cpuCores: -1 } });
} catch (error) {
  if (error instanceof SandboxError) {
    console.error(error.kind, error.field, error.message);
  } else {
    throw error;
  }
}
```

## Providers

| Package                | Role              |
| ---------------------- | ----------------- |
| `@sandbox-kit/daytona` | Daytona sandboxes |
| `@sandbox-kit/modal`   | Modal sandboxes   |
| `@sandbox-kit/blaxel`  | Blaxel sandboxes  |

## Examples

Runnable projects:
[examples/typescript](https://github.com/sandbox-kit/kit/tree/main/examples/typescript).

## License

MIT

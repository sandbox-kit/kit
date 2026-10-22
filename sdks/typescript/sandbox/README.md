# `@sandbox-kit/sandbox`

Create and manage cloud sandboxes from one TypeScript API. Pick a provider
package (`@sandbox-kit/daytona` or `@sandbox-kit/modal`), construct a client,
then create or open sandboxes through the same methods.

Closing a client or handle releases **local** SDK resources only. It does not
stop, terminate, or delete the remote sandbox.

## Install

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/daytona
# or
npm install @sandbox-kit/sandbox @sandbox-kit/modal
```

Requires Node.js 20, 22, or 24.15+.

## Quick start

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

| Field      | Purpose                                                         |
| ---------- | --------------------------------------------------------------- |
| `provider` | Provider factory (`daytona()`, `modal()`, …)                    |
| `auth`     | One credential shape (`apiKey`, `bearerToken`, `tokenPair`, …)  |
| `scope`    | App, environment, organization, or project (provider-dependent) |
| `endpoint` | API URL when the provider supports it                           |
| `region`   | Target or placement preference                                  |
| `timeout`  | Default operation deadline in **milliseconds**                  |

Omit `auth` when the provider SDK should resolve credentials from its own
environment. Provider-specific auth and scope rules live in each provider
package README.

## Create a sandbox

```ts
const sandbox = await client.create({
  source: { image: { reference: "daytonaio/playground:latest" } },
  provisioning: { timeout: 120_000 },
});
```

`create()` with no options uses the provider default when that is valid (for
example Daytona's default snapshot). Timeouts and delays are millisecond
numbers. Nested exclusive fields use shapes like `source: { image: … }` or
`auth: { apiKey: … }`.

## Open an existing sandbox

```ts
const attached = await client.attach(id); // provider must support attach
const fetched = await client.get(id); // provider-verified lookup
```

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
| `close`          | Drop the local handle (not the cloud sandbox) |

Unsupported operations throw before the provider call. Use `supports(...)` when
the sandbox class may differ (for example Daytona container vs VM).

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

## Examples

Runnable projects:
[examples/typescript](https://github.com/sandbox-kit/kit/tree/main/examples/typescript).

## License

MIT

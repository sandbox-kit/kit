# `@sandbox-kit/daytona`

Daytona provider for [`@sandbox-kit/sandbox`](https://www.npmjs.com/package/@sandbox-kit/sandbox).
Use the shared client interface; this package supplies the `daytona()` provider.

## Install

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/daytona
```

Requires Node.js 20, 22, or 24.15+, and a Daytona API key (or credentials your
Daytona SDK already resolves from the environment).

## Quick start

```ts
import { daytona } from "@sandbox-kit/daytona";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
});

// Default snapshot for the target:
const sandbox = await client.create();
console.log(sandbox.id);
```

Create from an image and size resources:

```ts
const sandbox = await client.create({
  source: { image: { reference: "daytonaio/playground:latest" } },
  resources: {
    cpuCores: 2,
    memoryMib: 4096n,
    diskMib: 8192n,
  },
  runtime: { language: "typescript" },
});
```

Closing the client releases local SDK resources only. Daytona keeps the sandbox
until you stop or delete it, or a lifetime policy fires.

## Configuration

| Setting         | Notes                                                         |
| --------------- | ------------------------------------------------------------- |
| API key         | `auth: { apiKey: { key } }`                                   |
| Bearer / JWT    | `auth: { bearerToken: { token } }` — requires organization ID |
| Omitted `auth`  | Daytona SDK environment credentials                           |
| Organization    | `scope: { organizationId }`                                   |
| Endpoint        | `endpoint` — defaults to `https://app.daytona.io/api`         |
| Region / target | `region`                                                      |

Remote endpoints require HTTPS. HTTP is allowed for localhost or loopback only.

## What you can do

**Create** from the default snapshot, a named snapshot, or a registry image.
Optional runtime language/user, network allowlists, secrets, labels, TTL,
ephemeral behavior, and stop/pause/archive/delete policies.

**Lifecycle** on the handle: `get`, `refresh`, `start`, `stop`, `pause`,
`archive`, `delete`, and `close`. Daytona accepts or rejects each lifecycle
call for that sandbox. `close` does not contact Daytona and does not change
the sandbox.

Daytona does not support unverified `attach`. Use `client.get(id)` for a
provider-verified handle.

CPU requests are whole cores. Memory and disk are whole GiB expressed as MiB.
Resource overrides on create require an **image** source; snapshot creates
inherit snapshot resources.

## Examples

| Example    | Project                                                                                                                               |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| Creation   | [create-daytona-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/daytona/create-daytona-sandbox)             |
| Policies   | [create-daytona-with-policies](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/daytona/create-daytona-with-policies) |
| Retrieval  | [get-daytona-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/typescript/daytona/get-daytona-sandbox)                   |
| Full index | [TypeScript examples](https://github.com/sandbox-kit/kit/tree/main/examples/typescript)                                               |

## License

MIT

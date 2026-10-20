# Sandbox Kit for TypeScript

Core client package (`@sandbox-kit/sandbox`). Providers such as
`@sandbox-kit/daytona` and `@sandbox-kit/modal` are published
separately and plug into this module. Apps and providers share public
`CreateOptions` / `ClientConfig` (`Config` plus `provider`): nested exclusive
fields (`auth: { apiKey: … }`, `source: { image: … }`, `readiness: { command: … }`),
plain maps/allowlists (`outboundCidrs: string[]`), and millisecond timeouts.
`ErrorInfo.field` uses canonical protobuf paths (`CreateField*` / `ConfigField*` /
`InfoField*`), matching Go. Closing the client releases local provider SDK
resources only; it does not stop or delete cloud sandboxes.

## Public input → wire boundary

Applications and provider backends speak the **public** shapes above. Inside
`Client`, requests are converted before Protovalidate and provider mapping:

1. **Public input** — ms numbers, nested oneofs, `string[]` allowlists, plain
   `providerOptions` maps.
2. **`materialize*`** — package-private conversion to protobuf wire messages.
3. **Validate** — shared Protovalidate rules (`secure_endpoint` stays handwritten),
   then provider checks.
4. **Provider mappings** — native SDK calls using the validated wire request.

Wire types stay package-private to the client.

## Usage

```ts
import {
  Client,
  type ClientConfig,
  type CreateOptions,
} from "@sandbox-kit/sandbox";
import { daytona } from "@sandbox-kit/daytona";

const config: ClientConfig = {
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
};

const createOptions: CreateOptions = {
  source: { image: { reference: "daytonaio/playground:latest" } },
  provisioning: { timeout: 120_000 },
};

await using client = new Client(config);
const sandbox = await client.create(createOptions);
console.log(sandbox.id);
```

Prefer `await using` when the runtime supports explicit resource management;
otherwise call `await client.close()` in `finally`. Both forms have identical
cleanup semantics.

## Examples

| Language   | Index                                                         |
| ---------- | ------------------------------------------------------------- |
| TypeScript | [examples/typescript](../../../examples/typescript/README.md) |
| Go         | [examples/go](../../../examples/go/README.md)                 |

Install a provider package such as `@sandbox-kit/daytona` or
`@sandbox-kit/modal` beside this client. See the
[root README](../../../README.md) for a full create snippet.

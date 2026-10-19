# Sandbox Kit for TypeScript

Core client package (`@sandbox-kit/sandbox`). Providers such as
`@sandbox-kit/provider-daytona` and `@sandbox-kit/provider-modal` are published
separately and plug into this module. Shared create and client-config validation
uses Protovalidate (`secure_endpoint` stays handwritten). Closing the client
releases local provider SDK resources only; it does not stop or delete cloud
sandboxes.

## Usage

```ts
import { Client } from "@sandbox-kit/sandbox";

const client = new Client(config);

try {
  const sandbox = await client.create();
  console.log(sandbox.id);
} finally {
  await client.close();
}
```

With TypeScript 5.2+, `ESNext.Disposable`, and a runtime that supports explicit
resource management:

```ts
await using client = new Client(config);
const sandbox = await client.create();
```

Both forms have identical cleanup semantics.
`client[Symbol.asyncDispose]()` runs automatically at the end of an `await using`
scope.

## Examples

| Language   | Index                                                         |
| ---------- | ------------------------------------------------------------- |
| TypeScript | [examples/typescript](../../../examples/typescript/README.md) |
| Go         | [examples/go](../../../examples/go/README.md)                 |

Install a provider package such as `@sandbox-kit/provider-daytona` or
`@sandbox-kit/provider-modal` beside this client. See the
[root README](../../../README.md) for a full create snippet.

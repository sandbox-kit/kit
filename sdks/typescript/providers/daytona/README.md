# Daytona provider for Sandbox Kit

Independent package over `@daytona/sdk@0.225.0`, published separately from the
core `@sandbox-kit/sandbox` client. Use it through that shared API.

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/daytona
```

## Usage

```ts
import { Client } from "@sandbox-kit/sandbox";
import { daytona } from "@sandbox-kit/daytona";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
});

const sandbox = await client.create({
  source: { image: { reference: "daytonaio/playground:latest" } },
  runtime: { language: "typescript" },
});
console.log(sandbox.id);
```

Auth may be omitted when the native Daytona SDK should resolve credentials from
its environment. Closing the client releases local SDK resources only.

## Examples

| Example    | Project                                                                                                        |
| ---------- | -------------------------------------------------------------------------------------------------------------- |
| Creation   | [create-daytona-sandbox](../../../../examples/typescript/daytona/create-daytona-sandbox/README.md)             |
| Policies   | [create-daytona-with-policies](../../../../examples/typescript/daytona/create-daytona-with-policies/README.md) |
| Retrieval  | [get-daytona-sandbox](../../../../examples/typescript/daytona/get-daytona-sandbox/README.md)                   |
| Full index | [TypeScript examples](../../../../examples/typescript/README.md)                                               |

## Implementation notes

Checks and portable field mappings are handwritten in `checks/provider_checks.ts`,
`mappings/provider_mappings.ts`, and `config/provider_config.ts`, kept in sync with
`specs/providers/daytona.yaml`. `src/daytona.ts` owns typed Daytona SDK call
ordering and lifecycle orchestration; `src/mappings.ts` owns remaining
provider-specific composition.

## Dependency audit

`npm audit` reports four high-severity findings in the pinned
`@daytona/sdk@0.225.0` dependency chain (`fast-glob` → `micromatch` → `braces`).
The registry reports no fix (`braces` latest is still vulnerable; Daytona still
depends on `fast-glob`). The provider does not override those transitive
packages, because overrides could create an unsupported SDK combination.
Reassess the audit when Daytona drops or patches that chain.

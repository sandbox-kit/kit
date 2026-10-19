# Daytona provider for Sandbox Kit

Optional package over `@daytona/sdk@0.224.0`. Use it through `@sandbox-kit/sandbox`.

## Usage

```ts
import { Client } from "@sandbox-kit/sandbox";
import { daytona } from "@sandbox-kit/provider-daytona";

const client = new Client({ provider: daytona() });

try {
  const sandbox = await client.create({
    runtime: { language: "typescript" },
  });
  console.log(sandbox.id);
} finally {
  await client.close();
}
```

The native Daytona SDK resolves credentials from its environment when `auth` is
omitted. Closing the client releases local SDK resources only.

## Examples

| Example    | Project                                                                                                        |
| ---------- | -------------------------------------------------------------------------------------------------------------- |
| Creation   | [create-daytona-sandbox](../../../../examples/typescript/daytona/create-daytona-sandbox/README.md)             |
| Policies   | [create-daytona-with-policies](../../../../examples/typescript/daytona/create-daytona-with-policies/README.md) |
| Retrieval  | [get-daytona-sandbox](../../../../examples/typescript/daytona/get-daytona-sandbox/README.md)                   |
| Full index | [TypeScript examples](../../../../examples/typescript/README.md)                                               |

## Implementation notes

Checks and portable field mappings are handwritten in `provider_checks.ts`,
`provider_mappings.ts`, and `provider_config.ts`, kept in sync with
`specs/providers/daytona.yaml`. `src/daytona.ts` owns typed Daytona SDK call
ordering and lifecycle orchestration; `src/mappings.ts` owns remaining
provider-specific composition.

## Dependency audit

As of 2026-10-14, `npm audit` reports four high-severity findings in the pinned
`@daytona/sdk@0.224.0` dependency chain (`fast-glob` → `micromatch` → `braces`).
The registry currently reports no fix. The provider does not override those
transitive packages because doing so could create an unsupported SDK
combination; reassess the audit when upgrading the Daytona SDK.

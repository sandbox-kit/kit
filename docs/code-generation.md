# Code generation

The target architecture is documented in
[protobuf SDK boundary](protobuf-sdk-boundary.md). Protobuf is the one source of
truth for portable data and operation contracts. Buf (`buf generate`) emits
messages, enums, reflection metadata, and service declarations via standard
plugins; handwritten SDK code owns client, handle, provider, checks, mappings,
and error behavior. Shared create and client-config validation is Protovalidate
options on those messages, evaluated by language Protovalidate libraries
(`secure_endpoint` stays handwritten).

Portable create/config/info types come from standard protobuf contracts plus
handwritten SDK aliases. Kit custom emitters are retired.

## Generate and verify

The repository is a Buf v2 module rooted at `proto/`, with
`buf.build/bufbuild/protovalidate` as a module dependency (see `buf.yaml` /
`buf.lock`). Buf checks schema formatting, lint, compilation, and compatibility,
and generates contract packages through `buf.gen.yaml`.

From `tooling/`, check the protobuf module:

```sh
GOWORK=off go run ./cmd/sandbox-kit test proto
GOWORK=off go run ./cmd/sandbox-kit test proto --against '.git#branch=main,subdir=proto'
```

The second command adds `buf breaking` against the selected Git or Buf input.
The initial lint baseline is `BASIC`: it catches structural mistakes without
renaming established enum defaults or the local `SandboxService` declaration.

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit generate contracts-go
GOWORK=off go run ./cmd/sandbox-kit generate contracts-typescript
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit test typescript
```

Contract generation requires the Buf CLI. Both generate commands run
`buf generate` for the plugins declared in `buf.gen.yaml`.
`contracts-typescript` also runs `npm run build` in
`sdks/typescript/contracts` so `@sandbox-kit/contracts` exports compiled
descriptors (including Protovalidate options) from `dist/`.

Contract packages are committed under `sdks/<language>/contracts/`. Edit
`proto/`, then regenerate contracts. SDK applications use native types and do
not load YAML or import protobuf messages at runtime.

## Inputs

| Input                             | Defines                                                                            |
| --------------------------------- | ---------------------------------------------------------------------------------- |
| `proto/kit/sandbox/v1/`           | Public portable types, service requests/responses, and operation contracts         |
| `specs/validation.yaml`           | Portable create-validation inventory (mirrors Protovalidate options in proto)      |
| `specs/client.yaml`               | Config validation inventory (mirrors Protovalidate; `secure_endpoint` handwritten) |
| `specs/contracts.yaml`            | Metadata/error contract limits                                                     |
| `specs/providers/<provider>.yaml` | Provider capability, check, mapping, and error-classification inventories          |

Specs use `version: 1`. Strict loading rejects unknown YAML keys, extra documents,
and unsupported versions. Shared fields use protobuf names; SDK symbols belong
inside the language's bindings.

## API declarations

Public data signatures and service request/response shapes belong in protobuf.
Runtime objects and behavior are handwritten in each language. YAML must not
declare types, constructors, methods, result shapes, visibility, or operation
sequencing.

## Response and error contracts

Protobuf owns portable `ErrorInfo` and `ErrorKind`. SDK runtimes alias those
types and handwrite wrappers (`NewError` / `WithErrorContext` in Go,
`sandboxError` in TypeScript) plus per-provider classification mappers.

Handwritten `CreateField*`, `ConfigField*`, and `InfoField*` constants are the
field paths used by validation, mappings, and errors.

`specs/contracts.yaml` identifies the shared metadata, error, and origin
declarations and the copy depth, visited-node, and payload-byte limits.
Portable messages come from standard protobuf generation into contract packages;
SDK runtimes alias them and keep handwritten clone helpers (for example Go
`clone.go` and TypeScript metadata helpers).

Provider error classification is not generated: Modal and Daytona mappers live
in handwritten `errors.go` / `errors.ts`. Response groups declare
`origin: provider|request` for mapped metadata; unavailable fields receive no
fabricated values.

See [responses and errors](responses-and-errors.md) for the public contract,
native cause preservation, and copy restrictions.

## Shared validation rules

```yaml
version: 1
messages:
  Resources:
    fields:
      cpu_cores: { finite: true, exclusive_minimum: 0 }
    constraints:
      - { op: limit, field: cpu_limit_cores, other: cpu_cores }
      - { op: requires_positive, field: cpu_limit_cores, other: cpu_cores }
```

Field rules include `minimum`, `exclusive_minimum`, `maximum`, `finite`,
`nonblank`, `min_items`, `format: http_url`, and `each` for list elements.
Serialization policy is separate from validation. Client representation metadata
marks credential fields `sensitive`, which excludes them from Go JSON output
without treating secrecy as a value-validity rule.

| Constraint          | Meaning                                       |
| ------------------- | --------------------------------------------- |
| `at_most_one`       | At most one listed field is selected          |
| `exactly_one`       | Exactly one listed field is selected          |
| `requires_any`      | At least one listed field is supplied         |
| `forbids`           | None of the listed fields is selected         |
| `limit`             | A positive limit must be at least its request |
| `requires_positive` | A positive limit requires a nonzero request   |

Conditions use `when: {field: mode, equals: 2}` or `not_equals`.
Selection uses presence for optional values/messages, nonempty collections,
true booleans, and nonzero/nonempty ordinary scalars.

Shared create validation is handwritten in the Go SDK
(`ValidateCreateOptions`). Keep `specs/validation.yaml` aligned with that
implementation. Both languages preserve explicit zero/false and supply no
defaults.

## Provider checks

Provider `checks` in `specs/providers/<provider>.yaml` inventory create-time
policy the integration enforces. Implementations are handwritten
(`provider_checks.go` / `provider_checks.ts`) and run before native creation
calls. Keep the YAML list and handwritten code in sync; checks are not emitted.

```yaml
checks:
  - { path: resources.cpu_cores, minimum: 0.125 }
  - { path: placement.cloud, allowed: [aws, gcp, oci, auto] }
  - { path: runtime.working_directory, format: absolute_path }
  - { path: network.outbound_cidrs.entries, format: cidr }
```

Supported formats include `absolute_path`, `http_url`, `cidr`, `domain`, and
`env_name`. Policy relationships use `exclusive_policies` and
`forbid_policy_with`. Dynamic account quotas, permissions, and capacity remain
provider-enforced.

## Typed field mappings

```yaml
groups:
  - name: mapResources
    direction: request
    source: { name: Resources }
    target:
      name: Resources
      bindings:
        go:
          import: github.com/daytona/clients/sdk-go/pkg/types
          name: Resources
          fields: { memory: Memory }
    fields:
      - {
          from: memory_mib,
          to: memory,
          transform: divide_exactly,
          factor: 1024,
          maximum: 2147483647,
          cast: integer,
        }
```

Mapping groups are handwritten in `provider_mappings.go` /
`provider_mappings.ts`. The YAML `groups:` list is an inventory: keep it aligned
with the implementations. Request helpers return the mapped native value and a
copy with mapped fields cleared. Providers reject remaining unsupported intent.
Response helpers return shared native structs.

| Transform          | Semantics                                                                            |
| ------------------ | ------------------------------------------------------------------------------------ |
| Copy               | Preserve the supplied value                                                          |
| `whole`            | Require an integral value                                                            |
| `divide_exactly`   | Reject remainders when converting units                                              |
| `multiply`         | Check sign/overflow before multiplying                                               |
| `scaled_integer`   | Require exact scaled precision; compensate for binary rounding before SDK truncation |
| `duration_seconds` | Require positive whole seconds and bound the seconds value                           |
| `policy_minutes`   | Map an automatic policy to an optional native minute interval                        |

`maximum` bounds the converted value. `policy_minutes` uses
`policy.disabled: zero|reject` and `policy.immediate: true|false`.
Absent/default policies stay omitted. See [creation policies](sandbox-creation.md#policy-semantics).

Portable casts are `integer`, `number`, `unsigned`, and `string`.
A group's `when` conditions support numeric `greater_than`; all must hold.
This controls availability of Daytona allocated-resource metadata.

## Client composition and ownership

A provider's `client` spec inventories how mapping groups compose:

- `settings` maps top-level configuration.
- `scope` maps scope fields.
- `auth` declares supported credential alternatives.
- `managed` marks Kit-owned configuration such as timeout that is not mapped
  into the native client.
- `retained` requires a matching backend state capture.
- `rejected` describes unsupported supplied configuration.
- `destination` attaches a mapped object using the native field and declared
  `objects` type, as with Modal OAuth credentials.

Handwritten `provider_config.ts` (and Go construction paths) implement that
composition. Keep full schema/native identities, identifiers, field ownership,
and destination collisions consistent with the inventory.

## Runtime bindings

`runtime.go.constructor.function` documents the native constructor taking a
pointer to mapped configuration and returning a client/error.
`runtime.go.cleanup` names the cleanup method and declares context/error behavior.

Backend state capture and metadata cloning are handwritten in the SDK runtime
(for example Go `clone.go` and TypeScript metadata helpers).

## Extend a provider or language

1. Declare shared types and method intent in protobuf.
2. Describe portable validation and conversions in YAML inventories, and
   operation intent on an annotated `SandboxService` method.
3. Implement native SDK integration, checks, and mappings in the provider
   package.
4. Add annotated service operations that use existing outcome shapes.
5. Add semantic boundary tests and update examples/docs.
6. Run the CLI verification command.

SDK orchestration that the inventory cannot express stays in handwritten
provider integration files (`sandbox.go` in Go and the provider-named `.ts`
plus `mappings.ts` files in TypeScript): lookup sequences, source selection,
readiness, and remaining provider semantics. Handwritten
`provider_mappings` / `provider_config` own field conversions and native client
configuration; handwritten `provider_checks` enforces create checks; handwritten
`errors` classifies native provider errors.

Lifecycle option and result types come from protobuf contracts. Handle methods,
typed provider dispatch, capability discovery, storage, and close coordination
remain handwritten runtime code shaped by annotated `SandboxService` methods.

See [tooling](../tooling/README.md), [naming](naming.md), and
[reference verification](provider-verification.md).

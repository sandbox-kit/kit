# Code generation

The target architecture is documented in
[protobuf SDK boundary](protobuf-sdk-boundary.md). Protobuf is the one source of
truth for portable data and operation contracts. Standard protobuf plugins emit
messages, enums, reflection metadata, and service declarations; handwritten SDK
code owns client, handle, provider, and error behavior.

The repository is currently being migrated from legacy Kit emitters to that
model. Do not extend the legacy API-declaration, template, or lifecycle emitter
paths while the migration is in progress.

## Generate and verify

The repository is a Buf v2 module rooted at `proto/`. Buf checks schema
formatting, lint, compilation, and compatibility; the Kit CLI continues to own
SDK generation so `specs/targets.yaml` remains the only target manifest.

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
GOWORK=off go run ./cmd/sandbox-kit generate go
GOWORK=off go run ./cmd/sandbox-kit generate contracts-go
GOWORK=off go run ./cmd/sandbox-kit generate contracts-typescript
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit generate typescript
GOWORK=off go run ./cmd/sandbox-kit test typescript
```

Generation requires `protoc`. Temporary binaries are cleaned up.

Generated source is committed under `sdks/<language>/`. Edit its schema, spec,
or emitter, then regenerate. SDK applications use native types and do not load
YAML or import protobuf messages.

## Inputs

| Input                             | Defines                                                                 |
| --------------------------------- | ----------------------------------------------------------------------- |
| `proto/kit/sandbox/v1/`           | Public portable types, service requests/responses, and operation contracts |
| `proto/kit/sandbox/internal/`     | Internal-only portable DTOs                                             |
| `proto/kit/providers/`            | Provider contract declarations                                          |
| `specs/targets.yaml`              | Language/provider output modules and source descriptors                 |
| `specs/validation.yaml`           | Portable validation rules                                               |
| `specs/providers/<provider>.yaml` | Provider capability, validation, mapping, and error-classification rules |

Specs use `version: 1`. Strict loading rejects unknown YAML keys, extra documents,
and unsupported versions. Shared fields use protobuf names; SDK symbols belong
inside the language's bindings.

## API declarations

Public data signatures and service request/response shapes belong in protobuf.
Private portable data belongs in the internal protobuf package. Runtime objects
and behavior are handwritten in each language. YAML must not declare types,
constructors, methods, result shapes, visibility, or operation sequencing.

## Generation plan and client operations

`specs/generation.yaml` selects semantic output phases for shared contracts,
clients, and providers. The Go dispatcher consumes this plan rather than loading
specifications independently for each emitter.

Client operation instructions also live in this file:

```yaml
operations:
  create:
    - require_context
    - require_client
    - prepare_request
    - resolve_deadline
    - invoke_creation
    - return_handle
```

The compiler verifies the version 1 vocabulary and required ordering. It rejects
missing/duplicate phases, unknown instructions, and unsafe operation sequences.
Initialization validates configuration before initializing the provider; creation
prepares an owned request before applying deadlines and invoking the backend.
Cleanup releases SDK resources without deleting sandboxes. The target manifest
drives provider/language outputs, so adding a provider does not require editing
the generation command's source list.

These are semantic instructions, not arbitrary expressions or YAML code snippets.
Constructor and method names come from `specs/api.yaml`. Output suffixes come
from the language profile. The Go emitter turns these instructions into
synchronous calls, context deadlines, and error results. The TypeScript emitter
uses promises, abort signals, timeouts, and thrown errors.

Handle lifecycle operations are registered once as annotated methods on
`SandboxService` in `proto/kit/sandbox/v1/client.proto`. The shared compiler
validates their portable operation metadata and option types. The Go lifecycle
emitter derives public handle methods, typed dispatch, timeouts, metadata
replacement, terminal updates, and capability discovery from that service. A
new operation using an existing outcome shape therefore does not require adding
its name to the Go emitter.

The model resolves shared validation references, shared mapping fields,
conversion vocabulary, policy relationships, and error kinds/protocol codes.
Language-specific SDK member checks and generated validator syntax remain in
the native emitters. Handwritten provider orchestration is still required where
noted below.

## Response and error contracts

Go emitters build error expressions through `internal/go/error_gen`; TypeScript
emitters use `internal/typescript/error_gen`. Each language's schema package
resolves field paths against full protobuf identities.
Generated `CreateField*`, `ConfigField*`, and `InfoField*` constants are the
field paths used by mappings and errors. Nested mapping sources can declare
`error_path`; unknown paths fail generation.

Response origin validation is generated from schema presence rules. It rejects
unknown keys and absent values, while accepting explicit zero/false and empty
collections. Context enrichment avoids an extra wrapper when details already match.

`specs/contracts.yaml` identifies the shared metadata, error, and origin
declarations and the copy depth, visited-node, and payload-byte limits.
Error and metadata-copy signatures come from `specs/api.yaml`. `types_gen/runtime.go` emits their behavior bodies
into `.errors.runtime.gen.go` and `.clone.runtime.gen.go`. Schema `Clone`
methods are separate: templates declare them, and `*.copy.gen.go` copies fields
directly, preserving numeric types and owned collections. Metadata has explicit
signed and unsigned alternatives instead of protobuf Struct's double-only
numeric representation.

`specs/errors.yaml` declares ordered cause rules, fallback classification, context
reuse, detail preservation, and provider classification precedence. Version 1
supports `canceled` and `deadline_exceeded` cause matches. Go emitters translate
these into `errors.Is` checks; the first matching cause overrides the supplied
kind. An invalid kind falls back to `unknown`.

Provider classification follows existing shared error, native type, HTTP status,
gRPC code, then fallback. Version 1 requires this order and detail preservation;
unsupported policies fail generation. `reuse_matching_context` controls whether
an error with matching context is returned unchanged.

Provider `error_statuses` and `error_grpc` map protocol numbers to semantic kinds.
Modal `error_rules` separates portable classifications from native bindings:

```yaml
error_rules:
  timeout_error: timeout
errors:
  go:
    types:
      - { import: github.com/modal-labs/modal-client/go, name: TimeoutError, rule: timeout_error }
```

Another language binds the same rule to its native SDK type. Native type names,
imports, pointer matching, and field members belong in `errors.<language>`.
Generation rejects unknown kinds, missing rules, invalid protocol codes, and
missing member bindings. Response groups declare `origin: provider|request` for mapped
metadata; unavailable fields receive no fabricated values.

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
marks credential fields `sensitive`, which excludes them from generated Go JSON
output without treating secrecy as a value-validity rule.

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

The Go emitter uses validator tags, custom checks, and struct callbacks. The
TypeScript emitter generates nested validation functions and enum membership
checks. Both preserve explicit zero/false and supply no defaults. Go optional
fields use `omitnil`; repeated Go message validation rejects nil elements.

## Provider checks

Provider `checks` validate common creation fields using schema paths:

```yaml
checks:
  - { path: resources.cpu_cores, minimum: 0.125 }
  - { path: placement.cloud, allowed: [aws, gcp, oci, auto] }
  - { path: runtime.working_directory, format: absolute_path }
  - { path: network.outbound_cidrs.entries, format: cidr }
```

Supported formats are `absolute_path`, `http_url`, `cidr`, `domain`, and
`env_name`. Checks support repeated strings, repeated-message paths, and
environment-map keys. `allow_empty` preserves documented sentinel behavior for
supported formats.

Policy relationships use `exclusive_policies` and `forbid_policy_with`.
Current relationships reject simultaneous positive stop/pause delays and
incompatible ephemeral/immediate-delete combinations.

The emitter validates paths, types, formats, and bounds. Generated checks run
before native creation calls. Dynamic account quotas, permissions, and capacity
remain provider-enforced.

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
      - { from: memory_mib, to: memory, transform: divide_exactly, factor: 1024, maximum: 2147483647, cast: integer }
```

Shared types resolve against descriptors. Logical native fields resolve through
`bindings.go.fields`. Request helpers return the mapped native value and a copy
with mapped fields cleared. Providers reject remaining unsupported intent.
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

A provider's `client` spec combines existing mapping groups:

- `settings` maps top-level configuration.
- `scope` maps scope fields.
- `auth` declares supported credential alternatives.
- `managed` permits Kit-owned provider selection and timeout fields.
- `retained` requires a matching backend state capture.
- `rejected` describes unsupported supplied configuration.
- `destination` attaches a mapped object using the native field and declared
  `objects` type, as with Modal OAuth credentials.

Composition checks full schema/native identities, identifiers, field ownership,
and destination collisions. Exclusive auth branches may share destinations with
each other, but cannot overwrite settings or scope. Implicit precedence and
overrides are rejected.

## Runtime bindings

`runtime.go.constructor.function` binds a native constructor taking a pointer
to mapped configuration and returning a client/error.
`runtime.go.cleanup` names the cleanup method and declares context/error behavior.

Backend state declares its native name, shared type, optional flag, and
`capture` path relative to `Config`. Generated copies preserve absence and avoid
aliasing caller data.

Capture currently supports optional strings and messages with scalar/nested
message fields. Collections, bytes, recursive messages, and real oneofs fail
generation until copying support is added.

## Extend a provider or language

1. Declare shared types and method intent in protobuf.
2. Describe portable validation and conversions in YAML, and operation intent
   on an annotated `SandboxService` method.
3. Add native SDK symbols to the language's bindings.
4. Add annotated service operations that use existing outcome shapes;
   extend the shared vocabulary and affected emitters only for a new shape.
5. Add semantic boundary tests, regenerate, and update examples/docs.
6. Run the CLI verification command.

SDK orchestration that the vocabulary cannot express stays in handwritten
provider integration files (`sandbox.go` in Go and the provider-named `.ts` plus
`mappings.ts` files in TypeScript): lookup sequences, source selection,
readiness, and remaining provider semantics. Bindings alone do not produce a
complete integration. Generated TypeScript `provider.gen.ts` files own the
provider/backend identity bases, `provider.config.gen.ts` composes native client
configuration and retained state, `provider.checks.gen.ts` enforces portable
provider checks, `provider.errors.gen.ts` classifies native errors, and
`provider.mappings.gen.ts` contains compiled YAML mapping groups. None of these
files are embedded source templates.

Lifecycle option and result types are generated from protobuf. The Go and
TypeScript lifecycle emitters generate handle methods and typed provider dispatch
from annotated `SandboxService` methods; the Go emitter also generates static capability
discovery. Handle storage, operation/close coordination, and provider-specific
lifecycle call sequences remain native runtime code. Optional capabilities
prevent a provider from advertising operations it cannot implement; Daytona
resolves additional class-dependent capabilities from current metadata.

See [tooling](../tooling/README.md), [naming](naming.md), and
[reference verification](provider-verification.md).

Client mappings can declare `environment` resolution rules with a schema field,
ordered environment-variable names, and a default. The shared compiler checks
that the target is an optional scalar string. Native emitters resolve omitted
values into a private configuration copy before validation and SDK construction.
Daytona uses this to validate native endpoint environment fallbacks.

The portable `secure_endpoint` format requires remote HTTPS and permits HTTP
only for local loopback testing; `image_reference` declares supported registry
reference syntax. The shared model rejects unknown format names before emission.

# Code generation

Protobuf and YAML define the shared contracts and portable rules. All generators
are written in Go under `tooling/`; only Go SDK output is implemented today.

```text
Protobuf + versioned YAML + native SDK bindings
                    ↓
      Language-neutral compiled model
                    ↓
             Language emitter
                    ↓
            Native SDK source
```

`tooling/internal/model/` compiles descriptor identities, field paths, local client
methods, validation references, mapping plans, and error classifications before
emission. It uses protobuf reflection descriptors, not `protogen`, Go naming, or
runtime SDK types. `internal/go/` translates that model into Go syntax and files.

## Generate and verify

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit generate go
GOWORK=off go run ./cmd/sandbox-kit test go
```

Generation requires `protoc`. The CLI builds temporary generator binaries,
bootstraps annotation bindings, and generates runtime source from the repository
root, where the plugin resolves `specs/`. Temporary binaries are cleaned up.

Generated source is committed under `sdks/go/`. Edit its schema, spec, or emitter,
then regenerate. SDK applications use native types and do not load YAML or import
protobuf messages.

## Inputs

| Input                             | Defines                                                             |
| --------------------------------- | ------------------------------------------------------------------- |
| `proto/kit/sandbox/v1/`           | Shared types, presence, enums, and local methods                    |
| `proto/kit/providers/`            | Provider and native SDK declarations                                |
| `specs/behaviors.yaml` | Required behavior inputs, outputs, owner, and execution flags |
| `specs/templates.yaml` | Reusable descriptor-driven signatures and field-path rules |
| `specs/api.yaml` | Runtime types, constructors, methods, parameters, and results |
| `specs/languages/<language>.yaml` | Native representation, builtins, externals, and output conventions |
| `specs/generation.yaml` | Generation phases and ordered client-operation instructions |
| `specs/client.yaml`               | Provider selection and client/auth validation                       |
| `specs/validation.yaml`           | Shared creation validation                                          |
| `specs/contracts.yaml`            | Metadata/error declarations and copy depth                          |
| `specs/errors.yaml`               | Shared cause classification and error-context policy                 |
| `specs/providers/<provider>.yaml` | Client assembly, provider checks, conversions, and runtime bindings |

Specs use `version: 1`. Strict loading rejects unknown YAML keys, extra documents,
and unsupported versions. Shared fields use protobuf names; SDK symbols belong
inside the language's bindings.

## API declarations

[`api.yaml`](../specs/api.yaml) declares runtime objects, interfaces, fields,
constructors, methods, parameter/result shapes, and behavior attachments.
[`languages/go.yaml`](../specs/languages/go.yaml) lowers references, presence,
cancellation, multiple results, and failures into Go conventions.
`internal/go/declaration_gen/` generates declarations from the compiled model.
Template bindings retain portable type shapes and native symbol identities.
Behavior bodies use explicit symbol references, with no textual identifier renaming.
The shared compiler checks their declarations against `specs/behaviors.yaml`.
Native output extensions and capability checks belong to each backend.

The same vocabulary supports aliases, integer enums, callbacks, generic objects
and functions, static methods, variadic parameters, and composed fields. Go
rejects unsupported shapes explicitly, including first-class unions/tuples and
independent generic method parameters. Protobuf still owns shared data fields.

Client/provider signatures, the error runtime, metadata-copy helpers, and `Value`
now use these declarations. Schema getters, clone methods, validators, mappings, provider checks, and origin
helpers use declaration templates from `specs/templates.yaml`; their bodies retain
specialized emitters. Read the
[declaration guide](api-declarations.md) for syntax, coverage, and extension rules.

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
Cleanup releases SDK resources without deleting sandboxes.

These are semantic instructions, not arbitrary expressions or YAML code snippets.
Go emits constructors, context deadlines, and error returns. Another emitter must
implement equivalent native behavior, such as async calls and exceptions where
appropriate. File names, folders, imports, and constructor names stay owned by
each language emitter and its schema bindings.

The model currently resolves shared validation references, shared mapping fields,
conversion vocabulary, policy relationships, and error kinds/protocol codes.
Go-specific SDK member checks and generated validator syntax remain in Go emitters.
Handwritten provider orchestration is still required where noted below.

## Response and error contracts

All Go emitters build error expressions through `internal/go/error_gen`.
`internal/go/schema` resolves field paths against full protobuf identities.
Generated `CreateField*`, `ConfigField*`, and `InfoField*` constants replace
handwritten path strings. Nested mapping sources can declare `error_path`;
unknown paths fail generation.

Response origin validation is generated from schema presence rules. It rejects
unknown keys and absent values, while accepting explicit zero/false and empty
collections. Context enrichment avoids an extra wrapper when details already match.

`specs/contracts.yaml` identifies the shared metadata/error/origin declarations
and the maximum copy depth. Go runtime helpers are emitted with `g.P(...)` in `types_gen/runtime.go`;
generated methods copy schema fields directly, preserving numeric types and
owned collections. Metadata has explicit signed/unsigned alternatives instead
of protobuf Struct's double-only numeric representation.

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
      - {import: github.com/modal-labs/modal-client/go, name: TimeoutError, rule: timeout_error}
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
      cpu_cores: {finite: true, exclusive_minimum: 0}
    constraints:
      - {op: limit, field: cpu_limit_cores, other: cpu_cores}
      - {op: requires_positive, field: cpu_limit_cores, other: cpu_cores}
```

Field rules include `minimum`, `exclusive_minimum`, `maximum`, `finite`,
`nonblank`, `min_items`, `format: http_url`, and `each` for list elements.
`sensitive` excludes a field from generated JSON output.

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

The Go emitter uses validator tags, custom checks, and struct callbacks.
Optional fields use `omitnil`; repeated message validation rejects nil elements.
Enum membership comes from protobuf descriptors, including numeric gaps.
Validation preserves explicit zero/false and supplies no defaults.

## Provider checks

Provider `checks` validate common creation fields using schema paths:

```yaml
checks:
  - {path: resources.cpu_cores, minimum: 0.125}
  - {path: placement.cloud, allowed: [aws, gcp, oci, auto]}
  - {path: runtime.working_directory, format: absolute_path}
  - {path: network.outbound_cidrs.entries, format: cidr}
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
    source: {name: Resources}
    target:
      name: Resources
      bindings:
        go:
          import: github.com/daytona/clients/sdk-go/pkg/types
          name: Resources
          fields: {memory: Memory}
    fields:
      - {from: memory_mib, to: memory, transform: divide_exactly, factor: 1024, maximum: 2147483647, cast: integer}
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

* `settings` maps top-level configuration.
* `scope` maps scope fields.
* `auth` declares supported credential alternatives.
* `managed` permits Kit-owned provider selection and timeout fields.
* `retained` requires a matching backend state capture.
* `rejected` describes unsupported supplied configuration.
* `destination` attaches a mapped object using the native field and declared
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
2. Describe portable validation, conversions, and operation intent in YAML.
3. Add native SDK symbols to the language's bindings.
4. Extend the shared compiler vocabulary and each affected emitter for a new operation.
5. Add semantic boundary tests, regenerate, and update examples/docs.
6. Run the CLI verification command.

SDK orchestration that the vocabulary cannot express stays in handwritten
`sandbox.go`: lookup sequences, source selection, readiness, and remaining
provider semantics. Another language needs its own emitter and typed integration
code for these operations; bindings alone do not produce a complete integration.

See [tooling](../tooling/README.md), [naming](naming.md), and
[reference verification](provider-verification.md).

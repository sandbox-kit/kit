# Code generation

Protobuf and YAML define the shared contracts and portable rules. All generators
are written in Go under `tooling/`; only Go SDK output is implemented today.

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
| `specs/client.yaml`               | Provider selection and client/auth validation                       |
| `specs/validation.yaml`           | Shared creation validation                                          |
| `specs/providers/<provider>.yaml` | Client assembly, provider checks, conversions, and runtime bindings |

Specs use `version: 1`. Strict loading rejects unknown YAML keys, extra documents,
and unsupported versions. Shared fields use protobuf names; SDK symbols belong
inside the language's bindings.

## Shared validation

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
2. Describe portable validation and conversions in YAML.
3. Add native SDK symbols to the language's bindings.
4. Extend an emitter when the vocabulary needs another operation.
5. Add semantic boundary tests, regenerate, and update examples/docs.
6. Run the CLI verification command.

SDK orchestration that the vocabulary cannot express stays in handwritten
`sandbox.go`: lookup sequences, source selection, readiness, and remaining
provider semantics. Another language needs its own emitter and typed integration
code for these operations; bindings alone do not produce a complete integration.

See [tooling](../tooling/README.md), [naming](naming.md), and
[reference verification](provider-verification.md).

# Generation specifications

## Sources of truth

* `proto/`: shared message types, optional presence, oneofs, enums, local method declarations, and adapter SDK bindings.
* `specs/client.yaml`: common client configuration, runtime provider selection, and auth/client validation.
* `specs/validation.yaml`: language-neutral field and cross-field validation rules.
* `specs/providers/<provider>.yaml`: field mappings, conversions, range checks, response availability, and target-language SDK bindings.

All generators are Go programs in `tooling/`. YAML is read at generation time;
SDK users do not load it. Generated Go structs, validators, and mappings are
committed source. No protobuf imports are required in application code.

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit generate go
GOWORK=off go run ./cmd/sandbox-kit test go
```

The Cobra command runs protoc from the repository root, where the plugin resolves
`specs/`. It bootstraps annotation types, then generates runtime source.

## Validation version 1

Rules are keyed by protobuf message name and protobuf field name, not Go member
names. Types remain in protobuf and are not repeated in YAML.

```yaml
version: 1
messages:
  Resources:
    fields:
      cpu_cores: {finite: true, exclusive_minimum: 0}
    constraints:
      - {op: limit, field: cpu_limit_cores, other: cpu_cores}
```

Supported field rules: `minimum`, `exclusive_minimum`, `maximum`, `finite`,
`nonblank`, `min_items`, `format: http_url`, and `each` (apply scalar rules to list elements).
Go translates these to go-playground/validator tags and generated custom tag
validators. Optional pointers use `omitnil`, preserving explicit zero and false.
Repeated messages use `dive,required`, rejecting nil elements and validating each.

Supported cross-field operations:

| Operation | Semantics |
| --- | --- |
| `at_most_one` | At most one listed alternative is selected |
| `exactly_one` | Exactly one listed alternative is selected |
| `requires_any` | At least one listed field is supplied |
| `forbids` | None of the listed fields is selected |
| `limit` | Positive limit must be at least the reservation; zero defers to provider semantics |

A constraint can have `when: {field: mode, equals: 2}` or `not_equals`.
Conditions currently compare numeric enum/scalar values from the protobuf schema.
A selection means non-nil for optional values/messages, nonempty for collections,
true for booleans, and nonzero/nonempty for other ordinary scalars.
`sensitive: true` excludes credential fields from native Go JSON output.
Go emits struct-level validators; future languages implement the same semantics
with their own validation mechanism. Validation does not fill defaults.

## Mapping version 1

```yaml
version: 1
provider: daytona
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
      - from: memory_mib
        to: memory
        transform: divide_exactly
        factor: 1024
        maximum: 2147483647
        cast: integer
```

Shared types resolve against the protobuf descriptors. Provider types use logical
field keys mapped to native members under `bindings.go`. A new SDK language adds
its own bindings and emitter; the field rules and conversion semantics stay shared.
Types/field names inside a language binding are SDK-specific, not portable rules.

Supported conversions: copy, `whole`, `divide_exactly`, and `multiply`.
`maximum` bounds the converted value. Division rejects remainders; multiplication
checks unsigned overflow and negative input. Portable casts are `integer`,
`number`, `unsigned`, and `string`; Go chooses the corresponding native type.

Request helpers preserve nil optional fields and return a remaining configuration
with mapped fields cleared. backends reject remaining unsupported intent before
provider calls. Response helpers return native shared structs and set provider
identity. A group's `when` list supports numeric `greater_than` conditions; all
must hold, otherwise no fields are mapped. Daytona uses this for allocated-resource
metadata availability.

Strict YAML decoding rejects unknown keys, extra documents, and unsupported
versions. Generators reject unknown protobuf fields, missing SDK member bindings,
unknown operations/casts, zero conversion factors, and duplicate destinations.
The compiler and adapter tests verify native SDK member types against pinned SDKs.

## Provider runtime bindings

`runtime.go.state` lists only the state fields that an adapter actually retains.
Modal declares app/environment scope and region; Daytona has no extra state beyond
its SDK client. Each state entry has a native field name, shared type (or `string`),
and optional pointer flag. Other languages can define their own runtime bindings.

`runtime.go.cleanup` declares the SDK cleanup method, whether it accepts context,
and whether it returns an error. The adapter emitter generates the typed SDK call
directly; separate handwritten cleanup forwarding functions are not needed.

## Current coverage and extensions

Common client and creation validation are generated. Config's native
Provider field requires Provider and is not serialized. Go emits
NewClient(Config), SDK initialization delegation, default timeout handling,
and Close delegation. Constructor names and initialisms are Go emitter choices.

Generated mappings cover both
providers' name/environment/labels, resource requests and precision checks, plus
Daytona sandbox metadata and allocated-resource conversions. They also cover
Modal token-pair/OAuth credentials and Daytona API-key/bearer credentials,
endpoint, region, and organization scope. Initialization helpers handle provider
support checks and typed SDK constructors.

Provider orchestration and remaining semantic mappings stay in typed `create.go`
files: app/image/secret resolution, source choice, policy semantics, networking,
readiness, and SDK invocation. These are not automatically portable yet. Do not
claim that adding a language binding alone generates a complete adapter.

Extend the mapping vocabulary with a reviewed operation when moving another
semantic rule into YAML. Add tests for its boundary values and unsupported intent,
then implement that operation in each target emitter. Avoid raw code snippets in
YAML: they prevent portable generation. Keep SDK calls as typed language bindings
until an explicit call-sequence model is introduced.

Harness adapters can later reuse this spec infrastructure; no harness generator
or non-Go SDK is implemented today.

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

Supported conversions: copy, `whole`, `divide_exactly`, `multiply`,
`scaled_integer`, `duration_seconds`, and `policy_minutes`.
`maximum` bounds the converted value. Division rejects remainders; multiplication
checks unsigned overflow and negative input. Portable casts are `integer`,
`number`, `unsigned`, and `string`; Go chooses the corresponding native type.

`scaled_integer` requires a nonnegative finite shared double exactly representable
at `factor` units per value; `maximum` bounds the scaled integer. It keeps the
native float value and compensates for binary rounding before SDK truncation.
Modal CPU uses factor 1000 and a uint32 maximum. `duration_seconds` requires a
positive whole-second duration, bounds its seconds with `maximum`, and retains
the native duration. Both operations preserve absent values and reject lossy
input before SDK calls. Their target-language emitters must reproduce these
semantics; SDK resource quotas still apply separately.

`policy_minutes` maps an `AutomaticAction` to a native optional integer interval.
Its `policy.disabled` declares `zero` or `reject`; `policy.immediate` declares
whether AFTER zero represents an immediate action. Absent/default stays omitted.
Delayed actions require representable whole minutes. Daytona declares these rules
per action. Stop/pause zero disables them; only delete zero requests immediate
action. Archive's maximum is declared in minutes. The service guide governs
semantics where SDK parameter comments conflict; see [verification](provider-verification.md).

Provider `checks` declare schema paths with numeric bounds, allowed string values,
and formats (absolute paths, HTTP URLs, CIDRs, domains, environment keys). Checks
preserve omitted values. Policy relationships reject conflicting idle delays and
ephemeral combinations before native SDK calls. Enum membership is generated from
protobuf descriptors. Resource constraints require positive limits to have a
corresponding request. Account quotas and placement capacity stay provider-enforced.

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

`runtime.go.constructor.function` names the SDK constructor taking a pointer to
the mapped configuration and returning a client plus error. The emitter generates
`newBackend`: normalize omitted config, map/validate it, call the constructor,
preserve constructor errors, and capture configured backend state. State `capture`
is a schema field path relative to Config (for example `region`, `scope`, or
`scope.organization_id`), not a Go expression. Generated typed copies preserve
nil and explicit values without aliasing caller data. Message copies follow all
schema fields, including organization/project scope. The current emitter supports
optional strings and messages containing scalar/nested message fields; collections,
bytes, recursive messages, and real oneofs fail generation until copying support
is added. Invalid paths and incompatible capture types also fail generation.
A new language emitter implements these ownership semantics and binds its native
constructor.

## Current coverage and extensions

Common client and creation validation are generated. Config's native
Provider field requires Provider and is not serialized. Go emits
NewClient(Config), SDK initialization delegation, default timeout handling,
and Close delegation. Constructor names and initialisms are Go emitter choices.

Generated mappings cover both
providers' name/environment/labels, resource requests and precision checks, plus
Daytona sandbox metadata and allocated-resource conversions. They also cover
Modal token-pair/OAuth credentials and Daytona API-key/bearer credentials,
endpoint, region, and organization scope. Configuration assembly, provider support
checks, native SDK constructors, and state capture are generated.

The provider's `client` spec composes request mapping groups into initialization
parameters. `settings` maps top-level config; `scope` maps scope; `auth` declares
supported credential variants. Components reference a mapping `group`, and can
declare additional `retained` fields captured by the backend. A `destination`
attaches a mapped object by pointer using a named native binding and its declared
`objects` type (Modal OAuth credentials). `managed` permits only Kit's provider
selection and timeout; `retained` requires a matching state capture. Mapped fields
and auth/scope components own their input automatically, while
`rejected` identifies unsupported fields with diagnostic messages. Unknown fields,
missing groups, incompatible full schema/native identities, conflicting ownership,
and duplicate native destinations fail generation. Exclusive auth alternatives
may share destinations with each other, but may not overwrite common settings or
scope. Overrides are rejected; no implicit precedence exists. The resulting
`provider.client.gen.go` validates common configuration,
rejects unmapped intent, and leaves omitted credentials to SDK defaults.

Provider orchestration and remaining semantic mappings stay in typed `sandbox.go`
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

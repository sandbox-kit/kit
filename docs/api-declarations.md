# API declarations and language profiles

Runtime API declarations live in [`specs/api.yaml`](../specs/api.yaml).
Shared data fields, presence, serialized names, and enum values stay in protobuf.
The shared compiler validates declarations and behavior compatibility before a
language emitter renders them.

## Describe a method

```yaml
- id: create
  name: Create
  behavior: create
  receiver: c
  cancellable: true
  params:
    - id: request
      name: request
      type: {ref: schema.CreateOptions, reference: true, nullable: true}
  results:
    - id: instance
      name: instance
      type: {ref: sandbox_handle, reference: true}
  fallible: true
  named_results: true
```

`id` connects the declaration to behavior instructions. `name` is its emitted
name. Parameter IDs let the Go behavior emitter resolve explicit parameter references
when names change. The declaration above generates:

```go
func (c *Client) Create(ctx context.Context, request *CreateOptions) (instance *Sandbox, result error)
```

Names, parameters, results, fields, and receiver names come from the declaration.
The language profile supplies context injection and the extra error result.
Behavior bodies are rendered by language emitters using the shared instruction
vocabulary; YAML does not contain executable Go snippets.

## Declaration vocabulary

| Declaration | Meaning | Go representation |
| --- | --- | --- |
| `object` | Runtime object with fields and methods | Struct |
| `interface` | Method contract without storage | Interface |
| `alias` | Alternate name for another type | Type alias |
| `enum` | Named integer type and numeric members | Defined type and constants |
| `callback` | Named function signature | Function type |
| Constructor | Builds a declared object; may fail | Factory function |
| Static method | Operation associated with a type | Package function named `TypeMethod` |
| Generic function or object | Type parameters with constraints | Go generics |
| Variadic parameter | Final sequence of arguments | `...T` |
| Promoted field | Exposes composed object's fields | Embedded named field |

Runtime-only `attachments` add native fields to a schema-backed object. For example,
`Config.Provider` is declared in API YAML and excluded from serialization. Shared
serialized fields remain in protobuf.

Object fields are instance variables. Constructors declare `kind: constructor`
and `constructs: <object-id>`. `effect: pure|io` records execution intent for
future language emitters; current Go operations execute synchronously.
`doc` becomes a generated Go comment. `fields_from: provider_state` incorporates
provider state from the existing provider YAML without duplicating storage rules.

Type references support builtins, named references, generic arguments, sequences,
maps, function types, tuples, and unions. A type reference selects exactly one
shape. `schema.<Message>` refers to protobuf; language externals refer to native
bindings in the profile. Unknown or unsupported references fail generation.

## Presence, identity, and ownership

These are independent meanings:

- `optional`: a value can be omitted.
- `nullable`: explicit null is permitted.
- `reference`: identity or shared access is required.
- `ownership: owned|borrowed|shared`: how the value is held.

The Go profile represents optional scalars and object references with pointers.
Interfaces, callbacks, maps, and slices already have native nil representations.
It rejects combined optional-and-nullable slots because a single nil cannot
preserve both states. Ownership metadata does not itself introduce cloning or
cleanup; the associated behavior instructions must implement it.

A union describes alternatives, not error handling. `fallible` independently
expresses that an operation can fail. Several successful results are declared
through `results`, not by treating errors as another success alternative.

## Language profiles and output

[`specs/languages/go.yaml`](../specs/languages/go.yaml) selects native object and
interface forms, reference/presence rules, error results, multiple-return syntax,
cancellation parameters, builtins, external symbols, and output suffixes.
The shared loader checks safe output suffixes without imposing Go extensions.
Go validates its required output roles and `.gen.go` suffixes in the backend.
A `.gen.ts` profile can be loaded without implementing a TypeScript emitter.

Output remains inside the directories declared by the Cobra generation command
and protobuf package bindings. Profiles cannot redirect output through suffixes.

Go supports multiple return values but has no first-class tuple or union value
type. The emitter rejects those type shapes explicitly; use a named result object
or a tagged schema alternative. Go also rejects independent generic method
parameters; methods can use their generic receiver's parameters.

For a future TypeScript emitter, optional parameters, `undefined`, `null`, unions,
callbacks, and `Promise` results need their own lowering rules. TypeScript unions
are alternatives in a single value, while tuple or named-object results can carry
several values. A fallible asynchronous initializer needs a factory convention.
These distinctions follow the official [function reference](https://www.typescriptlang.org/docs/handbook/2/functions.html)
and [type reference](https://www.typescriptlang.org/docs/handbook/2/everyday-types.html).
Only Go emission is implemented and verified today.

## Current migration coverage

Client/provider objects, interfaces, constructors, methods, fields, parameters,
and result signatures are declared in YAML. The common error object, its helpers,
metadata-copy helpers, and generic `Value` helper also use this declaration model.

Protobuf annotations now identify source contracts and provider SDK bindings.
Their retired client/backend naming fields are reserved; there is no naming
fallback to the previous generator implementation.

Descriptor-driven signatures for getters, clone methods, validators, provider
mappings/checks, client configuration mappings, and response-origin helpers now
come from [`specs/templates.yaml`](../specs/templates.yaml). Each template binds
names and types from a schema descriptor or compiled mapping; a YAML entry for
every field or provider mapping is unnecessary.

Field-path traversal roots, constant name patterns, native leaves, and excluded
origin fields also come from this specification. Constant references in provider
mappings use the same rules as the constant declarations.

Shared data storage, validator callbacks, conversions, and copy/presence logic
remain specialized emitter code. Provider orchestration and the sandbox handle
remain native integration code. This vocabulary does not claim to represent every
Go construct.

## Descriptor-driven templates

```yaml
getter:
  owner: subject
  declaration:
    id: getter
    name: 'Get{field}'
    receiver: x
    behavior: getter
    results: [{type: {ref: field_value}}]
```

The shared compiler validates a pattern, then instantiates its name with
`field: CPUCores`. The Go adapter binds `subject` to `Resources` and `field_value`
to the descriptor's getter return type. Bound types retain portable shape references and selected native symbol identities;
Go type strings are not parsed back into syntax. YAML contains no Go expressions. Missing bindings and invalid names fail generation.

A request-mapping template declares the source parameter, native target result,
remaining request result, and possible failure. Go lowers these to multiple
return values. Another language can implement its own result-object convention.
Body references explicitly select parameters, results, receivers, fields, or local
variables. The emitter does not scan or rename body text. A local variable and a
parameter can share a canonical spelling without being confused.

## Behavior contracts

[`specs/behaviors.yaml`](../specs/behaviors.yaml) defines each supported behavior's
owner category, ordered parameter/result shapes, generic inputs, failure intent,
and cancellation intent. The compiler checks API declarations and reusable
templates against those requirements before emission. Native names can change;
required IDs and type shapes must remain compatible.

For example, creation requires a creation request, a sandbox result, failure
handling, and cancellation. Removing the request, changing it to a string, adding
an unsupported argument, or dropping failure handling is rejected during compilation.
Contract changes also require matching behavior implementations in each backend.

Go-specific requirements remain in the Go backend. For example, deferred error
context needs named Go results; a future language can implement equivalent
context enrichment through its own calling conventions.

## Structured bindings and explicit references

`TemplateBindings.Types` carries `TypeRef` values. Native SDK types retain their
import path, name, and kind in a separate external-symbol table. Schema field
projection preserves numeric widths, enum/message identities, container shapes,
and getter presence semantics before Go lowering.

Behavior emitters resolve parameter/receiver references explicitly, for example
`d.Param(callable, "request")` and `d.Receiver(callable)`. Fields are resolved with
a declared owner and field ID; locals have independent identities. Missing
references fail generation. Raw body fragments are written verbatim and never
searched for identifiers to replace.

## Change and verify

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit generate go
GOWORK=off go run ./cmd/sandbox-kit test go
```

Edit declarations, profiles, schemas, or behavior rules, then regenerate.
Declaration tests type-check representative aliases, enums, callbacks, generics,
and multiple-result functions. Runtime/provider/example tests verify the current
SDK behavior. API changes also require matching examples and integration bindings.

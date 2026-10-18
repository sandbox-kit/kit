# API declarations and language profiles

Public portable types and service metadata live in protobuf under `proto/`.
SDK runtimes alias those generated contract types and keep constructors,
facades, and orchestration handwritten.

YAML no longer declares a parallel public API surface. Remaining specs describe
provider rules, validation, mappings, templates, and language profiles:

| Spec                                             | Role                                                                |
| ------------------------------------------------ | ------------------------------------------------------------------- |
| `specs/generation.yaml`                          | Remaining Kit emission phases (paths, validation, checks, mappings) |
| `specs/templates.yaml`                           | Helper signatures and field-path rules for those emitters           |
| `specs/behaviors.yaml`                           | Signature checks for helper templates still emitted                 |
| `specs/languages/<language>.yaml`                | Native representation and output conventions                        |
| `specs/providers/*.yaml`                         | Provider mappings, checks, and bindings                             |
| `specs/validation.yaml` / `specs/contracts.yaml` | Shared validation and contract limits                               |

Error payloads (`ErrorInfo` / `ErrorKind`) live in protobuf. Wrappers and
provider classification are handwritten per language.

## Language profiles

[`specs/languages/go.yaml`](../specs/languages/go.yaml) and
[`specs/languages/typescript.yaml`](../specs/languages/typescript.yaml) control
builtins, externals, and output suffixes. Duration maps to `time.Duration` in Go
and protobuf `Duration` in TypeScript (`duration(ms)` helpers build values for
callers).

## What stays handwritten

- `Client` / `Sandbox` facades and lifecycle dispatch
- Provider selection, native SDK initialization, and cleanup
- Error wrappers, cause chaining, and provider-native error classification
- Metadata clone helpers and request materialization at language boundaries
- Provider orchestration that the mapping vocabulary cannot express

See [Protobuf SDK boundary](protobuf-sdk-boundary.md) and
[code generation](code-generation.md) for the full split between contracts,
generated helpers, and handwritten runtime code.

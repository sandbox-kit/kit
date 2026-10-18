# API declarations

Public portable types and service metadata live in protobuf under `proto/`.
SDK runtimes alias those generated contract types and keep constructors,
facades, and orchestration handwritten.

YAML does not declare a parallel public API surface. Remaining specs inventory
provider rules, validation, and contracts:

| Spec                                          | Role                                            |
| --------------------------------------------- | ----------------------------------------------- |
| `specs/providers/*.yaml`                      | Provider mapping/check inventories and bindings |
| `specs/validation.yaml` / `specs/client.yaml` | Portable create/config validation inventories   |
| `specs/contracts.yaml`                        | Portable metadata/error contract limits         |

Error payloads (`ErrorInfo` / `ErrorKind`) live in protobuf. Wrappers and
provider classification are handwritten per language.

## What stays handwritten

- `Client` / `Sandbox` facades and lifecycle dispatch
- Provider selection, native SDK initialization, and cleanup
- Error wrappers, cause chaining, and provider-native error classification
- Metadata clone helpers and request materialization at language boundaries
- Shared validation, provider checks, field mappings, and config assembly
- Provider orchestration that the mapping inventory cannot express

See [Protobuf SDK boundary](protobuf-sdk-boundary.md) and
[code generation](code-generation.md) for the full split between contracts and
handwritten runtime code.

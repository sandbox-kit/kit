# API declarations

Public portable types live in protobuf under `proto/`. SDK runtimes alias those
generated contract types and keep constructors, facades, provider backends, and
orchestration handwritten. TypeScript exposes handwritten public create/config
facades; see [protobuf SDK boundary](protobuf-sdk-boundary.md).

YAML does not declare a parallel public API surface. Remaining specs inventory
provider rules, validation, and contracts:

| Spec                                          | Role                                            |
| --------------------------------------------- | ----------------------------------------------- |
| `specs/providers/*.yaml`                      | Provider mapping/check inventories and bindings |
| `specs/validation.yaml` / `specs/client.yaml` | Portable create/config validation inventories   |
| `specs/contracts.yaml`                        | Portable metadata/error contract limits         |

Error payloads (`ErrorInfo` / `ErrorKind`) live in protobuf. Wrappers and
provider classification are handwritten per language.

## Handwritten runtime

- `Client` / `Sandbox` facades and lifecycle dispatch
- TypeScript public create/config input types and materialization to wire
- Provider selection, native SDK initialization, and cleanup
- Error wrappers, cause chaining, and provider-native error classification
- Metadata clone helpers
- Protovalidate evaluation for shared create/config rules (`secure_endpoint` is
  handwritten)
- Provider checks, field mappings, and config assembly
- Provider orchestration that the mapping inventory cannot express

See [Protobuf SDK boundary](protobuf-sdk-boundary.md) and
[code generation](code-generation.md) for the full split between contracts and
handwritten runtime code.

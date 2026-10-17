# API declarations and language profiles

Public portable types and service metadata live in protobuf under `proto/`.
SDK runtimes alias those generated contract types and keep constructors,
facades, and orchestration handwritten.

YAML no longer declares a parallel public API surface. Remaining specs describe
provider rules, validation, mappings, errors, templates, and language profiles:

| Spec | Role |
| --- | --- |
| `specs/generation.yaml` | Semantic generation phases and client operation vocabulary |
| `specs/templates.yaml` | Reusable descriptor-driven helper signatures and field-path rules |
| `specs/behaviors.yaml` | Behavior compatibility checks before emission |
| `specs/languages/<language>.yaml` | Native representation and output conventions |
| `specs/providers/*.yaml` | Provider mappings, checks, and bindings |
| `specs/validation.yaml` / `specs/errors.yaml` / `specs/contracts.yaml` | Shared validation, error policy, and contract limits |

## Language profiles

[`specs/languages/go.yaml`](../specs/languages/go.yaml) and
[`specs/languages/typescript.yaml`](../specs/languages/typescript.yaml) control
builtins, externals, and output suffixes. Duration maps to `time.Duration` in Go
and protobuf `Duration` in TypeScript (`duration(ms)` helpers build values for
callers).

## What stays handwritten

- `Client` / `Sandbox` facades and lifecycle dispatch
- Provider selection, native SDK initialization, and cleanup
- Metadata clone helpers and request materialization at language boundaries
- Provider orchestration that the mapping vocabulary cannot express

See [Protobuf SDK boundary](protobuf-sdk-boundary.md) and
[code generation](code-generation.md) for the full split between contracts,
generated helpers, and handwritten runtime code.

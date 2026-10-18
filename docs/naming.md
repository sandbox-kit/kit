# Naming and package boundaries

Use names that read clearly at the call site. Public Go APIs follow Go
[package naming](https://go.dev/blog/package-names) and
[initialism conventions](https://go.dev/wiki/CodeReviewComments#initialisms).

## Public vocabulary

| Concept                 | Go name               | Responsibility                                       |
| ----------------------- | --------------------- | ---------------------------------------------------- |
| Public SDK              | `sandbox`             | Configuration, lifecycle, and shared metadata        |
| Provider selection      | `Provider`            | `Name()` and `NewClient(*Config)`                    |
| Initialized integration | `Backend`             | Provider identity and local SDK cleanup              |
| Optional creation       | `SandboxCreator`      | `Create()` when the provider supports creation       |
| Optional lifecycle      | Capability interfaces | Typed provider operations such as start or terminate |
| Public client           | `Client`              | Owns the backend and delegates operations            |
| Client settings         | `Config`              | Provider, auth, scope, endpoint, region, timeout     |
| Creation request        | `CreateOptions`       | Source, runtime, resources, and creation policies    |
| Backend result          | `CreateResult`        | Shared sandbox metadata                              |
| Public handle           | `Sandbox`             | Identity, copied metadata, and lifecycle methods     |

Public portable types come from protobuf contracts (`kit.sandbox.v1`) aliased in
the SDK runtime. Descriptor-driven getter, validation, mapping, and origin
signatures come from `specs/templates.yaml`. `specs/languages/go.yaml` owns Go
representation and output suffixes. Protobuf annotations identify source
contracts and provider SDK bindings.

Provider packages expose `New() *Provider`. The initialized `backend` and native
SDK client remain private. Construction uses `sandbox.NewClient(sandbox.Config{...})`.

## Files and directories

```text
sdks/go/sandbox/            Public SDK module
sdks/go/providers/modal/    Optional Modal module
sdks/go/providers/daytona/  Optional Daytona module
proto/kit/sandbox/v1/       Shared contracts
proto/kit/providers/        Thin provider descriptor entrypoints
specs/                     Portable validation, mapping, and emit rules
tooling/internal/model/    Language-neutral compiled generation model
tooling/internal/go/       Go phase dispatcher and emitters
examples/go/<provider>/      One folder per provider, one module per example
sdks/typescript/sandbox/     Public TypeScript SDK package
sdks/typescript/providers/   Optional TypeScript provider packages
tooling/internal/typescript/ TypeScript phase dispatcher and emitters
examples/typescript/<provider>/ One package per example
```

Emitter directories use the `*_gen` suffix. Kit-generated helpers end in `.gen.go`
(or `.gen.ts`). Portable messages and enums are not re-emitted by kit generators;
they come from contract packages plus handwritten aliases such as
`sandbox_data_contracts.go` / `sandbox_data_contracts.ts`.

| File                         | Purpose                                            |
| ---------------------------- | -------------------------------------------------- |
| `*.validation.gen.go`        | Shared request validation                          |
| `*.errors.gen.go`            | Diagnostic field names and provider classification |
| `*.paths.gen.go`             | Canonical field-path constants                     |
| `provider.mappings.gen.go`   | Request and metadata field conversions             |
| `provider.validation.gen.go` | Provider-specific create checks                    |
| `provider.config.gen.ts`     | Native SDK configuration assembly                  |
| `provider.checks.gen.ts`     | Provider-specific create checks                    |
| `provider.mappings.gen.ts`   | Request and metadata field conversions             |

Handwritten runtime and provider files own client/orchestration behavior and
metadata clone helpers. Change the schema, spec, or emitter when changing
generated helpers.

## Language conventions

Field paths have generated names such as `CreateFieldResourcesCPUCores`,
`ConfigFieldScopeAppName`, and `InfoFieldResources`. Their values retain
canonical protobuf field names, including nested paths.

Protobuf and YAML use snake_case semantic field names. The Go emitter produces
names such as `CPUCores`, `MemoryMiB`, `OutboundCIDRs`, `TCPPort`, and `APIKey`.
Enums use typed constants such as `IsolationKindLinuxVM` and `PolicyModeDisabled`.
Native SDK members use the names from their pinned bindings.

The TypeScript emitter uses camelCase fields, discriminated unions for exclusive
alternatives, interfaces for named object contracts, type aliases for unions,
and generated files ending in `.gen.ts`. Preserve numeric enum values and
explicit serialized names in every language: the Go fields `Scope` and
`Provisioning` serialize as `context` and `creation`.

See [provider architecture](providers.md) and [generation rules](code-generation.md).

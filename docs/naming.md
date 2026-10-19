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
the SDK runtime. Field-path constants use prefixes such as `CreateField`,
`ConfigField`, and `InfoField`. Protobuf annotations identify source contracts
and operation metadata.

Provider packages expose `New() *Provider`. The initialized `backend` and native
SDK client remain private. Construction uses
`sandbox.NewClient(provider, &sandbox.Config{...})` in Go and
`new Client({ provider, ... })` in TypeScript.

## Files and directories

```text
sdks/go/sandbox/            Public SDK module
sdks/go/providers/modal/    Optional Modal module
sdks/go/providers/daytona/  Optional Daytona module
proto/kit/sandbox/v1/       Shared contracts
specs/                     Portable validation and mapping inventories
tooling/                   Developer CLI and YAML inventory loaders
examples/go/<provider>/      One folder per provider, one module per example
sdks/typescript/sandbox/     Public TypeScript SDK package
sdks/typescript/providers/   Optional TypeScript provider packages
examples/typescript/<provider>/ One package per example
```

Portable messages and enums come from standard protobuf contract packages plus
handwritten aliases such as `sandbox_data_contracts.go` /
`sandbox_data_contracts.ts`. Kit custom emitters no longer write SDK helpers.

| File                                            | Purpose                                        |
| ----------------------------------------------- | ---------------------------------------------- |
| `create_validation.go`                          | Shared CreateOptions validation                |
| `validation_fields.go`                          | Go↔protobuf field-name maps for validators     |
| `field_paths.go` / `config_paths.go`            | Canonical field-path constants                 |
| `provider_mappings.go` / `provider_mappings.ts` | Request and metadata field conversions         |
| `provider_checks.go` / `provider_checks.ts`     | Provider-specific create checks                |
| `provider_config.ts`                            | Native SDK configuration assembly (TypeScript) |

Handwritten runtime and provider files own client/orchestration behavior and
metadata clone helpers. Keep YAML inventories aligned when changing those files.

## Language conventions

Field-path constants use names such as `CreateFieldResourcesCpuCores`,
`ConfigFieldScopeAppName`, and `InfoFieldResources`. Their values retain
canonical protobuf field names, including nested paths.

Protobuf and YAML use snake_case semantic field names. Generated Go contract
types follow protobuf Go naming (`CpuCores`, `MemoryMib`, `OutboundCidrs`,
`TcpPort`, `ApiKey`). Enums use typed constants such as `IsolationKindLinuxVM`
and `PolicyModeDisabled`. Native SDK members use the names from their pinned
bindings.

TypeScript SDK types use camelCase fields, discriminated unions for exclusive
alternatives, interfaces for named object contracts, and type aliases for
unions. Preserve numeric enum values and explicit serialized names in every
language: the Go fields `Scope` and `Provisioning` serialize as `context` and
`creation`.

See [provider architecture](providers.md) and [generation rules](code-generation.md).

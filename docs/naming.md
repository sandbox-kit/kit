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
the SDK runtime. TypeScript app-facing `CreateOptions` / `Config` are handwritten
facades (ms timeouts, nested oneofs); wire messages stay private to the client.
Field-path constants use prefixes such as `CreateField`, `ConfigField`, and
`InfoField`. Protobuf annotations identify source contracts and operation
metadata.

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
sdks/typescript/contracts/   Generated protobuf wire types (@sandbox-kit/contracts)
sdks/typescript/sandbox/     Public TypeScript SDK (@sandbox-kit/sandbox)
  src/client/                Client runtime split across focused modules:
                               client.ts (Client + public re-exports),
                               sandbox.ts (handle), types.ts, errors.ts,
                               validation.ts, operation.ts (deadlines/signals)
  src/contracts/             Public type aliases over protobuf messages
  src/metadata/              Metadata map helpers
  src/internal/              Package-private duration/oneof/strip helpers
sdks/typescript/providers/   Optional providers (@sandbox-kit/daytona, @sandbox-kit/modal)
  src/mappings/              Create/info field mappings
  src/checks/                Provider create checks
  src/config/                Native SDK client configuration
examples/typescript/<provider>/ One package per example
```

Published TypeScript packages expose only their root export (`.`). Cross-folder
imports inside a package use private `#` paths (`#internal/...`, `#mappings/...`)
so apps never depend on package subfolders. Providers depend on
`@sandbox-kit/sandbox` only.

Portable messages and enums come from standard protobuf contract packages plus
handwritten aliases such as `sandbox_data_contracts.go` /
`sandbox_data_contracts.ts`. SDK helpers are handwritten.

| File                                                     | Purpose                                      |
| -------------------------------------------------------- | -------------------------------------------- |
| `create_validation.go`                                   | Shared CreateOptions validation              |
| `support.go`                                             | `RejectUnmapped` (protobuf tag → field path) |
| `field_paths.go` / `config_paths.go`                     | Canonical field-path constants (Go)          |
| `field_paths.ts` / `config_paths.ts`                     | Same protobuf paths for TS `ErrorInfo.field` |
| `provider_mappings.go` / `mappings/provider_mappings.ts` | Request and metadata field conversions       |
| `provider_checks.go` / `checks/provider_checks.ts`       | Provider-specific create checks              |
| `config/provider_config.ts`                              | Native SDK configuration assembly (TS)       |
| `capabilities.go`                                        | Optional backend operation interfaces (Go)   |
| `creation.go`                                            | Provider create planning (Go)                |
| `lifecycle.go`                                           | Provider lifecycle calls (Go)                |
| `metadata.go`                                            | Observed sandbox metadata (Go)               |
| `operation.go`                                           | Sandbox handle operation gate (Go)           |

Handwritten runtime and provider files own client/orchestration behavior and
metadata clone helpers. Go providers keep SDK calls in `sandbox.go` and the
create plan in `creation.go`. Daytona also keeps lifecycle calls in
`lifecycle.go` and observed metadata in `metadata.go`. Keep YAML inventories
aligned when changing those files.

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

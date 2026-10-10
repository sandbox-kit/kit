# Naming and package boundaries

Use names that read clearly at the call site. Public Go APIs follow Go
[package naming](https://go.dev/blog/package-names) and
[initialism conventions](https://go.dev/wiki/CodeReviewComments#initialisms).

## Public vocabulary

| Concept                 | Go name         | Responsibility                                    |
| ----------------------- | --------------- | ------------------------------------------------- |
| Public SDK              | `sandbox`       | Configuration, creation, and shared metadata      |
| Provider selection      | `Provider`      | `Name()` and `NewClient(*Config)`                 |
| Initialized integration | `Backend`       | `Name()`, `Create()`, and `Close()`               |
| Public client           | `Client`        | Owns the backend and delegates operations         |
| Client settings         | `Config`        | Provider, auth, scope, endpoint, region, timeout  |
| Creation request        | `CreateOptions` | Source, runtime, resources, and creation policies |
| Backend result          | `CreateResult`  | Shared sandbox metadata                           |
| Public handle           | `Sandbox`       | Identity and copied `SandboxInfo`                 |

Runtime declaration names come from `specs/api.yaml`. Descriptor-driven getter,
clone, validation, mapping, and origin signatures come from `specs/templates.yaml`.
`specs/languages/go.yaml` owns Go representation and output suffixes. Protobuf
annotations identify source contracts and provider SDK bindings.

Provider packages expose `New() *Provider`. The initialized `backend` and native
SDK client remain private. Construction uses `sandbox.NewClient(sandbox.Config{...})`.

## Files and directories

```text
sdks/go/sandbox/            Public SDK module
sdks/go/providers/modal/    Optional Modal module
sdks/go/providers/daytona/  Optional Daytona module
proto/kit/sandbox/v1/       Shared contracts
proto/kit/providers/        Provider declarations
specs/                     Portable API, behavior, and mapping rules
tooling/internal/model/    Language-neutral compiled generation model
tooling/internal/go/       Go phase dispatcher and emitters
examples/go/               Independent example modules
```

Emitter directories use the `*_gen` suffix. Generated files end in `.gen.go`:

| File                       | Purpose                                            |
| -------------------------- | -------------------------------------------------- |
| `client.gen.go`            | Public client and provider/backend contracts       |
| `*.types.gen.go`           | Native shared types                                |
| `*.validation.gen.go`      | Shared or provider validation                      |
| `*.errors.gen.go`          | Diagnostic field names and provider classification |
| `*.paths.gen.go`           | Canonical field-path constants                     |
| `*.copy.gen.go`            | Owned schema copies                                |
| `*.errors.runtime.gen.go`  | Error object, constructor, and context helpers     |
| `*.clone.runtime.gen.go`   | Metadata copy helpers                              |
| `provider.gen.go`          | Provider construction, state capture, and cleanup  |
| `provider.client.gen.go`   | Native SDK configuration assembly                  |
| `provider.mappings.gen.go` | Request and metadata field conversions             |

Output suffixes for client, provider, error-runtime, and metadata-copy files
come from `specs/languages/go.yaml`. Handwritten `sandbox.go` files contain
remaining provider creation orchestration. Change the schema, spec, or emitter
when changing generated behavior.

## Language conventions

Field paths have generated names such as `CreateFieldResourcesCPUCores`,
`ConfigFieldScopeAppName`, and `InfoFieldResources`. Their values retain
canonical protobuf field names, including nested paths.

Protobuf and YAML use snake_case semantic field names. The Go emitter produces
names such as `CPUCores`, `MemoryMiB`, `OutboundCIDRs`, `TCPPort`, and `APIKey`.
Enums use typed constants such as `IsolationKindLinuxVM` and `PolicyModeDisabled`.
Native SDK members use the names from their pinned bindings.

Future language emitters apply their own conventions to the same contracts.
Preserve numeric enum values and explicit serialized names: the Go fields
`Scope` and `Provisioning` serialize as `context` and `creation`.

See [provider architecture](providers.md) and [generation rules](code-generation.md).

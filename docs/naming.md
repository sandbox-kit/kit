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

Provider packages expose `New() *Provider`. The initialized `backend` and native
SDK client remain private. Construction uses `sandbox.NewClient(sandbox.Config{...})`.

## Files and directories

```text
sdks/go/sandbox/            Public SDK module
sdks/go/providers/modal/    Optional Modal module
sdks/go/providers/daytona/  Optional Daytona module
proto/kit/sandbox/v1/       Shared contracts
proto/kit/providers/        Provider declarations
specs/                     Portable validation and mapping rules
tooling/internal/go/       Go emitters
examples/go/               Independent example modules
```

Emitter directories use the `*_gen` suffix. Generated files end in `.gen.go`:

| File                       | Purpose                                           |
| -------------------------- | ------------------------------------------------- |
| `client.gen.go`            | Public client and provider/backend contracts      |
| `*.types.gen.go`           | Native shared types                               |
| `*.validation.gen.go`      | Shared or provider validation                     |
| `provider.gen.go`          | Provider construction, state capture, and cleanup |
| `provider.client.gen.go`   | Native SDK configuration assembly                 |
| `provider.mappings.gen.go` | Request and metadata field conversions            |

Handwritten `sandbox.go` files contain remaining provider creation orchestration.
Change the schema, spec, or emitter when changing generated behavior.

## Language conventions

Field paths have generated names such as `CreateFieldResourcesCPUCores`,
`ConfigFieldScopeAppName`, and `InfoFieldResources`. Their values retain
canonical protobuf field names, including nested paths.

Response/error output also includes `*.copy.gen.go`, `*.paths.gen.go` field-path constants,
provider `*.errors.gen.go` classifiers, and `*.runtime.gen.go` helpers.
Runtime helpers are emitted by `types_gen/runtime.go` using `specs/contracts.yaml` and `specs/errors.yaml`.

Protobuf and YAML use snake\_case semantic field names. The Go emitter produces
names such as `CPUCores`, `MemoryMiB`, `OutboundCIDRs`, `TCPPort`, and `APIKey`.
Enums use typed constants such as `IsolationKindLinuxVM` and `PolicyModeDisabled`.
Native SDK members use the names from their pinned bindings.

Future language emitters apply their own conventions to the same contracts.
Preserve numeric enum values and explicit serialized names: the Go fields
`Scope` and `Provisioning` serialize as `context` and `creation`.

See [provider architecture](providers.md) and [generation rules](code-generation.md).

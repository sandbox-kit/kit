# Naming and package boundaries

Names are chosen for how the public API reads at a call site. Go guidance
recommends short, meaningful [package names](https://go.dev/blog/package-names)
and consistent [initialisms](https://go.dev/wiki/CodeReviewComments#initialisms).
These conventions do not prescribe one universal name for every SDK abstraction.

## Public API

```go
client, err := sandbox.NewClient(sandbox.Config{
    Provider: modal.New(),
    Auth: &sandbox.AuthConfig{
        TokenPair: &sandbox.TokenPairCredentials{ID: tokenID, Secret: tokenSecret},
    },
    Scope: &sandbox.Scope{AppName: sandbox.Value("my-app")},
})
```

| Concept | Name | Responsibility |
| --- | --- | --- |
| Public SDK | `sandbox` | Client configuration and sandbox operations |
| Selected integration | `Provider` | `Name()` and `NewClient(*Config)` |
| Initialized integration | `Backend` | `Name()`, `Create()`, and `Close()` |
| Public client | `Client` | Owns the backend and delegates operations |
| Client settings | `Config` | Auth, scope, endpoint, region, timeout |
| Creation intent | `CreateOptions` | Source, runtime, resources, provisioning policy |
| Creation result | `CreateResult` | Shared sandbox metadata returned by a backend |
| User handle | `Sandbox` | Identity and copied `SandboxInfo` |

The provider acts as a factory, but callers select a provider rather than learn
an extra public `Factory` type. Provider packages expose `New() *Provider`; their
initialized `backend` implementation is private. The contracts are defined in the
SDK package that consumes them. There is no dependency-injection container or
runtime SDK reflection. The validated config object remains the construction API.

## Layout

```text
sdks/go/sandbox/           Public SDK module (package sandbox)
sdks/go/providers/modal/   Optional Modal integration (package modal)
sdks/go/providers/daytona/ Optional Daytona integration (package daytona)
proto/kit/sandbox/v1/      Shared SDK contracts
proto/kit/providers/       Provider generation declarations
tooling/internal/go/      Go emitters written in Go
examples/go/              Independent provider examples
```

Go emitter directories retain the explicit `*_gen` convention. `provider_gen`
replaces `adapter_gen`. Generated files consistently end in `.gen.go`:
`client.gen.go`, `client.types.gen.go`, `sandbox.types.gen.go`,
`provider.gen.go`, and `provider.mappings.gen.go`, for example.

`provider.client.gen.go` assembles native SDK configuration; `provider.gen.go`
constructs the SDK. `sandbox.go` implements
sandbox creation. Names describe responsibilities, not an architecture layer.

## Generated language conventions

Protobuf and YAML keep language-neutral snake_case field paths. The Go naming
emitter derives names such as `CPUCores`, `CPULimitCores`, `MemoryMiB`,
`OutboundCIDRs`, `OutboundProxyURL`, `TCPPort`, and `APIKey`.
Enums use typed MixedCaps constants such as `IsolationKindLinuxVM`,
`WaitConditionReady`, and `PolicyModeDisabled`; numeric values are unchanged.
Native SDK field names remain exactly as required by their pinned SDK bindings.

Only the Go emitter owns these spellings. Future emitters use the same semantic
contracts with their own naming conventions. Explicit protobuf JSON names preserve
serialized `context` and `creation` keys for the renamed `Scope` and `Provisioning`
Go fields. Existing scalar presence, units, validation and provider behavior stay
unchanged.

## Migration

These packages are unpublished. In-repository references have been migrated;
source imports and identifiers intentionally change, with no obsolete aliases.

| Previous name | Current name |
| --- | --- |
| `sdks/go/core`, `core` | `sdks/go/sandbox`, `sandbox` |
| `sdks/go/adapters` | `sdks/go/providers` |
| `ClientConfig` | `Config` |
| `ProviderFactory` | `Provider` |
| Runtime `Provider` | `Backend` |
| Provider package `Factory` | `Provider` |
| Provider package `Adapter` | Private `backend` |
| `Initialize` | `NewClient` |
| `ClientContext`, `Config.Context` | `Scope`, `Config.Scope` |
| `Authentication` | `AuthConfig` |
| `APIKeyAuth`, other auth payloads | `APIKeyCredentials`, corresponding credential types |
| `CreateSandboxRequest` / `CreateConfig` | `CreateOptions` |
| `CreateSandboxResponse` | `CreateResult` |
| `CreationOptions`, `CreateConfig.Creation` | `ProvisioningOptions`, `CreateOptions.Provisioning` |

`Client.ProviderName()` and `Sandbox.ProviderName()` remain descriptive public
identity methods. `Close` still releases SDK resources, not sandbox resources.

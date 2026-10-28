# Providers and backends

Sandbox Kit calls existing provider SDKs directly. Applications install the public
core `sandbox` module and only the optional provider modules they need.

## Publishing layout

| Artifact                                                | Role                                                                                                                                                                     |
| ------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Core client (`@sandbox-kit/sandbox`, `sdks/go/sandbox`) | Portable API, shared types, Protovalidate create/config rules, Client/Sandbox choreography. New providers plug into this module; it must not embed a specific cloud SDK. |
| Each provider package (per language)                    | Independently versioned and published. Owns native SDK construction, capability checks, and field mappings for one backend.                                              |

Add providers by shipping a new language-specific module that implements the
core `Provider` / backend contracts—not by growing the core sandbox package.

## Client construction

```go
client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
    Auth: &sandbox.AuthConfig{
        Credentials: &sandbox.TokenPair{
            TokenPair: &sandbox.TokenPairCredentials{Id: tokenID, Secret: tokenSecret},
        },
    },
    Scope: &sandbox.Scope{AppName: sandbox.Value("existing-app")},
})
```

The public client runs shared config validation, then the selected provider maps
configuration, constructs its official SDK, and returns a private backend. The
public client owns that backend and its SDK cleanup.

| Contract                        | Responsibility                                                       |
| ------------------------------- | -------------------------------------------------------------------- |
| `sandbox.Provider`              | Identify the integration and initialize its backend                  |
| `sandbox.Backend`               | Identify the provider and release SDK resources                      |
| `sandbox.SandboxCreator`        | Create sandboxes when the provider supports creation                 |
| Lifecycle capability interfaces | Implement only supported typed operations, such as refresh or delete |
| `sandbox.Client`                | Validate shared requests, apply deadlines, and delegate              |
| `sandbox.Sandbox`               | Hold copied metadata and dispatch supported lifecycle operations     |

These contracts live in the consuming `sandbox` package. SDK clients stay
private; standard operations require no application mapping callbacks.

In TypeScript, `Provider.initialize` and `Backend` methods take the same public
shapes applications use (`Config`, `CreateOptions`, and lifecycle options with
millisecond timeouts). The core `Client` materializes protobuf and runs shared
create/config validation before calling the backend.

## Generation boundary

Protobuf declares shared portable types (including `ErrorInfo` / `ErrorKind`),
service metadata, and Protovalidate create/config rules. Runtime constructors,
client facades, provider orchestration, checks, mappings, and error
classification stay handwritten. Provider YAML inventories describe checks,
mapping groups, configuration assembly, and native bindings for humans to keep
aligned with SDK code.

Provider Go files (`sandbox.go` for SDK calls, `creation.go` for the create
plan, and `lifecycle.go` / `metadata.go` where a provider observes native
sandbox state) and provider-named TypeScript files retain orchestration the
inventories cannot express: app/image/secret lookup, source selection,
readiness, and remaining network or policy semantics. Adding another language
requires handwritten typed integration for those operations.

Optional providers depend on their official SDKs. The core sandbox module
depends on neither provider SDK. Harness integrations, when present, stay
separate, independently published modules.

## Behavior

Creation, verified retrieval, or explicit attachment returns a shared sandbox
handle with sandbox metadata. Errors expose common kinds and preserve SDK or
validator causes. Metadata is copied before exposure.
Closing a client releases SDK resources and leaves cloud sandboxes running.
Handles must not be used after their owning client closes. Closing a sandbox
handle is also local-only and idempotent. Creation is an optional backend
capability, so a provider may expose management or attach without create.

Sandbox Kit introduces no middle service or generated gRPC transport. Native
provider SDKs use their own transports, including gRPC where applicable.

Out of scope for the public API: in-sandbox execution, GPU requests, storage,
and harness integrations.

See [creation coverage](sandbox-creation.md),
[configuration](configuration.md), and [generation](code-generation.md).

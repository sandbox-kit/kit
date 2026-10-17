# Providers and backends

Sandbox Kit calls existing provider SDKs directly. Applications install the public
`sandbox` module and only the optional provider modules they need.

## Client construction

```go
client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
    Auth: &sandbox.AuthConfig{
        Credentials: &sandbox.AuthConfig_TokenPair{
            TokenPair: &sandbox.TokenPairCredentials{Id: tokenID, Secret: tokenSecret},
        },
    },
    Scope: &sandbox.Scope{AppName: sandbox.Value("existing-app")},
})
```

The selected provider validates and maps configuration, constructs its official
SDK, and returns a private backend. The public client owns that backend and its
SDK cleanup.

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

## Generation boundary

Protobuf declares shared types and local methods. `specs/api.yaml` declares
runtime constructors and methods, and `specs/behaviors.yaml` checks those
declarations before emission. Other YAML declares validation, configuration
assembly, field conversions, provider checks, error classification, and native
bindings. Shared error policy lives in `specs/errors.yaml`; provider
classifications remain separate from language-specific SDK bindings.
The shared compiler resolves these inputs into a language-neutral model.
Go and TypeScript emitters consume the compiled mappings, client methods, and
error classifications; `specs/generation.yaml` declares generation phases and
client operation sequences.

Provider `sandbox.go` files and provider-named TypeScript files retain
orchestration that the current vocabulary cannot express: app/image/secret
lookup, source selection, readiness, and remaining network or policy semantics.
Adding another language requires an emitter and typed integration code for those
operations.

Optional providers depend on their official SDKs. The public SDK depends on
neither provider SDK. Future harness integrations will remain separate modules.

## Current behavior

Creation, verified retrieval, or explicit attachment returns a shared sandbox
handle with metadata origins. Errors expose common kinds and preserve SDK or
validator causes. Metadata is copied before exposure.
Closing a client releases SDK resources and leaves cloud sandboxes running.
Handles must not be used after their owning client closes. Closing a sandbox
handle is also local-only and idempotent. Creation is an optional backend capability, allowing future
management-only or attach-only integrations.

Sandbox Kit introduces no middle service or generated gRPC transport. Native
provider SDKs use their own transports, including gRPC where applicable.

Execution, GPU requests, storage, and harness integrations remain future work.
See [creation coverage](sandbox-creation.md),
[configuration](configuration.md), and [generation](code-generation.md).

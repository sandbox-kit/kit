# Providers and backends

Sandbox Kit calls existing provider SDKs directly. Applications install the public
`sandbox` module and only the optional provider modules they need.

## Client construction

```go
client, err := sandbox.NewClient(sandbox.Config{
    Provider: modal.New(),
    Auth: &sandbox.AuthConfig{
        TokenPair: &sandbox.TokenPairCredentials{ID: tokenID, Secret: tokenSecret},
    },
    Scope: &sandbox.Scope{AppName: sandbox.Value("existing-app")},
})
```

The selected provider validates and maps configuration, constructs its official
SDK, and returns a private backend. The public client owns that backend and its
SDK cleanup.

| Contract           | Responsibility                                          |
| ------------------ | ------------------------------------------------------- |
| `sandbox.Provider` | Identify the integration and initialize its backend     |
| `sandbox.Backend`  | Create sandboxes and release SDK resources              |
| `sandbox.Client`   | Validate shared requests, apply deadlines, and delegate |

These contracts live in the consuming `sandbox` package. SDK clients stay
private; standard operations require no application mapping callbacks.

## Generation boundary

Protobuf declares shared types and local methods. YAML declares validation,
configuration assembly, field conversions, provider checks, error classification,
and native bindings. Shared error policy lives in `specs/errors.yaml`; provider
classifications remain separate from language-specific SDK bindings.
The Go emitters generate those mechanical parts.

Provider `sandbox.go` files retain orchestration that the current vocabulary
cannot express: app/image/secret lookup, source selection, readiness, and remaining
network or policy semantics. Adding another language requires an emitter and
typed integration code for those operations.

Optional providers depend on their official SDKs. The public SDK depends on
neither provider SDK. Future harness integrations will remain separate modules.

## Current behavior

Creation returns a shared sandbox handle with metadata origins. Errors expose
common kinds and preserve SDK/validator causes. Metadata is copied before exposure.
Closing a client releases SDK resources and leaves cloud sandboxes running.

Sandbox Kit introduces no middle service or generated gRPC transport. Native
provider SDKs use their own transports, including gRPC where applicable.

Lifecycle methods, execution, GPU requests, storage, and harness integrations
remain future work. See [creation coverage](sandbox-creation.md),
[configuration](configuration.md), and [generation](code-generation.md).

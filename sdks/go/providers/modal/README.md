# Modal provider

Use Modal through the shared Sandbox Kit client. This optional module depends on
the official Modal Go SDK v0.11.0. The public `sandbox` module has no Modal SDK
dependency.

Start with the [creation example](../../../../examples/go/create-modal-sandbox/README.md)
or [resources example](../../../../examples/go/create-modal-with-resources/README.md).
Modules currently use local development replacements.

## Configure the client

Import `github.com/sandbox-kit/kit/sdks/go/providers/modal` and the public
`github.com/sandbox-kit/kit/sdks/go/sandbox` package.

```go
config := sandbox.Config{
    Provider: modal.New(),
    Auth: &sandbox.AuthConfig{
        TokenPair: &sandbox.TokenPairCredentials{
            ID: os.Getenv("MODAL_TOKEN_ID"),
            Secret: os.Getenv("MODAL_TOKEN_SECRET"),
        },
    },
    Scope: &sandbox.Scope{
        AppName: sandbox.Value("existing-app"),
        Environment: sandbox.Value("main"),
    },
    Timeout: sandbox.Value(2 * time.Minute),
}
client, err := sandbox.NewClient(config)
```

Check `err` before using the client. Call `client.Close(ctx)` when finished.
See the runnable examples for imports and cleanup-error handling.

The app must exist in the selected workspace/environment. App names and token
display names identify different objects. Kit does not create apps automatically.

| Client setting    | Support                                                 |
| ----------------- | ------------------------------------------------------- |
| Token pair        | Supported                                               |
| OAuth refresh     | Supported with client ID and exactly one secret/JWT key |
| Omitted auth      | Native SDK environment/profile resolution               |
| Scope             | App name and environment                                |
| Region            | Default sandbox placement preference                    |
| Endpoint override | Rejected through the pinned public constructor          |

See [common configuration](../../../../docs/configuration.md) for the full matrix.

## Create a sandbox

With a caller context `ctx` and an initialized client:

```go
instance, err := client.Create(ctx, &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{
        Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
    },
    Resources: &sandbox.Resources{
        CPUCores: sandbox.Value(0.5),
        MemoryMiB: sandbox.Value(uint64(512)),
    },
    Lifetime: &sandbox.LifetimePolicy{
        MaximumLifetime: sandbox.Value(5 * time.Minute),
    },
})
```

Check `err`, then use `instance.ID()` or `instance.Info()`.

Supplied requests require at least 0.125 physical CPU cores and 128 MiB memory.
CPU precision is 0.001 cores; maximum lifetime is positive whole seconds up to
24 hours. Positive resource limits require a request and cannot be smaller.
Omitted values retain provider defaults.

The integration also maps supported runtime, isolation, placement, security,
network, readiness, and observability fields. Default/snapshot/warm-pool sources
and unsupported fields are rejected. See [creation coverage](../../../../docs/sandbox-creation.md).

## Ownership and verification

Errors expose shared kinds and preserve native SDK causes. Account permissions, resource quotas,
and placement availability remain provider-enforced. Closing the client releases
SDK resources without terminating the sandbox.

From this directory:

```sh
GOWORK=off go test ./...
```

Tests use native SDK doubles and generated mapping checks without cloud
provisioning. [Provider YAML](../../../../specs/providers/modal.yaml) generates
configuration assembly, SDK construction, state capture, validation, and field
conversions. `sandbox.go` contains remaining creation orchestration.

Common error behavior comes from [shared error policy](../../../../specs/errors.yaml).
Provider YAML declares classifications and native bindings; generated errors
preserve the original SDK cause. See [responses and errors](../../../../docs/responses-and-errors.md).

See [generation](../../../../docs/code-generation.md) and
[reference verification](../../../../docs/provider-verification.md).

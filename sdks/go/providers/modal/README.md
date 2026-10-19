# Modal provider

Optional Go module over the official Modal SDK v0.11.0. Use it through the
shared Sandbox Kit client. Requires Go 1.26.9 or newer and this checkout (local
module replacements).

## Usage

```go
import (
    "context"
    "os"
    "time"

    "github.com/sandbox-kit/kit/sdks/go/providers/modal"
    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

config := &sandbox.Config{
    Auth: &sandbox.AuthConfig{
        Credentials: &sandbox.AuthConfig_TokenPair{
            TokenPair: &sandbox.TokenPairCredentials{
                Id:     os.Getenv("MODAL_TOKEN_ID"),
                Secret: os.Getenv("MODAL_TOKEN_SECRET"),
            },
        },
    },
    Scope: &sandbox.Scope{
        AppName:     sandbox.Value("existing-app"),
        Environment: sandbox.Value("main"),
    },
    Timeout: sandbox.Duration(2 * time.Minute),
}
client, err := sandbox.NewClient(modal.New(), config)
if err != nil {
    return err
}
defer client.Close(context.Background())

instance, err := client.Create(ctx, &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{
        Kind: &sandbox.SandboxSource_Image{
            Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
        },
    },
    Resources: &sandbox.Resources{
        CpuCores:  sandbox.Value(0.5),
        MemoryMib: sandbox.Value(uint64(512)),
    },
    Lifetime: &sandbox.LifetimePolicy{
        MaximumLifetime: sandbox.Duration(5 * time.Minute),
    },
})
```

Check `err`, then use `instance.ID()` or `instance.Info()`. Closing the client
releases local SDK resources only. The app must already exist; Kit does not
create apps.

## Examples

| Example    | Project                                                                                            |
| ---------- | -------------------------------------------------------------------------------------------------- |
| Creation   | [create-modal-sandbox](../../../../examples/go/modal/create-modal-sandbox/README.md)               |
| Resources  | [create-modal-with-resources](../../../../examples/go/modal/create-modal-with-resources/README.md) |
| Attach     | [attach-modal-sandbox](../../../../examples/go/modal/attach-modal-sandbox/README.md)               |
| Full index | [Go Modal examples](../../../../examples/go/README.md)                                             |

## Configuration

| Client setting    | Support                                                 |
| ----------------- | ------------------------------------------------------- |
| Token pair        | Supported                                               |
| OAuth refresh     | Supported with client ID and exactly one secret/JWT key |
| Omitted auth      | Native SDK environment/profile resolution               |
| Scope             | App name and environment                                |
| Region            | Default sandbox placement preference                    |
| Endpoint override | Rejected through the pinned public constructor          |

See [common configuration](../../../../docs/configuration.md).

Supplied requests require at least 0.125 physical CPU cores and 128 MiB memory.
CPU precision is 0.001 cores; maximum lifetime is positive whole seconds up to
24 hours. Positive resource limits require a request and cannot be smaller.
Omitted values retain provider defaults.

The integration also maps supported runtime, isolation, placement, security,
network, readiness, and observability fields. Default/snapshot/warm-pool sources
and unsupported fields are rejected. See [creation coverage](../../../../docs/sandbox-creation.md).

## Ownership and verification

Errors expose shared kinds and preserve native SDK causes. Account permissions,
quotas, and placement availability remain provider-enforced. `Client.Attach`
builds a handle from a Modal ID without a remote lookup; handles support
`Terminate`. Provider-verified `Get`/`Refresh`, restartable stop, start, pause,
archive, and resource deletion return `ErrorKindUnsupported`. Successful
termination marks the handle terminal.

```sh
GOWORK=off go test ./...
```

Tests use native SDK doubles without cloud provisioning.
[Provider YAML](../../../../specs/providers/modal.yaml) inventories checks and
mapping groups in `provider_checks.go` and `provider_mappings.go`. `sandbox.go`
owns remaining creation orchestration. Error classification is handwritten in
`errors.go`.

See [responses and errors](../../../../docs/responses-and-errors.md),
[generation](../../../../docs/code-generation.md), and
[reference verification](../../../../docs/provider-verification.md).

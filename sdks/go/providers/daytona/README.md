# Daytona provider

Independent Go module over the official Daytona SDK v0.222.0. Published
separately from the core `sdks/go/sandbox` client; use it through that shared
API. Requires Go 1.26.9 or newer and this checkout (local module replacements).

## Usage

```go
import (
    "context"
    "os"
    "time"

    "github.com/sandbox-kit/kit/sdks/go/providers/daytona"
    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

config := &sandbox.Config{
    Auth: &sandbox.AuthConfig{
        Credentials: &sandbox.APIKey{
            ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
        },
    },
    Timeout: sandbox.Duration(2 * time.Minute),
}
client, err := sandbox.NewClient(daytona.New(), config)
if err != nil {
    return err
}
defer client.Close(context.Background())

// Default snapshot:
instance, err := client.Create(ctx, nil)

// Or an image with resource requests:
instance, err := client.Create(ctx, &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{
        Kind: &sandbox.Image{
            Image: &sandbox.ImageSource{Reference: "python:3.11"},
        },
    },
    Resources: &sandbox.Resources{
        CpuCores:  sandbox.Value(2.0),
        MemoryMib: sandbox.Value(uint64(4096)),
        DiskMib:   sandbox.Value(uint64(8192)),
    },
})
```

Check `err`, then use `instance.ID()` or `instance.Info()`. Closing the client
releases local SDK resources only.

## Examples

| Example    | Project                                                                                                |
| ---------- | ------------------------------------------------------------------------------------------------------ |
| Creation   | [create-daytona-sandbox](../../../../examples/go/daytona/create-daytona-sandbox/README.md)             |
| Policies   | [create-daytona-with-policies](../../../../examples/go/daytona/create-daytona-with-policies/README.md) |
| Retrieval  | [get-daytona-sandbox](../../../../examples/go/daytona/get-daytona-sandbox/README.md)                   |
| Full index | [Go Daytona examples](../../../../examples/go/README.md)                                               |

## Configuration

| Client setting   | Support                           |
| ---------------- | --------------------------------- |
| API key          | Supported                         |
| Bearer/JWT token | Supported; organization required  |
| Omitted auth     | Native SDK environment resolution |
| Scope            | Organization ID                   |
| Endpoint         | `Config.Endpoint` maps to API URL |
| Region           | `Config.Region` maps to target    |

Organization and target may use SDK environment defaults when omitted. An omitted
endpoint defaults to `https://app.daytona.io/api`. Explicit `Config.Endpoint`
takes precedence. Kit does not read `DAYTONA_API_URL` or `DAYTONA_SERVER_URL`.
Remote endpoints require HTTPS; HTTP is accepted for localhost or literal loopback
IPs. See [common configuration](../../../../docs/configuration.md).

CPU requests must be whole cores. Memory/disk requests must be whole GiB
expressed as MiB. Snapshot creation inherits resources, so overrides require an
image. Account quotas remain provider-enforced.

Toolbox languages are `python`, `javascript`, and `typescript`. Per-request
placement is unsupported; use the client's region/target.

## Creation policies

| Request                      | Mapping                                      |
| ---------------------------- | -------------------------------------------- |
| Default/absent policy        | Omitted; provider default applies            |
| Disabled stop/pause          | Native zero                                  |
| Immediate delete             | Native zero                                  |
| Immediate stop/pause/archive | Rejected                                     |
| Disabled archive/delete      | Rejected by the current creation integration |

Archive delays cannot exceed 30 days. Conflicting idle actions and ephemeral
policies fail local validation. See [creation semantics](../../../../docs/sandbox-creation.md).

## Ownership and verification

Errors expose shared kinds and preserve native SDK causes. Handles support
`Start`, `Stop`, `Pause`, `Archive`, `Delete`, and `Refresh`; class-dependent
operations reject known unsupported combinations before the lifecycle call.
Call `instance.Capabilities()` or `instance.Supports(...)` for the resolved
container, VM, Windows, or GPU sandbox. Use `Client.Get` for a provider-verified
handle. Daytona does not support unverified `Client.Attach`.

```sh
GOWORK=off go test ./...
```

Tests use a fake HTTP transport without cloud provisioning.
[Provider YAML](../../../../specs/providers/daytona.yaml) inventories checks and
mapping groups in `provider_checks.go` and `provider_mappings.go`. `sandbox.go`
owns remaining creation orchestration. Error classification is handwritten in
`errors.go`.

See [responses and errors](../../../../docs/responses-and-errors.md),
[generation](../../../../docs/code-generation.md), and
[reference verification](../../../../docs/provider-verification.md).

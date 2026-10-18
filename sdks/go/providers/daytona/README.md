# Daytona provider

Use Daytona through the shared Sandbox Kit client. This optional module depends
on the official Daytona Go SDK v0.222.0. The public `sandbox` module has no
Daytona SDK dependency.

Start with the [creation example](../../../../examples/go/daytona/create-daytona-sandbox/README.md),
[policy example](../../../../examples/go/daytona/create-daytona-with-policies/README.md),
or [retrieval example](../../../../examples/go/daytona/get-daytona-sandbox/README.md).
Requires Go 1.26.9 or newer and this checkout. Modules currently use local
development replacements.

## Configure the client

Import `github.com/sandbox-kit/kit/sdks/go/providers/daytona` and the public
`github.com/sandbox-kit/kit/sdks/go/sandbox` package.

```go
config := sandbox.Config{
    Provider: daytona.New(),
    Auth: &sandbox.AuthConfig{
        APIKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
    },
    Timeout: sandbox.Value(2 * time.Minute),
}
client, err := sandbox.NewClient(config)
```

Check `err` before using the client and call `client.Close(ctx)` when finished.
The runnable examples include imports and cleanup-error handling.

| Client setting   | Support                           |
| ---------------- | --------------------------------- |
| API key          | Supported                         |
| Bearer/JWT token | Supported; organization required  |
| Omitted auth     | Native SDK environment resolution |
| Scope            | Organization ID                   |
| Endpoint         | `Config.Endpoint` maps to API URL |
| Region           | `Config.Region` maps to target    |

Organization and target may use SDK environment defaults when omitted. An omitted
endpoint resolves from `DAYTONA_API_URL`, then `DAYTONA_SERVER_URL`, then
`https://app.daytona.io/api`. Explicit `Config.Endpoint` takes precedence.
Remote endpoints require HTTPS; HTTP is accepted for localhost or literal loopback
IPs. See [common configuration](../../../../docs/configuration.md).

## Create a sandbox

With a caller context `ctx` and an initialized client, use the default snapshot:

```go
instance, err := client.Create(ctx, nil)
```

Or create from an image with resource requests:

```go
instance, err := client.Create(ctx, &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{
        Image: &sandbox.ImageSource{Reference: "python:3.11"},
    },
    Resources: &sandbox.Resources{
        CPUCores: sandbox.Value(2.0),
        MemoryMiB: sandbox.Value(uint64(4096)),
        DiskMiB: sandbox.Value(uint64(8192)),
    },
})
```

Check `err`, then use `instance.ID()` or `instance.Info()`.

CPU requests must be whole cores. Memory/disk requests must be whole GiB
expressed as MiB. Snapshot creation inherits resources, so overrides require an
image. Account quotas remain provider-enforced; standard-tier limits are not
hardcoded into the SDK.

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
policies fail local validation. See [creation semantics](../../../../docs/sandbox-creation.md)
for the full matrix and the SDK/service documentation conflict.

## Ownership and verification

Errors expose shared kinds and preserve native SDK causes. Closing the client
releases SDK resources without stopping or deleting the sandbox. Handles support
`Start`, `Stop`, `Pause`, `Archive`, `Delete`, and `Refresh`; class-dependent
operations reject known unsupported combinations before the lifecycle call.
Call `instance.Capabilities()` or `instance.Supports(...)` to inspect support
for the resolved container, VM, Windows, or GPU sandbox. Lifecycle responses
update the handle's normalized state and provider metadata.
Use `Client.Get` to retrieve a provider-verified handle for an existing Daytona
sandbox. Daytona does not support unverified `Client.Attach`.

From this directory:

```sh
GOWORK=off go test ./...
```

Tests use a fake HTTP transport without cloud provisioning.
[Provider YAML](../../../../specs/providers/daytona.yaml) inventories checks and
mapping groups implemented in `provider_checks.go` and `provider_mappings.go`.
`sandbox.go` owns remaining creation orchestration.

Error classification is handwritten in `errors.go` and preserves the original
SDK cause. Shared payloads use protobuf `ErrorInfo` / `ErrorKind`. See
[responses and errors](../../../../docs/responses-and-errors.md).

See [generation](../../../../docs/code-generation.md) and
[reference verification](../../../../docs/provider-verification.md).

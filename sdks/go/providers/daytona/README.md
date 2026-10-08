# Daytona provider

Use Daytona through the common Sandbox Kit client. This optional Go module
depends on the official Daytona Go SDK v0.222.0; the public `sandbox` module does
not. Packages are unpublished; start with [the checkout example](../../../../examples/go/create-daytona-sandbox/README.md).

## Configure the client

```go
import (
    "context"
    "os"

    "github.com/sandbox-kit/kit/sdks/go/providers/daytona"
    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)
```

Inside your application:

```go
client, err := sandbox.NewClient(sandbox.Config{
    Provider: daytona.New(),
    Auth: &sandbox.AuthConfig{
        APIKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
    },
})
if err != nil {
    return err
}
defer client.Close(context.Background())
```

Supply a context to creation. Omitted `Auth` uses
SDK environment resolution. Bearer/JWT auth is also supported; it requires an
organization in `Config.Scope.OrganizationID` or SDK environment settings.
Endpoint and target can be set through `Config.Endpoint` and `Config.Region`, or
left to the SDK's defaults. See [configuration](../../../../docs/configuration.md).

## Create a sandbox

Use the default snapshot:

```go
instance, err := client.Create(ctx, nil)
```

Or request an image and resources:

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
if err != nil {
    return err
}
id := instance.ID()
```

Check `err` for default creation too. Shared units are MiB; requests must represent
whole CPU cores and whole GiB allocations. Snapshot sources inherit resources,
so resource overrides require an image. Per-request placement is rejected;
configure the client's region/target instead. See [creation support](../../../../docs/sandbox-creation.md).

`client.Close(ctx)` releases SDK resources, not the sandbox. Execution and sandbox
lifecycle methods remain future work. Provider SDK errors pass through unchanged.

## Development

Run `GOWORK=off go test ./...` from this module. Local replacements resolve the
public SDK. [Provider YAML](../../../../specs/providers/daytona.yaml) generates
configuration, credential, resource and metadata mappings; typed SDK construction
and semantic mappings live in `client.go` and `sandbox.go`.

# Modal provider

Use Modal through the common Sandbox Kit client. This optional Go module depends
on the official Modal Go SDK v0.11.0; the public `sandbox` module does not.
Packages are unpublished; start with [the checkout example](../../../../examples/go/create-modal-sandbox/README.md).

## Configure the client

```go
import (
    "context"
    "os"

    "github.com/sandbox-kit/kit/sdks/go/providers/modal"
    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)
```

Inside your application:

```go
client, err := sandbox.NewClient(sandbox.Config{
    Provider: modal.New(),
    Auth: &sandbox.AuthConfig{
        TokenPair: &sandbox.TokenPairCredentials{
            ID: os.Getenv("MODAL_TOKEN_ID"),
            Secret: os.Getenv("MODAL_TOKEN_SECRET"),
        },
    },
    Scope: &sandbox.Scope{
        AppName: sandbox.Value("my-existing-app"),
        Environment: sandbox.Value("main"),
    },
})
if err != nil {
    return err
}
defer client.Close(context.Background())
```

Supply a context to creation. App name
and token display name are separate. The app must exist in the selected workspace
and environment. [The example guide](../../../../examples/go/create-modal-sandbox/README.md)
explains credentials and one-time app creation.

Explicit token-pair and OAuth-refresh auth are supported. Omit `Auth` to use
Modal's SDK environment/profile resolution. App scope is still needed for creation.
`Config.Endpoint` is rejected because the pinned SDK has no public per-client
override. `Config.Region` becomes a default sandbox placement preference.
See [the complete configuration matrix](../../../../docs/configuration.md).

## Create a sandbox

```go
instance, err := client.Create(ctx, &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{
        Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
    },
})
if err != nil {
    return err
}
id := instance.ID()
```

The current mapping requires an explicit registry image. Default/snapshot/warm-pool
sources are not implemented. CPU in 0.001-core increments, supported memory limits, runtime
command/working directory, and other mapped creation settings are available;
unsupported intent is rejected. See [creation support](../../../../docs/sandbox-creation.md).

Requests require at least 0.125 CPU cores and 128 MiB memory when supplied.
Maximum lifetime is capped at 24 hours. Omitted resource values retain Modal's
defaults; account quotas and placement availability are checked by Modal.

`client.Close(ctx)` releases the SDK, not the sandbox. Sandbox lifecycle methods
remain future work. Provider SDK errors pass through unchanged.

## Development

From this module, run `GOWORK=off go test ./...`. Local replacements resolve the
public SDK from the checkout. Generated provider/client-field mappings come from
[provider YAML](../../../../specs/providers/modal.yaml); typed SDK construction
and scope capture are generated. Configuration assembly is in `provider.client.gen.go`;
creation orchestration and semantic mappings remain in `sandbox.go`.

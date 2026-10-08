# Sandbox Kit

Create sandboxes through one Go API, with separately installable provider
integrations. Sandbox Kit currently supports **Modal** and **Daytona**. Providers
initialize their official SDK internally from common client configuration; your
application imports Sandbox Kit and the provider it needs.

**Current scope:** authentication, sandbox creation, and shared identity/metadata.
Execution, sandbox stop/delete, snapshots, and harness integrations are not yet
implemented. Other languages are planned; Go is the implemented SDK today.

## Start with a working example

Requires **Go 1.26.1 or newer**. Modules are not published yet; use the development
checkout. The current implementation is on the feature branch:

```sh
git clone --branch feat/unified-sandbox-sdk https://github.com/sandbox-kit/kit.git
cd kit
```

Choose one project:

| Provider | Project | Required `.env` values |
| --- | --- | --- |
| Daytona | [create-daytona-sandbox](examples/go/create-daytona-sandbox/README.md) | `DAYTONA_API_KEY` |
| Modal | [create-modal-sandbox](examples/go/create-modal-sandbox/README.md) | `MODAL_TOKEN_ID`, `MODAL_TOKEN_SECRET`, `MODAL_APP_NAME` |

For Daytona:

```sh
cd examples/go/create-daytona-sandbox
cp .env.example .env   # Skip if you already have .env.
# Edit .env and set DAYTONA_API_KEY from your dashboard.
go run .
```

For Modal, run the same commands from `examples/go/create-modal-sandbox` and fill
in its token pair and an existing app name. `MODAL_ENVIRONMENT` defaults to `main`.
See [Modal setup, including app creation](examples/go/create-modal-sandbox/README.md).

Each project has its own Go module and imports only its selected provider.
There are no flags or provider selectors. Running it creates **one real sandbox**
and prints its ID. `.env` is Git-ignored; commit only the blank `.env.example`.

## Use the SDK in Go

The main steps are select a provider, construct a client, and create a sandbox:

```go
import (
    "context"
    "fmt"
    "os"

    "github.com/sandbox-kit/kit/sdks/go/providers/daytona"
    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func createSandbox(ctx context.Context) error {
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

    instance, err := client.Create(ctx, nil) // Daytona's default snapshot.
    if err != nil {
        return err
    }
    fmt.Println(instance.ID(), instance.ProviderName())
    return nil
}
```

Applications must supply environment variables themselves. Only the example
projects load `.env` using `godotenv`; the SDK does not load files automatically.
The runnable examples also propagate SDK cleanup errors.

For Modal, select `modal.New()` and configure its supported credentials/scope:

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
        AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME")),
        Environment: sandbox.Value("main"),
    },
}
```

Import `github.com/sandbox-kit/kit/sdks/go/providers/modal`. Pass `config` to
`sandbox.NewClient`, then create with an explicit image:

```go
instance, err := client.Create(ctx, &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{
        Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
    },
})
```

Modal requires an existing app and image for the current mapping. Daytona can
use its default snapshot with `nil` creation options. Both return the same
`sandbox.Sandbox` handle. See [Modal](sdks/go/providers/modal/README.md) and
[Daytona](sdks/go/providers/daytona/README.md) for provider-specific usage.

## Configure creation

Use `sandbox.CreateOptions` for image/snapshot source, runtime, isolation,
resources, environment variables, labels, placement, and supported policies.
Fields express intent; providers reject settings they cannot honor.

```go
options := &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{
        Image: &sandbox.ImageSource{Reference: "python:3.11"},
    },
    Resources: &sandbox.Resources{
        CPUCores: sandbox.Value(2.0),
        MemoryMiB: sandbox.Value(uint64(4096)),
    },
    Environment: map[string]string{"MODE": "development"},
    Labels: map[string]string{"team": "sandbox-kit"},
}
```

`sandbox.Value` marks an optional value as supplied, including zero or false.
Memory/disk units are MiB. Daytona resource overrides require an image source,
whole CPU cores, and whole GiB expressed as MiB. Modal also supports fractional
CPU. See [creation support and semantics](docs/sandbox-creation.md).

Client settings live in `sandbox.Config`: `Provider`, `Auth`, `Scope`, `Endpoint`,
`Region`, and `Timeout`. Supported auth modes and context differ by provider.
For example, Modal's pinned Go SDK does not expose a per-client endpoint override.
See [the configuration matrix](docs/configuration.md).

`Config.Timeout` is a default operation deadline. Per-operation
`CreateOptions.Provisioning.Timeout` overrides it; zero adds no Kit deadline.
Caller context deadlines still apply. Provider SDK errors pass through unchanged.

## Results and cleanup

* `instance.ID()` — sandbox ID.
* `instance.ProviderName()` — selected provider.
* `instance.Info()` — copied shared metadata, including available provider state/resources.
* `client.Close(ctx)` — release SDK resources.

**Closing the client does not stop or delete a sandbox.** Use the provider's
management tools/dashboard for sandbox cleanup until lifecycle methods are added.
No additional gRPC service or middle service is introduced by Sandbox Kit.

## Development

The repository contains separate SDK/provider/example modules. `go.work` and local
`replace` directives wire the checkout together; v0.0.0 requirements are development
placeholders, not published releases. For your own application before publication,
use local replacements for the SDK and selected provider, as the example modules do.

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit generate go  # Requires protoc.
```

Tests compile both examples without executing cloud creation. Protobuf and YAML
are generation inputs; users receive native Go types without protobuf imports.

```text
sdks/go/sandbox/       Public Go SDK
sdks/go/providers/     Optional provider modules
examples/go/           Separate Modal and Daytona projects
proto/                 Shared types and generation declarations
specs/                 Validation and provider mapping rules
tooling/               Go generators and Cobra development commands
docs/                  Configuration, semantics, and design
```

Further reading: [code generation](docs/code-generation.md),
[provider architecture](docs/providers.md), [naming](docs/naming.md).

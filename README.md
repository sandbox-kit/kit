# Sandbox Kit

Create sandboxes through one Go API. Install the public SDK and the provider you
need; each provider initializes its official SDK from shared configuration.

Go integrations for Modal and Daytona currently support authentication, sandbox
creation, validation, shared metadata, and common error kinds. Execution, lifecycle methods,
storage, GPU requests, harness integrations, and other language SDKs are planned.

## Get started

Requires Go 1.26.1 or newer. Modules use development versions and local
replacements, so run the examples from this repository checkout:

```sh
git clone https://github.com/sandbox-kit/kit.git
cd kit/examples/go/create-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY in .env.
go run .
```

Skip copying the template if you already have a configured `.env`.
Running a creation example creates one real sandbox. Error-handling examples
exercise local failures without cloud creation. Closing the client releases SDK
resources; use the provider's tools to stop or delete the sandbox.

| Example                                                                | Demonstrates                                                   |
| ---------------------------------------------------------------------- | -------------------------------------------------------------- |
| [Daytona creation](examples/go/create-daytona-sandbox/README.md)       | API-key authentication and the default snapshot                |
| [Modal creation](examples/go/create-modal-sandbox/README.md)           | Token-pair authentication, app/environment scope, and an image |
| [Daytona policies](examples/go/create-daytona-with-policies/README.md) | Disabled auto-pause and delayed deletion                       |
| [Modal resources](examples/go/create-modal-with-resources/README.md)   | CPU, memory, and a bounded lifetime                            |
| [Modal errors](examples/go/handle-modal-errors/README.md) | Handle an unsupported disk override with `ErrorKindUnsupported` |
| [Daytona errors](examples/go/handle-daytona-errors/README.md) | Handle an invalid CPU request with `ErrorKindInvalidArgument` |

Each project has its own module, `.env.example`, and selected provider. Modal
requires an existing app; its [setup guide](examples/go/create-modal-sandbox/README.md)
explains token and app creation. Credentials stay in Git-ignored `.env` files.

## Use the Go SDK

```go
import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/sandbox-kit/kit/sdks/go/providers/daytona"
    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func createSandbox() (result error) {
    client, err := sandbox.NewClient(sandbox.Config{
        Provider: daytona.New(),
        Auth: &sandbox.AuthConfig{
            APIKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
        },
        Timeout: sandbox.Value(2 * time.Minute),
    })
    if err != nil {
        return err
    }
    defer func() {
        if closeErr := client.Close(context.Background()); closeErr != nil {
            if result == nil {
                result = closeErr
            } else {
                handleError(closeErr)
            }
        }
    }()

    instance, err := client.Create(context.Background(), nil)
    if err != nil {
        return err
    }
    fmt.Println(instance.ID(), instance.ProviderName())
    return nil
}
```

Applications supply credentials and environment variables. The examples load
`.env` with `godotenv`; the SDK itself does not read files.

For your own application, add local replacements for the public SDK and selected
provider as shown in the example modules. The `v0.0.0` requirements are
development placeholders, not published release versions.

## Handle errors

Use `errors.As` to extract `*sandbox.Error`, then switch on its `Kind`.
The same categories work with Modal and Daytona. No native SDK imports are needed.
Call this handler when `createSandbox()` returns an error; the cleanup code above
also uses it when creation and cleanup both fail.

```go
import (
    "errors"
    "fmt"
    "os"

    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func handleError(err error) {
    var sandboxError *sandbox.Error
    if !errors.As(err, &sandboxError) {
        fmt.Fprintln(os.Stderr, err) // For example, a missing .env file.
        return
    }

    switch sandboxError.Kind {
    case sandbox.ErrorKindInvalidArgument, sandbox.ErrorKindUnsupported:
        fmt.Fprintf(os.Stderr, "Check %s: %s\n", sandboxError.Field, sandboxError.Message)
    case sandbox.ErrorKindAuthentication:
        fmt.Fprintln(os.Stderr, "Check your provider credentials.")
    case sandbox.ErrorKindPermissionDenied:
        fmt.Fprintln(os.Stderr, "Your account does not have permission for this operation.")
    case sandbox.ErrorKindRateLimited, sandbox.ErrorKindResourceExhausted:
        fmt.Fprintln(os.Stderr, "The provider's rate or resource limit was reached.")
    case sandbox.ErrorKindCanceled:
        fmt.Fprintln(os.Stderr, "The operation was canceled.")
    case sandbox.ErrorKindTimeout:
        fmt.Fprintln(os.Stderr, "The operation timed out.")
    default:
        fmt.Fprintln(os.Stderr, err)
    }
}
```

`Field` identifies the setting involved when available. Error kinds do not imply
safe retries: creation can fail after a sandbox has already been provisioned.
See the independent [Modal](examples/go/handle-modal-errors/README.md) and
[Daytona](examples/go/handle-daytona-errors/README.md) error examples, or the
[full contract guide](docs/responses-and-errors.md).

## Configuration and results

* `sandbox.Config` selects the provider and configures auth, scope, endpoint,
  region, and the default operation timeout.
* `sandbox.CreateOptions` describes the source, resources, runtime, network,
  and supported creation policies.
* `sandbox.Value(v)` returns a pointer that marks an optional value as supplied.
  It performs no validation itself.
* `Sandbox.ID()`, `ProviderName()`, and `Info()` expose shared identity and copied
  metadata. Metadata availability varies by provider.

Omitted values defer to provider defaults. Explicit values are validated for
supported syntax, units, precision, and combinations. Unsupported settings fail
before provider calls; account quotas and current capacity remain provider-enforced.
Errors expose shared categories and field paths while preserving native causes.
Response origins distinguish provider-reported metadata from requested settings.
See [responses and errors](docs/responses-and-errors.md).

Memory and disk use MiB. Modal accepts fractional physical CPU requests in
0.001-core increments from 0.125 cores; Daytona requires whole CPU cores and
whole GiB allocations expressed as MiB.

## Documentation

| Guide                                                   | Covers                                                    |
| ------------------------------------------------------- | --------------------------------------------------------- |
| [Client configuration](docs/configuration.md)           | Auth, scope, defaults, deadlines, and ownership           |
| [Sandbox creation](docs/sandbox-creation.md)            | Supported fields, units, and policy semantics             |
| [Provider architecture](docs/providers.md)              | Optional modules, contracts, and SDK boundaries           |
| [Responses and errors](docs/responses-and-errors.md) | Presence, origins, error details, and native causes |
| [API declarations](docs/api-declarations.md) | Runtime declarations and language representation rules |
| [Code generation](docs/code-generation.md)              | Protobuf, YAML rules, bindings, and extension workflow    |
| [Reference verification](docs/provider-verification.md) | Sources, conflicting documentation, and validation limits |
| [Naming](docs/naming.md)                                | Public names, package layout, and generated files         |

## Develop and contribute

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit generate go
```

Generation also requires `protoc`. Verification tests the SDK, tooling, providers,
and all six example modules. Example tests use dummy credentials and local
servers; they create no cloud resources.

Put contracts in `proto/`, portable rules in `specs/`, and all generators in the
Go module under `tooling/`. A shared compiler in `tooling/internal/model/` resolves
these inputs before language emission. `specs/generation.yaml` defines generation
phases and client-operation sequences. Runtime modules live in `sdks/<language>/`.
Update examples and docs with public changes, regenerate affected output, and
verify before submitting. See [tooling](tooling/README.md) and
[repository instructions](AGENTS.md).

## License

Sandbox Kit uses the [MIT license](LICENSE).

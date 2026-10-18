# Sandbox Kit

Sandbox Kit is a lightweight SDK for creating and managing sandboxes across providers.
Install the public SDK and the provider you need; each provider initializes its
official SDK from shared configuration.

Go and TypeScript support shared contracts, the public client, Daytona, and
Modal. Both languages include authentication, sandbox creation, supported
lifecycle operations, validation, copied metadata, common error kinds, and a
complete matrix of independent provider examples. Execution, storage, GPU
requests, and harness integrations are planned.

## Get started

Requires Go 1.26.9 or newer. Modules use development versions and local
replacements, so run the examples from this repository checkout:

```sh
git clone https://github.com/sandbox-kit/kit.git
cd kit/examples/go/daytona/create-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY in .env.
go run .
```

For a TypeScript example, use Node.js 24.15.0, npm 11.19.0, and one of the
[provider projects](examples/typescript/README.md). For Daytona:

```sh
cd examples/typescript/daytona/create-daytona-sandbox
cp .env.example .env
npm install
npm run build
npm start
```

Skip copying the template if you already have a configured `.env`.
Creation examples create real sandboxes: one each, except the linked Daytona
example, which creates a parent and a child. Lifecycle examples operate on the
ID supplied in `.env`; destructive examples stop, terminate, archive, or delete
that sandbox as their names indicate. Error-handling examples exercise local
failures without cloud creation. Closing the client or a handle releases only
local resources and does not change a cloud sandbox.

| Example                                                                            | Demonstrates                                                    |
| ---------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| [Daytona creation](examples/go/daytona/create-daytona-sandbox/README.md)           | API-key authentication and the default snapshot                 |
| [Daytona snapshot](examples/go/daytona/create-daytona-from-snapshot/README.md)     | Stock `daytona-small` snapshot                                  |
| [Daytona image](examples/go/daytona/create-daytona-from-image/README.md)           | Public image with CPU, memory, and disk                         |
| [Daytona runtime](examples/go/daytona/create-daytona-with-runtime/README.md)       | Toolbox language, version, and user                             |
| [Daytona ephemeral](examples/go/daytona/create-daytona-ephemeral/README.md)        | Delete-on-stop from the default snapshot                        |
| [Daytona linked](examples/go/daytona/create-daytona-linked/README.md)              | Parent sandbox and an ephemeral linked child                    |
| [Daytona policies](examples/go/daytona/create-daytona-with-policies/README.md)     | Disabled auto-pause and delayed deletion                        |
| [Daytona get](examples/go/daytona/get-daytona-sandbox/README.md)                   | Retrieve an existing sandbox and its resolved capabilities      |
| [Daytona refresh](examples/go/daytona/refresh-daytona-sandbox/README.md)           | Refresh metadata on an existing handle                          |
| [Daytona start](examples/go/daytona/start-daytona-sandbox/README.md)               | Start a stopped sandbox                                         |
| [Daytona stop](examples/go/daytona/stop-daytona-sandbox/README.md)                 | Stop a restartable sandbox                                      |
| [Daytona pause](examples/go/daytona/pause-daytona-sandbox/README.md)               | Pause a supported VM sandbox                                    |
| [Daytona archive](examples/go/daytona/archive-daytona-sandbox/README.md)           | Archive a supported container sandbox                           |
| [Daytona delete](examples/go/daytona/delete-daytona-sandbox/README.md)             | Delete an existing sandbox                                      |
| [Daytona close handle](examples/go/daytona/close-daytona-handle/README.md)         | Release a local handle without changing the cloud sandbox       |
| [Modal creation](examples/go/modal/create-modal-sandbox/README.md)                 | Token-pair authentication, app/environment scope, and an image  |
| [Modal resources](examples/go/modal/create-modal-with-resources/README.md)         | CPU, memory, and a bounded lifetime                             |
| [Modal runtime](examples/go/modal/create-modal-with-runtime/README.md)             | Working directory, entrypoint, and PTY                          |
| [Modal gVisor](examples/go/modal/create-modal-gvisor/README.md)                    | Container isolation                                             |
| [Modal Linux VM](examples/go/modal/create-modal-linux-vm/README.md)                | Linux VM isolation                                              |
| [Modal placement](examples/go/modal/create-modal-with-placement/README.md)         | Cloud and optional region preference                            |
| [Modal labels](examples/go/modal/create-modal-with-labels/README.md)               | Environment variables and tags                                  |
| [Modal network](examples/go/modal/create-modal-with-network/README.md)             | Blocked outbound access                                         |
| [Modal idle timeout](examples/go/modal/create-modal-idle/README.md)                | Termination after 60 idle seconds                               |
| [Modal readiness](examples/go/modal/create-modal-ready/README.md)                  | Wait until a readiness command succeeds                         |
| [Modal verbose](examples/go/modal/create-modal-verbose/README.md)                  | Verbose provider logging                                        |
| [Modal identity](examples/go/modal/create-modal-with-identity/README.md)           | Workload identity token                                         |
| [Modal attach](examples/go/modal/attach-modal-sandbox/README.md)                   | Construct a handle from an ID without a remote lookup           |
| [Modal terminate](examples/go/modal/terminate-modal-sandbox/README.md)             | Terminate an attached sandbox                                   |
| [Modal termination request](examples/go/modal/request-modal-termination/README.md) | Request termination without waiting for an exit result          |
| [Modal close handle](examples/go/modal/close-modal-handle/README.md)               | Release a local handle without terminating the sandbox          |
| [Modal errors](examples/go/modal/handle-modal-errors/README.md)                    | Handle an unsupported disk override with `ErrorKindUnsupported` |
| [Daytona errors](examples/go/daytona/handle-daytona-errors/README.md)              | Handle an invalid CPU request with `ErrorKindInvalidArgument`   |

Each project has its own module, `.env.example`, and selected provider. Modal
requires an existing app; its [setup guide](examples/go/modal/create-modal-sandbox/README.md)
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
safe retries: creation can fail after a sandbox has been provisioned, and a
lifecycle request can fail after the provider has begun or completed a transition.
See the independent [Modal](examples/go/modal/handle-modal-errors/README.md) and
[Daytona](examples/go/daytona/handle-daytona-errors/README.md) error examples, or the
[full contract guide](docs/responses-and-errors.md).

## Configuration and results

- `sandbox.Config` selects the provider and configures auth, scope, endpoint,
  region, and the default operation timeout.
- `sandbox.CreateOptions` describes the source, resources, runtime, network,
  and supported creation policies.
- `sandbox.Value(v)` returns a pointer that marks an optional value as supplied.
  It performs no validation itself.
- `Sandbox.ID()`, `ProviderName()`, and `Info()` expose shared identity and copied
  metadata. Metadata availability varies by provider.
- `Client.Get` retrieves current provider metadata when supported;
  `Client.Attach` constructs an explicitly unverified handle when supported.
- `Sandbox.Capabilities` and `Supports` report lifecycle operations available
  for the provider and resolved sandbox class.

Omitted values defer to provider defaults. Explicit values are validated for
supported syntax, units, precision, and combinations. Unsupported settings fail
before provider calls; account quotas and current capacity remain provider-enforced.
Errors expose shared categories and field paths while preserving native causes.
Response origins distinguish provider-reported metadata from requested settings.
See [responses and errors](docs/responses-and-errors.md).

Creation is an optional provider capability. `Client.Close` and
`Sandbox.Close` release local resources only; neither changes a remote sandbox.
Do not use handles after closing their owning client; the provider SDK resources
they depend on have been released.

Memory and disk use MiB. Modal accepts fractional physical CPU requests in
0.001-core increments from 0.125 cores; Daytona requires whole CPU cores and
whole GiB allocations expressed as MiB.

## Documentation

| Guide                                                   | Covers                                                    |
| ------------------------------------------------------- | --------------------------------------------------------- |
| [Client configuration](docs/configuration.md)           | Auth, scope, defaults, deadlines, and ownership           |
| [Sandbox creation](docs/sandbox-creation.md)            | Supported fields, units, and policy semantics             |
| [Sandbox lifecycle](docs/sandbox-lifecycle.md)          | Retrieval, lifecycle semantics, and provider capabilities |
| [Provider architecture](docs/providers.md)              | Optional modules, contracts, and SDK boundaries           |
| [Responses and errors](docs/responses-and-errors.md)    | Presence, origins, error details, and native causes       |
| [API declarations](docs/api-declarations.md)            | Runtime declarations and language representation rules    |
| [Code generation](docs/code-generation.md)              | Protobuf, YAML rules, bindings, and extension workflow    |
| [Reference verification](docs/provider-verification.md) | Sources, conflicting documentation, and validation limits |
| [Naming](docs/naming.md)                                | Public names, package layout, and generated files         |

## Develop and contribute

Requires Go 1.26.9 or newer. From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit generate contracts-go
GOWORK=off go run ./cmd/sandbox-kit test typescript
GOWORK=off go run ./cmd/sandbox-kit generate contracts-typescript
```

Contract generation requires `protoc`. Verification tests the SDK, tooling,
providers, and every example module. Example tests use dummy credentials and
local servers; they create no cloud resources. Each module is checked
independently with `GOWORK=off`. See [tooling setup](tooling/README.md) for
workspace synchronization and dependency checksum updates.

Put contracts in `proto/`, portable rule inventories in `specs/`, and the
developer CLI in `tooling/`. Portable public types come from protobuf contracts
with handwritten SDK aliases. Runtime modules live in `sdks/<language>/`.
Update examples and docs with public changes, regenerate contracts when `proto/`
changes, and verify before submitting. See [tooling](tooling/README.md) and
[repository instructions](AGENTS.md).

## License

Sandbox Kit uses the [MIT license](LICENSE).

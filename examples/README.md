# Sandbox Kit: client initialization

One runnable example shows the complete initialization flow for Modal and
Daytona:

1. Your application initializes the official SDK with its authentication configuration.
2. Pass that client to the selected adapter's `New` constructor.
3. Pass the adapter to `core.NewClient` and use the common client API.

Read [main.go](main.go). The provider-specific functions keep the three steps
explicit and handle errors. Modal's SDK client is closed by the application,
after the core client is finished.

## Run

From the repository root:

```sh
cd examples
go run . --help
```

Requires Go 1.26.1 or newer. Dependencies may be downloaded on the first run.
Help does not initialize a provider or require credentials.

For Modal, configure its SDK profile or set `MODAL_TOKEN_ID` and
`MODAL_TOKEN_SECRET`, then run:

```sh
go run . modal
```

For Daytona, set `DAYTONA_API_KEY`, or set `DAYTONA_JWT_TOKEN` together with
`DAYTONA_ORGANIZATION_ID`, then run:

```sh
go run . daytona
```

Expected output for the selected provider:

```text
Sandbox Kit ready: modal
```

or `Sandbox Kit ready: daytona`.

## The important lines

Once your SDK client is initialized:

These are the application imports for the Modal path; the complete program also
imports Modal's SDK to initialize the client:

```go
import (
    "fmt"

    modalAdapter "github.com/sandbox-kit/kit/adapters/modal"
    "github.com/sandbox-kit/kit/core"
)
```

Inside your application function:

```go
provider, err := modalAdapter.New(sdkClient)
if err != nil {
    return err
}
client, err := core.NewClient(provider)
if err != nil {
    return err
}
fmt.Println(client.ProviderName())
```

The Daytona path imports `github.com/sandbox-kit/kit/adapters/daytona` and uses
`daytonaAdapter.New(sdkClient)`. Both produce the same
`*core.Client`. Core and adapters do not import official SDKs; only the application
does. This example imports both providers to demonstrate both paths in one
program. Your application can depend on only its chosen provider.

This increment covers initialization and provider identity. It does not create
sandboxes or demonstrate shared operation responses yet. The adapter borrows the
SDK pointer without reauthenticating or changing its configuration.

## Verify

```sh
GOWORK=off go test ./...
```

The tests cover help, invalid selection, and injection of both SDK client types
without authenticating or contacting a provider.
Local `replace` directives resolve core and adapters from this checkout, so the
example also builds with workspace mode disabled. Credentialed runs use the
provider SDK's normal initialization behavior.

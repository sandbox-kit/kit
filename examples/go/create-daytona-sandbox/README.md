# Create a Daytona sandbox

This standalone Go project authenticates with Daytona through Sandbox Kit,
creates one sandbox from the default snapshot, and prints its ID.
It imports only the public SDK and the Daytona provider. Requires Go 1.26.1+.

## 1. Configure credentials

From the repository root:

```sh
cd examples/go/create-daytona-sandbox
cp .env.example .env   # Skip if .env already exists.
```

Edit `.env` with an API key from your Daytona dashboard:

```dotenv
DAYTONA_API_KEY=your_api_key
```

No provider selector, app name, or command-line flags are required.
See [.env.example](.env.example) for the blank template.

## 2. Run

```sh
go run .
```

Expected output:

```text
Daytona sandbox created: <sandbox-id>
```

The program loads `.env`, passes the API key through `sandbox.Config.Auth`, and
calls `client.Create(ctx, nil)` for the provider's default snapshot, following
[Daytona's Go quickstart](https://www.daytona.io/docs/en/go-sdk/).
See [main.go](main.go) for the complete error-handled code.
Existing shell environment variables take precedence over `.env`; the SDK itself
does not load `.env` files.

The example uses Daytona's normal endpoint/target defaults. If necessary, add
`DAYTONA_API_URL` and `DAYTONA_TARGET` to `.env`; the underlying SDK reads those
settings when the corresponding `sandbox.Config` values are omitted.
For explicit configuration, see [the provider guide](../../../sdks/go/providers/daytona/README.md).

## Scope and cleanup

The example leaves lifetime policies at Daytona defaults. To explicitly disable
auto-pause, replace its `client.Create(context.Background(), nil)` call with:

```go
instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
    Lifetime: &sandbox.LifetimePolicy{
        IdlePause: &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled},
    },
})
```

Auto-stop and auto-pause support disabling with native zero. Explicit archive/delete
disabling is rejected; immediate deletion uses zero. See [policy semantics](../../../docs/sandbox-creation.md).

This example only authenticates and creates a sandbox. It does not execute
commands or run sandbox lifecycle operations. Running it provisions a real sandbox.
Deferred `client.Close` releases SDK resources but does not delete the sandbox;
manage cleanup through Daytona's dashboard/tools.

`.env` is Git-ignored; only the blank `.env.example` belongs in Git.
Compile without provisioning:

```sh
GOWORK=off go test ./...
```

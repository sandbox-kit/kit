# Close a Daytona sandbox handle

Retrieve one Daytona sandbox, then close only its local Sandbox Kit handle.
`Sandbox.Close` releases the handle and rejects later lifecycle calls through
that handle. It does not stop, archive, or delete the cloud sandbox.

The retrieved metadata still reports whether the sandbox is a container,
Linux VM, Windows sandbox, or GPU-backed container. Handle cleanup behaves the
same for every Daytona class because it is local rather than a provider
lifecycle operation.

## Run

Requires Go 1.26.9 or newer and this checkout.

```sh
cd examples/go/daytona/close-daytona-handle
cp .env.example .env
# Set DAYTONA_API_KEY and DAYTONA_SANDBOX_ID, then:
go run .
```

`sandbox.NewClient` initializes the Daytona provider. The deferred
`Client.Close` releases the private Daytona SDK client after the handle closes;
neither close method changes the remote sandbox.

## Verify locally

```sh
GOWORK=off go test ./...
```

This compiles the example without contacting Daytona. Running it contacts
Daytona to retrieve the configured sandbox.

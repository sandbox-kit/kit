# Get a Daytona sandbox

Load a provider-verified handle for `DAYTONA_SANDBOX_ID` with
`Client.Get` / `client.get` and print identity, metadata, and capabilities.

## What this means

`Get` performs a remote lookup and returns current Daytona metadata (state,
class, labels, …). Use it when you already have an ID and need a fresh handle.
Capabilities list the lifecycle methods Daytona implements.

## Edge cases

- Unknown or deleted IDs fail at the provider (`NotFound` or similar).
- The resolved class can differ from what you assumed at create time.
- Does not change the cloud sandbox.

## Run

```sh
cd examples/go/daytona/get-daytona-sandbox
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
DAYTONA_API_KEY=
DAYTONA_SANDBOX_ID=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Needs an **existing sandbox ID** in `.env`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/daytona/get-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

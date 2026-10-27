# Close a Daytona sandbox handle

`Get` a sandbox, then call **`Sandbox.Close` / `sandbox.close()`**.

## What this means

Daytona has no sandbox close, so handle close returns without a provider
request. It does not stop, archive, or delete the sandbox. Client close
releases the SDK connection.

## Edge cases

- Safe against a sandbox you want to keep running or stopped.
- A later call on the handle goes to Daytona.
- Client close releases the SDK connection. A later call returns Daytona's result.

## Run

```sh
cd examples/go/daytona/close-daytona-handle
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
- Handle close does not change the cloud sandbox. Client close releases the SDK connection.

## Related

- [typescript twin](../../../typescript/daytona/close-daytona-handle/README.md)

- Catalog: [examples/go/README.md](../../README.md)

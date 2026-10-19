# Refresh a Daytona sandbox

`Get` a sandbox, confirm **Refresh** is supported, then refresh cached metadata.

## What this means

`Refresh` re-fetches Daytona metadata onto the handle. It does not start, stop,
or delete the sandbox. Useful after another process may have changed state.

## Edge cases

- Capabilities can change after refresh if class or terminal state changed.
- Not a substitute for lifecycle operations.
- Requires a live sandbox ID.

## Run

```sh
cd examples/go/daytona/refresh-daytona-sandbox
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.9+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
DAYTONA_API_KEY=
DAYTONA_SANDBOX_ID=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Needs an **existing sandbox ID** in `.env`.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/daytona/refresh-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

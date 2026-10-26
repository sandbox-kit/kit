# Close a Daytona sandbox handle

`Get` a sandbox, then **`Sandbox.Close` / `sandbox.close()`** to drop the local
handle only.

## What this means

Handle close does not stop, archive, or delete the Daytona sandbox. Client close
likewise only releases the SDK connection.

## Edge cases

- Safe against a sandbox you want to keep running or stopped.
- Do not use the handle after close (or after client close).
- Close is idempotent and waits for in-flight ops on that handle.

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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/daytona/close-daytona-handle/README.md)

- Catalog: [examples/go/README.md](../../README.md)

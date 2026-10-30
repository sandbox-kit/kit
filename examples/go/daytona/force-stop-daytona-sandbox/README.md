# Force-stop a Daytona sandbox

`Get` a sandbox, check **Stop**, then stop it with `Force` set.

## What this means

A normal stop asks the sandbox to shut down. Force stop tells Daytona to stop
when a graceful stop may not finish. `Force` sits on the same options value as
the call timeout.

## Edge cases

- Use [stop](../stop-daytona-sandbox/README.md) when a clean shutdown is enough.
- Force stop still requires Stop on this provider and a live sandbox ID.
- The call changes the sandbox named in `.env`.

## Run

```sh
cd examples/go/daytona/force-stop-daytona-sandbox
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
- **Changes or destroys** that sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Daytona.

## Related

- [typescript twin](../../../typescript/daytona/force-stop-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

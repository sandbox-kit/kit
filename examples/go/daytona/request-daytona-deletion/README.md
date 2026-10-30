# Request Daytona sandbox deletion without waiting

`Get` a sandbox, check **Delete**, then delete with `Wait` set to false.

## What this means

[delete-daytona-sandbox](../delete-daytona-sandbox/README.md) waits until
Daytona finishes. This call returns after Daytona accepts the delete. The
handle stores the metadata from that response.

## Edge cases

- Call `Get` or `Refresh` later when you need to confirm the sandbox is gone.
- The call destroys the sandbox named in `.env`.

## Run

```sh
cd examples/go/daytona/request-daytona-deletion
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

- [typescript twin](../../../typescript/daytona/request-daytona-deletion/README.md)

- Catalog: [examples/go/README.md](../../README.md)

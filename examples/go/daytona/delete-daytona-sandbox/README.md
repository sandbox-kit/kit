# Delete a Daytona sandbox

`Get` a sandbox, check **Delete**, then delete with **`wait: true`** until Daytona
finishes removal.

## What this means

Delete asks Daytona to remove the sandbox. The handle stores the metadata
Daytona returns. A later lifecycle call goes to Daytona, and Kit returns that
result.

## Edge cases

- **Destructive.** The ID in `.env` will be gone.
- `wait: false` returns earlier without proving completion; this example waits.
- Deleting a parent also deletes linked children (cascade).

## Run

```sh
cd examples/go/daytona/delete-daytona-sandbox
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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/daytona/delete-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

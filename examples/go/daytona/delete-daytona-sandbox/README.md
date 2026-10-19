# Delete a Daytona sandbox

`Get` a sandbox, check **Delete**, then delete with **`wait: true`** until Daytona
finishes removal.

## What this means

Delete irreversibly removes the sandbox resource. After success, Kit marks the
handle terminal and rejects further lifecycle calls.

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
- **Changes or destroys** that sandbox.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/daytona/delete-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

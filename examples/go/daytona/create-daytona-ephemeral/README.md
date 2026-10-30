# Create an ephemeral Daytona sandbox

Create from the default snapshot with **`lifetime.ephemeral = true`**.

## What this means

Daytona sandboxes persist by default. An **ephemeral** sandbox is deleted as soon
as it stops, discarding state and avoiding stopped-state disk billing. Daytona
sets auto-delete to immediate (`0`) for ephemeral creates. See
[persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- Stopping an ephemeral sandbox deletes it; do not expect to start it again.
- Ephemeral conflicts with positive idle-pause, archive, or delayed-delete
  policies—Kit rejects those combinations locally.
- [Linked children](../create-daytona-linked/README.md) must be ephemeral.

## Run

```sh
cd examples/go/daytona/create-daytona-ephemeral
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
DAYTONA_API_KEY=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/daytona/create-daytona-ephemeral/README.md)

- Catalog: [examples/go/README.md](../../README.md)

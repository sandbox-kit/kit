# Create linked Daytona sandboxes

Create a parent sandbox, then an ephemeral child with
`network.linkedSandbox = parent.id`. Prints both IDs.

## What this means

Linked sandboxes attach ephemeral children to a parent. Daytona schedules each
child on the same runner and joins them on a shared link network so the group can
reach each other by sandbox ID or name (DNS aliases). Children are always
ephemeral (auto-delete `0`). Deleting the parent cascades to linked children.
See [sandboxes → linked sandboxes](https://www.daytona.io/docs/en/sandboxes/).

This example creates **two** cloud sandboxes.

## Edge cases

- Kit requires `lifetime.ephemeral = true` on the child or rejects before the API.
- If parent create succeeds and child create fails, the parent still exists.
- Clean up both IDs when you are done (or delete the parent to cascade).

## Run

```sh
cd examples/go/daytona/create-daytona-linked
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

- [typescript twin](../../../typescript/daytona/create-daytona-linked/README.md)

- Catalog: [examples/go/README.md](../../README.md)

# Create a Daytona sandbox from a snapshot

Create one sandbox from the stock **`daytona-small`** snapshot
(`source.snapshot.reference`).

## What this means

Snapshots are named boot images. Daytona documents stock sizes such as
`daytona-small` (1 vCPU / 1 GiB / 3 GiB, container class). Any number of sandboxes
can start from the same snapshot. See
[snapshots](https://www.daytona.io/docs/en/snapshots/) and the snapshot table in
[sandboxes](https://www.daytona.io/docs/en/sandboxes/).

Kit sends the snapshot reference string. Explicit snapshot-kind enums and
warm-pool sources are not mapped and are rejected if supplied.

## Edge cases

- The snapshot name must exist for your organization/target; missing names fail
  at the provider.
- Resource overrides still need an image source; snapshot create inherits the
  snapshot's CPU/memory/disk.
- The snapshot supplies the sandbox class as metadata. Capabilities stay the methods Daytona implements.

## Run

```sh
cd examples/go/daytona/create-daytona-from-snapshot
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

- [typescript twin](../../../typescript/daytona/create-daytona-from-snapshot/README.md)

- Catalog: [examples/go/README.md](../../README.md)

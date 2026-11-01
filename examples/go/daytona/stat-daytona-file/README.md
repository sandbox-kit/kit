# Stat a Daytona path

Print the path, type, and size of `DAYTONA_REMOTE_PATH`. The program loads that sandbox with `client.get`.

## What this means

The call prints the path, type, and size. Daytona reports a file or a directory.

## Edge cases

- A relative path is rejected before Daytona is called.
- A missing path comes back from Daytona as not found.

## Run

```sh
cd examples/go/daytona/stat-daytona-file
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
DAYTONA_REMOTE_PATH=/tmp/note.txt

# Optional Daytona endpoint and target overrides.
# DAYTONA_API_URL=https://app.daytona.io/api
# DAYTONA_TARGET=

```

## Warnings

- Needs an existing sandbox and a path at `DAYTONA_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/daytona/stat-daytona-file/README.md)
- Catalog: [examples/go/README.md](../../README.md)

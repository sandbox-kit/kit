# Remove a Daytona path

Delete `DAYTONA_REMOTE_PATH`. The program loads that sandbox with `client.get`.

## What this means

The call deletes `DAYTONA_REMOTE_PATH`. With `DAYTONA_RECURSIVE` unset, it removes a file or an empty directory. Set `DAYTONA_RECURSIVE=true` to delete a directory and its contents.

## Edge cases

- A relative path is rejected before Daytona is called.
- A missing path comes back from Daytona as not found.
- A directory that still has entries fails when `DAYTONA_RECURSIVE` is unset.

## Run

```sh
cd examples/go/daytona/remove-daytona-file
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
# DAYTONA_RECURSIVE=true

```

## Warnings

- Deletes the remote path. A recursive delete removes the directory and everything inside it.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/daytona/remove-daytona-file/README.md)
- Catalog: [examples/go/README.md](../../README.md)

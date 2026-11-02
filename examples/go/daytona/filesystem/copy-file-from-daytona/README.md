# Copy a file from Daytona

Copy `DAYTONA_REMOTE_PATH` from the sandbox to `DAYTONA_LOCAL_PATH` on this machine. The program loads that sandbox with `client.get`.

## What this means

The call copies `DAYTONA_REMOTE_PATH` from the sandbox to `DAYTONA_LOCAL_PATH` on this machine.

## Edge cases

- The sandbox path must be absolute. A blank local path is rejected before Daytona is called. The error field is `local_path`.
- A missing remote file comes back from Daytona as not found.
- The local parent directory must exist. A missing parent is not found before the file is downloaded, and the field is `local_path`. A local directory is an invalid argument with the same field.
- An existing local file is replaced.

## Run

```sh
cd examples/go/daytona/filesystem/copy-file-from-daytona
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
DAYTONA_LOCAL_PATH=./note.txt

```

## Warnings

- Overwrites the local file when it already exists.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/daytona/filesystem/copy-file-from-daytona/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

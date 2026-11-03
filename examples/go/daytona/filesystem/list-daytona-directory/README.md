# List a Daytona directory

Print the names in the directory at `DAYTONA_REMOTE_PATH`. The program loads that sandbox with `client.get`.

## What this means

The call prints the names in one directory. It does not walk into subdirectories.

## Edge cases

- A relative path is rejected before Daytona is called.
- A missing directory comes back from Daytona as not found.
- A file path fails at the provider.

## Run

```sh
cd examples/go/daytona/filesystem/list-daytona-directory
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
DAYTONA_REMOTE_PATH=/tmp

# Optional Daytona endpoint and target overrides.
# DAYTONA_API_URL=https://app.daytona.io/api
# DAYTONA_TARGET=

```

## Warnings

- Needs an existing sandbox and a directory at `DAYTONA_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/daytona/filesystem/list-daytona-directory/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

# Run a Daytona command with environment variables

Run `printenv KIT_DEMO` with environment variables for that call only. The
program loads that sandbox with `client.get`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Daytona's basic exec
returns stdout only; `stderr` is empty.

## Edge cases

- An omitted `environment` map leaves the Daytona default. An empty map is
  explicit.
- These variables apply to this command, not to later calls on the same handle.

## Run

```sh
cd examples/go/daytona/process/run-daytona-with-env
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

# Optional Daytona endpoint and target overrides.
# DAYTONA_API_URL=https://app.daytona.io/api
# DAYTONA_TARGET=

```

## Warnings

- Needs an existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes
  your connection to the provider.

## Related

- [typescript twin](../../../../typescript/daytona/process/run-daytona-with-env/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

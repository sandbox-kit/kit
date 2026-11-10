# Check a Blaxel non-zero exit

Run a command that exits non-zero and read the exit code from the result.

## What this means

A non-zero exit is not an SDK error. Applications check `exitCode` / `GetExitCode`.

## Edge cases

- Transport and validation failures still raise errors.

## Run

```sh
cd examples/go/blaxel/process/run-blaxel-check-exit
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
BL_API_KEY=
BL_WORKSPACE=
# Sandbox name (Blaxel id is the name).
BL_SANDBOX_NAME=
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/process/run-blaxel-check-exit/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

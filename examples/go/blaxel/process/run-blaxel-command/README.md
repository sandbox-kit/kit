# Run a Blaxel command

Run one argv command in a Blaxel sandbox and print exit code plus streams.

## What this means

`Process.Run` is blocking. Set exactly one of `argv` or `shell`. Argv is sent as a POSIX-quoted shell string.

## Edge cases

- Non-zero exit codes are success for the SDK—read `exitCode`.
- Blaxel returns stdout and stderr.

## Run

```sh
cd examples/go/blaxel/process/run-blaxel-command
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
BL_SANDBOX_ID=
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/process/run-blaxel-command/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

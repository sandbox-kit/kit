# Run a Blaxel shell command

Run one shell string in a Blaxel sandbox.

## What this means

Uses `shell` instead of `argv`. Exactly one of the two must be set.

## Edge cases

- Setting both argv and shell is rejected.

## Run

```sh
cd examples/go/blaxel/process/run-blaxel-shell-command
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

- [typescript twin](../../../../typescript/blaxel/process/run-blaxel-shell-command/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

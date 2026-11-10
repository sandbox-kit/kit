# Run a Blaxel command with env

Run a command with environment variables for that process.

## What this means

Process environment is separate from create-time sandbox environment.

## Edge cases

- Invalid names follow the same env-name rules as create.

## Run

```sh
cd examples/go/blaxel/process/run-blaxel-with-env
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

- [typescript twin](../../../../typescript/blaxel/process/run-blaxel-with-env/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

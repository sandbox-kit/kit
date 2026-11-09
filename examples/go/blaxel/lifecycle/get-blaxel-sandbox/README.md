# Get a Blaxel sandbox

Look up an existing Blaxel sandbox by name and print metadata.

## What this means

`Get` retrieves the sandbox by name (the Blaxel id). Use the name printed by a create example as `BL_SANDBOX_ID`.

## Edge cases

- Unknown names return not found.
- Attach is unsupported; use get.

## Run

```sh
cd examples/go/blaxel/lifecycle/get-blaxel-sandbox
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

- [typescript twin](../../../../typescript/blaxel/lifecycle/get-blaxel-sandbox/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

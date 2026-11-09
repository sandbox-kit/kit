# Refresh a Blaxel sandbox

Refresh metadata for an existing Blaxel sandbox by name.

## What this means

`Refresh` re-fetches provider metadata onto the handle.

## Edge cases

- Unsupported when the handle does not report the refresh capability.

## Run

```sh
cd examples/go/blaxel/lifecycle/refresh-blaxel-sandbox
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

- [typescript twin](../../../../typescript/blaxel/lifecycle/refresh-blaxel-sandbox/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

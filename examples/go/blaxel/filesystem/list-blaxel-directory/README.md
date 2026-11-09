# List a Blaxel directory

List entries in one absolute directory path.

## What this means

Returns shared directory entries with names and types.

## Edge cases

- Missing directories return not found.

## Run

```sh
cd examples/go/blaxel/filesystem/list-blaxel-directory
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
# Absolute path inside the sandbox.
BL_REMOTE_PATH=/tmp
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/filesystem/list-blaxel-directory/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

# Read a Blaxel file

Read one absolute path from a Blaxel sandbox.

## What this means

Uses `filesystem.readFile` / `ReadFile`. Paths in the sandbox are absolute.

## Edge cases

- Missing files return not found from Blaxel.

## Run

```sh
cd examples/go/blaxel/filesystem/read-blaxel-file
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
BL_REMOTE_PATH=/tmp/kit-example.txt
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/filesystem/read-blaxel-file/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

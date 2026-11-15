# Copy a file from Blaxel

Copy a sandbox file out to a local path.

## What this means

`copyToLocal` creates a missing local parent directory.

## Edge cases

- Remote path must exist.

## Run

```sh
cd examples/go/blaxel/filesystem/copy-file-from-blaxel
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
# Absolute path inside the sandbox.
BL_REMOTE_PATH=/tmp/kit-example.txt
# Absolute local path.
BL_LOCAL_PATH=./local-example.txt
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/filesystem/copy-file-from-blaxel/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

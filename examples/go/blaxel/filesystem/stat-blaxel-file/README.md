# Stat a Blaxel path

Stat one absolute path in a Blaxel sandbox.

## What this means

`Stat` returns `Type` (file or directory), `Size`, and usually `Permissions`,
`Owner`, `Group`, and `ModifiedAt`. `Name` and `Path` are whatever Blaxel
reports. Size may be zero when Blaxel does not return it.

## Edge cases

- Missing paths return not found.

## Run

```sh
cd examples/go/blaxel/filesystem/stat-blaxel-file
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
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/filesystem/stat-blaxel-file/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

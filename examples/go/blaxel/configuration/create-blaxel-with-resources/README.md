# Create a Blaxel sandbox with memory

Create one Blaxel sandbox with an explicit memory size in MiB. Blaxel derives CPU from memory.

## What this means

Only `resources.memoryMib` is mapped. CPU cores, limits, and disk MiB overrides are rejected. Omitted memory keeps the Blaxel default.

## Edge cases

- Memory values that Blaxel rejects are returned as provider errors.

## Run

```sh
cd examples/go/blaxel/configuration/create-blaxel-with-resources
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
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/configuration/create-blaxel-with-resources/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

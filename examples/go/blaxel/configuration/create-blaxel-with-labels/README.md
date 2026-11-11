# Create a Blaxel sandbox with labels and env

Create one Blaxel sandbox with labels and environment variables.

## What this means

Set `labels` and create-time environment variables. Environment variables are
available inside the sandbox. `Info().Labels` is set only when Blaxel reports
labels back — Kit does not copy them from the create request.

## Edge cases

- Invalid environment names are rejected before Blaxel is called.

## Run

```sh
cd examples/go/blaxel/configuration/create-blaxel-with-labels
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

- [typescript twin](../../../../typescript/blaxel/configuration/create-blaxel-with-labels/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

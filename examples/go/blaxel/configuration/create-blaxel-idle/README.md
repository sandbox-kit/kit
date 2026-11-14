# Create a Blaxel sandbox with idle delete

Create one Blaxel sandbox that deletes after an idle period.

## What this means

Set `lifetime.idleTerminate` with `after` and a delay (≥1 minute). The sandbox deletes after that idle period. You cannot disable idle terminate.

## Edge cases

- Disabling idle terminate is not representable and is rejected.
- Minimum idle duration follows Blaxel rules (1 minute).

## Run

```sh
cd examples/go/blaxel/configuration/create-blaxel-idle
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

- [typescript twin](../../../../typescript/blaxel/configuration/create-blaxel-idle/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

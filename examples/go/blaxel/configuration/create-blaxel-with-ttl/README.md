# Create a Blaxel sandbox with maximum lifetime

Create one Blaxel sandbox that deletes after a maximum lifetime TTL.

## What this means

Maps `lifetime.maximumLifetime` to Blaxel max-age TTL. Durations must be whole seconds and meet Blaxel minima (5 minutes for max-age).

## Edge cases

- Fractional-second durations are rejected.
- Idle cleanup uses `idleTerminate`, not maximum lifetime.

## Run

```sh
cd examples/go/blaxel/configuration/create-blaxel-with-ttl
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

- [typescript twin](../../../../typescript/blaxel/configuration/create-blaxel-with-ttl/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

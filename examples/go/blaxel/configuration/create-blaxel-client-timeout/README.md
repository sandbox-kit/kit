# Create a Blaxel sandbox with a client timeout

Create one Blaxel sandbox using a shorter client-wide timeout.

## What this means

Sets `Config.Timeout` / `timeout` on the shared client. This bounds local waiting for SDK calls, including create.

## Edge cases

- Create provisioning timeouts are not mapped on Go; use the client timeout instead.
- TypeScript can also pass a create `provisioning.timeout` (1–50 whole seconds).

## Run

```sh
cd examples/go/blaxel/configuration/create-blaxel-client-timeout
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

- [typescript twin](../../../../typescript/blaxel/configuration/create-blaxel-client-timeout/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

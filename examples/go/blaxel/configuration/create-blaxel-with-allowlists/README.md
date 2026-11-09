# Create a Blaxel sandbox with outbound domains

Create one Blaxel sandbox that allows selected outbound domains.

## What this means

Maps `network.outboundDomains` to Blaxel `allowedDomains`. Other egress controls are unsupported.

## Edge cases

- Restricted egress modes and CIDR allowlists are rejected.
- Domain entries follow Blaxel validation.

## Run

```sh
cd examples/go/blaxel/configuration/create-blaxel-with-allowlists
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

- [typescript twin](../../../../typescript/blaxel/configuration/create-blaxel-with-allowlists/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

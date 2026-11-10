# Create a Blaxel sandbox with ports

Create one Blaxel sandbox that publishes selected ports.

## What this means

Maps `network.ports` to Blaxel runtime ports. Default leaves the protocol to
Blaxel (HTTP). Plain → TCP, TLS → TLS. HTTP/2 over TLS is rejected.

## Edge cases

- Other network settings (egress mode, CIDRs, private network, proxies) are rejected.
- Port numbers must be in range.
- Blaxel reserves ports 80, 443, and 8080; this example uses 3000 and 3001.

## Run

```sh
cd examples/go/blaxel/configuration/create-blaxel-with-ports
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

- [typescript twin](../../../../typescript/blaxel/configuration/create-blaxel-with-ports/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

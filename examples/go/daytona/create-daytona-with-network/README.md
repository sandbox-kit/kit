# Create a Daytona sandbox with restricted egress

Create from an image with `EgressModeRestricted`,
`OutboundCidrs` `192.0.2.0/24`, and `OutboundDomains` `example.com`.

## What this means

Restricted egress becomes Daytona network allowlists. Kit checks CIDR and
domain entries before the provider call.

## Edge cases

- An empty allowlist is rejected.
- An invalid CIDR or domain fails before the provider call.
- Ports, inbound CIDRs, and a linked sandbox are separate requests.

## Run

```sh
cd examples/go/daytona/create-daytona-with-network
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
DAYTONA_API_KEY=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Daytona.

## Related

- [typescript twin](../../../typescript/daytona/create-daytona-with-network/README.md)

- Catalog: [examples/go/README.md](../../README.md)

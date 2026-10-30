# Create a Daytona sandbox with restricted network egress

Create from an image with `EgressMode.Restricted`,
`outboundCidrs: ["192.0.2.0/24"]`, and `outboundDomains: ["example.com"]`.

## What this means

Daytona maps restricted egress to network allowlists. Kit validates CIDR/domain entries locally before the provider call.

## Edge cases

- Empty allowlists are rejected by the Daytona mapper.
- Invalid CIDRs/domains fail Kit provider checks before the native SDK runs.
- Other network fields (ports, inbound CIDRs, linked sandbox) are covered by other examples or unsupported here.

## Run

```sh
cd examples/typescript/daytona/create-daytona-with-network
cp .env.example .env
# Fill in the values below, then:
npm install
npm run build
npm start
```

Requires Node.js 20/22/24.15+, npm 11.19.0, and this checkout. Uses `await using`
so local client cleanup runs on success or failure.

### Environment

```dotenv
DAYTONA_API_KEY=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)

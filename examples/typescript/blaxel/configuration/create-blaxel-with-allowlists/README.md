# Create a Blaxel sandbox with outbound domains

Create one Blaxel sandbox that allows selected outbound domains.

## What this means

Set `network.outboundDomains` to allow those hosts. Other egress controls are unsupported on Blaxel.

## Edge cases

- Restricted egress modes and CIDR allowlists are rejected.
- Domain entries follow Blaxel validation.

## Run

```sh
cd examples/typescript/blaxel/configuration/create-blaxel-with-allowlists
cp .env.example .env
# Fill in the values below, then:
npm install
npm run build
npm start
```

Requires Node.js 20/22/24.15+, npm 11.19.0, and this checkout. Skip copying if
`.env` already exists.

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

- [go twin](../../../../go/blaxel/configuration/create-blaxel-with-allowlists/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)

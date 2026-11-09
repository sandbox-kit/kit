# Create a Blaxel sandbox with ports

Create one Blaxel sandbox that publishes selected ports.

## What this means

Maps `network.ports` to Blaxel runtime ports. Default/plain → TCP, TLS → TLS, HTTP/2 over TLS → HTTP.

## Edge cases

- Other network settings (egress mode, CIDRs, private network, proxies) are rejected.
- Port numbers must be in range.
- Blaxel reserves ports 80, 443, and 8080; this example uses 3000 and 3001.

## Run

```sh
cd examples/typescript/blaxel/configuration/create-blaxel-with-ports
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

- [go twin](../../../../go/blaxel/configuration/create-blaxel-with-ports/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)

# Handle Daytona create errors

Call create with **`cpuCores: -1`** so local validation fails with
`ErrorKindInvalidArgument`. Prints shared error fields. No sandbox is created
when validation rejects the request.

## What this means

Kit checks Daytona resource rules before calling the API (CPU must be a positive
whole core count). Failures expose portable `kind`, `field`, `operation`, and
`message`. Native causes remain available when an SDK error is wrapped.

## Edge cases

- This path usually never reaches Daytona's HTTP API.
- Other failures can occur after provision; kinds are not safe retry signals.

## Run

```sh
cd examples/typescript/daytona/handle-daytona-errors
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

- Intended to fail **locally** before a provider create call.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/daytona/handle-daytona-errors/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)

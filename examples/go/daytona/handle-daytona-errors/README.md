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
- Missing `.env` is an ordinary error in Go examples that load dotenv.

## Run

```sh
cd examples/go/daytona/handle-daytona-errors
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.9+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
DAYTONA_API_KEY=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Intended to fail **locally** before a provider create call.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/daytona/handle-daytona-errors/README.md)

- Catalog: [examples/go/README.md](../../README.md)

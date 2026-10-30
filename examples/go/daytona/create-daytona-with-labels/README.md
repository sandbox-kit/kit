# Create a Daytona sandbox with labels and environment

Create from an image with `Labels` and `Environment` maps.

## What this means

Labels become Daytona sandbox labels. Environment variables are passed into
the sandbox process environment.

## Edge cases

- Environment keys must be valid env names (`[A-Za-z_][A-Za-z0-9_]*`).
- Label values are opaque strings to Kit. Daytona enforces its own size limits.

## Run

```sh
cd examples/go/daytona/create-daytona-with-labels
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

- [typescript twin](../../../typescript/daytona/create-daytona-with-labels/README.md)

- Catalog: [examples/go/README.md](../../README.md)

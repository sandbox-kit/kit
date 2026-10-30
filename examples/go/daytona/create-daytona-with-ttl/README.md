# Create a Daytona sandbox with a maximum lifetime

Create from an image with `Lifetime.MaximumLifetime` set to one hour.

## What this means

Daytona maps that lifetime to a TTL in minutes. Kit requires a positive
whole number of minutes.

## Edge cases

- A duration that is not a whole number of minutes is rejected locally.
- Idle stop, idle pause, and stopped archive or delete policies are in
  [create-daytona-with-policies](../create-daytona-with-policies/README.md).

## Run

```sh
cd examples/go/daytona/create-daytona-with-ttl
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

- [typescript twin](../../../typescript/daytona/create-daytona-with-ttl/README.md)

- Catalog: [examples/go/README.md](../../README.md)

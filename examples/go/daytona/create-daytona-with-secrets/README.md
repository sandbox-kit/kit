# Create a Daytona sandbox with an injected secret

Create from an image and inject an existing Daytona secret as `GITHUB_TOKEN`.

## What this means

The secret needs an environment variable name and egress-placeholder
injection. It must already exist in the Daytona organization.

## Edge cases

- `DAYTONA_SECRET_NAME` must name a secret that already exists.
- A duplicate environment variable in `Security.Secrets` is rejected locally.

## Run

```sh
cd examples/go/daytona/create-daytona-with-secrets
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
DAYTONA_SECRET_NAME=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Daytona.

## Related

- [typescript twin](../../../typescript/daytona/create-daytona-with-secrets/README.md)

- Catalog: [examples/go/README.md](../../README.md)

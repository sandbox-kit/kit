# Start a Daytona sandbox with a timeout

`Get` a sandbox, check **Start**, then start it with a 60-second deadline.

## What this means

The example sets both a 60-second context deadline and
`StartOptions.Timeout` of one minute. That options value is the deadline Kit
uses and the start timeout Daytona receives, in whole seconds.

## Edge cases

- Omitting `Timeout` leaves Daytona's 60-second start timeout in place. The client `Timeout` still limits how long Kit waits.
- An ephemeral sandbox that has already stopped cannot be started again.
- The call changes the sandbox named in `.env`.

## Run

```sh
cd examples/go/daytona/start-daytona-with-timeout
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
DAYTONA_API_KEY=
DAYTONA_SANDBOX_ID=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Needs an **existing sandbox ID** in `.env`.
- **Changes** that sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Daytona.

## Related

- [typescript twin](../../../typescript/daytona/start-daytona-with-timeout/README.md)

- Catalog: [examples/go/README.md](../../README.md)

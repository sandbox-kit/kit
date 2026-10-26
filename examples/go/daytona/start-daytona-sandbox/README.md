# Start a Daytona sandbox

`Get` a sandbox, check **Start**, then start it from a retained stopped (or
otherwise startable) state.

## What this means

Daytona sandboxes that were stopped (and not deleted) keep their filesystem and
can be started again. See [persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- Ephemeral sandboxes that already stopped are gone—`Get` fails.
- GPU start may still be denied by account/features even when the capability appears.
- Mutates the sandbox named in `.env`.

## Run

```sh
cd examples/go/daytona/start-daytona-sandbox
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
- **Changes or destroys** that sandbox.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/daytona/start-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

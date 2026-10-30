# Stop a Daytona sandbox

`Get` a sandbox, check **Stop**, then stop it.

## What this means

Stop asks Daytona to stop the sandbox. A retained sandbox can usually be started
again. An ephemeral sandbox can be deleted by its policy. See
[persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- Daytona applies its own stop rules, including ephemeral deletion.
- Ephemeral + stop deletes the sandbox.
- **Mutates** the ID in `.env`.

## Run

```sh
cd examples/go/daytona/stop-daytona-sandbox
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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/daytona/stop-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

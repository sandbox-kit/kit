# Archive a Daytona sandbox

`Get` a sandbox, then archive it.

## What this means

Archive asks Daytona to move the sandbox filesystem to object storage. Daytona
accepts or rejects that sandbox. See
[persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- A rejected sandbox returns Daytona's error.
- Archive is not delete; recovery rules are Daytona-specific afterward.
- Mutates the target sandbox.

## Run

```sh
cd examples/go/daytona/archive-daytona-sandbox
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

- [typescript twin](../../../typescript/daytona/archive-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

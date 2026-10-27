# Pause a Daytona sandbox

`Get` a sandbox, then pause it.

## What this means

Pause asks Daytona to freeze the sandbox. Daytona accepts or rejects that
sandbox. See Daytona's lifecycle docs under
[sandboxes](https://www.daytona.io/docs/en/sandboxes/).

## Edge cases

- A rejected sandbox returns Daytona's error.
- Inspect refreshed state after success for the resolved paused/running value.
- Mutates the target sandbox.

## Run

```sh
cd examples/go/daytona/pause-daytona-sandbox
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

- [typescript twin](../../../typescript/daytona/pause-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

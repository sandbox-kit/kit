# Pause a Daytona sandbox

`Get` a sandbox, check **Pause**, then pause it.

## What this means

Pause is for VM sandbox classes (Linux VM / Windows): it preserves filesystem and
memory. Container and GPU classes do not support pause—Kit rejects via
capabilities before the API call. See Daytona's lifecycle docs under
[sandboxes](https://www.daytona.io/docs/en/sandboxes/).

## Edge cases

- Using a container sandbox ID fails the capability check.
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

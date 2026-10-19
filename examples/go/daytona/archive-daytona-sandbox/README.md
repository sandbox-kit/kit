# Archive a Daytona sandbox

`Get` a sandbox, check **Archive**, then archive retained container filesystem
state.

## What this means

Archive moves a stopped container sandbox's filesystem to object storage
(Daytona auto-archive defaults to 7 days stopped, max 30 days). VM/GPU classes
do not support archive. See
[persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- Wrong class → unsupported before the call.
- Archive is not delete; recovery rules are Daytona-specific afterward.
- Mutates the target sandbox.

## Run

```sh
cd examples/go/daytona/archive-daytona-sandbox
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.9+ and this checkout. Skip copying if `.env` already exists.
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

- [typescript twin](../../../typescript/daytona/archive-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

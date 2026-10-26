# Stop a Daytona sandbox

`Get` a sandbox, check **Stop**, then stop it.

## What this means

Stop leaves a retained sandbox that can usually be started again. Containers keep
the filesystem; VMs clear memory on stop. **Ephemeral** sandboxes are deleted when
they stop. GPU sandboxes may be destroyed on stop by Daytona design. See
[persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- Always check `Supports(Stop)` for the resolved class.
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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/daytona/stop-daytona-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

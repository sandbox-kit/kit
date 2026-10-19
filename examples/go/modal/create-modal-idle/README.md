# Create a Modal sandbox with idle termination

Create with **`idleTerminate` after 60 seconds** of inactivity.

## What this means

Modal idle timeouts automatically terminate a sandbox after inactivity (no
active exec/tunnels/etc., per Modal's idle rules). Separate from the maximum
lifetime `timeout`. See
[Sandboxes → Idle Timeouts](https://modal.com/docs/guide/sandboxes).

## Edge cases

- Delay must be a positive whole number of seconds when using `After`.
- The sandbox may disappear soon after create if it goes idle—expected.
- `Disabled` turns idle terminate off; zero delay is not “disabled.”

## Run

```sh
cd examples/go/modal/create-modal-idle
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.9+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
MODAL_TOKEN_ID=
MODAL_TOKEN_SECRET=
MODAL_APP_NAME=
MODAL_ENVIRONMENT=main
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/modal/create-modal-idle/README.md)

- Catalog: [examples/go/README.md](../../README.md)

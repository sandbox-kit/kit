# Terminate a Modal sandbox (wait)

Attach by ID, check **Terminate**, then terminate with **`wait: true`** and print
the exit code.

## What this means

`terminate` ends the sandbox execution. Modal also ends sandboxes when max
timeout or idle timeout fires, or when the entrypoint exits. See
[Sandboxes](https://modal.com/docs/guide/sandboxes). Waiting returns an exit
result when available. A later call goes to Modal, and Kit returns that result.

## Edge cases

- **Destructive** to the running sandbox.
- Waiting can still error after termination has begun.
- Use [request-modal-termination](../request-modal-termination/README.md) when you
  do not need the exit code.

## Run

```sh
cd examples/go/modal/terminate-modal-sandbox
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
MODAL_TOKEN_ID=
MODAL_TOKEN_SECRET=
MODAL_APP_NAME=
MODAL_ENVIRONMENT=main
MODAL_SANDBOX_ID=
```

## Warnings

- Needs an **existing sandbox ID** in `.env`.
- **Changes or destroys** that sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/terminate-modal-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

# Attach to a Modal sandbox by ID

Build a handle from `MODAL_SANDBOX_ID` with `Client.Attach` / `client.attach`.

## What this means

Modal's SDKs expose ID-based sandbox handles (for example `Sandbox.from_id` /
equivalent). Kit's `Attach` matches that: construct a handle from an ID without
claiming a full remote metadata refresh. You can then terminate or close the
handle.

## Edge cases

- A wrong ID may only fail when you invoke an operation, not at attach time.
- Attach is not proof the sandbox is still running.
- Does not terminate the sandbox by itself.

## Run

```sh
cd examples/go/modal/attach-modal-sandbox
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
MODAL_SANDBOX_ID=
```

## Warnings

- Needs an **existing sandbox ID** in `.env`.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/modal/attach-modal-sandbox/README.md)

- Catalog: [examples/go/README.md](../../README.md)

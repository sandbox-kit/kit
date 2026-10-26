# Handle Modal unsupported-intent errors

Attempt create with **`diskMib` set**. Modal does not allow increasing sandbox
disk at create, so Kit returns **`ErrorKindUnsupported`** before the API call.

## What this means

Modal documents that Sandbox disk size cannot be increased (VM root disk is a
fixed size). Kit rejects `diskMib` as unsupported intent and exposes portable
error fields. See
[sandbox resources → Disk](https://modal.com/docs/guide/sandbox-resources).

## Edge cases

- No sandbox should be created on this intentional failure path.
- Good template for switching on `kind` / `field` / `operation`.
- Kinds are not safe automatic retry signals for other failures.

## Run

```sh
cd examples/go/modal/handle-modal-errors
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
```

## Warnings

- Intended to fail **locally** before a provider create call.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/modal/handle-modal-errors/README.md)

- Catalog: [examples/go/README.md](../../README.md)

# Handle Modal unsupported-intent errors

Attempt create with **`DiskMib` set**. Modal does not allow a disk override at
creation, so Kit returns **`ErrorKindUnsupported`** and field
**`resources.disk_mib`** before the API call.

## What this means

Modal documents that Sandbox disk size cannot be increased (VM root disk is a
fixed size). Kit rejects the override before calling Modal. The portable field
is `resources.disk_mib`. See
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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/handle-modal-errors/README.md)

- Catalog: [examples/go/README.md](../../README.md)

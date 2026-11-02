# Stat a Modal path

Print the path, type, and size of `MODAL_REMOTE_PATH`. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call prints the path, type, and size. A symlink result describes the symlink. Modal reports `mode` as a four-digit octal string when it has one.

## Edge cases

- A relative path is rejected before Modal is called.
- A missing path comes back from Modal as not found.

## Run

```sh
cd examples/go/modal/filesystem/stat-modal-file
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
MODAL_REMOTE_PATH=/tmp/note.txt

```

## Warnings

- Needs an existing sandbox and a path at `MODAL_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/modal/filesystem/stat-modal-file/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

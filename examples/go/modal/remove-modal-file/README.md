# Remove a Modal path

Delete `MODAL_REMOTE_PATH`. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call deletes `MODAL_REMOTE_PATH`. With `MODAL_RECURSIVE` unset, it removes a file or an empty directory. Set `MODAL_RECURSIVE=true` to delete a directory and its contents.

## Edge cases

- A relative path is rejected before Modal is called.
- A missing path comes back from Modal as not found.
- A directory that still has entries fails when `MODAL_RECURSIVE` is unset.

## Run

```sh
cd examples/go/modal/remove-modal-file
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
# MODAL_RECURSIVE=true

```

## Warnings

- Deletes the remote path. A recursive delete removes the directory and everything inside it.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/remove-modal-file/README.md)
- Catalog: [examples/go/README.md](../../README.md)

# Write a Modal file

Replace the file at `MODAL_REMOTE_PATH` with the bytes of `hello from sandbox-kit`. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call replaces the file with the bytes of `hello from sandbox-kit` and creates any missing parent directories.

## Edge cases

- A relative sandbox path is rejected before Modal is called. The error field is `path`.
- An existing file at that path is replaced.
- A path that names a directory fails at the provider.

## Run

```sh
cd examples/go/modal/write-modal-file
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

- Overwrites the file at `MODAL_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/write-modal-file/README.md)
- Catalog: [examples/go/README.md](../../README.md)

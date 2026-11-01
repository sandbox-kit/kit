# Copy a file into Modal

Copy `MODAL_LOCAL_PATH` on this machine to `MODAL_REMOTE_PATH` in the sandbox. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call copies `MODAL_LOCAL_PATH` on this machine to `MODAL_REMOTE_PATH` in the sandbox. Missing parent directories of the remote file are created.

## Edge cases

- The sandbox path must be absolute. A blank local path is rejected before Modal is called. The error field is `local_path`.
- A missing local file is not found. The error field is `local_path`. A directory at that path is an invalid argument with the same field.
- An existing remote file is replaced.

## Run

```sh
cd examples/go/modal/copy-file-into-modal
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
MODAL_LOCAL_PATH=./note.txt

```

## Warnings

- Overwrites the remote file when it already exists.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/copy-file-into-modal/README.md)
- Catalog: [examples/go/README.md](../../README.md)

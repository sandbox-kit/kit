# Read a Modal file

Read the bytes at `MODAL_REMOTE_PATH` with `Filesystem.ReadFile` / `sandbox.filesystem.readFile`. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call returns the file contents as bytes. Decode text in your own program. The sandbox path is absolute.

## Edge cases

- A relative path, a blank path, or a path that contains a NUL byte is rejected before Modal is called. The error field is `path`.
- A missing file comes back from Modal as not found.
- A directory path fails at the provider.

## Run

```sh
cd examples/go/modal/read-modal-file
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

- Needs an existing sandbox and a file at `MODAL_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/read-modal-file/README.md)
- Catalog: [examples/go/README.md](../../README.md)

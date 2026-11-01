# List a Modal directory

Print the names in the directory at `MODAL_REMOTE_PATH`. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call prints the names in one directory. It does not walk into subdirectories.

## Edge cases

- A relative path is rejected before Modal is called.
- A missing directory comes back from Modal as not found.
- A file path fails at the provider.

## Run

```sh
cd examples/go/modal/list-modal-directory
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
MODAL_REMOTE_PATH=/tmp

```

## Warnings

- Needs an existing sandbox and a directory at `MODAL_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/list-modal-directory/README.md)
- Catalog: [examples/go/README.md](../../README.md)

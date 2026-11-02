# Create a Modal directory

Create the directory at `MODAL_REMOTE_PATH`, including any missing parents. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call creates `MODAL_REMOTE_PATH` and any missing parents.

## Edge cases

- A relative path is rejected before Modal is called.
- `Parents` defaults to true: missing parents are created, and an existing directory is accepted.
- `Parents` false requires the parent directory to exist.

## Run

```sh
cd examples/go/modal/filesystem/make-modal-directory
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
MODAL_REMOTE_PATH=/tmp/kit-files

```

## Warnings

- Creates the directory inside the existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/modal/filesystem/make-modal-directory/README.md)
- Catalog: [examples/go/README.md](../../../README.md)

# Create a Modal sandbox with runtime settings

Create with **working directory**, **entrypoint**, and **PTY** enabled.

## What this means

These fields configure the sandbox process: where it starts (`workdir`), what
command runs (`entrypoint`), and whether a PTY is attached. See Modal
[Sandboxes](https://modal.com/docs/guide/sandboxes) for workdir/entrypoint usage.

## Edge cases

- Working directory must be an absolute path.
- Entrypoint replaces the image default command for this sandbox.
- PTY affects interactive IO for exec sessions.

## Run

```sh
cd examples/go/modal/create-modal-with-runtime
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
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/modal/create-modal-with-runtime/README.md)

- Catalog: [examples/go/README.md](../../README.md)

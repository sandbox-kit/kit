# Create a Modal Linux VM sandbox

Create with **`isolation.kind = LinuxVM`**.

## What this means

Requests Modal's Linux VM sandbox runtime. End the sandbox with terminate, idle
timeout, or max lifetime. VM root disk size is fixed per Modal's sandbox resource
docs.

## Edge cases

- Workspace permissions may block VM isolation.
- GPU is not available on the VM runtime path in Modal's docs.

## Run

```sh
cd examples/go/modal/create-modal-linux-vm
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

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/modal/create-modal-linux-vm/README.md)

- Catalog: [examples/go/README.md](../../README.md)

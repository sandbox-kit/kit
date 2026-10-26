# Create a Modal gVisor sandbox

Create with **`isolation.kind = Container`** (Modal's gVisor/container runtime).

## What this means

Modal supports `gvisor` (container) and Linux VM runtimes. GPU sandboxes require
the gVisor runtime. See [Sandboxes](https://modal.com/docs/guide/sandboxes) and
[sandbox resources](https://modal.com/docs/guide/sandbox-resources).

## Edge cases

- Making container isolation explicit does not change terminate/attach support.
- Account permissions can still reject a runtime.

## Run

```sh
cd examples/go/modal/create-modal-gvisor
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

- [typescript twin](../../../typescript/modal/create-modal-gvisor/README.md)

- Catalog: [examples/go/README.md](../../README.md)

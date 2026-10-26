# Create a Modal sandbox with resources

Create with **0.5 physical CPU cores**, **512 MiB** memory, and a **5-minute
maximum lifetime**; print whether allocated resources appear in metadata.

## What this means

Modal `cpu` is physical cores (floating point); `memory` is MiB. Sandboxes can
burst above the request when the host has spare capacity; you pay for
`max(request, actual)`. Optional CPU/memory **limits** cap bursting. See
[sandbox resources](https://modal.com/docs/guide/sandbox-resources) and
[reserving CPU and memory](https://modal.com/docs/guide/resources).

Kit enforces documented floors (0.125 cores, 128 MiB), 0.001-core precision, and
a maximum lifetime up to 24 hours (whole seconds).

## Edge cases

- Sandbox disk size cannot be increased; VM runtime root disk is fixed (512 GiB).
- Metadata resources are provider-reported when present—not an echo of the request.
- Values that pass Kit can still fail Modal account maxima.

## Run

```sh
cd examples/go/modal/create-modal-with-resources
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

- [typescript twin](../../../typescript/modal/create-modal-with-resources/README.md)

- Catalog: [examples/go/README.md](../../README.md)

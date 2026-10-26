# Create a Modal sandbox and wait until ready

Create with an **exec readiness probe** (`true`) and
`provisioning.waitFor = Ready`, plus a positive creation timeout.

## What this means

Modal readiness probes (TCP or exec) run inside the sandbox until they succeed.
Kit's ready wait maps to waiting until the probe passes before returning the
handle. See
[Sandboxes → Readiness Probes](https://modal.com/docs/guide/sandboxes).

Kit requires a probe and a positive create timeout for ready wait.

## Edge cases

- Missing probe or non-positive timeout → local rejection.
- A failing/slow probe can time out create; check the Modal dashboard if unsure
  whether a sandbox remained.
- Entrypoint `sleep 300` keeps the process alive after ready.

## Run

```sh
cd examples/go/modal/create-modal-ready
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

- [typescript twin](../../../typescript/modal/create-modal-ready/README.md)

- Catalog: [examples/go/README.md](../../README.md)

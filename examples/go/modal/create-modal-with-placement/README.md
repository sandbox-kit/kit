# Create a Modal sandbox with placement

Create with a **cloud** preference and optional **region** list.

## What this means

Modal can schedule sandboxes with cloud/region preferences. Allowed cloud names
in Kit checks include `aws`, `gcp`, `oci`, and `auto`. See
[region selection](https://modal.com/docs/guide/region-selection).

## Edge cases

- Invalid cloud names fail local checks.
- Preferences can be ignored under capacity; failures are provider-side.
- Optional env: `MODAL_CLOUD` (default `auto`), `MODAL_REGION`.

## Run

```sh
cd examples/go/modal/create-modal-with-placement
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

Optional: `MODAL_CLOUD` (default `auto`), `MODAL_REGION`.

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/modal/create-modal-with-placement/README.md)

- Catalog: [examples/go/README.md](../../README.md)

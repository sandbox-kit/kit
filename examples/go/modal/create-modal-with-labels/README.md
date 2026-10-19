# Create a Modal sandbox with labels and environment

Create with process **environment** variables and **labels/tags**.

## What this means

Environment variables are injected into the sandbox process. Labels/tags help
organize sandboxes in the Modal UI and APIs. Kit validates environment key syntax
when mapping.

## Edge cases

- Do not put secrets in labels or example env values committed to git.
- Origins on name/labels may be request-derived in returned metadata.

## Run

```sh
cd examples/go/modal/create-modal-with-labels
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

- [typescript twin](../../../typescript/modal/create-modal-with-labels/README.md)

- Catalog: [examples/go/README.md](../../README.md)

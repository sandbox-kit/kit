# Create a Modal sandbox with labels and environment

Create with process **environment** variables and **labels/tags**.

## What this means

Environment variables are injected into the sandbox process. Labels/tags help
organize sandboxes in the Modal UI and APIs. Kit validates environment key syntax
when mapping.

## Edge cases

- Do not put secrets in labels or example env values committed to git.
- Name and labels on the returned metadata are the values from the create request.

## Run

```sh
cd examples/go/modal/create-modal-with-labels
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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/create-modal-with-labels/README.md)

- Catalog: [examples/go/README.md](../../README.md)

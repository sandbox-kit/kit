# Create a Modal sandbox with verbose logging

Create with **`observability.verbose = true`**.

## What this means

Enables verbose provider logging for the sandbox session so Modal surfaces more
detail in logs/UI for debugging.

## Edge cases

- Verbose logs may include sensitive command or environment detail.
- Does not change terminate/attach behavior.

## Run

```sh
cd examples/go/modal/create-modal-verbose
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

- [typescript twin](../../../typescript/modal/create-modal-verbose/README.md)

- Catalog: [examples/go/README.md](../../README.md)

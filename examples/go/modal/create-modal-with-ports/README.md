# Create a Modal sandbox with exposed ports

Create with a TLS port on 443 and a plain port on 8080.

## What this means

Each port sets `PortTransport`. TLS maps to an encrypted port. Plain maps to
an unencrypted port. HTTP/2 TLS maps to an h2 port.

## Edge cases

- Omitting `Transport` fails locally.
- Public HTTP access is expressed with ports.

## Run

```sh
cd examples/go/modal/create-modal-with-ports
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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Modal.

## Related

- [typescript twin](../../../typescript/modal/create-modal-with-ports/README.md)

- Catalog: [examples/go/README.md](../../README.md)

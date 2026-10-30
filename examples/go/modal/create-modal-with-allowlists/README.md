# Create a Modal sandbox with outbound allowlists

Create with `EgressModeRestricted`, outbound CIDR and domain allowlists, and
`PrivateNetwork` set.

## What this means

The request limits outbound traffic and turns on Modal private networking.

## Edge cases

- An empty inbound CIDR list is rejected.
- A linked sandbox and an outbound proxy URL are rejected on Modal.

## Run

```sh
cd examples/go/modal/create-modal-with-allowlists
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

- [typescript twin](../../../typescript/modal/create-modal-with-allowlists/README.md)

- Catalog: [examples/go/README.md](../../README.md)

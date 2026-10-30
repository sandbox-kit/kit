# Create a Modal sandbox with a network proxy

Create with **`network.proxyReference`** set to an existing Modal Proxy name
(`MODAL_PROXY_NAME`).

## What this means

Attaches the sandbox to a named Modal Proxy in the selected environment so
outbound traffic can be routed through that proxy. The proxy must already exist.

## Edge cases

- Missing proxy → provider error at create.
- Requires `MODAL_PROXY_NAME` in addition to token/app vars.
- Empty egress allowlists mean deny-all when you use allowlist fields; this
  example only sets the proxy reference.

## Run

```sh
cd examples/go/modal/create-modal-with-network
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

Required: `MODAL_PROXY_NAME` (existing proxy in the environment).

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/create-modal-with-network/README.md)

- Catalog: [examples/go/README.md](../../README.md)

# Create a Modal sandbox with a managed proxy

Create one sandbox from `alpine:3.21` and route its outbound traffic through an
existing Modal Proxy selected by `MODAL_PROXY_NAME`.

## Run

Requires Go 1.26.9 or newer, this checkout, an existing Modal app, and an
existing Modal Proxy in the selected environment.
See [Modal creation](../create-modal-sandbox/README.md) for the token and app.

```sh
cd examples/go/modal/create-modal-with-network
cp .env.example .env
go run .
```

Closing the client does not terminate the sandbox.

```sh
GOWORK=off go test ./...
```

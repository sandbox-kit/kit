# Create a Modal sandbox with outbound access blocked

Create one sandbox from `alpine:3.21` and block all outbound network access.

## Run

Requires Go 1.26.9 or newer, this checkout, and an existing Modal app.
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

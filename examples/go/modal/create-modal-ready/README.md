# Create a Modal sandbox and wait until ready

Create one sandbox from `alpine:3.21`, start `sleep 300`, and wait until the
command `true` succeeds. The wait uses a two-minute creation timeout.

## Run

Requires Go 1.26.9 or newer, this checkout, and an existing Modal app.
See [Modal creation](../create-modal-sandbox/README.md) for the token and app.

```sh
cd examples/go/modal/create-modal-ready
cp .env.example .env
go run .
```

Closing the client does not terminate the sandbox.

```sh
GOWORK=off go test ./...
```

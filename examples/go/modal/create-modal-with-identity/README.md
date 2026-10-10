# Create a Modal sandbox with a workload identity

Create one sandbox from `alpine:3.21` and request a Modal identity token for
the workload.

## Run

Requires Go 1.26.9 or newer, this checkout, and an existing Modal app.
See [Modal creation](../create-modal-sandbox/README.md) for the token and app.

```sh
cd examples/go/modal/create-modal-with-identity
cp .env.example .env
go run .
```

Closing the client does not terminate the sandbox.

```sh
GOWORK=off go test ./...
```

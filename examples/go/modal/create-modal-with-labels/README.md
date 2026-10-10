# Create a Modal sandbox with environment and tags

Create one sandbox from `alpine:3.21` with `SANDBOX_KIT=example` and the tag
`example=labels`.

## Run

Requires Go 1.26.1 or newer, this checkout, and an existing Modal app.
See [Modal creation](../create-modal-sandbox/README.md) for the token and app.

```sh
cd examples/go/modal/create-modal-with-labels
cp .env.example .env
go run .
```

Closing the client does not terminate the sandbox.

```sh
GOWORK=off go test ./...
```

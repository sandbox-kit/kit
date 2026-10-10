# Create a Modal sandbox with placement

Create one sandbox from `alpine:3.21` with a cloud preference. `MODAL_CLOUD`
defaults to `auto`. Set `MODAL_REGION` to add a region preference.

## Run

Requires Go 1.26.9 or newer, this checkout, and an existing Modal app.
See [Modal creation](../create-modal-sandbox/README.md) for the token and app.

```sh
cd examples/go/modal/create-modal-with-placement
cp .env.example .env
go run .
```

A preference is not a scheduling guarantee. Closing the client does not
terminate the sandbox.

```sh
GOWORK=off go test ./...
```

# Request Modal termination

Attach to an existing Modal sandbox and request termination without waiting for
an exit code. `sandbox.NewClient` selects the Modal provider and scopes the
request to an existing app and environment.

## Run

Requires Go 1.26.9 or newer, this checkout, and an existing Modal sandbox.
See [Modal creation](../create-modal-sandbox/README.md) for token and app setup.

```sh
cd examples/go/modal/request-modal-termination
cp .env.example .env
# Fill in the token, app, environment, and sandbox ID.
go run .
```

The operation works for Modal sandboxes on either the gVisor container or Linux
VM runtime. Attachment does not fetch the runtime class or current state. Use
the [waited termination example](../terminate-modal-sandbox/README.md) when an
exit code is required.

```sh
GOWORK=off go test ./...
```

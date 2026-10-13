# Terminate a Modal sandbox

This example attaches to an existing Modal sandbox ID without claiming a remote
lookup, checks termination support, and waits for its exit code.
`sandbox.NewClient` selects the Modal provider and scopes the request to an
existing app and environment.

## Run

Requires Go 1.26.9 or newer, this checkout, and an existing Modal sandbox.
See [Modal creation](../create-modal-sandbox/README.md) for token and app setup.

```sh
cd examples/go/modal/terminate-modal-sandbox
cp .env.example .env
# Fill in the token, app, environment, and sandbox ID.
go run .
```

Termination works for both Modal runtime classes: gVisor containers and Linux
VMs. Attachment does not retrieve the runtime class or current state. A
successful waited termination marks this local handle terminal and returns the
exit code; `client.Close` only releases local SDK resources.

```sh
GOWORK=off go test ./...
```

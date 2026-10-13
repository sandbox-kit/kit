# Attach to a Modal sandbox

Construct a Sandbox Kit handle from an existing Modal sandbox ID and inspect
its capabilities. `sandbox.NewClient` selects the Modal provider and scopes the
handle to an existing Modal app and environment.

## Run

Requires Go 1.26.9 or newer, this checkout, an existing Modal app, and a sandbox
ID from that app. See [Modal creation](../create-modal-sandbox/README.md) for
token and app setup.

```sh
cd examples/go/modal/attach-modal-sandbox
cp .env.example .env
# Fill in the token, app, environment, and sandbox ID.
go run .
```

Attachment performs no remote lookup, so it reports no current state or
resolved isolation class. Modal's attached handle supports termination for both
the gVisor container and Linux VM runtimes. Closing the client does not terminate
the remote sandbox.

```sh
GOWORK=off go test ./...
```

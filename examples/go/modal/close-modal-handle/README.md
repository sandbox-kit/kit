# Close a local Modal sandbox handle

Attach a Sandbox Kit handle to an existing Modal sandbox, then release that
handle with `Sandbox.Close`. This is local cleanup: it does not terminate or
otherwise change the remote sandbox.

## Run

Requires Go 1.26.9 or newer, this checkout, an existing Modal app, and a sandbox
ID from that app. See [Modal creation](../create-modal-sandbox/README.md) for
token and app setup.

```sh
cd examples/go/modal/close-modal-handle
cp .env.example .env
# Fill in the token, app, environment, and sandbox ID.
go run .
```

`sandbox.NewClient` selects Modal and supplies the app and environment required
to attach. Attachment does not retrieve the sandbox's isolation class, and
`Sandbox.Close` behaves the same for gVisor containers and Linux VMs. After the
close, the local handle exposes no lifecycle capabilities. Use `Terminate` when
you intend to stop the remote sandbox.

```sh
GOWORK=off go test ./...
```

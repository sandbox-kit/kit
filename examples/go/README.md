# Go examples

Independent modules with local replacements for the public SDK and one provider.
Requires Go 1.26.9 or newer and this checkout.

## Run one

```sh
cd examples/go/daytona/create-daytona-sandbox
cp .env.example .env
# Fill in the values, then:
go run .
```

Skip copying if `.env` already exists. Programs load `.env` with `godotenv`;
shell variables win. The SDK does not load files. Modal app setup:
[create-modal-sandbox](modal/create-modal-sandbox/README.md).

Creation programs provision real sandboxes (linked Daytona creates a parent and
child). Lifecycle programs use the configured sandbox ID—read the guide before
destructive stop, terminate, archive, or delete. Error-handling programs fail
locally. Closing a client or handle is local-only.

## Catalog

| Project                                                                        | Required configuration                                | Local verification                              |
| ------------------------------------------------------------------------------ | ----------------------------------------------------- | ----------------------------------------------- |
| [create-modal-sandbox](modal/create-modal-sandbox/README.md)                   | Modal token pair, existing app, optional environment  | Error-kind handling tests                       |
| [create-modal-with-resources](modal/create-modal-with-resources/README.md)     | Modal token pair, existing app, optional environment  | Validation and scope tests                      |
| [create-modal-with-runtime](modal/create-modal-with-runtime/README.md)         | Modal token pair, existing app, optional environment  | Compiles                                        |
| [create-modal-gvisor](modal/create-modal-gvisor/README.md)                     | Modal token pair, existing app, optional environment  | Compiles                                        |
| [create-modal-linux-vm](modal/create-modal-linux-vm/README.md)                 | Modal token pair, existing app, optional environment  | Compiles                                        |
| [create-modal-with-placement](modal/create-modal-with-placement/README.md)     | Modal token pair, app, optional cloud and region      | Compiles                                        |
| [create-modal-with-labels](modal/create-modal-with-labels/README.md)           | Modal token pair, existing app, optional environment  | Compiles                                        |
| [create-modal-with-network](modal/create-modal-with-network/README.md)         | Modal token pair, existing app and proxy              | Compiles                                        |
| [create-modal-idle](modal/create-modal-idle/README.md)                         | Modal token pair, existing app, optional environment  | Compiles                                        |
| [create-modal-ready](modal/create-modal-ready/README.md)                       | Modal token pair, existing app, optional environment  | Compiles                                        |
| [create-modal-verbose](modal/create-modal-verbose/README.md)                   | Modal token pair, existing app, optional environment  | Compiles                                        |
| [create-modal-with-identity](modal/create-modal-with-identity/README.md)       | Modal token pair, existing app, optional environment  | Compiles                                        |
| [create-daytona-sandbox](daytona/create-daytona-sandbox/README.md)             | Daytona API key, optional endpoint/target             | Error-kind handling tests                       |
| [create-daytona-from-snapshot](daytona/create-daytona-from-snapshot/README.md) | Daytona API key, optional endpoint/target             | Compiles                                        |
| [create-daytona-from-image](daytona/create-daytona-from-image/README.md)       | Daytona API key, optional endpoint/target             | Compiles                                        |
| [create-daytona-with-runtime](daytona/create-daytona-with-runtime/README.md)   | Daytona API key, optional endpoint/target             | Compiles                                        |
| [create-daytona-ephemeral](daytona/create-daytona-ephemeral/README.md)         | Daytona API key, optional endpoint/target             | Compiles                                        |
| [create-daytona-linked](daytona/create-daytona-linked/README.md)               | Daytona API key, optional endpoint/target             | Compiles                                        |
| [create-daytona-with-policies](daytona/create-daytona-with-policies/README.md) | Daytona API key, optional endpoint/target             | SDK request and policy tests                    |
| [handle-modal-errors](modal/handle-modal-errors/README.md)                     | Modal token pair and app name                         | Unsupported disk override and kind-switch tests |
| [handle-daytona-errors](daytona/handle-daytona-errors/README.md)               | Daytona API key                                       | Local validation and kind-switch tests          |
| [get-daytona-sandbox](daytona/get-daytona-sandbox/README.md)                   | Daytona API key, sandbox ID, optional endpoint/target | Compiles                                        |
| [refresh-daytona-sandbox](daytona/refresh-daytona-sandbox/README.md)           | Daytona API key, sandbox ID, optional endpoint/target | Compiles                                        |
| [start-daytona-sandbox](daytona/start-daytona-sandbox/README.md)               | Daytona API key, sandbox ID, optional endpoint/target | Compiles                                        |
| [stop-daytona-sandbox](daytona/stop-daytona-sandbox/README.md)                 | Daytona API key, sandbox ID, optional endpoint/target | Compiles                                        |
| [pause-daytona-sandbox](daytona/pause-daytona-sandbox/README.md)               | Daytona API key, sandbox ID, optional endpoint/target | Compiles                                        |
| [archive-daytona-sandbox](daytona/archive-daytona-sandbox/README.md)           | Daytona API key, sandbox ID, optional endpoint/target | Compiles                                        |
| [delete-daytona-sandbox](daytona/delete-daytona-sandbox/README.md)             | Daytona API key, sandbox ID, optional endpoint/target | Compiles                                        |
| [close-daytona-handle](daytona/close-daytona-handle/README.md)                 | Daytona API key, sandbox ID, optional endpoint/target | Compiles                                        |
| [attach-modal-sandbox](modal/attach-modal-sandbox/README.md)                   | Modal token pair, app, environment, sandbox ID        | Compiles                                        |
| [terminate-modal-sandbox](modal/terminate-modal-sandbox/README.md)             | Modal token pair, app, environment, sandbox ID        | Compiles                                        |
| [request-modal-termination](modal/request-modal-termination/README.md)         | Modal token pair, app, environment, sandbox ID        | Compiles                                        |
| [close-modal-handle](modal/close-modal-handle/README.md)                       | Modal token pair, app, environment, sandbox ID        | Compiles                                        |

## Verify without credentials

```sh
GOWORK=off go test -v ./...
```

Creation and error examples include handler tests. Resources and policies also
use dummy credentials and local servers (Modal: numeric rejection and scope;
Daytona: serialized policies and rejection before creation). Local tests may
bind a port; they do not schedule real provider work.

Repository-wide: `sandbox-kit test go` from `tooling/`.

## Error handling

Creation and error examples keep provider setup in `main.go` and a small handler
in `errors.go` that switches on `*sandbox.Error` `Kind`. Other create projects
print and exit. Cleanup failures do not replace creation failures.

Modal errors: unsupported disk. Daytona errors: invalid CPU.

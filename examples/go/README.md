# Go examples

Each project uses the sandbox client and one provider from this repository.
Projects sit in `filesystem`, `process`, `configuration`, `lifecycle`, or
`errors`. Requires Go 1.26.0 or newer and this checkout.

## Run one

```sh
cd examples/go/daytona/configuration/create-daytona-sandbox
cp .env.example .env
# Fill in the values, then:
go run .
```

Skip copying if `.env` already exists. Programs load `.env` with `godotenv`;
shell variables win. The SDK does not load files. Modal app setup:
[create-modal-sandbox](modal/configuration/create-modal-sandbox/README.md).

Creation programs provision real sandboxes (linked Daytona creates a parent and
child). Filesystem, process, and lifecycle programs use the configured sandbox
ID. Read the guide before a destructive stop, terminate, archive, delete, or
remove.
Error-handling programs fail locally. Closing a handle leaves the sandbox as it is. Closing the client
closes your connection to Daytona or Modal.

## Catalog

### Daytona

| Project                                                                                      | Demonstrates                            |
| -------------------------------------------------------------------------------------------- | --------------------------------------- |
| [create-daytona-sandbox](daytona/configuration/create-daytona-sandbox/README.md)             | default Daytona snapshot                |
| [run-daytona-command](daytona/process/run-daytona-command/README.md)                         | run one command                         |
| [run-daytona-shell-command](daytona/process/run-daytona-shell-command/README.md)             | run one shell string                    |
| [run-daytona-with-cwd](daytona/process/run-daytona-with-cwd/README.md)                       | run with working directory              |
| [run-daytona-with-env](daytona/process/run-daytona-with-env/README.md)                       | run with environment variables          |
| [run-daytona-check-exit](daytona/process/run-daytona-check-exit/README.md)                   | non-zero exit in the result             |
| [read-daytona-file](daytona/filesystem/read-daytona-file/README.md)                          | read one file                           |
| [write-daytona-file](daytona/filesystem/write-daytona-file/README.md)                        | write one file                          |
| [copy-file-into-daytona](daytona/filesystem/copy-file-into-daytona/README.md)                | copy a local file in                    |
| [copy-file-from-daytona](daytona/filesystem/copy-file-from-daytona/README.md)                | copy a sandbox file out                 |
| [make-daytona-directory](daytona/filesystem/make-daytona-directory/README.md)                | create a directory                      |
| [list-daytona-directory](daytona/filesystem/list-daytona-directory/README.md)                | list one directory                      |
| [stat-daytona-file](daytona/filesystem/stat-daytona-file/README.md)                          | stat one path                           |
| [remove-daytona-file](daytona/filesystem/remove-daytona-file/README.md)                      | remove one path                         |
| [create-daytona-from-snapshot](daytona/configuration/create-daytona-from-snapshot/README.md) | named snapshot source                   |
| [create-daytona-from-image](daytona/configuration/create-daytona-from-image/README.md)       | image source and resource sizing        |
| [create-daytona-with-runtime](daytona/configuration/create-daytona-with-runtime/README.md)   | runtime language and user               |
| [create-daytona-ephemeral](daytona/configuration/create-daytona-ephemeral/README.md)         | delete-on-stop behavior                 |
| [create-daytona-linked](daytona/configuration/create-daytona-linked/README.md)               | linked parent and ephemeral child       |
| [create-daytona-with-policies](daytona/configuration/create-daytona-with-policies/README.md) | pause and deletion policies             |
| [create-daytona-with-network](daytona/configuration/create-daytona-with-network/README.md)   | restricted egress allowlists            |
| [create-daytona-with-labels](daytona/configuration/create-daytona-with-labels/README.md)     | labels and environment variables        |
| [create-daytona-with-secrets](daytona/configuration/create-daytona-with-secrets/README.md)   | injected organization secret            |
| [create-daytona-with-kvm](daytona/configuration/create-daytona-with-kvm/README.md)           | nested virtualization (KVM)             |
| [create-daytona-with-ttl](daytona/configuration/create-daytona-with-ttl/README.md)           | maximum lifetime TTL                    |
| [handle-daytona-errors](daytona/errors/handle-daytona-errors/README.md)                      | typed validation errors                 |
| [get-daytona-sandbox](daytona/lifecycle/get-daytona-sandbox/README.md)                       | remote retrieval and metadata           |
| [refresh-daytona-sandbox](daytona/lifecycle/refresh-daytona-sandbox/README.md)               | metadata refresh                        |
| [start-daytona-sandbox](daytona/lifecycle/start-daytona-sandbox/README.md)                   | start lifecycle operation               |
| [start-daytona-with-timeout](daytona/lifecycle/start-daytona-with-timeout/README.md)         | start with operation timeout            |
| [stop-daytona-sandbox](daytona/lifecycle/stop-daytona-sandbox/README.md)                     | stop lifecycle operation                |
| [force-stop-daytona-sandbox](daytona/lifecycle/force-stop-daytona-sandbox/README.md)         | force stop                              |
| [pause-daytona-sandbox](daytona/lifecycle/pause-daytona-sandbox/README.md)                   | pause lifecycle operation               |
| [archive-daytona-sandbox](daytona/lifecycle/archive-daytona-sandbox/README.md)               | archive lifecycle operation             |
| [delete-daytona-sandbox](daytona/lifecycle/delete-daytona-sandbox/README.md)                 | waited deletion                         |
| [request-daytona-deletion](daytona/lifecycle/request-daytona-deletion/README.md)             | non-waited deletion                     |
| [close-daytona-handle](daytona/lifecycle/close-daytona-handle/README.md)                     | close the handle; sandbox keeps running |

### Modal

| Project                                                                                    | Demonstrates                               |
| ------------------------------------------------------------------------------------------ | ------------------------------------------ |
| [create-modal-sandbox](modal/configuration/create-modal-sandbox/README.md)                 | image creation with resources and policies |
| [run-modal-command](modal/process/run-modal-command/README.md)                             | run one command                            |
| [run-modal-shell-command](modal/process/run-modal-shell-command/README.md)                 | run one shell string                       |
| [run-modal-with-cwd](modal/process/run-modal-with-cwd/README.md)                           | run with working directory                 |
| [run-modal-with-env](modal/process/run-modal-with-env/README.md)                           | run with environment variables             |
| [run-modal-check-exit](modal/process/run-modal-check-exit/README.md)                       | non-zero exit in the result                |
| [read-modal-file](modal/filesystem/read-modal-file/README.md)                              | read one file                              |
| [write-modal-file](modal/filesystem/write-modal-file/README.md)                            | write one file                             |
| [copy-file-into-modal](modal/filesystem/copy-file-into-modal/README.md)                    | copy a local file in                       |
| [copy-file-from-modal](modal/filesystem/copy-file-from-modal/README.md)                    | copy a sandbox file out                    |
| [make-modal-directory](modal/filesystem/make-modal-directory/README.md)                    | create a directory                         |
| [list-modal-directory](modal/filesystem/list-modal-directory/README.md)                    | list one directory                         |
| [stat-modal-file](modal/filesystem/stat-modal-file/README.md)                              | stat one path                              |
| [remove-modal-file](modal/filesystem/remove-modal-file/README.md)                          | remove one path                            |
| [create-modal-with-resources](modal/configuration/create-modal-with-resources/README.md)   | fractional CPU, memory, and lifetime       |
| [create-modal-with-runtime](modal/configuration/create-modal-with-runtime/README.md)       | working directory, entrypoint, and PTY     |
| [create-modal-gvisor](modal/configuration/create-modal-gvisor/README.md)                   | gVisor/container isolation                 |
| [create-modal-linux-vm](modal/configuration/create-modal-linux-vm/README.md)               | Linux VM isolation                         |
| [create-modal-with-placement](modal/configuration/create-modal-with-placement/README.md)   | cloud and region placement                 |
| [create-modal-with-labels](modal/configuration/create-modal-with-labels/README.md)         | environment variables and tags             |
| [create-modal-with-network](modal/configuration/create-modal-with-network/README.md)       | managed outbound proxy                     |
| [create-modal-idle](modal/configuration/create-modal-idle/README.md)                       | idle termination                           |
| [create-modal-ready](modal/configuration/create-modal-ready/README.md)                     | command readiness                          |
| [create-modal-tcp-ready](modal/configuration/create-modal-tcp-ready/README.md)             | TCP port readiness                         |
| [create-modal-client-timeout](modal/configuration/create-modal-client-timeout/README.md)   | ready-wait via client timeout              |
| [create-modal-with-ports](modal/configuration/create-modal-with-ports/README.md)           | TLS and plain port bindings                |
| [create-modal-with-allowlists](modal/configuration/create-modal-with-allowlists/README.md) | outbound CIDR/domain allowlists            |
| [create-modal-with-secrets](modal/configuration/create-modal-with-secrets/README.md)       | named Modal secret                         |
| [create-modal-verbose](modal/configuration/create-modal-verbose/README.md)                 | verbose provider logging                   |
| [create-modal-with-identity](modal/configuration/create-modal-with-identity/README.md)     | workload identity                          |
| [handle-modal-errors](modal/errors/handle-modal-errors/README.md)                          | typed unsupported-intent errors            |
| [attach-modal-sandbox](modal/lifecycle/attach-modal-sandbox/README.md)                     | ID-only attachment                         |
| [terminate-modal-sandbox](modal/lifecycle/terminate-modal-sandbox/README.md)               | waited termination                         |
| [terminate-modal-with-timeout](modal/lifecycle/terminate-modal-with-timeout/README.md)     | waited terminate with timeout              |
| [request-modal-termination](modal/lifecycle/request-modal-termination/README.md)           | non-waited termination                     |
| [close-modal-handle](modal/lifecycle/close-modal-handle/README.md)                         | close the handle; sandbox keeps running    |

## Verify without credentials

```sh
GOWORK=off go test -v ./...
```

Most examples compile only. Local tests live in `handle-modal-errors`,
`handle-daytona-errors`, `create-modal-with-resources`, and
`create-daytona-with-policies` (dummy credentials and local servers; they do not
schedule real provider work). Local tests may bind a port.

Repository-wide: `sandbox-kit test go` from `tooling/`.

## Error handling

[handle-modal-errors](modal/errors/handle-modal-errors/README.md) and
[handle-daytona-errors](daytona/errors/handle-daytona-errors/README.md) switch on
`*sandbox.Error` `Kind` in `main.go`. Other create projects print and exit.
Cleanup failures do not replace creation failures.

Modal errors: unsupported disk. Daytona errors: invalid CPU.

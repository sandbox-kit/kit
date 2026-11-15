# Examples

Independent Go and TypeScript projects. Each selects one provider, loads its own
`.env`, and demonstrates one operation. Projects for a provider sit in
`filesystem`, `process`, `configuration`, `lifecycle`, or `errors`. Every
project README explains what the example does for that provider, with edge
cases drawn from Daytona, Modal, and Blaxel behavior.

## Run one

```sh
# Go — requires Go 1.26.0+
cd examples/go/daytona/configuration/create-daytona-sandbox
cp .env.example .env   # set DAYTONA_API_KEY
go run .

# TypeScript — requires Node.js 20/22/24.15+ and npm 11.19.0
cd examples/typescript/daytona/configuration/create-daytona-sandbox
cp .env.example .env   # set DAYTONA_API_KEY
npm install && npm run build && npm start
```

Modal needs a token pair and an existing app:
[Go](go/modal/configuration/create-modal-sandbox/README.md) ·
[TypeScript](typescript/modal/configuration/create-modal-sandbox/README.md).

Blaxel needs an API key, a workspace, and a client region (examples use
`BL_REGION`, default `us-pdx-1`):
[Go](go/blaxel/configuration/create-blaxel-sandbox/README.md) ·
[TypeScript](typescript/blaxel/configuration/create-blaxel-sandbox/README.md).

Creation examples provision real resources. Filesystem, process, and lifecycle
examples use the sandbox id in `.env` (`DAYTONA_SANDBOX_ID`, `MODAL_SANDBOX_ID`,
or Blaxel `BL_SANDBOX_NAME` — Blaxel id is the create name). A filesystem
example runs one file operation (read, write, copy, mkdir, list, stat, or
remove) on the path you set. A process example runs one command. Stop,
terminate, archive, and delete can destroy the sandbox.
Error-handling examples fail locally. `Client.Close` closes your connection to
Daytona, Modal, or Blaxel. `Sandbox.Close` does not change the sandbox, and you can keep
using the handle.

Indexes: [Go](go/README.md) · [TypeScript](typescript/README.md).

## Catalog

| Example                   | Go                                                                    | TypeScript                                                                            |
| ------------------------- | --------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| Daytona creation          | [go](go/daytona/configuration/create-daytona-sandbox/README.md)       | [typescript](typescript/daytona/configuration/create-daytona-sandbox/README.md)       |
| Daytona run command       | [go](go/daytona/process/run-daytona-command/README.md)                | [typescript](typescript/daytona/process/run-daytona-command/README.md)                |
| Daytona shell command     | [go](go/daytona/process/run-daytona-shell-command/README.md)          | [typescript](typescript/daytona/process/run-daytona-shell-command/README.md)          |
| Daytona run with cwd      | [go](go/daytona/process/run-daytona-with-cwd/README.md)               | [typescript](typescript/daytona/process/run-daytona-with-cwd/README.md)               |
| Daytona run with env      | [go](go/daytona/process/run-daytona-with-env/README.md)               | [typescript](typescript/daytona/process/run-daytona-with-env/README.md)               |
| Daytona run with timeout  | [go](go/daytona/process/run-daytona-with-timeout/README.md)           | [typescript](typescript/daytona/process/run-daytona-with-timeout/README.md)           |
| Daytona check exit        | [go](go/daytona/process/run-daytona-check-exit/README.md)             | [typescript](typescript/daytona/process/run-daytona-check-exit/README.md)             |
| Daytona read file         | [go](go/daytona/filesystem/read-daytona-file/README.md)               | [typescript](typescript/daytona/filesystem/read-daytona-file/README.md)               |
| Daytona write file        | [go](go/daytona/filesystem/write-daytona-file/README.md)              | [typescript](typescript/daytona/filesystem/write-daytona-file/README.md)              |
| Daytona copy in           | [go](go/daytona/filesystem/copy-file-into-daytona/README.md)          | [typescript](typescript/daytona/filesystem/copy-file-into-daytona/README.md)          |
| Daytona copy out          | [go](go/daytona/filesystem/copy-file-from-daytona/README.md)          | [typescript](typescript/daytona/filesystem/copy-file-from-daytona/README.md)          |
| Daytona directory         | [go](go/daytona/filesystem/make-daytona-directory/README.md)          | [typescript](typescript/daytona/filesystem/make-daytona-directory/README.md)          |
| Daytona list              | [go](go/daytona/filesystem/list-daytona-directory/README.md)          | [typescript](typescript/daytona/filesystem/list-daytona-directory/README.md)          |
| Daytona stat              | [go](go/daytona/filesystem/stat-daytona-file/README.md)               | [typescript](typescript/daytona/filesystem/stat-daytona-file/README.md)               |
| Daytona remove            | [go](go/daytona/filesystem/remove-daytona-file/README.md)             | [typescript](typescript/daytona/filesystem/remove-daytona-file/README.md)             |
| Daytona snapshot          | [go](go/daytona/configuration/create-daytona-from-snapshot/README.md) | [typescript](typescript/daytona/configuration/create-daytona-from-snapshot/README.md) |
| Daytona image             | [go](go/daytona/configuration/create-daytona-from-image/README.md)    | [typescript](typescript/daytona/configuration/create-daytona-from-image/README.md)    |
| Daytona runtime           | [go](go/daytona/configuration/create-daytona-with-runtime/README.md)  | [typescript](typescript/daytona/configuration/create-daytona-with-runtime/README.md)  |
| Daytona ephemeral         | [go](go/daytona/configuration/create-daytona-ephemeral/README.md)     | [typescript](typescript/daytona/configuration/create-daytona-ephemeral/README.md)     |
| Daytona linked            | [go](go/daytona/configuration/create-daytona-linked/README.md)        | [typescript](typescript/daytona/configuration/create-daytona-linked/README.md)        |
| Daytona policies          | [go](go/daytona/configuration/create-daytona-with-policies/README.md) | [typescript](typescript/daytona/configuration/create-daytona-with-policies/README.md) |
| Daytona network           | [go](go/daytona/configuration/create-daytona-with-network/README.md)  | [typescript](typescript/daytona/configuration/create-daytona-with-network/README.md)  |
| Daytona labels            | [go](go/daytona/configuration/create-daytona-with-labels/README.md)   | [typescript](typescript/daytona/configuration/create-daytona-with-labels/README.md)   |
| Daytona secrets           | [go](go/daytona/configuration/create-daytona-with-secrets/README.md)  | [typescript](typescript/daytona/configuration/create-daytona-with-secrets/README.md)  |
| Daytona KVM               | [go](go/daytona/configuration/create-daytona-with-kvm/README.md)      | [typescript](typescript/daytona/configuration/create-daytona-with-kvm/README.md)      |
| Daytona TTL               | [go](go/daytona/configuration/create-daytona-with-ttl/README.md)      | [typescript](typescript/daytona/configuration/create-daytona-with-ttl/README.md)      |
| Daytona get               | [go](go/daytona/lifecycle/get-daytona-sandbox/README.md)              | [typescript](typescript/daytona/lifecycle/get-daytona-sandbox/README.md)              |
| Daytona refresh           | [go](go/daytona/lifecycle/refresh-daytona-sandbox/README.md)          | [typescript](typescript/daytona/lifecycle/refresh-daytona-sandbox/README.md)          |
| Daytona start             | [go](go/daytona/lifecycle/start-daytona-sandbox/README.md)            | [typescript](typescript/daytona/lifecycle/start-daytona-sandbox/README.md)            |
| Daytona start timeout     | [go](go/daytona/lifecycle/start-daytona-with-timeout/README.md)       | [typescript](typescript/daytona/lifecycle/start-daytona-with-timeout/README.md)       |
| Daytona stop              | [go](go/daytona/lifecycle/stop-daytona-sandbox/README.md)             | [typescript](typescript/daytona/lifecycle/stop-daytona-sandbox/README.md)             |
| Daytona force stop        | [go](go/daytona/lifecycle/force-stop-daytona-sandbox/README.md)       | [typescript](typescript/daytona/lifecycle/force-stop-daytona-sandbox/README.md)       |
| Daytona pause             | [go](go/daytona/lifecycle/pause-daytona-sandbox/README.md)            | [typescript](typescript/daytona/lifecycle/pause-daytona-sandbox/README.md)            |
| Daytona archive           | [go](go/daytona/lifecycle/archive-daytona-sandbox/README.md)          | [typescript](typescript/daytona/lifecycle/archive-daytona-sandbox/README.md)          |
| Daytona delete            | [go](go/daytona/lifecycle/delete-daytona-sandbox/README.md)           | [typescript](typescript/daytona/lifecycle/delete-daytona-sandbox/README.md)           |
| Daytona delete request    | [go](go/daytona/lifecycle/request-daytona-deletion/README.md)         | [typescript](typescript/daytona/lifecycle/request-daytona-deletion/README.md)         |
| Daytona close handle      | [go](go/daytona/lifecycle/close-daytona-handle/README.md)             | [typescript](typescript/daytona/lifecycle/close-daytona-handle/README.md)             |
| Daytona errors            | [go](go/daytona/errors/handle-daytona-errors/README.md)               | [typescript](typescript/daytona/errors/handle-daytona-errors/README.md)               |
| Modal creation            | [go](go/modal/configuration/create-modal-sandbox/README.md)           | [typescript](typescript/modal/configuration/create-modal-sandbox/README.md)           |
| Modal run command         | [go](go/modal/process/run-modal-command/README.md)                    | [typescript](typescript/modal/process/run-modal-command/README.md)                    |
| Modal shell command       | [go](go/modal/process/run-modal-shell-command/README.md)              | [typescript](typescript/modal/process/run-modal-shell-command/README.md)              |
| Modal run with cwd        | [go](go/modal/process/run-modal-with-cwd/README.md)                   | [typescript](typescript/modal/process/run-modal-with-cwd/README.md)                   |
| Modal run with env        | [go](go/modal/process/run-modal-with-env/README.md)                   | [typescript](typescript/modal/process/run-modal-with-env/README.md)                   |
| Modal run with timeout    | [go](go/modal/process/run-modal-with-timeout/README.md)               | [typescript](typescript/modal/process/run-modal-with-timeout/README.md)               |
| Modal check exit          | [go](go/modal/process/run-modal-check-exit/README.md)                 | [typescript](typescript/modal/process/run-modal-check-exit/README.md)                 |
| Modal read file           | [go](go/modal/filesystem/read-modal-file/README.md)                   | [typescript](typescript/modal/filesystem/read-modal-file/README.md)                   |
| Modal write file          | [go](go/modal/filesystem/write-modal-file/README.md)                  | [typescript](typescript/modal/filesystem/write-modal-file/README.md)                  |
| Modal copy in             | [go](go/modal/filesystem/copy-file-into-modal/README.md)              | [typescript](typescript/modal/filesystem/copy-file-into-modal/README.md)              |
| Modal copy out            | [go](go/modal/filesystem/copy-file-from-modal/README.md)              | [typescript](typescript/modal/filesystem/copy-file-from-modal/README.md)              |
| Modal directory           | [go](go/modal/filesystem/make-modal-directory/README.md)              | [typescript](typescript/modal/filesystem/make-modal-directory/README.md)              |
| Modal list                | [go](go/modal/filesystem/list-modal-directory/README.md)              | [typescript](typescript/modal/filesystem/list-modal-directory/README.md)              |
| Modal stat                | [go](go/modal/filesystem/stat-modal-file/README.md)                   | [typescript](typescript/modal/filesystem/stat-modal-file/README.md)                   |
| Modal remove              | [go](go/modal/filesystem/remove-modal-file/README.md)                 | [typescript](typescript/modal/filesystem/remove-modal-file/README.md)                 |
| Modal resources           | [go](go/modal/configuration/create-modal-with-resources/README.md)    | [typescript](typescript/modal/configuration/create-modal-with-resources/README.md)    |
| Modal runtime             | [go](go/modal/configuration/create-modal-with-runtime/README.md)      | [typescript](typescript/modal/configuration/create-modal-with-runtime/README.md)      |
| Modal gVisor              | [go](go/modal/configuration/create-modal-gvisor/README.md)            | [typescript](typescript/modal/configuration/create-modal-gvisor/README.md)            |
| Modal Linux VM            | [go](go/modal/configuration/create-modal-linux-vm/README.md)          | [typescript](typescript/modal/configuration/create-modal-linux-vm/README.md)          |
| Modal placement           | [go](go/modal/configuration/create-modal-with-placement/README.md)    | [typescript](typescript/modal/configuration/create-modal-with-placement/README.md)    |
| Modal labels              | [go](go/modal/configuration/create-modal-with-labels/README.md)       | [typescript](typescript/modal/configuration/create-modal-with-labels/README.md)       |
| Modal network             | [go](go/modal/configuration/create-modal-with-network/README.md)      | [typescript](typescript/modal/configuration/create-modal-with-network/README.md)      |
| Modal idle timeout        | [go](go/modal/configuration/create-modal-idle/README.md)              | [typescript](typescript/modal/configuration/create-modal-idle/README.md)              |
| Modal readiness           | [go](go/modal/configuration/create-modal-ready/README.md)             | [typescript](typescript/modal/configuration/create-modal-ready/README.md)             |
| Modal TCP readiness       | [go](go/modal/configuration/create-modal-tcp-ready/README.md)         | [typescript](typescript/modal/configuration/create-modal-tcp-ready/README.md)         |
| Modal client timeout wait | [go](go/modal/configuration/create-modal-client-timeout/README.md)    | [typescript](typescript/modal/configuration/create-modal-client-timeout/README.md)    |
| Modal ports               | [go](go/modal/configuration/create-modal-with-ports/README.md)        | [typescript](typescript/modal/configuration/create-modal-with-ports/README.md)        |
| Modal allowlists          | [go](go/modal/configuration/create-modal-with-allowlists/README.md)   | [typescript](typescript/modal/configuration/create-modal-with-allowlists/README.md)   |
| Modal secrets             | [go](go/modal/configuration/create-modal-with-secrets/README.md)      | [typescript](typescript/modal/configuration/create-modal-with-secrets/README.md)      |
| Modal verbose             | [go](go/modal/configuration/create-modal-verbose/README.md)           | [typescript](typescript/modal/configuration/create-modal-verbose/README.md)           |
| Modal identity            | [go](go/modal/configuration/create-modal-with-identity/README.md)     | [typescript](typescript/modal/configuration/create-modal-with-identity/README.md)     |
| Modal attach              | [go](go/modal/lifecycle/attach-modal-sandbox/README.md)               | [typescript](typescript/modal/lifecycle/attach-modal-sandbox/README.md)               |
| Modal terminate           | [go](go/modal/lifecycle/terminate-modal-sandbox/README.md)            | [typescript](typescript/modal/lifecycle/terminate-modal-sandbox/README.md)            |
| Modal terminate timeout   | [go](go/modal/lifecycle/terminate-modal-with-timeout/README.md)       | [typescript](typescript/modal/lifecycle/terminate-modal-with-timeout/README.md)       |
| Modal termination request | [go](go/modal/lifecycle/request-modal-termination/README.md)          | [typescript](typescript/modal/lifecycle/request-modal-termination/README.md)          |
| Modal close handle        | [go](go/modal/lifecycle/close-modal-handle/README.md)                 | [typescript](typescript/modal/lifecycle/close-modal-handle/README.md)                 |
| Modal errors              | [go](go/modal/errors/handle-modal-errors/README.md)                   | [typescript](typescript/modal/errors/handle-modal-errors/README.md)                   |
| Blaxel named image create | [go](go/blaxel/configuration/create-blaxel-sandbox/README.md)         | [typescript](typescript/blaxel/configuration/create-blaxel-sandbox/README.md)         |
| Blaxel memory MiB         | [go](go/blaxel/configuration/create-blaxel-with-resources/README.md)  | [typescript](typescript/blaxel/configuration/create-blaxel-with-resources/README.md)  |
| Blaxel labels and environment | [go](go/blaxel/configuration/create-blaxel-with-labels/README.md)     | [typescript](typescript/blaxel/configuration/create-blaxel-with-labels/README.md)     |
| Blaxel maximum lifetime TTL | [go](go/blaxel/configuration/create-blaxel-with-ttl/README.md)        | [typescript](typescript/blaxel/configuration/create-blaxel-with-ttl/README.md)        |
| Blaxel idle delete policy | [go](go/blaxel/configuration/create-blaxel-idle/README.md)            | [typescript](typescript/blaxel/configuration/create-blaxel-idle/README.md)            |
| Blaxel client timeout     | [go](go/blaxel/configuration/create-blaxel-client-timeout/README.md)  | [typescript](typescript/blaxel/configuration/create-blaxel-client-timeout/README.md)  |
| Blaxel published ports    | [go](go/blaxel/configuration/create-blaxel-with-ports/README.md)      | [typescript](typescript/blaxel/configuration/create-blaxel-with-ports/README.md)      |
| Blaxel outbound domains   | [go](go/blaxel/configuration/create-blaxel-with-allowlists/README.md) | [typescript](typescript/blaxel/configuration/create-blaxel-with-allowlists/README.md) |
| Blaxel run one command    | [go](go/blaxel/process/run-blaxel-command/README.md)                  | [typescript](typescript/blaxel/process/run-blaxel-command/README.md)                  |
| Blaxel run one shell string | [go](go/blaxel/process/run-blaxel-shell-command/README.md)            | [typescript](typescript/blaxel/process/run-blaxel-shell-command/README.md)            |
| Blaxel run with working directory | [go](go/blaxel/process/run-blaxel-with-cwd/README.md)                 | [typescript](typescript/blaxel/process/run-blaxel-with-cwd/README.md)                 |
| Blaxel run with environment variables | [go](go/blaxel/process/run-blaxel-with-env/README.md)                 | [typescript](typescript/blaxel/process/run-blaxel-with-env/README.md)                 |
| Blaxel run with a whole-second timeout | [go](go/blaxel/process/run-blaxel-with-timeout/README.md)             | [typescript](typescript/blaxel/process/run-blaxel-with-timeout/README.md)             |
| Blaxel non-zero exit in the result | [go](go/blaxel/process/run-blaxel-check-exit/README.md)               | [typescript](typescript/blaxel/process/run-blaxel-check-exit/README.md)               |
| Blaxel read one file      | [go](go/blaxel/filesystem/read-blaxel-file/README.md)                 | [typescript](typescript/blaxel/filesystem/read-blaxel-file/README.md)                 |
| Blaxel write one file     | [go](go/blaxel/filesystem/write-blaxel-file/README.md)                | [typescript](typescript/blaxel/filesystem/write-blaxel-file/README.md)                |
| Blaxel copy a local file in | [go](go/blaxel/filesystem/copy-file-into-blaxel/README.md)            | [typescript](typescript/blaxel/filesystem/copy-file-into-blaxel/README.md)            |
| Blaxel copy a sandbox file out | [go](go/blaxel/filesystem/copy-file-from-blaxel/README.md)            | [typescript](typescript/blaxel/filesystem/copy-file-from-blaxel/README.md)            |
| Blaxel create a directory | [go](go/blaxel/filesystem/make-blaxel-directory/README.md)            | [typescript](typescript/blaxel/filesystem/make-blaxel-directory/README.md)            |
| Blaxel list one directory | [go](go/blaxel/filesystem/list-blaxel-directory/README.md)            | [typescript](typescript/blaxel/filesystem/list-blaxel-directory/README.md)            |
| Blaxel stat one path      | [go](go/blaxel/filesystem/stat-blaxel-file/README.md)                 | [typescript](typescript/blaxel/filesystem/stat-blaxel-file/README.md)                 |
| Blaxel remove one path    | [go](go/blaxel/filesystem/remove-blaxel-file/README.md)               | [typescript](typescript/blaxel/filesystem/remove-blaxel-file/README.md)               |
| Blaxel remote retrieval and metadata | [go](go/blaxel/lifecycle/get-blaxel-sandbox/README.md)                | [typescript](typescript/blaxel/lifecycle/get-blaxel-sandbox/README.md)                |
| Blaxel metadata refresh   | [go](go/blaxel/lifecycle/refresh-blaxel-sandbox/README.md)            | [typescript](typescript/blaxel/lifecycle/refresh-blaxel-sandbox/README.md)            |
| Blaxel delete by name     | [go](go/blaxel/lifecycle/delete-blaxel-sandbox/README.md)             | [typescript](typescript/blaxel/lifecycle/delete-blaxel-sandbox/README.md)             |
| Blaxel close the handle; sandbox keeps running | [go](go/blaxel/lifecycle/close-blaxel-handle/README.md)               | [typescript](typescript/blaxel/lifecycle/close-blaxel-handle/README.md)               |
| Blaxel typed validation errors | [go](go/blaxel/errors/handle-blaxel-errors/README.md)                 | [typescript](typescript/blaxel/errors/handle-blaxel-errors/README.md)                 |

## Verify locally

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit test typescript
```

Go verification compiles every example module; resources, policies, and error
projects also run local handler and request tests with dummy credentials.
TypeScript verification builds the SDK packages and every example `package.json`.
Neither command creates cloud resources.

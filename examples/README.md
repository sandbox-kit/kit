# Examples

Independent Go and TypeScript projects. Each selects one provider, loads its own
`.env`, and demonstrates one creation, lifecycle, cleanup, or error-handling
operation. Every project README explains what the example does for that provider,
with edge cases drawn from Daytona and Modal behavior.

## Run one

```sh
# Go — requires Go 1.26.0+
cd examples/go/daytona/create-daytona-sandbox
cp .env.example .env   # set DAYTONA_API_KEY
go run .

# TypeScript — requires Node.js 20/22/24.15+ and npm 11.19.0
cd examples/typescript/daytona/create-daytona-sandbox
cp .env.example .env   # set DAYTONA_API_KEY
npm install && npm run build && npm start
```

Modal needs a token pair and an existing app:
[Go](go/modal/create-modal-sandbox/README.md) ·
[TypeScript](typescript/modal/create-modal-sandbox/README.md).

Creation examples provision real resources. Lifecycle examples use the sandbox
ID in `.env` (stop / terminate / archive / delete can destroy it).
Error-handling examples fail locally. `Client.Close` closes the provider SDK.
`Sandbox.Close` calls the provider handle close and does not change the cloud
sandbox. Daytona returns without a provider request. Modal detaches the local
sandbox object.

Indexes: [Go](go/README.md) · [TypeScript](typescript/README.md).

## Catalog

| Example                   | Go                                                      | TypeScript                                                              |
| ------------------------- | ------------------------------------------------------- | ----------------------------------------------------------------------- |
| Daytona creation          | [go](go/daytona/create-daytona-sandbox/README.md)       | [typescript](typescript/daytona/create-daytona-sandbox/README.md)       |
| Daytona snapshot          | [go](go/daytona/create-daytona-from-snapshot/README.md) | [typescript](typescript/daytona/create-daytona-from-snapshot/README.md) |
| Daytona image             | [go](go/daytona/create-daytona-from-image/README.md)    | [typescript](typescript/daytona/create-daytona-from-image/README.md)    |
| Daytona runtime           | [go](go/daytona/create-daytona-with-runtime/README.md)  | [typescript](typescript/daytona/create-daytona-with-runtime/README.md)  |
| Daytona ephemeral         | [go](go/daytona/create-daytona-ephemeral/README.md)     | [typescript](typescript/daytona/create-daytona-ephemeral/README.md)     |
| Daytona linked            | [go](go/daytona/create-daytona-linked/README.md)        | [typescript](typescript/daytona/create-daytona-linked/README.md)        |
| Daytona policies          | [go](go/daytona/create-daytona-with-policies/README.md) | [typescript](typescript/daytona/create-daytona-with-policies/README.md) |
| Daytona network           | [go](go/daytona/create-daytona-with-network/README.md)  | [typescript](typescript/daytona/create-daytona-with-network/README.md)  |
| Daytona labels            | [go](go/daytona/create-daytona-with-labels/README.md)   | [typescript](typescript/daytona/create-daytona-with-labels/README.md)   |
| Daytona secrets           | [go](go/daytona/create-daytona-with-secrets/README.md)  | [typescript](typescript/daytona/create-daytona-with-secrets/README.md)  |
| Daytona KVM               | [go](go/daytona/create-daytona-with-kvm/README.md)      | [typescript](typescript/daytona/create-daytona-with-kvm/README.md)      |
| Daytona TTL               | [go](go/daytona/create-daytona-with-ttl/README.md)      | [typescript](typescript/daytona/create-daytona-with-ttl/README.md)      |
| Daytona get               | [go](go/daytona/get-daytona-sandbox/README.md)          | [typescript](typescript/daytona/get-daytona-sandbox/README.md)          |
| Daytona refresh           | [go](go/daytona/refresh-daytona-sandbox/README.md)      | [typescript](typescript/daytona/refresh-daytona-sandbox/README.md)      |
| Daytona start             | [go](go/daytona/start-daytona-sandbox/README.md)        | [typescript](typescript/daytona/start-daytona-sandbox/README.md)        |
| Daytona start timeout     | [go](go/daytona/start-daytona-with-timeout/README.md)   | [typescript](typescript/daytona/start-daytona-with-timeout/README.md)   |
| Daytona stop              | [go](go/daytona/stop-daytona-sandbox/README.md)         | [typescript](typescript/daytona/stop-daytona-sandbox/README.md)         |
| Daytona force stop        | [go](go/daytona/force-stop-daytona-sandbox/README.md)   | [typescript](typescript/daytona/force-stop-daytona-sandbox/README.md)   |
| Daytona pause             | [go](go/daytona/pause-daytona-sandbox/README.md)        | [typescript](typescript/daytona/pause-daytona-sandbox/README.md)        |
| Daytona archive           | [go](go/daytona/archive-daytona-sandbox/README.md)      | [typescript](typescript/daytona/archive-daytona-sandbox/README.md)      |
| Daytona delete            | [go](go/daytona/delete-daytona-sandbox/README.md)       | [typescript](typescript/daytona/delete-daytona-sandbox/README.md)       |
| Daytona delete request    | [go](go/daytona/request-daytona-deletion/README.md)     | [typescript](typescript/daytona/request-daytona-deletion/README.md)     |
| Daytona close handle      | [go](go/daytona/close-daytona-handle/README.md)         | [typescript](typescript/daytona/close-daytona-handle/README.md)         |
| Daytona errors            | [go](go/daytona/handle-daytona-errors/README.md)        | [typescript](typescript/daytona/handle-daytona-errors/README.md)        |
| Modal creation            | [go](go/modal/create-modal-sandbox/README.md)           | [typescript](typescript/modal/create-modal-sandbox/README.md)           |
| Modal resources           | [go](go/modal/create-modal-with-resources/README.md)    | [typescript](typescript/modal/create-modal-with-resources/README.md)    |
| Modal runtime             | [go](go/modal/create-modal-with-runtime/README.md)      | [typescript](typescript/modal/create-modal-with-runtime/README.md)      |
| Modal gVisor              | [go](go/modal/create-modal-gvisor/README.md)            | [typescript](typescript/modal/create-modal-gvisor/README.md)            |
| Modal Linux VM            | [go](go/modal/create-modal-linux-vm/README.md)          | [typescript](typescript/modal/create-modal-linux-vm/README.md)          |
| Modal placement           | [go](go/modal/create-modal-with-placement/README.md)    | [typescript](typescript/modal/create-modal-with-placement/README.md)    |
| Modal labels              | [go](go/modal/create-modal-with-labels/README.md)       | [typescript](typescript/modal/create-modal-with-labels/README.md)       |
| Modal network             | [go](go/modal/create-modal-with-network/README.md)      | [typescript](typescript/modal/create-modal-with-network/README.md)      |
| Modal idle timeout        | [go](go/modal/create-modal-idle/README.md)              | [typescript](typescript/modal/create-modal-idle/README.md)              |
| Modal readiness           | [go](go/modal/create-modal-ready/README.md)             | [typescript](typescript/modal/create-modal-ready/README.md)             |
| Modal TCP readiness       | [go](go/modal/create-modal-tcp-ready/README.md)         | [typescript](typescript/modal/create-modal-tcp-ready/README.md)         |
| Modal client timeout wait | [go](go/modal/create-modal-client-timeout/README.md)    | [typescript](typescript/modal/create-modal-client-timeout/README.md)    |
| Modal ports               | [go](go/modal/create-modal-with-ports/README.md)        | [typescript](typescript/modal/create-modal-with-ports/README.md)        |
| Modal allowlists          | [go](go/modal/create-modal-with-allowlists/README.md)   | [typescript](typescript/modal/create-modal-with-allowlists/README.md)   |
| Modal secrets             | [go](go/modal/create-modal-with-secrets/README.md)      | [typescript](typescript/modal/create-modal-with-secrets/README.md)      |
| Modal verbose             | [go](go/modal/create-modal-verbose/README.md)           | [typescript](typescript/modal/create-modal-verbose/README.md)           |
| Modal identity            | [go](go/modal/create-modal-with-identity/README.md)     | [typescript](typescript/modal/create-modal-with-identity/README.md)     |
| Modal attach              | [go](go/modal/attach-modal-sandbox/README.md)           | [typescript](typescript/modal/attach-modal-sandbox/README.md)           |
| Modal terminate           | [go](go/modal/terminate-modal-sandbox/README.md)        | [typescript](typescript/modal/terminate-modal-sandbox/README.md)        |
| Modal terminate timeout   | [go](go/modal/terminate-modal-with-timeout/README.md)   | [typescript](typescript/modal/terminate-modal-with-timeout/README.md)   |
| Modal termination request | [go](go/modal/request-modal-termination/README.md)      | [typescript](typescript/modal/request-modal-termination/README.md)      |
| Modal close handle        | [go](go/modal/close-modal-handle/README.md)             | [typescript](typescript/modal/close-modal-handle/README.md)             |
| Modal errors              | [go](go/modal/handle-modal-errors/README.md)            | [typescript](typescript/modal/handle-modal-errors/README.md)            |

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

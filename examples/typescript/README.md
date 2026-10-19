# TypeScript examples

Independent Node.js projects. Each selects one provider and demonstrates one
operation. Requires Node.js 24.15.0 (or Node.js 20/22), npm 11.19.0, and this
checkout.

## Run one

```sh
cd examples/typescript/daytona/create-daytona-sandbox
cp .env.example .env
# Fill in DAYTONA_API_KEY, then:
npm install
npm run build
npm start
```

Creation examples provision real cloud resources. Lifecycle examples act on the
sandbox ID in `.env`. Closing a client or handle is local-only. Modal setup:
[create-modal-sandbox](modal/create-modal-sandbox/README.md).

## Catalog

### Daytona

| Project                                                                        | Demonstrates                      |
| ------------------------------------------------------------------------------ | --------------------------------- |
| [create-daytona-sandbox](daytona/create-daytona-sandbox/README.md)             | default Daytona snapshot          |
| [create-daytona-from-snapshot](daytona/create-daytona-from-snapshot/README.md) | named snapshot source             |
| [create-daytona-from-image](daytona/create-daytona-from-image/README.md)       | image source and resource sizing  |
| [create-daytona-with-runtime](daytona/create-daytona-with-runtime/README.md)   | runtime language and user         |
| [create-daytona-ephemeral](daytona/create-daytona-ephemeral/README.md)         | delete-on-stop behavior           |
| [create-daytona-linked](daytona/create-daytona-linked/README.md)               | linked parent and ephemeral child |
| [create-daytona-with-policies](daytona/create-daytona-with-policies/README.md) | pause and deletion policies       |
| [create-daytona-with-network](daytona/create-daytona-with-network/README.md)   | restricted egress allowlists      |
| [create-daytona-with-labels](daytona/create-daytona-with-labels/README.md)     | labels and environment variables  |
| [create-daytona-with-secrets](daytona/create-daytona-with-secrets/README.md)   | injected organization secret      |
| [create-daytona-with-kvm](daytona/create-daytona-with-kvm/README.md)           | nested virtualization (KVM)       |
| [create-daytona-with-ttl](daytona/create-daytona-with-ttl/README.md)           | maximum lifetime TTL              |
| [handle-daytona-errors](daytona/handle-daytona-errors/README.md)               | typed validation errors           |
| [get-daytona-sandbox](daytona/get-daytona-sandbox/README.md)                   | remote retrieval and metadata     |
| [refresh-daytona-sandbox](daytona/refresh-daytona-sandbox/README.md)           | metadata refresh                  |
| [start-daytona-sandbox](daytona/start-daytona-sandbox/README.md)               | start lifecycle operation         |
| [start-daytona-with-timeout](daytona/start-daytona-with-timeout/README.md)     | start with operation timeout      |
| [stop-daytona-sandbox](daytona/stop-daytona-sandbox/README.md)                 | stop lifecycle operation          |
| [force-stop-daytona-sandbox](daytona/force-stop-daytona-sandbox/README.md)     | force stop                        |
| [pause-daytona-sandbox](daytona/pause-daytona-sandbox/README.md)               | pause lifecycle operation         |
| [archive-daytona-sandbox](daytona/archive-daytona-sandbox/README.md)           | archive lifecycle operation       |
| [delete-daytona-sandbox](daytona/delete-daytona-sandbox/README.md)             | waited deletion                   |
| [request-daytona-deletion](daytona/request-daytona-deletion/README.md)         | non-waited deletion               |
| [close-daytona-handle](daytona/close-daytona-handle/README.md)                 | local handle cleanup              |

### Modal

| Project                                                                      | Demonstrates                               |
| ---------------------------------------------------------------------------- | ------------------------------------------ |
| [create-modal-sandbox](modal/create-modal-sandbox/README.md)                 | image creation with resources and policies |
| [create-modal-with-resources](modal/create-modal-with-resources/README.md)   | fractional CPU, memory, and lifetime       |
| [create-modal-with-runtime](modal/create-modal-with-runtime/README.md)       | working directory, entrypoint, and PTY     |
| [create-modal-gvisor](modal/create-modal-gvisor/README.md)                   | gVisor/container isolation                 |
| [create-modal-linux-vm](modal/create-modal-linux-vm/README.md)               | Linux VM isolation                         |
| [create-modal-with-placement](modal/create-modal-with-placement/README.md)   | cloud and region placement                 |
| [create-modal-with-labels](modal/create-modal-with-labels/README.md)         | environment variables and tags             |
| [create-modal-with-network](modal/create-modal-with-network/README.md)       | managed outbound proxy                     |
| [create-modal-idle](modal/create-modal-idle/README.md)                       | idle termination                           |
| [create-modal-ready](modal/create-modal-ready/README.md)                     | command readiness                          |
| [create-modal-tcp-ready](modal/create-modal-tcp-ready/README.md)             | TCP port readiness                         |
| [create-modal-client-timeout](modal/create-modal-client-timeout/README.md)   | ready-wait via client timeout              |
| [create-modal-with-ports](modal/create-modal-with-ports/README.md)           | TLS and plain port bindings                |
| [create-modal-with-allowlists](modal/create-modal-with-allowlists/README.md) | outbound CIDR/domain allowlists            |
| [create-modal-with-secrets](modal/create-modal-with-secrets/README.md)       | named Modal secret                         |
| [create-modal-verbose](modal/create-modal-verbose/README.md)                 | verbose provider logging                   |
| [create-modal-with-identity](modal/create-modal-with-identity/README.md)     | workload identity                          |
| [handle-modal-errors](modal/handle-modal-errors/README.md)                   | typed unsupported-intent errors            |
| [attach-modal-sandbox](modal/attach-modal-sandbox/README.md)                 | ID-only attachment                         |
| [terminate-modal-sandbox](modal/terminate-modal-sandbox/README.md)           | waited termination                         |
| [terminate-modal-with-timeout](modal/terminate-modal-with-timeout/README.md) | waited terminate with timeout              |
| [request-modal-termination](modal/request-modal-termination/README.md)       | non-waited termination                     |
| [close-modal-handle](modal/close-modal-handle/README.md)                     | local handle cleanup                       |

## Verify all TypeScript projects

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit test typescript
```

Builds and tests the SDK packages first, then every example with a `package.json`.

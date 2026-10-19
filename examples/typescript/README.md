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
| [handle-daytona-errors](daytona/handle-daytona-errors/README.md)               | typed validation errors           |
| [get-daytona-sandbox](daytona/get-daytona-sandbox/README.md)                   | remote retrieval and metadata     |
| [refresh-daytona-sandbox](daytona/refresh-daytona-sandbox/README.md)           | metadata refresh                  |
| [start-daytona-sandbox](daytona/start-daytona-sandbox/README.md)               | start lifecycle operation         |
| [stop-daytona-sandbox](daytona/stop-daytona-sandbox/README.md)                 | stop lifecycle operation          |
| [pause-daytona-sandbox](daytona/pause-daytona-sandbox/README.md)               | pause lifecycle operation         |
| [archive-daytona-sandbox](daytona/archive-daytona-sandbox/README.md)           | archive lifecycle operation       |
| [delete-daytona-sandbox](daytona/delete-daytona-sandbox/README.md)             | waited deletion                   |
| [close-daytona-handle](daytona/close-daytona-handle/README.md)                 | local handle cleanup              |

### Modal

| Project                                                                    | Demonstrates                               |
| -------------------------------------------------------------------------- | ------------------------------------------ |
| [create-modal-sandbox](modal/create-modal-sandbox/README.md)               | image creation with resources and policies |
| [create-modal-with-resources](modal/create-modal-with-resources/README.md) | fractional CPU, memory, and lifetime       |
| [create-modal-with-runtime](modal/create-modal-with-runtime/README.md)     | working directory, entrypoint, and PTY     |
| [create-modal-gvisor](modal/create-modal-gvisor/README.md)                 | gVisor/container isolation                 |
| [create-modal-linux-vm](modal/create-modal-linux-vm/README.md)             | Linux VM isolation                         |
| [create-modal-with-placement](modal/create-modal-with-placement/README.md) | cloud and region placement                 |
| [create-modal-with-labels](modal/create-modal-with-labels/README.md)       | environment variables and tags             |
| [create-modal-with-network](modal/create-modal-with-network/README.md)     | managed outbound proxy                     |
| [create-modal-idle](modal/create-modal-idle/README.md)                     | idle termination                           |
| [create-modal-ready](modal/create-modal-ready/README.md)                   | command readiness                          |
| [create-modal-verbose](modal/create-modal-verbose/README.md)               | verbose provider logging                   |
| [create-modal-with-identity](modal/create-modal-with-identity/README.md)   | workload identity                          |
| [handle-modal-errors](modal/handle-modal-errors/README.md)                 | typed unsupported-intent errors            |
| [attach-modal-sandbox](modal/attach-modal-sandbox/README.md)               | ID-only attachment                         |
| [terminate-modal-sandbox](modal/terminate-modal-sandbox/README.md)         | waited termination                         |
| [request-modal-termination](modal/request-modal-termination/README.md)     | non-waited termination                     |
| [close-modal-handle](modal/close-modal-handle/README.md)                   | local handle cleanup                       |

## Verify all TypeScript projects

From `tooling/`:

```sh
GOWORK=off go run ./cmd/sandbox-kit test typescript
```

Builds and tests the SDK packages first, then every example with a `package.json`.

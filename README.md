# Sandbox Kit

Sandbox Kit is being built incrementally as a unified interface over existing
sandbox SDKs. Shared operation requests and responses are the next design step.
Applications supply their already initialized provider SDK clients.

Core, Modal, and Daytona adapters are separate Go modules. Harness adapters are
planned and will also be optional.

## Implemented today

* `core.NewClient(provider)` creates the common `*core.Client`.
* `modal.New(sdkClient)` and `daytona.New(sdkClient)` borrow an existing SDK pointer.
* `client.ProviderName()` identifies the selected provider.
* Constructors reject nil injection and make no provider calls.

Core has no external dependencies. Each adapter depends only on core. Official
SDK imports and authentication configuration stay in the application. The caller
owns SDK cleanup, and the common API does not expose the raw SDK client.

Sandbox creation, execution, lifecycle operations, unified operation responses,
and harness integrations are not implemented yet.

## Try the single example

From the repository root:

```sh
cd examples
go run . --help
go run . modal
# Or: go run . daytona
```

Help requires no credentials. Provider runs need the chosen SDK's configuration;
see [the example guide](examples/README.md) and [main.go](examples/main.go).

The same initialization pattern applies to either provider. See the
[Modal guide](adapters/modal/README.md), [Daytona guide](adapters/daytona/README.md),
and [architecture notes](docs/adapter-design.md).

## Generate and verify

Requires Go 1.26.1 or newer and `protoc`. From the repository root:

```sh
make generate
make test
```

The repository contains several Go modules. Use `make test` from the root;
inside an individual module, use `GOWORK=off go test ./...`.

[Client declarations](proto/kit/core/v1/client.proto) and
[adapter declarations](proto/kit/adapters/) drive our local
`protoc-gen-kit-go` plugin. Generation produces the interfaces, structs,
constructors, and provider identity methods. Protobuf is build-time metadata for
this increment; core and adapters need no protobuf or gRPC runtime.

`make generate` first bootstraps the annotation types, then builds our plugin.
The protobuf Go generator version is pinned in `tooling/go.mod`. Generated source
belongs in version control and should be changed through its declarations and
generator, rather than edited manually.

`go.work` groups the local modules for editor support. Module-level `replace`
directives also allow independent builds inside this checkout. The v0.0.0
requirements are development placeholders; these modules are not published yet.

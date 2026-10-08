# Sandbox Kit

Sandbox Kit is being built incrementally as a unified interface over existing
sandbox SDKs. Shared operation requests and responses are the next design step.
Applications supply their already initialized provider SDK clients.

Go is the first implemented language. Core, Modal, and Daytona adapters are
separate Go modules. Harness adapters are planned and will also be optional.

## Repository layout

```text
proto/                   Shared contracts and generation declarations
sdks/go/core/            Go client runtime
sdks/go/adapters/        Separately installable Go provider adapters
tooling/                 One Go module for all SDK generators
tooling/internal/go/     Emitters for Go SDK output
examples/go/             One runnable Go example
docs/                    Shared architecture and design notes
```

New languages get sibling `sdks/<language>/` and `examples/<language>/`
directories when implemented. Their emitters go under `tooling/internal/<language>/`
and are also written in Go. All emitters reuse `proto/` and apply the output
language's constructor, packaging, and async/streaming conventions.

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

## Try the Go example

From the repository root:

```sh
cd examples/go
go run . --help
go run . modal
# Or: go run . daytona
```

Help requires no credentials. Provider runs need the chosen SDK's configuration;
see [the Go example guide](examples/go/README.md) and [main.go](examples/go/main.go).

See the [Modal guide](sdks/go/adapters/modal/README.md),
[Daytona guide](sdks/go/adapters/daytona/README.md), and
[architecture notes](docs/adapter-design.md) for details.

## Generate and verify

Requires Go 1.26.1 or newer; generation also requires `protoc`.
From the repository root:

```sh
cd tooling
GOWORK=off go install ./cmd/sandbox-kit
sandbox-kit --help
sandbox-kit generate go
sandbox-kit test go
```

These reusable Cobra commands own the workflows directly. The repository
contains several Go modules; inside one module, use `GOWORK=off go test ./...`.

Installation writes the executable to `GOBIN`, or `$(go env GOPATH)/bin` when
`GOBIN` is unset. That directory must be on your shell's `PATH`. Without
installing, use `GOWORK=off go run ./cmd/sandbox-kit generate go` from `tooling/`.
After installation, commands can also run from the repository root. From an
unrelated directory, pass `--repo /absolute/path/to/kit`.

[Client declarations](proto/kit/core/v1/client.proto) and
[adapter declarations](proto/kit/adapters/) drive the Go
`protoc-gen-kit-go` plugin. Generation produces interfaces, structs,
constructors, and provider identity methods. Protobuf is build-time metadata for
this increment; core and adapters need no protobuf or gRPC runtime.

Declaration names are language-neutral. Constructor naming belongs to the target
generator. Standard `go_package` options direct Go output without changing shared
contract semantics. Future generators will apply their native language conventions.

Generation first bootstraps the annotation types, then builds the Go plugin.
Generator binaries are built in a system temporary directory and removed when
generation finishes. Go uses its normal external build cache; no `work/` folder
is needed in the repository.
The protobuf Go generator version is pinned in `tooling/go.mod`. Generated
source belongs in version control and should be changed through its declarations
and generator, rather than edited manually.

The CLI discovers the repository from the current directory. Use `--repo`,
`--go-binary`, or `--protoc` to select explicit paths. Generation cleans up its
temporary tool directory on success or failure.

Go is required to build the generator toolchain. Users of future non-Go SDKs will
install their native packages without needing Go or this tooling module.

`go.work` groups the local Go modules for editor support. Module-local
`replace` directives also allow independent builds inside this checkout.
The v0.0.0 requirements are development placeholders; modules are not published yet.

Go import paths now include `sdks/go`, for example
`github.com/sandbox-kit/kit/sdks/go/core` and
`github.com/sandbox-kit/kit/sdks/go/adapters/modal`.

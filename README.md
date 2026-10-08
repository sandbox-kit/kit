# Sandbox Kit

Sandbox Kit is being built incrementally as a unified interface over existing
sandbox SDKs. A shared creation contract and client API are now implemented.
Applications supply a common client configuration; optional providers initialize
their provider SDK internally.

Go is the first implemented language. Sandbox, Modal, and Daytona adapters are
separate Go modules. Harness adapters are planned and will also be optional.

## Repository layout

```text
proto/                   Shared contracts and generation declarations
specs/                   Portable validation and provider mapping YAML
sdks/go/sandbox/            Go client runtime
sdks/go/providers/        Separately installable Go provider integrations
tooling/                 One Go module for all SDK generators
tooling/internal/go/     Emitters for Go SDK output
examples/go/             Independent provider examples
docs/                    Shared architecture and design notes
```

New languages get sibling `sdks/<language>/` and `examples/<language>/`
directories when implemented. Their emitters go under `tooling/internal/<language>/`
and are also written in Go. All emitters reuse `proto/` and apply the output
language's constructor, packaging, and async/streaming conventions.

## Implemented today

* `sandbox.NewClient(sandbox.Config{...})` validates common configuration and initializes a provider SDK.
* `modal.New()` and `daytona.New()` select an optional provider.
* `client.ProviderName()` identifies the selected provider.
* `client.Create(ctx, request)` delegates to a creation binding and returns a
  common `Sandbox` handle with `ID()`, `ProviderName()`, and `Info()`.
* Optional typed adapters depend on their respective official SDK and own creation mappings.
* Common auth, endpoint, region, context, and operation timeout settings are generated.
* `client.Close(ctx)` releases the initialized SDK; sandbox deletion is separate.
* Native Go configuration, validator tags/cross-field checks, and resource mappings
  are generated from protobuf and language-neutral YAML specs.

Sandbox has no provider SDK or protobuf runtime dependency. It uses
`go-playground/validator/v10` for generated validation. Installing an adapter also
installs its official SDK. Factories map common settings into SDK initialization;
the client owns SDK cleanup. Lifecycle operations and harness integrations remain future increments.
See [client configuration](docs/configuration.md) and
[creation semantics](docs/sandbox-creation.md).

## Try the Go examples

Choose one standalone project:

```sh
cd examples/go/create-modal-sandbox
# Or: cd examples/go/create-daytona-sandbox
cp .env.example .env  # Skip if .env already exists.
# Fill in the provider's credentials and required settings.
go run .
```

Each project authenticates and creates one real sandbox using our SDK. There are
no flags or provider selectors. See [the examples index](examples/go/README.md).

See the [Modal guide](sdks/go/providers/modal/README.md),
[Daytona guide](sdks/go/providers/daytona/README.md), and
[architecture notes](docs/providers.md) for details.

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

[Client declarations](proto/kit/sandbox/v1/client.proto) and
[adapter declarations](proto/kit/providers/) drive the Go
`protoc-gen-kit-go` plugin. Generation produces interfaces, structs,
constructors, native creation data types, and local creation methods. YAML specs
generate validation and selected typed mappings. See [generation specifications](docs/code-generation.md)
for coverage and extension instructions. No gRPC service/client or additional network hop is added.

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
`github.com/sandbox-kit/kit/sdks/go/sandbox` and
`github.com/sandbox-kit/kit/sdks/go/providers/modal`.

See [naming and migration](docs/naming.md) for package, type, field and generated naming conventions.

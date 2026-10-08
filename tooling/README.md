# SDK generation tooling

All generators are written in this Go module. Shared inputs are `proto/` and
`specs/`; SDK runtimes do not need Go tooling or YAML files installed.

```text
cmd/sandbox-kit/              Cobra developer CLI
cmd/protoc-gen-kit-go/        Go output plugin
internal/spec/               Strict YAML loader and portable rule model
internal/go/types_gen/       Native Go structs and validator tags
internal/go/validation_gen/  Shared field and cross-field validation
internal/go/mapping_gen/     Typed request/response mappings and conversions
internal/go/client_gen/      Local client methods and interfaces
internal/go/provider_gen/     Provider constructors, state, SDK cleanup and delegation
internal/gen/               Protobuf annotation bindings used only by tooling
```

From this directory:

```sh
GOWORK=off go run ./cmd/sandbox-kit generate go
GOWORK=off go run ./cmd/sandbox-kit test go
```

Or install with `GOWORK=off go install ./cmd/sandbox-kit`, then use
`sandbox-kit generate go`. Requires Go 1.26.1+ and protoc for generation.
The CLI discovers the repository or accepts `--repo`; `--go-binary` and
`--protoc` select executables. Subprocesses receive cancellation and `GOWORK=off`.
Temporary generator binaries are cleaned up outside the checkout; no `work/`
directory is needed. The normal external Go cache is used.

Generation bootstraps protobuf annotation bindings, then emits native runtime
source. Protobuf is a build-time schema input; the public API has no protobuf
objects and no generated gRPC transport. Client generation emits the common config and provider/backend contracts, a configuration-based
constructor, default operation deadlines, and SDK cleanup delegation. `specs/client.yaml`
defines initialization validation and native provider selection. Provider YAML includes
SDK initialization mappings. Go runtime validation uses
`go-playground/validator/v10` with generated tags and struct-level validators.

See [the spec format and extension workflow](../docs/code-generation.md).
Future emitters belong in `internal/<language>/`; add language SDK bindings while
reusing semantic rules. Other language and harness emitters are not implemented.
Generated source should be changed through specs/generators, then regenerated.

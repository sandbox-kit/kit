# SDK generation tooling

All SDK generators are implemented in this one Go module, regardless of the
language they emit. Shared inputs live in the repository's `proto/` directory.

```text
cmd/protoc-gen-kit-go/    Current Go-output plugin
cmd/sandbox-kit/         Cobra developer CLI
internal/commands/       Reusable generation and verification commands
internal/go/client_gen/  Go client emitter
internal/go/adapter_gen/ Go adapter emitter
internal/gen/            Generated Go bindings for shared annotation metadata
```

Only Go SDK generation exists today. Add future target emitters under
`internal/<language>/` and their entry points under `cmd/`, using Go throughout.
Each emitter owns the target language's naming conventions, syntax, and package
layout. SDK runtime dependencies belong in that language's SDK, not this tooling.

Run the reusable Cobra commands from this directory:

```sh
GOWORK=off go install ./cmd/sandbox-kit
sandbox-kit --help
sandbox-kit generate go
sandbox-kit test go
```

The CLI finds the repository by walking up from the working directory, or accepts
an explicit `--repo` path. `--go-binary` and `--protoc` select the tool executables.
Subprocesses receive cancellation from the command context and run with
`GOWORK=off`. Cobra is a tooling dependency; it is not an SDK runtime dependency.

Ensure `GOBIN` (or `$(go env GOPATH)/bin` when unset) is on `PATH`. For a run
without installation, use `GOWORK=off go run ./cmd/sandbox-kit generate go`.
The installed binary can run from the repository root, or anywhere with an
explicit `--repo /absolute/path/to/kit`.

Generation bootstraps the annotation bindings before building the plugin. Those
bindings are tooling implementation details, not runtime SDK contracts. The
protobuf dependency and generator version are pinned in this module's `go.mod`.

Generator binaries live in a system temporary directory and are cleaned up on
exit. Go's build cache stays outside the checkout. Generated source is the only
generation output kept in the repository.

Go is needed to build and run this toolchain in development or CI. SDK packages
are not published yet. Future generated packages will target their native
runtimes; non-Go SDK consumers will not need Go or the tooling module.

# Go SDK

Go modules under `sdks/go/`:

| Module path                                             | Role                           |
| ------------------------------------------------------- | ------------------------------ |
| `github.com/sandbox-kit/kit/sdks/go/internal/contracts` | Generated contracts (internal) |
| `github.com/sandbox-kit/kit/sdks/go/internal/sdk`       | Shared SDK helpers (internal)  |
| `github.com/sandbox-kit/kit/sdks/go/sandbox`            | Core client                    |
| `github.com/sandbox-kit/kit/sdks/go/providers/daytona`  | Daytona provider               |
| `github.com/sandbox-kit/kit/sdks/go/providers/modal`    | Modal provider                 |
| `github.com/sandbox-kit/kit/sdks/go/providers/blaxel`   | Blaxel provider                |

Callers install `sandbox` plus one provider. Internal modules are not a public
install surface; they ride along as module dependencies.

## Local development

From the repo root, the workspace `go.work` loads these modules for the editor.
CLI and tests use `GOWORK=off` so each module resolves through its own
`go.mod` (with local `replace` directives where present):

```sh
GOWORK=off go run ./tooling/cmd/sandbox-kit test go
```

Change dependencies inside the module you are editing (`GOWORK=off go get` /
`go mod tidy`). Do not use `go work sync` to drive release requires.

## Manual release

Go publishes by tagging module paths. There is no registry upload step. Do not
tag until you intend callers to resolve the new versions via `go get` /
`@latest`.

TypeScript npm cuts are a separate track
([`sdks/typescript/README.md`](../typescript/README.md)).

### 1. Preflight (on `main`)

- Confirm the feature branch for this cut is merged.
- `GOWORK=off go run ./tooling/cmd/sandbox-kit test go`
- Confirm `go.work.sum` is clean (the `test go` workspace-sum check).
- Spot-check that providers still build against the bumped sandbox line.

### 2. Version / require bumps

Bump shared modules in lockstep when the public API or contracts change.
A new provider may start at `v0.1.0` while existing modules take the next
patch (for example sandbox `v0.1.6` with blaxel `v0.1.0`).

Typical order for `go.mod` `require` lines (use `GOWORK=off go get` +
`go mod tidy` per module, then refresh example modules the same way):

1. `internal/contracts`
2. `internal/sdk` (requires contracts)
3. `sandbox` (requires contracts / sdk as needed)
4. Each provider (requires sandbox, and sdk when direct)
5. Example modules under `examples/go/`

Commit the require bumps on `main` (or a short release PR) **before** tagging.
Prior cuts used commit messages like `Require Go module v0.1.N so the release
tags can be fetched.`

### 3. Tag order

Tag the release commit (annotated tags match prior cuts). Prefix each tag with
the module directory under the repo:

```text
sdks/go/internal/contracts/vX.Y.Z
sdks/go/internal/sdk/vX.Y.Z
sdks/go/sandbox/vX.Y.Z
sdks/go/providers/daytona/vX.Y.Z
sdks/go/providers/modal/vX.Y.Z
sdks/go/providers/blaxel/vX.Y.Z
```

Push tags to `origin`. Allow a few minutes for the module proxy / sumdb, then:

```sh
go list -m github.com/sandbox-kit/kit/sdks/go/sandbox@vX.Y.Z
go list -m github.com/sandbox-kit/kit/sdks/go/providers/blaxel@vX.Y.Z
```

### 4. Consumer smoke

In a clean temp module (no `replace`):

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@vX.Y.Z
go get github.com/sandbox-kit/kit/sdks/go/providers/<provider>@vX.Y.Z
```

Build a tiny program that constructs `sandbox.NewClient` with that provider.
Confirm `go get …@latest` resolves to the new tags after the proxy updates.
pkg.go.dev pages can lag behind the proxy.

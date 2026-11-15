# TypeScript SDK

Workspace for the published TypeScript packages:

| Package                  | Role                                             |
| ------------------------ | ------------------------------------------------ |
| `@sandbox-kit/contracts` | Internal wire types (dependency of sandbox only) |
| `@sandbox-kit/sandbox`   | Core client                                      |
| `@sandbox-kit/daytona`   | Daytona provider                                 |
| `@sandbox-kit/modal`     | Modal provider                                   |
| `@sandbox-kit/blaxel`    | Blaxel provider                                  |

## Local development

```sh
cd sdks/typescript
npm install
```

Then build and test from the repo root:

```sh
GOWORK=off go run ./tooling/cmd/sandbox-kit test typescript
```

Examples under `examples/typescript/` keep `file:` dependencies on these
packages for in-checkout runs. Run `npm install` in this workspace first so
provider and contracts dependencies resolve through the workspace
`node_modules`.

## Manual publish

Do not publish until you intend a registry release. Packages are public scoped
packages under the `sandbox-kit` npm org. Go module tags are a separate track
([`sdks/go/README.md`](../go/README.md)).

### 1. Preflight (on `main`)

- Confirm the feature branch for this cut is merged.
- `cd sdks/typescript && npm install`
- `GOWORK=off go run ./tooling/cmd/sandbox-kit test typescript`
- Dry-run each pack and confirm the tarball lists `dist/`, `README.md`,
  `LICENSE`, and versioned `@sandbox-kit/*` deps (no `file:`):

```sh
cd sdks/typescript/contracts && npm pack --dry-run
cd ../sandbox && npm pack --dry-run
cd ../providers/daytona && npm pack --dry-run
cd ../providers/modal && npm pack --dry-run
cd ../providers/blaxel && npm pack --dry-run
```

### 2. Version bumps

Bump packages in lockstep when shared contracts or the core client change. A
new provider may keep `0.1.0` on first publish while existing packages take the
next patch (for example sandbox `0.1.8` with blaxel `0.1.0`).

Update each `package.json` and the workspace `package-lock.json`, then point
each provider at the new `@sandbox-kit/sandbox` version. Refresh example
lockfiles when they pin registry versions (in-checkout examples use `file:`
today).

Commit before publishing (prior pattern:
`Release TypeScript SDK vX.Y.Z.`).

### 3. Publish order

From `sdks/typescript/` after `npm install`:

```sh
cd contracts && npm publish
cd ../sandbox && npm publish
cd ../providers/daytona && npm publish
cd ../providers/modal && npm publish
cd ../providers/blaxel && npm publish
```

`prepublishOnly` deletes `dist/` then rebuilds so stale outputs are not packed.

### 4. Verify

```sh
npm view @sandbox-kit/contracts version
npm view @sandbox-kit/sandbox version
npm view @sandbox-kit/daytona version
npm view @sandbox-kit/modal version
npm view @sandbox-kit/blaxel version
```

In a fresh directory:

```sh
npm install @sandbox-kit/sandbox@X.Y.Z @sandbox-kit/blaxel@X.Y.Z
```

Construct `new Client({ provider: blaxel(), … })` and confirm install snippets
in the root README and website match the published names.

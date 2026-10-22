# TypeScript SDK

Workspace for the published TypeScript packages:

| Package                  | Role                                             |
| ------------------------ | ------------------------------------------------ |
| `@sandbox-kit/contracts` | Internal wire types (dependency of sandbox only) |
| `@sandbox-kit/sandbox`   | Core client                                      |
| `@sandbox-kit/daytona`   | Daytona provider                                 |
| `@sandbox-kit/modal`     | Modal provider                                   |

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

Do not publish until you intend a registry release. Packages are already
configured for public scoped publish under the `sandbox-kit` npm org.

Publish order (contracts first):

```sh
cd sdks/typescript
npm install

cd contracts && npm publish
cd ../sandbox && npm publish
cd ../providers/daytona && npm publish
cd ../providers/modal && npm publish
```

Before a real publish, dry-run each package:

```sh
cd sdks/typescript/contracts && npm pack --dry-run
cd ../sandbox && npm pack --dry-run
cd ../providers/daytona && npm pack --dry-run
cd ../providers/modal && npm pack --dry-run
```

`prepublishOnly` deletes `dist/` then rebuilds so stale outputs are not packed.
Confirm each tarball includes `dist/`, `README.md`, and `LICENSE`, and that
package dependencies use version numbers (not `file:`).

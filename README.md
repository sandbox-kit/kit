# Sandbox Kit

A multi-provider SDK for creating and managing cloud sandboxes. Choose a
provider, install the core SDK and its provider package, then create and manage
sandboxes through a consistent client interface. Daytona and Modal are
supported.

## Get started

Packages are not published yet. The commands below are the post-release install
commands; until publication, use the runnable projects in [Examples](#examples).

### Go

Requires Go 1.26.0 or later. Install the core SDK and one provider:

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@latest
go get github.com/sandbox-kit/kit/sdks/go/providers/daytona@latest
# or: go get github.com/sandbox-kit/kit/sdks/go/providers/modal@latest
```

Set `DAYTONA_API_KEY`, then create a sandbox:

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
			ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
		}},
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		panic(err)
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(instance.ID())
}
```

### TypeScript

Requires Node.js 20, 22, or 24.15+. Install the core SDK and one provider:

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/daytona
# or: npm install @sandbox-kit/sandbox @sandbox-kit/modal
```

Set `DAYTONA_API_KEY`, then create a sandbox:

```ts
import { daytona } from "@sandbox-kit/daytona";
import { Client } from "@sandbox-kit/sandbox";

await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
});

const sandbox = await client.create();
console.log(sandbox.id);
```

Use `await client.close()` in a `finally` block when `await using` is not
available. Timeouts are milliseconds in TypeScript.

### What gets installed

| Language   | Core SDK                                     | Providers                                      |
| ---------- | -------------------------------------------- | ---------------------------------------------- |
| Go         | `github.com/sandbox-kit/kit/sdks/go/sandbox` | `.../providers/daytona`, `.../providers/modal` |
| TypeScript | `@sandbox-kit/sandbox`                       | `@sandbox-kit/daytona`, `@sandbox-kit/modal`   |

Install only the provider you use. The core SDK exposes the shared client
interface; the provider package connects it to the provider's native SDK. Do
not install or import the contracts package directly. The Go wire package is
under `internal`, so an application import does not compile.

Closing a client or sandbox handle releases local SDK resources. It does not
stop, terminate, archive, or delete the cloud sandbox.

## Examples

Runnable projects live under `examples/go/` and `examples/typescript/`. Each
selects one provider, loads its own `.env`, and demonstrates one operation.

```sh
git clone https://github.com/sandbox-kit/kit.git

# Go (Daytona)
cd kit/examples/go/daytona/create-daytona-sandbox
cp .env.example .env   # set DAYTONA_API_KEY
go run .

# TypeScript (Daytona)
cd kit/examples/typescript/daytona/create-daytona-sandbox
cp .env.example .env   # set DAYTONA_API_KEY
npm install && npm run build && npm start
```

Modal needs a token pair and an existing app:
[Go](examples/go/modal/create-modal-sandbox/README.md) ·
[TypeScript](examples/typescript/modal/create-modal-sandbox/README.md).

Creation examples provision real sandboxes. Lifecycle examples use the ID in
`.env` (stop / terminate / archive / delete can destroy that resource).
Error-handling examples fail locally and create nothing. Skip copying
`.env.example` if you already have a configured `.env`.

Full catalog: [examples/README.md](examples/README.md) ·
[Go](examples/go/README.md) · [TypeScript](examples/typescript/README.md).

## Handle errors

Same portable kinds in both languages. Native causes stay on Go `Unwrap` /
`errors.As` and the TypeScript `cause` property. Kinds do not imply safe retries.

### Go

```go
import (
    "errors"
    "fmt"
    "os"

    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func handleError(err error) {
    var sandboxError *sandbox.Error
    if !errors.As(err, &sandboxError) {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    switch sandboxError.Kind {
    case sandbox.ErrorKindInvalidArgument, sandbox.ErrorKindUnsupported:
        fmt.Fprintf(os.Stderr, "Check %s: %s\n", sandboxError.Field, sandboxError.Message)
    case sandbox.ErrorKindAuthentication:
        fmt.Fprintln(os.Stderr, "Check your provider credentials.")
    case sandbox.ErrorKindPermissionDenied:
        fmt.Fprintln(os.Stderr, "Your account does not have permission for this operation.")
    case sandbox.ErrorKindRateLimited, sandbox.ErrorKindResourceExhausted:
        fmt.Fprintln(os.Stderr, "The provider's rate or resource limit was reached.")
    case sandbox.ErrorKindCanceled:
        fmt.Fprintln(os.Stderr, "The operation was canceled.")
    case sandbox.ErrorKindTimeout:
        fmt.Fprintln(os.Stderr, "The operation timed out.")
    default:
        fmt.Fprintln(os.Stderr, err)
    }
}
```

### TypeScript

```ts
import { ErrorKind, SandboxError } from "@sandbox-kit/sandbox";

function handleError(error: unknown): void {
  if (!(error instanceof SandboxError)) {
    console.error(error);
    return;
  }

  switch (error.kind) {
    case ErrorKind.InvalidArgument:
    case ErrorKind.Unsupported:
      console.error(`Check ${error.field}: ${error.message}`);
      break;
    case ErrorKind.Authentication:
      console.error("Check your provider credentials.");
      break;
    case ErrorKind.PermissionDenied:
      console.error(
        "Your account does not have permission for this operation.",
      );
      break;
    case ErrorKind.RateLimited:
    case ErrorKind.ResourceExhausted:
      console.error("The provider's rate or resource limit was reached.");
      break;
    case ErrorKind.Canceled:
      console.error("The operation was canceled.");
      break;
    case ErrorKind.Timeout:
      console.error("The operation timed out.");
      break;
    default:
      console.error(error);
  }
}
```

More detail: [responses and errors](docs/responses-and-errors.md).

## Configuration

- Go: `NewClient(provider, &Config{...})`. TypeScript: `new Client({ provider, ... })`.
- Config covers auth, scope, endpoint, region, and default timeout.
- Create options cover source, resources, runtime, network, and policies.
- Go helpers: `Value(v)` for optional scalars, `Duration(d)` for timeouts.
  TypeScript: property presence; duration fields are milliseconds.
- Unsupported settings fail before provider calls; quotas stay provider-enforced.
- Memory and disk are MiB. Modal CPU: fractional cores from 0.125 (0.001 steps).
  Daytona: whole cores and whole GiB expressed as MiB.

Guides: [configuration](docs/configuration.md) ·
[creation](docs/sandbox-creation.md) ·
[lifecycle](docs/sandbox-lifecycle.md).

## Documentation

| Guide                                                                                                          | Covers                                         |
| -------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| [Client configuration](docs/configuration.md)                                                                  | Auth, scope, defaults, deadlines               |
| [Sandbox creation](docs/sandbox-creation.md)                                                                   | Fields, units, policies                        |
| [Sandbox lifecycle](docs/sandbox-lifecycle.md)                                                                 | Get, attach, capabilities                      |
| [Provider architecture](docs/providers.md)                                                                     | Optional modules and SDK boundaries            |
| [Responses and errors](docs/responses-and-errors.md)                                                           | Origins, error details, native causes          |
| [Code generation](docs/code-generation.md)                                                                     | Protobuf, YAML inventories, extension workflow |
| [Go Daytona](sdks/go/providers/daytona/README.md) / [Modal](sdks/go/providers/modal/README.md)                 | Go provider notes                              |
| [TS Daytona](sdks/typescript/providers/daytona/README.md) / [Modal](sdks/typescript/providers/modal/README.md) | TypeScript provider notes                      |

## Develop and contribute

From `tooling/` (needs Go 1.26.0+; contract generation needs the Buf CLI):

```sh
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit generate contracts-go
GOWORK=off go run ./cmd/sandbox-kit test typescript
GOWORK=off go run ./cmd/sandbox-kit generate contracts-typescript
```

| Path               | Contents                                       |
| ------------------ | ---------------------------------------------- |
| `proto/`           | Portable protobuf contracts                    |
| `specs/`           | YAML validation and provider inventories       |
| `sdks/go/`         | Go SDK, providers, contracts                   |
| `sdks/typescript/` | TypeScript SDK, providers, contracts           |
| `examples/`        | Independent Go and TypeScript example projects |
| `tooling/`         | `sandbox-kit` CLI                              |
| `docs/`            | Guides                                         |

See [tooling](tooling/README.md) and [AGENTS.md](AGENTS.md).

## License

Sandbox Kit uses the [MIT license](LICENSE).

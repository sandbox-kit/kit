# Sandbox Kit

A multi-provider SDK for creating and managing cloud sandboxes from one portable
API. Install the core client and the provider you need; each provider is a
separate package that initializes its official SDK. Go and TypeScript share
contracts, Protovalidate create and client-config rules, error kinds, and the
same Daytona / Modal example matrix. More providers plug into the core sandbox
module without being bundled inside it.

Closing a client or handle releases local SDK resources only. It does not stop,
terminate, archive, or delete a cloud sandbox. The public API covers client
setup, creation, and lifecycle against Daytona and Modal. Execution, storage,
GPU requests, and harness integrations are out of scope.

## Usage

### Go

```go
import (
    "context"
    "fmt"
    "os"
    "time"

    "github.com/sandbox-kit/kit/sdks/go/providers/daytona"
    "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func createSandbox() (result error) {
    client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
        Auth: &sandbox.AuthConfig{
            Credentials: &sandbox.AuthConfig_ApiKey{
                ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
            },
        },
        Timeout: sandbox.Duration(2 * time.Minute),
    })
    if err != nil {
        return err
    }
    defer func() {
        if closeErr := client.Close(context.Background()); closeErr != nil {
            if result == nil {
                result = closeErr
            } else {
                handleError(closeErr)
            }
        }
    }()

    instance, err := client.Create(context.Background(), nil)
    if err != nil {
        return err
    }
    fmt.Println(instance.ID(), instance.ProviderName())
    return nil
}
```

Pass the provider as the first `NewClient` argument. Wire local `go.mod`
replacements as in the example modules. Requires Go 1.26.9 or newer.

### TypeScript

```ts
import { daytona } from "@sandbox-kit/daytona";
import { Client, type ClientConfig } from "@sandbox-kit/sandbox";

function clientConfig(): ClientConfig {
  return {
    provider: daytona(),
    auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
    timeout: 120_000,
  };
}

await using client = new Client(clientConfig());
const sandbox = await client.create();
console.log(sandbox.id, client.providerName);
```

Install the core client and a provider:

```sh
npm install @sandbox-kit/sandbox @sandbox-kit/daytona
# or: npm install @sandbox-kit/sandbox @sandbox-kit/modal
```

`provider` lives on `ClientConfig`. Prefer `await using`; otherwise call
`await client.close()` in `finally`. Requires Node.js 20 / 22 / 24.15+ and
npm 11.19.0. In this checkout, examples use `file:` dependencies on the local
packages under `sdks/typescript/`.

| Package                  | Role                                                        |
| ------------------------ | ----------------------------------------------------------- |
| `@sandbox-kit/sandbox`   | Core client (publish once; more providers depend on this)   |
| `@sandbox-kit/daytona`   | Daytona provider — published independently (`@daytona/sdk`) |
| `@sandbox-kit/modal`     | Modal provider — published independently (`modal`)          |
| `@sandbox-kit/contracts` | Internal wire types — do not depend on this directly        |

TypeScript apps use nested auth/source fields and millisecond timeouts on public
`ClientConfig` / `CreateOptions`. The client materializes protobuf privately
before validation.

Go uses the same module split: `sdks/go/sandbox` is the core module;
`sdks/go/providers/<name>` are independent modules.

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

From `tooling/` (needs Go 1.26.9+; contract generation needs the Buf CLI):

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

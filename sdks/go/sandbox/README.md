# Sandbox Kit for Go

Create and manage cloud sandboxes with a consistent Go client interface. Install
a provider module (`daytona` or `modal`), pass it to `sandbox.NewClient`, then
create or open sandboxes through the same client.

Closing a client or handle releases **local** SDK resources only. It does not
stop, terminate, or delete the remote sandbox.

Requires Go 1.26.0 or newer.

## Install

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox
go get github.com/sandbox-kit/kit/sdks/go/providers/daytona
# or
go get github.com/sandbox-kit/kit/sdks/go/providers/modal
```

## Quick start

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
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.APIKey{
				ApiKey: &sandbox.APIKeyCredentials{
					Key: os.Getenv("DAYTONA_API_KEY"),
				},
			},
		},
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
	fmt.Println(instance.ID(), instance.ProviderName())
}
```

## Client configuration

Pass the provider as the first `NewClient` argument, then shared config:

| Field      | Purpose                                                         |
| ---------- | --------------------------------------------------------------- |
| Provider   | `daytona.New()`, `modal.New()`, …                               |
| `Auth`     | One credential variant; omit for SDK environment credentials    |
| `Scope`    | App, environment, organization, or project (provider-dependent) |
| `Endpoint` | API URL when the provider supports it                           |
| `Region`   | Target or placement preference                                  |
| `Timeout`  | Default operation deadline (`sandbox.Duration`)                 |

Use `sandbox.Value(v)` for optional scalars (including explicit zero) and
`sandbox.Duration(d)` for timeouts. Provider-specific auth and scope rules live
in each provider module README.

## Create a sandbox

```go
instance, err := client.Create(ctx, &sandbox.CreateOptions{
	Source: &sandbox.SandboxSource{
		Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "python:3.11"},
		},
	},
	Resources: &sandbox.Resources{
		CpuCores:  sandbox.Value(2.0),
		MemoryMib: sandbox.Value(uint64(4096)),
	},
})
```

`Create(ctx, nil)` uses the provider default when that is valid (for example
Daytona's default snapshot).

## Open an existing sandbox

```go
attached, err := client.Attach(ctx, id) // provider must support Attach
fetched, err := client.Get(ctx, id)     // provider-verified lookup
```

## Sandbox handle

| Method         | Purpose                                       |
| -------------- | --------------------------------------------- |
| `ID`           | Sandbox identifier                            |
| `Info`         | Copied metadata                               |
| `Capabilities` | Lifecycle operations this handle can run      |
| `Supports`     | Check one capability                          |
| `Refresh`      | Refresh metadata from the provider            |
| `Start`        | Start                                         |
| `Stop`         | Stop                                          |
| `Pause`        | Pause                                         |
| `Archive`      | Archive                                       |
| `Delete`       | Delete                                        |
| `Terminate`    | Terminate                                     |
| `Close`        | Drop the local handle (not the cloud sandbox) |

Unsupported operations return an error before the provider call. Use
`Supports(...)` when the sandbox class may differ.

## Errors

```go
var detail *sandbox.Error
if errors.As(err, &detail) {
	fmt.Println(detail.Kind, detail.Field, detail.Message)
}
```

## Providers

| Module                                                 | Role              |
| ------------------------------------------------------ | ----------------- |
| `github.com/sandbox-kit/kit/sdks/go/providers/daytona` | Daytona sandboxes |
| `github.com/sandbox-kit/kit/sdks/go/providers/modal`   | Modal sandboxes   |

## Examples

Runnable projects:
[examples/go](https://github.com/sandbox-kit/kit/tree/main/examples/go).

## License

MIT

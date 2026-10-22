# Daytona provider for Go

Daytona provider for
[`sdks/go/sandbox`](https://github.com/sandbox-kit/kit/tree/main/sdks/go/sandbox).
Use the shared client API; this module supplies `daytona.New()`.

Requires Go 1.26.9 or newer and a Daytona API key (or credentials your Daytona
SDK already resolves from the environment).

## Install

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox
go get github.com/sandbox-kit/kit/sdks/go/providers/daytona
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

	// Default snapshot for the target:
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(instance.ID())
}
```

Create from an image and size resources:

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
		DiskMib:   sandbox.Value(uint64(8192)),
	},
})
```

Closing the client releases local SDK resources only. Daytona keeps the sandbox
until you stop or delete it, or a lifetime policy fires.

## Configuration

| Setting         | Notes                                                        |
| --------------- | ------------------------------------------------------------ |
| API key         | `Auth` with `APIKey`                                         |
| Bearer / JWT    | Supported; organization ID required                          |
| Omitted `Auth`  | Daytona SDK environment credentials                          |
| Organization    | `Scope.OrganizationId`                                       |
| Endpoint        | `Config.Endpoint` — defaults to `https://app.daytona.io/api` |
| Region / target | `Config.Region`                                              |

Remote endpoints require HTTPS. HTTP is allowed for localhost or loopback only.

## What you can do

**Create** from the default snapshot, a named snapshot, or a registry image.
Optional runtime language/user, network allowlists, secrets, labels, TTL,
ephemeral behavior, and stop/pause/archive/delete policies.

**Lifecycle** on the handle: `Get`, `Refresh`, `Start`, `Stop`, `Pause`,
`Archive`, `Delete`, and local `Close`. Capability depends on the sandbox class
(container, VM, Windows, GPU). Check with `instance.Supports(...)` or
`instance.Capabilities()`.

Daytona does not support unverified `Attach`. Use `client.Get(ctx, id)` for a
provider-verified handle.

CPU requests are whole cores. Memory and disk are whole GiB expressed as MiB.
Resource overrides on create require an **image** source; snapshot creates
inherit snapshot resources.

## Examples

| Example    | Project                                                                                                                       |
| ---------- | ----------------------------------------------------------------------------------------------------------------------------- |
| Creation   | [create-daytona-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/daytona/create-daytona-sandbox)             |
| Policies   | [create-daytona-with-policies](https://github.com/sandbox-kit/kit/tree/main/examples/go/daytona/create-daytona-with-policies) |
| Retrieval  | [get-daytona-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/daytona/get-daytona-sandbox)                   |
| Full index | [Go examples](https://github.com/sandbox-kit/kit/tree/main/examples/go)                                                       |

## License

MIT

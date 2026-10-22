# Modal provider for Go

Modal provider for
[`sdks/go/sandbox`](https://github.com/sandbox-kit/kit/tree/main/sdks/go/sandbox).
Use the shared client API; this module supplies `modal.New()`.

Requires Go 1.26.9 or newer, and a Modal token pair.

## Install

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox
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

	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.TokenPair{
				TokenPair: &sandbox.TokenPairCredentials{
					Id:     os.Getenv("MODAL_TOKEN_ID"),
					Secret: os.Getenv("MODAL_TOKEN_SECRET"),
				},
			},
		},
		Scope: &sandbox.Scope{
			AppName:     sandbox.Value(os.Getenv("MODAL_APP_NAME")),
			Environment: sandbox.Value("main"),
		},
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		panic(err)
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{
			Kind: &sandbox.Image{
				Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
			},
		},
		Resources: &sandbox.Resources{
			CpuCores:  sandbox.Value(0.5),
			MemoryMib: sandbox.Value(uint64(512)),
		},
		Lifetime: &sandbox.LifetimePolicy{
			MaximumLifetime: sandbox.Duration(5 * time.Minute),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(instance.ID())
}
```

Closing the client releases local SDK resources only. Reclaim sandboxes with
idle/max lifetime policies or `Terminate`.

## Configuration

| Setting        | Notes                                                      |
| -------------- | ---------------------------------------------------------- |
| Token pair     | `Auth` with `TokenPair` (`Id` + `Secret`)                  |
| OAuth refresh  | Supported with client ID and exactly one secret or JWT key |
| Omitted `Auth` | Modal SDK environment / profile credentials                |
| App            | `Scope.AppName` — required                                 |
| Environment    | `Scope.Environment` (often `main`)                         |
| Region         | `Config.Region` — default placement preference             |
| Endpoint       | Not supported                                              |

## What you can do

**Create** from a registry image with runtime (command, workdir, PTY), gVisor or
Linux VM isolation, CPU/memory, placement, network policy and ports, secrets,
OIDC workload identity, readiness probes, timeouts, labels, and verbose mode.

**Lifecycle** on the handle: `Attach`, `Terminate`, and local `Close`.
Provider-verified `Get` / `Refresh`, start, stop, pause, archive, and delete are
unsupported.

When you set resources, Modal requires at least 0.125 CPU cores and 128 MiB
memory. Omitted values keep Modal defaults. Disk size cannot be raised at create.

## Examples

| Example    | Project                                                                                                                   |
| ---------- | ------------------------------------------------------------------------------------------------------------------------- |
| Creation   | [create-modal-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/modal/create-modal-sandbox)               |
| Resources  | [create-modal-with-resources](https://github.com/sandbox-kit/kit/tree/main/examples/go/modal/create-modal-with-resources) |
| Attach     | [attach-modal-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/modal/attach-modal-sandbox)               |
| Full index | [Go examples](https://github.com/sandbox-kit/kit/tree/main/examples/go)                                                   |

## License

MIT

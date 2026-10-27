# Modal

Create and manage Modal sandboxes with the shared sandbox client.

Requires Go 1.26.0 or later and an existing Modal app. Set `Auth` to a token
pair or OAuth refresh credentials. If `Auth` is omitted, the Modal SDK uses
credentials it already finds in the environment.

## Install

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@latest
go get github.com/sandbox-kit/kit/sdks/go/providers/modal@latest
```

## Create a sandbox

```go
client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
    Auth: &sandbox.AuthConfig{Credentials: &sandbox.TokenPair{
        TokenPair: &sandbox.TokenPairCredentials{
            Id:     os.Getenv("MODAL_TOKEN_ID"),
            Secret: os.Getenv("MODAL_TOKEN_SECRET"),
        },
    }},
    Scope: &sandbox.Scope{
        AppName:     sandbox.Value(os.Getenv("MODAL_APP_NAME")),
        Environment: sandbox.Value("main"),
    },
    Timeout: sandbox.Duration(2 * time.Minute),
})
if err != nil {
    log.Fatal(err)
}
defer client.Close(context.Background())

instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
        Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
    }},
    Resources: &sandbox.Resources{
        CpuCores:  sandbox.Value(0.5),
        MemoryMib: sandbox.Value(uint64(512)),
    },
    Lifetime: &sandbox.LifetimePolicy{
        MaximumLifetime: sandbox.Duration(5 * time.Minute),
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(instance.ID())
```

Closing the client does not delete the sandbox. Set a lifetime policy, or call
`Terminate`.

## Settings

| Setting       | How to set it                     |
| ------------- | --------------------------------- |
| Token pair    | `Auth` with `TokenPair`           |
| OAuth refresh | `Auth` with `OAuthRefresh`        |
| App           | `Scope.AppName`, required         |
| Environment   | `Scope.Environment`, often `main` |
| Region        | `Config.Region`                   |

`Config.Endpoint` is not supported.

## Operations

Create from a registry image. You can set the command, working directory, PTY,
gVisor or Linux VM isolation, CPU and memory, placement, network policy, ports,
secrets, readiness, timeout, and labels.

On the handle: `Attach`, `Terminate`, and `Close`. `Close` does not terminate
the sandbox. `Get`, `Refresh`, `Start`, `Stop`, `Pause`, `Archive`, and
`Delete` are not supported.

When you set resources, Modal requires at least 0.125 CPU cores and 128 MiB of
memory. Omitted values keep Modal's defaults. Disk size cannot be raised at
create time.

## Examples

| Operation | Program                                                                                                                   |
| --------- | ------------------------------------------------------------------------------------------------------------------------- |
| Create    | [create-modal-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/modal/create-modal-sandbox)               |
| Resources | [create-modal-with-resources](https://github.com/sandbox-kit/kit/tree/main/examples/go/modal/create-modal-with-resources) |
| Attach    | [attach-modal-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/modal/attach-modal-sandbox)               |
| All       | [Go examples](https://github.com/sandbox-kit/kit/tree/main/examples/go)                                                   |

## License

MIT

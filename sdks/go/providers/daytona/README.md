# Daytona

Create and manage Daytona sandboxes with the shared sandbox client.

Requires Go 1.26.0 or later and a Daytona API key. If `Auth` is omitted, the
Daytona SDK uses credentials it already finds in the environment.

## Install

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@latest
go get github.com/sandbox-kit/kit/sdks/go/providers/daytona@latest
```

## Create a sandbox

```go
client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
    Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
        ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
    }},
    Timeout: sandbox.Duration(2 * time.Minute),
})
if err != nil {
    log.Fatal(err)
}
defer client.Close(context.Background())

instance, err := client.Create(context.Background(), nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(instance.ID())
```

A nil create request uses Daytona's default snapshot. To choose an image and
a size:

```go
instance, err = client.Create(ctx, &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
        Image: &sandbox.ImageSource{Reference: "python:3.11"},
    }},
    Resources: &sandbox.Resources{
        CpuCores:  sandbox.Value(2.0),
        MemoryMib: sandbox.Value(uint64(4096)),
        DiskMib:   sandbox.Value(uint64(8192)),
    },
})
```

Closing the client does not delete the sandbox. Stop or delete it, or set a
lifetime policy.

## Settings

| Setting      | How to set it                                                  |
| ------------ | -------------------------------------------------------------- |
| API key      | `Auth` with `APIKey`                                           |
| Bearer token | `Auth` with `BearerToken`, and `Scope.OrganizationId`          |
| Endpoint     | `Config.Endpoint`. The default is `https://app.daytona.io/api` |
| Region       | `Config.Region`                                                |

A remote endpoint uses HTTPS. HTTP is accepted for localhost and loopback
addresses.

## Operations

Create from the default snapshot, a named snapshot, or a registry image. You
can set the runtime user and language, network allowlists, secrets, labels,
TTL, ephemeral mode, and stop, pause, archive, and delete policies.

On the handle: `Get`, `Refresh`, `Start`, `Stop`, `Pause`, `Archive`,
`Delete`, and `Close`. `Get` checks the ID with Daytona. `Attach` is not
supported. Daytona accepts or rejects each lifecycle call for that sandbox.
`Close` does not contact Daytona and does not change the sandbox.

CPU is a whole number of cores. Memory and disk are a whole number of GiB,
written as MiB. A resource override requires an image source. A snapshot create
uses the snapshot's resources.

## Examples

| Operation | Program                                                                                                                       |
| --------- | ----------------------------------------------------------------------------------------------------------------------------- |
| Create    | [create-daytona-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/daytona/create-daytona-sandbox)             |
| Policies  | [create-daytona-with-policies](https://github.com/sandbox-kit/kit/tree/main/examples/go/daytona/create-daytona-with-policies) |
| Look up   | [get-daytona-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/daytona/get-daytona-sandbox)                   |
| All       | [Go examples](https://github.com/sandbox-kit/kit/tree/main/examples/go)                                                       |

## License

MIT

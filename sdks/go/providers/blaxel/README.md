# Blaxel

Create and manage Blaxel sandboxes with the shared sandbox client.

Requires Go 1.26.0 or later and a Blaxel API key. Set `Auth` to an API key and
`Scope.Workspace` to the workspace name. If either is omitted, the Blaxel SDK
uses `BL_API_KEY` and `BL_WORKSPACE` from the environment when present.
`Config.Region` is required for create.

## Install

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@latest
go get github.com/sandbox-kit/kit/sdks/go/providers/blaxel@latest
```

## Create a sandbox

```go
client, err := sandbox.NewClient(blaxel.New(), &sandbox.Config{
    Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
        ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("BL_API_KEY")},
    }},
    Scope: &sandbox.Scope{
        Workspace: sandbox.Value(os.Getenv("BL_WORKSPACE")),
    },
    Region:  sandbox.Value("us-pdx-1"),
    Timeout: sandbox.Duration(2 * time.Minute),
})
if err != nil {
    log.Fatal(err)
}
defer client.Close(context.Background())

instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
    Name: sandbox.Value("my-sandbox"),
    Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
        Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"},
    }},
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(instance.ID())
```

Create requires a name and an image. Closing the client does not delete the
sandbox. Call `Delete`, or set a lifetime policy.

## Settings

| Setting   | How to set it                         |
| --------- | ------------------------------------- |
| API key   | `Auth` with `APIKey`                  |
| Workspace | `Scope.Workspace`                     |
| Region    | `Config.Region` (required for create) |
| Timeout   | `Config.Timeout` (create/op deadline) |

`Config.Endpoint` is not supported. Token pair, OAuth, and bearer auth are
rejected. `Scope` fields other than `Workspace` are rejected.

## Operations

Create from a registry image with an explicit name. The name must be lowercase
alphanumeric with single hyphens and at most 49 characters. You can set labels,
environment variables, memory, ports, outbound domains, maximum lifetime, and
idle terminate. Max-age TTL must be at least 5m; idle TTL must be at least 1m.
Ports 80, 443, and 8080 are rejected (Blaxel-reserved). Default port transport
leaves the protocol to Blaxel (HTTP); plain maps to TCP. Kit
`Config.Timeout` / create call timeout is the create deadline. The Go Blaxel
SDK has no separate creation-timeout argument. `WaitConditionReady` is rejected.

On the handle: `Get`, `Refresh`, `Delete`, and `Close`. `Get` and `Refresh` look
up the sandbox by name. `Attach`, `Start`, `Stop`, `Pause`, `Archive`, and
`Terminate` are not supported. `Close` does not contact Blaxel and does not
change the sandbox.

`Filesystem` supports the shared file calls: `ReadFile`, `WriteFile`,
`CopyFromLocal`, `CopyToLocal`, `MakeDirectory`, `ListDirectory`, `Stat`, and
`Remove`. Paths in the sandbox are absolute. A file timeout is the context
deadline. `Parents` defaults to true. `Parents` set to false is rejected. A
missing file passed to `CopyFromLocal` is not found, and the error field is
`local_path`. A local directory is an invalid argument with the same field.
`CopyToLocal` creates a missing local parent directory.

`Process().Run` runs one command and waits for the result. Set exactly one of
`Argv` or `Shell`. `Argv` is sent as a POSIX-quoted shell string. A non-zero
exit code is success for the SDK—read `GetExitCode`. Blaxel returns stdout and
stderr. Omit or zero `Timeout` for no Blaxel process auto-kill. The blocking
wait can still end at Blaxel's edge hold (about one to two minutes) while the
command keeps running. Command timeouts must be whole seconds. Cancel the
context to stop waiting.

`resources.memory_mib` maps 1:1 onto Blaxel's memory field. Blaxel labels that
integer MB; Kit does not convert. Blaxel derives CPU from that value, so
`cpu_cores` and disk overrides are rejected. Kit has no portable volume
fields; Blaxel volumes are outside create.

## Examples

| Operation                | Program                                                                                                                                      |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------- |
| Create                   | [create-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/feat/blaxel-provider/examples/go/blaxel/configuration/create-blaxel-sandbox) |
| Run a command            | [run-blaxel-command](https://github.com/sandbox-kit/kit/tree/feat/blaxel-provider/examples/go/blaxel/process/run-blaxel-command)             |
| Shell / cwd / env / exit | [process examples](https://github.com/sandbox-kit/kit/tree/feat/blaxel-provider/examples/go/blaxel/process)                                  |
| Read a file              | [read-blaxel-file](https://github.com/sandbox-kit/kit/tree/feat/blaxel-provider/examples/go/blaxel/filesystem/read-blaxel-file)              |
| Filesystem               | [filesystem examples](https://github.com/sandbox-kit/kit/tree/feat/blaxel-provider/examples/go/blaxel/filesystem)                            |
| Look up                  | [get-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/feat/blaxel-provider/examples/go/blaxel/lifecycle/get-blaxel-sandbox)           |
| Delete                   | [delete-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/feat/blaxel-provider/examples/go/blaxel/lifecycle/delete-blaxel-sandbox)     |
| All                      | [Go examples](https://github.com/sandbox-kit/kit/tree/feat/blaxel-provider/examples/go)                                                      |

pkg.go.dev also shows `ExampleNew_*` samples in this package for create,
lifecycle, process, and filesystem.

## License

MIT

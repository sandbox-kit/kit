# Blaxel

Create and manage Blaxel sandboxes with the shared sandbox client.

Requires Go 1.26.0 or later and a Blaxel API key. Set `Auth` to an API key and
`Scope.Workspace` to the workspace name. If either is omitted, the client uses
`BL_API_KEY` and `BL_WORKSPACE` from the environment when present.
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

Create requires a name and an image. `ID()` is that name — use it with `Get`
and example env `BL_SANDBOX_NAME`. Closing the client does not delete the
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

## What you can do

### Create

You need a name and an image. Names are lowercase alphanumeric with single
hyphens (max 49 characters). You can also set labels, environment variables,
memory, ports, outbound domains, maximum lifetime (≥5 minutes), and idle
terminate (≥1 minute). You cannot disable idle terminate. Ports 80, 443, and
8080 are reserved. Default port transport is HTTP; plain is TCP; TLS is TLS;
HTTP/2 over TLS is not supported.

Your client/create timeout is how long create may take. Ready wait and
readiness probes are not supported (`WaitFor` only Default or Scheduled).

Set `Resources.MemoryMib` for memory size. Blaxel documents that value as MB;
Kit does not convert units. CPU follows memory on Blaxel — do not set CPU or
disk overrides. If you omit memory, Blaxel applies its default. Volumes are not
part of create.

If the name already exists, you get `AlreadyExists`.

### Lifecycle

You can `Get`, `Refresh`, `Delete`, and `Close`. Lookup is by name
(`ID()` equals the create name). Attach, start, stop, pause, archive, and
terminate are not available. `Close` only drops the local handle.

After a successful `Delete`, the handle reports `Deleted` and `terminal`.
Blaxel can still be cleaning up — recreating the same name right away may
return `Conflict`.

On `Get` / `Refresh` / `Info`, an idle sandbox shows as `Paused`; an active
one shows as `Running`.

### Filesystem

You get the usual file APIs: read, write, copy in/out, mkdir, list, stat, and
remove. Paths must be absolute. `Parents` defaults to true; `Parents: false`
is rejected.

`Stat` returns whether the path is a file or a directory, plus size and mode
when Blaxel provides them. Missing local files for `CopyFromLocal` return not
found with field `local_path`.

### Process

`Process().Run` runs one command and waits for it to finish. Set exactly one of
`Argv` or `Shell`. `Argv` is sent as one quoted shell command string. A
non-zero exit code is not an error — read `GetExitCode`. You get both stdout
and stderr.

Omit or zero `Timeout` to wait without a Blaxel process time limit (your
context can still cancel). A positive timeout must be whole seconds and caps
how long the command may run. Canceling the context stops waiting; the remote
command may keep running. If the wait times out, you get a timeout error.

### Errors you will see

| Situation                                  | Kind            |
| ------------------------------------------ | --------------- |
| Create when the name already exists        | `AlreadyExists` |
| Conflict while a delete is still finishing | `Conflict`      |

## Examples

| Operation                | Program                                                                                                                      |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| Create                   | [create-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/blaxel/configuration/create-blaxel-sandbox) |
| Run a command            | [run-blaxel-command](https://github.com/sandbox-kit/kit/tree/main/examples/go/blaxel/process/run-blaxel-command)             |
| Shell / cwd / env / exit | [process examples](https://github.com/sandbox-kit/kit/tree/main/examples/go/blaxel/process)                                  |
| Read a file              | [read-blaxel-file](https://github.com/sandbox-kit/kit/tree/main/examples/go/blaxel/filesystem/read-blaxel-file)              |
| Filesystem               | [filesystem examples](https://github.com/sandbox-kit/kit/tree/main/examples/go/blaxel/filesystem)                            |
| Look up                  | [get-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/blaxel/lifecycle/get-blaxel-sandbox)           |
| Delete                   | [delete-blaxel-sandbox](https://github.com/sandbox-kit/kit/tree/main/examples/go/blaxel/lifecycle/delete-blaxel-sandbox)     |
| All                      | [Go examples](https://github.com/sandbox-kit/kit/tree/main/examples/go)                                                      |

pkg.go.dev also shows `ExampleNew_*` samples in this package for create,
lifecycle, process, and filesystem.

## License

MIT

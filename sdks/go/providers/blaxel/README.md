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
config := &sandbox.Config{
    Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
        ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("BL_API_KEY")},
    }},
    Region:  sandbox.Value("us-pdx-1"),
    Timeout: sandbox.Duration(2 * time.Minute),
}
if workspace := strings.TrimSpace(os.Getenv("BL_WORKSPACE")); workspace != "" {
    config.Scope = &sandbox.Scope{Workspace: sandbox.Value(workspace)}
}
client, err := sandbox.NewClient(blaxel.New(), config)
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
| Workspace | `Scope.Workspace` (or `BL_WORKSPACE`) |
| Region    | `Config.Region` (required for create) |
| Timeout   | `Config.Timeout` (create/op deadline) |

`Config.Endpoint` is not supported. Token pair, OAuth, and bearer auth are
rejected. `Scope` fields other than `Workspace` are rejected.

## What you can do

### Create

You need a name and an image. Name is required (blank or whitespace-only is
rejected). Names are lowercase alphanumeric with single hyphens (max 49
characters). You can also set labels, environment variables,
memory, ports, outbound domains, maximum lifetime (≥5 minutes), and idle
terminate (≥1 minute). You cannot disable idle terminate. Runtime and
isolation are not supported. `Info().Labels` is set only when Blaxel reports
labels. Ports 80, 443, and 8080 are reserved. Default port transport is HTTP;
plain is TCP; TLS is TLS; HTTP/2 over TLS is not supported.

Set client or create timeout for how long create may wait. With a 1–50s
timeout, a timed-out create frees the name. Ready/Started/Submitted wait and
readiness probes are not supported (`WaitFor` only Default or Scheduled).

If Blaxel’s edge returns HTTP 504 during create, Go returns a timeout as soon
as that response arrives. TypeScript may keep waiting until the sandbox is
ready or your create timeout ends. A Go timeout does not always mean the
sandbox was never created — call `Get` by name before recreating.

Set `Resources.MemoryMib` for memory size. Kit converts MiB to Blaxel's decimal
MB value, rounding up so a request is never under-provisioned. CPU follows
memory on Blaxel — do not set CPU or
disk overrides. If you omit memory, Blaxel applies its default. Volumes are not
part of create.

If the name already exists, you get `AlreadyExists`.

### Lifecycle

You can `Get`, `Refresh`, `Delete`, and `Close`. Lookup is by name
(`ID()` equals the create name). Attach, start, stop, pause, archive, and
terminate are not available. `Close` only drops the local handle.

After a successful `Delete`, `Info().State` is `Deleted` and `Terminal` is
true. Omit `Wait` to return as soon as Blaxel accepts the delete — recreating
the same name too soon may return `Conflict`. Pass `Wait: true` to wait until
removal finishes before the call returns.

### Filesystem

You get the usual file APIs: read, write, copy in/out, mkdir, list, stat, and
remove. Paths must be absolute. `Parents` defaults to true; Blaxel does not
support `Parents: false`.

`Stat` returns `Type` (file or directory), `Size`, and usually `Permissions`,
`Owner`, `Group`, and `ModifiedAt`. `Name` and `Path` are whatever Blaxel
reports. Missing local files for `CopyFromLocal`, and local write failures from
`CopyToLocal`, use field `local_path`.

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

| Situation                                  | Kind              |
| ------------------------------------------ | ----------------- |
| Create when the name already exists        | `AlreadyExists`   |
| Conflict while a delete is still finishing | `Conflict`        |
| Lifetime below floor or not whole seconds  | `InvalidArgument` |
| Unsupported create field (runtime, …)      | `Unsupported`     |

Blaxel control-plane codes appear on `ProviderCode` when present.

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

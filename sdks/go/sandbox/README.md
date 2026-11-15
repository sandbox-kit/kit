# Sandbox client

Create and manage cloud sandboxes with one Go client. Install this module and
one provider (Daytona, Modal, or Blaxel).

Requires Go 1.26.0 or later.

## Install

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@latest
go get github.com/sandbox-kit/kit/sdks/go/providers/daytona@latest
```

For Modal or Blaxel, install
`github.com/sandbox-kit/kit/sdks/go/providers/modal` or
`github.com/sandbox-kit/kit/sdks/go/providers/blaxel` instead of the Daytona
provider.

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

`Create` with a nil request uses the provider default when that default is
valid. Daytona uses its default snapshot. Modal requires an image. Blaxel
requires a name and an image, plus client `Region`.

`client.Close` closes your connection to Daytona, Modal, or Blaxel.
`instance.Close` does not stop, terminate, or delete the sandbox, and you can
keep using the handle. Reclaim a Daytona sandbox with Stop, Archive, or
Delete; a Modal sandbox with Terminate or lifetime; a Blaxel sandbox with
Delete or TTL/idle. Blaxel has no Stop, Terminate, or Archive.

## Run a command

```go
result, err := instance.Process().Run(context.Background(), &sandbox.ExecOptions{
    Argv: []string{"echo", "hello from kit"},
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.GetExitCode(), string(result.GetStdout()))
```

Set exactly one of `Argv` or `Shell`. A non-zero exit code is not an SDK
error—check `GetExitCode`. Optional `WorkingDirectory` must be absolute.
`Environment` applies to that call only. Omit `Timeout` to use the client
deadline; use `Duration(0)` for no extra deadline on this call. Cancel the
context to stop waiting. See
[Process](https://github.com/sandbox-kit/kit/blob/main/docs/process.md).

## Read and write files

```go
files := instance.Filesystem()
if err := files.WriteFile(context.Background(), "/tmp/note.txt", []byte("hello\n"), nil); err != nil {
    log.Fatal(err)
}
data, err := files.ReadFile(context.Background(), "/tmp/note.txt", nil)
if err != nil {
    log.Fatal(err)
}
fmt.Print(string(data))
```

Paths inside the sandbox are absolute. The same handle also supports
`CopyFromLocal`, `CopyToLocal`, `MakeDirectory`, `ListDirectory`, `Stat`, and
`Remove`. See
[Filesystem](https://github.com/sandbox-kit/kit/blob/main/docs/filesystem.md).

## Choose a provider

| Provider                                                                           | Import          | Open an existing sandbox |
| ---------------------------------------------------------------------------------- | --------------- | ------------------------ |
| [Daytona](https://pkg.go.dev/github.com/sandbox-kit/kit/sdks/go/providers/daytona) | `daytona.New()` | `client.Get`             |
| [Modal](https://pkg.go.dev/github.com/sandbox-kit/kit/sdks/go/providers/modal)     | `modal.New()`   | `client.Attach`          |
| [Blaxel](https://pkg.go.dev/github.com/sandbox-kit/kit/sdks/go/providers/blaxel)   | `blaxel.New()`  | `client.Get` (by name)   |

`Supports` tells you whether this handle can attempt that operation. On
Blaxel, you can `Get`, `Refresh`, `Delete`, and `Close`. Other lifecycle
methods are not available. If the provider rejects a supported call, you get
that provider's error.

## Read an error

```go
var detail *sandbox.Error
if errors.As(err, &detail) {
    fmt.Println(detail.Kind, detail.Field, detail.Message)
}
```

`Kind` is the shared category, such as an authentication failure or an
unsupported setting. `Field` is the setting to change. `Unwrap` returns the
provider's own error when there is one.

## Examples

Runnable programs for each operation are in
[examples/go](https://github.com/sandbox-kit/kit/tree/main/examples/go).
Provider packages ship `ExampleNew_*` samples for pkg.go.dev (create,
lifecycle, process, and filesystem).

## License

MIT

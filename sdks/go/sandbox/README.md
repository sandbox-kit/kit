# Sandbox client

Create and manage cloud sandboxes with one Go client. Install this module and
either the Daytona provider or the Modal provider.

Requires Go 1.26.0 or later.

## Install

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@latest
go get github.com/sandbox-kit/kit/sdks/go/providers/daytona@latest
```

For Modal, install `github.com/sandbox-kit/kit/sdks/go/providers/modal` instead
of the Daytona provider.

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
valid. Daytona uses its default snapshot. Modal requires an image.

`client.Close` and `instance.Close` release local resources. The cloud sandbox
keeps running until you stop, terminate, archive, or delete it.

## Choose a provider

| Provider                                                                           | Import          | Open an existing sandbox |
| ---------------------------------------------------------------------------------- | --------------- | ------------------------ |
| [Daytona](https://pkg.go.dev/github.com/sandbox-kit/kit/sdks/go/providers/daytona) | `daytona.New()` | `client.Get`             |
| [Modal](https://pkg.go.dev/github.com/sandbox-kit/kit/sdks/go/providers/modal)     | `modal.New()`   | `client.Attach`          |

Call `instance.Supports` before a lifecycle method when the sandbox class
varies. An unsupported method returns an error and does not call the provider.

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

## License

MIT

# Sandbox Kit for Go

Core client module (`github.com/sandbox-kit/kit/sdks/go/sandbox`). Provider
modules under `sdks/go/providers/<name>` are published independently and plug
into this package. Shared create and client-config validation uses Protovalidate
(`secure_endpoint` stays handwritten). Closing the client releases local
provider SDK resources only; it does not stop or delete cloud sandboxes.

## Usage

```go
client, err := sandbox.NewClient(provider, &sandbox.Config{ /* ... */ })
if err != nil {
    return err
}
defer client.Close(context.Background())

instance, err := client.Create(context.Background(), nil)
```

See the [root README](../../../README.md) and
[provider architecture](../../../docs/providers.md).

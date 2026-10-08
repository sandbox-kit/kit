# Modal provider

Install this optional module and select `modal.New()` in `sandbox.Config.Provider`.
The public `sandbox.NewClient` constructs the official SDK from common settings;
its private backend maps sandbox creation and cleanup operations.

```go
client, err := sandbox.NewClient(sandbox.Config{
    Provider: modal.New(),
})
```

Omitted auth uses SDK environment/profile resolution. Check `err` and call
`client.Close(ctx)` when finished. Close releases SDK resources, not sandboxes.

[Provider YAML](../../../../specs/providers/modal.yaml) generates supported
configuration, credential, resource and response mappings. SDK construction and
remaining semantic mappings live in `client.go` and `sandbox.go`.

See [configuration](../../../../docs/configuration.md),
[sandbox creation](../../../../docs/sandbox-creation.md), and
[the runnable example](../../../../examples/go/README.md).

Verify with `GOWORK=off go test ./...` from this module. Packages are unpublished;
local module replacements and v0.0.0 requirements are development settings.

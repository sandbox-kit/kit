# Modal adapter

This module depends only on core. The application installs and initializes the
official Modal SDK, then passes its client pointer to `New`.

Use these imports in the application:

```go
import (
    modalAdapter "github.com/sandbox-kit/kit/sdks/go/adapters/modal"
    "github.com/sandbox-kit/kit/sdks/go/core"
)
```

Inside the application function, with `sdkClient` already initialized:

```go
provider, err := modalAdapter.New(sdkClient)
if err != nil {
    return err
}
client, err := core.NewClient(provider)
if err != nil {
    return err
}
name := client.ProviderName() // "modal"
```

`New` borrows the pointer through a type parameter and rejects nil. It does not
authenticate, modify configuration, or own SDK cleanup. The application retains its SDK reference and closes it when done.
The borrowed client is private to the adapter.

This increment covers initialization and provider identity. SDK-specific
operations and shared request/response mappings are planned. The constructor
does not validate an arbitrary client's SDK method signatures.

[The protobuf declaration](../../../../proto/kit/adapters/modal/v1/adapter.proto)
generates the adapter shell. From this module directory, verify it with:

```sh
GOWORK=off go test ./...
```

The local core replacement and v0.0.0 dependency are development settings;
the module has not been published. See [the single runnable example](../../../../examples/go/README.md)
for complete SDK initialization, authentication setup, and execution instructions.

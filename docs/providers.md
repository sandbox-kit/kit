# Providers

Install the core package and the provider you use.

|            | Daytona                                                                  | Modal                                                                  |
| ---------- | ------------------------------------------------------------------------ | ---------------------------------------------------------------------- |
| TypeScript | `@sandbox-kit/sandbox` and `@sandbox-kit/daytona`                        | `@sandbox-kit/sandbox` and `@sandbox-kit/modal`                        |
| Go         | `github.com/sandbox-kit/kit/sdks/go/sandbox` and `.../providers/daytona` | `github.com/sandbox-kit/kit/sdks/go/sandbox` and `.../providers/modal` |

```go
client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
    Auth: &sandbox.AuthConfig{
        Credentials: &sandbox.TokenPair{
            TokenPair: &sandbox.TokenPairCredentials{Id: tokenID, Secret: tokenSecret},
        },
    },
    Scope: &sandbox.Scope{AppName: sandbox.Value("existing-app")},
})
```

```ts
await using client = new Client({
  provider: modal(),
  auth: { tokenPair: { id: tokenId, secret: tokenSecret } },
  scope: { appName: "existing-app" },
});
```

## What each provider can do

| Operation                           | Modal                              | Daytona                                    |
| ----------------------------------- | ---------------------------------- | ------------------------------------------ |
| Create                              | Registry image and an existing app | Default snapshot, named snapshot, or image |
| Open an existing sandbox            | `attach` (ID only)                 | `get` (provider lookup)                    |
| Start, stop, pause, archive, delete | Unsupported                        | Supported                                  |
| Terminate                           | Supported                          | Unsupported                                |

`supports` reports whether the provider implements that operation. A missing
operation fails before the provider call. Daytona still accepts or rejects a
supported call for a particular sandbox, and Kit returns that error.

Set Daytona's target with `region` on the client. A placement block on the
create request is rejected. Modal accepts cloud and region on the create
request, and rejects a placement target.

Closing the client closes your connection to the provider. The sandboxes keep
running. Closing a handle leaves the sandbox as it is. You can keep using the
handle. The next call returns the provider's result.

The public API covers create and lifecycle. In-sandbox execution, GPU
requests, resizing, storage, browsers, and macOS or Windows creation are
outside it.

See [configuration](configuration.md), [creation](sandbox-creation.md), and
[lifecycle](sandbox-lifecycle.md).

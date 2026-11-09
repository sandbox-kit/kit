# Providers

Install the core package and the provider you use.

|            | Daytona                                                                  | Modal                                                                  | Blaxel                                                                  |
| ---------- | ------------------------------------------------------------------------ | ---------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| TypeScript | `@sandbox-kit/sandbox` and `@sandbox-kit/daytona`                        | `@sandbox-kit/sandbox` and `@sandbox-kit/modal`                        | `@sandbox-kit/sandbox` and `@sandbox-kit/blaxel`                        |
| Go         | `github.com/sandbox-kit/kit/sdks/go/sandbox` and `.../providers/daytona` | `github.com/sandbox-kit/kit/sdks/go/sandbox` and `.../providers/modal` | `github.com/sandbox-kit/kit/sdks/go/sandbox` and `.../providers/blaxel` |

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


```go
client, err := sandbox.NewClient(blaxel.New(), &sandbox.Config{
    Auth: &sandbox.AuthConfig{
        Credentials: &sandbox.APIKey{
            ApiKey: &sandbox.APIKeyCredentials{Key: apiKey},
        },
    },
    Scope:  &sandbox.Scope{Workspace: sandbox.Value("my-workspace")},
    Region: sandbox.Value("us-pdx-1"),
})
```

```ts
await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: apiKey } },
  scope: { workspace: "my-workspace" },
  region: "us-pdx-1",
});
```

## What each provider can do

| Operation                           | Modal                              | Daytona                                    | Blaxel                                      |
| ----------------------------------- | ---------------------------------- | ------------------------------------------ | ------------------------------------------- |
| Create                              | Registry image and an existing app | Default snapshot, named snapshot, or image | Named image create                          |
| Open an existing sandbox            | `attach` (ID only)                 | `get` (provider lookup)                    | `get` / `refresh` by name                   |
| Start, stop, pause, archive         | Unsupported                        | Supported                                  | Unsupported                                 |
| Delete                              | Unsupported                        | Supported                                  | Supported                                   |
| Terminate                           | Supported                          | Unsupported                                | Unsupported                                 |
| Filesystem                          | Same eight shared methods          | Same eight shared methods                  | Same eight shared methods                   |
| Process                             | Same blocking `Run`                | Same blocking `Run`                        | Same blocking `Run`                         |

`supports` reports whether the provider implements that operation. A missing
operation fails before the provider call. Daytona still accepts or rejects a
supported call for a particular sandbox, and Kit returns that error.

Set Daytona's target with `region` on the client. A placement block on the
create request is rejected. Modal accepts cloud and region on the create
request, and rejects a placement target. Blaxel requires `region` on the
client for create and rejects per-request placement.

Closing the client closes your connection to the provider. The sandboxes keep
running. Closing a handle leaves the sandbox as it is. You can keep using the
handle. The next call returns the provider's result.

The public API covers create, lifecycle, files, and blocking command execution
inside the sandbox. GPU requests, resizing, storage volumes, browsers,
streaming process handles, sessions, PTY, and macOS or Windows creation are
outside it.

See [configuration](configuration.md), [creation](sandbox-creation.md),
[lifecycle](sandbox-lifecycle.md), [filesystem](filesystem.md), and
[process](process.md).

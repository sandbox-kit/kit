# Providers

Install the core package and the provider you use.

|            | Daytona                                                                  | Modal                                                                  | Blaxel                                                                  |
| ---------- | ------------------------------------------------------------------------ | ---------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| TypeScript | `@sandbox-kit/sandbox` and `@sandbox-kit/daytona`                        | `@sandbox-kit/sandbox` and `@sandbox-kit/modal`                        | `@sandbox-kit/sandbox` and `@sandbox-kit/blaxel`                        |
| Go         | `github.com/sandbox-kit/kit/sdks/go/sandbox` and `.../providers/daytona` | `github.com/sandbox-kit/kit/sdks/go/sandbox` and `.../providers/modal` | `github.com/sandbox-kit/kit/sdks/go/sandbox` and `.../providers/blaxel` |

## TypeScript

```ts
await using client = new Client({
  provider: modal(),
  auth: { tokenPair: { id: tokenId, secret: tokenSecret } },
  scope: { appName: "existing-app" },
});
```

```ts
await using client = new Client({
  provider: blaxel(),
  auth: { apiKey: { key: apiKey } },
  scope: { workspace: "my-workspace" },
  region: "us-pdx-1",
});
```

Use one Blaxel configuration per process. A client with different credentials
or workspace is rejected.

## Go

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

## What each provider can do

| Operation                   | Modal                              | Daytona                                    | Blaxel                                 |
| --------------------------- | ---------------------------------- | ------------------------------------------ | -------------------------------------- |
| Create                      | Registry image and an existing app | Default snapshot, named snapshot, or image | Named image create (Blaxel name rules) |
| Open an existing sandbox    | `attach` (ID only)                 | `get` (provider lookup)                    | `get` / `refresh` by name              |
| Start, stop, pause, archive | Unsupported                        | Supported                                  | Unsupported                            |
| Delete                      | Unsupported                        | Supported                                  | Supported                              |
| Terminate                   | Supported                          | Unsupported                                | Unsupported                            |
| Filesystem                  | Same eight shared methods          | Same eight shared methods                  | Same eight shared methods              |
| Process                     | Same blocking `Run`                | Same blocking `Run`                        | Same blocking `Run`                    |

`supports` / `Supports` reports whether the provider implements that operation.
A missing operation fails before the provider call. Daytona may still reject a
supported call for a particular sandbox.

### Placement

- **Daytona** — set the target with `region` on the client. Placement on the
  create request is rejected.
- **Modal** — set cloud and region on the create request. A placement target
  on the client is rejected.
- **Blaxel** — `region` on the client is required for create. Per-request
  placement is rejected.

Closing the client closes your connection to the provider. The sandboxes keep
running. Closing a handle leaves the sandbox as it is.

GPU requests, resizing, storage volumes, browsers, streaming process handles,
sessions, PTY, and macOS or Windows creation are outside this API.

See [configuration](configuration.md), [creation](sandbox-creation.md),
[lifecycle](sandbox-lifecycle.md), [filesystem](filesystem.md), and
[process](process.md).

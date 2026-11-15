# Client configuration

Pass a provider and shared settings to the client.

## TypeScript

Invalid settings throw at construction. Timeouts are milliseconds.

```ts
await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
});
```

Call `await client.close()`, or prefer `await using`. That closes your
connection to the provider. Sandboxes keep running. Closing one handle does
not close the client.

Auth is one nested shape: `auth: { apiKey: { key } }`, `tokenPair`,
`bearerToken`, or `oauthRefresh`.

Use one Blaxel configuration per process. A client with different credentials
or workspace is rejected.

A TypeScript config is a plain object, so `JSON.stringify` includes credential
fields — keep it out of logs.

## Go

`NewClient` returns `(*Client, error)`.

```go
client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
	Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
		ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
	}},
	Timeout: sandbox.Duration(2 * time.Minute),
})
```

Call `client.Close(ctx)`. That closes your connection to the provider.
Sandboxes keep running. Closing one handle does not close the client.

## Configuration fields

| Setting    | Purpose                                               |
| ---------- | ----------------------------------------------------- |
| Provider   | Daytona, Modal, or Blaxel constructor                 |
| Auth       | One credential shape. Omit to use provider env creds  |
| Scope      | App, environment, organization, project, or workspace |
| Endpoint   | Provider API URL, where supported                     |
| Region     | Provider target or default placement                  |
| Timeout    | Default deadline for calls that omit their own        |

Changing the config object after you create the client does not change the
client.

## Provider support

| Setting          | Modal                                                      | Daytona                           | Blaxel                         |
| ---------------- | ---------------------------------------------------------- | --------------------------------- | ------------------------------ |
| API key          | Unsupported                                                | Supported                         | Supported                      |
| Token pair       | Supported                                                  | Unsupported                       | Unsupported                    |
| Bearer/JWT token | Unsupported                                                | Supported; organization required  | Unsupported                    |
| OAuth refresh    | Refresh token, client ID, and one client secret or JWT key | Unsupported                       | Unsupported                    |
| Endpoint         | Unsupported                                                | API URL                           | Unsupported                    |
| Region           | Default sandbox region preference                          | Target                            | Required for create            |
| App name         | Required for app lookup                                    | Unsupported                       | Unsupported                    |
| Environment      | Client, app, and secret lookup environment                 | Unsupported                       | Unsupported                    |
| Organization ID  | Unsupported                                                | Organization                      | Unsupported                    |
| Project ID       | Unsupported                                                | Unsupported                       | Unsupported                    |
| Workspace        | Unsupported                                                | Unsupported                       | Blaxel workspace (`BL_WORKSPACE` when omitted) |

A setting the provider does not support fails before the provider client is
created. Omit auth when the provider should use credentials from the
environment.

### Daytona endpoints

A remote endpoint uses HTTPS. HTTP is accepted for `localhost` or a loopback
address. A URL with a username, password, or fragment is rejected. Daytona
uses `endpoint` when you set it, otherwise `https://app.daytona.io/api`. Kit
does not read a Daytona endpoint from the environment.

Construction does not check that credentials work. The provider checks on the
first create or open.

## Deadlines

The client timeout is the default deadline. Create and lifecycle calls can
override it. Explicit zero means no deadline from the client. A context
deadline or abort signal that fires earlier still wins.

### Daytona

Create keeps its 60-second default when the timeout is omitted or zero.
Start, stop, pause, and delete use that 60-second default when omitted, and
treat explicit zero as no operation timeout.

### Blaxel

Set a client or create timeout for how long create may wait.

**TypeScript.** If Blaxel’s edge returns HTTP 504 during create, the client
may keep waiting until the sandbox is ready or that timeout ends. A timeout
does not always mean the sandbox was never created — call `get` by name
before recreating.

**Go.** If Blaxel’s edge returns HTTP 504 during create, the call returns a
timeout as soon as that response arrives. Call `Get` by name before
recreating.

## Credentials

Keep configuration out of logs and stored files. Construction, operation, and
cleanup failures use the same error details. See
[responses and errors](responses-and-errors.md).

See [Go examples](../examples/go/README.md) and
[TypeScript examples](../examples/typescript/README.md).

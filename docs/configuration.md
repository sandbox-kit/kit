# Client configuration

Pass a provider and shared settings to the client.

Go returns `(*Client, error)`. TypeScript throws when the settings are invalid.

```go
client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
    Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
        ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
    }},
    Timeout: sandbox.Duration(2 * time.Minute),
})
```

```ts
await using client = new Client({
  provider: daytona(),
  auth: { apiKey: { key: process.env.DAYTONA_API_KEY! } },
  timeout: 120_000,
});
```

Call `client.Close(ctx)` in Go, or `await client.close()` or `await using` in
TypeScript. That closes your connection to Daytona, Modal, or Blaxel. The sandboxes keep
running. You can keep using a handle. The next call returns the provider's
result. Closing one handle does not close the client.

## Configuration fields

| Go / TypeScript         | Purpose                                                             |
| ----------------------- | ------------------------------------------------------------------- |
| Provider / `provider`   | `daytona.New()` / `daytona()`, `modal.New()` / `modal()`, or `blaxel.New()` / `blaxel()` |
| `Auth` / `auth`         | One credential shape. Omit it to use the provider's own credentials |
| `Scope` / `scope`       | App, environment, organization, project, or workspace               |
| `Endpoint` / `endpoint` | Provider API URL, where the provider supports one                   |
| `Region` / `region`     | Provider target or default placement                                |
| `Timeout` / `timeout`   | Default deadline. Go: `sandbox.Duration`. TypeScript: milliseconds  |

TypeScript auth is one nested shape: `auth: { apiKey: { key } }`, `tokenPair`,
`bearerToken`, or `oauthRefresh`.

Changing the config object after you create the client does not change the
client. Editing the create options you passed, or the metadata you read back,
does not change the handle.

## Provider support

The bindings use Modal v0.11.1, Daytona v0.225.0, `@blaxel/core@0.3.36`, and
`github.com/blaxel-ai/sdk-go@v0.31.0`.

| Setting          | Modal                                                      | Daytona                           | Blaxel                            |
| ---------------- | ---------------------------------------------------------- | --------------------------------- | --------------------------------- |
| API key          | Unsupported                                                | `APIKey`                          | `APIKey`                          |
| Token pair       | `TokenID` and `TokenSecret`                                | Unsupported                       | Unsupported                       |
| Bearer/JWT token | Unsupported                                                | `JWTToken`; organization required | Unsupported                       |
| OAuth refresh    | Refresh token, client ID, and one client secret or JWT key | Unsupported                       | Unsupported                       |
| Endpoint         | Unsupported                                                | API URL                           | Unsupported                       |
| Region           | Default sandbox region preference                          | Target                            | Required for create               |
| App name         | Required for app lookup                                    | Unsupported                       | Unsupported                       |
| Environment      | Client, app, and secret lookup environment                 | Unsupported                       | Unsupported                       |
| Organization ID  | Unsupported                                                | Organization                      | Unsupported                       |
| Project ID       | Unsupported                                                | Unsupported                       | Unsupported                       |
| Workspace        | Unsupported                                                | Unsupported                       | Required Blaxel workspace name    |

A setting the provider does not support fails before the provider client is
created. Omit `auth` when the provider should use credentials it already finds
in the environment. Daytona can take the organization from that environment
when you use a bearer token and leave organization unset.

A remote endpoint uses HTTPS. HTTP is accepted for `localhost` or a loopback
address. A URL with a username, password, or fragment is rejected. Daytona uses
`endpoint` when you set it, and otherwise uses `https://app.daytona.io/api`.
Set that URL on the client. Kit does not read a Daytona endpoint from the
environment.

A client that constructs successfully has not checked that the credentials
work. The provider checks permissions and capacity when it runs a request.

## Deadlines

`timeout` is the default deadline for provider operations. A create request
can override it with `provisioning.timeout`. A lifecycle call can override it
with that call's `timeout`. Explicit zero adds no Kit deadline. A context
deadline or abort signal that arrives earlier still applies.

Daytona create keeps its 60-second default when the timeout is omitted and
when it is explicit zero. Daytona start, stop, pause, and delete use that
60-second default when the timeout is omitted, and treat explicit zero as no
operation timeout.

## Credentials

Keep configuration out of logs and stored files. It holds credentials. A
TypeScript config is a plain object, so `JSON.stringify` includes those fields.

Construction, operation, and cleanup failures use the same error details.
The provider's own error stays available as the cause. See
[responses and errors](responses-and-errors.md).

See [Go examples](../examples/go/README.md) and
[TypeScript examples](../examples/typescript/README.md).

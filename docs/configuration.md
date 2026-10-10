# Client configuration

Construct a client with `sandbox.NewClient(sandbox.Config{...})`.
The optional provider initializes its official SDK from this shared configuration.
Applications use generated native Go types.

## Configuration fields

| Field      | Purpose                                                                |
| ---------- | ---------------------------------------------------------------------- |
| `Provider` | Integration selected with `modal.New()` or `daytona.New()`             |
| `Auth`     | One explicit credential variant; omit to use SDK credential resolution |
| `Scope`    | App, environment, organization, or project scope                       |
| `Endpoint` | Provider API URL, where supported                                      |
| `Region`   | Provider target or default placement preference                        |
| `Timeout`  | Default Kit operation deadline                                         |

A provider must implement `sandbox.Provider` and return a complete
`sandbox.Backend`. Bare SDK clients and arbitrary interfaces do not satisfy
that contract. The backend's identity must match the selected provider;
Kit closes and rejects mismatched backends.

## Provider support

Bindings target Modal Go v0.11.0 and Daytona Go v0.222.0.

| Setting          | Modal                                                      | Daytona                           |
| ---------------- | ---------------------------------------------------------- | --------------------------------- |
| API key          | Unsupported                                                | `APIKey`                          |
| Token pair       | `TokenID` and `TokenSecret`                                | Unsupported                       |
| Bearer/JWT token | Unsupported                                                | `JWTToken`; organization required |
| OAuth refresh    | Refresh token, client ID, and one client secret or JWT key | Unsupported                       |
| Endpoint         | Unsupported through the pinned public constructor          | `APIUrl`                          |
| Region           | Default sandbox `Regions` preference                       | `Target`                          |
| App name         | Retained for app lookup                                    | Unsupported                       |
| Environment      | Client, app, and secret lookup environment                 | Unsupported                       |
| Organization ID  | Unsupported                                                | `OrganizationID`                  |
| Project ID       | Unsupported                                                | Unsupported                       |

Unsupported supplied settings fail before SDK construction. Omitted auth uses
the native SDK's environment/profile resolution. Daytona may resolve organization
scope from its environment when using bearer auth.

Remote endpoints require HTTPS. HTTP is accepted only for `localhost` or a
literal loopback IP, to support local development and mock-server tests. URL
user information and fragments are rejected. Daytona resolves an omitted endpoint
from `DAYTONA_API_URL`, then `DAYTONA_SERVER_URL`, then its pinned SDK default
`https://app.daytona.io/api`; the resolved value is validated before SDK construction.
Explicit endpoints take precedence. Remote self-hosted endpoints need HTTPS.

Initialization does not prove remote credentials are valid. The provider enforces
account permissions and resource availability when it processes requests.

## Ownership and deadlines

The client owns its initialized SDK. Call `client.Close(ctx)` when finished;
this releases SDK resources without stopping or deleting sandboxes.

Generated state capture copies supplied scope and region values, preserving
`nil` and explicit values. Later mutations of the original configuration do not
change captured state. Kit also copies creation requests and returned metadata.

`Config.Timeout` applies when a creation request omits
`CreateOptions.Provisioning.Timeout`. Explicit zero adds no Kit deadline.
An earlier caller-context deadline still applies, and native SDK transport or
operation timeouts remain in force.

## Credentials and errors

Credential fields marked `sensitive` and the provider object are excluded from
generated JSON output. This does not protect arbitrary logging or formatting;
keep credentials out of logs and source control.

Construction, operation, and cleanup failures expose shared error details.
Native SDK and validator errors remain accessible through the cause chain.
See [responses and errors](responses-and-errors.md).

Configuration types and contracts come from
[`client.proto`](../proto/kit/sandbox/v1/client.proto),
[client rules](../specs/client.yaml), and provider YAML.
See [examples](../examples/go/README.md) and [generation](code-generation.md).

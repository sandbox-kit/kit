# Client configuration

Construct a Go client with `sandbox.NewClient(provider, &sandbox.Config{...})`
or a TypeScript client with `new Client({...})`. The selected provider
initializes its official SDK from generated protobuf configuration types. This
guide uses Go field spelling; TypeScript exposes the same fields in camelCase.

Provider selection is never part of the portable `Config` message. Go passes the
provider as a separate `NewClient` argument; TypeScript attaches it on
`ClientConfig` beside the protobuf config fields.

The constructors have the same ownership and validation semantics while
following each language's conventions: Go returns `(*Client, error)`, and
TypeScript throws on invalid synchronous construction. Go callers defer
`client.Close(ctx)`; TypeScript callers use `await client.close()` or scope the
client with `await using`. None of these cleanup forms changes the remote
sandbox lifecycle.

## Configuration fields

| Field      | Purpose                                                                   |
| ---------- | ------------------------------------------------------------------------- |
| `Provider` | Go: pass `modal.New()` / `daytona.New()` to `NewClient`. Not in `Config`. |
| `Auth`     | One explicit credential variant; omit to use SDK credential resolution    |
| `Scope`    | App, environment, organization, or project scope                          |
| `Endpoint` | Provider API URL, where supported                                         |
| `Region`   | Provider target or default placement preference                           |
| `Timeout`  | Default Kit operation deadline (Go `sandbox.Duration`; TS `duration(ms)`) |

A provider must implement `sandbox.Provider` and return a `sandbox.Backend`.
The backend supplies identity and cleanup, then implements typed optional
capabilities such as creation, refresh, or termination for its supported
operations. Bare SDK clients and arbitrary interfaces do not satisfy these
contracts. Both languages verify that the backend identity matches the selected
provider. Go closes a mismatched backend before returning. Because TypeScript
construction is synchronous while backend cleanup is asynchronous, TypeScript
starts cleanup before throwing and suppresses only that cleanup promise's error.

## Provider support

Go bindings target Modal v0.11.0 and Daytona v0.222.0. TypeScript bindings target
`modal@0.11.0` and `@daytona/sdk@0.224.0`.

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

The client owns its initialized SDK. Call `client.Close(ctx)` in Go or
`await client.close()` in TypeScript when finished. This releases SDK resources
without stopping or deleting sandboxes. Close is idempotent and waits for client
operations already running through that client. Do not use its handles after
client closure. Closing one sandbox handle does not close the client or change
the remote sandbox.

Generated state capture copies supplied scope and region values, preserving
`nil` and explicit values. Later mutations of the original configuration do not
change captured state. Kit also copies creation requests and returned metadata.

`Config.Timeout` is the default Kit deadline for provider operations. A creation
request can override it with `CreateOptions.Provisioning.Timeout`; lifecycle
options can override it with their operation-specific `Timeout`. Explicit zero
adds no Kit deadline. An earlier caller-context deadline still applies, and
native SDK transport or operation timeouts remain in force.

## Credentials and errors

Go credential fields marked `sensitive` and the provider object are excluded
from generated JSON output. TypeScript configuration uses plain structural
objects, so `JSON.stringify` includes enumerable credential fields. Do not
serialize, log, or persist configuration objects in either language.

Construction, operation, and cleanup failures expose shared error details.
Native SDK and validator errors remain accessible through the cause chain.
See [responses and errors](responses-and-errors.md).

Configuration types and contracts come from
[`client.proto`](../proto/kit/sandbox/v1/client.proto) (Protovalidate options;
`secure_endpoint` stays handwritten),
[client rules](../specs/client.yaml), and provider YAML.
See [examples](../examples/go/README.md) and [generation](code-generation.md).

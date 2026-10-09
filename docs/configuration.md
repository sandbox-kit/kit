# Common client initialization

`sandbox.NewClient(sandbox.Config{...})` initializes a provider SDK through the
selected optional provider. `modal.New()` and `daytona.New()` return
provider selections; SDK backends are private. Sandbox imports neither official provider SDK.

`Config` and authentication/context types are native Go output from
`proto/kit/sandbox/v1/client.proto`. `specs/client.yaml` attaches the native provider
contract and defines portable validation. Provider YAML generates the supported
credential/client field mappings, configuration assembly, and typed SDK construction.

## Common configuration

| Field | Meaning |
| --- | --- |
| `Provider` | Native `Provider`; initializes and returns a full operation backend |
| `Auth` | Exactly one explicit credential variant, or absent for SDK defaults |
| `Endpoint` | Optional provider API endpoint, HTTP/HTTPS URL |
| `Region` | Provider target or default sandbox placement |
| `Scope` | Application/environment/organization/project scope |
| `Timeout` | Default Kit operation deadline; not a replacement for SDK transport timeouts |

Sandbox validates common settings before calling the provider. A bare SDK pointer,
name-only provider, or arbitrary interface cannot be supplied to `NewClient`.
Custom providers must implement the same initialization contract and return a
backend supporting provider identity, creation, and SDK cleanup. The initialized
backend identity must match the selected provider; mismatches are closed and rejected.

## Pinned SDK mappings

| Setting | Modal Go v0.11.0 | Daytona Go v0.222.0 |
| --- | --- | --- |
| API key | Rejected | SDK APIKey |
| Token pair | SDK TokenID/TokenSecret | Rejected |
| Bearer token | Rejected | SDK JWTToken; organization required by SDK |
| OAuth refresh | SDK refresh token/client ID plus one client secret or JWT key | Rejected |
| Endpoint | Rejected: no public per-client override | SDK APIUrl |
| Region | Default creation Regions preference | SDK Target |
| App name | Stored for sandbox app lookup | Rejected |
| Environment | SDK Environment and app/secret lookup environment | Rejected |
| Organization ID | Rejected | SDK OrganizationID |
| Project ID | Rejected | Rejected |

Unsupported settings fail before SDK construction. Omitted auth invokes the
SDK's normal environment/profile resolution. Account permissions and roles are
still enforced by the provider; this contract does not create or modify them.
An initialized client is not proof that credentials have been accepted remotely.

Modal's pinned constructor exposes environment and credentials but no public
endpoint field. The adapter does not mutate process environment to emulate one.
Existing SDK environment defaults still apply. Adding another auth mode requires
provider support, mapping rules, and validation coverage.

## Ownership, deadlines, and sensitive values

The client owns its initialized SDK. Call `client.Close(ctx)` to release its
resources. Close does not stop/delete sandboxes. Native initialization, operation,
and cleanup errors pass through; local validation errors remain ordinary Go errors.

`Config.Timeout` applies when an operation does not supply its own timeout.
Explicit zero means no additional Kit deadline. Neither zero nor a longer timeout
removes an earlier caller-context deadline or changes native SDK transport limits.
The default timeout is snapshotted during construction, and creation requests are
copied before defaults are applied. App/environment/region values are captured by
the adapter, so later caller mutations do not change its scope.

Credential fields marked `sensitive` are excluded from generated Go JSON output.
Provider objects are also excluded. This is not a persisted credential/config-file
format; YAML specs contain rules and bindings, never application credential values.
Avoid printing credentials through other formatting mechanisms.

See [the runnable example](../examples/go/README.md) and
[generator specification](code-generation.md).

# Sandbox creation

Create returns a sandbox handle. Creation settings are separate from lifecycle
calls on an existing sandbox.

## TypeScript

```ts
const sandbox = await client.create(options, { timeout: 60_000 });
```

Create options use nested exclusive fields (`source: { image: … }`,
`readiness: { tcpPort }` or `{ command }`), plain `environment` / `labels`
maps, and millisecond durations. A missing property means the field was
omitted.

## Go

```go
instance, err := client.Create(ctx, options)
```

Use `sandbox.Value(v)` for optional scalars (including zero or false) and
`sandbox.Duration(d)` for timeouts.

## Supported creation settings

| Setting                | Modal                                                                                  | Daytona                                                      | Blaxel                                                                          |
| ---------------------- | -------------------------------------------------------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------------------------- |
| Source                 | Explicit registry image and an existing app                                            | Default/specified snapshot or registry image                 | Named create (Blaxel name rules); registry image required                       |
| Runtime                | Entrypoint, absolute working directory, PTY                                            | Toolbox language and user                                    | Unsupported                                                                     |
| Isolation              | Default, gVisor container, Linux VM                                                    | Nested KVM flag; class inherited from source                 | Unsupported                                                                     |
| Resources              | Fractional CPU, CPU limit, memory request/limit                                        | Whole CPU and GiB memory/disk; image required for overrides  | Memory in MiB; CPU/disk overrides not supported                                 |
| Placement              | Cloud and region preferences                                                           | Client target; request placement unsupported                 | Client `region` required; request placement unsupported                         |
| Environment and labels | Environment map and tags                                                               | Environment map and labels                                   | Environment map and labels                                                      |
| Security               | Named secrets and workload identity                                                    | Egress-placeholder secrets and public access                 | Unsupported                                                                     |
| Network                | Egress allowlists, inbound CIDRs, private network, managed proxy, ports, custom domain | Egress allowlists, outbound proxy, linked sandbox            | Ports (not 80/443/8080; default → HTTP); outbound domains                       |
| Lifetime               | Maximum lifetime and idle termination                                                  | TTL, ephemeral, stop/pause/archive/delete policies           | Maximum lifetime (≥5m) and idle terminate (≥1m); idle disable not supported     |
| Observability          | Verbose output                                                                         | Telemetry endpoint                                           | Unsupported                                                                     |

Volumes are not available on create for any provider.

Unsupported fields fail before the provider call. The error names the field
(for example `resources.disk_mib`).

### Wait behavior

**Modal.** Default or scheduled wait. Ready wait needs a probe and a positive
create timeout.

**Daytona (TypeScript).** Default or started wait. Submitted wait is not
supported — create returns after startup.

**Daytona (Go).** Default, submitted, or started wait.

**Blaxel.** Ready, Started, and Submitted wait are not supported. Set a create
timeout; with 1–50 seconds, a timed-out create frees the name.

**Blaxel TypeScript.** If the edge returns HTTP 504 during create, the client
may keep waiting until the sandbox is ready or the create timeout ends. Call
`get` by name before recreating.

**Blaxel Go.** If the edge returns HTTP 504 during create, the call returns a
timeout as soon as that response arrives. Call `Get` by name before recreating.

## Units and limits

Omitted fields keep the provider default. Invalid combinations fail in the
client before the provider call.

| Resource         | Modal                                           | Daytona                                      | Blaxel                                      |
| ---------------- | ----------------------------------------------- | -------------------------------------------- | ------------------------------------------- |
| CPU              | Physical cores; minimum 0.125, precision 0.001  | Whole CPU cores                              | Derived by Blaxel from memory; overrides rejected |
| Memory           | MiB; minimum request 128                        | Whole GiB expressed as MiB                   | Pass memory in MiB; omit for Blaxel default |
| Disk             | Creation override unsupported                   | Whole GiB expressed as MiB                   | Creation override unsupported               |
| Maximum lifetime | Positive whole seconds, up to 24 hours (TS: ms) | Whole-minute TTL (TS: ms); zero disables TTL | At least 5 minutes when set                 |

Account quotas, image existence, regions, and capacity remain
provider-enforced. Resources and labels on the handle are what the provider
reported.

## Policy semantics

| Mode     | Meaning                                      |
| -------- | -------------------------------------------- |
| Default  | Defer to the provider                        |
| Disabled | Explicitly disable that policy               |
| After    | Apply after a delay (zero = immediate action)|

TypeScript delays are milliseconds.

### Daytona

| Action          | Disabled    | After zero  | Positive delay                 |
| --------------- | ----------- | ----------- | ------------------------------ |
| Idle stop       | Native zero | Rejected    | Whole minutes                  |
| Idle pause      | Native zero | Rejected    | Whole minutes                  |
| Stopped archive | Rejected    | Rejected    | Whole minutes, at most 30 days |
| Stopped delete  | Rejected    | Native zero | Whole minutes                  |

Positive stop and pause delays are mutually exclusive. Ephemeral cannot combine
with positive pause, archive, or delayed-delete. Linked sandboxes must be
ephemeral.

### Modal

Idle termination supports disabled or positive whole-second delays. Maximum
lifetime and idle delay are separate settings.

### Blaxel

Idle terminate supports only After with a delay of at least 1m, or provider
default. Disabled idle terminate is rejected. Maximum lifetime, when set, must
be at least 5m. Ephemeral, idle stop/pause, stopped archive/delete, runtime,
and isolation are rejected.

## Networking

A present empty allowlist means deny-all. Empty inbound Modal allowlists and
empty Daytona allowlists are rejected. Secrets are references to
provider-managed secrets, not authentication credentials.

## Results and cleanup

`info()` / `Info()` returns a copy. Editing that copy leaves the handle
unchanged. Closing the client or handle does not delete the sandbox.

Lifecycle: [lifecycle guide](sandbox-lifecycle.md). Files:
[filesystem](filesystem.md). Commands: [process](process.md).

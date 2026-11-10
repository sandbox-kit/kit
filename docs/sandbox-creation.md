# Sandbox creation

Go's `Client.Create(ctx, options)` and TypeScript's `client.create(options)`
return a sandbox handle with copied metadata. Creation settings are distinct
from lifecycle operations on an existing sandbox.

## Supported creation settings

Go and TypeScript integrations target Modal v0.11.1, Daytona v0.225.0, and
Blaxel (`modal@0.11.1`, `@daytona/sdk@0.225.0`, `@blaxel/core@0.3.36` /
`github.com/blaxel-ai/sdk-go@v0.31.0`).

| Setting                | Modal                                                                                  | Daytona                                                            | Blaxel                                                                                         |
| ---------------------- | -------------------------------------------------------------------------------------- | ------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------- |
| Source                 | Explicit registry image and an existing app                                            | Default/specified snapshot or registry image                       | Named create (Blaxel name rules); registry image required                                      |
| Runtime                | Entrypoint, absolute working directory, PTY                                            | Toolbox language and user                                          | Unsupported                                                                                    |
| Isolation              | Default, gVisor container, Linux VM                                                    | Nested KVM flag; class inherited from source                       | Unsupported                                                                                    |
| Resources              | Fractional CPU, CPU limit, memory request/limit                                        | Whole CPU and GiB memory/disk; image source required for overrides | Memory 1:1 (`memory_mib`); CPU/disk overrides rejected                                         |
| Placement              | Cloud and region preferences                                                           | Client target; request placement unsupported                       | Client `region` required for create; request placement unsupported                             |
| Environment and labels | Environment map and tags                                                               | Environment map and labels                                         | Environment map and labels                                                                     |
| Security               | Named secrets and workload identity                                                    | Egress-placeholder secrets and public access                       | Unsupported                                                                                    |
| Network                | Egress allowlists, inbound CIDRs, private network, managed proxy, ports, custom domain | Egress allowlists, outbound proxy, linked sandbox                  | Ports (not 80/443/8080; Default → Blaxel HTTP); outbound domains                               |
| Lifetime               | Maximum lifetime and idle termination                                                  | TTL, ephemeral behavior, stop/pause/archive/delete policies        | Maximum lifetime (≥5m) and idle terminate (≥1m); explicit idle disable rejected                |
| Wait                   | Default/scheduled; ready requires a probe and positive timeout                         | Go: default, submitted, or started. TypeScript: default or started | Kit client/create timeout only (native create timeout / SandboxCreateOptions not forwarded); waitFor Default/Scheduled only |
| Observability          | Verbose output                                                                         | Telemetry endpoint                                                 | Unsupported                                                                                    |

Kit has no portable volume fields; Blaxel volumes stay out of create the same
way Daytona and Modal volumes do.

Modal default, snapshot, and warm-pool sources are unavailable. Daytona warm
pools and explicit snapshot kinds are unavailable. A field the provider does
not support is rejected before the provider call. The error names that field,
such as `resources.disk_mib`.

The Daytona TypeScript SDK returns only after startup, so a submitted wait is
rejected there. Go can return once the sandbox is submitted.

## Presence, units, and limits

Omitted fields keep the provider default. In Go, `sandbox.Value(v)` marks an
optional scalar as set, including zero or false, and `sandbox.Duration(d)`
sets a duration.

TypeScript create options use nested exclusive fields (`source: { image: … }`,
`readiness: { tcpPort }` or `{ command }`), `string[]` allowlists, plain
`environment` and `labels` maps, and millisecond durations (`timeout`,
`after`, `maximumLifetime`). A missing property, or `undefined`, means the
field was omitted.

Invalid combinations fail in the client before the provider call. A field the
selected provider does not support fails there too.

| Resource         | Modal                                           | Daytona                                      | Blaxel                                                           |
| ---------------- | ----------------------------------------------- | -------------------------------------------- | ---------------------------------------------------------------- |
| CPU              | Physical cores; minimum 0.125, precision 0.001  | Whole CPU cores                              | Derived by Blaxel from memory; Kit rejects CPU overrides         |
| Memory           | MiB; minimum request 128                        | Whole GiB expressed as MiB                   | `memory_mib` / `memoryMib` passed 1:1; Blaxel labels the unit MB |
| Disk             | Creation override unsupported                   | Whole GiB expressed as MiB                   | Creation override unsupported                                    |
| Maximum lifetime | Positive whole seconds, up to 24 hours (TS: ms) | Whole-minute TTL (TS: ms); zero disables TTL | Blaxel TTL string (at least 5m when set)                         |

Daytona request conversion divides MiB exactly by 1024. Allocated metadata
converts GiB back to MiB. Snapshot creation inherits snapshot resources.
Positive CPU or memory limits require a corresponding request and cannot be
lower than that request. Blaxel memory is a 1:1 integer pass-through; Kit
does not convert MiB to MB.

Local checks cover syntax, known enums, documented bounds, precision, and
supported combinations. Account quotas, class-specific permissions, image
existence, regions, and capacity remain provider-enforced. Metadata availability
varies; resource metadata represents SDK-reported allocation, not echoed requests.

## Policy semantics

Default mode defers to the provider (`PolicyModeDefault` / `PolicyMode.Default`).
Disabled mode requests explicit disabling (`PolicyModeDisabled` /
`PolicyMode.Disabled`). After mode supplies a delay (`PolicyModeAfter` /
`PolicyMode.After`); a zero delay requests an immediate action and must not be
interpreted as disabling. TypeScript delays are millisecond numbers.

Daytona applies these policy values:

| Action          | Disabled    | AFTER zero  | Positive delay                 |
| --------------- | ----------- | ----------- | ------------------------------ |
| Idle stop       | Native zero | Rejected    | Whole minutes                  |
| Idle pause      | Native zero | Rejected    | Whole minutes                  |
| Stopped archive | Rejected    | Rejected    | Whole minutes, at most 30 days |
| Stopped delete  | Rejected    | Native zero | Whole minutes                  |

A zero archive delay is rejected. Daytona treats that zero as 30 days, so it
is not an immediate archive. A negative interval, the value Daytona uses to
disable deletion, is rejected at creation.

Positive stop and pause delays are mutually exclusive. Ephemeral requests cannot
combine with positive pause, archive, or delayed-delete policies. Positive pause
also conflicts with immediate deletion.

Daytona linked sandboxes must be ephemeral. Kit rejects a linked-sandbox request
before provider calls unless `lifetime.ephemeral` is explicitly true.

Modal idle termination supports disabled or positive whole-second delays.
Its maximum lifetime and idle delay remain separate settings.

Blaxel idle terminate supports only `PolicyModeAfter` with a delay of at least
1m, or provider default. `PolicyModeDisabled` on idle terminate is rejected.
Maximum lifetime, when set, must be at least 5m.

Image references support lowercase repository paths, optional registry hostnames
and ports, tags up to 128 characters, and SHA-256 digests, within 512 characters.
IPv6 registry literals and other digest algorithms are unsupported.

## Networking and waiting

A present empty allowlist means deny-all. Modal outbound wrappers preserve this.
Empty inbound Modal allowlists and empty Daytona allowlists are rejected because
their native meanings do not match that intent.

Mapped network values receive CIDR/domain checks. Modal validates environment
keys; Daytona validates secret environment-variable names. Cloud names and
toolbox languages use declared supported values. Secrets are references to
provider-managed secrets, not authentication credentials.

Creation timeout, queue timeout, and lifetime are distinct. Omitted creation
timeout inherits `Config.Timeout` / `timeout`; explicit zero adds no Kit
deadline. A caller-context deadline still applies. A Daytona creation timeout
must be a whole number of seconds. An omitted Daytona creation timeout and an
explicit zero both leave Daytona's 60-second create default in place. Modal
ready waiting also requires a readiness probe and a positive creation timeout.

## Results and cleanup

`Info()` returns a copy of the metadata. Editing that copy leaves the handle
unchanged. Errors include a kind, field, and message. The provider's own error
stays available as the cause. Go forwards cancellation into the provider call.
TypeScript rejects the `await` on abort; create itself is not queued behind
other client calls. Kit does not retry, and it does not delete a sandbox when a
call fails.

Closing the client closes your connection to the provider. The sandbox keeps
running. Closing a handle leaves the sandbox as it is. You can keep using the
handle. The next call returns the provider's result.
Lifecycle operations are described in the [lifecycle guide](sandbox-lifecycle.md).
Files inside a sandbox are described in the [filesystem guide](filesystem.md).

Creation covers the sandbox itself. GPU requests, resizing, storage volumes
(including Blaxel volumes), browsers, streaming process handles, sessions,
PTY, and macOS or Windows creation are outside it. Blocking command execution
is covered in [process](process.md).

See [configuration](configuration.md), [lifecycle](sandbox-lifecycle.md), and
[filesystem](filesystem.md).

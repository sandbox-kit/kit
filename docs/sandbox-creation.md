# Sandbox creation

`Client.Create(ctx, options)` returns a `sandbox.Sandbox` with `ID()`,
`ProviderName()`, and copied `Info()` metadata. Creation settings are distinct
from future operations on an existing sandbox.

## Supported creation settings

The integrations target Modal Go v0.11.0 and Daytona Go v0.222.0.

| Setting                | Modal                                                                   | Daytona                                                            |
| ---------------------- | ----------------------------------------------------------------------- | ------------------------------------------------------------------ |
| Source                 | Explicit registry image and an existing app                             | Default/specified snapshot or registry image                       |
| Runtime                | Entrypoint, absolute working directory, PTY                             | Toolbox language and user                                          |
| Isolation              | Default, gVisor container, Linux VM                                     | Nested KVM flag; class inherited from source                       |
| Resources              | Fractional CPU, CPU limit, memory request/limit                         | Whole CPU and GiB memory/disk; image source required for overrides |
| Placement              | Cloud and region preferences                                            | Client target; request placement unsupported                       |
| Environment and labels | Environment map and tags                                                | Environment map and labels                                         |
| Security               | Named secrets and workload identity                                     | Egress-placeholder secrets and public access                       |
| Network                | Egress allowlists, inbound CIDRs, private network, ports, custom domain | Egress allowlists, outbound proxy, linked sandbox                  |
| Lifetime               | Maximum lifetime and idle termination                                   | TTL, ephemeral behavior, stop/pause/archive/delete policies        |
| Wait                   | Default/scheduled; ready requires a probe and positive timeout          | Default/submitted/started                                          |
| Observability          | Verbose output                                                          | Telemetry endpoint                                                 |

Modal default/snapshot/warm-pool sources are not mapped. Daytona warm pools and
explicit snapshot kinds are not mapped. Supplied unsupported fields and
`ProviderOptions` are rejected before native creation calls.

## Presence, units, and limits

Omitted fields retain provider defaults. `sandbox.Value(v)` marks an optional
value as supplied, including zero or false; validation determines whether that
explicit value can be honored.

| Resource         | Modal                                          | Daytona                             |
| ---------------- | ---------------------------------------------- | ----------------------------------- |
| CPU              | Physical cores; minimum 0.125, precision 0.001 | Whole CPU cores                     |
| Memory           | MiB; minimum request 128                       | Whole GiB expressed as MiB          |
| Disk             | Creation override unsupported                  | Whole GiB expressed as MiB          |
| Maximum lifetime | Positive whole seconds, up to 24 hours         | Whole-minute TTL; zero disables TTL |

Daytona request conversion divides MiB exactly by 1024. Allocated metadata
converts GiB back to MiB. Snapshot creation inherits snapshot resources.
Positive CPU or memory limits require a corresponding request and cannot be
lower than that request.

Local checks cover syntax, known enums, documented bounds, precision, and
supported combinations. Account quotas, class-specific permissions, image
existence, regions, and capacity remain provider-enforced. Metadata availability
varies; resource metadata represents SDK-reported allocation, not echoed requests.

## Policy semantics

`PolicyModeDefault` defers to the provider. `PolicyModeDisabled` requests explicit
disabling. `PolicyModeAfter` supplies a delay; a zero delay requests an immediate
action and must not be interpreted as disabling.

Current Daytona mappings follow the service guide where SDK comments conflict:

| Action          | Disabled    | AFTER zero  | Positive delay                 |
| --------------- | ----------- | ----------- | ------------------------------ |
| Idle stop       | Native zero | Rejected    | Whole minutes                  |
| Idle pause      | Native zero | Rejected    | Whole minutes                  |
| Stopped archive | Rejected    | Rejected    | Whole minutes, at most 30 days |
| Stopped delete  | Rejected    | Native zero | Whole minutes                  |

Daytona archive zero selects 30 days rather than immediate archival. The service
uses -1 to disable deletion, but the pinned creation SDK rejects negative
intervals, so that request is unsupported at creation.

Positive stop and pause delays are mutually exclusive. Ephemeral requests cannot
combine with positive pause, archive, or delayed-delete policies. Positive pause
also conflicts with immediate deletion.

Modal idle termination supports disabled or positive whole-second delays.
Its maximum lifetime and idle delay remain separate settings.

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
timeout inherits `Config.Timeout`; explicit zero adds no Kit deadline.
Caller deadlines and native SDK timeouts still apply. Modal ready waiting also
requires a readiness probe and a positive creation timeout.

## Results, cleanup, and future work

Kit uses generated typed copies for requests and result metadata. Errors expose
shared details while preserving native causes. Cancellation is forwarded to the
SDK. Kit performs no automatic retries or cloud cleanup.

Closing the client releases SDK resources. Use provider tools for sandbox
cleanup. Execution, lifecycle methods, GPU requests, resizing, storage, harnesses,
browsers, computer use, and macOS/Windows support remain future work.

See [configuration](configuration.md), [provider reference verification](provider-verification.md),
and [generation rules](code-generation.md).

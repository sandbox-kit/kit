# Sandbox creation contract

The shared protobuf schema generates native Go `sandbox.CreateOptions` and response
structs, plus a local `Client.Create` method. It does not generate network RPCs.
Creation returns a `Sandbox` handle with `ID()`, `ProviderName()`, and copied `Info()`.
Lifecycle methods are not implemented in this increment.

## Implemented provider mappings

Compatibility is pinned to Modal Go v0.11.0 and Daytona Go v0.222.0. Refer to the
[Modal sandbox reference](https://modal.com/docs/sdk/go/latest/Sandbox),
[Daytona client reference](https://www.daytona.io/docs/en/go-sdk/daytona/), and
[Daytona types](https://www.daytona.io/docs/en/go-sdk/types/) when upgrading.

| Group | Modal | Daytona |
| --- | --- | --- |
| Source | Explicit registry image; existing app name required | Default/specified snapshot, or registry image |
| Runtime | Entrypoint, working directory, PTY | Toolbox language and user |
| Isolation | Default, gVisor container, Linux VM | Nested KVM flag; class inherited from source |
| Resources | Fractional CPU, CPU limit, MiB memory and memory limit | Whole CPU and whole GiB memory/disk, image source only |
| Placement | Cloud and region preferences | Target selected on initialized SDK; per-request placement rejected |
| Environment/labels | Env map and tags | Env map and labels |
| Secrets/security | Named secrets, workload identity | Egress-placeholder secrets, public access |
| Network | CIDR/domain egress, inbound CIDRs, private network, explicit port transport, custom domain | CIDR/domain egress, outbound proxy URL, linked sandbox |
| Lifetime | Positive whole-second maximum lifetime and idle termination | Whole-minute TTL, ephemeral, auto-stop/pause/archive/delete |
| Wait | Default/scheduled; ready requires probe and positive creation timeout | Default/submitted/started |
| Observability | Verbose | Telemetry endpoint |

Unsupported options are rejected before SDK calls. Default/snapshot/warm-pool
creation is not implemented for Modal; warm pools and explicit snapshot kinds
are rejected for Daytona. Provider-options extensions are rejected until mapped.

All shared validation is generated from YAML, along with common field copies,
resource mappings, and Daytona metadata conversions. Other semantics and typed
SDK call sequences remain explicit in adapters. See [generation specifications](code-generation.md).

## Presence and units

Absent optional values defer to the provider. Explicit zero/false is preserved.
Source alternatives are mutually exclusive. Runtime language is a hint, not a
promise to install a language, and does not determine isolation.

Shared memory/disk units are MiB. Daytona requests divide exactly by 1024; its
pinned sandbox metadata documents memory/disk in GiB, so response mapping
multiplies by 1024. Fractional GiB and out-of-range allocations are rejected.
Reservations and limits differ; unsupported limits are not dropped. Allocated
resources are mapped from SDK metadata, not echoed from the request.

Creation timeout, queue timeout, and sandbox lifetime are separate. Omitted creation timeout inherits `Config.Timeout`. Zero adds no Kit
deadline; it does not remove the caller's context deadline.
Policy DEFAULT defers, DISABLED disables, and AFTER supplies a delay. Zero AFTER
means immediate; adapters reject it where native numeric zero means disabled.

Present empty allowlists mean deny-all. Modal outbound wrappers preserve this;
unsupported inbound/daytona empty-list mappings are rejected. Secrets contain
references, not credentials. AuthConfig is configured through `Config.Auth`, while account permissions
are enforced by the provider. Modal app/environment now lives in
`Config.Scope`; it is not a per-sandbox creation field. See
[client configuration](configuration.md).

## Ownership and scope

The client owns its initialized SDK and clones declarative requests/returned metadata. Provider errors pass
through unchanged; local validation errors retain validator details through wrapping.
Cancellation is checked before delegation and forwarded to the provider SDK.
There is no extra retry layer or automatic cleanup of created cloud resources.

GPU, resizing, storage/volumes, harnesses, browsers, computer use, macOS/Windows,
and lifecycle operations remain future extensions. Tests use SDK service doubles
and a fake HTTP transport, with no live provisioning.

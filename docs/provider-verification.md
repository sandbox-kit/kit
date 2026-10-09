# Provider reference verification

The Go integrations are pinned to Modal v0.11.0 and Daytona v0.222.0. Latest
documentation is used to discover APIs; pinned source and compilation determine
the signatures and serialization used by this checkout. Verification covers
initialization and sandbox creation, not all features or live provisioning.

## References checked

- Modal: [Client](https://modal.com/docs/sdk/go/latest/Client),
  [App](https://modal.com/docs/sdk/go/latest/App),
  [Image](https://modal.com/docs/sdk/go/latest/Image),
  [Sandbox](https://modal.com/docs/sdk/go/latest/Sandbox), and
  [Errors](https://modal.com/docs/sdk/go/latest/Errors).
  Creation dependencies were also checked through
  [Secret](https://modal.com/docs/sdk/go/latest/Secret),
  [Probe](https://modal.com/docs/sdk/go/latest/Probe),
  [Allowlist](https://modal.com/docs/sdk/go/latest/Allowlist), and
  [SandboxRuntime](https://modal.com/docs/sdk/go/latest/SandboxRuntime).
- Daytona: [client and sandbox](https://www.daytona.io/docs/en/go-sdk/daytona/),
  [types](https://www.daytona.io/docs/en/go-sdk/types/),
  [options](https://www.daytona.io/docs/en/go-sdk/options/), and
  [errors](https://www.daytona.io/docs/en/go-sdk/errors/).

Modal uses `NewClientWithOptions(*ClientParams)`. Creation looks up an existing
app and registry image before `Sandboxes.Create`. Credentials, environment,
creation fields, and readiness signatures were compared with the reference and
pinned source. The pinned SDK serializes CPU as uint32 mill CPUs and lifetime/idle
delays as uint32 seconds; generated mappings reject precision loss and overflow.

Daytona uses `NewClientWithConfig(*types.DaytonaConfig)`. Its `Client.Create`
accepts snapshot/image params and variadic creation options. `WithTimeout` and
`WithWaitForStart` retain SDK behavior. SDK errors pass through, preserving native
classification and fields; common error classification remains future work.

## Policy correction

The Go parameter comments conflict with the
[sandbox guide](https://www.daytona.io/docs/en/sandboxes/). The guide and sandbox
object documentation say zero disables stop/pause; the guide says archive zero
means 30 days, not immediate archive. Generated mappings follow these rules:
disabled stop/pause maps to zero; immediate deletion maps to zero. Immediate
stop/pause/archive and explicit disabling of archive/delete are rejected.
Auto-delete disabling uses -1 in the service guide, but the pinned creation SDK
rejects negative intervals, so it cannot be exposed faithfully at creation.
Archive delays are bounded to 30 days. Omitted/default policies stay omitted.

## Static validation and provider defaults

Provider `checks` are generated from documented creation rules. Modal enforces
its minimum 0.125-core/128-MiB request and 24-hour maximum lifetime, supported
cloud names, absolute workdirs, and environment/network syntax. Daytona checks
its three supported toolbox languages, network syntax, mutually exclusive idle
actions, and ephemeral policy conflicts. Common enum values are validated against
schema descriptors, and positive resource limits require a corresponding request.

Modal defaults are 0.125 physical cores, 128 MiB, and a five-minute lifetime.
Daytona image creation has documented defaults of 1 vCPU, 1 GiB, and 3 GiB disk;
snapshot creation inherits snapshot resources. These defaults are not injected
by Kit. See [Modal resources](https://modal.com/docs/guide/resources),
[Modal sandboxes](https://modal.com/docs/guide/sandboxes), and
[Daytona resources](https://www.daytona.io/docs/en/sandboxes/).

Per-account quotas, custom regions, available hardware, image/snapshot existence,
permissions, and class-specific features need provider-side validation. Standard
Daytona limits (4 vCPUs, 8 GiB, 10 GiB) are organization limits that can increase,
not global SDK type limits. They are not hardcoded. Local acceptance is not a
promise that a provider can allocate a request for a particular account.

## Documentation ambiguity

The Daytona client reference includes an image example using `Memory: 4096`
without stating units there. Sandbox field documentation explicitly states GiB.
That example alone is insufficient to change the existing memory conversion.
Request units should be checked against the API contract when revisiting it.

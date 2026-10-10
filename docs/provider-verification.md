# Provider reference verification

Bindings target Modal Go v0.11.0 and Daytona Go v0.222.0. Latest documentation
helps discover APIs; pinned SDK source and compilation determine the signatures
and serialization used by this checkout.

This review covers authentication and mapped creation settings. Local tests
verify conversions, validation, SDK request serialization, and ownership.
They do not prove live scheduling, account capacity, or lifecycle behavior.

## Reference sources

| Provider | API references                                                                                                                                                                                                                                                              | Service guides                                                                                                                                                                                                                          |
| -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Modal    | [Client](https://modal.com/docs/sdk/go/latest/Client), [App](https://modal.com/docs/sdk/go/latest/App), [Image](https://modal.com/docs/sdk/go/latest/Image), [Sandbox](https://modal.com/docs/sdk/go/latest/Sandbox), [Errors](https://modal.com/docs/sdk/go/latest/Errors) | [Resources](https://modal.com/docs/guide/resources), [Sandbox resources](https://modal.com/docs/guide/sandbox-resources), [Sandboxes](https://modal.com/docs/guide/sandboxes), [Regions](https://modal.com/docs/guide/region-selection) |
| Daytona  | [Client/sandbox](https://www.daytona.io/docs/en/go-sdk/daytona/), [Types](https://www.daytona.io/docs/en/go-sdk/types/), [Options](https://www.daytona.io/docs/en/go-sdk/options/), [Errors](https://www.daytona.io/docs/en/go-sdk/errors/)                                 | [Sandboxes](https://www.daytona.io/docs/en/sandboxes/), [Limits](https://www.daytona.io/docs/limits)                                                                                                                                    |

Modal creation dependencies also use the [Secret](https://modal.com/docs/sdk/go/latest/Secret),
[Probe](https://modal.com/docs/sdk/go/latest/Probe),
[Allowlist](https://modal.com/docs/sdk/go/latest/Allowlist), and
[SandboxRuntime](https://modal.com/docs/sdk/go/latest/SandboxRuntime) references.

## Binding decisions

Modal uses `NewClientWithOptions(*ClientParams)`. Creation looks up an existing
app and a registry image before `Sandboxes.Create`. The pinned SDK serializes
CPU as uint32 mill CPUs and lifetime/idle delays as uint32 seconds. Generated
mappings prevent precision loss and overflow; provider checks enforce documented
minimum requests and the 24-hour lifetime limit.

Daytona uses `NewClientWithConfig(*types.DaytonaConfig)`. Its `Create` method
accepts snapshot/image parameters plus creation options. Timeout and wait options
retain native behavior. Generated errors expose common categories while preserving
native classification through the original cause.

## Conflicting Daytona policy documentation

The `SandboxBaseParams` comments label zero as immediate stop/archive/delete.
The service guide and sandbox-object documentation disagree for stop/archive:

* Zero disables auto-stop and auto-pause.
* Archive zero selects 30 days; it does not request immediate archival.
* Delete zero requests deletion immediately after stopping.
* Delete disabling uses -1 in the service API, which the pinned creation SDK rejects.

The generated policies follow these verified service semantics. Explicit
archive/delete disabling and immediate stop/pause/archive are rejected.
See the [policy matrix](sandbox-creation.md#policy-semantics).

## Defaults and dynamic limits

| Provider                  | Documented resource defaults         | Lifetime behavior            |
| ------------------------- | ------------------------------------ | ---------------------------- |
| Modal                     | 0.125 physical CPU cores and 128 MiB | Five-minute maximum lifetime |
| Daytona image creation    | 1 vCPU, 1 GiB memory, 3 GiB disk     | Provider/class defaults      |
| Daytona snapshot creation | Snapshot resource allocation         | Provider/class defaults      |

Kit leaves omitted fields to the provider. It does not inject these documented
defaults into requests.

Standard Daytona organization limits can increase. Account quotas, custom
regions, available hardware, image/snapshot existence, permissions, and
class-specific features require provider-side validation. Local acceptance
does not guarantee allocation.

The Daytona client reference contains an image example with `Memory: 4096`
without specifying units. The sandbox reference states GiB. The integration
uses the documented API units; that example alone is insufficient evidence to
change the conversion.

## Updating a binding

1. Check the detailed method/type references and service guides.
2. Compare them with the exact SDK version in the provider's `go.mod`.
3. Record conflicting or account-dependent behavior.
4. Update specs, generate output, and add boundary/serialization tests.
5. Verify the relevant examples and document what local tests cannot establish.

See [generation](code-generation.md) and [example verification](../examples/go/README.md).

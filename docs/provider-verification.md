# Provider reference verification

Go and TypeScript bindings target Modal v0.11.1 and Daytona v0.225.0
(`modal@0.11.1`, `@daytona/sdk@0.225.0`). Latest documentation helps discover
APIs; pinned SDK source and compilation determine the signatures and
serialization used by this checkout.

This review covers authentication, mapped creation settings, and lifecycle
bindings. Local tests verify conversions, validation, SDK request serialization,
and ownership.
They do not prove live scheduling or account capacity. Lifecycle contract tests
cover capability dispatch, metadata updates, class checks, attachment semantics,
provider-class capability resolution, and error translation without provisioning
cloud resources. Provider tests use the pinned SDK's transport/service seams;
they do not expose those native clients through the public SDK.

## Reference sources

| Provider | API references                                                                                                                                                                                  | Service guides                                                                                                                                                                                                                          |
| -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Modal    | [Go client](https://modal.com/docs/sdk/go/latest/Client), [Go sandbox](https://modal.com/docs/sdk/go/latest/Sandbox), [JavaScript/Go SDK guide](https://modal.com/docs/guide/sdk-javascript-go) | [Resources](https://modal.com/docs/guide/resources), [Sandbox resources](https://modal.com/docs/guide/sandbox-resources), [Sandboxes](https://modal.com/docs/guide/sandboxes), [Regions](https://modal.com/docs/guide/region-selection) |
| Daytona  | [Go client/sandbox](https://www.daytona.io/docs/en/go-sdk/daytona/), [Go types](https://www.daytona.io/docs/en/go-sdk/types/), [TypeScript SDK](https://www.daytona.io/docs/en/typescript-sdk/) | [Sandboxes](https://www.daytona.io/docs/en/sandboxes/), [Limits](https://www.daytona.io/docs/limits)                                                                                                                                    |

Modal creation dependencies also use the [Secret](https://modal.com/docs/sdk/go/latest/Secret),
[Probe](https://modal.com/docs/sdk/go/latest/Probe),
[Allowlist](https://modal.com/docs/sdk/go/latest/Allowlist), and
[SandboxRuntime](https://modal.com/docs/sdk/go/latest/SandboxRuntime) references.

## Binding decisions

Modal uses `NewClientWithOptions(*ClientParams)`. Creation looks up an existing
app and a registry image before `Sandboxes.Create`. The pinned SDK serializes
CPU as uint32 mill CPUs and lifetime/idle delays as uint32 seconds. Handwritten
mappings prevent precision loss and overflow; provider checks enforce documented
minimum requests and the 24-hour lifetime limit.

Daytona uses `NewClientWithConfig(*types.DaytonaConfig)`. Its `Create` method
accepts snapshot/image parameters plus creation options. Timeout and wait options
retain native behavior. Handwritten error classification exposes shared
`ErrorKind` values while preserving native classification through the original
cause.

The TypeScript adapters construct `ModalClient` and `Daytona` directly. Modal
looks up the app, image, secrets, and attached sandboxes through the pinned
client services. Daytona uses its asynchronous create/get and lifecycle methods.
Handwritten mappings and error classification preserve the same portable
semantics in both languages.

## Conflicting Daytona policy documentation

The `SandboxBaseParams` comments label zero as immediate stop/archive/delete.
The service guide and sandbox-object documentation disagree for stop/archive:

- Zero disables auto-stop and auto-pause.
- Archive zero selects 30 days; it does not request immediate archival.
- Delete zero requests deletion immediately after stopping.
- Delete disabling uses -1 in the service API, which the pinned creation SDK rejects.

Handwritten Daytona policy mappings follow these verified service semantics.
Explicit archive/delete disabling and immediate stop/pause/archive are rejected.
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
2. Compare them with the pinned SDK version (`go.mod` or `package.json`).
3. Record conflicting or account-dependent behavior.
4. Update specs and handwritten checks, mappings, and classification; add
   boundary/serialization tests.
5. Verify the relevant examples and document what local tests cannot establish.

See [generation](code-generation.md),
[Go examples](../examples/go/README.md), and
[TypeScript examples](../examples/typescript/README.md).

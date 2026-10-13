# Sandbox lifecycle

`CreateOptions` describes creation intent and automatic policies. Lifecycle
methods operate on the resolved sandbox returned by `Client.Create`, a
provider-verified `Client.Get`, or an explicitly unverified `Client.Attach`;
they do not reuse or mutate the creation request.

```go
instance, err := client.Create(ctx, options)
if instance.Supports(sandbox.LifecycleCapabilityStop) {
    err = instance.Stop(ctx, &sandbox.StopOptions{Timeout: sandbox.Value(time.Minute)})
}
err = instance.Start(ctx, nil)
```

An operation-specific timeout overrides `Config.Timeout`. Nil inherits the
client default, and explicit zero adds no Kit deadline. Closing the client does
not change the cloud sandbox and invalidates handles that depend on its SDK.
`Sandbox.Close` closes only the local handle; it never changes the remote
sandbox. Client and handle closure are idempotent and wait for operations already
running through that object.

## Operation semantics

`Start` enters the running state from a provider-supported retained state.
`Stop` requests a restartable stop where the sandbox class supports one; an
ephemeral provider resource can still be destroyed by its configured policy.
`Pause` requests a provider memory-preserving transition; successful return
means the provider left its pausing transition and callers should inspect the
returned handle metadata for the resolved state. `Terminate` permanently finishes
an execution without promising deletion of its provider record. `Archive` moves
retained state to provider archive storage. `Delete` irreversibly removes the
provider sandbox resource. These meanings are not aliases.

Unsupported operations return `ErrorKindUnsupported` before an operation call
when Kit knows the selected provider or resolved sandbox class cannot honor the
intent. Account permissions and dynamically enabled features remain
provider-enforced.

`Capabilities` returns support for the selected provider and the sandbox's
resolved isolation class. `Supports` is the convenient single-capability check.
Capabilities can change after `Refresh` when the provider reports a different
resolved class or terminal state.

`SandboxInfo.State` is Kit's portable lifecycle state. `ProviderState` retains
the provider's exact value. After successful `Terminate` or `Delete`, Kit marks
the handle terminal, exposes no further lifecycle capabilities, and rejects
subsequent lifecycle calls with `ErrorKindFailedPrecondition`.

Lifecycle transitions on one handle are serialized. This does not serialize a
different handle or process referring to the same remote sandbox. A caller
waiting behind an in-progress transition can still cancel its context or reach
its deadline.

## Provider support

| Operation | Modal | Daytona container | Daytona Linux/Windows VM | Daytona GPU |
| --- | --- | --- | --- | --- |
| Attach | Supported; ID only, no provider lookup | Unsupported | Unsupported | Unsupported |
| Refresh/Get | Unsupported | Supported | Supported | Supported |
| Start | Unsupported | Supported | Supported | Provider-enforced |
| Stop | Unsupported | Supported; filesystem retained unless ephemeral | Supported; memory cleared | Destructive by provider design |
| Pause | Unsupported | Unsupported | Supported; memory retained | Unsupported |
| Terminate | Supported | Unsupported | Unsupported | Unsupported |
| Archive | Unsupported | Supported | Unsupported | Unsupported |
| Delete | Unsupported | Supported | Supported | Supported |

Daytona reports its resolved sandbox class after creation. Kit maps container
and Linux VM to shared isolation kinds and preserves Windows as a provider-specific
kind. GPU allocation is provider metadata and does not replace the underlying
isolation class. Snapshot references can
select a class that was not knowable from `CreateOptions`, so lifecycle checks
use current provider metadata rather than the original request.

Modal gVisor and Linux VM sandboxes share Modal's terminal execution lifecycle.
`Terminate` maps to Modal termination; it is deliberately not presented as
`Stop` or `Delete`. Modal's native `FromID` performs no remote lookup, so Kit
exposes it as `Client.Attach`; Modal does not advertise `Get` or `Refresh`.

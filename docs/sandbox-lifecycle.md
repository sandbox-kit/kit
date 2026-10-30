# Sandbox lifecycle

`CreateOptions` describes creation intent and automatic policies. Lifecycle
methods operate on the resolved sandbox returned by `Client.Create`, a
provider-verified `Client.Get`, or an explicitly unverified `Client.Attach`;
they do not reuse or mutate the creation request.

```go
instance, err := client.Create(ctx, options)
if err != nil {
    return err
}
if !instance.Supports(sandbox.LifecycleCapabilityStop) {
    return fmt.Errorf("sandbox %s cannot be stopped", instance.ID())
}
return instance.Stop(ctx, &sandbox.StopOptions{Timeout: sandbox.Duration(time.Minute)})
```

An operation-specific timeout overrides `Config.Timeout` / `timeout` as Kit's
deadline for that call. An omitted timeout inherits the client deadline.
Explicit zero adds no Kit deadline of its own. The provider receives a timeout
only when the call sets one.

Daytona start, stop, pause, and delete use the native 60-second default when
the timeout is omitted. An explicit zero is passed through, and those methods
treat zero as no operation timeout. A timeout that is not a whole number of
seconds is rejected before the native call. A terminate timeout is how long
Kit waits. Modal does not take that duration as its own terminate timeout.
With wait, Modal polls until the sandbox exits or the deadline fires.
`Client.Close` / `client.close()` closes your connection to Daytona or Modal
and returns that result. It does not stop or delete the sandbox. A later call
uses that connection and returns the provider's result.
`Sandbox.Close` / `sandbox.close()` does not change the sandbox. You can keep
using the handle. Closing a Daytona handle does not contact Daytona. Closing a
Modal handle does not terminate the sandbox.

TypeScript clients and sandbox handles implement `AsyncDisposable`. Prefer a
scoped client when its lifetime matches the current block. Lifecycle call
options use millisecond `timeout` (for example
`await sandbox.stop({ force: true, timeout: 60_000 })`):

```ts
await using client = new Client(config);
const sandbox = await client.create(request);
await sandbox.stop({ timeout: 60_000 });
```

Disposing a client closes your connection to Daytona or Modal. Disposing a
handle does not stop, terminate, or delete the sandbox. You can keep using the
handle.

## Operation semantics

`Start` enters the running state from a provider-supported retained state.
`Stop` requests a stop. An ephemeral sandbox can still be destroyed by its
configured policy.
`Pause` requests a provider memory-preserving transition; successful return
means the provider left its pausing transition and callers should inspect the
returned handle metadata for the resolved state. `Terminate` permanently finishes
an execution without promising deletion of its provider record. `Archive` moves
retained state to provider archive storage. `Delete` irreversibly removes the
provider sandbox resource. These meanings are not aliases.

Unsupported operations return `ErrorKindUnsupported` before a provider call
when the selected provider does not implement them. Daytona accepts or rejects
a call for a particular sandbox, and Kit returns that error. Account
permissions and dynamically enabled features remain provider-enforced.

`Capabilities` lists the operations the selected provider implements.
`Supports` checks one of them. Close does not clear the list.

`SandboxInfo.State` is Kit's portable lifecycle state. `ProviderState` retains
the provider's exact value. `Terminal` is present when the provider reports
that the sandbox is finished. A later lifecycle call still goes to the provider,
and Kit returns that result.

Each SDK serializes lifecycle transitions on one handle. This does not
serialize a different handle or process referring to the same remote sandbox. A
caller waiting behind an in-progress transition can still cancel its context or
abort signal, or reach its deadline. Go forwards that context into the native
SDK. TypeScript checks the abort signal before each native call, so a canceled
operation does not start the next provider request.

## Provider support

| Operation   | Modal                                  | Daytona     |
| ----------- | -------------------------------------- | ----------- |
| Attach      | Supported; ID only, no provider lookup | Unsupported |
| Refresh/Get | Unsupported                            | Supported   |
| Start       | Unsupported                            | Supported   |
| Stop        | Unsupported                            | Supported   |
| Pause       | Unsupported                            | Supported   |
| Terminate   | Supported                              | Unsupported |
| Archive     | Unsupported                            | Supported   |
| Delete      | Unsupported                            | Supported   |

Daytona reports its resolved sandbox class after creation. Kit maps container
and Linux VM to shared isolation kinds and preserves Windows as a provider-specific
kind. GPU allocation is provider metadata and does not replace the underlying
isolation class. Whether a sandbox accepts pause, archive, or stop is Daytona's
decision.

Modal gVisor and Linux VM sandboxes share Modal's terminal execution lifecycle.
`Terminate` maps to Modal termination; it is deliberately not presented as
`Stop` or `Delete`. Modal's native `FromID` performs no remote lookup, so Kit
exposes it as `Client.Attach`; Modal does not advertise `Get` or `Refresh`.

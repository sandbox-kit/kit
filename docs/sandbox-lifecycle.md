# Sandbox lifecycle

Lifecycle methods operate on a handle from `Create`, a provider-verified
`Get`, or Modal `Attach`. They do not reuse the creation request.

## TypeScript

```ts
await using client = new Client(config);
const sandbox = await client.create(request);
if (!sandbox.supports(LifecycleCapability.Stop)) {
  throw new Error(`sandbox ${sandbox.id} cannot be stopped`);
}
await sandbox.stop({ timeout: 60_000 });
```

Clients and handles implement `AsyncDisposable`. Prefer `await using` when the
lifetime matches the block. Timeouts are milliseconds.

Abort rejects the `await`. The handle stays busy until the provider call
settles, so the next call on the same handle does not overlap it.

## Go

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

Cancel or a deadline on the context ends the wait. The handle stays busy until
the provider call returns.

## Timeouts

An operation-specific timeout overrides the client default. Omitted timeout
inherits the client deadline. Explicit zero adds no client deadline.

**Daytona.** Start, stop, pause, and delete use a native 60-second default when
the timeout is omitted. Explicit zero means no operation timeout. Non-whole
seconds are rejected.

**Modal.** A terminate timeout is how long the SDK waits. With wait, Modal
polls until the sandbox exits or the deadline fires.

## Close

`client.close()` / `Client.Close` closes your connection to the provider. It
does not stop or delete the sandbox.

`sandbox.close()` / `Sandbox.Close` releases local handle resources only. You
can keep using the handle afterward for Daytona and Blaxel local close; Modal
close does not terminate the sandbox.

## Operation meanings

| Method    | Meaning                                                              |
| --------- | -------------------------------------------------------------------- |
| Start     | Enter the running state from a retained state                        |
| Stop      | Request a stop (ephemeral policy may still destroy the sandbox)      |
| Pause     | Memory-preserving transition; inspect returned metadata for state    |
| Terminate | Permanently finish execution; does not promise record deletion       |
| Archive   | Move retained state to archive storage                               |
| Delete    | Irreversibly remove the provider sandbox                             |

Unsupported operations return `Unsupported` before a provider call.
`supports` / `Supports` checks whether the provider implements an operation.

Prefer portable `state` / `State` and `terminal` / `Terminal` on the handle.
A later lifecycle call still goes to the provider.

## Provider support

| Operation   | Modal                                  | Daytona     | Blaxel                             |
| ----------- | -------------------------------------- | ----------- | ---------------------------------- |
| Attach      | Supported; ID only, no provider lookup | Unsupported | Unsupported                        |
| Refresh/Get | Unsupported                            | Supported   | Supported (lookup by sandbox name) |
| Start       | Unsupported                            | Supported   | Unsupported                        |
| Stop        | Unsupported                            | Supported   | Unsupported                        |
| Pause       | Unsupported                            | Supported   | Unsupported                        |
| Terminate   | Supported                              | Unsupported | Unsupported                        |
| Archive     | Unsupported                            | Supported   | Unsupported                        |
| Delete      | Unsupported                            | Supported   | Supported                          |

### Daytona

Whether a sandbox accepts pause, archive, or stop is Daytona’s decision.

### Modal

`Terminate` is not presented as stop or delete. `Attach` builds a handle from
an id with no remote lookup. `Get` and `Refresh` are not available.

### Blaxel

Lookup is by name for `Get` and `Refresh`. After a successful `Delete`,
`state` is deleted and `terminal` is true. Omit wait to return when Blaxel
accepts the delete; recreating the same name too soon may return `Conflict`.
Pass wait to wait until removal finishes. Attach, start, stop, pause, archive,
and terminate are not available.

File calls: [filesystem](filesystem.md).

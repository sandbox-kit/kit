# Process

Run a command inside an open sandbox and wait for it to finish. Use the same
options on Daytona and Modal.

| You want to                           | Go    | TypeScript |
| ------------------------------------- | ----- | ---------- |
| Run a command and wait for the result | `Run` | `run`      |

Open a handle with `Create`, Daytona `Get`, or Modal `Attach`, then call
`Process()` / `sandbox.process`.

Calls on one handle run one at a time (including filesystem and lifecycle). A
second handle for the same sandbox does not wait on the first.

## Command shape

Set exactly one of `argv` or `shell`.

- `argv` — program and arguments, for example `["echo", "hello"]`
- `shell` — one shell command string, for example `"echo hello"`

Setting both or neither fails with field `argv`. A blank or whitespace-only
`shell` (after trim) fails with field `shell`.

How providers run the command:

| You pass | Daytona                   | Modal                    |
| -------- | ------------------------- | ------------------------ |
| `argv`   | POSIX-quoted shell string | argument list as-is      |
| `shell`  | shell string as-is        | `bash -c` with that text |

## Working directory and environment

`workingDirectory` / `WorkingDirectory`, when set, must be an absolute path
such as `/workspace`. Relative, blank, or NUL-containing paths fail with field
`working_directory`.

`environment` / `Environment` adds variables for that call only:

- Omit the map to leave the provider default alone.
- Pass an empty map to add no variables for this call.

## Result

When the command finishes you get:

- `exitCode` / `exit_code`
- `stdout`
- `stderr`

A non-zero exit code is **not** an SDK error. Check the exit code in your
code.

| Provider | Output                                                            |
| -------- | ----------------------------------------------------------------- |
| Daytona  | stdout only; `stderr` is empty (Go returns a non-nil empty slice) |
| Modal    | both stdout and stderr                                            |

## Timeouts and cancel

| Setting                      | Meaning                                                                                                                           |
| ---------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| Omit `timeout`               | Use the client default deadline for this call                                                                                     |
| `timeout: 0` / `Duration(0)` | Do not add an extra client deadline for this call                                                                                 |
| Positive timeout             | Cap this call; must be a whole number of seconds (`5000` ms or `5 * time.Second` is fine; `1500` ms is rejected, field `timeout`) |

Omit or zero also means **no separate provider command limit**. Provider
quirks:

- **Daytona** — if you omit a native timeout, Daytona defaults to 10 seconds.
  This client therefore uses timeout `0` for omit/zero so the command is not
  capped at 10 seconds.
- **Modal** — omit/zero adds no Modal command limit. In TypeScript the client
  omits `timeoutMs` (sending `0` is rejected by Modal).

Cancel behavior:

- **Go** — cancel or deadline on the context ends the wait and returns an
  error. The handle stays busy until that provider call returns.
- **TypeScript** — abort rejects the `await`. The remote command may keep
  running; Daytona and Modal exec do not take an `AbortSignal`. The handle
  still stays busy until that provider call settles, so the next call on the
  same handle does not overlap it.

```ts
const result = await sandbox.process.run({
  argv: ["echo", "hello"],
  timeout: 5_000,
});
```

```go
result, err := instance.Process().Run(ctx, &sandbox.ExecOptions{
	Argv:    []string{"echo", "hello"},
	Timeout: sandbox.Duration(5 * time.Second),
})
```

## Run a command

```ts
const result = await sandbox.process.run({ argv: ["echo", "hello from kit"] });
console.log(result.exitCode, new TextDecoder().decode(result.stdout));
```

```go
result, err := instance.Process().Run(ctx, &sandbox.ExecOptions{
	Argv: []string{"echo", "hello from kit"},
})
```

[Go Daytona](../examples/go/daytona/process/run-daytona-command/README.md) ·
[TypeScript Daytona](../examples/typescript/daytona/process/run-daytona-command/README.md) ·
[Go Modal](../examples/go/modal/process/run-modal-command/README.md) ·
[TypeScript Modal](../examples/typescript/modal/process/run-modal-command/README.md).

## More examples

| Topic              | Go Daytona                                                                                      | TypeScript Daytona                                                                                      | Go Modal                                                                                  | TypeScript Modal                                                                                  |
| ------------------ | ----------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| Shell string       | [run-daytona-shell-command](../examples/go/daytona/process/run-daytona-shell-command/README.md) | [run-daytona-shell-command](../examples/typescript/daytona/process/run-daytona-shell-command/README.md) | [run-modal-shell-command](../examples/go/modal/process/run-modal-shell-command/README.md) | [run-modal-shell-command](../examples/typescript/modal/process/run-modal-shell-command/README.md) |
| Working directory  | [run-daytona-with-cwd](../examples/go/daytona/process/run-daytona-with-cwd/README.md)           | [run-daytona-with-cwd](../examples/typescript/daytona/process/run-daytona-with-cwd/README.md)           | [run-modal-with-cwd](../examples/go/modal/process/run-modal-with-cwd/README.md)           | [run-modal-with-cwd](../examples/typescript/modal/process/run-modal-with-cwd/README.md)           |
| Environment        | [run-daytona-with-env](../examples/go/daytona/process/run-daytona-with-env/README.md)           | [run-daytona-with-env](../examples/typescript/daytona/process/run-daytona-with-env/README.md)           | [run-modal-with-env](../examples/go/modal/process/run-modal-with-env/README.md)           | [run-modal-with-env](../examples/typescript/modal/process/run-modal-with-env/README.md)           |
| Timeout            | [run-daytona-with-timeout](../examples/go/daytona/process/run-daytona-with-timeout/README.md)   | [run-daytona-with-timeout](../examples/typescript/daytona/process/run-daytona-with-timeout/README.md)   | [run-modal-with-timeout](../examples/go/modal/process/run-modal-with-timeout/README.md)   | [run-modal-with-timeout](../examples/typescript/modal/process/run-modal-with-timeout/README.md)   |
| Non-zero exit code | [run-daytona-check-exit](../examples/go/daytona/process/run-daytona-check-exit/README.md)       | [run-daytona-check-exit](../examples/typescript/daytona/process/run-daytona-check-exit/README.md)       | [run-modal-check-exit](../examples/go/modal/process/run-modal-check-exit/README.md)       | [run-modal-check-exit](../examples/typescript/modal/process/run-modal-check-exit/README.md)       |

## Limits

Streaming stdin/stdout, long-lived sessions, PTY, and Daytona `code_run` are
not part of this API.

Closing the handle leaves the sandbox running. Closing the client closes your
connection to the provider.

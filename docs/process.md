# Process

Command calls live on the sandbox handle, next to lifecycle and filesystem.
Daytona and Modal implement the same `Run` method. Open the handle with
`Create`, Daytona `Get`, or Modal `Attach`, then call `Process` /
`sandbox.process`.

| You want to                         | Go  | TypeScript |
| ----------------------------------- | --- | ---------- |
| Run a command and wait for the result | `Run` | `run`    |

Calls on one handle run one at a time, including lifecycle and filesystem
calls. A second handle for the same sandbox does not wait.

## Command shape

Set exactly one of `argv` or `shell`.

- `argv` is the program and its arguments, such as `["echo", "hello"]`.
- `shell` is one shell command string, such as `"echo hello"`.

Setting both, or neither, fails before the provider is called. The error field
is `argv`. A blank `shell` fails with field `shell`.

Daytona receives a shell string. Kit maps `argv` to a POSIX-quoted command.
Modal receives an argument vector. Kit maps `shell` to `bash -c`.

## Working directory and environment

A working directory, when set, must be absolute: `/workspace`. A relative
path, a blank path, or a path that contains a NUL byte is rejected before the
provider is called. The error field is `working_directory`.

`environment` sets variables for that command. An omitted map leaves the
provider default. An empty map is explicit.

## Result

`Run` returns when the command finishes:

- `exit_code` / `exitCode`
- `stdout`
- `stderr`

A non-zero exit code is not an SDK error. Check the exit code in your program.

Daytona's basic exec returns stdout only. `stderr` is empty on Daytona. Modal
pipes both streams and Kit collects them.

## Timeouts

Leave the timeout unset to use the client deadline. Set it to zero when this
call should add no Kit deadline. TypeScript records the timeout in
milliseconds. Go records it as a `time.Duration`.

An omitted command timeout, and an explicit zero, are passed to the provider as
no command limit. Daytona's omitted native default is 10 seconds, so Kit passes
`0` instead. Provider command timeouts must be whole seconds. `2000`
milliseconds is two seconds. `1500` is rejected. The error field is `timeout`.
When the timeout is omitted, the abort signal or context still carries the
client deadline.

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

## Provider notes

- Streaming stdin/stdout, long-lived sessions, PTY, and Daytona `code_run` are
  outside this surface.
- Closing the handle leaves the sandbox as it is. Closing the client closes
  your connection to the provider.

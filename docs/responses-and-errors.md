# Responses and errors

A failed call returns an error with a stable kind. Validation failures also
include a `rule` such as `finite`, `nonblank`, `exclusive_minimum`, or
`known_enum`, so you can branch on that value. `Info()` returns a copy of the
sandbox metadata.

## What you get back

Create returns a sandbox handle. `info()` / `Info()` returns an independent
copy. Resource numbers on the handle are what the provider reported.

Field paths use the same string in both languages (for example
`resources.cpu_limit_cores`). A rule about several settings names their parent
(for example `network`).

## TypeScript

Operations reject with `SandboxError`. Use `instanceof` before reading it:

```ts
import {
  CreateFieldResourcesCpuCores,
  SandboxError,
} from "@sandbox-kit/sandbox";

try {
  await client.create(request);
} catch (error) {
  if (error instanceof SandboxError) {
    if (error.field === CreateFieldResourcesCpuCores) {
      console.error("CPU setting rejected:", error.message);
    } else {
      console.error(error.kind, error.field, error.operation);
    }
  } else {
    throw error;
  }
}
```

The provider’s own error stays on `cause`.

Examples:
[handle-daytona-errors](../examples/typescript/daytona/errors/handle-daytona-errors/README.md),
[handle-modal-errors](../examples/typescript/modal/errors/handle-modal-errors/README.md),
[handle-blaxel-errors](../examples/typescript/blaxel/errors/handle-blaxel-errors/README.md).

## Go

```go
info := instance.Info()
if info.Resources != nil {
	fmt.Println(info.Resources.GetCpuCores())
}
```

```go
var detail *sandbox.Error
if errors.As(err, &detail) {
	switch detail.Kind {
	case sandbox.ErrorKindInvalidArgument, sandbox.ErrorKindUnsupported:
		fmt.Println("Check configuration:", detail.Field)
	case sandbox.ErrorKindAuthentication:
		fmt.Println("Check provider credentials")
	case sandbox.ErrorKindAlreadyExists:
		fmt.Println("Named resource already exists")
	case sandbox.ErrorKindConflict:
		fmt.Println("Current state rejects this operation")
	}
}
```

`Unwrap` keeps the provider error, so `errors.Is` / `errors.As` still match
cancellation and native types:

```go
if errors.Is(err, context.Canceled) {
	// The original cancellation still matches.
}
var native *sdkerrors.DaytonaError
if errors.As(err, &native) {
	// Native code, source, status, and headers stay available.
}
```

Examples:
[handle-modal-errors](../examples/go/modal/errors/handle-modal-errors/README.md),
[handle-daytona-errors](../examples/go/daytona/errors/handle-daytona-errors/README.md),
[handle-blaxel-errors](../examples/go/blaxel/errors/handle-blaxel-errors/README.md).

## Error shape

Failed calls return `Kind`, `Provider`, `Operation`, `Field`, `Rule`, and
`Message`. `StatusCode`, `ProviderCode`, and `ProviderSource` carry provider
codes when present. `NotFound` covers a missing sandbox or file path.

A kind is a category, not a retry signal. The SDK does not retry for you.

### Blaxel

| Situation                                         | Kind              |
| ------------------------------------------------- | ----------------- |
| Create when the name already exists               | `AlreadyExists`   |
| Delete still finishing / other conflict           | `Conflict`        |
| Lifetime below floor or not whole seconds         | `InvalidArgument` |
| Unsupported create field                          | `Unsupported`     |

Control-plane codes appear on `providerCode` / `ProviderCode` when present.

## After a failure

An error does not prove that no sandbox exists. Readiness or metadata handling
can fail after the provider has created one. A lifecycle call can fail after
the provider has already changed the sandbox. Refresh the handle when the
provider supports refresh before you decide to retry. When creation fails
after the provider has allocated a sandbox, the error does not include that
sandbox ID.

See [creation](sandbox-creation.md) and [lifecycle](sandbox-lifecycle.md).

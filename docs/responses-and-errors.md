# Responses and errors

A failed call returns an error with a stable kind. Validation failures also
include a `rule` such as `finite`, `nonblank`, `exclusive_minimum`, or
`known_enum`, so you can branch on that value. `Info()` returns a copy of the
sandbox metadata.

## What you get back

`Create` returns a sandbox handle. `Info()` returns an independent copy.
Optional fields that the provider left unset stay absent. Resource numbers on
the handle are what the provider reported, which can differ from the request.

```go
info := instance.Info()
if info.Resources != nil {
    fmt.Println(info.Resources.GetCpuCores())
}
```

A present optional value keeps explicit zero, false, and empty strings. A
present map or slice keeps an empty collection. A create or lookup fails in
the client when the provider returns no sandbox ID. When the provider sets
`provider`, it must match the provider you selected.

## Read an error

Initialization, creation, lifecycle calls, and cleanup return the same error
shape: `Kind`, `Provider`, `Operation`, `Field`, `Rule`, and `Message`.
`StatusCode`, `ProviderCode`, and `ProviderSource` carry the provider's own
codes when it sent them.

```go
var detail *sandbox.Error
if errors.As(err, &detail) {
    switch detail.Kind {
    case sandbox.ErrorKindInvalidArgument, sandbox.ErrorKindUnsupported:
        fmt.Println("Check configuration:", detail.Field)
    case sandbox.ErrorKindAuthentication:
        fmt.Println("Check provider credentials")
    }
}
```

The same switch works for either provider. Runnable examples:
[Modal unsupported settings](../examples/go/modal/handle-modal-errors/README.md)
and
[Daytona invalid configuration](../examples/go/daytona/handle-daytona-errors/README.md).

Go and TypeScript use the same field string. A rule about one setting names
that setting, such as `resources.cpu_limit_cores`. A rule about several
settings names their parent, such as `network`.

TypeScript operations reject with `SandboxError`. Use `instanceof` before
reading it, and compare `error.field` with the exported field constants:

```ts
import {
  CreateFieldResourcesCpuCores,
  ErrorKind,
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

TypeScript examples:
[handle-daytona-errors](../examples/typescript/daytona/handle-daytona-errors/README.md)
and
[handle-modal-errors](../examples/typescript/modal/handle-modal-errors/README.md).

## Provider errors

The provider's own error stays attached. In Go, `Unwrap` returns it, so
`errors.Is` and `errors.As` still match cancellation and the provider type:

```go
if errors.Is(err, context.Canceled) {
    // The original cancellation still matches.
}
var native *sdkerrors.DaytonaError
if errors.As(err, &native) {
    // Native code, source, status, and headers stay available.
}
```

The last snippet imports Daytona's error package as `sdkerrors`. You can
branch on `Kind` without importing that package.

TypeScript keeps the provider error on the standard `cause` property of
`SandboxError`.

An error that Kit could not classify keeps an unknown kind and the original
cause. A kind is a category. It is not a signal that retrying is safe, and
Kit does not retry for you.

## After a failure

An error does not prove that no sandbox exists. Readiness or metadata handling
can fail after the provider has created one. A lifecycle call can fail after
the provider has already changed the sandbox. Refresh the handle when the
provider supports refresh before you decide to retry. When creation fails
after the provider has allocated a sandbox, the error does not include that
sandbox ID.

See [creation](sandbox-creation.md) and [lifecycle](sandbox-lifecycle.md).

# Responses and errors

A failed call gives you `ErrorInfo` with a stable `ErrorKind`. Validation
failures also include a `rule` such as `finite`, `nonblank`,
`exclusive_minimum`, or `known_enum`, so you can branch on that instead of the
message text. `Info()` returns `SandboxInfo`.

## Response presence

`Client.Create` returns a `Sandbox`. `Info()` returns an independent metadata
copy. Unavailable optional fields stay absent; allocations are not filled from
the request.

```go
info := instance.Info()
if info.Resources != nil {
    fmt.Println(info.Resources.GetCpuCores())
}
```

Optional pointers preserve explicit zero, false, and empty strings; non-nil maps
and slices preserve empty collections. Non-optional strings require a nonempty
value. A create or lookup fails locally when the provider returns no sandbox id.
When the provider sets `provider`, it must match the selected provider.

## Handle errors in an application

Initialization, creation, lifecycle operations, and cleanup return `*sandbox.Error`
with `Kind`, `Provider`, `Operation`, `Field`, `Rule`, and `Message`. Optional
`StatusCode`, `ProviderCode`, and `ProviderSource` retain native machine details
when available.

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

The same switch works for either provider. Runnable examples demonstrate
[Modal unsupported intent](../examples/go/modal/handle-modal-errors/README.md) and
[Daytona invalid configuration](../examples/go/daytona/handle-daytona-errors/README.md).

TypeScript operations reject with `SandboxError`. Use `instanceof` before
reading its shared details. `error.field` is a canonical protobuf path; compare
it to exported constants such as `CreateFieldResourcesCpuCores`:

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

See also the TypeScript examples
[handle-daytona-errors](../examples/typescript/daytona/handle-daytona-errors/README.md)
and [handle-modal-errors](../examples/typescript/modal/handle-modal-errors/README.md).

## Native causes

Local errors are classified at their source rather than by parsing messages.
Validator errors and original SDK errors remain accessible through `Unwrap()`:

```go
if errors.Is(err, context.Canceled) {
    // Original cancellation remains matchable.
}
var native *sdkerrors.DaytonaError
if errors.As(err, &native) {
    // Native code, source, status, and headers remain accessible.
}
```

The last snippet imports Daytona's native errors package as `sdkerrors`.
Applications can use common kinds without importing provider error types.
Shared error details serialize without the runtime cause.

TypeScript preserves the native error as the standard `cause` property on
`SandboxError`.

Modal uses pinned native error types and gRPC codes. Daytona uses native HTTP
status and code/source fields. Unclassified errors retain an unknown kind and
their original cause; no speculative retry classification is added.

## Constructing errors

Applications normally receive errors from the client. Custom integrations can
construct a common error with a protobuf `ErrorInfo` and a runtime cause:

```go
err := sandbox.NewError(&sandbox.ErrorInfo{
    Kind:      sandbox.ErrorKindUnsupported,
    Provider:  "modal",
    Operation: "create",
    Field:     sandbox.CreateFieldResourcesDiskMib,
    Message:   "Disk overrides are unsupported",
}, cause)
```

`NewError` is handwritten. The constructor owns a copy of the details,
including optional pointers. The native cause stays separate from the
serializable contract. Cancellation and deadline causes override the supplied
kind (`canceled` / `timeout`). Provider modules add provider and operation
context through `WithErrorContext` in `sdks/go/internal/sdk`. That helper
reuses an error when its provider and operation already match, and copies the
error only when the context must change.

## Classification

Classification is handwritten in each language and provider so native SDK type
names, `errors.As` / `instanceof`, gRPC codes, and HTTP status maps stay local:

- Precedence: existing shared error, then native type, then protocol status/code,
  then unknown.
- Modal: native Modal error types, then gRPC status codes.
- Daytona: `DaytonaError` provider codes, then HTTP statuses.

See [generation](code-generation.md#response-and-error-contracts) for how
Protovalidate create/config checks (plus handwritten `secure_endpoint`), and
mappings construct `NewError` call sites.

## Exact metadata copying

Generated `Clone()` methods copy typed fields directly, preserving numeric types,
optional presence, and empty collections. Copying uses typed generated methods.

TypeScript returns independent metadata through `structuredClone`, preserving
its native `bigint`, `Date`, byte-array, collection, and optional-property
representations.

`MetadataObject` and `MetadataValue` describe signed and unsigned integers,
numbers, strings, booleans, bytes, lists, and objects. An empty value represents
null. Provider adapters fill `MetadataObject.Values` directly.
`MetadataObjectAsMap` reads that object as a Go map. TypeScript create options
accept a plain object and materialize it to the same protobuf message.

Ordinary JSON decoding into `map[string]any` can still change numeric types.
The owned copy keeps the types stored on `MetadataValue`. No new transport or
runtime serialization codec is introduced.

## Completion and retries

An error does not prove that no sandbox exists: readiness or metadata handling
can fail after creation. A lifecycle call can also fail after the provider has
begun or completed a transition. Refresh the handle when the provider supports
it before deciding whether to retry. Kit does not surface a partial-success
creation identity when creation fails after the provider has allocated a
resource. Error kinds do not imply safe retries, and Kit performs no automatic
retries.

See [creation](sandbox-creation.md), [generation](code-generation.md), and
[provider verification](provider-verification.md).

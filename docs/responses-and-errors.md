# Responses and errors

Protobuf defines `SandboxInfo`, `ValueOrigin`, `ErrorInfo`, and `ErrorKind`.
Validation failures include an optional stable `rule` identifier so applications
do not need to parse human-readable messages. Examples include `finite`,
`nonblank`, `exclusive_minimum`, `known_enum`, and cross-field operation names.
Error wrappers, cause chaining, and provider classification are handwritten per
language. YAML inventories describe response mappings and validation for humans
to keep aligned with that handwritten code.

## Response presence and origin

`Client.Create` returns a `Sandbox`. `Info()` returns an independent metadata
copy. Unavailable optional fields stay absent; allocations are not filled from
the request.

```go
info := instance.Info()
if info.Resources != nil {
    fmt.Println(info.Resources.GetCpuCores())
}
switch info.Origins[sandbox.InfoFieldName] {
case sandbox.ValueOriginProvider:
    // Reported by the SDK.
case sandbox.ValueOriginRequest:
    // Derived from creation settings.
default:
    // Origin is unknown.
}
```

Origin keys are shared field names: `id`, `name`, `labels`, `resources`, and so on.
Nested field paths such as `resources.cpu_cores` are also supported. Origins must
identify known fields that are present. Optional pointers preserve explicit zero,
false, and empty strings; non-nil maps and slices preserve empty collections.
Non-optional strings require a nonempty value. Unknown keys or absent targets produce
an invalid-response error.
Missing entries mean unknown. Modal identity is provider-reported, while its
current name/labels are request-derived. Daytona maps SDK-reported metadata.

## Handle errors in an application

Initialization, creation, lifecycle operations, and cleanup expose `*sandbox.Error`.
Its protobuf `ErrorInfo` contains `Kind`, `Provider`, `Operation`, `Field`, `Rule`,
and `Message`. Optional `StatusCode`, `ProviderCode`, and `ProviderSource`
retain native machine details when available.

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
reading its shared details:

```ts
try {
  await client.create(request);
} catch (error) {
  if (error instanceof SandboxError) {
    console.error(error.kind, error.field, error.operation);
  } else {
    throw error;
  }
}
```

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
err := sandbox.NewError(sandbox.ErrorInfo{
    Kind:      sandbox.ErrorKindUnsupported,
    Provider:  "modal",
    Operation: "create",
    Field:     sandbox.CreateFieldResourcesDiskMib,
    Message:   "Disk overrides are unsupported",
}, cause)
```

`NewError` and `WithErrorContext` are handwritten. The constructor owns a copy
of the details, including optional pointers. The native cause stays separate
from the serializable contract. Cancellation and deadline causes override the
supplied kind (`canceled` / `timeout`). Context enrichment reuses an error when
its provider/operation already match and copies only when context must change.

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

`MetadataObject` and `MetadataValue` describe exact signed/unsigned integers,
finite numbers, strings, booleans, bytes, lists, and objects. An empty value
represents null. Go exposes metadata as `map[string]any` and supports primitive
integer/float types, `[]byte`, `[]string`, `[]any`, and nested `map[string]any`.
Unsupported objects, nonfinite values, cycles, and excessive depth are rejected.
`specs/contracts.yaml` limits each metadata map copy to depth 64, 16,384 visited
values (including containers and string-list entries), and 1 MiB of string,
map-key, and byte-slice payload. These are traversal and payload limits;
container allocation overhead is separate.

Ordinary JSON decoding into `map[string]any` can still change numeric types;
the guarantee applies to native values and their owned copies. No new transport
or runtime serialization codec is introduced.

## Completion and retries

An error does not prove that no sandbox exists: readiness or metadata handling
can fail after creation. A lifecycle call can also fail after the provider has
begun or completed a transition. Refresh the handle when the provider supports
it before deciding whether to retry. Partial-success creation identity reporting
remains deferred. Error kinds do not imply safe retries, and Kit performs no
automatic retries.

See [creation](sandbox-creation.md), [generation](code-generation.md), and
[provider verification](provider-verification.md).

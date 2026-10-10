# Responses and errors

Protobuf defines `SandboxInfo`, `ValueOrigin`, `ErrorInfo`, and `ErrorKind`.
YAML declares provider classification and response mappings. Generators emit
native Go types, owned copies, error helpers, and provider mappings.

## Response presence and origin

`Client.Create` returns a `Sandbox`. `Info()` returns an independent metadata
copy. Unavailable optional fields stay absent; allocations are not filled from
the request.

```go
info := instance.Info()
if info.Resources != nil {
    fmt.Println(info.Resources.GetCPUCores())
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

Initialization, creation, and cleanup expose `*sandbox.Error`.
Its generated `ErrorInfo` contains `Kind`, `Provider`, `Operation`, `Field`,
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

Modal uses pinned native error types and gRPC codes. Daytona uses native HTTP
status and code/source fields. Unclassified errors retain an unknown kind and
their original cause; no speculative retry classification is added.

## Constructing errors

Applications normally receive errors from the client. Custom integrations can
construct a common error with a generated details object and a runtime cause:

```go
err := sandbox.NewError(sandbox.ErrorInfo{
    Kind:      sandbox.ErrorKindUnsupported,
    Provider:  "modal",
    Operation: "create",
    Field:     sandbox.CreateFieldResourcesDiskMiB,
    Message:   "Disk overrides are unsupported",
}, cause)
```

The constructor owns a copy of the details, including optional pointers.
The native cause stays separate from the serializable contract.
Context enrichment reuses an error when its provider/operation already match.
It creates an owned copy only when context needs to change.

## Error policy sources

[Shared error policy](../specs/errors.yaml) declares cancellation/deadline
classification, unknown fallback, context reuse, detail preservation, and
provider classification precedence. Ordered cause rules take priority over the
supplied kind: cancellation becomes `canceled`, deadline expiry becomes `timeout`.

Provider YAML contains protocol mappings and portable rules. Modal binds its
portable `error_rules` to native Go types through `errors.go.types`; Daytona
binds native status/code/source members through `errors.go.target`. Go emitters
implement matching, owned copies, and cause chaining.

Version 1 requires existing shared errors to take priority over native types,
then HTTP statuses, gRPC codes, and fallback. Unsupported policy alternatives
fail generation. See [generation](code-generation.md#response-and-error-contracts)
for the specification vocabulary.

## Exact metadata copying

Generated `Clone()` methods copy typed fields directly, preserving numeric types,
optional presence, and empty collections. Copying uses typed generated methods.

`MetadataObject` and `MetadataValue` describe exact signed/unsigned integers,
finite numbers, strings, booleans, bytes, lists, and objects. An empty value
represents null. Go exposes metadata as `map[string]any` and supports primitive
integer/float types, `[]byte`, `[]string`, `[]any`, and nested `map[string]any`.
Unsupported objects, nonfinite values, cycles, and excessive depth are rejected.
`specs/contracts.yaml` sets the maximum copy depth, currently 64.

Ordinary JSON decoding into `map[string]any` can still change numeric types;
the guarantee applies to native values and their owned copies. No new transport
or runtime serialization codec is introduced.

## Completion and retries

An error does not prove that no sandbox exists: readiness or metadata handling
can fail after creation. Partial-success identity reporting remains deferred.
Error kinds do not imply safe retries, and Kit performs no automatic retries.

See [creation](sandbox-creation.md), [generation](code-generation.md), and
[provider verification](provider-verification.md).

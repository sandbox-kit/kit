# Protobuf SDK boundary

Sandbox Kit has one portable contract source: protobuf. YAML describes
provider-specific rules only; it never declares an SDK's public API, native
types, constructors, methods, errors, or control flow.

## Contract packages

`kit.sandbox.v1` contains public, portable data contracts:

- configuration and creation options;
- sandbox metadata, lifecycle options, results, and portable error details.

Kit does not declare a protobuf RPC service. Language SDKs expose handwritten
`Client` / `Sandbox` APIs and provider `Backend` adapters over vendor SDKs.

`kit.sandbox.internal.v1` is reserved for implementation-only contract data used
by the SDK runtime and provider adapters. Public contracts live in
`client.proto` / `sandbox.proto`. Internal generated Go output belongs under a
Go `internal/` directory; TypeScript must not re-export that module from a
package entrypoint.

An internal protobuf message is still a data contract. It must not contain a
native SDK client, callback, mutex, function, or another runtime object.

## Generated versus handwritten code

Standard protobuf plugins generate messages, enums, and reflection metadata
into each language's contract package. Those types are the wire and validation
representation of portable data.

Handwritten language runtime code owns behavior:

- public `Client` and `Sandbox` facades;
- TypeScript public create/config input types (materialized to wire privately);
- provider selection, native SDK initialization, cleanup, and deadlines;
- private adapter interfaces and private native SDK state;
- error behavior and conversion at provider boundaries;
- provider checks, field mappings, and config assembly.

The runtime consumes generated request and response messages for Protovalidate
and provider mapping. It does not regenerate parallel structs or assemble
methods with a custom emitter.

## YAML scope

YAML remains only for rules that vary by provider and should be reviewed as
data:

- supported capabilities and unsupported intent;
- field-level mapping to an official provider SDK;
- validation and representation rules.

Error payloads live in protobuf (`ErrorInfo` / `ErrorKind`). Wrappers, cause
chaining, and provider-native classification are handwritten per language.

It must not prescribe method signatures, constructors, response shapes,
visibility, or operation sequences. Validation, mappings, and SDK control flow
are handwritten; YAML inventories stay documentation for humans.

## Boundary

Contract packages come from Buf-driven standard protobuf plugins
(`buf generate`). Shared create and client-config validation is declared with
Protovalidate on portable messages and evaluated in each SDK
(`secure_endpoint` is handwritten); provider checks and Client/Sandbox
choreography are handwritten.

### TypeScript public surface

Applications use handwritten public types from `@sandbox-kit/sandbox`
(`CreateOptions`, `Config` / `ClientConfig`): nested exclusive fields, millisecond
timeouts, and plain allowlists/maps. The client **materializes** those values into
generated protobuf messages privately for Protovalidate, then runs provider
checks and mappings. Application code depends on `@sandbox-kit/sandbox` (and a
provider); `@sandbox-kit/contracts` is an internal wire dependency of the SDK.

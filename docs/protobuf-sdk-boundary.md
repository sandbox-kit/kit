# Protobuf SDK boundary

Sandbox Kit has one portable contract source: protobuf.  YAML describes
provider-specific rules only; it never declares an SDK's public API, native
types, constructors, methods, errors, or control flow.

## Contract packages

`kit.sandbox.v1` contains public, portable data and service contracts:

- configuration and creation options;
- sandbox metadata, lifecycle options, results, and portable error details;
- `SandboxService` request/response messages and its operation graph.

`kit.sandbox.internal.v1` is reserved for implementation-only contract data used
by the SDK runtime and provider adapters.  It is not required for the current
public `client.proto` / `sandbox.proto` adoption.  When added, generated Go
output belongs under a Go `internal/` directory, and TypeScript must not
re-export that module from a package entrypoint.

An internal protobuf message is still a data contract.  It must not contain a
native SDK client, callback, mutex, function, or another runtime object.

## Generated versus handwritten code

Standard protobuf plugins generate messages, enums, reflection metadata, and
service declarations.  Those generated types are the sole representation of
portable data in each SDK.

Handwritten language runtime code owns behavior:

- public `Client` and `Sandbox` facades;
- provider selection, native SDK initialization, cleanup, and deadlines;
- private adapter interfaces and private native SDK state;
- error behavior and conversion at provider boundaries.

The runtime consumes generated request and response messages.  It does not
regenerate its own lookalike structs or use a generator to assemble methods.

## YAML scope

YAML remains only for rules that vary by provider and should be reviewed as
data:

- supported capabilities and unsupported intent;
- field-level mapping to an official provider SDK;
- validation and representation rules;
- provider-native error classification.

It must not prescribe method signatures, constructors, response shapes,
visibility, or operation sequences.  A future mapping/validation compiler may
generate narrow pure conversion and validation helpers from these rules.  It
must not generate SDK control flow.

## Migration rule

Do not add to the legacy custom SDK emitters.  Each migrated area is converted
to standard protobuf output plus handwritten runtime code, then its legacy
specification and generated output are removed in the same change.

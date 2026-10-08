// Package sandbox assembles an explicitly selected sandbox provider.
//
// Provider implementations are installed separately. Client exposes a shared API;
// Client.Create accepts generated creation contracts and returns a Sandbox handle.
// Optional typed providers own SDK mappings. Public data types are native Go;
// protobuf and YAML specs are generation inputs.
// NewClient accepts Config and initializes an SDK through an selected provider.
// Client.Close releases that SDK; sandbox lifecycle operations remain separate.
package sandbox

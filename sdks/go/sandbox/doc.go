// Package sandbox creates and manages cloud sandboxes through one client.
//
// Install this module and one provider. Pass daytona.New or modal.New to
// [NewClient]. [Client.Create] returns a [Sandbox] handle.
//
// [Client.Close] and [Sandbox.Close] release local resources. They leave the
// cloud sandbox running. Stop, terminate, archive, or delete it with a
// lifecycle method the provider supports.
//
// When a call fails, use errors.As with [*Error]. Kind is the shared category.
// Field names the setting to change, such as "resources.cpu_cores" or "endpoint".
// Compare it with [CreateFieldResourcesCpuCores], [ConfigFieldEndpoint], and the other field constants.
//
// [Client.Create] with a nil request uses Daytona's default snapshot. Modal requires an image.
package sandbox

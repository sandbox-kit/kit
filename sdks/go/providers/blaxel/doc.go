// Package blaxel connects the sandbox client to Blaxel.
//
// Pass [New] to sandbox.NewClient. Create requires a Blaxel-valid sandbox name
// (lowercase alphanumeric with single hyphens, at most 49 characters), an
// image source, and [sandbox.Config.Region]. Set Scope.Workspace and an API
// key, or omit them so the Blaxel SDK can use BL_API_KEY and BL_WORKSPACE from
// the environment. The backend implements shared filesystem operations through
// [sandbox.SandboxFiles] and blocking command execution through
// [sandbox.SandboxProcess]. Get and Refresh look up a sandbox by name. Delete
// removes the Blaxel resource. Closing the client releases local provider
// state and leaves the Blaxel sandbox running. Start, Stop, Pause, Archive,
// Attach, and Terminate are unsupported. resources.memory_mib maps 1:1 onto
// Blaxel's memory field; Kit has no portable volume fields.
package blaxel

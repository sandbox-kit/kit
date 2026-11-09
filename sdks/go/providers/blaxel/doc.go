// Package blaxel connects the sandbox client to Blaxel.
//
// Pass [New] to sandbox.NewClient. Create requires a sandbox name and an image
// source. The backend implements shared filesystem operations through
// [sandbox.SandboxFiles] and blocking command execution through
// [sandbox.SandboxProcess]. Closing the client does not delete the Blaxel
// sandbox.
package blaxel

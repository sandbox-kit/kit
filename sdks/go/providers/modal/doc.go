// Package modal connects the sandbox client to Modal.
//
// Pass [New] to sandbox.NewClient. Create requires an image. Set Scope.AppName
// to an existing Modal app. Closing the client releases the local SDK and
// leaves the Modal sandbox running. Reclaim it with Terminate or a lifetime
// policy.
package modal

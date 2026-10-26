// Package daytona connects the sandbox client to Daytona.
//
// Pass [New] to sandbox.NewClient. [sandbox.Client.Create] with a nil request
// uses Daytona's default snapshot. Closing the client releases the local SDK
// and leaves the Daytona sandbox running.
package daytona

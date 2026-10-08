// Package daytona maps shared client and creation configuration to its official SDK.
// New selects a factory; sandbox.NewClient initializes the SDK from sandbox.Config.
// The sandbox client owns SDK cleanup through Close.
package daytona

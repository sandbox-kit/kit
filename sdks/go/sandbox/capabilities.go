package sandbox

import "context"

// SandboxCreator creates sandboxes. Daytona and Modal implement it.
// Call Client.Create.
type SandboxCreator interface {
	// Create creates one sandbox.
	// A nil request is Daytona's default snapshot. Modal requires an image in the request.
	Create(context.Context, *CreateOptions) (*CreateResult, error)
}

// SandboxAttacher returns metadata for an existing sandbox ID without a lookup.
// Modal implements it. Call Client.Attach.
type SandboxAttacher interface {
	// AttachSandbox returns metadata for id.
	AttachSandbox(context.Context, string) (*SandboxInfo, error)
}

// SandboxRefresher loads current metadata for a sandbox ID.
// Daytona implements it. Call Client.Get or Sandbox.Refresh.
type SandboxRefresher interface {
	// RefreshSandbox loads current metadata for id.
	RefreshSandbox(context.Context, string) (*SandboxInfo, error)
}

// SandboxStarter starts a stopped sandbox.
// Daytona implements it. Call Sandbox.Start.
type SandboxStarter interface {
	// StartSandbox starts the sandbox identified by id.
	StartSandbox(context.Context, string, *StartOptions) (*SandboxInfo, error)
}

// SandboxStopper stops a sandbox.
// Daytona implements it. Call Sandbox.Stop.
type SandboxStopper interface {
	// StopSandbox stops the sandbox identified by id.
	StopSandbox(context.Context, string, *StopOptions) (*SandboxInfo, error)
}

// SandboxPauser pauses a sandbox.
// Daytona implements it and accepts or rejects the sandbox. Modal does not.
type SandboxPauser interface {
	// PauseSandbox pauses the sandbox identified by id.
	PauseSandbox(context.Context, string, *PauseOptions) (*SandboxInfo, error)
}

// SandboxTerminator ends a sandbox and can wait for its exit code.
// Modal implements it. Call Sandbox.Terminate.
type SandboxTerminator interface {
	// TerminateSandbox ends the sandbox identified by id.
	TerminateSandbox(context.Context, string, *TerminateOptions) (*ExitResult, error)
}

// SandboxArchiver archives a sandbox.
// Daytona implements it and accepts or rejects the sandbox. Modal does not.
type SandboxArchiver interface {
	// ArchiveSandbox archives the sandbox identified by id.
	ArchiveSandbox(context.Context, string, *ArchiveOptions) (*SandboxInfo, error)
}

// SandboxDeleter deletes a sandbox.
// Daytona implements it. Call Sandbox.Delete. Close does not delete the cloud sandbox.
type SandboxDeleter interface {
	// DeleteSandbox deletes the sandbox identified by id.
	DeleteSandbox(context.Context, string, *DeleteOptions) (*SandboxInfo, error)
}

package sandbox

import "context"

// SandboxCreator is the optional create capability on a provider Backend.
type SandboxCreator interface {
	Create(context.Context, *CreateOptions) (*CreateResult, error)
}

// SandboxAttacher is the optional attach capability on a provider Backend.
type SandboxAttacher interface {
	AttachSandbox(context.Context, string) (*SandboxInfo, error)
}

// SandboxRefresher is the optional refresh capability on a provider Backend.
type SandboxRefresher interface {
	RefreshSandbox(context.Context, string) (*SandboxInfo, error)
}

// SandboxStarter is the optional start capability on a provider Backend.
type SandboxStarter interface {
	StartSandbox(context.Context, string, *StartOptions) (*SandboxInfo, error)
}

// SandboxStopper is the optional stop capability on a provider Backend.
type SandboxStopper interface {
	StopSandbox(context.Context, string, *StopOptions) (*SandboxInfo, error)
}

// SandboxPauser is the optional pause capability on a provider Backend.
type SandboxPauser interface {
	PauseSandbox(context.Context, string, *PauseOptions) (*SandboxInfo, error)
}

// SandboxTerminator is the optional terminate capability on a provider Backend.
type SandboxTerminator interface {
	TerminateSandbox(context.Context, string, *TerminateOptions) (*ExitResult, error)
}

// SandboxArchiver is the optional archive capability on a provider Backend.
type SandboxArchiver interface {
	ArchiveSandbox(context.Context, string, *ArchiveOptions) (*SandboxInfo, error)
}

// SandboxDeleter is the optional delete capability on a provider Backend.
type SandboxDeleter interface {
	DeleteSandbox(context.Context, string, *DeleteOptions) (*SandboxInfo, error)
}

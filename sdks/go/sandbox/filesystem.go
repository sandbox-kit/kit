package sandbox

import (
	"context"
	"reflect"
	"strings"

	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"google.golang.org/protobuf/proto"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// Filesystem reads and writes files inside one sandbox.
// Obtain it from Sandbox.Filesystem. Paths in the sandbox are absolute.
// A nil timeout inherits Config.Timeout. Explicit zero adds no Kit deadline.
type Filesystem struct {
	owner *Sandbox
}

// Filesystem returns the file operations for this sandbox.
// Methods on a nil handle return ErrorKindFailedPrecondition.
func (s *Sandbox) Filesystem() *Filesystem {
	if s == nil {
		return &Filesystem{}
	}
	return &Filesystem{owner: s}
}

// ReadFile returns the bytes of the file at path.
// path is absolute. A nil context returns ErrorKindInvalidArgument.
func (f *Filesystem) ReadFile(ctx context.Context, path string, options *FileReadOptions) ([]byte, error) {
	var data []byte
	err := f.call(ctx, "read_file", options.GetTimeout(), path, nil, func(operationContext context.Context, id string) error {
		files, err := f.files("read_file")
		if err != nil {
			return err
		}
		data, err = files.ReadFile(operationContext, id, path, cloneMessage(options))
		return err
	})
	return data, err
}

// WriteFile replaces the file at path with data and creates missing parent directories.
// path is absolute. A nil context returns ErrorKindInvalidArgument.
func (f *Filesystem) WriteFile(ctx context.Context, path string, data []byte, options *FileWriteOptions) error {
	return f.call(ctx, "write_file", options.GetTimeout(), path, nil, func(operationContext context.Context, id string) error {
		files, err := f.files("write_file")
		if err != nil {
			return err
		}
		return files.WriteFile(operationContext, id, path, data, cloneMessage(options))
	})
}

// CopyFromLocal copies the file at localPath into remotePath.
// remotePath is absolute. localPath is a path on the caller's machine.
func (f *Filesystem) CopyFromLocal(ctx context.Context, localPath, remotePath string, options *FileCopyOptions) error {
	return f.call(ctx, "copy_from_local", options.GetTimeout(), remotePath, &localPath, func(operationContext context.Context, id string) error {
		files, err := f.files("copy_from_local")
		if err != nil {
			return err
		}
		return files.CopyFromLocal(operationContext, id, localPath, remotePath, cloneMessage(options))
	})
}

// CopyToLocal copies the file at remotePath to localPath.
// remotePath is absolute. localPath is a path on the caller's machine.
func (f *Filesystem) CopyToLocal(ctx context.Context, remotePath, localPath string, options *FileCopyOptions) error {
	return f.call(ctx, "copy_to_local", options.GetTimeout(), remotePath, &localPath, func(operationContext context.Context, id string) error {
		files, err := f.files("copy_to_local")
		if err != nil {
			return err
		}
		return files.CopyToLocal(operationContext, id, remotePath, localPath, cloneMessage(options))
	})
}

// MakeDirectory creates the directory at path.
// An omitted Parents flag is true: missing parents are created. Set false with Value.
// Daytona rejects false.
// path is absolute. A nil context returns ErrorKindInvalidArgument.
func (f *Filesystem) MakeDirectory(ctx context.Context, path string, options *MakeDirectoryOptions) error {
	return f.call(ctx, "make_directory", options.GetTimeout(), path, nil, func(operationContext context.Context, id string) error {
		files, err := f.files("make_directory")
		if err != nil {
			return err
		}
		return files.MakeDirectory(operationContext, id, path, cloneMessage(options))
	})
}

// ListDirectory lists the entries in the directory at path.
// path is absolute. The result is one directory level.
func (f *Filesystem) ListDirectory(ctx context.Context, path string, options *ListDirectoryOptions) ([]*FileInfo, error) {
	var entries []*FileInfo
	err := f.call(ctx, "list_directory", options.GetTimeout(), path, nil, func(operationContext context.Context, id string) error {
		files, err := f.files("list_directory")
		if err != nil {
			return err
		}
		entries, err = files.ListDirectory(operationContext, id, path, cloneMessage(options))
		return err
	})
	return entries, err
}

// Stat returns metadata for path.
// path is absolute. A symlink result describes the symlink.
func (f *Filesystem) Stat(ctx context.Context, path string, options *StatOptions) (*FileInfo, error) {
	var info *FileInfo
	err := f.call(ctx, "stat", options.GetTimeout(), path, nil, func(operationContext context.Context, id string) error {
		files, err := f.files("stat")
		if err != nil {
			return err
		}
		info, err = files.Stat(operationContext, id, path, cloneMessage(options))
		return err
	})
	if err != nil {
		return nil, err
	}
	return cloneMessage(info), nil
}

// Remove deletes the file at path.
// Set RemoveOptions.Recursive with Value(true) to delete a directory and its contents.
// An omitted recursive flag removes a file or an empty directory.
func (f *Filesystem) Remove(ctx context.Context, path string, options *RemoveOptions) error {
	return f.call(ctx, "remove", options.GetTimeout(), path, nil, func(operationContext context.Context, id string) error {
		files, err := f.files("remove")
		if err != nil {
			return err
		}
		return files.Remove(operationContext, id, path, cloneMessage(options))
	})
}

func (f *Filesystem) call(ctx context.Context, operation string, timeout *durationpb.Duration, sandboxPath string, localPath *string, invoke func(context.Context, string) error) error {
	owner := f.owner
	if owner == nil {
		if ctx == nil {
			return NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Operation: operation, Field: "context", Message: "sandbox-kit: context is required"}, nil)
		}
		return NewError(&ErrorInfo{Kind: ErrorKindFailedPrecondition, Operation: operation, Message: "sandbox-kit: sandbox handle is not initialized"}, nil)
	}
	if err := absoluteSandboxPath(owner, operation, sandboxPath); err != nil {
		return err
	}
	if localPath != nil {
		if err := callerPath(owner, operation, *localPath); err != nil {
			return err
		}
	}
	operationContext, cancel, id, err := owner.operationContext(ctx, operation, sdk.OptionalDuration(timeout))
	if err != nil {
		return err
	}
	defer cancel()
	if err := invoke(operationContext, id); err != nil {
		return sdk.WithErrorContext(err, owner.ProviderName(), operation)
	}
	return nil
}

func (f *Filesystem) files(operation string) (SandboxFiles, error) {
	files, ok := f.owner.operations.(SandboxFiles)
	if !ok {
		return nil, f.owner.unsupported(operation)
	}
	return files, nil
}

func absoluteSandboxPath(owner *Sandbox, operation, path string) error {
	if path == "" || !strings.HasPrefix(path, "/") || strings.ContainsRune(path, 0) {
		return NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidArgument,
			Provider:  owner.ProviderName(),
			Operation: operation,
			Field:     FileFieldPath,
			Message:   "sandbox-kit: sandbox path must be absolute",
		}, nil)
	}
	return nil
}

func callerPath(owner *Sandbox, operation, path string) error {
	if path == "" || strings.ContainsRune(path, 0) {
		return NewError(&ErrorInfo{
			Kind:      ErrorKindInvalidArgument,
			Provider:  owner.ProviderName(),
			Operation: operation,
			Field:     FileFieldLocalPath,
			Message:   "sandbox-kit: local path is required",
		}, nil)
	}
	return nil
}

func cloneMessage[T proto.Message](message T) T {
	value := reflect.ValueOf(message)
	if !value.IsValid() || (value.Kind() == reflect.Pointer && value.IsNil()) {
		return message
	}
	cloned, _ := proto.Clone(message).(T)
	return cloned
}

package daytona

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ sandbox.SandboxFiles = (*backend)(nil)

func (b *backend) ReadFile(ctx context.Context, id, path string, _ *sandbox.FileReadOptions) ([]byte, error) {
	instance, err := b.fileSandbox(ctx, id, "read_file")
	if err != nil {
		return nil, err
	}
	data, err := instance.FileSystem.DownloadFile(ctx, path, nil)
	if err != nil {
		return nil, mapProviderError(err, "read_file")
	}
	return data, nil
}

func (b *backend) WriteFile(ctx context.Context, id, path string, data []byte, _ *sandbox.FileWriteOptions) error {
	instance, err := b.fileSandbox(ctx, id, "write_file")
	if err != nil {
		return err
	}
	if len(data) == 0 {
		data = []byte{}
	}
	if err := instance.FileSystem.UploadFile(ctx, data, path); err != nil {
		return mapProviderError(err, "write_file")
	}
	return nil
}

func (b *backend) CopyFromLocal(ctx context.Context, id, localPath, remotePath string, _ *sandbox.FileCopyOptions) error {
	if err := localSourceFile(localPath, "copy_from_local"); err != nil {
		return err
	}
	instance, err := b.fileSandbox(ctx, id, "copy_from_local")
	if err != nil {
		return err
	}
	if err := instance.FileSystem.UploadFile(ctx, localPath, remotePath); err != nil {
		return mapCallerPathError(err, "copy_from_local", localPath)
	}
	return nil
}

func (b *backend) CopyToLocal(ctx context.Context, id, remotePath, localPath string, _ *sandbox.FileCopyOptions) error {
	if err := localDestination(localPath, "copy_to_local"); err != nil {
		return err
	}
	instance, err := b.fileSandbox(ctx, id, "copy_to_local")
	if err != nil {
		return err
	}
	if _, err := instance.FileSystem.DownloadFile(ctx, remotePath, &localPath); err != nil {
		return mapCallerPathError(err, "copy_to_local", localPath)
	}
	return nil
}

func (b *backend) MakeDirectory(ctx context.Context, id, path string, options *sandbox.MakeDirectoryOptions) error {
	if err := rejectParents(options); err != nil {
		return err
	}
	instance, err := b.fileSandbox(ctx, id, "make_directory")
	if err != nil {
		return err
	}
	if err := instance.FileSystem.CreateFolder(ctx, path); err != nil {
		return mapProviderError(err, "make_directory")
	}
	return nil
}

func (b *backend) ListDirectory(ctx context.Context, id, path string, _ *sandbox.ListDirectoryOptions) ([]*sandbox.FileInfo, error) {
	instance, err := b.fileSandbox(ctx, id, "list_directory")
	if err != nil {
		return nil, err
	}
	entries, err := instance.FileSystem.ListFiles(ctx, path)
	if err != nil {
		return nil, mapProviderError(err, "list_directory")
	}
	result := make([]*sandbox.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := daytonaFileInfo(entry, "list_directory")
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}

func (b *backend) Stat(ctx context.Context, id, path string, _ *sandbox.StatOptions) (*sandbox.FileInfo, error) {
	instance, err := b.fileSandbox(ctx, id, "stat")
	if err != nil {
		return nil, err
	}
	info, err := instance.FileSystem.GetFileInfo(ctx, path)
	if err != nil {
		return nil, mapProviderError(err, "stat")
	}
	return daytonaFileInfo(info, "stat")
}

func (b *backend) Remove(ctx context.Context, id, path string, options *sandbox.RemoveOptions) error {
	instance, err := b.fileSandbox(ctx, id, "remove")
	if err != nil {
		return err
	}
	recursive := options != nil && options.GetRecursive()
	if err := instance.FileSystem.DeleteFile(ctx, path, recursive); err != nil {
		return mapProviderError(err, "remove")
	}
	return nil
}

func (b *backend) fileSandbox(ctx context.Context, id, operation string) (*daytona.Sandbox, error) {
	instance, err := b.nativeSandbox(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	if instance == nil || instance.FileSystem == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "daytona",
			Operation: operation,
			Message:   "sandbox-kit daytona: SDK returned no filesystem",
		}, nil)
	}
	return instance, nil
}

func rejectParents(options *sandbox.MakeDirectoryOptions) error {
	if options == nil || options.Parents == nil || *options.Parents {
		return nil
	}
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      sandbox.ErrorKindUnsupported,
		Provider:  "daytona",
		Operation: "make_directory",
		Field:     sandbox.FileFieldParents,
		Message:   "sandbox-kit daytona: parents false is not available",
	}, nil)
}

func localSourceFile(path, operation string) error {
	info, err := os.Stat(path)
	if err != nil {
		return localPathError(err, operation)
	}
	if info.IsDir() {
		return localDirectoryError(operation)
	}
	return nil
}

func localDestination(path, operation string) error {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return localParent(path, operation)
	}
	if err != nil {
		return localPathError(err, operation)
	}
	if info.IsDir() {
		return localDirectoryError(operation)
	}
	return nil
}

func localParent(path, operation string) error {
	info, err := os.Stat(filepath.Dir(path))
	if err != nil {
		return localPathError(err, operation)
	}
	if !info.IsDir() {
		return sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "daytona",
			Operation: operation,
			Field:     sandbox.FileFieldLocalPath,
			Message:   "sandbox-kit daytona: local parent is not a directory",
		}, nil)
	}
	return nil
}

func localDirectoryError(operation string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      sandbox.ErrorKindInvalidArgument,
		Provider:  "daytona",
		Operation: operation,
		Field:     sandbox.FileFieldLocalPath,
		Message:   "sandbox-kit daytona: local path is a directory",
	}, nil)
}

func localPathError(err error, operation string) error {
	kind := sandbox.ErrorKindInvalidArgument
	if errors.Is(err, os.ErrNotExist) {
		kind = sandbox.ErrorKindNotFound
	}
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      kind,
		Provider:  "daytona",
		Operation: operation,
		Field:     sandbox.FileFieldLocalPath,
		Message:   err.Error(),
	}, err)
}

func mapCallerPathError(err error, operation, localPath string) error {
	var pathError *os.PathError
	if errors.As(err, &pathError) && filepath.Clean(pathError.Path) == filepath.Clean(localPath) {
		return localPathError(pathError, operation)
	}
	return mapProviderError(err, operation)
}

func daytonaFileInfo(info *types.FileInfo, operation string) (*sandbox.FileInfo, error) {
	if info == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "daytona",
			Operation: operation,
			Field:     sandbox.FileFieldPath,
			Message:   "sandbox-kit daytona: SDK returned no file metadata",
		}, nil)
	}
	kind := sandbox.FileTypeFile
	if info.IsDirectory {
		kind = sandbox.FileTypeDirectory
	}
	result := &sandbox.FileInfo{
		Name: info.Name,
		Path: info.Path,
		Type: kind,
		Size: info.Size,
	}
	if mode := info.Mode; mode != "" {
		result.Mode = sandbox.Value(mode)
	}
	if !info.ModifiedTime.IsZero() {
		result.ModifiedAt = timestamppb.New(info.ModifiedTime)
	}
	return result, nil
}

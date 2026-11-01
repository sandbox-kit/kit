package daytona

import (
	"context"
	"time"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ sandbox.SandboxFiles = (*backend)(nil)

func (b *backend) ReadFile(ctx context.Context, id, path string, options *sandbox.FileReadOptions) ([]byte, error) {
	if err := wholeSecondFileTimeout("read_file", options.GetTimeout()); err != nil {
		return nil, err
	}
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

func (b *backend) WriteFile(ctx context.Context, id, path string, data []byte, options *sandbox.FileWriteOptions) error {
	if err := wholeSecondFileTimeout("write_file", options.GetTimeout()); err != nil {
		return err
	}
	instance, err := b.fileSandbox(ctx, id, "write_file")
	if err != nil {
		return err
	}
	if data == nil {
		data = []byte{}
	}
	if err := instance.FileSystem.UploadFile(ctx, data, path); err != nil {
		return mapProviderError(err, "write_file")
	}
	return nil
}

func (b *backend) CopyFromLocal(ctx context.Context, id, localPath, remotePath string, options *sandbox.FileCopyOptions) error {
	if err := wholeSecondFileTimeout("copy_from_local", options.GetTimeout()); err != nil {
		return err
	}
	instance, err := b.fileSandbox(ctx, id, "copy_from_local")
	if err != nil {
		return err
	}
	if err := instance.FileSystem.UploadFile(ctx, localPath, remotePath); err != nil {
		return mapProviderError(err, "copy_from_local")
	}
	return nil
}

func (b *backend) CopyToLocal(ctx context.Context, id, remotePath, localPath string, options *sandbox.FileCopyOptions) error {
	if err := wholeSecondFileTimeout("copy_to_local", options.GetTimeout()); err != nil {
		return err
	}
	instance, err := b.fileSandbox(ctx, id, "copy_to_local")
	if err != nil {
		return err
	}
	if _, err := instance.FileSystem.DownloadFile(ctx, remotePath, &localPath); err != nil {
		return mapProviderError(err, "copy_to_local")
	}
	return nil
}

func (b *backend) MakeDirectory(ctx context.Context, id, path string, options *sandbox.MakeDirectoryOptions) error {
	if err := wholeSecondFileTimeout("make_directory", options.GetTimeout()); err != nil {
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

func (b *backend) ListDirectory(ctx context.Context, id, path string, options *sandbox.ListDirectoryOptions) ([]*sandbox.FileInfo, error) {
	if err := wholeSecondFileTimeout("list_directory", options.GetTimeout()); err != nil {
		return nil, err
	}
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

func (b *backend) Stat(ctx context.Context, id, path string, options *sandbox.StatOptions) (*sandbox.FileInfo, error) {
	if err := wholeSecondFileTimeout("stat", options.GetTimeout()); err != nil {
		return nil, err
	}
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
	if err := wholeSecondFileTimeout("remove", options.GetTimeout()); err != nil {
		return err
	}
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

func wholeSecondFileTimeout(operation string, value *durationpb.Duration) error {
	timeout := sdk.OptionalDuration(value)
	if timeout == nil {
		return nil
	}
	if *timeout%time.Second != 0 {
		return sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "daytona",
			Operation: operation,
			Field:     sandbox.ConfigFieldTimeout,
			Message:   "sandbox-kit daytona: timeout must be a whole number of seconds",
		}, nil)
	}
	return nil
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
	if info.Mode != "" {
		result.Mode = sandbox.Value(info.Mode)
	}
	if !info.ModifiedTime.IsZero() {
		result.ModifiedAt = timestamppb.New(info.ModifiedTime)
	}
	return result, nil
}

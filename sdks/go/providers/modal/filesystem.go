package modal

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ sandbox.SandboxFiles = (*backend)(nil)

func (b *backend) ReadFile(ctx context.Context, id, path string, _ *sandbox.FileReadOptions) ([]byte, error) {
	instance, err := b.fileSandbox(ctx, id, "read_file")
	if err != nil {
		return nil, err
	}
	data, err := instance.Filesystem.ReadBytes(ctx, path, nil)
	if err != nil {
		return nil, mapFileError(err, "read_file", "")
	}
	return data, nil
}

func (b *backend) WriteFile(ctx context.Context, id, path string, data []byte, _ *sandbox.FileWriteOptions) error {
	instance, err := b.fileSandbox(ctx, id, "write_file")
	if err != nil {
		return err
	}
	if data == nil {
		data = []byte{}
	}
	if err := instance.Filesystem.WriteBytes(ctx, data, path, nil); err != nil {
		return mapFileError(err, "write_file", "")
	}
	return nil
}

func (b *backend) CopyFromLocal(ctx context.Context, id, localPath, remotePath string, _ *sandbox.FileCopyOptions) error {
	instance, err := b.fileSandbox(ctx, id, "copy_from_local")
	if err != nil {
		return err
	}
	if err := instance.Filesystem.CopyFromLocal(ctx, localPath, remotePath, nil); err != nil {
		return mapFileError(err, "copy_from_local", sandbox.FileFieldLocalPath)
	}
	return nil
}

func (b *backend) CopyToLocal(ctx context.Context, id, remotePath, localPath string, _ *sandbox.FileCopyOptions) error {
	instance, err := b.fileSandbox(ctx, id, "copy_to_local")
	if err != nil {
		return err
	}
	if err := instance.Filesystem.CopyToLocal(ctx, remotePath, localPath, nil); err != nil {
		return mapFileError(err, "copy_to_local", sandbox.FileFieldLocalPath)
	}
	return nil
}

func (b *backend) MakeDirectory(ctx context.Context, id, path string, options *sandbox.MakeDirectoryOptions) error {
	instance, err := b.fileSandbox(ctx, id, "make_directory")
	if err != nil {
		return err
	}
	if err := instance.Filesystem.MakeDirectory(ctx, path, makeDirectoryParams(options)); err != nil {
		return mapFileError(err, "make_directory", "")
	}
	return nil
}

func makeDirectoryParams(options *sandbox.MakeDirectoryOptions) *sdk.SandboxFilesystemMakeDirectoryParams {
	parents := true
	if options != nil && options.Parents != nil {
		parents = *options.Parents
	}
	return &sdk.SandboxFilesystemMakeDirectoryParams{CreateParents: &parents}
}

func (b *backend) ListDirectory(ctx context.Context, id, path string, _ *sandbox.ListDirectoryOptions) ([]*sandbox.FileInfo, error) {
	instance, err := b.fileSandbox(ctx, id, "list_directory")
	if err != nil {
		return nil, err
	}
	entries, err := instance.Filesystem.ListFiles(ctx, path, nil)
	if err != nil {
		return nil, mapFileError(err, "list_directory", "")
	}
	result := make([]*sandbox.FileInfo, 0, len(entries))
	for _, entry := range entries {
		result = append(result, modalFileInfo(entry))
	}
	return result, nil
}

func (b *backend) Stat(ctx context.Context, id, path string, _ *sandbox.StatOptions) (*sandbox.FileInfo, error) {
	instance, err := b.fileSandbox(ctx, id, "stat")
	if err != nil {
		return nil, err
	}
	info, err := instance.Filesystem.Stat(ctx, path, nil)
	if err != nil {
		return nil, mapFileError(err, "stat", "")
	}
	if info == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "modal",
			Operation: "stat",
			Field:     sandbox.FileFieldPath,
			Message:   "sandbox-kit modal: SDK returned no file metadata",
		}, nil)
	}
	return modalFileInfo(*info), nil
}

func (b *backend) Remove(ctx context.Context, id, path string, options *sandbox.RemoveOptions) error {
	instance, err := b.fileSandbox(ctx, id, "remove")
	if err != nil {
		return err
	}
	recursive := options != nil && options.GetRecursive()
	if err := instance.Filesystem.Remove(ctx, path, &sdk.SandboxFilesystemRemoveParams{Recursive: recursive}); err != nil {
		return mapFileError(err, "remove", "")
	}
	return nil
}

func (b *backend) fileSandbox(ctx context.Context, id, operation string) (*sdk.Sandbox, error) {
	instance, err := b.client.Sandboxes.FromID(ctx, id, nil)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	if instance == nil || instance.Filesystem == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "modal",
			Operation: operation,
			Message:   "sandbox-kit modal: SDK returned no filesystem",
		}, nil)
	}
	return instance, nil
}

func mapFileError(err error, operation, localField string) error {
	if err == nil {
		return nil
	}
	var pathError *os.PathError
	if errors.As(err, &pathError) {
		kind := sandbox.ErrorKindInvalidArgument
		if errors.Is(pathError.Err, os.ErrNotExist) {
			kind = sandbox.ErrorKindNotFound
		}
		field := localField
		if field == "" {
			field = sandbox.FileFieldLocalPath
		}
		return sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      kind,
			Provider:  "modal",
			Operation: operation,
			Field:     field,
			Message:   err.Error(),
		}, err)
	}
	return mapProviderError(err, operation)
}

func modalFileInfo(info sdk.FileInfo) *sandbox.FileInfo {
	kind := sandbox.FileTypeFile
	switch info.Type {
	case sdk.FileTypeDirectory:
		kind = sandbox.FileTypeDirectory
	case sdk.FileTypeSymlink:
		kind = sandbox.FileTypeSymlink
	}
	result := &sandbox.FileInfo{
		Name: info.Name,
		Path: info.Path,
		Type: kind,
		Size: info.Size,
		Mode: sandbox.Value(fmt.Sprintf("%04o", info.Mode&0o7777)),
	}
	if info.Permissions != "" {
		result.Permissions = sandbox.Value(info.Permissions)
	}
	if info.Owner != "" {
		result.Owner = sandbox.Value(info.Owner)
	}
	if info.Group != "" {
		result.Group = sandbox.Value(info.Group)
	}
	if info.SymlinkTarget != nil && *info.SymlinkTarget != "" {
		result.SymlinkTarget = info.SymlinkTarget
	}
	if !math.IsNaN(info.ModifiedTime) && !math.IsInf(info.ModifiedTime, 0) {
		whole := int64(info.ModifiedTime)
		frac := info.ModifiedTime - float64(whole)
		result.ModifiedAt = timestamppb.New(time.Unix(whole, int64(frac*1e9)).UTC())
	}
	return result
}

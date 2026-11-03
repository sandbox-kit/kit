package daytona

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	sdkerrors "github.com/daytona/clients/sdk-go/pkg/errors"
	"github.com/daytona/clients/toolbox-api-client-go"
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
	if data == nil {
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
	if err := ensureLocalParent(localPath, "copy_to_local"); err != nil {
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
	client, err := b.toolbox(ctx, id, "list_directory")
	if err != nil {
		return nil, err
	}
	entries, httpResp, err := client.FileSystemAPI.ListFiles(ctx).Path(path).Execute()
	if err != nil {
		return nil, mapProviderError(sdkerrors.ConvertToolboxError(err, httpResp), "list_directory")
	}
	result := make([]*sandbox.FileInfo, 0, len(entries))
	for i := range entries {
		info, err := daytonaFileInfo(&entries[i], "list_directory")
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}

func (b *backend) Stat(ctx context.Context, id, path string, _ *sandbox.StatOptions) (*sandbox.FileInfo, error) {
	client, err := b.toolbox(ctx, id, "stat")
	if err != nil {
		return nil, err
	}
	info, httpResp, err := client.FileSystemAPI.GetFileInfo(ctx).Path(path).Execute()
	if err != nil {
		return nil, mapProviderError(sdkerrors.ConvertToolboxError(err, httpResp), "stat")
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

func (b *backend) toolbox(ctx context.Context, id, operation string) (*toolbox.APIClient, error) {
	instance, err := b.fileSandbox(ctx, id, operation)
	if err != nil {
		return nil, err
	}
	if instance.ToolboxClient == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "daytona",
			Operation: operation,
			Message:   "sandbox-kit daytona: SDK returned no filesystem",
		}, nil)
	}
	return instance.ToolboxClient, nil
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

// localDestination rejects a path that is already a directory. A missing path
// is left for ensureLocalParent and Daytona's write.
func localDestination(path, operation string) error {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return localPathError(err, operation)
	}
	if info.IsDir() {
		return localDirectoryError(operation)
	}
	return nil
}

// ensureLocalParent creates missing local parents. Daytona's DownloadFile uses
// os.WriteFile, which does not create them.
func ensureLocalParent(path, operation string) error {
	parent := filepath.Dir(path)
	if parent == "." || parent == string(filepath.Separator) {
		return nil
	}
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return localPathError(err, operation)
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

func daytonaFileInfo(info *toolbox.FileInfo, operation string) (*sandbox.FileInfo, error) {
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
	if info.GetIsDir() {
		kind = sandbox.FileTypeDirectory
	}
	result := &sandbox.FileInfo{
		Name: info.GetName(),
		Path: info.GetPath(),
		Type: kind,
		Size: int64(info.GetSize()),
	}
	if mode := info.GetMode(); mode != "" {
		result.Mode = sandbox.Value(mode)
	}
	if permissions := info.GetPermissions(); permissions != "" {
		result.Permissions = sandbox.Value(permissions)
	}
	if owner := info.GetOwner(); owner != "" {
		result.Owner = sandbox.Value(owner)
	}
	if group := info.GetGroup(); group != "" {
		result.Group = sandbox.Value(group)
	}
	if modified := daytonaModifiedAt(info); modified != nil {
		result.ModifiedAt = modified
	}
	return result, nil
}

// modifiedAt wins when it is present. A fractional RFC3339 timestamp is kept.
// modTime is used only when modifiedAt is empty, matching the TypeScript mapper.
func daytonaModifiedAt(info *toolbox.FileInfo) *timestamppb.Timestamp {
	raw := info.GetModifiedAt()
	if raw == "" {
		raw = info.GetModTime()
	}
	if raw == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed)
}

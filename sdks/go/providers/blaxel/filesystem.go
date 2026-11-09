package blaxel

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/blaxel-ai/sdk-go/option"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ sandbox.SandboxFiles = (*backend)(nil)

func (b *backend) ReadFile(ctx context.Context, id, path string, _ *sandbox.FileReadOptions) ([]byte, error) {
	instance, err := b.fileSandbox(ctx, id, "read_file")
	if err != nil {
		return nil, err
	}
	data, err := instance.FS.ReadBinary(ctx, path)
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
	if _, err := instance.FS.WriteBinary(ctx, path, data, ""); err != nil {
		return mapProviderError(err, "write_file")
	}
	return nil
}

func (b *backend) CopyFromLocal(ctx context.Context, id, localPath, remotePath string, _ *sandbox.FileCopyOptions) error {
	if err := localSourceFile(localPath, "copy_from_local"); err != nil {
		return err
	}
	data, err := os.ReadFile(localPath)
	if err != nil {
		return localPathError(err, "copy_from_local")
	}
	instance, err := b.fileSandbox(ctx, id, "copy_from_local")
	if err != nil {
		return err
	}
	if _, err := instance.FS.WriteBinary(ctx, remotePath, data, ""); err != nil {
		return mapProviderError(err, "copy_from_local")
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
	if err := instance.FS.Download(ctx, remotePath, localPath); err != nil {
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
	if _, err := instance.FS.Mkdir(ctx, path, ""); err != nil {
		return mapProviderError(err, "make_directory")
	}
	return nil
}

func (b *backend) ListDirectory(ctx context.Context, id, path string, _ *sandbox.ListDirectoryOptions) ([]*sandbox.FileInfo, error) {
	instance, err := b.fileSandbox(ctx, id, "list_directory")
	if err != nil {
		return nil, err
	}
	listing, err := instance.FS.LS(ctx, path)
	if err != nil {
		return nil, mapProviderError(err, "list_directory")
	}
	result := make([]*sandbox.FileInfo, 0, len(listing.Files)+len(listing.Subdirectories))
	for _, entry := range listing.Files {
		info, err := blaxelFileInfo(entry, sandbox.FileTypeFile, "list_directory")
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	for _, entry := range listing.Subdirectories {
		result = append(result, &sandbox.FileInfo{
			Name: entry.Name,
			Path: entry.Path,
			Type: sandbox.FileTypeDirectory,
		})
	}
	return result, nil
}

func (b *backend) Stat(ctx context.Context, id, path string, _ *sandbox.StatOptions) (*sandbox.FileInfo, error) {
	instance, err := b.fileSandbox(ctx, id, "stat")
	if err != nil {
		return nil, err
	}
	fs := b.filesystemService(instance)
	res, err := fs.Get(ctx, path, blaxel.SandboxFilesystemGetParams{})
	if err != nil {
		return nil, mapProviderError(err, "stat")
	}
	if res == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: "stat",
			Field:     sandbox.FileFieldPath,
			Message:   "sandbox-kit blaxel: SDK returned no file metadata",
		}, nil)
	}
	if res.JSON.Files.Valid() || res.JSON.Subdirectories.Valid() {
		dir := res.AsDirectory()
		return &sandbox.FileInfo{
			Name: dir.Name,
			Path: dir.Path,
			Type: sandbox.FileTypeDirectory,
		}, nil
	}
	file := res.AsFilesystemReadWithContent()
	return blaxelFileInfo(blaxel.FilesystemRead{
		Name:         file.Name,
		Path:         file.Path,
		Size:         file.Size,
		Permissions:  file.Permissions,
		Owner:        file.Owner,
		Group:        file.Group,
		LastModified: file.LastModified,
	}, sandbox.FileTypeFile, "stat")
}

func (b *backend) Remove(ctx context.Context, id, path string, options *sandbox.RemoveOptions) error {
	instance, err := b.fileSandbox(ctx, id, "remove")
	if err != nil {
		return err
	}
	recursive := options != nil && options.GetRecursive()
	if _, err := instance.FS.RM(ctx, path, recursive); err != nil {
		return mapProviderError(err, "remove")
	}
	return nil
}

func (b *backend) fileSandbox(ctx context.Context, id, operation string) (*blaxel.SandboxInstance, error) {
	instance, err := b.client.Sandboxes.GetInstance(ctx, id)
	if err != nil {
		return nil, mapProviderError(err, operation)
	}
	if instance == nil || instance.FS == nil {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: operation,
			Message:   "sandbox-kit blaxel: SDK returned no filesystem",
		}, nil)
	}
	return instance, nil
}

func (b *backend) filesystemService(instance *blaxel.SandboxInstance) blaxel.SandboxFilesystemService {
	opts := append([]option.RequestOption{}, b.client.Options...)
	if instance != nil && instance.Metadata.URL != "" {
		opts = append(opts, option.WithBaseURL(instance.Metadata.URL))
	}
	return blaxel.NewSandboxFilesystemService(opts...)
}

func rejectParents(options *sandbox.MakeDirectoryOptions) error {
	if options == nil || options.Parents == nil || *options.Parents {
		return nil
	}
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      sandbox.ErrorKindUnsupported,
		Provider:  "blaxel",
		Operation: "make_directory",
		Field:     sandbox.FileFieldParents,
		Message:   "sandbox-kit blaxel: parents false is not available",
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
		Provider:  "blaxel",
		Operation: operation,
		Field:     sandbox.FileFieldLocalPath,
		Message:   "sandbox-kit blaxel: local path is a directory",
	}, nil)
}

func localPathError(err error, operation string) error {
	kind := sandbox.ErrorKindInvalidArgument
	if errors.Is(err, os.ErrNotExist) {
		kind = sandbox.ErrorKindNotFound
	}
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      kind,
		Provider:  "blaxel",
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

func blaxelFileInfo(entry blaxel.FilesystemRead, kind sandbox.FileType, operation string) (*sandbox.FileInfo, error) {
	if entry.Path == "" && entry.Name == "" {
		return nil, sandbox.NewError(&sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "blaxel",
			Operation: operation,
			Field:     sandbox.FileFieldPath,
			Message:   "sandbox-kit blaxel: SDK returned no file metadata",
		}, nil)
	}
	result := &sandbox.FileInfo{
		Name: entry.Name,
		Path: entry.Path,
		Type: kind,
		Size: entry.Size,
	}
	if entry.Permissions != "" {
		result.Permissions = sandbox.Value(entry.Permissions)
	}
	if entry.Owner != "" {
		result.Owner = sandbox.Value(entry.Owner)
	}
	if entry.Group != "" {
		result.Group = sandbox.Value(entry.Group)
	}
	if modified := parseOptionalTime(entry.LastModified); modified != nil {
		result.ModifiedAt = modified
	}
	return result, nil
}

func parseOptionalTime(raw string) *timestamppb.Timestamp {
	if raw == "" {
		return nil
	}
	parsed, err := parseBlaxelTime(raw)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed)
}

func parseBlaxelTime(raw string) (time.Time, error) {
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		time.RFC1123,
		time.RFC1123Z,
	} {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, errors.New("unrecognized timestamp")
}

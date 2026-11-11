package blaxel

import (
	"context"
	"errors"
	"os"
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/blaxel-ai/sdk-go/option"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ sandbox.SandboxFiles = (*backend)(nil)

func (b *backend) ReadFile(ctx context.Context, id, path string, _ *sandbox.FileReadOptions) ([]byte, error) {
	return withFS(b, ctx, fsCall{ID: id, Operation: "read_file"}, func(fs *blaxel.SandboxInstanceFSService) ([]byte, error) {
		return fs.ReadBinary(ctx, path)
	})
}

func (b *backend) WriteFile(ctx context.Context, id, path string, data []byte, _ *sandbox.FileWriteOptions) error {
	if data == nil {
		data = []byte{}
	}
	_, err := withFS(b, ctx, fsCall{ID: id, Operation: "write_file"}, func(fs *blaxel.SandboxInstanceFSService) (struct{}, error) {
		_, err := fs.WriteBinary(ctx, path, data, "")
		return struct{}{}, err
	})
	return err
}

func (b *backend) CopyFromLocal(ctx context.Context, id, localPath, remotePath string, _ *sandbox.FileCopyOptions) error {
	if err := localSourceFile(localPath, "copy_from_local"); err != nil {
		return err
	}
	data, err := os.ReadFile(localPath)
	if err != nil {
		return localPathError(err, "copy_from_local")
	}
	_, err = withFS(b, ctx, fsCall{ID: id, Operation: "copy_from_local"}, func(fs *blaxel.SandboxInstanceFSService) (struct{}, error) {
		_, err := fs.WriteBinary(ctx, remotePath, data, "")
		return struct{}{}, err
	})
	return err
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
	_, err := withFS(b, ctx, fsCall{ID: id, Operation: "make_directory"}, func(fs *blaxel.SandboxInstanceFSService) (struct{}, error) {
		_, err := fs.Mkdir(ctx, path, "")
		return struct{}{}, err
	})
	return err
}

func (b *backend) ListDirectory(ctx context.Context, id, path string, _ *sandbox.ListDirectoryOptions) ([]*sandbox.FileInfo, error) {
	return withFS(b, ctx, fsCall{ID: id, Operation: "list_directory"}, func(fs *blaxel.SandboxInstanceFSService) ([]*sandbox.FileInfo, error) {
		listing, err := fs.LS(ctx, path)
		if err != nil {
			return nil, err
		}
		result := make([]*sandbox.FileInfo, 0, len(listing.Files)+len(listing.Subdirectories))
		// Directories first, then files — matches TypeScript Blaxel listing order.
		for _, entry := range listing.Subdirectories {
			result = append(result, &sandbox.FileInfo{
				Name: entry.Name,
				Path: entry.Path,
				Type: sandbox.FileTypeDirectory,
			})
		}
		for _, entry := range listing.Files {
			info, err := blaxelFileInfo(entry, sandbox.FileTypeFile, "list_directory")
			if err != nil {
				return nil, err
			}
			result = append(result, info)
		}
		return result, nil
	})
}

func (b *backend) Stat(ctx context.Context, id, path string, _ *sandbox.StatOptions) (*sandbox.FileInfo, error) {
	instance, err := b.fileSandbox(ctx, id, "stat")
	if err != nil {
		return nil, err
	}
	// SandboxInstanceFSService has Read/LS/… but no public Get/Stat. LS always
	// calls AsDirectory(), so a file path would look like an empty directory.
	// Rebuild the lower SandboxFilesystemService with the instance base URL
	// (same construction NewInstance uses for instance.FS) and call Get.
	fs := b.filesystemService(instance)
	res, err := fs.Get(ctx, path, blaxel.SandboxFilesystemGetParams{})
	if err != nil {
		return nil, mapFilesystemError(err, "stat")
	}
	if res == nil {
		return nil, invalidResponse("stat", sandbox.FileFieldPath, "SDK returned no file metadata")
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
	recursive := options != nil && options.GetRecursive()
	_, err := withFS(b, ctx, fsCall{ID: id, Operation: "remove"}, func(fs *blaxel.SandboxInstanceFSService) (struct{}, error) {
		_, err := fs.RM(ctx, path, recursive)
		return struct{}{}, err
	})
	return err
}

type fsCall struct {
	ID        string
	Operation string
}

func withFS[T any](
	b *backend,
	ctx context.Context,
	call fsCall,
	fn func(*blaxel.SandboxInstanceFSService) (T, error),
) (T, error) {
	var zero T
	instance, err := b.fileSandbox(ctx, call.ID, call.Operation)
	if err != nil {
		return zero, err
	}
	result, err := fn(instance.FS)
	if err != nil {
		// Preserve SandboxErrors raised inside the callback (e.g. list mapping).
		var local *sandbox.Error
		if errors.As(err, &local) {
			return zero, err
		}
		return zero, mapFilesystemError(err, call.Operation)
	}
	return result, nil
}

// mapFilesystemError attaches Field=path on NotFound so portable callers can
// branch the same way as TypeScript Blaxel Stat/list failures.
func mapFilesystemError(err error, operation string) error {
	mapped := mapProviderError(err, operation)
	var detail *sandbox.Error
	if !errors.As(mapped, &detail) || detail.Kind != sandbox.ErrorKindNotFound || detail.Field != "" {
		return mapped
	}
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:           detail.Kind,
		Provider:       detail.Provider,
		Operation:      detail.Operation,
		Field:          sandbox.FileFieldPath,
		Message:        detail.Message,
		StatusCode:     detail.StatusCode,
		ProviderCode:   detail.ProviderCode,
		ProviderSource: detail.ProviderSource,
	}, err)
}

func (b *backend) fileSandbox(ctx context.Context, id, operation string) (*blaxel.SandboxInstance, error) {
	instance, err := b.getInstance(ctx, id, operation)
	if err != nil {
		return nil, err
	}
	if instance == nil || instance.FS == nil {
		return nil, invalidResponse(operation, "", "SDK returned no filesystem")
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
	return kitError(kitFault{
		Kind:      sandbox.ErrorKindUnsupported,
		Operation: "make_directory",
		Field:     sandbox.FileFieldParents,
		Message:   "parents false is not available",
	})
}

func blaxelFileInfo(entry blaxel.FilesystemRead, kind sandbox.FileType, operation string) (*sandbox.FileInfo, error) {
	if entry.Path == "" && entry.Name == "" {
		return nil, invalidResponse(operation, sandbox.FileFieldPath, "SDK returned no file metadata")
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

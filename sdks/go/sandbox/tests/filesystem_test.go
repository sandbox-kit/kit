package sandbox_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type fileBackend struct {
	read func(context.Context, string, string, *sandbox.FileReadOptions) ([]byte, error)
}

func (*fileBackend) Name() string                { return "files" }
func (*fileBackend) Close(context.Context) error { return nil }
func (*fileBackend) CloseSandbox(context.Context, string) error {
	return nil
}
func (*fileBackend) Create(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{Id: "created"}}, nil
}
func (b *fileBackend) ReadFile(ctx context.Context, id, path string, options *sandbox.FileReadOptions) ([]byte, error) {
	return b.read(ctx, id, path, options)
}
func (*fileBackend) WriteFile(context.Context, string, string, []byte, *sandbox.FileWriteOptions) error {
	return nil
}
func (*fileBackend) CopyFromLocal(context.Context, string, string, string, *sandbox.FileCopyOptions) error {
	return nil
}
func (*fileBackend) CopyToLocal(context.Context, string, string, string, *sandbox.FileCopyOptions) error {
	return nil
}
func (*fileBackend) MakeDirectory(context.Context, string, string, *sandbox.MakeDirectoryOptions) error {
	return nil
}
func (*fileBackend) ListDirectory(context.Context, string, string, *sandbox.ListDirectoryOptions) ([]*sandbox.FileInfo, error) {
	return nil, nil
}
func (*fileBackend) Stat(context.Context, string, string, *sandbox.StatOptions) (*sandbox.FileInfo, error) {
	return &sandbox.FileInfo{}, nil
}
func (*fileBackend) Remove(context.Context, string, string, *sandbox.RemoveOptions) error {
	return nil
}

type fileProvider struct{ backend sandbox.Backend }

func (*fileProvider) Name() string { return "files" }
func (p *fileProvider) NewClient(*sandbox.Config) (sandbox.Backend, error) {
	return p.backend, nil
}

func TestFilesystemNilHandle(t *testing.T) {
	_, err := (*sandbox.Sandbox)(nil).Filesystem().ReadFile(context.Background(), "/tmp/note.txt", nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindFailedPrecondition {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFilesystemRejectsRelativePathBeforeProvider(t *testing.T) {
	calls := 0
	backend := &fileBackend{read: func(context.Context, string, string, *sandbox.FileReadOptions) ([]byte, error) {
		calls++
		return nil, nil
	}}
	client, err := sandbox.NewClient(&fileProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instance.Filesystem().ReadFile(context.Background(), "relative.txt", nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.FileFieldPath || detail.Operation != "read_file" {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 0 {
		t.Fatal("provider was called for a relative path")
	}
}

func TestFilesystemReadsAbsolutePath(t *testing.T) {
	backend := &fileBackend{read: func(_ context.Context, id, path string, _ *sandbox.FileReadOptions) ([]byte, error) {
		if id != "created" || path != "/workspace/note.txt" {
			t.Fatalf("unexpected read %s %s", id, path)
		}
		return []byte("hello"), nil
	}}
	client, err := sandbox.NewClient(&fileProvider{backend: backend}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	data, err := instance.Filesystem().ReadFile(context.Background(), "/workspace/note.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected bytes %q", data)
	}
}

func TestFilesystemCopyRequiresLocalPath(t *testing.T) {
	client, err := sandbox.NewClient(&fileProvider{backend: &fileBackend{}}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Filesystem().CopyFromLocal(context.Background(), "", "/workspace/note.txt", nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.FileFieldLocalPath {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFilesystemUnsupportedWithoutCapability(t *testing.T) {
	client, err := sandbox.NewClient(&lifecycleProvider{backend: &lifecycleBackend{}}, &sandbox.Config{})
	if err != nil {
		t.Fatal(err)
	}
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = instance.Filesystem().Stat(context.Background(), "/workspace", nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Operation != "stat" {
		t.Fatalf("unexpected error: %v", err)
	}
}

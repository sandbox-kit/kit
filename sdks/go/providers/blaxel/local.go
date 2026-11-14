package blaxel

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

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
	return kitError(kitFault{
		Kind:      sandbox.ErrorKindInvalidArgument,
		Operation: operation,
		Field:     sandbox.FileFieldLocalPath,
		Message:   "local path is a directory",
	})
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

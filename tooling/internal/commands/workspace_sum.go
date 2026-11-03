package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// checkWorkspaceSum fails when go.work.sum contains checksums that loading the
// workspace packages does not write. go work sync and go list -m all load the
// full requirement graph, fetch go.mod files the build never imports, and
// record those checksums. Module tidy does not remove them.
func checkWorkspaceSum(ctx context.Context, root, goBinary string) (err error) {
	workPath := filepath.Join(root, "go.work")
	if _, statErr := os.Stat(workPath); statErr != nil {
		if os.IsNotExist(statErr) {
			return nil
		}
		return fmt.Errorf("inspect go.work: %w", statErr)
	}
	sumPath := filepath.Join(root, "go.work.sum")
	committed, readErr := os.ReadFile(sumPath)
	if readErr != nil {
		return fmt.Errorf("read go.work.sum: %w", readErr)
	}
	modules, parseErr := workspaceModules(workPath)
	if parseErr != nil {
		return parseErr
	}
	if removeErr := os.Remove(sumPath); removeErr != nil {
		return fmt.Errorf("reset go.work.sum: %w", removeErr)
	}
	defer func() {
		if writeErr := os.WriteFile(sumPath, committed, 0o644); writeErr != nil {
			err = errors.Join(err, fmt.Errorf("restore go.work.sum: %w", writeErr))
		}
	}()

	rebuilt, loadErr := reloadWorkspaceSum(ctx, root, goBinary, modules)
	if loadErr != nil {
		return loadErr
	}
	extra, missing := workspaceSumDrift(committed, rebuilt)
	if len(extra) == 0 && len(missing) == 0 {
		return nil
	}
	return errors.New(formatWorkspaceSumDrift(extra, missing))
}

func workspaceModules(workPath string) ([]string, error) {
	data, err := os.ReadFile(workPath)
	if err != nil {
		return nil, fmt.Errorf("read go.work: %w", err)
	}
	var modules []string
	inUse := false
	for line := range strings.SplitSeq(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if !inUse {
			if trimmed == "use (" {
				inUse = true
			}
			continue
		}
		if trimmed == ")" {
			break
		}
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}
		modules = append(modules, strings.TrimPrefix(trimmed, "./"))
	}
	if len(modules) == 0 {
		return nil, errors.New("go.work has no use block")
	}
	return modules, nil
}

func reloadWorkspaceSum(ctx context.Context, root, goBinary string, modules []string) ([]byte, error) {
	env := environWithout("GOWORK")
	for _, module := range modules {
		cmd := exec.CommandContext(ctx, goBinary, "list", "./...")
		cmd.Dir = filepath.Join(root, filepath.FromSlash(module))
		cmd.Env = env
		cmd.Stdout = io.Discard
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			detail := strings.TrimSpace(stderr.String())
			if detail == "" {
				return nil, fmt.Errorf("load packages in %s: %w", module, err)
			}
			return nil, fmt.Errorf("load packages in %s: %w\n%s", module, err, detail)
		}
	}
	rebuilt, err := os.ReadFile(filepath.Join(root, "go.work.sum"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read rebuilt go.work.sum: %w", err)
	}
	return rebuilt, nil
}

func environWithout(key string) []string {
	prefix := key + "="
	env := os.Environ()
	filtered := make([]string, 0, len(env))
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func workspaceSumDrift(committed, rebuilt []byte) (extra, missing []string) {
	have := sumEntries(committed)
	want := sumEntries(rebuilt)
	for entry := range have {
		if _, ok := want[entry]; !ok {
			extra = append(extra, entry)
		}
	}
	for entry := range want {
		if _, ok := have[entry]; !ok {
			missing = append(missing, entry)
		}
	}
	sort.Strings(extra)
	sort.Strings(missing)
	return extra, missing
}

func sumEntries(data []byte) map[string]struct{} {
	entries := make(map[string]struct{})
	for line := range bytes.SplitSeq(data, []byte("\n")) {
		fields := bytes.Fields(line)
		if len(fields) < 2 {
			continue
		}
		entries[string(fields[0])+" "+string(fields[1])] = struct{}{}
	}
	return entries
}

func formatWorkspaceSumDrift(extra, missing []string) string {
	var b strings.Builder
	if len(extra) > 0 {
		b.WriteString("go.work.sum records modules that loading the workspace packages does not fetch:\n")
		for _, entry := range extra {
			b.WriteString("  ")
			b.WriteString(entry)
			b.WriteString("\n")
		}
		b.WriteString("\ngo work sync and go list -m all read the full requirement graph, including test requirements of google.golang.org/grpc and the modules those requirements name. Package loads do not fetch them, and go mod tidy does not remove the checksums. Delete go.work.sum and rerun sandbox-kit test go.\n")
	}
	if len(missing) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString("go.work.sum is missing checksums from loading the workspace packages:\n")
		for _, entry := range missing {
			b.WriteString("  ")
			b.WriteString(entry)
			b.WriteString("\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

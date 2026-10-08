package commands

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type recordingRunner struct {
	steps  []step
	failAt int
	err    error
}

func (r *recordingRunner) run(_ context.Context, task step, _, _ io.Writer) error {
	r.steps = append(r.steps, task)
	if r.failAt == len(r.steps) {
		return r.err
	}
	return nil
}

func repositoryFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range []string{"proto/kit/codegen/v1/options.proto", "tooling/go.mod"} {
		file := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte("fixture"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestGenerationUsesAndCleansExternalTempDirectory(t *testing.T) {
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "build failure"}[failure], func(t *testing.T) {
			repo := repositoryFixture(t)
			runner := &recordingRunner{}
			if failure {
				runner.failAt, runner.err = 1, errors.New("build failure")
			}
			command := newRootCommand(runner)
			command.SetOut(io.Discard)
			command.SetErr(io.Discard)
			command.SetArgs([]string{"generate", "go", "--repo", repo, "--go-binary", "custom-go", "--protoc", "custom-protoc"})
			err := command.ExecuteContext(context.Background())
			if failure && !errors.Is(err, runner.err) {
				t.Fatalf("expected build failure, got %v", err)
			}
			if !failure && err != nil {
				t.Fatal(err)
			}
			if len(runner.steps) == 0 || runner.steps[0].name != "custom-go" {
				t.Fatal("generation did not use the configured executable")
			}
			if !failure && len(runner.steps) != 7 {
				t.Fatalf("expected bootstrap and three SDK generations, got %d steps", len(runner.steps))
			}
			toolsDir := filepath.Dir(runner.steps[0].args[2])
			relative, err := filepath.Rel(repo, toolsDir)
			if err != nil || (relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator))) {
				t.Fatal("temporary tools were created inside the repository")
			}
			if _, err := os.Stat(toolsDir); !os.IsNotExist(err) {
				t.Fatal("temporary tools directory was not cleaned up")
			}
			if _, err := os.Stat(filepath.Join(repo, "work")); !os.IsNotExist(err) {
				t.Fatal("generation created a repository-local work folder")
			}
		})
	}
}

func TestVerificationIncludesEveryGoModuleAndExample(t *testing.T) {
	repo := repositoryFixture(t)
	runner := &recordingRunner{}
	command := newRootCommand(runner)
	command.SetArgs([]string{"test", "go", "--repo", repo})
	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(runner.steps) != 6 {
		t.Fatalf("expected five module tests and example help, got %d", len(runner.steps))
	}
	for i, module := range []string{"sdks/go/sandbox", "tooling", "sdks/go/providers/modal", "sdks/go/providers/daytona", "examples/go"} {
		if runner.steps[i].dir != filepath.Join(repo, filepath.FromSlash(module)) || runner.steps[i].args[0] != "test" {
			t.Fatalf("incorrect verification step for %s", module)
		}
	}
	if strings.Join(runner.steps[5].args, " ") != "run . --help" {
		t.Fatal("verification must only run example help, not a credentialed provider")
	}
}

func TestHelpAndInvalidArgumentsDoNotExecuteTools(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"generate", "go", "extra"}} {
		runner := &recordingRunner{}
		command := newRootCommand(runner)
		command.SetOut(&bytes.Buffer{})
		command.SetErr(io.Discard)
		command.SetArgs(args)
		err := command.ExecuteContext(context.Background())
		if args[0] == "--help" && err != nil {
			t.Fatal(err)
		}
		if args[0] != "--help" && err == nil {
			t.Fatal("unexpected positional arguments were accepted")
		}
		if len(runner.steps) != 0 {
			t.Fatal("help/invalid arguments executed tools")
		}
	}
}

func TestFindsRootFromToolingDirectory(t *testing.T) {
	repo := repositoryFixture(t)
	t.Chdir(filepath.Join(repo, "tooling"))
	root, err := resolveRoot("")
	if err != nil || root != repo {
		t.Fatalf("expected repository %q, got %q: %v", repo, root, err)
	}
}

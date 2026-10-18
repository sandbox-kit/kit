package commands

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
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
	targets := `version: 1
languages:
  go:
    targets:
      - {source: shared.proto, output: sdks/go/sandbox, module: example.com/sandbox}
      - {source: modal.proto, output: sdks/go/providers/modal, module: example.com/modal, provider: modal}
      - {source: daytona.proto, output: sdks/go/providers/daytona, module: example.com/daytona, provider: daytona}
  typescript:
    targets:
      - {source: modal.proto, output: sdks/typescript/providers/modal/src, module: '@example/modal', provider: modal}
      - {source: daytona.proto, output: sdks/typescript/providers/daytona/src, module: '@example/daytona', provider: daytona}
`
	path := filepath.Join(root, "specs", "targets.yaml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(targets), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestTypeScriptGenerationUsesDeclaredTargets(t *testing.T) {
	repo := repositoryFixture(t)
	runner := &recordingRunner{}
	command := newRootCommand(runner)
	command.SetOut(io.Discard)
	command.SetErr(io.Discard)
	command.SetArgs([]string{"generate", "typescript", "--repo", repo, "--go-binary", "custom-go", "--protoc", "custom-protoc"})
	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(runner.steps) != 3 {
		t.Fatalf("expected plugin build and two provider generations, got %d", len(runner.steps))
	}
	if runner.steps[0].name != "custom-go" || runner.steps[1].name != "custom-protoc" {
		t.Fatal("TypeScript generation did not use configured executables")
	}
	for _, task := range runner.steps[1:] {
		if !strings.HasPrefix(task.args[2], "--plugin=protoc-gen-kit-typescript=") {
			t.Fatalf("missing TypeScript plugin binding: %#v", task.args)
		}
	}
}

func TestGenerationSkipsEmptyTargets(t *testing.T) {
	repo := repositoryFixture(t)
	empty := `version: 1
languages:
  go:
    targets: []
  typescript:
    targets: []
`
	if err := os.WriteFile(filepath.Join(repo, "specs", "targets.yaml"), []byte(empty), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"generate", "go", "--repo", repo},
		{"generate", "typescript", "--repo", repo},
	} {
		runner := &recordingRunner{}
		command := newRootCommand(runner)
		out := &bytes.Buffer{}
		command.SetOut(out)
		command.SetErr(io.Discard)
		command.SetArgs(args)
		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		if len(runner.steps) != 0 {
			t.Fatalf("%v executed tools with empty targets: %#v", args, runner.steps)
		}
		if !strings.Contains(out.String(), "handwritten SDK sources are authoritative") {
			t.Fatalf("%v missing empty-target message: %q", args, out.String())
		}
	}
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
			if !failure && len(runner.steps) != 6 {
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
	if len(runner.steps) != len(goTestModules) {
		t.Fatalf("expected %d module tests, got %d", len(goTestModules), len(runner.steps))
	}
	for i, module := range goTestModules {
		if runner.steps[i].dir != filepath.Join(repo, filepath.FromSlash(module)) || runner.steps[i].args[0] != "test" {
			t.Fatalf("incorrect verification step for %s", module)
		}
	}
}

func TestProtoVerificationUsesBufAndOptionalBaseline(t *testing.T) {
	repo := repositoryFixture(t)
	runner := &recordingRunner{}
	command := newRootCommand(runner)
	command.SetArgs([]string{"test", "proto", "--repo", repo, "--buf", "custom-buf", "--against", ".git#branch=main,subdir=proto"})
	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"format", "--diff", "--exit-code"},
		{"lint"},
		{"build"},
		{"breaking", "--against", ".git#branch=main,subdir=proto"},
	}
	if len(runner.steps) != len(want) {
		t.Fatalf("expected %d protobuf checks, got %d", len(want), len(runner.steps))
	}
	for i, args := range want {
		if runner.steps[i].dir != repo || runner.steps[i].name != "custom-buf" || !reflect.DeepEqual(runner.steps[i].args, args) {
			t.Fatalf("protobuf check %d = %#v, want %q %#v", i, runner.steps[i], "custom-buf", args)
		}
	}
}

func TestTypeScriptVerificationBuildsBeforeTests(t *testing.T) {
	repo := repositoryFixture(t)
	for _, path := range []string{
		"examples/typescript/daytona/create-daytona-sandbox",
		"examples/typescript/modal/create-modal-sandbox",
	} {
		dir := filepath.Join(repo, filepath.FromSlash(path))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte("{}\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runner := &recordingRunner{}
	command := newRootCommand(runner)
	command.SetArgs([]string{"test", "typescript", "--repo", repo, "--npm-binary", "custom-npm"})
	if err := command.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(runner.steps) != 8 {
		t.Fatalf("unexpected TypeScript verification steps: %#v", runner.steps)
	}
	for _, task := range runner.steps {
		if task.name != "custom-npm" {
			t.Fatalf("TypeScript verification did not use configured npm: %#v", runner.steps)
		}
	}
	if strings.Join(runner.steps[0].args, " ") != "run build" || strings.Join(runner.steps[1].args, " ") != "test" || strings.Join(runner.steps[7].args, " ") != "run build" {
		t.Fatalf("TypeScript verification order is incorrect: %#v", runner.steps)
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

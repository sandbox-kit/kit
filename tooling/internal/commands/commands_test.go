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
	for _, path := range []string{"proto/kit/sandbox/v1/client.proto", "tooling/go.mod"} {
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

func TestContractsGenerationUsesBuf(t *testing.T) {
	t.Run("contracts-go", func(t *testing.T) {
		repo := repositoryFixture(t)
		runner := &recordingRunner{}
		command := newRootCommand(runner)
		command.SetOut(io.Discard)
		command.SetErr(io.Discard)
		command.SetArgs([]string{"generate", "contracts-go", "--repo", repo, "--buf", "custom-buf"})
		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatal(err)
		}
		if len(runner.steps) != 1 || runner.steps[0].name != "custom-buf" {
			t.Fatalf("expected buf generate, got %#v", runner.steps)
		}
		if strings.Join(runner.steps[0].args, " ") != "generate" {
			t.Fatalf("unexpected buf args: %#v", runner.steps[0].args)
		}
		if _, err := os.Stat(filepath.Join(repo, "work")); !os.IsNotExist(err) {
			t.Fatal("generation created a repository-local work folder")
		}
	})
	t.Run("contracts-typescript", func(t *testing.T) {
		repo := repositoryFixture(t)
		runner := &recordingRunner{}
		command := newRootCommand(runner)
		command.SetOut(io.Discard)
		command.SetErr(io.Discard)
		command.SetArgs([]string{"generate", "contracts-typescript", "--repo", repo, "--buf", "custom-buf", "--npm-binary", "custom-npm"})
		if err := command.ExecuteContext(context.Background()); err != nil {
			t.Fatal(err)
		}
		if len(runner.steps) != 2 {
			t.Fatalf("expected buf generate then npm build, got %#v", runner.steps)
		}
		if runner.steps[0].name != "custom-buf" || strings.Join(runner.steps[0].args, " ") != "generate" {
			t.Fatalf("unexpected buf step: %#v", runner.steps[0])
		}
		if runner.steps[1].name != "custom-npm" || strings.Join(runner.steps[1].args, " ") != "run build" {
			t.Fatalf("unexpected npm step: %#v", runner.steps[1])
		}
		if _, err := os.Stat(filepath.Join(repo, "work")); !os.IsNotExist(err) {
			t.Fatal("generation created a repository-local work folder")
		}
	})
}

func TestContractsGenerationSurfacesBufFailure(t *testing.T) {
	repo := repositoryFixture(t)
	runner := &recordingRunner{failAt: 1, err: errors.New("buf failure")}
	command := newRootCommand(runner)
	command.SetOut(io.Discard)
	command.SetErr(io.Discard)
	command.SetArgs([]string{"generate", "contracts-go", "--repo", repo, "--buf", "custom-buf"})
	if err := command.ExecuteContext(context.Background()); !errors.Is(err, runner.err) {
		t.Fatalf("expected buf failure, got %v", err)
	}
}

func TestWorkspaceSumCheckSkipsRepoWithoutGoWork(t *testing.T) {
	repo := repositoryFixture(t)
	if err := checkWorkspaceSum(context.Background(), repo, "go"); err != nil {
		t.Fatal(err)
	}
}

func TestWorkspaceModulesReadsUseBlock(t *testing.T) {
	repo := repositoryFixture(t)
	work := "go 1.26.0\n\nuse (\n\t./sdks/go/sandbox\n\t// comment\n\n\t./tooling\n)\n"
	if err := os.WriteFile(filepath.Join(repo, "go.work"), []byte(work), 0o644); err != nil {
		t.Fatal(err)
	}
	modules, err := workspaceModules(filepath.Join(repo, "go.work"))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"sdks/go/sandbox", "tooling"}
	if !reflect.DeepEqual(modules, want) {
		t.Fatalf("modules = %#v, want %#v", modules, want)
	}
}

func TestWorkspaceSumDriftReportsGraphOnlyChecksums(t *testing.T) {
	committed := []byte("cel.dev/expr v0.25.2/go.mod h1:def\ngonum.org/v1/plot v0.16.0/go.mod h1:abc\n")
	rebuilt := []byte("cel.dev/expr v0.25.2/go.mod h1:def\n")
	extra, missing := workspaceSumDrift(committed, rebuilt)
	if len(missing) != 0 || len(extra) != 1 || extra[0] != "gonum.org/v1/plot v0.16.0/go.mod" {
		t.Fatalf("extra=%v missing=%v", extra, missing)
	}
	message := formatWorkspaceSumDrift(extra, missing)
	if !strings.Contains(message, "gonum.org/v1/plot v0.16.0/go.mod") || !strings.Contains(message, "go work sync") {
		t.Fatalf("drift message = %q", message)
	}
}

func TestWorkspaceSumMatchesPackageLoads(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if !isRepository(root) {
		t.Fatalf("repository root not found at %s", root)
	}
	if err := checkWorkspaceSum(context.Background(), root, "go"); err != nil {
		t.Fatal(err)
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
		"examples/typescript/daytona/configuration/create-daytona-sandbox",
		"examples/typescript/modal/configuration/create-modal-sandbox",
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
	if len(runner.steps) != 11 {
		t.Fatalf("unexpected TypeScript verification steps: %#v", runner.steps)
	}
	for _, task := range runner.steps {
		if task.name != "custom-npm" {
			t.Fatalf("TypeScript verification did not use configured npm: %#v", runner.steps)
		}
	}
	last := len(runner.steps) - 1
	if strings.Join(runner.steps[0].args, " ") != "run build" || strings.Join(runner.steps[1].args, " ") != "run build" || strings.Join(runner.steps[2].args, " ") != "test" || strings.Join(runner.steps[last].args, " ") != "run build" {
		t.Fatalf("TypeScript verification order is incorrect: %#v", runner.steps)
	}
}

func TestHelpAndInvalidArgumentsDoNotExecuteTools(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"generate", "contracts-go", "extra"}} {
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

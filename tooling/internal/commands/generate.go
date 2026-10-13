package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sandbox-kit/kit/tooling/internal/spec"

	"github.com/spf13/cobra"
)

func newGenerateCommand(config *settings, runner runner) *cobra.Command {
	command := &cobra.Command{Use: "generate", Short: "Generate SDK source"}
	command.AddCommand(&cobra.Command{
		Use: "go", Short: "Generate Go SDKs from protobuf and YAML declarations",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			return generateGo(cmd, root, config, runner)
		},
	})
	return command
}

func generateGo(cmd *cobra.Command, root string, config *settings, runner runner) (result error) {
	// Only this newly allocated directory is removed; caller-owned paths are untouched.
	toolsDir, err := os.MkdirTemp("", "sandbox-kit-tools-")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(toolsDir); err != nil {
			result = errors.Join(result, fmt.Errorf("clean temporary generator tools: %w", err))
		}
	}()

	toolingDir := filepath.Join(root, "tooling")
	goPlugin := filepath.Join(toolsDir, "protoc-gen-go")
	kitPlugin := filepath.Join(toolsDir, "protoc-gen-kit-go")
	steps := []step{
		{toolingDir, config.goBinary, []string{"build", "-o", goPlugin, "google.golang.org/protobuf/cmd/protoc-gen-go"}},
		{root, config.protoc, []string{
			"-I", filepath.Join(root, "proto"), "--plugin=protoc-gen-go=" + goPlugin,
			"--go_out=" + toolingDir, "--go_opt=module=github.com/sandbox-kit/kit/tooling",
			"kit/codegen/v1/options.proto",
		}},
		{toolingDir, config.goBinary, []string{"build", "-o", kitPlugin, "./cmd/protoc-gen-kit-go"}},
	}
	targets, err := spec.LoadTargets(filepath.Join(root, "specs"))
	if err != nil {
		return err
	}
	goTargets, ok := targets.Languages["go"]
	if !ok {
		return fmt.Errorf("generation target manifest has no Go targets")
	}
	for _, target := range goTargets.Targets {
		steps = append(steps, step{root, config.protoc, []string{
			"-I", filepath.Join(root, "proto"), "--plugin=protoc-gen-kit-go=" + kitPlugin,
			"--kit-go_out=" + filepath.Join(root, filepath.FromSlash(target.Output)),
			"--kit-go_opt=module=" + target.Module, target.Source,
		}})
	}
	for _, task := range steps {
		if err := runner.run(cmd.Context(), task, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), "Generated Go client and providers.")
	return err
}

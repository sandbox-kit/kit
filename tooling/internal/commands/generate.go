package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
	for _, target := range []struct{ directory, module, source string }{
		{"sdks/go/sandbox", "github.com/sandbox-kit/kit/sdks/go/sandbox", "kit/sandbox/v1/sandbox.proto"},
		{"sdks/go/sandbox", "github.com/sandbox-kit/kit/sdks/go/sandbox", "kit/sandbox/v1/client.proto"},
		{"sdks/go/providers/modal", "github.com/sandbox-kit/kit/sdks/go/providers/modal", "kit/providers/modal/v1/provider.proto"},
		{"sdks/go/providers/daytona", "github.com/sandbox-kit/kit/sdks/go/providers/daytona", "kit/providers/daytona/v1/provider.proto"},
	} {
		steps = append(steps, step{root, config.protoc, []string{
			"-I", filepath.Join(root, "proto"), "--plugin=protoc-gen-kit-go=" + kitPlugin,
			"--kit-go_out=" + filepath.Join(root, filepath.FromSlash(target.directory)),
			"--kit-go_opt=module=" + target.module, target.source,
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

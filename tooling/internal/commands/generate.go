package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newGenerateCommand(config *settings, runner runner) *cobra.Command {
	command := &cobra.Command{Use: "generate", Short: "Generate protobuf contract packages"}
	command.AddCommand(&cobra.Command{
		Use: "contracts-go", Short: "Generate Go protobuf contracts",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			return generateGoContracts(cmd, root, config, runner)
		},
	})
	command.AddCommand(&cobra.Command{
		Use: "contracts-typescript", Short: "Generate TypeScript protobuf contracts",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			return generateTypeScriptContracts(cmd, root, config, runner)
		},
	})
	return command
}

// generateGoContracts uses the official protobuf generator for portable data
// contracts. Runtime SDK behavior is intentionally not emitted here.
func generateGoContracts(cmd *cobra.Command, root string, config *settings, runner runner) (result error) {
	toolsDir, err := os.MkdirTemp("", "sandbox-kit-contract-tools-")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(toolsDir); err != nil {
			result = errors.Join(result, fmt.Errorf("clean temporary generator tools: %w", err))
		}
	}()

	plugin := filepath.Join(toolsDir, "protoc-gen-go")
	contracts := filepath.Join(root, "sdks", "go", "contracts")
	steps := []step{
		{filepath.Join(root, "tooling"), config.goBinary, []string{"build", "-o", plugin, "google.golang.org/protobuf/cmd/protoc-gen-go"}},
		{root, config.protoc, []string{
			"-I", filepath.Join(root, "proto"), "--plugin=protoc-gen-go=" + plugin,
			"--go_out=" + contracts, "--go_opt=module=github.com/sandbox-kit/kit/sdks/go/contracts",
			"--go_opt=Mkit/sandbox/v1/sandbox.proto=github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1;sandboxv1",
			"--go_opt=Mkit/sandbox/v1/client.proto=github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1;sandboxv1",
			"--go_opt=Mkit/codegen/v1/options.proto=github.com/sandbox-kit/kit/sdks/go/contracts/internal/codegen/v1;codegenv1",
			"kit/codegen/v1/options.proto", "kit/sandbox/v1/sandbox.proto", "kit/sandbox/v1/client.proto",
		}},
	}
	for _, task := range steps {
		if err := runner.run(cmd.Context(), task, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), "Generated Go protobuf contracts.")
	return err
}

// generateTypeScriptContracts invokes the pinned Protobuf-ES plugin installed
// by the contracts package. It deliberately emits only portable protobuf data.
func generateTypeScriptContracts(cmd *cobra.Command, root string, config *settings, runner runner) error {
	contracts := filepath.Join(root, "sdks", "typescript", "contracts")
	if err := runner.run(cmd.Context(), step{contracts, config.npmBinary, []string{"run", "generate"}}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
		return err
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), "Generated TypeScript protobuf contracts.")
	return err
}

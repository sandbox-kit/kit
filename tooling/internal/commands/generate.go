package commands

import (
	"fmt"
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
			return generateContracts(cmd, root, config, runner, false, "Generated Go protobuf contracts.")
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
			return generateContracts(cmd, root, config, runner, true, "Generated TypeScript protobuf contracts.")
		},
	})
	return command
}

// generateContracts runs buf generate for both language plugins declared in
// buf.gen.yaml. Portable create/config rules come from Protovalidate options in
// proto/; SDK runtimes evaluate them with language Protovalidate libraries.
// When buildTypeScript is true, compile contracts/src into contracts/dist so
// Protovalidate descriptors are what @sandbox-kit/contracts exports.
func generateContracts(cmd *cobra.Command, root string, config *settings, runner runner, buildTypeScript bool, done string) error {
	if err := runner.run(cmd.Context(), step{root, config.bufBinary, []string{"generate"}}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
		return err
	}
	if buildTypeScript {
		contracts := filepath.Join(root, "sdks", "typescript", "contracts")
		if err := runner.run(cmd.Context(), step{contracts, config.npmBinary, []string{"run", "build"}}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), done)
	return err
}

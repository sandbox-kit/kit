package commands

import (
	"fmt"

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
			return generateContracts(cmd, root, config, runner, "Generated Go protobuf contracts.")
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
			return generateContracts(cmd, root, config, runner, "Generated TypeScript protobuf contracts.")
		},
	})
	return command
}

// generateContracts runs buf generate for both language plugins declared in
// buf.gen.yaml. Portable create/config rules come from Protovalidate options in
// proto/; SDK runtimes evaluate them with language Protovalidate libraries.
func generateContracts(cmd *cobra.Command, root string, config *settings, runner runner, done string) error {
	if err := runner.run(cmd.Context(), step{root, config.bufBinary, []string{"generate"}}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
		return err
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), done)
	return err
}

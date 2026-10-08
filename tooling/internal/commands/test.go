package commands

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

func newTestCommand(config *settings, runner runner) *cobra.Command {
	command := &cobra.Command{Use: "test", Short: "Verify SDK modules and examples"}
	command.AddCommand(&cobra.Command{
		Use: "go", Short: "Test Go SDKs, tooling, and the runnable example", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			for _, module := range []string{"sdks/go/sandbox", "tooling", "sdks/go/providers/modal", "sdks/go/providers/daytona", "examples/go"} {
				if err := runner.run(cmd.Context(), step{
					filepath.Join(root, filepath.FromSlash(module)), config.goBinary, []string{"test", "./..."},
				}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
					return err
				}
			}
			return runner.run(cmd.Context(), step{
				filepath.Join(root, "examples", "go"), config.goBinary, []string{"run", ".", "--help"},
			}, cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	})
	return command
}

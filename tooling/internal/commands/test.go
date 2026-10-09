package commands

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

func newTestCommand(config *settings, runner runner) *cobra.Command {
	command := &cobra.Command{Use: "test", Short: "Verify SDK modules and examples"}
	command.AddCommand(&cobra.Command{
		Use: "go", Short: "Test Go SDKs, tooling, and the example projects", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			for _, module := range []string{"sdks/go/sandbox", "tooling", "sdks/go/providers/modal", "sdks/go/providers/daytona", "examples/go/create-modal-sandbox", "examples/go/create-daytona-sandbox", "examples/go/create-modal-with-resources", "examples/go/create-daytona-with-policies"} {
				if err := runner.run(cmd.Context(), step{
					filepath.Join(root, filepath.FromSlash(module)), config.goBinary, []string{"test", "./..."},
				}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
					return err
				}
			}
			return nil
		},
	})
	return command
}

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
			for _, module := range []string{"sdks/go/sandbox", "tooling", "sdks/go/providers/modal", "sdks/go/providers/daytona", "examples/go/modal/create-modal-sandbox", "examples/go/daytona/create-daytona-sandbox", "examples/go/modal/create-modal-with-resources", "examples/go/daytona/create-daytona-with-policies", "examples/go/daytona/create-daytona-from-snapshot", "examples/go/daytona/create-daytona-from-image", "examples/go/daytona/create-daytona-ephemeral", "examples/go/daytona/create-daytona-linked", "examples/go/modal/create-modal-with-runtime", "examples/go/modal/create-modal-gvisor", "examples/go/modal/create-modal-linux-vm", "examples/go/modal/create-modal-with-placement", "examples/go/modal/create-modal-with-labels", "examples/go/modal/create-modal-with-network", "examples/go/modal/create-modal-idle", "examples/go/modal/create-modal-ready", "examples/go/modal/create-modal-verbose", "examples/go/modal/create-modal-with-identity", "examples/go/modal/handle-modal-errors", "examples/go/daytona/handle-daytona-errors"} {
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

package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var goTestModules = []string{
	"sdks/go/sandbox", "tooling", "sdks/go/providers/modal", "sdks/go/providers/daytona",
	"examples/go/modal/create-modal-sandbox", "examples/go/daytona/create-daytona-sandbox",
	"examples/go/modal/create-modal-with-resources", "examples/go/daytona/create-daytona-with-policies",
	"examples/go/daytona/create-daytona-from-snapshot", "examples/go/daytona/create-daytona-from-image",
	"examples/go/daytona/create-daytona-with-runtime",
	"examples/go/daytona/create-daytona-ephemeral", "examples/go/daytona/create-daytona-linked",
	"examples/go/modal/create-modal-with-runtime", "examples/go/modal/create-modal-gvisor",
	"examples/go/modal/create-modal-linux-vm", "examples/go/modal/create-modal-with-placement",
	"examples/go/modal/create-modal-with-labels", "examples/go/modal/create-modal-with-network",
	"examples/go/modal/create-modal-idle", "examples/go/modal/create-modal-ready",
	"examples/go/modal/create-modal-verbose", "examples/go/modal/create-modal-with-identity",
	"examples/go/modal/handle-modal-errors", "examples/go/daytona/handle-daytona-errors",
	"examples/go/daytona/stop-daytona-sandbox", "examples/go/modal/terminate-modal-sandbox",
	"examples/go/daytona/get-daytona-sandbox", "examples/go/daytona/refresh-daytona-sandbox",
	"examples/go/daytona/start-daytona-sandbox", "examples/go/daytona/pause-daytona-sandbox",
	"examples/go/daytona/archive-daytona-sandbox", "examples/go/daytona/delete-daytona-sandbox",
	"examples/go/modal/attach-modal-sandbox", "examples/go/modal/request-modal-termination",
	"examples/go/daytona/close-daytona-handle", "examples/go/modal/close-modal-handle",
}

var typescriptTestModules = []struct {
	path string
	test bool
}{
	{"sdks/typescript/sandbox", true},
	{"sdks/typescript/providers/daytona", true},
	{"sdks/typescript/providers/modal", true},
}

func typescriptModules(root string) ([]struct {
	path string
	test bool
}, error) {
	modules := append([]struct {
		path string
		test bool
	}{}, typescriptTestModules...)
	examplesRoot := filepath.Join(root, "examples", "typescript")
	providers, err := os.ReadDir(examplesRoot)
	if err != nil {
		return nil, fmt.Errorf("read TypeScript examples: %w", err)
	}
	var examplePaths []string
	for _, provider := range providers {
		if !provider.IsDir() {
			continue
		}
		projects, err := os.ReadDir(filepath.Join(examplesRoot, provider.Name()))
		if err != nil {
			return nil, fmt.Errorf("read TypeScript %s examples: %w", provider.Name(), err)
		}
		for _, project := range projects {
			if !project.IsDir() {
				continue
			}
			manifest := filepath.Join(examplesRoot, provider.Name(), project.Name(), "package.json")
			if _, err := os.Stat(manifest); err == nil {
				examplePaths = append(examplePaths, filepath.ToSlash(filepath.Join("examples", "typescript", provider.Name(), project.Name())))
			} else if !os.IsNotExist(err) {
				return nil, fmt.Errorf("inspect TypeScript example %s: %w", project.Name(), err)
			}
		}
	}
	sort.Strings(examplePaths)
	for _, path := range examplePaths {
		modules = append(modules, struct {
			path string
			test bool
		}{path: path})
	}
	return modules, nil
}

func newTestCommand(config *settings, runner runner) *cobra.Command {
	command := &cobra.Command{Use: "test", Short: "Verify SDK modules and examples"}
	command.AddCommand(&cobra.Command{
		Use: "go", Short: "Test Go SDKs, tooling, and the example projects", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			for _, module := range goTestModules {
				if err := runner.run(cmd.Context(), step{
					filepath.Join(root, filepath.FromSlash(module)), config.goBinary, []string{"test", "./..."},
				}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
					return err
				}
			}
			return nil
		},
	})
	command.AddCommand(&cobra.Command{
		Use: "typescript", Short: "Build and test TypeScript SDK modules", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			modules, err := typescriptModules(root)
			if err != nil {
				return err
			}
			for _, module := range modules {
				dir := filepath.Join(root, filepath.FromSlash(module.path))
				if err := runner.run(cmd.Context(), step{dir, config.npmBinary, []string{"run", "build"}}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
					return err
				}
				if module.test {
					if err := runner.run(cmd.Context(), step{dir, config.npmBinary, []string{"test"}}, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
						return err
					}
				}
			}
			return nil
		},
	})
	return command
}

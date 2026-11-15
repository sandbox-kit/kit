package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var goTestModules = []string{
	"sdks/go/sandbox", "tooling", "sdks/go/providers/modal", "sdks/go/providers/daytona", "sdks/go/providers/blaxel",
	"examples/go/modal/configuration/create-modal-sandbox", "examples/go/daytona/configuration/create-daytona-sandbox",
	"examples/go/modal/configuration/create-modal-with-resources", "examples/go/daytona/configuration/create-daytona-with-policies",
	"examples/go/daytona/configuration/create-daytona-from-snapshot", "examples/go/daytona/configuration/create-daytona-from-image",
	"examples/go/daytona/configuration/create-daytona-with-runtime",
	"examples/go/daytona/configuration/create-daytona-ephemeral", "examples/go/daytona/configuration/create-daytona-linked",
	"examples/go/daytona/configuration/create-daytona-with-kvm", "examples/go/daytona/configuration/create-daytona-with-labels",
	"examples/go/daytona/configuration/create-daytona-with-network", "examples/go/daytona/configuration/create-daytona-with-secrets",
	"examples/go/daytona/configuration/create-daytona-with-ttl",
	"examples/go/modal/configuration/create-modal-with-runtime", "examples/go/modal/configuration/create-modal-gvisor",
	"examples/go/modal/configuration/create-modal-linux-vm", "examples/go/modal/configuration/create-modal-with-placement",
	"examples/go/modal/configuration/create-modal-with-labels", "examples/go/modal/configuration/create-modal-with-network",
	"examples/go/modal/configuration/create-modal-idle", "examples/go/modal/configuration/create-modal-ready",
	"examples/go/modal/configuration/create-modal-verbose", "examples/go/modal/configuration/create-modal-with-identity",
	"examples/go/modal/errors/handle-modal-errors", "examples/go/daytona/errors/handle-daytona-errors",
	"examples/go/daytona/lifecycle/stop-daytona-sandbox", "examples/go/modal/lifecycle/terminate-modal-sandbox",
	"examples/go/daytona/lifecycle/get-daytona-sandbox", "examples/go/daytona/lifecycle/refresh-daytona-sandbox",
	"examples/go/daytona/lifecycle/start-daytona-sandbox", "examples/go/daytona/lifecycle/pause-daytona-sandbox",
	"examples/go/daytona/lifecycle/archive-daytona-sandbox", "examples/go/daytona/lifecycle/delete-daytona-sandbox",
	"examples/go/daytona/lifecycle/force-stop-daytona-sandbox", "examples/go/daytona/lifecycle/request-daytona-deletion",
	"examples/go/daytona/lifecycle/start-daytona-with-timeout",
	"examples/go/modal/lifecycle/attach-modal-sandbox", "examples/go/modal/lifecycle/request-modal-termination",
	"examples/go/daytona/lifecycle/close-daytona-handle", "examples/go/modal/lifecycle/close-modal-handle",
	"examples/go/modal/configuration/create-modal-client-timeout", "examples/go/modal/configuration/create-modal-tcp-ready",
	"examples/go/modal/configuration/create-modal-with-allowlists", "examples/go/modal/configuration/create-modal-with-ports",
	"examples/go/modal/configuration/create-modal-with-secrets", "examples/go/modal/lifecycle/terminate-modal-with-timeout",
	"examples/go/daytona/filesystem/copy-file-from-daytona",
	"examples/go/daytona/filesystem/copy-file-into-daytona",
	"examples/go/daytona/filesystem/list-daytona-directory",
	"examples/go/daytona/filesystem/make-daytona-directory",
	"examples/go/daytona/filesystem/read-daytona-file",
	"examples/go/daytona/filesystem/remove-daytona-file",
	"examples/go/daytona/filesystem/stat-daytona-file",
	"examples/go/daytona/filesystem/write-daytona-file",
	"examples/go/daytona/process/run-daytona-check-exit",
	"examples/go/daytona/process/run-daytona-command",
	"examples/go/daytona/process/run-daytona-shell-command",
	"examples/go/daytona/process/run-daytona-with-cwd",
	"examples/go/daytona/process/run-daytona-with-env",
	"examples/go/daytona/process/run-daytona-with-timeout",
	"examples/go/modal/filesystem/copy-file-from-modal",
	"examples/go/modal/filesystem/copy-file-into-modal",
	"examples/go/modal/filesystem/list-modal-directory",
	"examples/go/modal/filesystem/make-modal-directory",
	"examples/go/modal/filesystem/read-modal-file",
	"examples/go/modal/filesystem/remove-modal-file",
	"examples/go/modal/filesystem/stat-modal-file",
	"examples/go/modal/filesystem/write-modal-file",
	"examples/go/modal/process/run-modal-check-exit",
	"examples/go/modal/process/run-modal-command",
	"examples/go/modal/process/run-modal-shell-command",
	"examples/go/modal/process/run-modal-with-cwd",
	"examples/go/modal/process/run-modal-with-env",
	"examples/go/modal/process/run-modal-with-timeout",
	"examples/go/blaxel/configuration/create-blaxel-sandbox",
	"examples/go/blaxel/configuration/create-blaxel-with-resources",
	"examples/go/blaxel/configuration/create-blaxel-with-labels",
	"examples/go/blaxel/configuration/create-blaxel-with-ttl",
	"examples/go/blaxel/configuration/create-blaxel-idle",
	"examples/go/blaxel/configuration/create-blaxel-client-timeout",
	"examples/go/blaxel/configuration/create-blaxel-with-ports",
	"examples/go/blaxel/configuration/create-blaxel-with-allowlists",
	"examples/go/blaxel/errors/handle-blaxel-errors",
	"examples/go/blaxel/lifecycle/get-blaxel-sandbox",
	"examples/go/blaxel/lifecycle/refresh-blaxel-sandbox",
	"examples/go/blaxel/lifecycle/delete-blaxel-sandbox",
	"examples/go/blaxel/lifecycle/close-blaxel-handle",
	"examples/go/blaxel/filesystem/copy-file-from-blaxel",
	"examples/go/blaxel/filesystem/copy-file-into-blaxel",
	"examples/go/blaxel/filesystem/list-blaxel-directory",
	"examples/go/blaxel/filesystem/make-blaxel-directory",
	"examples/go/blaxel/filesystem/read-blaxel-file",
	"examples/go/blaxel/filesystem/remove-blaxel-file",
	"examples/go/blaxel/filesystem/stat-blaxel-file",
	"examples/go/blaxel/filesystem/write-blaxel-file",
	"examples/go/blaxel/process/run-blaxel-check-exit",
	"examples/go/blaxel/process/run-blaxel-command",
	"examples/go/blaxel/process/run-blaxel-shell-command",
	"examples/go/blaxel/process/run-blaxel-with-cwd",
	"examples/go/blaxel/process/run-blaxel-with-env",
	"examples/go/blaxel/process/run-blaxel-with-timeout",
}

var typescriptTestModules = []struct {
	path string
	test bool
}{
	{"sdks/typescript/contracts", false},
	{"sdks/typescript/sandbox", true},
	{"sdks/typescript/providers/daytona", true},
	{"sdks/typescript/providers/modal", true},
	{"sdks/typescript/providers/blaxel", true},
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
		areas, err := os.ReadDir(filepath.Join(examplesRoot, provider.Name()))
		if err != nil {
			return nil, fmt.Errorf("read TypeScript %s examples: %w", provider.Name(), err)
		}
		for _, area := range areas {
			if !area.IsDir() {
				continue
			}
			projects, err := os.ReadDir(filepath.Join(examplesRoot, provider.Name(), area.Name()))
			if err != nil {
				return nil, fmt.Errorf("read TypeScript %s %s examples: %w", provider.Name(), area.Name(), err)
			}
			for _, project := range projects {
				if !project.IsDir() {
					continue
				}
				manifest := filepath.Join(examplesRoot, provider.Name(), area.Name(), project.Name(), "package.json")
				if _, err := os.Stat(manifest); err == nil {
					examplePaths = append(examplePaths, filepath.ToSlash(filepath.Join("examples", "typescript", provider.Name(), area.Name(), project.Name())))
				} else if !os.IsNotExist(err) {
					return nil, fmt.Errorf("inspect TypeScript example %s: %w", project.Name(), err)
				}
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
	var against string
	proto := &cobra.Command{
		Use: "proto", Short: "Check protobuf formatting, lint, compilation, and optional compatibility", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			steps := []step{
				{root, config.bufBinary, []string{"format", "--diff", "--exit-code"}},
				{root, config.bufBinary, []string{"lint"}},
				{root, config.bufBinary, []string{"build"}},
			}
			if against != "" {
				steps = append(steps, step{root, config.bufBinary, []string{"breaking", "--against", against}})
			}
			for _, task := range steps {
				if err := runner.run(cmd.Context(), task, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
					return err
				}
			}
			return nil
		},
	}
	proto.Flags().StringVar(&against, "against", "", "Optional Buf or Git input used as the breaking-change baseline")
	command.AddCommand(proto)
	command.AddCommand(&cobra.Command{
		Use: "go", Short: "Test Go SDKs, tooling, and the example projects", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			if err := checkWorkspaceSum(cmd.Context(), root, config.goBinary); err != nil {
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

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
	command.AddCommand(&cobra.Command{
		Use: "typescript", Short: "Generate TypeScript SDKs from protobuf and YAML declarations",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := resolveRoot(config.repoRoot)
			if err != nil {
				return err
			}
			return generateTypeScript(cmd, root, config, runner)
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
			"kit/sandbox/v1/sandbox.proto",
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

func generateTypeScript(cmd *cobra.Command, root string, config *settings, runner runner) (result error) {
	toolsDir, err := os.MkdirTemp("", "sandbox-kit-typescript-tools-")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(toolsDir); err != nil {
			result = errors.Join(result, fmt.Errorf("clean temporary generator tools: %w", err))
		}
	}()

	plugin := filepath.Join(toolsDir, "protoc-gen-kit-typescript")
	steps := []step{{filepath.Join(root, "tooling"), config.goBinary, []string{"build", "-o", plugin, "./cmd/protoc-gen-kit-typescript"}}}
	targets, err := spec.LoadTargets(filepath.Join(root, "specs"))
	if err != nil {
		return err
	}
	group, ok := targets.Languages["typescript"]
	if !ok {
		return fmt.Errorf("generation target manifest has no TypeScript targets")
	}
	for _, target := range group.Targets {
		steps = append(steps, step{root, config.protoc, []string{
			"-I", filepath.Join(root, "proto"), "--plugin=protoc-gen-kit-typescript=" + plugin,
			"--kit-typescript_out=" + filepath.Join(root, filepath.FromSlash(target.Output)), target.Source,
		}})
	}
	for _, task := range steps {
		if err := runner.run(cmd.Context(), task, cmd.OutOrStdout(), cmd.ErrOrStderr()); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), "Generated TypeScript SDK sources.")
	return err
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

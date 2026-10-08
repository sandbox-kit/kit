// Package commands provides reusable Cobra commands for SDK development.
package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type settings struct {
	repoRoot string
	goBinary string
	protoc   string
}

// NewRootCommand returns an independent command tree with no global flag state.
// Command construction performs no builds, network calls, or filesystem writes.
func NewRootCommand() *cobra.Command {
	return newRootCommand(execRunner{})
}

func newRootCommand(runner runner) *cobra.Command {
	config := &settings{}
	root := &cobra.Command{
		Use: "sandbox-kit", Short: "Generate and verify Sandbox Kit SDKs",
		SilenceUsage: true, SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&config.repoRoot, "repo", "", "Repository root (discovered from the current directory by default)")
	root.PersistentFlags().StringVar(&config.goBinary, "go-binary", "go", "Go executable")
	root.PersistentFlags().StringVar(&config.protoc, "protoc", "protoc", "Protobuf compiler executable")
	root.AddCommand(newGenerateCommand(config, runner), newTestCommand(config, runner))
	return root
}

func resolveRoot(configured string) (string, error) {
	start := configured
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	root, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if isRepository(root) {
			return root, nil
		}
		parent := filepath.Dir(root)
		if configured != "" || parent == root {
			return "", fmt.Errorf("Sandbox Kit repository not found from %q; set --repo", start)
		}
		root = parent
	}
}

func isRepository(root string) bool {
	for _, path := range []string{"proto/kit/codegen/v1/options.proto", "tooling/go.mod"} {
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil || info.IsDir() {
			return false
		}
	}
	return true
}

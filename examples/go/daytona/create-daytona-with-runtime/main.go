// Create one Daytona sandbox with explicit toolbox language and user settings.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() (result error) {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("copy .env.example to .env and fill in its values: %w", err)
	}
	config := &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.AuthConfig_ApiKey{ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")}}},
		Timeout: sandbox.Duration(2 * time.Minute),
	}
	if endpoint := os.Getenv("DAYTONA_API_URL"); endpoint != "" {
		config.Endpoint = sandbox.Value(endpoint)
	}
	if target := os.Getenv("DAYTONA_TARGET"); target != "" {
		config.Region = sandbox.Value(target)
	}
	client, err := sandbox.NewClient(daytona.New(), config)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(context.Background()); result == nil {
			result = err
		}
	}()

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{
			Snapshot: &sandbox.SnapshotSource{Reference: "daytona-small"},
		},
		Runtime: &sandbox.RuntimeConfig{
			Language: sandbox.Value("python"),
			User:     sandbox.Value("daytona"),
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("Daytona Python sandbox created: %s (%s)\n", instance.ID(), instance.Info().GetIsolation().GetKind())
	return nil
}

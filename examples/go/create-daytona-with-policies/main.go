// Create a Daytona sandbox with auto-pause disabled and delayed deletion.
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
		handleError(os.Stderr, err)
		os.Exit(1)
	}
}

func clientConfig() sandbox.Config {
	config := sandbox.Config{Provider: daytona.New(), Auth: &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")}}, Timeout: sandbox.Value(2 * time.Minute)}
	if endpoint := os.Getenv("DAYTONA_API_URL"); endpoint != "" {
		config.Endpoint = sandbox.Value(endpoint)
	}
	if target := os.Getenv("DAYTONA_TARGET"); target != "" {
		config.Region = sandbox.Value(target)
	}
	return config
}

func createOptions() *sandbox.CreateOptions {
	return &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{
		IdlePause:     &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled},
		StoppedDelete: &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(10 * time.Minute)},
	}}
}

func run() (result error) {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("copy .env.example to .env and set credentials: %w", err)
	}
	client, err := sandbox.NewClient(clientConfig())
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := client.Close(context.Background()); closeErr != nil {
			if result == nil {
				result = closeErr
			} else {
				handleError(os.Stderr, closeErr)
			}
		}
	}()
	instance, err := client.Create(context.Background(), createOptions())
	if err != nil {
		return err
	}
	fmt.Printf("Daytona sandbox created: %s (auto-pause disabled, deletion 10 minutes after stopping)\n", instance.ID())
	info := instance.Info()
	fmt.Printf("ID origin: %s; resource origin: %s\n", info.Origins["id"], info.Origins["resources"])
	return nil
}

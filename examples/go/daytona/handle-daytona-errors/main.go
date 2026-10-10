// Demonstrate common error handling with the Daytona integration.
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

func run() (result error) {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("copy .env.example to .env and fill in its values: %w", err)
	}
	config := sandbox.Config{
		Provider: daytona.New(),
		Auth:     &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")}},
		Timeout:  sandbox.Value(2 * time.Minute),
	}
	if endpoint := os.Getenv("DAYTONA_API_URL"); endpoint != "" {
		config.Endpoint = sandbox.Value(endpoint)
	}
	if target := os.Getenv("DAYTONA_TARGET"); target != "" {
		config.Region = sandbox.Value(target)
	}
	client, err := sandbox.NewClient(config)
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

	// Invalid CPU demonstrates how an application handles configuration errors.
	_, err = client.Create(context.Background(), &sandbox.CreateOptions{
		Resources: &sandbox.Resources{CPUCores: sandbox.Value(-1.0)},
	})
	if err != nil {
		handleError(os.Stdout, err)
		return nil
	}
	fmt.Fprintln(os.Stdout, "Sandbox created.")
	return nil
}

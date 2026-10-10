// Create one Daytona sandbox from a public image with explicit resources.
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
	client, err := newClient()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := client.Close(context.Background()); closeErr != nil {
			if result == nil {
				result = closeErr
			} else {
				fmt.Fprintln(os.Stderr, closeErr)
			}
		}
	}()
	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "ubuntu:22.04"}},
		Resources: &sandbox.Resources{
			CPUCores:  sandbox.Value(1.0),
			MemoryMiB: sandbox.Value(uint64(1024)),
			DiskMiB:   sandbox.Value(uint64(3072)),
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("Daytona sandbox created from ubuntu:22.04: %s\n", instance.ID())
	return nil
}

func newClient() (*sandbox.Client, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("copy .env.example to .env and fill in its values: %w", err)
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
	return sandbox.NewClient(config)
}

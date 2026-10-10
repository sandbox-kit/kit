// Create a Modal sandbox with fractional CPU, memory, and a bounded lifetime.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	if err := run(); err != nil {
		handleError(os.Stderr, err)
		os.Exit(1)
	}
}

func clientConfig() sandbox.Config {
	environment := os.Getenv("MODAL_ENVIRONMENT")
	if environment == "" {
		environment = "main"
	}
	return sandbox.Config{
		Provider: modal.New(),
		Auth:     &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: os.Getenv("MODAL_TOKEN_ID"), Secret: os.Getenv("MODAL_TOKEN_SECRET")}},
		Scope:    &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME")), Environment: sandbox.Value(environment)},
		Timeout:  sandbox.Value(2 * time.Minute),
	}
}

func createOptions() *sandbox.CreateOptions {
	return &sandbox.CreateOptions{
		Source:    &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}},
		Resources: &sandbox.Resources{CPUCores: sandbox.Value(0.5), MemoryMiB: sandbox.Value(uint64(512))},
		Lifetime:  &sandbox.LifetimePolicy{MaximumLifetime: sandbox.Value(5 * time.Minute)},
	}
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
	fmt.Printf("Modal sandbox created: %s (requested 0.5 CPU cores, 512 MiB, 5-minute lifetime)\n", instance.ID())
	info := instance.Info()
	fmt.Printf("ID origin: %s; allocated resources reported: %t\n", info.Origins["id"], info.Resources != nil)
	return nil
}

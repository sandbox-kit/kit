// Demonstrate common error handling with the Modal integration.
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

func run() (result error) {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("copy .env.example to .env and fill in its values: %w", err)
	}
	environment := os.Getenv("MODAL_ENVIRONMENT")
	if environment == "" {
		environment = "main"
	}
	client, err := sandbox.NewClient(sandbox.Config{
		Provider: modal.New(),
		Auth: &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{
			ID: os.Getenv("MODAL_TOKEN_ID"), Secret: os.Getenv("MODAL_TOKEN_SECRET"),
		}},
		Scope: &sandbox.Scope{
			AppName:     sandbox.Value(os.Getenv("MODAL_APP_NAME")),
			Environment: sandbox.Value(environment),
		},
		Timeout: sandbox.Value(2 * time.Minute),
	})
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

	// Modal does not support a disk allocation override at creation.
	_, err = client.Create(context.Background(), &sandbox.CreateOptions{
		Source:    &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}},
		Resources: &sandbox.Resources{DiskMiB: sandbox.Value(uint64(1024))},
	})
	if err != nil {
		handleError(os.Stdout, err)
		return nil
	}
	fmt.Fprintln(os.Stdout, "Sandbox created.")
	return nil
}

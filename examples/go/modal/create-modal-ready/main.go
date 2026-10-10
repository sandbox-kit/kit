// Create one Modal sandbox and wait until its readiness command succeeds.
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
		Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}},
		Runtime: &sandbox.RuntimeConfig{
			Entrypoint: []string{"sleep", "300"},
		},
		Readiness: &sandbox.ReadinessProbe{Command: &sandbox.CommandProbe{Argv: []string{"true"}}},
		Provisioning: &sandbox.ProvisioningOptions{
			Timeout: sandbox.Value(2 * time.Minute),
			WaitFor: sandbox.WaitConditionReady,
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("Modal sandbox ready: %s\n", instance.ID())
	return nil
}

func newClient() (*sandbox.Client, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("copy .env.example to .env and fill in its values: %w", err)
	}
	environment := os.Getenv("MODAL_ENVIRONMENT")
	if environment == "" {
		environment = "main"
	}
	return sandbox.NewClient(sandbox.Config{
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
}

// Attach to and terminate one existing Modal sandbox through Sandbox Kit.
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
			AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME")), Environment: sandbox.Value(environment),
		},
		Timeout: sandbox.Value(2 * time.Minute),
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(context.Background()); result == nil {
			result = err
		}
	}()
	instance, err := client.Attach(context.Background(), os.Getenv("MODAL_SANDBOX_ID"))
	if err != nil {
		return err
	}
	if !instance.Supports(sandbox.LifecycleCapabilityTerminate) {
		return fmt.Errorf("sandbox %s does not support terminate", instance.ID())
	}
	exit, err := instance.Terminate(context.Background(), &sandbox.TerminateOptions{Wait: sandbox.Value(true)})
	if err != nil {
		return err
	}
	fmt.Printf("Modal sandbox terminated: %s (exit %d)\n", instance.ID(), exit.GetExitCode())
	return nil
}

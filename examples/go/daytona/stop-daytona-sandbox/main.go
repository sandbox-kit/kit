// Retrieve and stop one existing Daytona sandbox through Sandbox Kit.
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
		return fmt.Errorf("load .env: %w", err)
	}
	client, err := sandbox.NewClient(sandbox.Config{Provider: daytona.New(), Auth: &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")}}, Timeout: sandbox.Value(2 * time.Minute)})
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(context.Background()); result == nil {
			result = err
		}
	}()
	instance, err := client.Get(context.Background(), os.Getenv("DAYTONA_SANDBOX_ID"))
	if err != nil {
		return err
	}
	if !instance.Supports(sandbox.LifecycleCapabilityStop) {
		return fmt.Errorf("sandbox %s does not support stop", instance.ID())
	}
	if err := instance.Stop(context.Background(), nil); err != nil {
		return err
	}
	fmt.Printf("Daytona sandbox stopped: %s (%s)\n", instance.ID(), instance.Info().State)
	return nil
}

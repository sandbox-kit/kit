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
	instance, err := client.Get(context.Background(), os.Getenv("DAYTONA_SANDBOX_ID"))
	if err != nil {
		return err
	}
	if !instance.Supports(sandbox.LifecycleCapabilityDelete) {
		return fmt.Errorf("sandbox %s does not support delete", instance.ID())
	}
	if err := instance.Delete(context.Background(), &sandbox.DeleteOptions{Wait: sandbox.Value(true)}); err != nil {
		return err
	}
	fmt.Printf("Daytona sandbox deleted: %s (terminal=%t)\n", instance.ID(), instance.Info().GetTerminal())
	return nil
}

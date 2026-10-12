package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"os"
	"time"
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
	if !instance.Supports(sandbox.LifecycleCapabilityDelete) {
		return fmt.Errorf("sandbox %s does not support delete", instance.ID())
	}
	if err := instance.Delete(context.Background(), &sandbox.DeleteOptions{Wait: sandbox.Value(true)}); err != nil {
		return err
	}
	fmt.Printf("Daytona sandbox deleted: %s (terminal=%t)\n", instance.ID(), instance.Info().GetTerminal())
	return nil
}

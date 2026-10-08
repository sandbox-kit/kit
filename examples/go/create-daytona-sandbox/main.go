// Authenticate and create one Daytona sandbox through Sandbox Kit.
package main

import (
	"context"
	"errors"
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
	client, err := sandbox.NewClient(sandbox.Config{
		Provider: daytona.New(),
		Auth:     &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")}},
		Timeout:  sandbox.Value(2 * time.Minute),
	})
	if err != nil {
		return err
	}
	defer func() { result = errors.Join(result, client.Close(context.Background())) }()

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		return err
	}
	fmt.Printf("Daytona sandbox created: %s\n", instance.ID())
	return nil
}

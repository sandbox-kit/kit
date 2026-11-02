package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(fmt.Errorf("copy .env.example to .env and fill in its values: %w", err))
	}
	config := &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.APIKey{
				ApiKey: &sandbox.APIKeyCredentials{
					Key: os.Getenv("DAYTONA_API_KEY"),
				},
			},
		}, Timeout: sandbox.Duration(2 * time.Minute)}
	if value := os.Getenv("DAYTONA_API_URL"); value != "" {
		config.Endpoint = sandbox.Value(value)
	}
	if value := os.Getenv("DAYTONA_TARGET"); value != "" {
		config.Region = sandbox.Value(value)
	}
	client, err := sandbox.NewClient(daytona.New(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.Background())
	instance, err := client.Get(context.Background(), os.Getenv("DAYTONA_SANDBOX_ID"))
	if err != nil {
		log.Fatal(err)
	}
	if !instance.Supports(sandbox.LifecycleCapabilityDelete) {
		log.Fatalf("sandbox %s does not support delete", instance.ID())
	}
	if err := instance.Delete(context.Background(), &sandbox.DeleteOptions{Wait: sandbox.Value(false)}); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Daytona deletion requested: %s (terminal=%t)\n", instance.ID(), instance.Info().GetTerminal())
}

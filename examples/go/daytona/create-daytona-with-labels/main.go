package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
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
		},
		Timeout: sandbox.Duration(2 * time.Minute),
	}
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
	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "ubuntu:22.04"}}},
		Labels: map[string]string{
			"team":    "platform",
			"purpose": "kit-example",
		},
		Environment: map[string]string{
			"NODE_ENV": "development",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Daytona sandbox created with labels and environment: %s\n", instance.ID())
}

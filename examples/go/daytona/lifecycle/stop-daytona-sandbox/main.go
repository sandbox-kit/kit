// Retrieve and stop one existing Daytona sandbox through Sandbox Kit.
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
		log.Fatalf("copy .env.example to .env and fill in its values: %v", err)
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
	if endpoint := os.Getenv("DAYTONA_API_URL"); endpoint != "" {
		config.Endpoint = sandbox.Value(endpoint)
	}
	if target := os.Getenv("DAYTONA_TARGET"); target != "" {
		config.Region = sandbox.Value(target)
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
	if !instance.Supports(sandbox.LifecycleCapabilityStop) {
		log.Fatalf("sandbox %s does not support stop", instance.ID())
	}
	if err := instance.Stop(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Daytona sandbox stopped: %s (%s)\n", instance.ID(), instance.Info().State)
}

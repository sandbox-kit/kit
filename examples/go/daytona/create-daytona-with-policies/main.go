// Create a Daytona sandbox with auto-pause disabled and delayed deletion.
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
		log.Fatalf("copy .env.example to .env and set credentials: %v", err)
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
	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Lifetime: &sandbox.LifetimePolicy{
			IdlePause: &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled},
			StoppedDelete: &sandbox.AutomaticAction{
				Mode:  sandbox.PolicyModeAfter,
				After: sandbox.Duration(10 * time.Minute),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Daytona sandbox created: %s (auto-pause disabled, deletion 10 minutes after stopping)\n", instance.ID())
	info := instance.Info()
	fmt.Printf("id %s; resources reported: %t\n", info.GetId(), info.Resources != nil)
}

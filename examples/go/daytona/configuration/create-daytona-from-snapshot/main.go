// Create one Daytona sandbox from the stock daytona-small snapshot.
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
		Source: &sandbox.SandboxSource{Kind: &sandbox.Snapshot{Snapshot: &sandbox.SnapshotSource{Reference: "daytona-small"}}},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Daytona sandbox created from daytona-small: %s\n", instance.ID())
}

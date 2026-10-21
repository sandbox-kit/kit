// Create a Daytona parent sandbox and one ephemeral sandbox linked to it.
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
	parent, err := client.Create(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	child, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Network:  &sandbox.NetworkConfig{LinkedSandbox: sandbox.Value(parent.ID())},
		Lifetime: &sandbox.LifetimePolicy{Ephemeral: sandbox.Value(true)},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Daytona parent %s; linked ephemeral child %s\n", parent.ID(), child.ID())
}

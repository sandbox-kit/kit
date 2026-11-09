package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/blaxel"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("copy .env.example to .env and fill in its values: %v", err)
	}
	region := os.Getenv("BL_REGION")
	if region == "" {
		region = "us-pdx-1"
	}
	config := &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.APIKey{
				ApiKey: &sandbox.APIKeyCredentials{
					Key: os.Getenv("BL_API_KEY"),
				},
			},
		},
		Scope: &sandbox.Scope{
			Workspace: sandbox.Value(os.Getenv("BL_WORKSPACE")),
		},
		Region:  sandbox.Value(region),
		Timeout: sandbox.Duration(90 * time.Second),
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.Background())
	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Name:   sandbox.Value("kit-create-blaxel-client-timeout"),
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"}}},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Blaxel sandbox created with client timeout: %s\n", instance.ID())
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(fmt.Errorf("copy .env.example to .env and fill in its values: %w", err))
	}
	environment := os.Getenv("MODAL_ENVIRONMENT")
	if environment == "" {
		environment = "main"
	}
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.TokenPair{
				TokenPair: &sandbox.TokenPairCredentials{
					Id:     os.Getenv("MODAL_TOKEN_ID"),
					Secret: os.Getenv("MODAL_TOKEN_SECRET"),
				},
			},
		},
		Scope: &sandbox.Scope{
			AppName:     sandbox.Value(os.Getenv("MODAL_APP_NAME")),
			Environment: sandbox.Value(environment),
		},
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.Background())
	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}}},
		Security: &sandbox.SecurityConfig{
			Secrets: []*sandbox.SecretReference{{
				Reference: os.Getenv("MODAL_SECRET_NAME"),
				Injection: sandbox.SecretInjectionEnvironment,
			}},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Modal sandbox created with secret: %s\n", instance.ID())
}

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
	instance, err := client.Attach(context.Background(), os.Getenv("MODAL_SANDBOX_ID"))
	if err != nil {
		log.Fatal(err)
	}
	if !instance.Supports(sandbox.LifecycleCapabilityTerminate) {
		log.Fatalf("sandbox %s does not support terminate", instance.ID())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	exit, err := instance.Terminate(ctx, &sandbox.TerminateOptions{
		Wait:    sandbox.Value(true),
		Timeout: sandbox.Duration(time.Minute),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Modal sandbox terminated with 60s wait timeout: %s (exit %d)\n", instance.ID(), exit.GetExitCode())
}

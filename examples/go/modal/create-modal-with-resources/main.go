// Create a Modal sandbox with fractional CPU, memory, and a bounded lifetime.
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
		log.Fatalf("copy .env.example to .env and set credentials: %v", err)
	}
	client, err := sandbox.NewClient(clientConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.Background())
	instance, err := client.Create(context.Background(), createOptions())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Modal sandbox created: %s (requested 0.5 CPU cores, 512 MiB, 5-minute lifetime)\n", instance.ID())
	info := instance.Info()
	fmt.Printf("ID origin: %s; allocated resources reported: %t\n", info.Origins["id"], info.Resources != nil)
}

func clientConfig() (sandbox.Provider, *sandbox.Config) {
	environment := os.Getenv("MODAL_ENVIRONMENT")
	if environment == "" {
		environment = "main"
	}
	return modal.New(), &sandbox.Config{
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
	}
}

func createOptions() *sandbox.CreateOptions {
	return &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{
			Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}},
		},
		Resources: &sandbox.Resources{
			CpuCores:  sandbox.Value(0.5),
			MemoryMib: sandbox.Value(uint64(512)),
		},
		Lifetime: &sandbox.LifetimePolicy{MaximumLifetime: sandbox.Duration(5 * time.Minute)},
	}
}
